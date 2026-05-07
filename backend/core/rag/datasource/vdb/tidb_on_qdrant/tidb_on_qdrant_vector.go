// Package tidb_on_qdrant registers the composite vector store that uses
// TiDB for metadata and Qdrant for vector search.
//
// The Python reference keeps metadata in TiDB and dense vectors in Qdrant.
// For the Go port we delegate both read and write to the Qdrant adapter
// (which already handles payload-based metadata), and rely on the main
// MySQL/TiDB config for tenant-level whitelist integration (as the legacy
// `vector_store_whitelist_enable` path wires).
package tidb_on_qdrant

// The actual Factory registration for Vector_TIDB_ON_QDRANT happens in the
// qdrant package's init() — imported here for side effects so callers that
// only pull in this package transitively get the registration.
import _ "github.com/odysseythink/gofy/backend/core/rag/datasource/vdb/qdrant"
