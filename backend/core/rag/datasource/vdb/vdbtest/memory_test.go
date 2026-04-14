package vdbtest

import (
	"context"
	"testing"

	ragentities "mlib.com/gofy/server/entities/rag"
)

// Apply the contract suite to the reference in-memory implementation.
// Failing this means the harness itself is broken.
func TestContractAgainstMemory(t *testing.T) {
	RunContractSuite(t, func(t *testing.T, collection string) (ragentities.IVector, func()) {
		m := NewMemory(collection)
		return m, func() { _ = m.Delete(context.Background()) }
	}, Capabilities{
		MetadataFilter: true,
		MetadataDelete: true,
		FullText:       false,
	})
}
