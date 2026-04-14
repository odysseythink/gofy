// Package tablestore is a stub adapter for AliCloud Tablestore (OTS) with
// vector search.
//
// Requires github.com/aliyun/aliyun-tablestore-go-sdk which is not vendored.
// The Factory registers so configuration is preserved; New returns a
// descriptive error until the SDK is added.
package tablestore

import (
	"context"
	"errors"

	vdb "mlib.com/gofy/server/core/rag/datasource/vdb"
	ragentities "mlib.com/gofy/server/entities/rag"
	vectorenumtypes "mlib.com/gofy/server/enum_types/rag/vector"
)

func init() { vdb.Register(vectorenumtypes.Vector_TABLESTORE, &Factory{}) }

type Factory struct{}

func (Factory) New(ctx context.Context, cfg vdb.AdapterConfig) (ragentities.IVector, error) {
	return nil, errors.New("tablestore: adapter not yet implemented. Vendor github.com/aliyun/aliyun-tablestore-go-sdk to enable")
}
