// Package matrixone is the MatrixOne IVector adapter. MatrixOne speaks MySQL
// protocol and supports a native vector type; we reuse the tidb_vector
// transport.
package matrixone

import (
	"context"

	"github.com/odysseythink/confy"

	vdb "mlib.com/gofy/server/core/rag/datasource/vdb"
	"mlib.com/gofy/server/core/rag/datasource/vdb/tidb_vector"
	ragentities "mlib.com/gofy/server/entities/rag"
	vectorenumtypes "mlib.com/gofy/server/enum_types/rag/vector"
)

func init() { vdb.Register(vectorenumtypes.Vector_MATRIXONE, &Factory{}) }

type Factory struct{}

func (Factory) New(ctx context.Context, cfg vdb.AdapterConfig) (ragentities.IVector, error) {
	c := tidb_vector.Config{
		Host:     confy.GetWithDefault[string]("matrixone.host", "127.0.0.1"),
		Port:     confy.GetWithDefault[int]("matrixone.port", 6001),
		User:     confy.GetWithDefault[string]("matrixone.user", "dump"),
		Password: confy.GetWithDefault[string]("matrixone.password", "111"),
		Database: confy.GetWithDefault[string]("matrixone.database", "test"),
	}
	return tidb_vector.Open(ctx, c, cfg.CollectionName)
}
