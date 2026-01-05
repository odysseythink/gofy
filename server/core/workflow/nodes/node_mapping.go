package nodes

import (
	"fmt"

	"mlib.com/gofy/server/core/exceptions"
	"mlib.com/gofy/server/core/workflow/graph"
	"mlib.com/gofy/server/core/workflow/nodes/answer"
	"mlib.com/gofy/server/core/workflow/nodes/base"
	"mlib.com/gofy/server/core/workflow/nodes/code"
	documentextractor "mlib.com/gofy/server/core/workflow/nodes/document_extractor"
	"mlib.com/gofy/server/core/workflow/nodes/end"
	httprequest "mlib.com/gofy/server/core/workflow/nodes/http_request"
	ifelse "mlib.com/gofy/server/core/workflow/nodes/if_else"
	"mlib.com/gofy/server/core/workflow/nodes/iteration"
	knowledgeretrieval "mlib.com/gofy/server/core/workflow/nodes/knowledge_retrieval"
	listoperator "mlib.com/gofy/server/core/workflow/nodes/list_operator"
	"mlib.com/gofy/server/core/workflow/nodes/llm"
	"mlib.com/gofy/server/core/workflow/nodes/loop"
	parameterextractor "mlib.com/gofy/server/core/workflow/nodes/parameter_extractor"
	questionclassifier "mlib.com/gofy/server/core/workflow/nodes/question_classifier"
	"mlib.com/gofy/server/core/workflow/nodes/start"
	templatetransform "mlib.com/gofy/server/core/workflow/nodes/template_transform"
	"mlib.com/gofy/server/core/workflow/nodes/tool"
	variableaggregator "mlib.com/gofy/server/core/workflow/nodes/variable_aggregator"
	variableassigner "mlib.com/gofy/server/core/workflow/nodes/variable_assigner"
	graphengineentities "mlib.com/gofy/server/entities/graph_engine"
	answernodesentities "mlib.com/gofy/server/entities/nodes/answer"
	codenodesentities "mlib.com/gofy/server/entities/nodes/code"
	documentextractornodesentities "mlib.com/gofy/server/entities/nodes/document_extractor"
	endnodesentities "mlib.com/gofy/server/entities/nodes/end"
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
	nodesenumtypes "mlib.com/gofy/server/enum_types/nodes"
	"mlib.com/mlog"
)

func NewNode(
	id string,
	config map[string]any,
	graphInitParams *graphengineentities.GraphInitParams,
	graph *graph.Graph,
	graphRuntimeState *graphengineentities.GraphRuntimeState,
	previousNodeID string,
	nodeType nodesenumtypes.NodeType,
) base.Noder {
	switch nodeType {
	case nodesenumtypes.Node_START:
		ndata := startnodesentities.New()
		bn := base.NewBaseNode(id, config, graphInitParams, graph, graphRuntimeState, previousNodeID, ndata)
		n := &start.StartNode{
			BaseNode: bn,
		}
		return n
	case nodesenumtypes.Node_END:
		ndata := endnodesentities.New()
		bn := base.NewBaseNode(id, config, graphInitParams, graph, graphRuntimeState, previousNodeID, ndata)
		n := &end.EndNode{
			BaseNode: bn,
		}
		return n
	case nodesenumtypes.Node_ANSWER:
		ndata := answernodesentities.New()
		bn := base.NewBaseNode(id, config, graphInitParams, graph, graphRuntimeState, previousNodeID, ndata)
		n := &answer.AnswerNode{
			BaseNode: bn,
		}
		return n
	case nodesenumtypes.Node_LLM:
		ndata := llmnodesentities.New()
		bn := base.NewBaseNode(id, config, graphInitParams, graph, graphRuntimeState, previousNodeID, ndata)
		n := &llm.LLMNode{
			BaseNode: bn,
		}
		return n
	case nodesenumtypes.Node_KNOWLEDGE_RETRIEVAL:
		ndata := new(knowledgeretrievalnodesentities.KnowledgeRetrievalNodeData)
		bn := base.NewBaseNode(id, config, graphInitParams, graph, graphRuntimeState, previousNodeID, ndata)
		n := &knowledgeretrieval.KnowledgeRetrievalNode{
			BaseNode: bn,
		}
		return n
	case nodesenumtypes.Node_IF_ELSE:
		ndata := ifelsenodesentities.New()
		bn := base.NewBaseNode(id, config, graphInitParams, graph, graphRuntimeState, previousNodeID, ndata)
		n := &ifelse.IfElseNode{
			BaseNode: bn,
		}
		return n
	case nodesenumtypes.Node_CODE:
		ndata := codenodesentities.New()
		bn := base.NewBaseNode(id, config, graphInitParams, graph, graphRuntimeState, previousNodeID, ndata)
		n := &code.CodeNode{
			BaseNode: bn,
		}
		return n
	case nodesenumtypes.Node_TEMPLATE_TRANSFORM:
		ndata := new(templatetransformnodesentities.TemplateTransformNodeData)
		bn := base.NewBaseNode(id, config, graphInitParams, graph, graphRuntimeState, previousNodeID, ndata)
		n := &templatetransform.TemplateTransformNode{
			BaseNode: bn,
		}
		return n
	case nodesenumtypes.Node_QUESTION_CLASSIFIER:
		ndata := questionclassifiernodesentities.New()
		bn := base.NewBaseNode(id, config, graphInitParams, graph, graphRuntimeState, previousNodeID, ndata)
		n := &questionclassifier.QuestionClassifierNode{
			BaseNode: bn,
		}
		return n
	case nodesenumtypes.Node_HTTP_REQUEST:
		ndata := httprequestnodesentities.New()
		bn := base.NewBaseNode(id, config, graphInitParams, graph, graphRuntimeState, previousNodeID, ndata)
		n := &httprequest.HttpRequestNode{
			BaseNode: bn,
		}
		return n
	case nodesenumtypes.Node_TOOL:
		ndata := new(toolnodesentities.ToolNodeData)
		bn := base.NewBaseNode(id, config, graphInitParams, graph, graphRuntimeState, previousNodeID, ndata)
		n := &tool.ToolNode{
			BaseNode: bn,
		}
		return n
	case nodesenumtypes.Node_VARIABLE_AGGREGATOR:
		ndata := new(variableaggregatornodesentities.VariableAggregatorNodeData)
		bn := base.NewBaseNode(id, config, graphInitParams, graph, graphRuntimeState, previousNodeID, ndata)
		n := &variableaggregator.VariableAggregatorNode{
			BaseNode: bn,
		}
		return n
	case nodesenumtypes.Node_LOOP:
		ndata := new(loopnodesentities.LoopNodeData)
		bn := base.NewBaseNode(id, config, graphInitParams, graph, graphRuntimeState, previousNodeID, ndata)
		n := &loop.LoopNode{
			BaseNode: bn,
		}
		return n
	case nodesenumtypes.Node_ITERATION:
		ndata := new(iterationnodesentities.IterationNodeData)
		bn := base.NewBaseNode(id, config, graphInitParams, graph, graphRuntimeState, previousNodeID, ndata)
		n := &iteration.IterationNode{
			BaseNode: bn,
		}
		return n
	case nodesenumtypes.Node_PARAMETER_EXTRACTOR:
		ndata := new(parameterextractornodesentities.ParameterExtractorNodeData)
		bn := base.NewBaseNode(id, config, graphInitParams, graph, graphRuntimeState, previousNodeID, ndata)
		n := &parameterextractor.ParameterExtractorNode{
			BaseNode: bn,
		}
		return n
	case nodesenumtypes.Node_VARIABLE_ASSIGNER:
		ndata := new(variableassignernodesentities.VariableAssignerNodeData)
		bn := base.NewBaseNode(id, config, graphInitParams, graph, graphRuntimeState, previousNodeID, ndata)
		n := &variableassigner.VariableAssignerNode{
			BaseNode: bn,
		}
		return n
	case nodesenumtypes.Node_DOCUMENT_EXTRACTOR:
		ndata := new(documentextractornodesentities.DocumentExtractorNodeData)
		bn := base.NewBaseNode(id, config, graphInitParams, graph, graphRuntimeState, previousNodeID, ndata)
		n := &documentextractor.DocumentExtractorNode{
			BaseNode: bn,
		}
		return n
	case nodesenumtypes.Node_LIST_OPERATOR:
		ndata := new(listoperatornodesentities.ListOperatorNodeData)
		bn := base.NewBaseNode(id, config, graphInitParams, graph, graphRuntimeState, previousNodeID, ndata)
		n := &listoperator.ListOperatorNode{
			BaseNode: bn,
		}
		return n
	}
	mlog.Errorf("unsurported node_type(%v)", nodeType)
	panic(exceptions.NewNotImplementedError(fmt.Sprintf("unsurported node_type(%v)", nodeType)))
}
