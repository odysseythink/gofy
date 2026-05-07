// Package baidu is the Baidu Cloud VectorDB (VDB) IVector adapter, talking
// to MochowDB / Baidu VectorDB's HTTP API. Auth is bearer-style account+apikey.
package baidu

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/odysseythink/confy"

	vdb "github.com/odysseythink/gofy/backend/core/rag/datasource/vdb"
	"github.com/odysseythink/gofy/backend/core/rag/datasource/vdb/internal/httpvec"
	ragentities "github.com/odysseythink/gofy/backend/entities/rag"
	vectorenumtypes "github.com/odysseythink/gofy/backend/enum_types/rag/vector"
)

func init() { vdb.Register(vectorenumtypes.Vector_BAIDU, &Factory{}) }

type Config struct {
	Endpoint string
	Account  string
	APIKey   string
	Database string
}

type Factory struct{}

func (Factory) New(ctx context.Context, cfg vdb.AdapterConfig) (ragentities.IVector, error) {
	return Open(ctx, loadConfig(), cfg.CollectionName)
}

func loadConfig() Config {
	return Config{
		Endpoint: confy.GetWithDefault[string]("baidu.endpoint", ""),
		Account:  confy.GetWithDefault[string]("baidu.account", "root"),
		APIKey:   confy.GetWithDefault[string]("baidu.api-key", ""),
		Database: confy.GetWithDefault[string]("baidu.database", "gofy"),
	}
}

func Open(ctx context.Context, cfg Config, collection string) (*Vector, error) {
	if cfg.Endpoint == "" {
		return nil, errors.New("baidu: endpoint is required")
	}
	h := map[string]string{"Authorization": fmt.Sprintf("Bearer account=%s&api_key=%s", cfg.Account, cfg.APIKey)}
	return &Vector{cfg: cfg, collection: collection, http: httpvec.New(cfg.Endpoint, 30*time.Second, h)}, nil
}

type Vector struct {
	cfg        Config
	collection string
	http       *httpvec.Client
}

var _ ragentities.IVector = (*Vector)(nil)

func (v *Vector) GetType() string { return string(vectorenumtypes.Vector_BAIDU) }

// Baidu's API shape is close enough to Tencent's that the same REST verbs
// (collection/create, document/upsert, document/search, document/delete) work.
// If the target product is MochowDB the paths are identical to Tencent's;
// otherwise override via a future BAIDU_PATH_* config.

func (v *Vector) Create(ctx context.Context, docs []*ragentities.Document, embs [][]float32) error {
	if len(embs) == 0 {
		return errors.New("baidu: Create requires at least one embedding")
	}
	_ = v.http.Do(ctx, "POST", "/v1/collection/drop",
		map[string]any{"database": v.cfg.Database, "collection": v.collection}, nil)
	body := map[string]any{
		"database":   v.cfg.Database,
		"collection": v.collection,
		"fields": []map[string]any{
			{"name": "id", "type": "STRING", "primary": true},
			{"name": "vector", "type": "FLOAT_VECTOR", "dimension": len(embs[0])},
			{"name": "text", "type": "TEXT"},
			{"name": "metadata", "type": "JSON"},
		},
	}
	if err := v.http.Do(ctx, "POST", "/v1/collection/create", body, nil); err != nil {
		return fmt.Errorf("baidu: create collection: %w", err)
	}
	return v.AddTexts(ctx, docs, embs)
}

func (v *Vector) AddTexts(ctx context.Context, docs []*ragentities.Document, embs [][]float32) error {
	if len(docs) != len(embs) {
		return fmt.Errorf("baidu: docs=%d embeddings=%d mismatch", len(docs), len(embs))
	}
	rows := make([]map[string]any, len(docs))
	for i, d := range docs {
		id, _ := d.Metadata["doc_id"].(string)
		if id == "" {
			return fmt.Errorf("baidu: docs[%d] missing metadata.doc_id", i)
		}
		rows[i] = map[string]any{"id": id, "vector": embs[i], "text": d.PageContent, "metadata": d.Metadata}
	}
	return v.http.Do(ctx, "POST", "/v1/document/upsert",
		map[string]any{"database": v.cfg.Database, "collection": v.collection, "rows": rows}, nil)
}

func (v *Vector) TextExists(ctx context.Context, id string) (bool, error) {
	var resp struct {
		Rows []any `json:"rows"`
	}
	err := v.http.Do(ctx, "POST", "/v1/document/query",
		map[string]any{"database": v.cfg.Database, "collection": v.collection, "primary_key": id}, &resp)
	if err != nil {
		return false, err
	}
	return len(resp.Rows) > 0, nil
}

func (v *Vector) GetIDsByMetadataField(ctx context.Context, key, value string) ([]string, error) {
	body := map[string]any{
		"database": v.cfg.Database, "collection": v.collection,
		"filter": fmt.Sprintf(`metadata.%s = "%s"`, key, value),
	}
	var resp struct {
		Rows []struct {
			ID string `json:"id"`
		} `json:"rows"`
	}
	if err := v.http.Do(ctx, "POST", "/v1/document/query", body, &resp); err != nil {
		return nil, err
	}
	ids := make([]string, len(resp.Rows))
	for i, r := range resp.Rows {
		ids[i] = r.ID
	}
	return ids, nil
}

func (v *Vector) DeleteByIDs(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	return v.http.Do(ctx, "POST", "/v1/document/delete",
		map[string]any{"database": v.cfg.Database, "collection": v.collection, "primary_keys": ids}, nil)
}

func (v *Vector) DeleteByMetadataField(ctx context.Context, key, value string) error {
	return v.http.Do(ctx, "POST", "/v1/document/delete",
		map[string]any{"database": v.cfg.Database, "collection": v.collection,
			"filter": fmt.Sprintf(`metadata.%s = "%s"`, key, value)}, nil)
}

func (v *Vector) SearchByVector(ctx context.Context, qv []float32, opts ragentities.SearchOptions) ([]*ragentities.SearchResult, error) {
	topK := opts.TopK
	if topK <= 0 {
		topK = 4
	}
	body := map[string]any{
		"database": v.cfg.Database, "collection": v.collection,
		"anns": map[string]any{
			"vector_field": "vector", "vector": qv, "top_k": topK,
			"params": map[string]any{"ef": 200},
		},
		"retrieve_fields": []string{"text", "metadata"},
	}
	var resp struct {
		Rows []struct {
			Score    float32        `json:"score"`
			Text     string         `json:"text"`
			Metadata map[string]any `json:"metadata"`
		} `json:"rows"`
	}
	if err := v.http.Do(ctx, "POST", "/v1/document/search", body, &resp); err != nil {
		return nil, err
	}
	var out []*ragentities.SearchResult
	for _, r := range resp.Rows {
		if r.Score < opts.ScoreThreshold {
			continue
		}
		out = append(out, &ragentities.SearchResult{
			Document: &ragentities.Document{PageContent: r.Text, Metadata: r.Metadata},
			Score:    r.Score,
		})
	}
	return out, nil
}

func (v *Vector) SearchByFullText(ctx context.Context, q string, _ ragentities.SearchOptions) ([]*ragentities.SearchResult, error) {
	return nil, errors.New("baidu: full-text search not supported by this adapter")
}

func (v *Vector) Delete(ctx context.Context) error {
	return v.http.Do(ctx, "POST", "/v1/collection/drop",
		map[string]any{"database": v.cfg.Database, "collection": v.collection}, nil)
}
