// Package myscale is the MyScale IVector adapter. MyScale is a ClickHouse
// fork with native vector search; this adapter talks ClickHouse HTTP.
package myscale

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/odysseythink/confy"

	vdb "mlib.com/gofy/server/core/rag/datasource/vdb"
	ragentities "mlib.com/gofy/server/entities/rag"
	vectorenumtypes "mlib.com/gofy/server/enum_types/rag/vector"
)

func init() { vdb.Register(vectorenumtypes.Vector_MYSCALE, &Factory{}) }

type Config struct {
	Endpoint string // http://host:8123
	Username string
	Password string
	Database string
}

type Factory struct{}

func (Factory) New(ctx context.Context, cfg vdb.AdapterConfig) (ragentities.IVector, error) {
	return Open(ctx, loadConfig(), cfg.CollectionName)
}

func loadConfig() Config {
	return Config{
		Endpoint: confy.GetWithDefault[string]("myscale.endpoint", "http://localhost:8123"),
		Username: confy.GetWithDefault[string]("myscale.username", "default"),
		Password: confy.GetWithDefault[string]("myscale.password", ""),
		Database: confy.GetWithDefault[string]("myscale.database", "default"),
	}
}

func Open(ctx context.Context, cfg Config, collection string) (*Vector, error) {
	if collection == "" {
		return nil, errors.New("myscale: collection is required")
	}
	return &Vector{cfg: cfg, table: sanitize(collection),
		client: &http.Client{Timeout: 30 * time.Second}}, nil
}

type Vector struct {
	cfg    Config
	table  string
	client *http.Client
}

var _ ragentities.IVector = (*Vector)(nil)

func (v *Vector) GetType() string { return string(vectorenumtypes.Vector_MYSCALE) }

func (v *Vector) fullTable() string { return v.cfg.Database + "." + v.table }

func (v *Vector) exec(ctx context.Context, sql string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", v.cfg.Endpoint, strings.NewReader(sql))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "text/plain")
	if v.cfg.Username != "" {
		req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(v.cfg.Username+":"+v.cfg.Password)))
	}
	resp, err := v.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	buf, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("myscale: %d: %s", resp.StatusCode, string(buf))
	}
	return buf, nil
}

func (v *Vector) Create(ctx context.Context, docs []*ragentities.Document, embs [][]float32) error {
	if len(embs) == 0 {
		return errors.New("myscale: Create requires at least one embedding")
	}
	if _, err := v.exec(ctx, fmt.Sprintf("DROP TABLE IF EXISTS %s", v.fullTable())); err != nil {
		return err
	}
	ddl := fmt.Sprintf(`CREATE TABLE %s (
		id String, text String, metadata String, vector Array(Float32),
		CONSTRAINT cons_vec_len CHECK length(vector) = %d,
		VECTOR INDEX vidx vector TYPE HNSWFLAT('metric_type=Cosine')
	) ENGINE = MergeTree ORDER BY id`, v.fullTable(), len(embs[0]))
	if _, err := v.exec(ctx, ddl); err != nil {
		return fmt.Errorf("myscale: create table: %w", err)
	}
	return v.AddTexts(ctx, docs, embs)
}

func (v *Vector) AddTexts(ctx context.Context, docs []*ragentities.Document, embs [][]float32) error {
	if len(docs) == 0 {
		return nil
	}
	if len(docs) != len(embs) {
		return fmt.Errorf("myscale: docs=%d embeddings=%d mismatch", len(docs), len(embs))
	}
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "INSERT INTO %s (id, text, metadata, vector) FORMAT JSONEachRow\n", v.fullTable())
	enc := json.NewEncoder(&buf)
	for i, d := range docs {
		id, _ := d.Metadata["doc_id"].(string)
		if id == "" {
			return fmt.Errorf("myscale: docs[%d] missing metadata.doc_id", i)
		}
		metaBytes, _ := json.Marshal(d.Metadata)
		_ = enc.Encode(map[string]any{
			"id": id, "text": d.PageContent, "metadata": string(metaBytes), "vector": embs[i],
		})
	}
	_, err := v.exec(ctx, buf.String())
	return err
}

func (v *Vector) TextExists(ctx context.Context, id string) (bool, error) {
	sql := fmt.Sprintf("SELECT id FROM %s WHERE id = '%s' LIMIT 1 FORMAT TabSeparated", v.fullTable(), escape(id))
	body, err := v.exec(ctx, sql)
	if err != nil {
		if strings.Contains(err.Error(), "Table") && strings.Contains(err.Error(), "doesn't exist") {
			return false, nil
		}
		return false, err
	}
	return len(bytes.TrimSpace(body)) > 0, nil
}

func (v *Vector) GetIDsByMetadataField(ctx context.Context, key, value string) ([]string, error) {
	sql := fmt.Sprintf("SELECT id FROM %s WHERE JSONExtractString(metadata, '%s') = '%s' FORMAT TabSeparated",
		v.fullTable(), escape(key), escape(value))
	body, err := v.exec(ctx, sql)
	if err != nil {
		return nil, err
	}
	var ids []string
	for _, line := range strings.Split(strings.TrimSpace(string(body)), "\n") {
		if line != "" {
			ids = append(ids, line)
		}
	}
	return ids, nil
}

func (v *Vector) DeleteByIDs(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	quoted := make([]string, len(ids))
	for i, id := range ids {
		quoted[i] = "'" + escape(id) + "'"
	}
	sql := fmt.Sprintf("ALTER TABLE %s DELETE WHERE id IN (%s)", v.fullTable(), strings.Join(quoted, ","))
	_, err := v.exec(ctx, sql)
	return err
}

func (v *Vector) DeleteByMetadataField(ctx context.Context, key, value string) error {
	sql := fmt.Sprintf("ALTER TABLE %s DELETE WHERE JSONExtractString(metadata, '%s') = '%s'",
		v.fullTable(), escape(key), escape(value))
	_, err := v.exec(ctx, sql)
	return err
}

func (v *Vector) SearchByVector(ctx context.Context, qv []float32, opts ragentities.SearchOptions) ([]*ragentities.SearchResult, error) {
	topK := opts.TopK
	if topK <= 0 {
		topK = 4
	}
	var parts []string
	for k, val := range opts.Filter {
		parts = append(parts, fmt.Sprintf("JSONExtractString(metadata, '%s') = '%s'", escape(k), escape(fmt.Sprint(val))))
	}
	if len(opts.DocumentIDs) > 0 {
		q := make([]string, len(opts.DocumentIDs))
		for i, id := range opts.DocumentIDs {
			q[i] = "'" + escape(id) + "'"
		}
		parts = append(parts, fmt.Sprintf("JSONExtractString(metadata, 'document_id') IN (%s)", strings.Join(q, ",")))
	}
	where := ""
	if len(parts) > 0 {
		where = " WHERE " + strings.Join(parts, " AND ")
	}
	vec := formatArray(qv)
	sql := fmt.Sprintf(`SELECT id, text, metadata, distance(vector, %s) AS d FROM %s%s ORDER BY d ASC LIMIT %d FORMAT JSONEachRow`,
		vec, v.fullTable(), where, topK)
	body, err := v.exec(ctx, sql)
	if err != nil {
		return nil, err
	}
	var out []*ragentities.SearchResult
	for _, line := range bytes.Split(body, []byte("\n")) {
		if len(bytes.TrimSpace(line)) == 0 {
			continue
		}
		var row struct {
			Text     string  `json:"text"`
			Metadata string  `json:"metadata"`
			Distance float64 `json:"d"`
		}
		if err := json.Unmarshal(line, &row); err != nil {
			continue
		}
		score := float32(1.0 - row.Distance)
		if score < opts.ScoreThreshold {
			continue
		}
		var meta map[string]any
		_ = json.Unmarshal([]byte(row.Metadata), &meta)
		out = append(out, &ragentities.SearchResult{
			Document: &ragentities.Document{PageContent: row.Text, Metadata: meta},
			Score:    score,
		})
	}
	return out, nil
}

func (v *Vector) SearchByFullText(ctx context.Context, q string, opts ragentities.SearchOptions) ([]*ragentities.SearchResult, error) {
	topK := opts.TopK
	if topK <= 0 {
		topK = 5
	}
	sql := fmt.Sprintf(`SELECT id, text, metadata FROM %s WHERE positionCaseInsensitive(text, '%s') > 0 LIMIT %d FORMAT JSONEachRow`,
		v.fullTable(), escape(q), topK)
	body, err := v.exec(ctx, sql)
	if err != nil {
		return nil, err
	}
	var out []*ragentities.SearchResult
	for _, line := range bytes.Split(body, []byte("\n")) {
		if len(bytes.TrimSpace(line)) == 0 {
			continue
		}
		var row struct {
			Text     string `json:"text"`
			Metadata string `json:"metadata"`
		}
		if err := json.Unmarshal(line, &row); err != nil {
			continue
		}
		var meta map[string]any
		_ = json.Unmarshal([]byte(row.Metadata), &meta)
		out = append(out, &ragentities.SearchResult{
			Document: &ragentities.Document{PageContent: row.Text, Metadata: meta},
			Score:    1.0, // myscale full-text has no score; flat 1.0
		})
	}
	return out, nil
}

func (v *Vector) Delete(ctx context.Context) error {
	_, err := v.exec(ctx, fmt.Sprintf("DROP TABLE IF EXISTS %s", v.fullTable()))
	return err
}

func formatArray(v []float32) string {
	var b strings.Builder
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

func escape(s string) string { return strings.ReplaceAll(s, "'", `\'`) }

func sanitize(s string) string {
	var b strings.Builder
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
