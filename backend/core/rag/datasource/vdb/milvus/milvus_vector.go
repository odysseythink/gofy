// Package milvus talks to Milvus >= 2.4 via its HTTP v2 API.
package milvus

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/odysseythink/confy"

	vdb "mlib.com/gofy/server/core/rag/datasource/vdb"
	"mlib.com/gofy/server/core/rag/datasource/vdb/internal/httpvec"
	ragentities "mlib.com/gofy/server/entities/rag"
	vectorenumtypes "mlib.com/gofy/server/enum_types/rag/vector"
)

func init() { vdb.Register(vectorenumtypes.Vector_MILVUS, &Factory{}) }

type Config struct {
	Endpoint string
	Token    string
	Database string
}

type Factory struct{}

func (Factory) New(ctx context.Context, cfg vdb.AdapterConfig) (ragentities.IVector, error) {
	return Open(ctx, loadConfig(), cfg.CollectionName)
}

func loadConfig() Config {
	return Config{
		Endpoint: confy.GetWithDefault[string]("milvus.endpoint", "http://localhost:19530"),
		Token:    confy.GetWithDefault[string]("milvus.token", ""),
		Database: confy.GetWithDefault[string]("milvus.database", "default"),
	}
}

func Open(ctx context.Context, cfg Config, collection string) (*Vector, error) {
	if collection == "" {
		return nil, errors.New("milvus: collection is required")
	}
	headers := map[string]string{}
	if cfg.Token != "" {
		headers["Authorization"] = "Bearer " + cfg.Token
	}
	return &Vector{
		cfg:        cfg,
		collection: sanitize(collection),
		http:       httpvec.New(cfg.Endpoint, 30*time.Second, headers),
	}, nil
}

type Vector struct {
	cfg        Config
	collection string
	http       *httpvec.Client
}

var _ ragentities.IVector = (*Vector)(nil)

func (v *Vector) GetType() string { return string(vectorenumtypes.Vector_MILVUS) }

func (v *Vector) Create(ctx context.Context, docs []*ragentities.Document, embs [][]float32) error {
	if len(embs) == 0 {
		return errors.New("milvus: Create requires at least one embedding")
	}
	_ = v.http.Do(ctx, "POST", "/v2/vectordb/collections/drop",
		map[string]any{"dbName": v.cfg.Database, "collectionName": v.collection}, nil)

	create := map[string]any{
		"dbName":           v.cfg.Database,
		"collectionName":   v.collection,
		"dimension":        len(embs[0]),
		"metricType":       "COSINE",
		"idType":           "VarChar",
		"primaryFieldName": "id",
		"vectorFieldName":  "vector",
	}
	if err := v.http.Do(ctx, "POST", "/v2/vectordb/collections/create", create, nil); err != nil {
		return fmt.Errorf("milvus: create collection: %w", err)
	}
	return v.AddTexts(ctx, docs, embs)
}

func (v *Vector) AddTexts(ctx context.Context, docs []*ragentities.Document, embs [][]float32) error {
	if len(docs) == 0 {
		return nil
	}
	if len(docs) != len(embs) {
		return fmt.Errorf("milvus: docs=%d embeddings=%d mismatch", len(docs), len(embs))
	}
	data := make([]map[string]any, len(docs))
	for i, d := range docs {
		id, _ := d.Metadata["doc_id"].(string)
		if id == "" {
			return fmt.Errorf("milvus: docs[%d] missing metadata.doc_id", i)
		}
		data[i] = map[string]any{
			"id":       id,
			"vector":   embs[i],
			"text":     d.PageContent,
			"metadata": d.Metadata,
		}
	}
	return v.http.Do(ctx, "POST", "/v2/vectordb/entities/upsert",
		map[string]any{"dbName": v.cfg.Database, "collectionName": v.collection, "data": data}, nil)
}

func (v *Vector) TextExists(ctx context.Context, id string) (bool, error) {
	var resp struct {
		Data []any `json:"data"`
	}
	err := v.http.Do(ctx, "POST", "/v2/vectordb/entities/get",
		map[string]any{"dbName": v.cfg.Database, "collectionName": v.collection, "id": []string{id}}, &resp)
	if err != nil {
		if httpvec.Is404(err) {
			return false, nil
		}
		return false, err
	}
	return len(resp.Data) > 0, nil
}

func (v *Vector) GetIDsByMetadataField(ctx context.Context, key, value string) ([]string, error) {
	var resp struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	body := map[string]any{
		"dbName":         v.cfg.Database,
		"collectionName": v.collection,
		"filter":         fmt.Sprintf(`metadata["%s"] == "%s"`, key, value),
		"outputFields":   []string{"id"},
		"limit":          1000,
	}
	if err := v.http.Do(ctx, "POST", "/v2/vectordb/entities/query", body, &resp); err != nil {
		return nil, err
	}
	ids := make([]string, len(resp.Data))
	for i, r := range resp.Data {
		ids[i] = r.ID
	}
	return ids, nil
}

func (v *Vector) DeleteByIDs(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	return v.http.Do(ctx, "POST", "/v2/vectordb/entities/delete",
		map[string]any{"dbName": v.cfg.Database, "collectionName": v.collection, "id": ids}, nil)
}

func (v *Vector) DeleteByMetadataField(ctx context.Context, key, value string) error {
	return v.http.Do(ctx, "POST", "/v2/vectordb/entities/delete",
		map[string]any{
			"dbName":         v.cfg.Database,
			"collectionName": v.collection,
			"filter":         fmt.Sprintf(`metadata["%s"] == "%s"`, key, value),
		}, nil)
}

func (v *Vector) SearchByVector(ctx context.Context, qv []float32, opts ragentities.SearchOptions) ([]*ragentities.SearchResult, error) {
	topK := opts.TopK
	if topK <= 0 {
		topK = 4
	}
	body := map[string]any{
		"dbName":         v.cfg.Database,
		"collectionName": v.collection,
		"data":           [][]float32{qv},
		"limit":          topK,
		"outputFields":   []string{"id", "text", "metadata"},
	}
	if filter := buildFilter(opts); filter != "" {
		body["filter"] = filter
	}
	var resp struct {
		Data []struct {
			Distance float32        `json:"distance"`
			Text     string         `json:"text"`
			Metadata map[string]any `json:"metadata"`
		} `json:"data"`
	}
	if err := v.http.Do(ctx, "POST", "/v2/vectordb/entities/search", body, &resp); err != nil {
		return nil, err
	}
	out := make([]*ragentities.SearchResult, 0, len(resp.Data))
	for _, r := range resp.Data {
		if r.Distance < opts.ScoreThreshold {
			continue
		}
		out = append(out, &ragentities.SearchResult{
			Document: &ragentities.Document{PageContent: r.Text, Metadata: r.Metadata},
			Score:    r.Distance,
		})
	}
	return out, nil
}

func (v *Vector) SearchByFullText(ctx context.Context, q string, _ ragentities.SearchOptions) ([]*ragentities.SearchResult, error) {
	return nil, errors.New("milvus: full-text search not supported by this adapter")
}

func (v *Vector) Delete(ctx context.Context) error {
	return v.http.Do(ctx, "POST", "/v2/vectordb/collections/drop",
		map[string]any{"dbName": v.cfg.Database, "collectionName": v.collection}, nil)
}

func buildFilter(opts ragentities.SearchOptions) string {
	var parts []string
	for k, v := range opts.Filter {
		parts = append(parts, fmt.Sprintf(`metadata["%s"] == "%v"`, k, v))
	}
	if len(opts.DocumentIDs) > 0 {
		quoted := make([]string, len(opts.DocumentIDs))
		for i, id := range opts.DocumentIDs {
			quoted[i] = `"` + id + `"`
		}
		parts = append(parts, fmt.Sprintf(`metadata["document_id"] in [%s]`, strings.Join(quoted, ",")))
	}
	return strings.Join(parts, " && ")
}

func sanitize(s string) string {
	var b strings.Builder
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9', r == '_':
			b.WriteRune(r)
		default:
			b.WriteByte('_')
		}
	}
	return b.String()
}
