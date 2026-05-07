// Package pgvector is the PostgreSQL + pgvector IVector adapter.
package pgvector

import (
	"context"
	"fmt"

	"github.com/odysseythink/confy"

	vdb "github.com/odysseythink/gofy/backend/core/rag/datasource/vdb"
	"github.com/odysseythink/gofy/backend/core/rag/datasource/vdb/sqlvec"
	ragentities "github.com/odysseythink/gofy/backend/entities/rag"
	vectorenumtypes "github.com/odysseythink/gofy/backend/enum_types/rag/vector"
)

func init() {
	vdb.Register(vectorenumtypes.Vector_PGVECTOR, &Factory{})
}

// Config is a convenience re-export that tests use.
type Config = sqlvec.Config

type Factory struct{}

func (Factory) New(ctx context.Context, cfg vdb.AdapterConfig) (ragentities.IVector, error) {
	c := loadConfig()
	return Open(ctx, c, cfg.CollectionName)
}

func loadConfig() Config {
	return Config{
		Host:          confy.GetWithDefault[string]("pgvector.host", "localhost"),
		Port:          confy.GetWithDefault[int]("pgvector.port", 5432),
		User:          confy.GetWithDefault[string]("pgvector.user", "postgres"),
		Password:      confy.GetWithDefault[string]("pgvector.password", ""),
		Database:      confy.GetWithDefault[string]("pgvector.database", "postgres"),
		MinConnection: int32(confy.GetWithDefault[int]("pgvector.min-connection", 1)),
		MaxConnection: int32(confy.GetWithDefault[int]("pgvector.max-connection", 10)),
	}
}

// Open is exposed for tests.
func Open(ctx context.Context, cfg Config, collection string) (*sqlvec.Vector, error) {
	return sqlvec.Open(ctx, cfg, collection, dialect{})
}

// dialect implements sqlvec.Dialect for pgvector.
type dialect struct{}

func (dialect) Name() string               { return string(vectorenumtypes.Vector_PGVECTOR) }
func (dialect) CreateExtensionSQL() string { return `CREATE EXTENSION IF NOT EXISTS vector` }
func (dialect) VectorCast() string         { return "vector" }
func (dialect) DistanceExpr() string       { return `embedding <=> $1::vector` }
func (dialect) FormatVector(v []float32) string {
	return sqlvec.FormatPGVector(v)
}

func (dialect) CreateTableSQL(table string, dim int) string {
	return fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		id TEXT PRIMARY KEY,
		text TEXT NOT NULL,
		meta JSONB NOT NULL,
		embedding vector(%d) NOT NULL
	)`, table, dim)
}

func (dialect) CreateIndexSQL(table, indexHash string, dim int) string {
	if dim > 2000 {
		return "" // pgvector hnsw tops out at 2000 dims
	}
	return fmt.Sprintf(`CREATE INDEX IF NOT EXISTS embedding_cosine_v1_idx_%s ON %s USING hnsw (embedding vector_cosine_ops) WITH (m = 16, ef_construction = 64)`,
		indexHash, table)
}
