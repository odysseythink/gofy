// Package vdbtest provides a contract test suite every IVector adapter must
// pass, plus an in-memory reference implementation used to validate the suite
// itself and to unblock higher layers before real adapters exist.
package vdbtest

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sort"
	"sync"

	ragentities "github.com/odysseythink/gofy/backend/entities/rag"
)

// Memory is a simple thread-safe in-memory IVector. It supports vector search
// (cosine similarity) and naive metadata equality filters. It does NOT support
// full-text search (returns ErrFullTextUnsupported).
type Memory struct {
	mu         sync.RWMutex
	collection string

	// indexed by doc_id
	items map[string]memItem
}

type memItem struct {
	doc    *ragentities.Document
	vector []float32
}

// ErrFullTextUnsupported is returned by SearchByFullText on backends that
// don't implement sparse retrieval. Adapters may return a wrapped version.
var ErrFullTextUnsupported = errors.New("vdb: full-text search not supported by this backend")

func NewMemory(collection string) *Memory {
	return &Memory{collection: collection, items: map[string]memItem{}}
}

func (m *Memory) GetType() string { return "memory" }

func (m *Memory) Create(ctx context.Context, texts []*ragentities.Document, embeddings [][]float32) error {
	m.mu.Lock()
	m.items = map[string]memItem{}
	m.mu.Unlock()
	if len(texts) == 0 {
		return nil
	}
	return m.AddTexts(ctx, texts, embeddings)
}

func (m *Memory) AddTexts(ctx context.Context, docs []*ragentities.Document, embeddings [][]float32) error {
	if len(docs) != len(embeddings) {
		return fmt.Errorf("memory: docs=%d embeddings=%d length mismatch", len(docs), len(embeddings))
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	for i, d := range docs {
		id, _ := d.Metadata["doc_id"].(string)
		if id == "" {
			return fmt.Errorf("memory: doc[%d] missing metadata.doc_id", i)
		}
		m.items[id] = memItem{doc: d, vector: embeddings[i]}
	}
	return nil
}

func (m *Memory) TextExists(ctx context.Context, id string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, ok := m.items[id]
	return ok, nil
}

func (m *Memory) GetIDsByMetadataField(ctx context.Context, key, value string) ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var out []string
	for id, it := range m.items {
		if s, ok := it.doc.Metadata[key].(string); ok && s == value {
			out = append(out, id)
		}
	}
	sort.Strings(out)
	return out, nil
}

func (m *Memory) DeleteByIDs(ctx context.Context, ids []string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, id := range ids {
		delete(m.items, id)
	}
	return nil
}

func (m *Memory) DeleteByMetadataField(ctx context.Context, key, value string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for id, it := range m.items {
		if s, ok := it.doc.Metadata[key].(string); ok && s == value {
			delete(m.items, id)
		}
	}
	return nil
}

func (m *Memory) SearchByVector(ctx context.Context, q []float32, opts ragentities.SearchOptions) ([]*ragentities.SearchResult, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	results := make([]*ragentities.SearchResult, 0, len(m.items))
	for _, it := range m.items {
		if !matchFilter(it.doc.Metadata, opts.Filter) {
			continue
		}
		if len(opts.DocumentIDs) > 0 && !containsDocumentID(it.doc.Metadata, opts.DocumentIDs) {
			continue
		}
		score := cosine(q, it.vector)
		if score < opts.ScoreThreshold {
			continue
		}
		results = append(results, &ragentities.SearchResult{Document: it.doc, Score: score})
	}
	sort.Slice(results, func(i, j int) bool { return results[i].Score > results[j].Score })
	if opts.TopK > 0 && len(results) > opts.TopK {
		results = results[:opts.TopK]
	}
	return results, nil
}

func (m *Memory) SearchByFullText(ctx context.Context, query string, opts ragentities.SearchOptions) ([]*ragentities.SearchResult, error) {
	return nil, ErrFullTextUnsupported
}

func (m *Memory) Delete(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.items = map[string]memItem{}
	return nil
}

func matchFilter(meta map[string]any, filter map[string]any) bool {
	for k, want := range filter {
		if meta[k] != want {
			return false
		}
	}
	return true
}

func containsDocumentID(meta map[string]any, ids []string) bool {
	got, _ := meta["document_id"].(string)
	for _, id := range ids {
		if id == got {
			return true
		}
	}
	return false
}

func cosine(a, b []float32) float32 {
	if len(a) == 0 || len(a) != len(b) {
		return 0
	}
	var dot, na, nb float64
	for i := range a {
		x, y := float64(a[i]), float64(b[i])
		dot += x * y
		na += x * x
		nb += y * y
	}
	if na == 0 || nb == 0 {
		return 0
	}
	return float32(dot / (math.Sqrt(na) * math.Sqrt(nb)))
}
