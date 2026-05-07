// Package pgvecto_rs is the PostgreSQL + pgvecto.rs IVector adapter.
//
// pgvecto.rs exposes a similar API to pgvector but uses the "vectors"
// extension, `vectors.vector` type, and different index ops.
package pgvecto_rs

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
	vdb.Register(vectorenumtypes.Vector_PGVECTO_RS, &Factory{})
}

type Config = sqlvec.Config

type Factory struct{}

func (Factory) New(ctx context.Context, cfg vdb.AdapterConfig) (ragentities.IVector, error) {
	return Open(ctx, loadConfig(), cfg.CollectionName)
}

func loadConfig() Config {
	return Config{
		Host:          confy.GetWithDefault[string]("pgvecto-rs.host", "localhost"),
		Port:          confy.GetWithDefault[int]("pgvecto-rs.port", 5432),
		User:          confy.GetWithDefault[string]("pgvecto-rs.user", "postgres"),
		Password:      confy.GetWithDefault[string]("pgvecto-rs.password", ""),
		Database:      confy.GetWithDefault[string]("pgvecto-rs.database", "postgres"),
		MinConnection: int32(confy.GetWithDefault[int]("pgvecto-rs.min-connection", 1)),
		MaxConnection: int32(confy.GetWithDefault[int]("pgvecto-rs.max-connection", 10)),
	}
}

func Open(ctx context.Context, cfg Config, collection string) (*sqlvec.Vector, error) {
	return sqlvec.Open(ctx, cfg, collection, dialect{})
}

type dialect struct{}

func (dialect) Name() string               { return string(vectorenumtypes.Vector_PGVECTO_RS) }
func (dialect) CreateExtensionSQL() string { return `CREATE EXTENSION IF NOT EXISTS vectors` }
func (dialect) VectorCast() string         { return "vectors.vector" }
func (dialect) DistanceExpr() string       { return `embedding <=> $1::vectors.vector` }
func (dialect) FormatVector(v []float32) string {
	return sqlvec.FormatPGVector(v)
}

func (dialect) CreateTableSQL(table string, dim int) string {
	return fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		id TEXT PRIMARY KEY,
		text TEXT NOT NULL,
		meta JSONB NOT NULL,
		embedding vectors.vector(%d) NOT NULL
	)`, table, dim)
}

func (dialect) CreateIndexSQL(table, indexHash string, dim int) string {
	return fmt.Sprintf(`CREATE INDEX IF NOT EXISTS embedding_cosine_idx_%s ON %s USING vectors (embedding vectors.vector_cos_ops)`,
		indexHash, table)
}
