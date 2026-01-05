package knowledgeretrieval

import (
	"iter"

	nodesconstants "mlib.com/gofy/server/constants/workflow/nodes"
	"mlib.com/gofy/server/core/workflow/nodes/base"
	knowledgeretrievalnodesentities "mlib.com/gofy/server/entities/nodes/knowledge_retrieval"
	workflowentities "mlib.com/gofy/server/entities/workflow"
	nodesenumtypes "mlib.com/gofy/server/enum_types/nodes"
	"mlib.com/gofy/server/models"
)

type KnowledgeRetrievalNode struct {
	*base.BaseNode[*knowledgeretrievalnodesentities.KnowledgeRetrievalNodeData]
}

func (n *KnowledgeRetrievalNode) Type() nodesenumtypes.NodeType {
	return nodesenumtypes.Node_KNOWLEDGE_RETRIEVAL
}

func (n *KnowledgeRetrievalNode) Run() (*workflowentities.NodeRunResult, iter.Seq[any]) {

	return &workflowentities.NodeRunResult{
		Status: models.WorkflowNodeExecutionStatus_SUCCEEDED,
	}, nil
}

func (n *KnowledgeRetrievalNode) ExtractVariableSelectorToVariableMapping(graph_config map[string]any, node_id string, node_data *knowledgeretrievalnodesentities.KnowledgeRetrievalNodeData) map[string][]string {

	return map[string][]string{}
}
func New() *KnowledgeRetrievalNode {
	return &KnowledgeRetrievalNode{
		BaseNode: &base.BaseNode[*knowledgeretrievalnodesentities.KnowledgeRetrievalNodeData]{},
	}
}
func init() {
	nodesconstants.Regist(New())
}
