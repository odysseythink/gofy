package pgvector_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"mlib.com/gofy/server/core/rag/datasource/vdb/pgvector"
	"mlib.com/gofy/server/core/rag/datasource/vdb/vdbtest"
	ragentities "mlib.com/gofy/server/entities/rag"
)

func TestPgvectorContract(t *testing.T) {
	host := os.Getenv("GOFY_TEST_PGVECTOR_HOST")
	if host == "" {
		t.Skip("GOFY_TEST_PGVECTOR_HOST not set")
	}
	cfg := pgvector.Config{
		Host:          host,
		Port:          envIntOr("GOFY_TEST_PGVECTOR_PORT", 5433),
		User:          envOr("GOFY_TEST_PGVECTOR_USER", "postgres"),
		Password:      envOr("GOFY_TEST_PGVECTOR_PASSWORD", "test"),
		Database:      envOr("GOFY_TEST_PGVECTOR_DATABASE", "postgres"),
		MinConnection: 1,
		MaxConnection: 4,
	}

	vdbtest.RunContractSuite(t, func(t *testing.T, name string) (ragentities.IVector, func()) {
		v, err := pgvector.Open(context.Background(), cfg, name+randSuffix())
		if err != nil {
			t.Fatalf("open: %v", err)
		}
		return v, func() {
			_ = v.Delete(context.Background())
			v.Close()
		}
	}, vdbtest.Capabilities{
		MetadataFilter: true,
		MetadataDelete: true,
		FullText:       true,
	})
}

func envOr(k, dflt string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return dflt
}
func envIntOr(k string, dflt int) int {
	var out int
	if _, err := fmt.Sscanf(os.Getenv(k), "%d", &out); err == nil && out > 0 {
		return out
	}
	return dflt
}

var suffixCounter int

func randSuffix() string { suffixCounter++; return fmt.Sprintf("_%d", suffixCounter) }
