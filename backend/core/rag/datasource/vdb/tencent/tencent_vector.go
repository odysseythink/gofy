// Package tencent is the Tencent Cloud VectorDB IVector adapter.
//
// Uses Tencent VectorDB's HTTP API. For production workloads, vendor
// github.com/tencent/vectordatabase-sdk-go — the REST path below handles
// simple bearer-token auth and covers the IVector surface.
package tencent

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

func init() { vdb.Register(vectorenumtypes.Vector_TENCENT, &Factory{}) }

type Config struct {
	Endpoint string
	Username string
	APIKey   string
	Database string
}

type Factory struct{}

func (Factory) New(ctx context.Context, cfg vdb.AdapterConfig) (ragentities.IVector, error) {
	return Open(ctx, loadConfig(), cfg.CollectionName)
}

func loadConfig() Config {
	return Config{
		Endpoint: confy.GetWithDefault[string]("tencent.endpoint", ""),
		Username: confy.GetWithDefault[string]("tencent.username", "root"),
		APIKey:   confy.GetWithDefault[string]("tencent.api-key", ""),
		Database: confy.GetWithDefault[string]("tencent.database", "gofy"),
	}
}

func Open(ctx context.Context, cfg Config, collection string) (*Vector, error) {
	if cfg.Endpoint == "" {
		return nil, errors.New("tencent: endpoint is required")
	}
	h := map[string]string{"Authorization": fmt.Sprintf("Bearer account=%s&api_key=%s", cfg.Username, cfg.APIKey)}
	return &Vector{cfg: cfg, collection: collection, http: httpvec.New(cfg.Endpoint, 30*time.Second, h)}, nil
}

type Vector struct {
	cfg        Config
	collection string
	http       *httpvec.Client
}

var _ ragentities.IVector = (*Vector)(nil)

func (v *Vector) GetType() string { return string(vectorenumtypes.Vector_TENCENT) }

func (v *Vector) Create(ctx context.Context, docs []*ragentities.Document, embs [][]float32) error {
	if len(embs) == 0 {
		return errors.New("tencent: Create requires at least one embedding")
	}
	_ = v.http.Do(ctx, "POST", "/collection/drop",
		map[string]any{"database": v.cfg.Database, "collection": v.collection}, nil)
	body := map[string]any{
		"database":   v.cfg.Database,
		"collection": v.collection,
		"indexes": []map[string]any{
			{"fieldName": "id", "fieldType": "string", "indexType": "primaryKey"},
			{"fieldName": "vector", "fieldType": "vector", "indexType": "HNSW",
				"dimension": len(embs[0]), "metricType": "COSINE",
				"params": map[string]any{"M": 16, "efConstruction": 200}},
		},
	}
	if err := v.http.Do(ctx, "POST", "/collection/create", body, nil); err != nil {
		return fmt.Errorf("tencent: create collection: %w", err)
	}
	return v.AddTexts(ctx, docs, embs)
}

func (v *Vector) AddTexts(ctx context.Context, docs []*ragentities.Document, embs [][]float32) error {
	if len(docs) != len(embs) {
		return fmt.Errorf("tencent: docs=%d embeddings=%d mismatch", len(docs), len(embs))
	}
	rows := make([]map[string]any, len(docs))
	for i, d := range docs {
		id, _ := d.Metadata["doc_id"].(string)
		if id == "" {
			return fmt.Errorf("tencent: docs[%d] missing metadata.doc_id", i)
		}
		rows[i] = map[string]any{"id": id, "vector": embs[i], "text": d.PageContent, "metadata": d.Metadata}
	}
	return v.http.Do(ctx, "POST", "/document/upsert",
		map[string]any{"database": v.cfg.Database, "collection": v.collection, "documents": rows}, nil)
}

func (v *Vector) TextExists(ctx context.Context, id string) (bool, error) {
	var resp struct {
		Documents []any `json:"documents"`
	}
	err := v.http.Do(ctx, "POST", "/document/query",
		map[string]any{"database": v.cfg.Database, "collection": v.collection,
			"query": map[string]any{"documentIds": []string{id}, "retrieveVector": false}}, &resp)
	if err != nil {
		return false, err
	}
	return len(resp.Documents) > 0, nil
}

func (v *Vector) GetIDsByMetadataField(ctx context.Context, key, value string) ([]string, error) {
	body := map[string]any{
		"database":   v.cfg.Database,
		"collection": v.collection,
		"query":      map[string]any{"filter": fmt.Sprintf(`metadata.%s = "%s"`, key, value)},
	}
	var resp struct {
		Documents []struct {
			ID string `json:"id"`
		} `json:"documents"`
	}
	if err := v.http.Do(ctx, "POST", "/document/query", body, &resp); err != nil {
		return nil, err
	}
	ids := make([]string, len(resp.Documents))
	for i, d := range resp.Documents {
		ids[i] = d.ID
	}
	return ids, nil
}

func (v *Vector) DeleteByIDs(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	return v.http.Do(ctx, "POST", "/document/delete",
		map[string]any{"database": v.cfg.Database, "collection": v.collection,
			"query": map[string]any{"documentIds": ids}}, nil)
}

func (v *Vector) DeleteByMetadataField(ctx context.Context, key, value string) error {
	return v.http.Do(ctx, "POST", "/document/delete",
		map[string]any{"database": v.cfg.Database, "collection": v.collection,
			"query": map[string]any{"filter": fmt.Sprintf(`metadata.%s = "%s"`, key, value)}}, nil)
}

func (v *Vector) SearchByVector(ctx context.Context, qv []float32, opts ragentities.SearchOptions) ([]*ragentities.SearchResult, error) {
	topK := opts.TopK
	if topK <= 0 {
		topK = 4
	}
	body := map[string]any{
		"database":   v.cfg.Database,
		"collection": v.collection,
		"search": map[string]any{
			"vectors":        [][]float32{qv},
			"limit":          topK,
			"retrieveVector": false,
			"outputFields":   []string{"text", "metadata"},
		},
	}
	var resp struct {
		Documents [][]struct {
			Score    float32        `json:"score"`
			Text     string         `json:"text"`
			Metadata map[string]any `json:"metadata"`
		} `json:"documents"`
	}
	if err := v.http.Do(ctx, "POST", "/document/search", body, &resp); err != nil {
		return nil, err
	}
	var out []*ragentities.SearchResult
	if len(resp.Documents) > 0 {
		for _, r := range resp.Documents[0] {
			if r.Score < opts.ScoreThreshold {
				continue
			}
			out = append(out, &ragentities.SearchResult{
				Document: &ragentities.Document{PageContent: r.Text, Metadata: r.Metadata},
				Score:    r.Score,
			})
		}
	}
	return out, nil
}

func (v *Vector) SearchByFullText(ctx context.Context, q string, _ ragentities.SearchOptions) ([]*ragentities.SearchResult, error) {
	return nil, errors.New("tencent: full-text search not supported by this adapter")
}

func (v *Vector) Delete(ctx context.Context) error {
	return v.http.Do(ctx, "POST", "/collection/drop",
		map[string]any{"database": v.cfg.Database, "collection": v.collection}, nil)
}
