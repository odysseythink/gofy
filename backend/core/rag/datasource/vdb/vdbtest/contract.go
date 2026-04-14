package vdbtest

import (
	"context"
	"errors"
	"testing"

	ragentities "mlib.com/gofy/server/entities/rag"
)

// NewAdapter produces a fresh IVector bound to a unique collection name.
// Each test case calls it to get an isolated instance; the second returned
// function is the cleanup hook (Delete + any teardown).
type NewAdapter func(t *testing.T, collection string) (ragentities.IVector, func())

// Capabilities lets adapters opt out of optional features. Required methods
// (Create/AddTexts/TextExists/SearchByVector/Delete) must always work.
type Capabilities struct {
	// FullText: adapter supports SearchByFullText. If false the suite expects
	// SearchByFullText to return an error (typically ErrFullTextUnsupported).
	FullText bool
	// MetadataFilter: adapter supports opts.Filter equality push-down. If
	// false the suite will skip filter-based assertions.
	MetadataFilter bool
	// MetadataDelete: adapter supports DeleteByMetadataField.
	MetadataDelete bool
}

// RunContractSuite exercises every IVector method against the adapter. Call it
// from an adapter-specific *_test.go:
//
//	func TestPgvectorContract(t *testing.T) {
//		if dsn := os.Getenv("PGVECTOR_TEST_DSN"); dsn == "" { t.Skip(...) }
//		vdbtest.RunContractSuite(t, func(t *testing.T, name string) (IVector, func()) {
//			impl := openPgvector(t, dsn, name)
//			return impl, func() { impl.Delete(context.Background()) }
//		}, vdbtest.Capabilities{MetadataFilter: true, MetadataDelete: true})
//	}
func RunContractSuite(t *testing.T, newAdapter NewAdapter, caps Capabilities) {
	t.Helper()

	t.Run("CreateAndAddTexts", func(t *testing.T) { testCreateAdd(t, newAdapter) })
	t.Run("TextExists", func(t *testing.T) { testTextExists(t, newAdapter) })
	t.Run("SearchByVector_Ranking", func(t *testing.T) { testSearchRanking(t, newAdapter) })
	t.Run("SearchByVector_TopK", func(t *testing.T) { testSearchTopK(t, newAdapter) })
	t.Run("DeleteByIDs", func(t *testing.T) { testDeleteByIDs(t, newAdapter) })
	if caps.MetadataFilter {
		t.Run("SearchByVector_Filter", func(t *testing.T) { testSearchFilter(t, newAdapter) })
	}
	if caps.MetadataDelete {
		t.Run("DeleteByMetadataField", func(t *testing.T) { testDeleteByMeta(t, newAdapter) })
	}
	t.Run("FullText", func(t *testing.T) { testFullText(t, newAdapter, caps.FullText) })
	t.Run("DeleteCollection", func(t *testing.T) { testDeleteCollection(t, newAdapter) })
}

// ---------- fixtures ----------

func fixtureDocs() ([]*ragentities.Document, [][]float32) {
	// Three docs arranged in a line in 4-D space so cosine similarity is
	// deterministic and ordering under a query is easy to reason about.
	docs := []*ragentities.Document{
		{PageContent: "alpha", Metadata: map[string]any{"doc_id": "a", "document_id": "doc-1", "lang": "en"}},
		{PageContent: "beta", Metadata: map[string]any{"doc_id": "b", "document_id": "doc-1", "lang": "en"}},
		{PageContent: "gamma", Metadata: map[string]any{"doc_id": "c", "document_id": "doc-2", "lang": "zh"}},
	}
	embeddings := [][]float32{
		{1, 0, 0, 0},
		{0.9, 0.1, 0, 0},
		{0, 0, 1, 0},
	}
	return docs, embeddings
}

// ---------- test cases ----------

func testCreateAdd(t *testing.T, newAdapter NewAdapter) {
	v, cleanup := newAdapter(t, "contract_create")
	defer cleanup()
	ctx := context.Background()
	docs, embs := fixtureDocs()
	must(t, v.Create(ctx, docs, embs), "Create")

	for _, d := range docs {
		ok, err := v.TextExists(ctx, d.Metadata["doc_id"].(string))
		must(t, err, "TextExists")
		if !ok {
			t.Fatalf("doc %v not indexed after Create", d.Metadata["doc_id"])
		}
	}
}

func testTextExists(t *testing.T, newAdapter NewAdapter) {
	v, cleanup := newAdapter(t, "contract_exists")
	defer cleanup()
	ctx := context.Background()
	docs, embs := fixtureDocs()
	must(t, v.Create(ctx, docs, embs), "Create")

	ok, err := v.TextExists(ctx, "a")
	must(t, err, "TextExists present")
	if !ok {
		t.Fatal("expected doc a to exist")
	}
	ok, err = v.TextExists(ctx, "does-not-exist")
	must(t, err, "TextExists missing")
	if ok {
		t.Fatal("expected missing doc to report absent")
	}
}

func testSearchRanking(t *testing.T, newAdapter NewAdapter) {
	v, cleanup := newAdapter(t, "contract_rank")
	defer cleanup()
	ctx := context.Background()
	docs, embs := fixtureDocs()
	must(t, v.Create(ctx, docs, embs), "Create")

	// Query aligned with doc a ("alpha"); expect a first, b second, c last.
	hits, err := v.SearchByVector(ctx, []float32{1, 0, 0, 0}, ragentities.SearchOptions{TopK: 3})
	must(t, err, "Search")
	if len(hits) == 0 {
		t.Fatal("no results")
	}
	if id := hits[0].Document.Metadata["doc_id"]; id != "a" {
		t.Fatalf("expected top hit 'a', got %v (full=%v)", id, idsOf(hits))
	}
	// Non-increasing order
	for i := 1; i < len(hits); i++ {
		if hits[i-1].Score < hits[i].Score {
			t.Fatalf("results not sorted desc by Score: %v", scoresOf(hits))
		}
	}
}

func testSearchTopK(t *testing.T, newAdapter NewAdapter) {
	v, cleanup := newAdapter(t, "contract_topk")
	defer cleanup()
	ctx := context.Background()
	docs, embs := fixtureDocs()
	must(t, v.Create(ctx, docs, embs), "Create")

	hits, err := v.SearchByVector(ctx, []float32{1, 0, 0, 0}, ragentities.SearchOptions{TopK: 1})
	must(t, err, "Search")
	if len(hits) != 1 {
		t.Fatalf("expected 1 hit, got %d: %v", len(hits), idsOf(hits))
	}
}

func testSearchFilter(t *testing.T, newAdapter NewAdapter) {
	v, cleanup := newAdapter(t, "contract_filter")
	defer cleanup()
	ctx := context.Background()
	docs, embs := fixtureDocs()
	must(t, v.Create(ctx, docs, embs), "Create")

	hits, err := v.SearchByVector(ctx, []float32{1, 0, 0, 0}, ragentities.SearchOptions{
		TopK:   10,
		Filter: map[string]any{"lang": "zh"},
	})
	must(t, err, "Search")
	for _, h := range hits {
		if h.Document.Metadata["lang"] != "zh" {
			t.Fatalf("filter violated: %v", h.Document.Metadata)
		}
	}
}

func testDeleteByIDs(t *testing.T, newAdapter NewAdapter) {
	v, cleanup := newAdapter(t, "contract_del_ids")
	defer cleanup()
	ctx := context.Background()
	docs, embs := fixtureDocs()
	must(t, v.Create(ctx, docs, embs), "Create")

	must(t, v.DeleteByIDs(ctx, []string{"a", "b"}), "DeleteByIDs")
	for _, id := range []string{"a", "b"} {
		ok, err := v.TextExists(ctx, id)
		must(t, err, "TextExists")
		if ok {
			t.Fatalf("%s should have been deleted", id)
		}
	}
	ok, err := v.TextExists(ctx, "c")
	must(t, err, "TextExists c")
	if !ok {
		t.Fatal("c should still exist")
	}
}

func testDeleteByMeta(t *testing.T, newAdapter NewAdapter) {
	v, cleanup := newAdapter(t, "contract_del_meta")
	defer cleanup()
	ctx := context.Background()
	docs, embs := fixtureDocs()
	must(t, v.Create(ctx, docs, embs), "Create")

	must(t, v.DeleteByMetadataField(ctx, "lang", "en"), "DeleteByMetadataField")
	for _, id := range []string{"a", "b"} {
		ok, _ := v.TextExists(ctx, id)
		if ok {
			t.Fatalf("%s should be gone (lang=en)", id)
		}
	}
}

func testFullText(t *testing.T, newAdapter NewAdapter, supported bool) {
	v, cleanup := newAdapter(t, "contract_fts")
	defer cleanup()
	ctx := context.Background()
	docs, embs := fixtureDocs()
	must(t, v.Create(ctx, docs, embs), "Create")

	hits, err := v.SearchByFullText(ctx, "alpha", ragentities.SearchOptions{TopK: 3})
	if supported {
		must(t, err, "SearchByFullText")
		if len(hits) == 0 {
			t.Fatal("expected at least one hit for 'alpha'")
		}
	} else {
		if err == nil {
			t.Fatal("expected ErrFullTextUnsupported, got nil error")
		}
		if !errors.Is(err, ErrFullTextUnsupported) {
			t.Logf("note: adapter returned a different sentinel for unsupported full-text: %v", err)
		}
	}
}

func testDeleteCollection(t *testing.T, newAdapter NewAdapter) {
	v, cleanup := newAdapter(t, "contract_drop")
	defer cleanup()
	ctx := context.Background()
	docs, embs := fixtureDocs()
	must(t, v.Create(ctx, docs, embs), "Create")

	must(t, v.Delete(ctx), "Delete collection")

	ok, err := v.TextExists(ctx, "a")
	// Either the collection is gone (error is acceptable) or it's empty.
	if err == nil && ok {
		t.Fatal("collection still contains data after Delete")
	}
}

// ---------- helpers ----------

func must(t *testing.T, err error, op string) {
	t.Helper()
	if err != nil {
		t.Fatalf("%s: %v", op, err)
	}
}

func idsOf(hits []*ragentities.SearchResult) []any {
	out := make([]any, len(hits))
	for i, h := range hits {
		out[i] = h.Document.Metadata["doc_id"]
	}
	return out
}

func scoresOf(hits []*ragentities.SearchResult) []float32 {
	out := make([]float32, len(hits))
	for i, h := range hits {
		out[i] = h.Score
	}
	return out
}
