package services

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/odysseythink/mlog"
	"mlib.com/gofy/server/core/exceptions"
	dbengine "mlib.com/gofy/server/db_engine"
	ragentities "mlib.com/gofy/server/entities/rag"
	"mlib.com/gofy/server/models"
	"mlib.com/gofy/server/utils/mapstruct"
)

type ExternalDatasetService struct {
}

func (s *ExternalDatasetService) FetchExternalKnowledgeRetrieval(
	tenant_id string,
	dataset_id string,
	query string,
	external_retrieval_parameters map[string]any,
	metadata_condition *ragentities.MetadataCondition,
) []map[string]any {
	external_knowledge_binding := new(models.ExternalKnowledgeBinding)
	err := dbengine.Instance().DB.Model(&models.ExternalKnowledgeBinding{}).Where("dataset_id=? and tenant_id = ?", dataset_id, tenant_id).First(external_knowledge_binding).Error
	if err != nil {
		mlog.Errorf("get ExternalKnowledgeBinding failed:%v", err)
		external_knowledge_binding = nil
	}
	if external_knowledge_binding == nil {
		panic(exceptions.NewValueError("external knowledge binding not found"))
	}
	external_knowledge_api := new(models.ExternalKnowledgeApi)
	err = dbengine.Instance().DB.Model(&models.ExternalKnowledgeApi{}).Where("id=?", external_knowledge_binding.ExternalKnowledgeApiID).First(external_knowledge_api).Error
	if err != nil {
		mlog.Errorf("get ExternalKnowledgeApi failed:%v", err)
		external_knowledge_api = nil
	}
	if external_knowledge_api == nil {
		panic(exceptions.NewValueError("external api template not found"))
	}

	var settings map[string]any
	err = json.Unmarshal([]byte(external_knowledge_api.Settings), &settings)
	if err != nil {
		mlog.Errorf("get ExternalKnowledgeApi failed:%v", err)
		external_knowledge_api = nil
	}
	if external_knowledge_api == nil {
		panic(exceptions.NewValueError("external api template not found"))
	}
	score_threshold_enabled := mapstruct.Get(external_retrieval_parameters, "score_threshold_enabled", false)
	score_threshold := mapstruct.Get(external_retrieval_parameters, "score_threshold", 0.0)
	if !score_threshold_enabled {
		score_threshold = 0.0
	}
	request_params := map[string]any{
		"retrieval_setting": map[string]any{
			"top_k":           external_retrieval_parameters["top_k"],
			"score_threshold": score_threshold,
		},
		"query":              query,
		"knowledge_id":       external_knowledge_binding.ExternalKnowledgeID,
		"metadata_condition": map[string]any{},
	}
	if metadata_condition != nil {
		request_params["metadata_condition"] = metadata_condition.ToDict()
	}
	payload_body, _ := json.Marshal(metadata_condition)
	req, err := http.NewRequest("POST", mapstruct.Get(settings, "endpoint", "")+"/retrieval", bytes.NewReader(payload_body))
	if err != nil {
		mlog.Errorf("create http request failed:", err)
		panic(exceptions.NewValueError("create http request failed"))
	}
	req.Header.Set("Content-Type", "application/json")
	api_key := mapstruct.Get(settings, "api_key", "")
	if api_key != "" {
		req.Header.Set("Authorization", "Bearer "+api_key)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		mlog.Errorf("http post failed:%v", err)
		panic(exceptions.NewValueError("http post failed"))
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		var res []map[string]any
		if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			mlog.Errorf("http response is not slice:%v", err)
			panic(exceptions.NewValueError("http response is not slice"))
		}
		return res
	}
	return []map[string]any{}
}
