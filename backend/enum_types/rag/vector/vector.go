package vector

type VectorType string

const (
	Vector_ANALYTICDB VectorType = "analyticdb"
	Vector_CHROMA     VectorType = "chroma"
	Vector_MILVUS     VectorType = "milvus"
	Vector_MYSCALE    VectorType = "myscale"
	Vector_PGVECTOR   VectorType = "pgvector"
	Vector_VASTBASE   VectorType = "vastbase"
	Vector_PGVECTO_RS VectorType = "pgvecto-rs"

	Vector_QDRANT           VectorType = "qdrant"
	Vector_RELYT            VectorType = "relyt"
	Vector_TIDB_VECTOR      VectorType = "tidb_vector"
	Vector_WEAVIATE         VectorType = "weaviate"
	Vector_OPENSEARCH       VectorType = "opensearch"
	Vector_TENCENT          VectorType = "tencent"
	Vector_ORACLE           VectorType = "oracle"
	Vector_ELASTICSEARCH    VectorType = "elasticsearch"
	Vector_ELASTICSEARCH_JA VectorType = "elasticsearch-ja"
	Vector_LINDORM          VectorType = "lindorm"
	Vector_COUCHBASE        VectorType = "couchbase"
	Vector_BAIDU            VectorType = "baidu"
	Vector_VIKINGDB         VectorType = "vikingdb"
	Vector_UPSTASH          VectorType = "upstash"
	Vector_TIDB_ON_QDRANT   VectorType = "tidb_on_qdrant"
	Vector_OCEANBASE        VectorType = "oceanbase"
	Vector_OPENGAUSS        VectorType = "opengauss"
	Vector_TABLESTORE       VectorType = "tablestore"
	Vector_HUAWEI_CLOUD     VectorType = "huawei_cloud"
	Vector_MATRIXONE        VectorType = "matrixone"
)
