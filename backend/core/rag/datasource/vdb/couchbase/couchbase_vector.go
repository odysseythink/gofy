// Package couchbase is a stub adapter for Couchbase vector search.
//
// Requires the gocb v2 SDK (github.com/couchbase/gocb/v2) which is not
// vendored in this repo yet. The Factory registers with the IVector registry
// so configuration plumbing works; New returns an error directing the user to
// vendor the SDK before enabling this backend.
package couchbase

import (
	"context"
	"errors"

	vdb "github.com/odysseythink/gofy/backend/core/rag/datasource/vdb"
	ragentities "github.com/odysseythink/gofy/backend/entities/rag"
	vectorenumtypes "github.com/odysseythink/gofy/backend/enum_types/rag/vector"
)

func init() { vdb.Register(vectorenumtypes.Vector_COUCHBASE, &Factory{}) }

type Factory struct{}

func (Factory) New(ctx context.Context, cfg vdb.AdapterConfig) (ragentities.IVector, error) {
	return nil, errors.New("couchbase: adapter not yet implemented. Vendor github.com/couchbase/gocb/v2 and implement NewFromGocb to enable")
}
