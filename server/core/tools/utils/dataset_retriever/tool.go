package datasetretriever

import (
	idxtoolcbhandler "mlib.com/gofy/server/core/callback_handler/index_tool"
	"mlib.com/gofy/server/core/tools/base"
	toolsentities "mlib.com/gofy/server/entities/tools"
)

type DatasetRetrieverToolor interface {
	Run(query string) string
}
type DatasetRetrieverBaseTool struct {
	Name           string                                              `json:"name"`
	Description    string                                              `json:"description"`
	TenantID       string                                              `json:"tenant_id"`
	TopK           int                                                 `json:"top_k"`
	ScoreThreshold float64                                             `json:"score_threshold"`
	HitCallbacks   []*idxtoolcbhandler.DatasetIndexToolCallbackHandler `json:"hit_callbacks"`
	ReturnResource bool                                                `json:"return_resource"`
	RetrieverFrom  string                                              `json:"retriever_from"`
	ModelConfig    map[string]any                                      `json:"model_config"`
}

func NewDatasetRetrieverBaseTool() *DatasetRetrieverBaseTool {
	return &DatasetRetrieverBaseTool{
		Name:         "dataset",
		Description:  "use this to retrieve a dataset. ",
		TopK:         2,
		HitCallbacks: make([]*idxtoolcbhandler.DatasetIndexToolCallbackHandler, 0),
	}
}

type DatasetRetrieverTool struct {
	*base.Tool
	RetrievalTool *DatasetRetrieverBaseTool
}

func NewDatasetRetrieverTool(entity *toolsentities.ToolEntity, runtime *base.ToolRuntime, retrieval_tool *DatasetRetrieverBaseTool) *DatasetRetrieverTool {
	if retrieval_tool == nil {
		retrieval_tool = NewDatasetRetrieverBaseTool()
	}
	return &DatasetRetrieverTool{
		Tool:          base.NewTool(entity, runtime),
		RetrievalTool: retrieval_tool,
	}
}
