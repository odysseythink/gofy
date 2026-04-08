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
	humaninput "mlib.com/gofy/server/core/workflow/nodes/human_input"
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
	humaninputnodesentities "mlib.com/gofy/server/entities/nodes/human_input"
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

const (
	LATEST_VERSION = "latest"
)

var (
	NODE_TYPE_CLASSES_MAPPING = map[nodesenumtypes.NodeType]map[string]base.Noder{
		nodesenumtypes.Node_START: {
			LATEST_VERSION: start.New(),
			"1":            start.New(),
		},
		nodesenumtypes.Node_END: {
			LATEST_VERSION: end.New(),
			"1":            end.New(),
		},
		nodesenumtypes.Node_ANSWER: {
			LATEST_VERSION: answer.New(),
			"1":            answer.New(),
		},
		nodesenumtypes.Node_LLM: {
			LATEST_VERSION: llm.New(),
			"1":            llm.New(),
		},
		nodesenumtypes.Node_KNOWLEDGE_RETRIEVAL: {
			LATEST_VERSION: knowledgeretrieval.New(),
			"1":            knowledgeretrieval.New(),
		},
		nodesenumtypes.Node_IF_ELSE: {
			LATEST_VERSION: ifelse.New(),
			"1":            ifelse.New(),
		},
		nodesenumtypes.Node_CODE: {
			LATEST_VERSION: code.New(),
			"1":            code.New(),
		},
		nodesenumtypes.Node_TEMPLATE_TRANSFORM: {
			LATEST_VERSION: templatetransform.New(),
			"1":            templatetransform.New(),
		},
		nodesenumtypes.Node_QUESTION_CLASSIFIER: {
			LATEST_VERSION: questionclassifier.New(),
			"1":            questionclassifier.New(),
		},
		nodesenumtypes.Node_HTTP_REQUEST: {
			LATEST_VERSION: httprequest.New(),
			"1":            httprequest.New(),
		},
		nodesenumtypes.Node_VARIABLE_AGGREGATOR: {
			LATEST_VERSION: variableaggregator.New(),
			"1":            variableaggregator.New(),
		},
		nodesenumtypes.Node_ITERATION: {
			LATEST_VERSION: iteration.New(),
			"1":            iteration.New(),
		},

		nodesenumtypes.Node_PARAMETER_EXTRACTOR: {
			LATEST_VERSION: parameterextractor.New(),
			"1":            parameterextractor.New(),
		},
		nodesenumtypes.Node_VARIABLE_ASSIGNER: {
			LATEST_VERSION: variableassigner.New(),
			"1":            variableassigner.New(),
			"2":            variableassigner.New(),
		},
		nodesenumtypes.Node_DOCUMENT_EXTRACTOR: {
			LATEST_VERSION: documentextractor.New(),
			"1":            documentextractor.New(),
		},
		nodesenumtypes.Node_LIST_OPERATOR: {
			LATEST_VERSION: listoperator.New(),
			"1":            listoperator.New(),
		},
		nodesenumtypes.Node_HUMAN_INPUT: {
			LATEST_VERSION: humaninput.New(),
			"1":            humaninput.New(),
		},
	}
)

func ExtractVariableSelectorToVariableMappingByNoder(graph_config map[string]any, node_id string, noder base.Noder) map[string][]string {
	if specific_noder, ok := any(noder).(base.SpecificNoder[*startnodesentities.StartNodeData]); ok {
		if real_node_data, ok := noder.GetNodeData().(*startnodesentities.StartNodeData); ok {
			return specific_noder.ExtractVariableSelectorToVariableMapping(graph_config, node_id, real_node_data)
		}
	} else if specific_noder, ok := any(noder).(base.SpecificNoder[*endnodesentities.EndNodeData]); ok {
		if real_node_data, ok := noder.GetNodeData().(*endnodesentities.EndNodeData); ok {
			return specific_noder.ExtractVariableSelectorToVariableMapping(graph_config, node_id, real_node_data)
		}
	} else if specific_noder, ok := any(noder).(base.SpecificNoder[*answernodesentities.AnswerNodeData]); ok {
		if real_node_data, ok := noder.GetNodeData().(*answernodesentities.AnswerNodeData); ok {
			return specific_noder.ExtractVariableSelectorToVariableMapping(graph_config, node_id, real_node_data)
		}
	} else if specific_noder, ok := any(noder).(base.SpecificNoder[*llmnodesentities.LLMNodeData]); ok {
		if real_node_data, ok := noder.GetNodeData().(*llmnodesentities.LLMNodeData); ok {
			return specific_noder.ExtractVariableSelectorToVariableMapping(graph_config, node_id, real_node_data)
		}
	} else if specific_noder, ok := any(noder).(base.SpecificNoder[*knowledgeretrievalnodesentities.KnowledgeRetrievalNodeData]); ok {
		if real_node_data, ok := noder.GetNodeData().(*knowledgeretrievalnodesentities.KnowledgeRetrievalNodeData); ok {
			return specific_noder.ExtractVariableSelectorToVariableMapping(graph_config, node_id, real_node_data)
		}
	} else if specific_noder, ok := any(noder).(base.SpecificNoder[*ifelsenodesentities.IfElseNodeData]); ok {
		if real_node_data, ok := noder.GetNodeData().(*ifelsenodesentities.IfElseNodeData); ok {
			return specific_noder.ExtractVariableSelectorToVariableMapping(graph_config, node_id, real_node_data)
		}
	} else if specific_noder, ok := any(noder).(base.SpecificNoder[*codenodesentities.CodeNodeData]); ok {
		if real_node_data, ok := noder.GetNodeData().(*codenodesentities.CodeNodeData); ok {
			return specific_noder.ExtractVariableSelectorToVariableMapping(graph_config, node_id, real_node_data)
		}
	} else if specific_noder, ok := any(noder).(base.SpecificNoder[*templatetransformnodesentities.TemplateTransformNodeData]); ok {
		if real_node_data, ok := noder.GetNodeData().(*templatetransformnodesentities.TemplateTransformNodeData); ok {
			return specific_noder.ExtractVariableSelectorToVariableMapping(graph_config, node_id, real_node_data)
		}
	} else if specific_noder, ok := any(noder).(base.SpecificNoder[*questionclassifiernodesentities.QuestionClassifierNodeData]); ok {
		if real_node_data, ok := noder.GetNodeData().(*questionclassifiernodesentities.QuestionClassifierNodeData); ok {
			return specific_noder.ExtractVariableSelectorToVariableMapping(graph_config, node_id, real_node_data)
		}
	} else if specific_noder, ok := any(noder).(base.SpecificNoder[*httprequestnodesentities.HttpRequestNodeData]); ok {
		if real_node_data, ok := noder.GetNodeData().(*httprequestnodesentities.HttpRequestNodeData); ok {
			return specific_noder.ExtractVariableSelectorToVariableMapping(graph_config, node_id, real_node_data)
		}
	} else if specific_noder, ok := any(noder).(base.SpecificNoder[*toolnodesentities.ToolNodeData]); ok {
		if real_node_data, ok := noder.GetNodeData().(*toolnodesentities.ToolNodeData); ok {
			return specific_noder.ExtractVariableSelectorToVariableMapping(graph_config, node_id, real_node_data)
		}
	} else if specific_noder, ok := any(noder).(base.SpecificNoder[*variableaggregatornodesentities.VariableAggregatorNodeData]); ok {
		if real_node_data, ok := noder.GetNodeData().(*variableaggregatornodesentities.VariableAggregatorNodeData); ok {
			return specific_noder.ExtractVariableSelectorToVariableMapping(graph_config, node_id, real_node_data)
		}
	} else if specific_noder, ok := any(noder).(base.SpecificNoder[*loopnodesentities.LoopNodeData]); ok {
		if real_node_data, ok := noder.GetNodeData().(*loopnodesentities.LoopNodeData); ok {
			return specific_noder.ExtractVariableSelectorToVariableMapping(graph_config, node_id, real_node_data)
		}
	} else if specific_noder, ok := any(noder).(base.SpecificNoder[*iterationnodesentities.IterationNodeData]); ok {
		if real_node_data, ok := noder.GetNodeData().(*iterationnodesentities.IterationNodeData); ok {
			return specific_noder.ExtractVariableSelectorToVariableMapping(graph_config, node_id, real_node_data)
		}
	} else if specific_noder, ok := any(noder).(base.SpecificNoder[*parameterextractornodesentities.ParameterExtractorNodeData]); ok {
		if real_node_data, ok := noder.GetNodeData().(*parameterextractornodesentities.ParameterExtractorNodeData); ok {
			return specific_noder.ExtractVariableSelectorToVariableMapping(graph_config, node_id, real_node_data)
		}
	} else if specific_noder, ok := any(noder).(base.SpecificNoder[*variableassignernodesentities.VariableAssignerNodeData]); ok {
		if real_node_data, ok := noder.GetNodeData().(*variableassignernodesentities.VariableAssignerNodeData); ok {
			return specific_noder.ExtractVariableSelectorToVariableMapping(graph_config, node_id, real_node_data)
		}
	} else if specific_noder, ok := any(noder).(base.SpecificNoder[*documentextractornodesentities.DocumentExtractorNodeData]); ok {
		if real_node_data, ok := noder.GetNodeData().(*documentextractornodesentities.DocumentExtractorNodeData); ok {
			return specific_noder.ExtractVariableSelectorToVariableMapping(graph_config, node_id, real_node_data)
		}
	} else if specific_noder, ok := any(noder).(base.SpecificNoder[*listoperatornodesentities.ListOperatorNodeData]); ok {
		if real_node_data, ok := noder.GetNodeData().(*listoperatornodesentities.ListOperatorNodeData); ok {
			return specific_noder.ExtractVariableSelectorToVariableMapping(graph_config, node_id, real_node_data)
		}
	} else if specific_noder, ok := any(noder).(base.SpecificNoder[*humaninputnodesentities.HumanInputNodeData]); ok {
		if real_node_data, ok := noder.GetNodeData().(*humaninputnodesentities.HumanInputNodeData); ok {
			return specific_noder.ExtractVariableSelectorToVariableMapping(graph_config, node_id, real_node_data)
		}
	}
	mlog.Errorf("unsurported noder:%#v", noder)
	panic(exceptions.NewNotImplementedError(fmt.Sprintf("unsurported noder:%#v", noder)))
}

func ExtractVariableSelectorToVariableMappingByNodeType(graph_config map[string]any, node_id string, node_type nodesenumtypes.NodeType) map[string][]string {
	var nodedata any
	switch node_type {
	case nodesenumtypes.Node_START:
		nodedata = new(startnodesentities.StartNodeData)
	case nodesenumtypes.Node_END:
		nodedata = new(endnodesentities.EndNodeData)
	case nodesenumtypes.Node_ANSWER:
		nodedata = new(answernodesentities.AnswerNodeData)
	case nodesenumtypes.Node_LLM:
		nodedata = new(llmnodesentities.LLMNodeData)
	case nodesenumtypes.Node_KNOWLEDGE_RETRIEVAL:
		nodedata = new(knowledgeretrievalnodesentities.KnowledgeRetrievalNodeData)
	case nodesenumtypes.Node_IF_ELSE:
		nodedata = new(ifelsenodesentities.IfElseNodeData)
	case nodesenumtypes.Node_CODE:
		nodedata = new(codenodesentities.CodeNodeData)
	case nodesenumtypes.Node_TEMPLATE_TRANSFORM:
		nodedata = new(templatetransformnodesentities.TemplateTransformNodeData)
	case nodesenumtypes.Node_QUESTION_CLASSIFIER:
		nodedata = new(questionclassifiernodesentities.QuestionClassifierNodeData)
	case nodesenumtypes.Node_HTTP_REQUEST:
		nodedata = new(httprequestnodesentities.HttpRequestNodeData)
	case nodesenumtypes.Node_TOOL:
		nodedata = new(toolnodesentities.ToolNodeData)
	case nodesenumtypes.Node_VARIABLE_AGGREGATOR:
		nodedata = new(variableaggregatornodesentities.VariableAggregatorNodeData)
	case nodesenumtypes.Node_LOOP:
		nodedata = new(loopnodesentities.LoopNodeData)
	case nodesenumtypes.Node_ITERATION:
		nodedata = new(iterationnodesentities.IterationNodeData)
	case nodesenumtypes.Node_PARAMETER_EXTRACTOR:
		nodedata = new(parameterextractornodesentities.ParameterExtractorNodeData)
	case nodesenumtypes.Node_VARIABLE_ASSIGNER:
		nodedata = new(variableassignernodesentities.VariableAssignerNodeData)
	case nodesenumtypes.Node_DOCUMENT_EXTRACTOR:
		nodedata = new(documentextractornodesentities.DocumentExtractorNodeData)
	case nodesenumtypes.Node_LIST_OPERATOR:
		nodedata = new(listoperatornodesentities.ListOperatorNodeData)
	case nodesenumtypes.Node_HUMAN_INPUT:
		nodedata = new(humaninputnodesentities.HumanInputNodeData)
	}
	return ExtractVariableSelectorToVariableMapping(graph_config, node_id, nodedata)
}

func ExtractVariableSelectorToVariableMapping(graph_config map[string]any, node_id string, node_data any) map[string][]string {
	switch data := node_data.(type) {
	case *startnodesentities.StartNodeData:
		if data == nil {
			data = new(startnodesentities.StartNodeData)
		}
		return start.New().ExtractVariableSelectorToVariableMapping(graph_config, node_id, data)

	case *endnodesentities.EndNodeData:
		if data == nil {
			data = new(endnodesentities.EndNodeData)
		}
		return end.New().ExtractVariableSelectorToVariableMapping(graph_config, node_id, data)

	case *answernodesentities.AnswerNodeData:
		if data == nil {
			data = new(answernodesentities.AnswerNodeData)
		}
		return answer.New().ExtractVariableSelectorToVariableMapping(graph_config, node_id, data)
	case *llmnodesentities.LLMNodeData:
		if data == nil {
			data = new(llmnodesentities.LLMNodeData)
		}
		return llm.New().ExtractVariableSelectorToVariableMapping(graph_config, node_id, data)
	case *knowledgeretrievalnodesentities.KnowledgeRetrievalNodeData:
		if data == nil {
			data = new(knowledgeretrievalnodesentities.KnowledgeRetrievalNodeData)
		}
		return knowledgeretrieval.New().ExtractVariableSelectorToVariableMapping(graph_config, node_id, data)
	case *ifelsenodesentities.IfElseNodeData:
		if data == nil {
			data = new(ifelsenodesentities.IfElseNodeData)
		}
		return ifelse.New().ExtractVariableSelectorToVariableMapping(graph_config, node_id, data)
	case *codenodesentities.CodeNodeData:
		if data == nil {
			data = new(codenodesentities.CodeNodeData)
		}
		return code.New().ExtractVariableSelectorToVariableMapping(graph_config, node_id, data)
	case *templatetransformnodesentities.TemplateTransformNodeData:
		if data == nil {
			data = new(templatetransformnodesentities.TemplateTransformNodeData)
		}
		return templatetransform.New().ExtractVariableSelectorToVariableMapping(graph_config, node_id, data)
	case *questionclassifiernodesentities.QuestionClassifierNodeData:
		if data == nil {
			data = new(questionclassifiernodesentities.QuestionClassifierNodeData)
		}
		return questionclassifier.New().ExtractVariableSelectorToVariableMapping(graph_config, node_id, data)
	case *httprequestnodesentities.HttpRequestNodeData:
		if data == nil {
			data = new(httprequestnodesentities.HttpRequestNodeData)
		}
		return httprequest.New().ExtractVariableSelectorToVariableMapping(graph_config, node_id, data)
	case *toolnodesentities.ToolNodeData:
		if data == nil {
			data = new(toolnodesentities.ToolNodeData)
		}
		return tool.New().ExtractVariableSelectorToVariableMapping(graph_config, node_id, data)
	case *variableaggregatornodesentities.VariableAggregatorNodeData:
		if data == nil {
			data = new(variableaggregatornodesentities.VariableAggregatorNodeData)
		}
		return variableaggregator.New().ExtractVariableSelectorToVariableMapping(graph_config, node_id, data)
	case *variableassignernodesentities.VariableAssignerNodeData:
		if data == nil {
			data = new(variableassignernodesentities.VariableAssignerNodeData)
		}
		return variableassigner.New().ExtractVariableSelectorToVariableMapping(graph_config, node_id, data)
	case *loopnodesentities.LoopNodeData:
		if data == nil {
			data = new(loopnodesentities.LoopNodeData)
		}
		return loop.New().ExtractVariableSelectorToVariableMapping(graph_config, node_id, data)
	case *iterationnodesentities.IterationNodeData:
		if data == nil {
			data = new(iterationnodesentities.IterationNodeData)
		}
		return iteration.New().ExtractVariableSelectorToVariableMapping(graph_config, node_id, data)
	case *parameterextractornodesentities.ParameterExtractorNodeData:
		if data == nil {
			data = new(parameterextractornodesentities.ParameterExtractorNodeData)
		}
		return parameterextractor.New().ExtractVariableSelectorToVariableMapping(graph_config, node_id, data)
	case *documentextractornodesentities.DocumentExtractorNodeData:
		if data == nil {
			data = new(documentextractornodesentities.DocumentExtractorNodeData)
		}
		return documentextractor.New().ExtractVariableSelectorToVariableMapping(graph_config, node_id, data)
	case *listoperatornodesentities.ListOperatorNodeData:
		if data == nil {
			data = new(listoperatornodesentities.ListOperatorNodeData)
		}
		return listoperator.New().ExtractVariableSelectorToVariableMapping(graph_config, node_id, data)
	case *humaninputnodesentities.HumanInputNodeData:
		if data == nil {
			data = new(humaninputnodesentities.HumanInputNodeData)
		}
		return humaninput.New().ExtractVariableSelectorToVariableMapping(graph_config, node_id, data)
	}
	mlog.Errorf("unsurported node_data:%#v", node_data)
	panic(exceptions.NewNotImplementedError(fmt.Sprintf("unsurported node_data:%#v", node_data)))
}

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
	case nodesenumtypes.Node_HUMAN_INPUT:
		ndata := humaninputnodesentities.New()
		bn := base.NewBaseNode(id, config, graphInitParams, graph, graphRuntimeState, previousNodeID, ndata)
		n := &humaninput.HumanInputNode{
			BaseNode: bn,
		}
		return n
	}
	mlog.Errorf("unsurported node_type(%v)", nodeType)
	panic(exceptions.NewNotImplementedError(fmt.Sprintf("unsurported node_type(%v)", nodeType)))
}
