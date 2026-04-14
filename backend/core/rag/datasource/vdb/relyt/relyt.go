// Package relyt is the Relyt (Hashdata / Greenplum-based managed PG) IVector adapter.
// Relyt is wire-compatible with PostgreSQL and uses the pgvector extension.
package relyt

import (
	"context"
	"fmt"

	"github.com/odysseythink/confy"

	vdb "mlib.com/gofy/server/core/rag/datasource/vdb"
	"mlib.com/gofy/server/core/rag/datasource/vdb/sqlvec"
	ragentities "mlib.com/gofy/server/entities/rag"
	vectorenumtypes "mlib.com/gofy/server/enum_types/rag/vector"
)

func init() { vdb.Register(vectorenumtypes.Vector_RELYT, &Factory{}) }

type Config = sqlvec.Config
type Factory struct{}

func (Factory) New(ctx context.Context, cfg vdb.AdapterConfig) (ragentities.IVector, error) {
	return Open(ctx, loadConfig(), cfg.CollectionName)
}

func loadConfig() Config {
	return Config{
		Host:          confy.GetWithDefault[string]("relyt.host", "localhost"),
		Port:          confy.GetWithDefault[int]("relyt.port", 5432),
		User:          confy.GetWithDefault[string]("relyt.user", "postgres"),
		Password:      confy.GetWithDefault[string]("relyt.password", ""),
		Database:      confy.GetWithDefault[string]("relyt.database", "postgres"),
		MinConnection: int32(confy.GetWithDefault[int]("relyt.min-connection", 1)),
		MaxConnection: int32(confy.GetWithDefault[int]("relyt.max-connection", 10)),
	}
}

func Open(ctx context.Context, cfg Config, collection string) (*sqlvec.Vector, error) {
	return sqlvec.Open(ctx, cfg, collection, dialect{})
}

type dialect struct{}

func (dialect) Name() string               { return string(vectorenumtypes.Vector_RELYT) }
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
	) DISTRIBUTED BY (id)`, table, dim)
}

func (dialect) CreateIndexSQL(table, indexHash string, dim int) string {
	if dim > 2000 {
		return ""
	}
	return fmt.Sprintf(`CREATE INDEX IF NOT EXISTS relyt_vec_idx_%s ON %s USING hnsw (embedding vector_cosine_ops) WITH (m = 16, ef_construction = 64)`,
		indexHash, table)
}
