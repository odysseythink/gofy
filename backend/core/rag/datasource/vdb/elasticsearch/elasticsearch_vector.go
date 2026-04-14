// Package elasticsearch is the IVector adapter for Elasticsearch 8+ with kNN.
// OpenSearch is handled by a thin re-registration below because its API is
// wire-compatible for the dense_vector subset we use.
package elasticsearch

import (
	"context"
	"encoding/base64"
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

func init() {
	vdb.Register(vectorenumtypes.Vector_ELASTICSEARCH, &Factory{cfgPrefix: "elasticsearch"})
	vdb.Register(vectorenumtypes.Vector_ELASTICSEARCH_JA, &Factory{cfgPrefix: "elasticsearch", analyzer: "kuromoji"})
	vdb.Register(vectorenumtypes.Vector_OPENSEARCH, &Factory{cfgPrefix: "opensearch"})
}

type Config struct {
	Endpoint string
	Username string
	Password string
	APIKey   string
	Analyzer string // optional; e.g. "kuromoji" for Japanese
}

type Factory struct {
	cfgPrefix string
	analyzer  string
}

func (f Factory) New(ctx context.Context, cfg vdb.AdapterConfig) (ragentities.IVector, error) {
	c := loadConfig(f.cfgPrefix)
	if f.analyzer != "" && c.Analyzer == "" {
		c.Analyzer = f.analyzer
	}
	return Open(ctx, c, cfg.CollectionName)
}

func loadConfig(prefix string) Config {
	return Config{
		Endpoint: confy.GetWithDefault[string](prefix+".endpoint", "http://localhost:9200"),
		Username: confy.GetWithDefault[string](prefix+".username", ""),
		Password: confy.GetWithDefault[string](prefix+".password", ""),
		APIKey:   confy.GetWithDefault[string](prefix+".api-key", ""),
		Analyzer: confy.GetWithDefault[string](prefix+".analyzer", ""),
	}
}

func Open(ctx context.Context, cfg Config, collection string) (*Vector, error) {
	if collection == "" {
		return nil, errors.New("elasticsearch: collection is required")
	}
	h := map[string]string{}
	if cfg.APIKey != "" {
		h["Authorization"] = "ApiKey " + cfg.APIKey
	} else if cfg.Username != "" {
		token := base64.StdEncoding.EncodeToString([]byte(cfg.Username + ":" + cfg.Password))
		h["Authorization"] = "Basic " + token
	}
	return &Vector{cfg: cfg, index: strings.ToLower(collection), http: httpvec.New(cfg.Endpoint, 30*time.Second, h)}, nil
}

type Vector struct {
	cfg   Config
	index string
	http  *httpvec.Client
}

var _ ragentities.IVector = (*Vector)(nil)

func (v *Vector) GetType() string { return string(vectorenumtypes.Vector_ELASTICSEARCH) }

func (v *Vector) Create(ctx context.Context, docs []*ragentities.Document, embs [][]float32) error {
	if len(embs) == 0 {
		return errors.New("elasticsearch: Create requires at least one embedding")
	}
	_ = v.http.Do(ctx, "DELETE", "/"+v.index, nil, nil)
	textProps := map[string]any{"type": "text"}
	if v.cfg.Analyzer != "" {
		textProps["analyzer"] = v.cfg.Analyzer
	}
	mapping := map[string]any{
		"mappings": map[string]any{
			"properties": map[string]any{
				"text":     textProps,
				"metadata": map[string]any{"type": "object"},
				"vector": map[string]any{
					"type":       "dense_vector",
					"dims":       len(embs[0]),
					"index":      true,
					"similarity": "cosine",
				},
			},
		},
	}
	if err := v.http.Do(ctx, "PUT", "/"+v.index, mapping, nil); err != nil {
		return fmt.Errorf("elasticsearch: create index: %w", err)
	}
	return v.AddTexts(ctx, docs, embs)
}

func (v *Vector) AddTexts(ctx context.Context, docs []*ragentities.Document, embs [][]float32) error {
	if len(docs) == 0 {
		return nil
	}
	if len(docs) != len(embs) {
		return fmt.Errorf("elasticsearch: docs=%d embeddings=%d mismatch", len(docs), len(embs))
	}
	for i, d := range docs {
		id, _ := d.Metadata["doc_id"].(string)
		if id == "" {
			return fmt.Errorf("elasticsearch: docs[%d] missing metadata.doc_id", i)
		}
		body := map[string]any{"text": d.PageContent, "metadata": d.Metadata, "vector": embs[i]}
		if err := v.http.Do(ctx, "PUT", "/"+v.index+"/_doc/"+id+"?refresh=true", body, nil); err != nil {
			return fmt.Errorf("elasticsearch: index %s: %w", id, err)
		}
	}
	return nil
}

func (v *Vector) TextExists(ctx context.Context, id string) (bool, error) {
	err := v.http.Do(ctx, "GET", "/"+v.index+"/_doc/"+id, nil, nil)
	if err == nil {
		return true, nil
	}
	if httpvec.Is404(err) {
		return false, nil
	}
	return false, err
}

func (v *Vector) GetIDsByMetadataField(ctx context.Context, key, value string) ([]string, error) {
	body := map[string]any{
		"_source": []string{"metadata.doc_id"},
		"query":   map[string]any{"term": map[string]any{"metadata." + key + ".keyword": value}},
		"size":    1000,
	}
	var resp struct {
		Hits struct {
			Hits []struct {
				Source map[string]any `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := v.http.Do(ctx, "POST", "/"+v.index+"/_search", body, &resp); err != nil {
		return nil, err
	}
	var out []string
	for _, h := range resp.Hits.Hits {
		if m, ok := h.Source["metadata"].(map[string]any); ok {
			if id, ok := m["doc_id"].(string); ok {
				out = append(out, id)
			}
		}
	}
	return out, nil
}

func (v *Vector) DeleteByIDs(ctx context.Context, ids []string) error {
	for _, id := range ids {
		_ = v.http.Do(ctx, "DELETE", "/"+v.index+"/_doc/"+id+"?refresh=true", nil, nil)
	}
	return nil
}

func (v *Vector) DeleteByMetadataField(ctx context.Context, key, value string) error {
	body := map[string]any{
		"query": map[string]any{"term": map[string]any{"metadata." + key + ".keyword": value}},
	}
	return v.http.Do(ctx, "POST", "/"+v.index+"/_delete_by_query?refresh=true", body, nil)
}

func (v *Vector) SearchByVector(ctx context.Context, qv []float32, opts ragentities.SearchOptions) ([]*ragentities.SearchResult, error) {
	topK := opts.TopK
	if topK <= 0 {
		topK = 4
	}
	knn := map[string]any{
		"field":          "vector",
		"query_vector":   qv,
		"k":              topK,
		"num_candidates": topK * 10,
	}
	if filter := buildESFilter(opts); filter != nil {
		knn["filter"] = filter
	}
	body := map[string]any{"knn": knn, "size": topK, "_source": []string{"text", "metadata"}}

	var resp struct {
		Hits struct {
			Hits []struct {
				Score  float32        `json:"_score"`
				Source map[string]any `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := v.http.Do(ctx, "POST", "/"+v.index+"/_search", body, &resp); err != nil {
		return nil, err
	}
	var out []*ragentities.SearchResult
	for _, h := range resp.Hits.Hits {
		if h.Score < opts.ScoreThreshold {
			continue
		}
		text, _ := h.Source["text"].(string)
		meta, _ := h.Source["metadata"].(map[string]any)
		out = append(out, &ragentities.SearchResult{
			Document: &ragentities.Document{PageContent: text, Metadata: meta},
			Score:    h.Score,
		})
	}
	return out, nil
}

func (v *Vector) SearchByFullText(ctx context.Context, q string, opts ragentities.SearchOptions) ([]*ragentities.SearchResult, error) {
	topK := opts.TopK
	if topK <= 0 {
		topK = 5
	}
	body := map[string]any{
		"query":   map[string]any{"match": map[string]any{"text": q}},
		"size":    topK,
		"_source": []string{"text", "metadata"},
	}
	var resp struct {
		Hits struct {
			Hits []struct {
				Score  float32        `json:"_score"`
				Source map[string]any `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := v.http.Do(ctx, "POST", "/"+v.index+"/_search", body, &resp); err != nil {
		return nil, err
	}
	var out []*ragentities.SearchResult
	for _, h := range resp.Hits.Hits {
		text, _ := h.Source["text"].(string)
		meta, _ := h.Source["metadata"].(map[string]any)
		out = append(out, &ragentities.SearchResult{
			Document: &ragentities.Document{PageContent: text, Metadata: meta},
			Score:    h.Score,
		})
	}
	return out, nil
}

func (v *Vector) Delete(ctx context.Context) error {
	return v.http.Do(ctx, "DELETE", "/"+v.index, nil, nil)
}

func buildESFilter(opts ragentities.SearchOptions) []map[string]any {
	var out []map[string]any
	for k, val := range opts.Filter {
		out = append(out, map[string]any{
			"term": map[string]any{"metadata." + k + ".keyword": fmt.Sprint(val)},
		})
	}
	if len(opts.DocumentIDs) > 0 {
		out = append(out, map[string]any{
			"terms": map[string]any{"metadata.document_id.keyword": opts.DocumentIDs},
		})
	}
	if len(out) == 0 {
		return nil
	}
	return out
}
