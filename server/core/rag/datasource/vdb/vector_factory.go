package vdb

import (
	"mlib.com/confy"
	"mlib.com/gofy/server/core/exceptions"
	"mlib.com/gofy/server/core/rag/embedding"
	dbengine "mlib.com/gofy/server/db_engine"
	ragentities "mlib.com/gofy/server/entities/rag"
	vectorenumtypes "mlib.com/gofy/server/enum_types/rag/vector"
	"mlib.com/gofy/server/models"
	"mlib.com/mlog"
)

type AbstractVectorFactory interface {
	InitVector(dataset *models.Dataset, attributes []any, embeddings embedding.IEmbeddings) ragentities.IVector
}

func GenIndexStructDict(vector_type vectorenumtypes.VectorType, collection_name string) map[string]any {
	return map[string]any{"type": vector_type, "vector_store": map[string]any{"class_prefix": collection_name}}
}

type Vector struct {
	_dataset          *models.Dataset
	_embeddings       embedding.IEmbeddings
	_attributes       []string
	_vector_processor ragentities.IVector
}

func NewVector(dataset *models.Dataset, attributes []string) *Vector {
	if len(attributes) == 0 {
		attributes = []string{"doc_id", "dataset_id", "document_id", "doc_hash"}
	}
	v := &Vector{
		_dataset:    dataset,
		_attributes: attributes,
	}
	v._embeddings = v._get_embeddings()
	v._vector_processor = v._init_vector()
	return v
}
func (v *Vector) _init_vector() ragentities.IVector {
	vector_type := confy.GetWithDefault[string]("vector_store", "")
	index_struct_dict := v._dataset.IndexStructDict()
	if len(index_struct_dict) > 0 {
		if _, ok := index_struct_dict["type"]; ok {
			if _, ok := index_struct_dict["type"].(string); ok {
				vector_type = index_struct_dict["type"].(string)
			}
		}
	} else {
		if confy.GetWithDefault[bool]("vector_store_whitelist_enable", false) {
			whitelist := new(models.Whitelist)
			err := dbengine.Instance().DB.Model(&models.Whitelist{}).Where("tenant_id = ? and category = ?", v._dataset.TenantID, "vector_db").First(whitelist).Error
			if err != nil {
				mlog.Errorf("get Whitelist failed:%v", err)
				whitelist = nil
			}
			if whitelist != nil {
				vector_type = string(vectorenumtypes.Vector_TIDB_ON_QDRANT)
			}
		}
	}
	if vector_type == "" {
		panic(exceptions.NewValueError("Vector store must be specified."))
	}
	// vector_factory_cls = v.get_vector_factory(vector_type)
	// return vector_factory_cls().init_vector(v._dataset, v._attributes, v._embeddings)
	return nil
}

func (v *Vector) get_vector_factory(vector_type string) AbstractVectorFactory {
	return nil
	// switch VectorType(vector_type) {
	//     case Vector_CHROMA:
	//         from core.rag.datasource.vdb.chroma.chroma_vector import ChromaVectorFactory

	//         return ChromaVectorFactory
	//     case Vector_MILVUS:
	//         from core.rag.datasource.vdb.milvus.milvus_vector import MilvusVectorFactory

	//         return MilvusVectorFactory
	//     case Vector_MYSCALE:
	//         from core.rag.datasource.vdb.myscale.myscale_vector import MyScaleVectorFactory

	//         return MyScaleVectorFactory
	//     case Vector_PGVECTOR:
	//         from core.rag.datasource.vdb.pgvector.pgvector import PGVectorFactory

	//         return PGVectorFactory
	//     case Vector_VASTBASE:
	//         from core.rag.datasource.vdb.pyvastbase.vastbase_vector import VastbaseVectorFactory

	//         return VastbaseVectorFactory
	//     case Vector_PGVECTO_RS:
	//         from core.rag.datasource.vdb.pgvecto_rs.pgvecto_rs import PGVectoRSFactory

	//         return PGVectoRSFactory
	//     case Vector_QDRANT:
	//         from core.rag.datasource.vdb.qdrant.qdrant_vector import QdrantVectorFactory

	//         return QdrantVectorFactory
	//     case Vector_RELYT:
	//         from core.rag.datasource.vdb.relyt.relyt_vector import RelytVectorFactory

	//         return RelytVectorFactory
	//     case Vector_ELASTICSEARCH:
	//         from core.rag.datasource.vdb.elasticsearch.elasticsearch_vector import ElasticSearchVectorFactory

	//         return ElasticSearchVectorFactory
	//     case Vector_ELASTICSEARCH_JA:
	//         from core.rag.datasource.vdb.elasticsearch.elasticsearch_ja_vector import (
	//             ElasticSearchJaVectorFactory,
	//         )

	//         return ElasticSearchJaVectorFactory
	//     case Vector_TIDB_VECTOR:
	//         from core.rag.datasource.vdb.tidb_vector.tidb_vector import TiDBVectorFactory

	//         return TiDBVectorFactory
	//     case Vector_WEAVIATE:
	//         from core.rag.datasource.vdb.weaviate.weaviate_vector import WeaviateVectorFactory

	//         return WeaviateVectorFactory
	//     case Vector_TENCENT:
	//         from core.rag.datasource.vdb.tencent.tencent_vector import TencentVectorFactory

	//         return TencentVectorFactory
	//     case Vector_ORACLE:
	//         from core.rag.datasource.vdb.oracle.oraclevector import OracleVectorFactory

	//         return OracleVectorFactory
	//     case Vector_OPENSEARCH:
	//         from core.rag.datasource.vdb.opensearch.opensearch_vector import OpenSearchVectorFactory

	//         return OpenSearchVectorFactory
	//     case Vector_ANALYTICDB:
	//         from core.rag.datasource.vdb.analyticdb.analyticdb_vector import AnalyticdbVectorFactory

	//         return AnalyticdbVectorFactory
	//     case Vector_COUCHBASE:
	//         from core.rag.datasource.vdb.couchbase.couchbase_vector import CouchbaseVectorFactory

	//         return CouchbaseVectorFactory
	//     case Vector_BAIDU:
	//         from core.rag.datasource.vdb.baidu.baidu_vector import BaiduVectorFactory

	//         return BaiduVectorFactory
	//     case Vector_VIKINGDB:
	//         from core.rag.datasource.vdb.vikingdb.vikingdb_vector import VikingDBVectorFactory

	//         return VikingDBVectorFactory
	//     case Vector_UPSTASH:
	//         from core.rag.datasource.vdb.upstash.upstash_vector import UpstashVectorFactory

	//         return UpstashVectorFactory
	//     case Vector_TIDB_ON_QDRANT:
	//         from core.rag.datasource.vdb.tidb_on_qdrant.tidb_on_qdrant_vector import TidbOnQdrantVectorFactory

	//         return TidbOnQdrantVectorFactory
	//     case Vector_LINDORM:
	//         from core.rag.datasource.vdb.lindorm.lindorm_vector import LindormVectorStoreFactory

	//         return LindormVectorStoreFactory
	//     case Vector_OCEANBASE:
	//         from core.rag.datasource.vdb.oceanbase.oceanbase_vector import OceanBaseVectorFactory

	//         return OceanBaseVectorFactory
	//     case Vector_OPENGAUSS:
	//         from core.rag.datasource.vdb.opengauss.opengauss import OpenGaussFactory

	//         return OpenGaussFactory
	//     case Vector_TABLESTORE:
	//         from core.rag.datasource.vdb.tablestore.tablestore_vector import TableStoreVectorFactory

	//         return TableStoreVectorFactory
	//     case Vector_HUAWEI_CLOUD:
	//         from core.rag.datasource.vdb.huawei.huawei_cloud_vector import HuaweiCloudVectorFactory

	//         return HuaweiCloudVectorFactory
	//     case Vector_MATRIXONE:
	//         from core.rag.datasource.vdb.matrixone.matrixone_vector import MatrixoneVectorFactory

	//         return MatrixoneVectorFactory
	//     default:
	//         panic(exceptions.NewValueError(fmt.Sprintf("Vector store {%s} is not supported.", vector_type)))
	//     }
}

// func(v *Vector) create(texts: Optional[list] = None, **kwargs):
//         if texts:
//             start = time.time()
//             logger.info("start embedding %s texts %s", len(texts), start)
//             batch_size = 1000
//             total_batches = len(texts) + batch_size - 1
//             for i in range(0, len(texts), batch_size):
//                 batch = texts[i : i + batch_size]
//                 batch_start = time.time()
//                 logger.info("Processing batch %s/%s (%s texts)", i // batch_size + 1, total_batches, len(batch))
//                 batch_embeddings = v._embeddings.embed_documents([document.page_content for document in batch])
//                 logger.info(
//                     "Embedding batch %s/%s took %s s", i // batch_size + 1, total_batches, time.time() - batch_start
//                 )
//                 v._vector_processor.create(texts=batch, embeddings=batch_embeddings, **kwargs)
//             logger.info("Embedding %s texts took %s s", len(texts), time.time() - start)

// }
// func(v *Vector) add_texts(documents: list[Document], **kwargs):
//         if kwargs.get("duplicate_check", False):
//             documents = v._filter_duplicate_texts(documents)

//         embeddings = v._embeddings.embed_documents([document.page_content for document in documents])
//         v._vector_processor.create(texts=documents, embeddings=embeddings, **kwargs)

// }
// func(v *Vector) text_exists(id string) -> bool:
//         return v._vector_processor.text_exists(id)

// }
// func(v *Vector) delete_by_ids(ids: list[str]) -> None:
//         v._vector_processor.delete_by_ids(ids)

// }
// func(v *Vector) delete_by_metadata_field(key string, value string) -> None:
//         v._vector_processor.delete_by_metadata_field(key, value)

// }
// func(v *Vector) search_by_vector(query string, **kwargs: Any) -> list[Document]:
//         query_vector = v._embeddings.embed_query(query)
//         return v._vector_processor.search_by_vector(query_vector, **kwargs)

// }
// func(v *Vector) search_by_full_text(query string, **kwargs: Any) -> list[Document]:
//         return v._vector_processor.search_by_full_text(query, **kwargs)

// }
// func(v *Vector) delete() -> None:
//         v._vector_processor.delete()
//         // delete collection redis cache
//         if v._vector_processor.collection_name:
//             collection_exist_cache_key = f"vector_indexing_{v._vector_processor.collection_name}"
//             redis_client.delete(collection_exist_cache_key)

// }
func (v *Vector) _get_embeddings() embedding.IEmbeddings {
	// model_manager = ModelManager()

	// embedding_model = model_manager.get_model_instance(
	//     tenant_id=v._dataset.tenant_id,
	//     provider=v._dataset.embedding_model_provider,
	//     model_type=ModelType.TEXT_EMBEDDING,
	//     model=v._dataset.embedding_model,
	// )
	// return CacheEmbedding(embedding_model)
	return nil
}

// func(v *Vector) _filter_duplicate_texts(texts: list[Document]) -> list[Document]:
//         for text in texts.copy():
//             if text.metadata is None:
//                 continue
//             doc_id = text.metadata["doc_id"]
//             if doc_id:
//                 exists_duplicate_node = v.text_exists(doc_id)
//                 if exists_duplicate_node:
//                     texts.remove(text)

//         return texts

// }
// func(v *Vector) __getattr__(name):
//         if v._vector_processor is not None:
//             method = getattr(v._vector_processor, name)
//             if callable(method):
//                 return method

//         raise AttributeError(f"'vector_processor' object has no attribute '{name}'")
// }
