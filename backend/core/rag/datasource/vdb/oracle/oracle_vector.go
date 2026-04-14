// Package oracle is a stub adapter for Oracle 23ai AI Vector Search.
//
// Requires the godror SDK (github.com/godror/godror) plus an Oracle Instant
// Client at runtime. Neither is vendored. The Factory registers so config
// plumbing works; New returns a descriptive error until the SDK is added.
package oracle

import (
	"context"
	"errors"

	vdb "mlib.com/gofy/server/core/rag/datasource/vdb"
	ragentities "mlib.com/gofy/server/entities/rag"
	vectorenumtypes "mlib.com/gofy/server/enum_types/rag/vector"
)

func init() { vdb.Register(vectorenumtypes.Vector_ORACLE, &Factory{}) }

type Factory struct{}

func (Factory) New(ctx context.Context, cfg vdb.AdapterConfig) (ragentities.IVector, error) {
	return nil, errors.New("oracle: adapter not yet implemented. Vendor github.com/godror/godror (requires Oracle Instant Client) to enable")
}
