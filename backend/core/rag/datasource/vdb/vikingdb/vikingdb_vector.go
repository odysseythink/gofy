// Package vikingdb is the ByteDance VikingDB IVector adapter.
//
// VikingDB uses AWS-V4-style request signing with volcengine's AK/SK. A full
// production impl should vendor github.com/volcengine/volcengine-go-sdk; this
// adapter uses a simplified bearer-token path suitable for private deployments
// and documents the TODO for cloud-SaaS use.
package vikingdb

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/odysseythink/confy"

	vdb "mlib.com/gofy/server/core/rag/datasource/vdb"
	"mlib.com/gofy/server/core/rag/datasource/vdb/internal/httpvec"
	ragentities "mlib.com/gofy/server/entities/rag"
	vectorenumtypes "mlib.com/gofy/server/enum_types/rag/vector"
)

func init() { vdb.Register(vectorenumtypes.Vector_VIKINGDB, &Factory{}) }

type Config struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Region    string
	Project   string
}

type Factory struct{}

func (Factory) New(ctx context.Context, cfg vdb.AdapterConfig) (ragentities.IVector, error) {
	return Open(ctx, loadConfig(), cfg.CollectionName)
}

func loadConfig() Config {
	return Config{
		Endpoint:  confy.GetWithDefault[string]("vikingdb.endpoint", ""),
		AccessKey: confy.GetWithDefault[string]("vikingdb.access-key", ""),
		SecretKey: confy.GetWithDefault[string]("vikingdb.secret-key", ""),
		Region:    confy.GetWithDefault[string]("vikingdb.region", "cn-beijing"),
		Project:   confy.GetWithDefault[string]("vikingdb.project", "default"),
	}
}

// Open constructs a VikingDB vector. NOTE: this adapter talks to VikingDB's
// inner HTTP API via bearer auth only; for production SaaS you need V4 signing.
func Open(ctx context.Context, cfg Config, collection string) (*Vector, error) {
	if cfg.Endpoint == "" {
		return nil, errors.New("vikingdb: endpoint is required")
	}
	h := map[string]string{"Authorization": "Bearer " + cfg.AccessKey + ":" + cfg.SecretKey}
	return &Vector{cfg: cfg, collection: collection, http: httpvec.New(cfg.Endpoint, 30*time.Second, h)}, nil
}

type Vector struct {
	cfg        Config
	collection string
	http       *httpvec.Client
}

var _ ragentities.IVector = (*Vector)(nil)

func (v *Vector) GetType() string { return string(vectorenumtypes.Vector_VIKINGDB) }

func (v *Vector) Create(ctx context.Context, docs []*ragentities.Document, embs [][]float32) error {
	if len(embs) == 0 {
		return errors.New("vikingdb: Create requires at least one embedding")
	}
	_ = v.http.Do(ctx, "POST", "/api/collection/delete",
		map[string]any{"collection_name": v.collection}, nil)
	body := map[string]any{
		"collection_name": v.collection,
		"fields": []map[string]any{
			{"field_name": "id", "field_type": "string", "primary_key": true},
			{"field_name": "vector", "field_type": "vector", "dim": len(embs[0])},
			{"field_name": "text", "field_type": "text"},
			{"field_name": "metadata", "field_type": "json"},
		},
	}
	if err := v.http.Do(ctx, "POST", "/api/collection/create", body, nil); err != nil {
		return fmt.Errorf("vikingdb: create collection: %w", err)
	}
	return v.AddTexts(ctx, docs, embs)
}

func (v *Vector) AddTexts(ctx context.Context, docs []*ragentities.Document, embs [][]float32) error {
	if len(docs) != len(embs) {
		return fmt.Errorf("vikingdb: docs=%d embeddings=%d mismatch", len(docs), len(embs))
	}
	rows := make([]map[string]any, len(docs))
	for i, d := range docs {
		id, _ := d.Metadata["doc_id"].(string)
		if id == "" {
			return fmt.Errorf("vikingdb: docs[%d] missing metadata.doc_id", i)
		}
		rows[i] = map[string]any{"fields": map[string]any{
			"id": id, "vector": embs[i], "text": d.PageContent, "metadata": d.Metadata,
		}}
	}
	return v.http.Do(ctx, "POST", "/api/data/upsert",
		map[string]any{"collection_name": v.collection, "data": rows}, nil)
}

func (v *Vector) TextExists(ctx context.Context, id string) (bool, error) {
	var resp struct {
		Data []any `json:"data"`
	}
	err := v.http.Do(ctx, "POST", "/api/data/fetch",
		map[string]any{"collection_name": v.collection, "primary_keys": []string{id}}, &resp)
	if err != nil {
		return false, err
	}
	return len(resp.Data) > 0, nil
}

func (v *Vector) GetIDsByMetadataField(ctx context.Context, key, value string) ([]string, error) {
	body := map[string]any{"collection_name": v.collection,
		"filter": fmt.Sprintf(`metadata.%s = "%s"`, key, value)}
	var resp struct {
		Data []struct {
			Fields struct {
				ID string `json:"id"`
			} `json:"fields"`
		} `json:"data"`
	}
	if err := v.http.Do(ctx, "POST", "/api/data/query", body, &resp); err != nil {
		return nil, err
	}
	ids := make([]string, len(resp.Data))
	for i, r := range resp.Data {
		ids[i] = r.Fields.ID
	}
	return ids, nil
}

func (v *Vector) DeleteByIDs(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	return v.http.Do(ctx, "POST", "/api/data/delete",
		map[string]any{"collection_name": v.collection, "primary_keys": ids}, nil)
}

func (v *Vector) DeleteByMetadataField(ctx context.Context, key, value string) error {
	return v.http.Do(ctx, "POST", "/api/data/delete",
		map[string]any{"collection_name": v.collection,
			"filter": fmt.Sprintf(`metadata.%s = "%s"`, key, value)}, nil)
}

func (v *Vector) SearchByVector(ctx context.Context, qv []float32, opts ragentities.SearchOptions) ([]*ragentities.SearchResult, error) {
	topK := opts.TopK
	if topK <= 0 {
		topK = 4
	}
	body := map[string]any{
		"collection_name": v.collection,
		"vector":          qv,
		"top_k":           topK,
		"output_fields":   []string{"text", "metadata"},
	}
	var resp struct {
		Data []struct {
			Score  float32 `json:"score"`
			Fields struct {
				Text     string         `json:"text"`
				Metadata map[string]any `json:"metadata"`
			} `json:"fields"`
		} `json:"data"`
	}
	if err := v.http.Do(ctx, "POST", "/api/data/search", body, &resp); err != nil {
		return nil, err
	}
	var out []*ragentities.SearchResult
	for _, r := range resp.Data {
		if r.Score < opts.ScoreThreshold {
			continue
		}
		out = append(out, &ragentities.SearchResult{
			Document: &ragentities.Document{PageContent: r.Fields.Text, Metadata: r.Fields.Metadata},
			Score:    r.Score,
		})
	}
	return out, nil
}

func (v *Vector) SearchByFullText(ctx context.Context, q string, _ ragentities.SearchOptions) ([]*ragentities.SearchResult, error) {
	return nil, errors.New("vikingdb: full-text search not supported by this adapter")
}

func (v *Vector) Delete(ctx context.Context) error {
	return v.http.Do(ctx, "POST", "/api/collection/delete",
		map[string]any{"collection_name": v.collection}, nil)
}
