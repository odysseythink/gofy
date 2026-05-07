// Package sqlvec is a shared IVector base for PG-compatible vector stores
// (pgvector, pgvecto_rs, opengauss, relyt, analyticdb, vastbase …).
//
// Dialect-specific bits (DDL, vector literal format, extension name) are
// supplied via the Dialect interface; the rest of the IVector surface —
// batch INSERT, cosine search, metadata filter, full-text — lives here.
package sqlvec

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	ragentities "github.com/odysseythink/gofy/backend/entities/rag"
)

// Dialect describes the per-backend differences for PG-compatible vector stores.
type Dialect interface {
	// Name is used as the IVector.GetType() value and for error messages.
	Name() string
	// CreateExtensionSQL is executed before CREATE TABLE. Return "" to skip.
	CreateExtensionSQL() string
	// CreateTableSQL returns a full "CREATE TABLE IF NOT EXISTS ..." statement.
	CreateTableSQL(table string, dimension int) string
	// CreateIndexSQL returns a full CREATE INDEX statement for the vector
	// column, or "" to skip. indexHash is a stable short hash for naming.
	CreateIndexSQL(table, indexHash string, dimension int) string
	// FormatVector returns the SQL literal for a float32 vector (e.g. "[1,2,3]").
	// When used with $N placeholder the caller still passes this as a string arg
	// and the SQL uses "$N::vector" cast; FormatVector is what goes in $N.
	FormatVector(v []float32) string
	// VectorCast is the type cast applied to the vector parameter (e.g. "vector",
	// "halfvec", "vectors.vector"). Used as "$N::<cast>".
	VectorCast() string
	// DistanceExpr returns the SQL expression for distance between the stored
	// column and the query parameter ($1). Lower distance = higher similarity.
	// Example: `embedding <=> $1::vector`.
	DistanceExpr() string
}

// Vector is the shared IVector implementation.
type Vector struct {
	dialect    Dialect
	pool       *pgxpool.Pool
	collection string
	table      string
	indexHash  string
	ownsPool   bool
}

var _ ragentities.IVector = (*Vector)(nil)

// Config carries the connection settings shared by all PG-compatible backends.
type Config struct {
	Host          string
	Port          int
	User          string
	Password      string
	Database      string
	MinConnection int32
	MaxConnection int32
	// Extra appended to the DSN query string, e.g. "sslmode=require".
	DSNQuery string
}

func (c Config) DSN() string {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", c.User, c.Password, c.Host, c.Port, c.Database)
	if c.DSNQuery != "" {
		dsn += "?" + c.DSNQuery
	}
	return dsn
}

// Open dials the database and constructs a Vector bound to the given collection.
func Open(ctx context.Context, cfg Config, collection string, d Dialect) (*Vector, error) {
	if collection == "" {
		return nil, errors.New("sqlvec: collection is required")
	}
	poolCfg, err := pgxpool.ParseConfig(cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("%s: parse dsn: %w", d.Name(), err)
	}
	if cfg.MinConnection > 0 {
		poolCfg.MinConns = cfg.MinConnection
	}
	if cfg.MaxConnection > 0 {
		poolCfg.MaxConns = cfg.MaxConnection
	}
	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("%s: open pool: %w", d.Name(), err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("%s: ping: %w", d.Name(), err)
	}
	return newWithPool(pool, collection, d, true), nil
}

// OpenWithPool wraps an existing pool (caller owns lifecycle).
func OpenWithPool(pool *pgxpool.Pool, collection string, d Dialect) *Vector {
	return newWithPool(pool, collection, d, false)
}

func newWithPool(pool *pgxpool.Pool, collection string, d Dialect, owns bool) *Vector {
	table := "embedding_" + sanitizeIdent(collection)
	sum := md5.Sum([]byte(table))
	return &Vector{
		dialect:    d,
		pool:       pool,
		collection: collection,
		table:      table,
		indexHash:  hex.EncodeToString(sum[:])[:8],
		ownsPool:   owns,
	}
}

func (v *Vector) Close() {
	if v.ownsPool {
		v.pool.Close()
	}
}

// Table exposes the backing table name (handy in tests).
func (v *Vector) Table() string { return v.table }

func (v *Vector) GetType() string { return v.dialect.Name() }

func (v *Vector) Create(ctx context.Context, docs []*ragentities.Document, embs [][]float32) error {
	if len(embs) == 0 {
		return fmt.Errorf("%s: Create requires at least one embedding", v.dialect.Name())
	}
	if err := v.createCollection(ctx, len(embs[0])); err != nil {
		return err
	}
	return v.AddTexts(ctx, docs, embs)
}

func (v *Vector) createCollection(ctx context.Context, dim int) error {
	conn, err := v.pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	if ext := v.dialect.CreateExtensionSQL(); ext != "" {
		if _, err := conn.Exec(ctx, ext); err != nil {
			return fmt.Errorf("%s: create extension: %w", v.dialect.Name(), err)
		}
	}
	if _, err := conn.Exec(ctx, v.dialect.CreateTableSQL(quoteIdent(v.table), dim)); err != nil {
		return fmt.Errorf("%s: create table: %w", v.dialect.Name(), err)
	}
	if idxSQL := v.dialect.CreateIndexSQL(quoteIdent(v.table), v.indexHash, dim); idxSQL != "" {
		if _, err := conn.Exec(ctx, idxSQL); err != nil {
			return fmt.Errorf("%s: create index: %w", v.dialect.Name(), err)
		}
	}
	return nil
}

func (v *Vector) AddTexts(ctx context.Context, docs []*ragentities.Document, embs [][]float32) error {
	if len(docs) == 0 {
		return nil
	}
	if len(docs) != len(embs) {
		return fmt.Errorf("%s: docs=%d embeddings=%d mismatch", v.dialect.Name(), len(docs), len(embs))
	}

	cast := v.dialect.VectorCast()
	stmt := fmt.Sprintf(
		`INSERT INTO %s (id, text, meta, embedding) VALUES ($1, $2, $3::jsonb, $4::%s)
		 ON CONFLICT (id) DO UPDATE SET text = EXCLUDED.text, meta = EXCLUDED.meta, embedding = EXCLUDED.embedding`,
		quoteIdent(v.table), cast)

	b := &pgx.Batch{}
	for i, d := range docs {
		id, _ := d.Metadata["doc_id"].(string)
		if id == "" {
			return fmt.Errorf("%s: docs[%d] missing metadata.doc_id", v.dialect.Name(), i)
		}
		metaBytes, err := json.Marshal(d.Metadata)
		if err != nil {
			return fmt.Errorf("%s: marshal metadata[%d]: %w", v.dialect.Name(), i, err)
		}
		b.Queue(stmt, id, d.PageContent, string(metaBytes), v.dialect.FormatVector(embs[i]))
	}
	br := v.pool.SendBatch(ctx, b)
	defer br.Close()
	for i := 0; i < len(docs); i++ {
		if _, err := br.Exec(); err != nil {
			return fmt.Errorf("%s: batch insert row %d: %w", v.dialect.Name(), i, err)
		}
	}
	return nil
}

func (v *Vector) TextExists(ctx context.Context, id string) (bool, error) {
	var got string
	err := v.pool.QueryRow(ctx,
		fmt.Sprintf(`SELECT id FROM %s WHERE id = $1`, quoteIdent(v.table)), id).Scan(&got)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, pgx.ErrNoRows) || isUndefinedTable(err) {
		return false, nil
	}
	return false, err
}

func (v *Vector) GetIDsByMetadataField(ctx context.Context, key, value string) ([]string, error) {
	rows, err := v.pool.Query(ctx,
		fmt.Sprintf(`SELECT id FROM %s WHERE meta->>$1 = $2`, quoteIdent(v.table)), key, value)
	if err != nil {
		if isUndefinedTable(err) {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()
	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func (v *Vector) DeleteByIDs(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	_, err := v.pool.Exec(ctx,
		fmt.Sprintf(`DELETE FROM %s WHERE id = ANY($1)`, quoteIdent(v.table)), ids)
	if err != nil && !isUndefinedTable(err) {
		return err
	}
	return nil
}

func (v *Vector) DeleteByMetadataField(ctx context.Context, key, value string) error {
	_, err := v.pool.Exec(ctx,
		fmt.Sprintf(`DELETE FROM %s WHERE meta->>$1 = $2`, quoteIdent(v.table)), key, value)
	if err != nil && !isUndefinedTable(err) {
		return err
	}
	return nil
}

func (v *Vector) SearchByVector(ctx context.Context, q []float32, opts ragentities.SearchOptions) ([]*ragentities.SearchResult, error) {
	topK := opts.TopK
	if topK <= 0 {
		topK = 4
	}
	where, args := buildSearchWhere(opts, 2)
	sql := fmt.Sprintf(
		`SELECT id, text, meta, %s AS distance FROM %s %s ORDER BY distance LIMIT %d`,
		v.dialect.DistanceExpr(), quoteIdent(v.table), where, topK)
	allArgs := append([]any{v.dialect.FormatVector(q)}, args...)

	rows, err := v.pool.Query(ctx, sql, allArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*ragentities.SearchResult
	for rows.Next() {
		var id, text string
		var metaJSON []byte
		var distance float64
		if err := rows.Scan(&id, &text, &metaJSON, &distance); err != nil {
			return nil, err
		}
		var meta map[string]any
		if len(metaJSON) > 0 {
			_ = json.Unmarshal(metaJSON, &meta)
		}
		score := float32(1.0 - distance)
		if score < opts.ScoreThreshold {
			continue
		}
		out = append(out, &ragentities.SearchResult{
			Document: &ragentities.Document{PageContent: text, Metadata: meta},
			Score:    score,
		})
	}
	return out, rows.Err()
}

func (v *Vector) SearchByFullText(ctx context.Context, query string, opts ragentities.SearchOptions) ([]*ragentities.SearchResult, error) {
	topK := opts.TopK
	if topK <= 0 {
		topK = 5
	}
	where, args := buildSearchWhere(opts, 2)
	textMatch := "to_tsvector(coalesce(text,'')) @@ plainto_tsquery($1)"
	if where == "" {
		where = "WHERE " + textMatch
	} else {
		where = strings.Replace(where, "WHERE", "WHERE "+textMatch+" AND", 1)
	}
	sql := fmt.Sprintf(
		`SELECT id, text, meta, ts_rank(to_tsvector(coalesce(text,'')), plainto_tsquery($1)) AS score
		 FROM %s %s ORDER BY score DESC LIMIT %d`,
		quoteIdent(v.table), where, topK)
	allArgs := append([]any{query}, args...)

	rows, err := v.pool.Query(ctx, sql, allArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*ragentities.SearchResult
	for rows.Next() {
		var id, text string
		var metaJSON []byte
		var score float64
		if err := rows.Scan(&id, &text, &metaJSON, &score); err != nil {
			return nil, err
		}
		var meta map[string]any
		if len(metaJSON) > 0 {
			_ = json.Unmarshal(metaJSON, &meta)
		}
		out = append(out, &ragentities.SearchResult{
			Document: &ragentities.Document{PageContent: text, Metadata: meta},
			Score:    float32(score),
		})
	}
	return out, rows.Err()
}

func (v *Vector) Delete(ctx context.Context) error {
	_, err := v.pool.Exec(ctx, fmt.Sprintf(`DROP TABLE IF EXISTS %s`, quoteIdent(v.table)))
	return err
}

// ---- helpers ----

func buildSearchWhere(opts ragentities.SearchOptions, startIdx int) (string, []any) {
	var parts []string
	var args []any
	nextArg := startIdx
	for k, val := range opts.Filter {
		parts = append(parts, fmt.Sprintf("meta->>$%d = $%d", nextArg, nextArg+1))
		args = append(args, k, fmt.Sprint(val))
		nextArg += 2
	}
	if len(opts.DocumentIDs) > 0 {
		parts = append(parts, fmt.Sprintf("meta->>'document_id' = ANY($%d)", nextArg))
		args = append(args, opts.DocumentIDs)
		nextArg++
	}
	if len(parts) == 0 {
		return "", nil
	}
	return "WHERE " + strings.Join(parts, " AND "), args
}

// FormatPGVector renders a float32 slice as the "[1,2,3]" literal that
// pgvector (and pgvector-compatible forks) accept when cast to ::vector.
func FormatPGVector(vec []float32) string {
	var b strings.Builder
	b.Grow(len(vec) * 8)
	b.WriteByte('[')
	for i, x := range vec {
		if i > 0 {
			b.WriteByte(',')
		}
		fmt.Fprintf(&b, "%g", x)
	}
	b.WriteByte(']')
	return b.String()
}

func quoteIdent(s string) string {
	return `"` + strings.ReplaceAll(s, `"`, `""`) + `"`
}

func sanitizeIdent(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range strings.ToLower(s) {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9', r == '_':
			b.WriteRune(r)
		default:
			b.WriteByte('_')
		}
	}
	return b.String()
}

func isUndefinedTable(err error) bool {
	msg := err.Error()
	return strings.Contains(msg, "SQLSTATE 42P01") ||
		(strings.Contains(msg, "does not exist") && strings.Contains(msg, "relation"))
}
