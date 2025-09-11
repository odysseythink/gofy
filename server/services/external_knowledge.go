package services

import (
	dbengine "mlib.com/gofy/server/db_engine"
	ragentities "mlib.com/gofy/server/entities/rag"
	"mlib.com/mlog"
)

type ExternalDatasetService struct {
}


func(s *ExternalDatasetService) FetchExternalKnowledgeRetrieval(
        tenant_id  string,
        dataset_id  string,
        query  string,
        external_retrieval_parameters map[string]any,
        metadata_condition *ragentities.MetadataCondition,
    ) []map[string]any{
        external_knowledge_binding := new(models.ExternalKnowledgeBinding)
        err := dbengine.Instance().DB.Model(&models.ExternalKnowledgeBindings{}).Where("dataset_id=? and tenant_id = ?",dataset_id, tenant_id).First(external_knowledge_binding).Error
		if err != nil {
			mlog.Errorf("get ExternalKnowledgeBinding failed:%v", err)
			external_knowledge_binding = nil
		}
        if external_knowledge_binding == nil{
            panic(exceptions.NewValueError("external knowledge binding not found"))
}
        external_knowledge_api := new(models.ExternalKnowledgeApi)
        err = dbengine.Instance().DB.Model(&models.ExternalKnowledgeApis{}).Where("id=?",external_knowledge_binding.ExternalKnowledgeApiID).First(external_knowledge_api).Error
        if not external_knowledge_api:
            raise ValueError("external api template not found")

        settings = json.loads(external_knowledge_api.settings)
        headers = {"Content-Type": "application/json"}
        if settings.get("api_key"):
            headers["Authorization"] = f"Bearer {settings.get('api_key')}"
        score_threshold_enabled = external_retrieval_parameters.get("score_threshold_enabled") or False
        score_threshold = external_retrieval_parameters.get("score_threshold", 0.0) if score_threshold_enabled else 0.0
        request_params = {
            "retrieval_setting": {
                "top_k": external_retrieval_parameters.get("top_k"),
                "score_threshold": score_threshold,
            },
            "query": query,
            "knowledge_id": external_knowledge_binding.external_knowledge_id,
            "metadata_condition": metadata_condition.model_dump() if metadata_condition else None,
        }

        response = ExternalDatasetService.process_external_api(
            ExternalKnowledgeApiSetting(
                url=f"{settings.get('endpoint')}/retrieval",
                request_method="post",
                headers=headers,
                params=request_params,
            ),
            None,
        )
        if response.status_code == 200:
            return cast(list[Any], response.json().get("records", []))
        return []
}