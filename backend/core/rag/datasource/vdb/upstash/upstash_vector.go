// Package upstash is the Upstash Vector IVector adapter.
// Upstash Vector is REST-first; auth is a bearer token.
package upstash

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

func init() { vdb.Register(vectorenumtypes.Vector_UPSTASH, &Factory{}) }

type Config struct {
	Endpoint string // e.g. https://xxx-us1-vector.upstash.io
	Token    string
}

type Factory struct{}

func (Factory) New(ctx context.Context, cfg vdb.AdapterConfig) (ragentities.IVector, error) {
	return Open(ctx, loadConfig(), cfg.CollectionName)
}

func loadConfig() Config {
	return Config{
		Endpoint: confy.GetWithDefault[string]("upstash.endpoint", ""),
		Token:    confy.GetWithDefault[string]("upstash.token", ""),
	}
}

func Open(ctx context.Context, cfg Config, collection string) (*Vector, error) {
	if cfg.Endpoint == "" {
		return nil, errors.New("upstash: endpoint is required")
	}
	h := map[string]string{"Authorization": "Bearer " + cfg.Token}
	return &Vector{ns: collection, http: httpvec.New(cfg.Endpoint, 30*time.Second, h)}, nil
}

// Upstash uses namespaces instead of collections; one namespace per IVector.
type Vector struct {
	ns   string
	http *httpvec.Client
}

var _ ragentities.IVector = (*Vector)(nil)

func (v *Vector) GetType() string { return string(vectorenumtypes.Vector_UPSTASH) }

func (v *Vector) Create(ctx context.Context, docs []*ragentities.Document, embs [][]float32) error {
	_ = v.Delete(ctx)
	return v.AddTexts(ctx, docs, embs)
}

func (v *Vector) AddTexts(ctx context.Context, docs []*ragentities.Document, embs [][]float32) error {
	if len(docs) != len(embs) {
		return fmt.Errorf("upstash: docs=%d embeddings=%d mismatch", len(docs), len(embs))
	}
	vectors := make([]map[string]any, len(docs))
	for i, d := range docs {
		id, _ := d.Metadata["doc_id"].(string)
		if id == "" {
			return fmt.Errorf("upstash: docs[%d] missing metadata.doc_id", i)
		}
		vectors[i] = map[string]any{
			"id":       id,
			"vector":   embs[i],
			"metadata": d.Metadata,
			"data":     d.PageContent,
		}
	}
	return v.http.Do(ctx, "POST", "/upsert/"+v.ns, vectors, nil)
}

func (v *Vector) TextExists(ctx context.Context, id string) (bool, error) {
	var resp struct {
		Result []any `json:"result"`
	}
	err := v.http.Do(ctx, "POST", "/fetch/"+v.ns,
		map[string]any{"ids": []string{id}}, &resp)
	if err != nil {
		return false, err
	}
	return len(resp.Result) > 0 && resp.Result[0] != nil, nil
}

func (v *Vector) GetIDsByMetadataField(ctx context.Context, key, value string) ([]string, error) {
	// Upstash doesn't offer listing by metadata in v1; rely on query.
	return nil, errors.New("upstash: GetIDsByMetadataField not supported")
}

func (v *Vector) DeleteByIDs(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	return v.http.Do(ctx, "DELETE", "/delete/"+v.ns, ids, nil)
}

func (v *Vector) DeleteByMetadataField(ctx context.Context, key, value string) error {
	body := map[string]any{"filter": fmt.Sprintf(`%s = "%s"`, key, value)}
	return v.http.Do(ctx, "DELETE", "/delete/"+v.ns, body, nil)
}

func (v *Vector) SearchByVector(ctx context.Context, qv []float32, opts ragentities.SearchOptions) ([]*ragentities.SearchResult, error) {
	topK := opts.TopK
	if topK <= 0 {
		topK = 4
	}
	body := map[string]any{
		"vector":          qv,
		"topK":            topK,
		"includeMetadata": true,
		"includeData":     true,
	}
	if filter := buildFilter(opts); filter != "" {
		body["filter"] = filter
	}
	var resp []struct {
		ID       string         `json:"id"`
		Score    float32        `json:"score"`
		Data     string         `json:"data"`
		Metadata map[string]any `json:"metadata"`
	}
	if err := v.http.Do(ctx, "POST", "/query/"+v.ns, body, &resp); err != nil {
		return nil, err
	}
	var out []*ragentities.SearchResult
	for _, r := range resp {
		if r.Score < opts.ScoreThreshold {
			continue
		}
		out = append(out, &ragentities.SearchResult{
			Document: &ragentities.Document{PageContent: r.Data, Metadata: r.Metadata},
			Score:    r.Score,
		})
	}
	return out, nil
}

func (v *Vector) SearchByFullText(ctx context.Context, q string, _ ragentities.SearchOptions) ([]*ragentities.SearchResult, error) {
	return nil, errors.New("upstash: full-text search not supported by this adapter")
}

func (v *Vector) Delete(ctx context.Context) error {
	return v.http.Do(ctx, "DELETE", "/reset/"+v.ns, nil, nil)
}

func buildFilter(opts ragentities.SearchOptions) string {
	var parts []string
	for k, val := range opts.Filter {
		parts = append(parts, fmt.Sprintf(`%s = "%v"`, k, val))
	}
	if len(opts.DocumentIDs) > 0 {
		ids := make([]string, len(opts.DocumentIDs))
		for i, id := range opts.DocumentIDs {
			ids[i] = `"` + id + `"`
		}
		parts = append(parts, fmt.Sprintf(`document_id IN (%s)`, joinWith(ids, ",")))
	}
	return joinWith(parts, " AND ")
}

func joinWith(ss []string, sep string) string {
	if len(ss) == 0 {
		return ""
	}
	out := ss[0]
	for _, s := range ss[1:] {
		out += sep + s
	}
	return out
}
