package register_test

import (
	"testing"

	vdb "mlib.com/gofy/server/core/rag/datasource/vdb"
	_ "mlib.com/gofy/server/core/rag/datasource/vdb/register"
	vectorenumtypes "mlib.com/gofy/server/enum_types/rag/vector"
)

// TestAllRegistered guards against someone adding a new VectorType enum value
// without also registering an adapter Factory for it. Fails loudly with the
// list of missing types.
func TestAllRegistered(t *testing.T) {
	expected := []vectorenumtypes.VectorType{
		vectorenumtypes.Vector_ANALYTICDB,
		vectorenumtypes.Vector_BAIDU,
		vectorenumtypes.Vector_CHROMA,
		vectorenumtypes.Vector_COUCHBASE,
		vectorenumtypes.Vector_ELASTICSEARCH,
		vectorenumtypes.Vector_ELASTICSEARCH_JA,
		vectorenumtypes.Vector_HUAWEI_CLOUD,
		vectorenumtypes.Vector_LINDORM,
		vectorenumtypes.Vector_MATRIXONE,
		vectorenumtypes.Vector_MILVUS,
		vectorenumtypes.Vector_MYSCALE,
		vectorenumtypes.Vector_OCEANBASE,
		vectorenumtypes.Vector_OPENGAUSS,
		vectorenumtypes.Vector_OPENSEARCH,
		vectorenumtypes.Vector_ORACLE,
		vectorenumtypes.Vector_PGVECTOR,
		vectorenumtypes.Vector_PGVECTO_RS,
		vectorenumtypes.Vector_QDRANT,
		vectorenumtypes.Vector_RELYT,
		vectorenumtypes.Vector_TABLESTORE,
		vectorenumtypes.Vector_TENCENT,
		vectorenumtypes.Vector_TIDB_ON_QDRANT,
		vectorenumtypes.Vector_TIDB_VECTOR,
		vectorenumtypes.Vector_UPSTASH,
		vectorenumtypes.Vector_VASTBASE,
		vectorenumtypes.Vector_VIKINGDB,
	}
	for _, t2 := range expected {
		if _, ok := vdb.Get(t2); !ok {
			t.Errorf("vector type %q not registered", t2)
		}
	}
	if t.Failed() {
		t.Logf("currently registered: %v", vdb.Registered())
	}
}
