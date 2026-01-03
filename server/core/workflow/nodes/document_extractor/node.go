package documentextractor

import (
	"iter"

	nodesconstants "mlib.com/gofy/server/constants/workflow/nodes"
	"mlib.com/gofy/server/core/workflow/nodes/base"
	documentextractornodesentities "mlib.com/gofy/server/entities/nodes/document_extractor"
	workflowentities "mlib.com/gofy/server/entities/workflow"
	nodesenumtypes "mlib.com/gofy/server/enum_types/nodes"
	"mlib.com/gofy/server/models"
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
func init() {
	nodesconstants.Regist(New())
}
