package datasetretriever

import (
	"github.com/odysseythink/gofy/backend/core/tools/base"
	toolsentities "github.com/odysseythink/gofy/backend/entities/tools"
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
