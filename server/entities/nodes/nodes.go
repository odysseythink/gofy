package nodes

import (
	"mlib.com/gofy/server/entities/nodes/answer"
	"mlib.com/gofy/server/entities/nodes/base"
	"mlib.com/gofy/server/entities/nodes/code"
	documentextractor "mlib.com/gofy/server/entities/nodes/document_extractor"
	"mlib.com/gofy/server/entities/nodes/end"
	httprequest "mlib.com/gofy/server/entities/nodes/http_request"
	ifelse "mlib.com/gofy/server/entities/nodes/if_else"
	"mlib.com/gofy/server/entities/nodes/iteration"
	knowledgeretrieval "mlib.com/gofy/server/entities/nodes/knowledge_retrieval"
	listoperator "mlib.com/gofy/server/entities/nodes/list_operator"
	"mlib.com/gofy/server/entities/nodes/llm"
	"mlib.com/gofy/server/entities/nodes/loop"
	parameterextractor "mlib.com/gofy/server/entities/nodes/parameter_extractor"
	questionclassifier "mlib.com/gofy/server/entities/nodes/question_classifier"
	"mlib.com/gofy/server/entities/nodes/start"
	templatetransform "mlib.com/gofy/server/entities/nodes/template_transform"
	"mlib.com/gofy/server/entities/nodes/tool"
	variableaggregator "mlib.com/gofy/server/entities/nodes/variable_aggregator"
	variableassigner "mlib.com/gofy/server/entities/nodes/variable_assigner"
	nodesenumtypes "mlib.com/gofy/server/enum_types/nodes"
	"mlib.com/mlog"
)

// GenericNodeData represents the generic node data type.
type GenericNodeData interface {
	*start.StartNodeData | *end.EndNodeData | *answer.AnswerNodeData | *code.CodeNodeData | *documentextractor.DocumentExtractorNodeData | *httprequest.HttpRequestNodeData | *ifelse.IfElseNodeData | *iteration.IterationNodeData | *knowledgeretrieval.KnowledgeRetrievalNodeData | *listoperator.ListOperatorNodeData | *llm.LLMNodeData | *loop.LoopNodeData | *loop.LoopStartNodeData | *parameterextractor.ParameterExtractorNodeData | *templatetransform.TemplateTransformNodeData | *variableaggregator.VariableAggregatorNodeData | *variableassigner.VariableAssignerNodeData | *questionclassifier.QuestionClassifierNodeData | *tool.ToolNodeData | *base.BaseNodeData
}

func NewNodeDataByNodeType(node_type nodesenumtypes.NodeType) any {
	switch node_type {
	case nodesenumtypes.Node_START:
		return start.New()
	case nodesenumtypes.Node_END:
		return end.New()
	case nodesenumtypes.Node_ANSWER:
		return answer.New()
	case nodesenumtypes.Node_LLM:
		return llm.New()
	case nodesenumtypes.Node_KNOWLEDGE_RETRIEVAL:
		return new(knowledgeretrieval.KnowledgeRetrievalNodeData)
	case nodesenumtypes.Node_IF_ELSE:
		return ifelse.New()
	case nodesenumtypes.Node_CODE:
		return code.New()
	case nodesenumtypes.Node_TEMPLATE_TRANSFORM:
		return new(templatetransform.TemplateTransformNodeData)
	case nodesenumtypes.Node_QUESTION_CLASSIFIER:
		return questionclassifier.New()
	case nodesenumtypes.Node_HTTP_REQUEST:
		return httprequest.New()
	case nodesenumtypes.Node_TOOL:
		return new(tool.ToolNodeData)
	case nodesenumtypes.Node_VARIABLE_AGGREGATOR:
		return new(variableaggregator.VariableAggregatorNodeData)
	case nodesenumtypes.Node_LOOP:
		return new(loop.LoopNodeData)
	case nodesenumtypes.Node_LOOP_START:
		return new(loop.LoopStartNodeData)
	case nodesenumtypes.Node_ITERATION:
		return new(iteration.IterationNodeData)
	case nodesenumtypes.Node_PARAMETER_EXTRACTOR:
		return parameterextractor.New()
	case nodesenumtypes.Node_VARIABLE_ASSIGNER:
		return new(variableassigner.VariableAssignerNodeData)
	case nodesenumtypes.Node_DOCUMENT_EXTRACTOR:
		return new(documentextractor.DocumentExtractorNodeData)
	case nodesenumtypes.Node_LIST_OPERATOR:
		return new(listoperator.ListOperatorNodeData)
	}
	mlog.Errorf("unsurported node_type=%v", node_type)
	return nil
}
