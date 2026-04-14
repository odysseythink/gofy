// Package qdrant is the Qdrant IVector adapter, talking to qdrant over its REST API.
//
// Uses only stdlib net/http to avoid pulling in a SDK.
package qdrant

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
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

func init() {
	vdb.Register(vectorenumtypes.Vector_QDRANT, &Factory{})
	vdb.Register(vectorenumtypes.Vector_TIDB_ON_QDRANT, &Factory{}) // same REST, different config prefix
}

type Config struct {
	Endpoint string // http://localhost:6333
	APIKey   string
	Timeout  time.Duration
}

type Factory struct{}

func (Factory) New(ctx context.Context, cfg vdb.AdapterConfig) (ragentities.IVector, error) {
	c := loadConfig(cfg.ConfigPrefix)
	return Open(ctx, c, cfg.CollectionName)
}

func loadConfig(prefix string) Config {
	if prefix == "" {
		prefix = "qdrant"
	}
	return Config{
		Endpoint: confy.GetWithDefault[string](prefix+".endpoint", "http://localhost:6333"),
		APIKey:   confy.GetWithDefault[string](prefix+".api-key", ""),
		Timeout:  time.Duration(confy.GetWithDefault[int](prefix+".timeout-seconds", 30)) * time.Second,
	}
}

func Open(ctx context.Context, cfg Config, collection string) (*Vector, error) {
	if collection == "" {
		return nil, errors.New("qdrant: collection is required")
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 30 * time.Second
	}
	return &Vector{
		cfg:        cfg,
		collection: sanitize(collection),
		client:     &http.Client{Timeout: cfg.Timeout},
	}, nil
}

type Vector struct {
	cfg        Config
	collection string
	client     *http.Client
}

var _ ragentities.IVector = (*Vector)(nil)

func (v *Vector) GetType() string { return string(vectorenumtypes.Vector_QDRANT) }

func (v *Vector) Create(ctx context.Context, docs []*ragentities.Document, embs [][]float32) error {
	if len(embs) == 0 {
		return errors.New("qdrant: Create requires at least one embedding")
	}
	// Upsert-style: delete-then-create so dimensions can change.
	_ = v.deleteCollection(ctx)
	body := map[string]any{
		"vectors": map[string]any{
			"size":     len(embs[0]),
			"distance": "Cosine",
		},
	}
	if err := v.do(ctx, http.MethodPut, "/collections/"+v.collection, body, nil); err != nil {
		return fmt.Errorf("qdrant: create collection: %w", err)
	}
	// Payload index on document_id for filter pushdown
	_ = v.do(ctx, http.MethodPut,
		"/collections/"+v.collection+"/index",
		map[string]any{"field_name": "document_id", "field_schema": "keyword"}, nil)

	return v.AddTexts(ctx, docs, embs)
}

func (v *Vector) AddTexts(ctx context.Context, docs []*ragentities.Document, embs [][]float32) error {
	if len(docs) == 0 {
		return nil
	}
	if len(docs) != len(embs) {
		return fmt.Errorf("qdrant: docs=%d embeddings=%d mismatch", len(docs), len(embs))
	}
	points := make([]map[string]any, 0, len(docs))
	for i, d := range docs {
		id, _ := d.Metadata["doc_id"].(string)
		if id == "" {
			return fmt.Errorf("qdrant: docs[%d] missing metadata.doc_id", i)
		}
		// Qdrant accepts UUIDs or positive integers as point IDs. Hash the doc_id
		// to a deterministic UUID5-like value so arbitrary string keys work.
		payload := map[string]any{"page_content": d.PageContent, "metadata": d.Metadata}
		// Also flatten document_id into payload for filter index above
		if did, ok := d.Metadata["document_id"]; ok {
			payload["document_id"] = did
		}
		points = append(points, map[string]any{
			"id":      hashUUID(id),
			"vector":  embs[i],
			"payload": payload,
		})
	}
	return v.do(ctx, http.MethodPut,
		"/collections/"+v.collection+"/points?wait=true",
		map[string]any{"points": points}, nil)
}

func (v *Vector) TextExists(ctx context.Context, id string) (bool, error) {
	var resp struct {
		Result []map[string]any `json:"result"`
	}
	err := v.do(ctx, http.MethodPost,
		"/collections/"+v.collection+"/points",
		map[string]any{"ids": []string{hashUUID(id)}, "with_payload": false, "with_vector": false},
		&resp)
	if err != nil {
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "Not Found") {
			return false, nil
		}
		return false, err
	}
	return len(resp.Result) > 0, nil
}

func (v *Vector) GetIDsByMetadataField(ctx context.Context, key, value string) ([]string, error) {
	var resp struct {
		Result struct {
			Points []struct {
				ID      any            `json:"id"`
				Payload map[string]any `json:"payload"`
			} `json:"points"`
		} `json:"result"`
	}
	body := map[string]any{
		"filter": map[string]any{
			"must": []map[string]any{
				{"key": "metadata." + key, "match": map[string]any{"value": value}},
			},
		},
		"limit":        1000,
		"with_payload": true,
	}
	if err := v.do(ctx, http.MethodPost, "/collections/"+v.collection+"/points/scroll", body, &resp); err != nil {
		return nil, err
	}
	var out []string
	for _, p := range resp.Result.Points {
		if m, ok := p.Payload["metadata"].(map[string]any); ok {
			if id, ok := m["doc_id"].(string); ok {
				out = append(out, id)
			}
		}
	}
	return out, nil
}

func (v *Vector) DeleteByIDs(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	pts := make([]string, len(ids))
	for i, id := range ids {
		pts[i] = hashUUID(id)
	}
	return v.do(ctx, http.MethodPost,
		"/collections/"+v.collection+"/points/delete?wait=true",
		map[string]any{"points": pts}, nil)
}

func (v *Vector) DeleteByMetadataField(ctx context.Context, key, value string) error {
	body := map[string]any{
		"filter": map[string]any{
			"must": []map[string]any{
				{"key": "metadata." + key, "match": map[string]any{"value": value}},
			},
		},
	}
	return v.do(ctx, http.MethodPost,
		"/collections/"+v.collection+"/points/delete?wait=true", body, nil)
}

func (v *Vector) SearchByVector(ctx context.Context, qv []float32, opts ragentities.SearchOptions) ([]*ragentities.SearchResult, error) {
	topK := opts.TopK
	if topK <= 0 {
		topK = 4
	}
	body := map[string]any{
		"vector":       qv,
		"limit":        topK,
		"with_payload": true,
	}
	if filter := buildFilter(opts); filter != nil {
		body["filter"] = filter
	}

	var resp struct {
		Result []struct {
			Score   float32        `json:"score"`
			Payload map[string]any `json:"payload"`
		} `json:"result"`
	}
	if err := v.do(ctx, http.MethodPost, "/collections/"+v.collection+"/points/search", body, &resp); err != nil {
		return nil, err
	}
	var out []*ragentities.SearchResult
	for _, r := range resp.Result {
		if r.Score < opts.ScoreThreshold {
			continue
		}
		page, _ := r.Payload["page_content"].(string)
		meta, _ := r.Payload["metadata"].(map[string]any)
		out = append(out, &ragentities.SearchResult{
			Document: &ragentities.Document{PageContent: page, Metadata: meta},
			Score:    r.Score,
		})
	}
	return out, nil
}

// SearchByFullText: qdrant needs full-text index setup per field and a MatchText
// filter; treat as unsupported for the initial port.
func (v *Vector) SearchByFullText(ctx context.Context, query string, opts ragentities.SearchOptions) ([]*ragentities.SearchResult, error) {
	return nil, errors.New("qdrant: full-text search not supported by this adapter")
}

func (v *Vector) Delete(ctx context.Context) error {
	return v.deleteCollection(ctx)
}

func (v *Vector) deleteCollection(ctx context.Context) error {
	return v.do(ctx, http.MethodDelete, "/collections/"+v.collection, nil, nil)
}

// ---- helpers ----

func (v *Vector) do(ctx context.Context, method, path string, body, out any) error {
	var reader io.Reader
	if body != nil {
		buf, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reader = bytes.NewReader(buf)
	}
	req, err := http.NewRequestWithContext(ctx, method, v.cfg.Endpoint+path, reader)
	if err != nil {
		return err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if v.cfg.APIKey != "" {
		req.Header.Set("api-key", v.cfg.APIKey)
	}
	resp, err := v.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	buf, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return fmt.Errorf("qdrant: %s %s -> %d: %s", method, path, resp.StatusCode, string(buf))
	}
	if out != nil && len(buf) > 0 {
		return json.Unmarshal(buf, out)
	}
	return nil
}

func buildFilter(opts ragentities.SearchOptions) map[string]any {
	var conds []map[string]any
	for k, val := range opts.Filter {
		conds = append(conds, map[string]any{
			"key":   "metadata." + k,
			"match": map[string]any{"value": fmt.Sprint(val)},
		})
	}
	if len(opts.DocumentIDs) > 0 {
		conds = append(conds, map[string]any{
			"key":   "document_id",
			"match": map[string]any{"any": opts.DocumentIDs},
		})
	}
	if len(conds) == 0 {
		return nil
	}
	return map[string]any{"must": conds}
}

// hashUUID produces a deterministic UUIDv5-style string from an arbitrary id,
// so qdrant (which requires UUID or uint point ids) accepts our doc_ids.
func hashUUID(s string) string {
	sum := sha1.Sum([]byte("gofy-qdrant:" + s))
	h := hex.EncodeToString(sum[:16])
	return fmt.Sprintf("%s-%s-%s-%s-%s", h[0:8], h[8:12], h[12:16], h[16:20], h[20:32])
}

func sanitize(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9', r == '_', r == '-':
			b.WriteRune(r)
		default:
			b.WriteByte('_')
		}
	}
	return b.String()
}
