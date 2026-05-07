// Package lindorm is the AliCloud Lindorm IVector adapter. Lindorm offers a
// search-engine mode that speaks OpenSearch-compatible REST; we delegate to
// the elasticsearch adapter with lindorm-specific config.
package lindorm

import (
	"context"
	"errors"

	"github.com/odysseythink/confy"

	vdb "github.com/odysseythink/gofy/backend/core/rag/datasource/vdb"
	"github.com/odysseythink/gofy/backend/core/rag/datasource/vdb/elasticsearch"
	ragentities "github.com/odysseythink/gofy/backend/entities/rag"
	vectorenumtypes "github.com/odysseythink/gofy/backend/enum_types/rag/vector"
)

func init() { vdb.Register(vectorenumtypes.Vector_LINDORM, &Factory{}) }

type Factory struct{}

func (Factory) New(ctx context.Context, cfg vdb.AdapterConfig) (ragentities.IVector, error) {
	esCfg := elasticsearch.Config{
		Endpoint: confy.GetWithDefault[string]("lindorm.endpoint", ""),
		Username: confy.GetWithDefault[string]("lindorm.username", ""),
		Password: confy.GetWithDefault[string]("lindorm.password", ""),
	}
	if esCfg.Endpoint == "" {
		return nil, errors.New("lindorm: endpoint is required")
	}
	v, err := elasticsearch.Open(ctx, esCfg, cfg.CollectionName)
	if err != nil {
		return nil, err
	}
	return v, nil
}
