// Package weaviate talks to Weaviate via its REST v1 API.
package weaviate

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

func init() { vdb.Register(vectorenumtypes.Vector_WEAVIATE, &Factory{}) }

type Config struct {
	Endpoint string
	APIKey   string
}

type Factory struct{}

func (Factory) New(ctx context.Context, cfg vdb.AdapterConfig) (ragentities.IVector, error) {
	return Open(ctx, loadConfig(), cfg.CollectionName)
}

func loadConfig() Config {
	return Config{
		Endpoint: confy.GetWithDefault[string]("weaviate.endpoint", "http://localhost:8080"),
		APIKey:   confy.GetWithDefault[string]("weaviate.api-key", ""),
	}
}

func Open(ctx context.Context, cfg Config, collection string) (*Vector, error) {
	if collection == "" {
		return nil, errors.New("weaviate: collection is required")
	}
	h := map[string]string{}
	if cfg.APIKey != "" {
		h["Authorization"] = "Bearer " + cfg.APIKey
	}
	// Weaviate class names must start with uppercase letter.
	class := sanitizeClass(collection)
	return &Vector{cfg: cfg, class: class, http: httpvec.New(cfg.Endpoint, 30*time.Second, h)}, nil
}

type Vector struct {
	cfg   Config
	class string
	http  *httpvec.Client
}

var _ ragentities.IVector = (*Vector)(nil)

func (v *Vector) GetType() string { return string(vectorenumtypes.Vector_WEAVIATE) }

func (v *Vector) Create(ctx context.Context, docs []*ragentities.Document, embs [][]float32) error {
	if len(embs) == 0 {
		return errors.New("weaviate: Create requires at least one embedding")
	}
	_ = v.http.Do(ctx, "DELETE", "/v1/schema/"+v.class, nil, nil)
	schema := map[string]any{
		"class":      v.class,
		"vectorizer": "none",
		"properties": []map[string]any{
			{"name": "text", "dataType": []string{"text"}},
			{"name": "doc_id", "dataType": []string{"text"}},
			{"name": "document_id", "dataType": []string{"text"}},
			{"name": "metadata", "dataType": []string{"text"}}, // stringified JSON
		},
	}
	if err := v.http.Do(ctx, "POST", "/v1/schema", schema, nil); err != nil {
		return fmt.Errorf("weaviate: create class: %w", err)
	}
	return v.AddTexts(ctx, docs, embs)
}

func (v *Vector) AddTexts(ctx context.Context, docs []*ragentities.Document, embs [][]float32) error {
	if len(docs) == 0 {
		return nil
	}
	if len(docs) != len(embs) {
		return fmt.Errorf("weaviate: docs=%d embeddings=%d mismatch", len(docs), len(embs))
	}
	objects := make([]map[string]any, len(docs))
	for i, d := range docs {
		id, _ := d.Metadata["doc_id"].(string)
		if id == "" {
			return fmt.Errorf("weaviate: docs[%d] missing metadata.doc_id", i)
		}
		did, _ := d.Metadata["document_id"].(string)
		objects[i] = map[string]any{
			"class":  v.class,
			"id":     hashUUID(id),
			"vector": embs[i],
			"properties": map[string]any{
				"text":        d.PageContent,
				"doc_id":      id,
				"document_id": did,
				"metadata":    toJSONString(d.Metadata),
			},
		}
	}
	return v.http.Do(ctx, "POST", "/v1/batch/objects",
		map[string]any{"objects": objects}, nil)
}

func (v *Vector) TextExists(ctx context.Context, id string) (bool, error) {
	err := v.http.Do(ctx, "GET", "/v1/objects/"+v.class+"/"+hashUUID(id), nil, nil)
	if err == nil {
		return true, nil
	}
	if httpvec.Is404(err) {
		return false, nil
	}
	return false, err
}

func (v *Vector) GetIDsByMetadataField(ctx context.Context, key, value string) ([]string, error) {
	q := fmt.Sprintf(`{Get{%s(where:{path:["%s"],operator:Equal,valueText:%q}){doc_id}}}`, v.class, key, value)
	var resp struct {
		Data struct {
			Get map[string][]struct {
				DocID string `json:"doc_id"`
			} `json:"Get"`
		} `json:"data"`
	}
	if err := v.http.Do(ctx, "POST", "/v1/graphql", map[string]any{"query": q}, &resp); err != nil {
		return nil, err
	}
	var ids []string
	for _, r := range resp.Data.Get[v.class] {
		ids = append(ids, r.DocID)
	}
	return ids, nil
}

func (v *Vector) DeleteByIDs(ctx context.Context, ids []string) error {
	for _, id := range ids {
		_ = v.http.Do(ctx, "DELETE", "/v1/objects/"+v.class+"/"+hashUUID(id), nil, nil)
	}
	return nil
}

func (v *Vector) DeleteByMetadataField(ctx context.Context, key, value string) error {
	body := map[string]any{
		"class": v.class,
		"where": map[string]any{"path": []string{key}, "operator": "Equal", "valueText": value},
	}
	return v.http.Do(ctx, "POST", "/v1/batch/objects/delete", body, nil)
}

func (v *Vector) SearchByVector(ctx context.Context, qv []float32, opts ragentities.SearchOptions) ([]*ragentities.SearchResult, error) {
	topK := opts.TopK
	if topK <= 0 {
		topK = 4
	}
	where := buildGraphQLWhere(opts)
	q := fmt.Sprintf(`{Get{%s(nearVector:{vector:%s} limit:%d %s){text doc_id document_id metadata _additional{certainty}}}}`,
		v.class, floatSlice(qv), topK, where)
	var resp struct {
		Data struct {
			Get map[string][]struct {
				Text       string         `json:"text"`
				DocID      string         `json:"doc_id"`
				DocumentID string         `json:"document_id"`
				Metadata   string         `json:"metadata"`
				Additional map[string]any `json:"_additional"`
			} `json:"Get"`
		} `json:"data"`
	}
	if err := v.http.Do(ctx, "POST", "/v1/graphql", map[string]any{"query": q}, &resp); err != nil {
		return nil, err
	}
	var out []*ragentities.SearchResult
	for _, r := range resp.Data.Get[v.class] {
		score := float32(0)
		if c, ok := r.Additional["certainty"].(float64); ok {
			score = float32(c)
		}
		if score < opts.ScoreThreshold {
			continue
		}
		meta := map[string]any{"doc_id": r.DocID, "document_id": r.DocumentID}
		_ = jsonUnmarshal(r.Metadata, &meta)
		out = append(out, &ragentities.SearchResult{
			Document: &ragentities.Document{PageContent: r.Text, Metadata: meta},
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
	query := fmt.Sprintf(`{Get{%s(bm25:{query:%q} limit:%d){text doc_id document_id metadata _additional{score}}}}`,
		v.class, q, topK)
	var resp struct {
		Data struct {
			Get map[string][]struct {
				Text       string         `json:"text"`
				DocID      string         `json:"doc_id"`
				DocumentID string         `json:"document_id"`
				Metadata   string         `json:"metadata"`
				Additional map[string]any `json:"_additional"`
			} `json:"Get"`
		} `json:"data"`
	}
	if err := v.http.Do(ctx, "POST", "/v1/graphql", map[string]any{"query": query}, &resp); err != nil {
		return nil, err
	}
	var out []*ragentities.SearchResult
	for _, r := range resp.Data.Get[v.class] {
		score := float32(0)
		if s, ok := r.Additional["score"].(float64); ok {
			score = float32(s)
		}
		meta := map[string]any{"doc_id": r.DocID, "document_id": r.DocumentID}
		_ = jsonUnmarshal(r.Metadata, &meta)
		out = append(out, &ragentities.SearchResult{
			Document: &ragentities.Document{PageContent: r.Text, Metadata: meta},
			Score:    score,
		})
	}
	return out, nil
}

func (v *Vector) Delete(ctx context.Context) error {
	return v.http.Do(ctx, "DELETE", "/v1/schema/"+v.class, nil, nil)
}
