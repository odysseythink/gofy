// Package pyvastbase is the VastBase (PostgreSQL-compatible) IVector adapter.
// Named pyvastbase to match the Python reference folder naming.
package pyvastbase

import (
	"context"
	"fmt"

	"github.com/odysseythink/confy"

	vdb "mlib.com/gofy/server/core/rag/datasource/vdb"
	"mlib.com/gofy/server/core/rag/datasource/vdb/sqlvec"
	ragentities "mlib.com/gofy/server/entities/rag"
	vectorenumtypes "mlib.com/gofy/server/enum_types/rag/vector"
)

func init() { vdb.Register(vectorenumtypes.Vector_VASTBASE, &Factory{}) }

type Config = sqlvec.Config
type Factory struct{}

func (Factory) New(ctx context.Context, cfg vdb.AdapterConfig) (ragentities.IVector, error) {
	return Open(ctx, loadConfig(), cfg.CollectionName)
}

func loadConfig() Config {
	return Config{
		Host:          confy.GetWithDefault[string]("vastbase.host", "localhost"),
		Port:          confy.GetWithDefault[int]("vastbase.port", 5432),
		User:          confy.GetWithDefault[string]("vastbase.user", "postgres"),
		Password:      confy.GetWithDefault[string]("vastbase.password", ""),
		Database:      confy.GetWithDefault[string]("vastbase.database", "postgres"),
		MinConnection: int32(confy.GetWithDefault[int]("vastbase.min-connection", 1)),
		MaxConnection: int32(confy.GetWithDefault[int]("vastbase.max-connection", 10)),
	}
}

func Open(ctx context.Context, cfg Config, collection string) (*sqlvec.Vector, error) {
	return sqlvec.Open(ctx, cfg, collection, dialect{})
}

type dialect struct{}

func (dialect) Name() string               { return string(vectorenumtypes.Vector_VASTBASE) }
func (dialect) CreateExtensionSQL() string { return `CREATE EXTENSION IF NOT EXISTS floatvector` }
func (dialect) VectorCast() string         { return "floatvector" }
func (dialect) DistanceExpr() string       { return `embedding <=> $1::floatvector` }
func (dialect) FormatVector(v []float32) string {
	return sqlvec.FormatPGVector(v)
}

func (dialect) CreateTableSQL(table string, dim int) string {
	return fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		id TEXT PRIMARY KEY,
		text TEXT NOT NULL,
		meta JSONB NOT NULL,
		embedding floatvector(%d) NOT NULL
	)`, table, dim)
}

func (dialect) CreateIndexSQL(table, indexHash string, dim int) string {
	return fmt.Sprintf(`CREATE INDEX IF NOT EXISTS vb_vec_idx_%s ON %s USING hnsw (embedding floatvector_cosine_ops)`,
		indexHash, table)
}
