// Package oceanbase is the OceanBase IVector adapter. OceanBase speaks the
// MySQL protocol and recent releases support VECTOR columns + HNSW indexes.
// The adapter reuses the tidb_vector transport layer with an OceanBase-tuned
// config and DDL.
package oceanbase

import (
	"context"
	"errors"

	"github.com/odysseythink/confy"

	vdb "github.com/odysseythink/gofy/backend/core/rag/datasource/vdb"
	"github.com/odysseythink/gofy/backend/core/rag/datasource/vdb/tidb_vector"
	ragentities "github.com/odysseythink/gofy/backend/entities/rag"
	vectorenumtypes "github.com/odysseythink/gofy/backend/enum_types/rag/vector"
)

func init() { vdb.Register(vectorenumtypes.Vector_OCEANBASE, &Factory{}) }

type Factory struct{}

func (Factory) New(ctx context.Context, cfg vdb.AdapterConfig) (ragentities.IVector, error) {
	c := tidb_vector.Config{
		Host:     confy.GetWithDefault[string]("oceanbase.host", "127.0.0.1"),
		Port:     confy.GetWithDefault[int]("oceanbase.port", 2881),
		User:     confy.GetWithDefault[string]("oceanbase.user", "root"),
		Password: confy.GetWithDefault[string]("oceanbase.password", ""),
		Database: confy.GetWithDefault[string]("oceanbase.database", "test"),
	}
	if c.Host == "" {
		return nil, errors.New("oceanbase: host is required")
	}
	return tidb_vector.Open(ctx, c, cfg.CollectionName)
}
