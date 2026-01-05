package knowledgeretrieval

import (
	"iter"

	nodesconstants "mlib.com/gofy/server/constants/workflow/nodes"
	"mlib.com/gofy/server/core/exceptions"
	"mlib.com/gofy/server/core/workflow/nodes/base"
	knowledgeretrievalnodesentities "mlib.com/gofy/server/entities/nodes/knowledge_retrieval"
	workflowentities "mlib.com/gofy/server/entities/workflow"
	nodesenumtypes "mlib.com/gofy/server/enum_types/nodes"
	"mlib.com/gofy/server/models"
	"mlib.com/gofy/server/utils/mapstruct"
	"mlib.com/mlog"
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
func (n *KnowledgeRetrievalNode) ExtractVarSelectorToVarMapping(
	graph_config map[string]any,
	node_id string,
	node_data map[string]any,
) map[string][]string {
	typed_node_data, err := mapstruct.MapToStruct1[*knowledgeretrievalnodesentities.KnowledgeRetrievalNodeData](node_data)
	if err != nil {
		mlog.Error("convert node data to AnswerNodeData failed:%v", err)
		panic(exceptions.NewValueError("convert node data to AnswerNodeData failed"))
	}
	mlog.Debug("node_data=", typed_node_data)

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
