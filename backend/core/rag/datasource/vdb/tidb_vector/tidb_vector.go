// Package tidb_vector is the TiDB-native vector IVector adapter.
//
// TiDB >= 7.5 supports native VECTOR column types, distance functions
// (VEC_COSINE_DISTANCE), and TiFlash HNSW indexes. This adapter talks
// MySQL wire protocol.
package tidb_vector

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/odysseythink/confy"

	vdb "mlib.com/gofy/server/core/rag/datasource/vdb"
	ragentities "mlib.com/gofy/server/entities/rag"
	vectorenumtypes "mlib.com/gofy/server/enum_types/rag/vector"
)

func init() { vdb.Register(vectorenumtypes.Vector_TIDB_VECTOR, &Factory{}) }

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	TLS      string // e.g. "true", "false", "skip-verify"
}

func (c Config) DSN() string {
	// MySQL DSN: user:pass@tcp(host:port)/db?tls=...&parseTime=true
	tls := c.TLS
	if tls == "" {
		tls = "false"
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?tls=%s&parseTime=true",
		c.User, c.Password, c.Host, c.Port, c.Database, tls)
}

type Factory struct{}

func (Factory) New(ctx context.Context, cfg vdb.AdapterConfig) (ragentities.IVector, error) {
	return Open(ctx, loadConfig(), cfg.CollectionName)
}

func loadConfig() Config {
	return Config{
		Host:     confy.GetWithDefault[string]("tidb-vector.host", "localhost"),
		Port:     confy.GetWithDefault[int]("tidb-vector.port", 4000),
		User:     confy.GetWithDefault[string]("tidb-vector.user", "root"),
		Password: confy.GetWithDefault[string]("tidb-vector.password", ""),
		Database: confy.GetWithDefault[string]("tidb-vector.database", "test"),
		TLS:      confy.GetWithDefault[string]("tidb-vector.tls", "false"),
	}
}

func Open(ctx context.Context, cfg Config, collection string) (*Vector, error) {
	if collection == "" {
		return nil, errors.New("tidb_vector: collection is required")
	}
	db, err := sql.Open("mysql", cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("tidb_vector: open: %w", err)
	}
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("tidb_vector: ping: %w", err)
	}
	return &Vector{db: db, table: "embedding_" + sanitize(collection)}, nil
}

type Vector struct {
	db    *sql.DB
	table string
}

var _ ragentities.IVector = (*Vector)(nil)

func (v *Vector) Close() error { return v.db.Close() }

func (v *Vector) GetType() string { return string(vectorenumtypes.Vector_TIDB_VECTOR) }

func (v *Vector) Create(ctx context.Context, docs []*ragentities.Document, embs [][]float32) error {
	if len(embs) == 0 {
		return errors.New("tidb_vector: Create requires at least one embedding")
	}
	ddl := fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s` ("+
		"id VARCHAR(64) NOT NULL PRIMARY KEY, "+
		"text LONGTEXT NOT NULL, "+
		"meta JSON NOT NULL, "+
		"embedding VECTOR(%d) NOT NULL"+
		") ENGINE=InnoDB", v.table, len(embs[0]))
	if _, err := v.db.ExecContext(ctx, ddl); err != nil {
		return fmt.Errorf("tidb_vector: create table: %w", err)
	}
	return v.AddTexts(ctx, docs, embs)
}

func (v *Vector) AddTexts(ctx context.Context, docs []*ragentities.Document, embs [][]float32) error {
	if len(docs) == 0 {
		return nil
	}
	if len(docs) != len(embs) {
		return fmt.Errorf("tidb_vector: docs=%d embeddings=%d mismatch", len(docs), len(embs))
	}
	stmt := fmt.Sprintf("INSERT INTO `%s` (id, text, meta, embedding) VALUES (?, ?, ?, ?) "+
		"ON DUPLICATE KEY UPDATE text=VALUES(text), meta=VALUES(meta), embedding=VALUES(embedding)",
		v.table)
	for i, d := range docs {
		id, _ := d.Metadata["doc_id"].(string)
		if id == "" {
			return fmt.Errorf("tidb_vector: docs[%d] missing metadata.doc_id", i)
		}
		meta, err := json.Marshal(d.Metadata)
		if err != nil {
			return fmt.Errorf("tidb_vector: marshal meta[%d]: %w", i, err)
		}
		if _, err := v.db.ExecContext(ctx, stmt, id, d.PageContent, string(meta), vectorLiteral(embs[i])); err != nil {
			return fmt.Errorf("tidb_vector: insert[%d]: %w", i, err)
		}
	}
	return nil
}

func (v *Vector) TextExists(ctx context.Context, id string) (bool, error) {
	var got string
	err := v.db.QueryRowContext(ctx,
		fmt.Sprintf("SELECT id FROM `%s` WHERE id = ?", v.table), id).Scan(&got)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, sql.ErrNoRows) || isMissingTable(err) {
		return false, nil
	}
	return false, err
}

func (v *Vector) GetIDsByMetadataField(ctx context.Context, key, value string) ([]string, error) {
	q := fmt.Sprintf("SELECT id FROM `%s` WHERE JSON_UNQUOTE(JSON_EXTRACT(meta, CONCAT('$.', ?))) = ?", v.table)
	rows, err := v.db.QueryContext(ctx, q, key, value)
	if err != nil {
		if isMissingTable(err) {
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
	placeholders := strings.Repeat("?,", len(ids))
	placeholders = placeholders[:len(placeholders)-1]
	args := make([]any, len(ids))
	for i, id := range ids {
		args[i] = id
	}
	_, err := v.db.ExecContext(ctx, fmt.Sprintf("DELETE FROM `%s` WHERE id IN (%s)", v.table, placeholders), args...)
	if err != nil && !isMissingTable(err) {
		return err
	}
	return nil
}

func (v *Vector) DeleteByMetadataField(ctx context.Context, key, value string) error {
	q := fmt.Sprintf("DELETE FROM `%s` WHERE JSON_UNQUOTE(JSON_EXTRACT(meta, CONCAT('$.', ?))) = ?", v.table)
	_, err := v.db.ExecContext(ctx, q, key, value)
	if err != nil && !isMissingTable(err) {
		return err
	}
	return nil
}

func (v *Vector) SearchByVector(ctx context.Context, qv []float32, opts ragentities.SearchOptions) ([]*ragentities.SearchResult, error) {
	topK := opts.TopK
	if topK <= 0 {
		topK = 4
	}
	where, args := buildWhere(opts)
	q := fmt.Sprintf(
		"SELECT id, text, meta, VEC_COSINE_DISTANCE(embedding, ?) AS distance FROM `%s` %s ORDER BY distance LIMIT %d",
		v.table, where, topK)
	args = append([]any{vectorLiteral(qv)}, args...)

	rows, err := v.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*ragentities.SearchResult
	for rows.Next() {
		var id, text, metaStr string
		var dist float64
		if err := rows.Scan(&id, &text, &metaStr, &dist); err != nil {
			return nil, err
		}
		var meta map[string]any
		_ = json.Unmarshal([]byte(metaStr), &meta)
		score := float32(1.0 - dist)
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

// SearchByFullText is not implemented (TiDB full-text is work-in-progress per
// upstream). Returning an error lets callers fall back to vector search.
func (v *Vector) SearchByFullText(ctx context.Context, query string, opts ragentities.SearchOptions) ([]*ragentities.SearchResult, error) {
	return nil, errors.New("tidb_vector: full-text search not supported")
}

func (v *Vector) Delete(ctx context.Context) error {
	_, err := v.db.ExecContext(ctx, fmt.Sprintf("DROP TABLE IF EXISTS `%s`", v.table))
	return err
}

// ---- helpers ----

func buildWhere(opts ragentities.SearchOptions) (string, []any) {
	var parts []string
	var args []any
	for k, val := range opts.Filter {
		parts = append(parts, "JSON_UNQUOTE(JSON_EXTRACT(meta, CONCAT('$.', ?))) = ?")
		args = append(args, k, fmt.Sprint(val))
	}
	if len(opts.DocumentIDs) > 0 {
		placeholders := strings.Repeat("?,", len(opts.DocumentIDs))
		placeholders = placeholders[:len(placeholders)-1]
		parts = append(parts, fmt.Sprintf("JSON_UNQUOTE(JSON_EXTRACT(meta, '$.document_id')) IN (%s)", placeholders))
		for _, id := range opts.DocumentIDs {
			args = append(args, id)
		}
	}
	if len(parts) == 0 {
		return "", nil
	}
	return "WHERE " + strings.Join(parts, " AND "), args
}

func vectorLiteral(v []float32) string {
	var b strings.Builder
	b.Grow(len(v) * 8)
	b.WriteByte('[')
	for i, x := range v {
		if i > 0 {
			b.WriteByte(',')
		}
		fmt.Fprintf(&b, "%g", x)
	}
	b.WriteByte(']')
	return b.String()
}

func sanitize(s string) string {
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

func isMissingTable(err error) bool {
	msg := err.Error()
	return strings.Contains(msg, "1146") || strings.Contains(msg, "doesn't exist")
}
