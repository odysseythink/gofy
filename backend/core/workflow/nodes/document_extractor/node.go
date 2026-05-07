package documentextractor

import (
	"iter"

	"github.com/odysseythink/gofy/backend/core/workflow/nodes/base"
	documentextractornodesentities "github.com/odysseythink/gofy/backend/entities/nodes/document_extractor"
	workflowentities "github.com/odysseythink/gofy/backend/entities/workflow"
	nodesenumtypes "github.com/odysseythink/gofy/backend/enum_types/nodes"
	"github.com/odysseythink/gofy/backend/models"
)

type DocumentExtractorNode struct {
	*base.BaseNode[*documentextractornodesentities.DocumentExtractorNodeData]
}

func (n *DocumentExtractorNode) Type() nodesenumtypes.NodeType {
	return nodesenumtypes.Node_CODE
}

func (n *DocumentExtractorNode) Run() (*workflowentities.NodeRunResult, iter.Seq[any]) {

	return &workflowentities.NodeRunResult{
		Status: models.WorkflowNodeExecutionStatus_SUCCEEDED,
	}, nil
}

func (n *DocumentExtractorNode) ExtractVariableSelectorToVariableMapping(graph_config map[string]any, node_id string, node_data *documentextractornodesentities.DocumentExtractorNodeData) map[string][]string {

	return map[string][]string{}
}

func New() *DocumentExtractorNode {
	return &DocumentExtractorNode{
		BaseNode: &base.BaseNode[*documentextractornodesentities.DocumentExtractorNodeData]{},
	}
}
