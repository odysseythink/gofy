package datasetretriever

import (
	"mlib.com/gofy/server/core/tools/base"
	toolsentities "mlib.com/gofy/server/entities/tools"
)

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
