package base

import (
	"fmt"
	"iter"
	"slices"

	"mlib.com/gofy/server/core/exceptions"
	"mlib.com/gofy/server/core/workflow/graph"
	graphengineentities "mlib.com/gofy/server/entities/graph_engine"
	nodesentities "mlib.com/gofy/server/entities/nodes"
	answernodesentities "mlib.com/gofy/server/entities/nodes/answer"
	"mlib.com/gofy/server/entities/nodes/base"
	basenodesentities "mlib.com/gofy/server/entities/nodes/base"
	codenodesentities "mlib.com/gofy/server/entities/nodes/code"
	documentextractornodesentities "mlib.com/gofy/server/entities/nodes/document_extractor"
	endnodesentities "mlib.com/gofy/server/entities/nodes/end"
	evententities "mlib.com/gofy/server/entities/nodes/event"
	httprequestnodesentities "mlib.com/gofy/server/entities/nodes/http_request"
	ifelsenodesentities "mlib.com/gofy/server/entities/nodes/if_else"
	iterationnodesentities "mlib.com/gofy/server/entities/nodes/iteration"
	knowledgeretrievalnodesentities "mlib.com/gofy/server/entities/nodes/knowledge_retrieval"
	listoperatornodesentities "mlib.com/gofy/server/entities/nodes/list_operator"
	llmnodesentities "mlib.com/gofy/server/entities/nodes/llm"
	loopnodesentities "mlib.com/gofy/server/entities/nodes/loop"
	parameterextractornodesentities "mlib.com/gofy/server/entities/nodes/parameter_extractor"
	questionclassifiernodesentities "mlib.com/gofy/server/entities/nodes/question_classifier"
	startnodesentities "mlib.com/gofy/server/entities/nodes/start"
	templatetransformnodesentities "mlib.com/gofy/server/entities/nodes/template_transform"
	toolnodesentities "mlib.com/gofy/server/entities/nodes/tool"
	variableaggregatornodesentities "mlib.com/gofy/server/entities/nodes/variable_aggregator"
	variableassignernodesentities "mlib.com/gofy/server/entities/nodes/variable_assigner"
	workflowentities "mlib.com/gofy/server/entities/workflow"
	appenumtypes "mlib.com/gofy/server/enum_types/app"
	nodesenumtypes "mlib.com/gofy/server/enum_types/nodes"
	"mlib.com/gofy/server/models"
	"mlib.com/gofy/server/utils"
	"mlib.com/mlog"
)

type BaseNode[T nodesentities.GenericNodeData] struct {
	id                string
	tenantID          string
	appID             string
	workflowType      models.WorkflowType
	workflowID        string
	graphConfig       map[string]any
	userID            string
	userFrom          models.UserFrom
	invokeFrom        appenumtypes.InvokeFrom
	workflowCallDepth int
	mGraph            *graph.Graph
	graphRuntimeState *graphengineentities.GraphRuntimeState
	previousNodeID    string
	nodeID            string
	NodeData          T
}

// NewBaseNode creates a new instance of BaseNode.
func NewBaseNode[T nodesentities.GenericNodeData](
	id string,
	config map[string]any,
	graphInitParams *graphengineentities.GraphInitParams,
	gf *graph.Graph,
	graphRuntimeState *graphengineentities.GraphRuntimeState,
	previousNodeID string,
	nodeData T,
) *BaseNode[T] {
	nodeID, ok := config["id"].(string)
	if !ok || nodeID == "" {
		panic(exceptions.NewValueError("Node id is required."))
	}

	var node_data_dict map[string]any
	if _, ok := config["data"]; ok {
		if _, ok := config["data"].(map[string]any); ok {
			node_data_dict = config["data"].(map[string]any)
		}
	}
	// bindata, _ := json.Marshal(node_data_dict)
	// err := json.Unmarshal(bindata, node.NodeData)
	// if err != nil {
	// 	return nil, exceptions.NewValueError(err.Error())
	// }
	mlog.Debugf("------node_data_dict=%#v", node_data_dict)
	if dataer, ok := any(nodeData).(base.NodeDataer); !ok {
		panic(exceptions.NewValueError("nodeData must implement NodeDataer."))
	} else {
		err := dataer.Marshal(node_data_dict, nodeData)
		if err != nil {
			mlog.Errorf("call NodeDataer marshal method failed:%v", err)
			panic(exceptions.NewValueError(fmt.Sprintf("call NodeDataer marshal method failed:%v", err)))
		}
		nodeData = any(dataer).(T)
		mlog.Debugf("------nodeData=%#v", nodeData)
	}
	node := &BaseNode[T]{
		id:                id,
		tenantID:          graphInitParams.TenantID,
		appID:             graphInitParams.AppID,
		workflowType:      graphInitParams.WorkflowType,
		workflowID:        graphInitParams.WorkflowID,
		graphConfig:       graphInitParams.GraphConfig,
		userID:            graphInitParams.UserID,
		userFrom:          graphInitParams.UserFrom,
		invokeFrom:        graphInitParams.InvokeFrom,
		workflowCallDepth: graphInitParams.CallDepth,
		mGraph:            gf,
		graphRuntimeState: graphRuntimeState,
		previousNodeID:    previousNodeID,
		NodeData:          nodeData,
	}
	node.nodeID = nodeID
	return node
}

func (n *BaseNode[T]) PublishTextChunk(text string, value_selector []string) {
	// """
	// Publish text chunk
	// :param text: chunk text
	// :param value_selector: value selector
	// :return:
	// """
	// if n.Callbacks != nil {
	// 	for _, callback := range n.Callbacks {
	// 		callback.OnNodeTextChunk(n.nodeID, text, map[string]any{
	// 			"node_type":               n.NodeType,
	// 			"is_answer_previous_node": n.IsAnswerPreviousNode,
	// 			"value_selector":          value_selector,
	// 		},
	// 		)
	// 	}
	// }
}

func (n *BaseNode[T]) GetDefaultConfig(filters map[string]any) map[string]any {
	// """
	// Get default config of node.
	// :param filters: filter by node config parameters.
	// :return:
	// """
	return map[string]any{}
}

func (n *BaseNode[T]) GetID() string                             { return n.id }
func (n *BaseNode[T]) SetID(val string)                          { n.id = val }
func (n *BaseNode[T]) GetTenantID() string                       { return n.tenantID }
func (n *BaseNode[T]) SetTenantID(val string)                    { n.tenantID = val }
func (n *BaseNode[T]) GetAppID() string                          { return n.appID }
func (n *BaseNode[T]) SetAppID(val string)                       { n.appID = val }
func (n *BaseNode[T]) GetWorkflowType() models.WorkflowType      { return n.workflowType }
func (n *BaseNode[T]) SetWorkflowType(val models.WorkflowType)   { n.workflowType = val }
func (n *BaseNode[T]) GetWorkflowID() string                     { return n.workflowID }
func (n *BaseNode[T]) SetWorkflowID(val string)                  { n.workflowID = val }
func (n *BaseNode[T]) GetGraphConfig() map[string]any            { return n.graphConfig }
func (n *BaseNode[T]) SetGraphConfig(val map[string]any)         { n.graphConfig = val }
func (n *BaseNode[T]) GetUserID() string                         { return n.userID }
func (n *BaseNode[T]) SetUserID(val string)                      { n.userID = val }
func (n *BaseNode[T]) GetUserFrom() models.UserFrom              { return n.userFrom }
func (n *BaseNode[T]) SetUserFrom(val models.UserFrom)           { n.userFrom = val }
func (n *BaseNode[T]) GetInvokeFrom() appenumtypes.InvokeFrom    { return n.invokeFrom }
func (n *BaseNode[T]) SetInvokeFrom(val appenumtypes.InvokeFrom) { n.invokeFrom = val }
func (n *BaseNode[T]) GetWorkflowCallDepth() int                 { return n.workflowCallDepth }
func (n *BaseNode[T]) SetWorkflowCallDepth(val int)              { n.workflowCallDepth = val }
func (n *BaseNode[T]) GetGraph() *graph.Graph                    { return n.mGraph }
func (n *BaseNode[T]) SetGraph(val *graph.Graph)                 { n.mGraph = val }
func (n *BaseNode[T]) GetGraphRuntimeState() *graphengineentities.GraphRuntimeState {
	return n.graphRuntimeState
}
func (n *BaseNode[T]) SetGraphRuntimeState(val *graphengineentities.GraphRuntimeState) {
	n.graphRuntimeState = val
}
func (n *BaseNode[T]) GetPreviousNodeID() string    { return n.previousNodeID }
func (n *BaseNode[T]) SetPreviousNodeID(val string) { n.previousNodeID = val }
func (n *BaseNode[T]) GetNodeID() string            { return n.nodeID }
func (n *BaseNode[T]) SetNodeID(val string)         { n.nodeID = val }
func (n *BaseNode[T]) GetNodeData() any             { return n.NodeData }
func (n *BaseNode[T]) GetBaseNodeData() *basenodesentities.BaseNodeData {
	switch data := any(n.NodeData).(type) {
	case *answernodesentities.AnswerNodeData:
		return data.BaseNodeData
	case *codenodesentities.CodeNodeData:
		return data.BaseNodeData
	case *endnodesentities.EndNodeData:
		return data.BaseNodeData
	case *iterationnodesentities.IterationNodeData:
		return data.BaseNodeData
	case *llmnodesentities.LLMNodeData:
		return data.BaseNodeData
	case *loopnodesentities.LoopNodeData:
		return data.BaseNodeData
	case *toolnodesentities.ToolNodeData:
		return data.BaseNodeData
	case *startnodesentities.StartNodeData:
		return data.BaseNodeData
	case *documentextractornodesentities.DocumentExtractorNodeData:
		return data.BaseNodeData
	case *httprequestnodesentities.HttpRequestNodeData:
		return data.BaseNodeData
	case *ifelsenodesentities.IfElseNodeData:
		return data.BaseNodeData
	case *knowledgeretrievalnodesentities.KnowledgeRetrievalNodeData:
		return data.BaseNodeData
	case *listoperatornodesentities.ListOperatorNodeData:
		return data.BaseNodeData
	case *parameterextractornodesentities.ParameterExtractorNodeData:
		return data.BaseNodeData
	case *questionclassifiernodesentities.QuestionClassifierNodeData:
		return data.BaseNodeData
	case *templatetransformnodesentities.TemplateTransformNodeData:
		return data.BaseNodeData
	case *variableaggregatornodesentities.VariableAggregatorNodeData:
		return data.BaseNodeData
	case *variableassignernodesentities.VariableAssignerNodeData:
		return data.BaseNodeData
	}
	return nil
}

func (n *BaseNode[T]) guessNodeType() nodesenumtypes.NodeType {
	switch any(n.NodeData).(type) {
	case *answernodesentities.AnswerNodeData:
		return nodesenumtypes.Node_ANSWER
	case *codenodesentities.CodeNodeData:
		return nodesenumtypes.Node_CODE
	case *endnodesentities.EndNodeData:
		return nodesenumtypes.Node_END
	case *iterationnodesentities.IterationNodeData:
		return nodesenumtypes.Node_ITERATION
	case *llmnodesentities.LLMNodeData:
		return nodesenumtypes.Node_LLM
	case *loopnodesentities.LoopNodeData:
		return nodesenumtypes.Node_LLM
	case *toolnodesentities.ToolNodeData:
		return nodesenumtypes.Node_TOOL
	case *startnodesentities.StartNodeData:
		return nodesenumtypes.Node_START
	case *documentextractornodesentities.DocumentExtractorNodeData:
		return nodesenumtypes.Node_DOCUMENT_EXTRACTOR
	case *httprequestnodesentities.HttpRequestNodeData:
		return nodesenumtypes.Node_HTTP_REQUEST
	case *ifelsenodesentities.IfElseNodeData:
		return nodesenumtypes.Node_IF_ELSE
	case *knowledgeretrievalnodesentities.KnowledgeRetrievalNodeData:
		return nodesenumtypes.Node_KNOWLEDGE_RETRIEVAL
	case *listoperatornodesentities.ListOperatorNodeData:
		return nodesenumtypes.Node_LIST_OPERATOR
	case *parameterextractornodesentities.ParameterExtractorNodeData:
		return nodesenumtypes.Node_PARAMETER_EXTRACTOR
	case *questionclassifiernodesentities.QuestionClassifierNodeData:
		return nodesenumtypes.Node_QUESTION_CLASSIFIER
	case *templatetransformnodesentities.TemplateTransformNodeData:
		return nodesenumtypes.Node_TEMPLATE_TRANSFORM
	case *variableaggregatornodesentities.VariableAggregatorNodeData:
		return nodesenumtypes.Node_VARIABLE_AGGREGATOR
	case *variableassignernodesentities.VariableAssignerNodeData:
		return nodesenumtypes.Node_VARIABLE_ASSIGNER
	}
	return nodesenumtypes.NodeType(-1)
}

func (n *BaseNode[T]) ShouldContinueOnError() bool {
	/*
		judge if should continue on error

		Returns:
			bool: if should continue on error
	*/
	return n.GetBaseNodeData().ErrorStrategy != "" && slices.Contains(nodesenumtypes.CONTINUE_ON_ERROR_NODE_TYPE, n.guessNodeType())
}
func (n *BaseNode[T]) ShouldRetry() bool {
	/*
		judge if should retry

		Returns:
			bool: if should retry
	*/
	return n.GetBaseNodeData().RetryConfig.RetryEnabled && slices.Contains(nodesenumtypes.RETRY_ON_ERROR_NODE_TYPE, n.guessNodeType())
}

func (n *BaseNode[T]) RunIter(ner Noder) iter.Seq[any] {
	return func(yield func(any) bool) {
		defer func() {
			if r := recover(); r != nil {
				mlog.Errorf("panic recover:%v\n%s", r, utils.GetCurrentGoroutineStack())
				if exp, ok := r.(error); ok {
					mlog.Errorf("Node %s failed to run", ner.GetNodeID())
					yield(&evententities.RunCompletedEvent{
						RunResult: &workflowentities.NodeRunResult{
							Status:    models.WorkflowNodeExecutionStatus_FAILED,
							Error:     exp.Error(),
							ErrorType: "WorkflowNodeError",
						},
					})
				} else {
					panic(r)
				}
			}
		}()
		result, iterresult := ner.Run()
		if result != nil {
			yield(&evententities.RunCompletedEvent{
				RunResult: result,
			})
			mlog.Debugf("----result.Inputs=%#v ", result.Inputs)
			return
		} else {
			for item := range iterresult {
				if !yield(item) {
					return
				}
			}
		}
	}
}
