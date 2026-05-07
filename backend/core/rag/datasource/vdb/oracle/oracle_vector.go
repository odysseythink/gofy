// Package oracle is a stub adapter for Oracle 23ai AI Vector Search.
//
// Requires the godror SDK (github.com/godror/godror) plus an Oracle Instant
// Client at runtime. Neither is vendored. The Factory registers so config
// plumbing works; New returns a descriptive error until the SDK is added.
package oracle

import (
	"context"
	"errors"

	vdb "github.com/odysseythink/gofy/backend/core/rag/datasource/vdb"
	ragentities "github.com/odysseythink/gofy/backend/entities/rag"
	vectorenumtypes "github.com/odysseythink/gofy/backend/enum_types/rag/vector"
)

func init() { vdb.Register(vectorenumtypes.Vector_ORACLE, &Factory{}) }

type Factory struct{}

func (Factory) New(ctx context.Context, cfg vdb.AdapterConfig) (ragentities.IVector, error) {
	return nil, errors.New("oracle: adapter not yet implemented. Vendor github.com/godror/godror (requires Oracle Instant Client) to enable")
}
