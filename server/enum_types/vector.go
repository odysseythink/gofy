package enumtypes

type VectorType string

const (
	Vector_ANALYTICDB    VectorType = "analyticdb"
	Vector_CHROMA        VectorType = "chroma"
	Vector_MILVUS        VectorType = "milvus"
	Vector_MYSCALE       VectorType = "myscale"
	Vector_PGVECTOR      VectorType = "pgvector"
	Vector_PGVECTO_RS    VectorType = "pgvecto-rs"
	Vector_QDRANT        VectorType = "qdrant"
	Vector_RELYT         VectorType = "relyt"
	Vector_TIDB_VECTOR   VectorType = "tidb_vector"
	Vector_WEAVIATE      VectorType = "weaviate"
	Vector_OPENSEARCH    VectorType = "opensearch"
	Vector_TENCENT       VectorType = "tencent"
	Vector_ORACLE        VectorType = "oracle"
	Vector_ELASTICSEARCH VectorType = "elasticsearch"
)
