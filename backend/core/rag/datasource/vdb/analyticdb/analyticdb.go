// Package analyticdb is the AliCloud AnalyticDB for PostgreSQL IVector adapter.
// AnalyticDB-PG is wire-compatible with PostgreSQL and ships the vector extension.
package analyticdb

import (
	"context"
	"fmt"

	"github.com/odysseythink/confy"

	vdb "github.com/odysseythink/gofy/backend/core/rag/datasource/vdb"
	"github.com/odysseythink/gofy/backend/core/rag/datasource/vdb/sqlvec"
	ragentities "github.com/odysseythink/gofy/backend/entities/rag"
	vectorenumtypes "github.com/odysseythink/gofy/backend/enum_types/rag/vector"
)

func init() { vdb.Register(vectorenumtypes.Vector_ANALYTICDB, &Factory{}) }

type Config = sqlvec.Config
type Factory struct{}

func (Factory) New(ctx context.Context, cfg vdb.AdapterConfig) (ragentities.IVector, error) {
	return Open(ctx, loadConfig(), cfg.CollectionName)
}

func loadConfig() Config {
	return Config{
		Host:          confy.GetWithDefault[string]("analyticdb.host", "localhost"),
		Port:          confy.GetWithDefault[int]("analyticdb.port", 5432),
		User:          confy.GetWithDefault[string]("analyticdb.user", "postgres"),
		Password:      confy.GetWithDefault[string]("analyticdb.password", ""),
		Database:      confy.GetWithDefault[string]("analyticdb.database", "postgres"),
		MinConnection: int32(confy.GetWithDefault[int]("analyticdb.min-connection", 1)),
		MaxConnection: int32(confy.GetWithDefault[int]("analyticdb.max-connection", 10)),
		DSNQuery:      "sslmode=disable",
	}
}

func Open(ctx context.Context, cfg Config, collection string) (*sqlvec.Vector, error) {
	return sqlvec.Open(ctx, cfg, collection, dialect{})
}

type dialect struct{}

func (dialect) Name() string               { return string(vectorenumtypes.Vector_ANALYTICDB) }
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
	return fmt.Sprintf(`CREATE INDEX IF NOT EXISTS adb_vec_idx_%s ON %s USING hnsw (embedding vector_cosine_ops) WITH (m = 16, ef_construction = 64)`,
		indexHash, table)
}
