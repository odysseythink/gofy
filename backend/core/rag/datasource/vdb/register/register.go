// Package register imports every VDB adapter purely for side-effects so a
// caller that imports this package picks up all 26 IVector Factory
// registrations in one line:
//
//	import _ "mlib.com/gofy/server/core/rag/datasource/vdb/register"
package register

import (
	// SQL family (Phase 1)
	_ "mlib.com/gofy/server/core/rag/datasource/vdb/analyticdb"
	_ "mlib.com/gofy/server/core/rag/datasource/vdb/opengauss"
	_ "mlib.com/gofy/server/core/rag/datasource/vdb/pgvecto_rs"
	_ "mlib.com/gofy/server/core/rag/datasource/vdb/pgvector"
	_ "mlib.com/gofy/server/core/rag/datasource/vdb/pyvastbase"
	_ "mlib.com/gofy/server/core/rag/datasource/vdb/relyt"
	_ "mlib.com/gofy/server/core/rag/datasource/vdb/tidb_vector"

	// Mature SDK family (Phase 2)
	_ "mlib.com/gofy/server/core/rag/datasource/vdb/chroma"
	_ "mlib.com/gofy/server/core/rag/datasource/vdb/couchbase"
	_ "mlib.com/gofy/server/core/rag/datasource/vdb/elasticsearch"
	_ "mlib.com/gofy/server/core/rag/datasource/vdb/milvus"
	_ "mlib.com/gofy/server/core/rag/datasource/vdb/opensearch"
	_ "mlib.com/gofy/server/core/rag/datasource/vdb/oracle"
	_ "mlib.com/gofy/server/core/rag/datasource/vdb/qdrant"
	_ "mlib.com/gofy/server/core/rag/datasource/vdb/upstash"
	_ "mlib.com/gofy/server/core/rag/datasource/vdb/weaviate"

	// Domestic cloud family (Phase 3)
	_ "mlib.com/gofy/server/core/rag/datasource/vdb/baidu"
	_ "mlib.com/gofy/server/core/rag/datasource/vdb/huawei"
	_ "mlib.com/gofy/server/core/rag/datasource/vdb/lindorm"
	_ "mlib.com/gofy/server/core/rag/datasource/vdb/tablestore"
	_ "mlib.com/gofy/server/core/rag/datasource/vdb/tencent"
	_ "mlib.com/gofy/server/core/rag/datasource/vdb/vikingdb"

	// Misc / composite (Phase 4)
	_ "mlib.com/gofy/server/core/rag/datasource/vdb/matrixone"
	_ "mlib.com/gofy/server/core/rag/datasource/vdb/myscale"
	_ "mlib.com/gofy/server/core/rag/datasource/vdb/oceanbase"
	_ "mlib.com/gofy/server/core/rag/datasource/vdb/tidb_on_qdrant"
)
