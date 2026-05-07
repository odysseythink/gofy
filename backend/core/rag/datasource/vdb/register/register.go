// Package register imports every VDB adapter purely for side-effects so a
// caller that imports this package picks up all 26 IVector Factory
// registrations in one line:
//
//	import _ "github.com/odysseythink/gofy/backend/core/rag/datasource/vdb/register"
package register

import (
	// SQL family (Phase 1)
	_ "github.com/odysseythink/gofy/backend/core/rag/datasource/vdb/analyticdb"
	_ "github.com/odysseythink/gofy/backend/core/rag/datasource/vdb/opengauss"
	_ "github.com/odysseythink/gofy/backend/core/rag/datasource/vdb/pgvecto_rs"
	_ "github.com/odysseythink/gofy/backend/core/rag/datasource/vdb/pgvector"
	_ "github.com/odysseythink/gofy/backend/core/rag/datasource/vdb/pyvastbase"
	_ "github.com/odysseythink/gofy/backend/core/rag/datasource/vdb/relyt"
	_ "github.com/odysseythink/gofy/backend/core/rag/datasource/vdb/tidb_vector"

	// Mature SDK family (Phase 2)
	_ "github.com/odysseythink/gofy/backend/core/rag/datasource/vdb/chroma"
	_ "github.com/odysseythink/gofy/backend/core/rag/datasource/vdb/couchbase"
	_ "github.com/odysseythink/gofy/backend/core/rag/datasource/vdb/elasticsearch"
	_ "github.com/odysseythink/gofy/backend/core/rag/datasource/vdb/milvus"
	_ "github.com/odysseythink/gofy/backend/core/rag/datasource/vdb/opensearch"
	_ "github.com/odysseythink/gofy/backend/core/rag/datasource/vdb/oracle"
	_ "github.com/odysseythink/gofy/backend/core/rag/datasource/vdb/qdrant"
	_ "github.com/odysseythink/gofy/backend/core/rag/datasource/vdb/upstash"
	_ "github.com/odysseythink/gofy/backend/core/rag/datasource/vdb/weaviate"

	// Domestic cloud family (Phase 3)
	_ "github.com/odysseythink/gofy/backend/core/rag/datasource/vdb/baidu"
	_ "github.com/odysseythink/gofy/backend/core/rag/datasource/vdb/huawei"
	_ "github.com/odysseythink/gofy/backend/core/rag/datasource/vdb/lindorm"
	_ "github.com/odysseythink/gofy/backend/core/rag/datasource/vdb/tablestore"
	_ "github.com/odysseythink/gofy/backend/core/rag/datasource/vdb/tencent"
	_ "github.com/odysseythink/gofy/backend/core/rag/datasource/vdb/vikingdb"

	// Misc / composite (Phase 4)
	_ "github.com/odysseythink/gofy/backend/core/rag/datasource/vdb/matrixone"
	_ "github.com/odysseythink/gofy/backend/core/rag/datasource/vdb/myscale"
	_ "github.com/odysseythink/gofy/backend/core/rag/datasource/vdb/oceanbase"
	_ "github.com/odysseythink/gofy/backend/core/rag/datasource/vdb/tidb_on_qdrant"
)
