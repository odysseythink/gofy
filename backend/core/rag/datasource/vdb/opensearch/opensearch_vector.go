// Package opensearch is a stub that points to the elasticsearch adapter,
// since OpenSearch's REST shape matches ES closely enough for our use. The
// actual Factory registration lives in the elasticsearch package.
package opensearch

// Side-effect import only — the elasticsearch package registers
// Vector_OPENSEARCH in its init().
import _ "github.com/odysseythink/gofy/backend/core/rag/datasource/vdb/elasticsearch"
