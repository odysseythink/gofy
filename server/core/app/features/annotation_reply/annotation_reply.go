package annotationreply

import (
	appenumtypes "mlib.com/gofy/server/enum_types/app"
	"mlib.com/gofy/server/models"
)

type AnnotationReplyFeature struct{}

func (feature *AnnotationReplyFeature) Query(
	app_record *models.App, message *models.Message, query string, user_id string, invoke_from appenumtypes.InvokeFrom,
) *models.MessageAnnotation {
	/*
	   Query app annotations to reply
	   :param app_record: app record
	   :param message: message
	   :param query: query
	   :param user_id: user id
	   :param invoke_from: invoke from
	   :return:
	*/
	// annotation_setting := new(models.AppAnnotationSetting)
	// err := dbengine.Instance().DB.Model(&models.AppAnnotationSetting{}).Where("app_id = ?", app_record.ID).Preload("CollectionBindingDetail").First(annotation_setting).Error
	// if err != nil {
	// 	mlog.Errorf("get AppAnnotationSetting failed:%v", err)
	// 	annotation_setting = nil
	// }

	// if annotation_setting == nil {
	// 	return nil
	// }
	// collection_binding_detail := annotation_setting.CollectionBindingDetail

	// // try:
	// score_threshold := annotation_setting.ScoreThreshold
	// if score_threshold == 0.0 {
	// 	score_threshold = 1.0
	// }
	// embedding_provider_name := collection_binding_detail.ProviderName
	// embedding_model_name := collection_binding_detail.ModelName

	// dataset_collection_binding := manager.Instance.Dataset.GetDatasetCollectionBinding(
	// 	embedding_provider_name, embedding_model_name, "annotation",
	// )

	// dataset := &models.Dataset{
	// 	ID:                     app_record.ID,
	// 	TenantID:               app_record.TenantID,
	// 	IndexingTechnique:      "high_quality",
	// 	EmbeddingModelProvider: embedding_provider_name,
	// 	EmbeddingModel:         embedding_model_name,
	// 	CollectionBindingID:    dataset_collection_binding.ID,
	// }

	// vector = Vector(dataset, attributes=["doc_id", "annotation_id", "app_id"])

	// documents = vector.search_by_vector(
	//     query=query, top_k=1, score_threshold=score_threshold, filter={"group_id": [dataset.id]}
	// )

	// if documents and documents[0].metadata:
	//     annotation_id = documents[0].metadata["annotation_id"]
	//     score = documents[0].metadata["score"]
	//     annotation = AppAnnotationService.get_annotation_by_id(annotation_id)
	//     if annotation:
	//         if invoke_from in {InvokeFrom.SERVICE_API, InvokeFrom.WEB_APP}:
	//             from_source = "api"
	// 		}else:
	//             from_source = "console"
	// 		}

	//         // insert annotation history
	//         AppAnnotationService.add_annotation_history(
	//             annotation.id,
	//             app_record.id,
	//             annotation.question,
	//             annotation.content,
	//             query,
	//             user_id,
	//             message.id,
	//             from_source,
	//             score,
	//         )

	//         return annotation
	// 	}
	// }
	// except Exception as e:
	//     logger.warning(f"Query annotation failed, exception: {str(e)}.")
	//     return None

	return nil
}
