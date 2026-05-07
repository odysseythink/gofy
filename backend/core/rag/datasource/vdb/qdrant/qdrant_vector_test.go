package qdrant_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/odysseythink/gofy/backend/core/rag/datasource/vdb/qdrant"
	"github.com/odysseythink/gofy/backend/core/rag/datasource/vdb/vdbtest"
	ragentities "github.com/odysseythink/gofy/backend/entities/rag"
)

func TestQdrantContract(t *testing.T) {
	endpoint := os.Getenv("GOFY_TEST_QDRANT_ENDPOINT")
	if endpoint == "" {
		t.Skip("GOFY_TEST_QDRANT_ENDPOINT not set")
	}
	cfg := qdrant.Config{Endpoint: endpoint, APIKey: os.Getenv("GOFY_TEST_QDRANT_API_KEY")}

	vdbtest.RunContractSuite(t, func(t *testing.T, name string) (ragentities.IVector, func()) {
		v, err := qdrant.Open(context.Background(), cfg, name+suffix())
		if err != nil {
			t.Fatalf("open: %v", err)
		}
		return v, func() { _ = v.Delete(context.Background()) }
	}, vdbtest.Capabilities{
		MetadataFilter: true,
		MetadataDelete: true,
		FullText:       false,
	})
}

var ctr int

func suffix() string { ctr++; return fmt.Sprintf("_%d", ctr) }
