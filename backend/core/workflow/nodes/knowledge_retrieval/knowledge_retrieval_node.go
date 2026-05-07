package knowledgeretrieval

import (
	"iter"

	"github.com/odysseythink/gofy/backend/core/workflow/nodes/base"
	knowledgeretrievalnodesentities "github.com/odysseythink/gofy/backend/entities/nodes/knowledge_retrieval"
	workflowentities "github.com/odysseythink/gofy/backend/entities/workflow"
	nodesenumtypes "github.com/odysseythink/gofy/backend/enum_types/nodes"
	"github.com/odysseythink/gofy/backend/models"
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
