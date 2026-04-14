// Package opengauss is the openGauss IVector adapter. openGauss is a
// PostgreSQL fork with a built-in `floatvector` / `vector` datatype.
package opengauss

import (
	"context"
	"fmt"

	"github.com/odysseythink/confy"

	vdb "mlib.com/gofy/server/core/rag/datasource/vdb"
	"mlib.com/gofy/server/core/rag/datasource/vdb/sqlvec"
	ragentities "mlib.com/gofy/server/entities/rag"
	vectorenumtypes "mlib.com/gofy/server/enum_types/rag/vector"
)

func init() {
	vdb.Register(vectorenumtypes.Vector_OPENGAUSS, &Factory{})
}

type Config = sqlvec.Config
type Factory struct{}

func (Factory) New(ctx context.Context, cfg vdb.AdapterConfig) (ragentities.IVector, error) {
	return Open(ctx, loadConfig(), cfg.CollectionName)
}

func loadConfig() Config {
	return Config{
		Host:          confy.GetWithDefault[string]("opengauss.host", "localhost"),
		Port:          confy.GetWithDefault[int]("opengauss.port", 6000),
		User:          confy.GetWithDefault[string]("opengauss.user", "gaussdb"),
		Password:      confy.GetWithDefault[string]("opengauss.password", ""),
		Database:      confy.GetWithDefault[string]("opengauss.database", "postgres"),
		MinConnection: int32(confy.GetWithDefault[int]("opengauss.min-connection", 1)),
		MaxConnection: int32(confy.GetWithDefault[int]("opengauss.max-connection", 10)),
	}
}

func Open(ctx context.Context, cfg Config, collection string) (*sqlvec.Vector, error) {
	return sqlvec.Open(ctx, cfg, collection, dialect{})
}

type dialect struct{}

func (dialect) Name() string               { return string(vectorenumtypes.Vector_OPENGAUSS) }
func (dialect) CreateExtensionSQL() string { return "" } // built-in
func (dialect) VectorCast() string         { return "floatvector" }
func (dialect) DistanceExpr() string       { return `embedding <=> $1::floatvector` }
func (dialect) FormatVector(v []float32) string {
	return sqlvec.FormatPGVector(v)
}

func (dialect) CreateTableSQL(table string, dim int) string {
	return fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		id VARCHAR(64) PRIMARY KEY,
		text TEXT NOT NULL,
		meta JSONB NOT NULL,
		embedding floatvector(%d) NOT NULL
	)`, table, dim)
}

func (dialect) CreateIndexSQL(table, indexHash string, dim int) string {
	return fmt.Sprintf(`CREATE INDEX IF NOT EXISTS og_vec_idx_%s ON %s USING hnsw (embedding cosine_ops)`,
		indexHash, table)
}
