package nodes

import (
	"github.com/odysseythink/gofy/backend/entities/nodes/answer"
	"github.com/odysseythink/gofy/backend/entities/nodes/base"
	"github.com/odysseythink/gofy/backend/entities/nodes/code"
	documentextractor "github.com/odysseythink/gofy/backend/entities/nodes/document_extractor"
	"github.com/odysseythink/gofy/backend/entities/nodes/end"
	httprequest "github.com/odysseythink/gofy/backend/entities/nodes/http_request"
	humaninput "github.com/odysseythink/gofy/backend/entities/nodes/human_input"
	ifelse "github.com/odysseythink/gofy/backend/entities/nodes/if_else"
	"github.com/odysseythink/gofy/backend/entities/nodes/iteration"
	knowledgeretrieval "github.com/odysseythink/gofy/backend/entities/nodes/knowledge_retrieval"
	listoperator "github.com/odysseythink/gofy/backend/entities/nodes/list_operator"
	"github.com/odysseythink/gofy/backend/entities/nodes/llm"
	"github.com/odysseythink/gofy/backend/entities/nodes/loop"
	parameterextractor "github.com/odysseythink/gofy/backend/entities/nodes/parameter_extractor"
	questionclassifier "github.com/odysseythink/gofy/backend/entities/nodes/question_classifier"
	"github.com/odysseythink/gofy/backend/entities/nodes/start"
	templatetransform "github.com/odysseythink/gofy/backend/entities/nodes/template_transform"
	"github.com/odysseythink/gofy/backend/entities/nodes/tool"
	variableaggregator "github.com/odysseythink/gofy/backend/entities/nodes/variable_aggregator"
	variableassigner "github.com/odysseythink/gofy/backend/entities/nodes/variable_assigner"
	nodesenumtypes "github.com/odysseythink/gofy/backend/enum_types/nodes"
	"github.com/odysseythink/mlog"
)

// GenericNodeData represents the generic node data type.
type GenericNodeData interface {
	*start.StartNodeData | *end.EndNodeData | *answer.AnswerNodeData | *code.CodeNodeData | *documentextractor.DocumentExtractorNodeData | *httprequest.HttpRequestNodeData | *humaninput.HumanInputNodeData | *ifelse.IfElseNodeData | *iteration.IterationNodeData | *knowledgeretrieval.KnowledgeRetrievalNodeData | *listoperator.ListOperatorNodeData | *llm.LLMNodeData | *loop.LoopNodeData | *parameterextractor.ParameterExtractorNodeData | *templatetransform.TemplateTransformNodeData | *variableaggregator.VariableAggregatorNodeData | *variableassigner.VariableAssignerNodeData | *questionclassifier.QuestionClassifierNodeData | *tool.ToolNodeData | *base.BaseNodeData
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
	case nodesenumtypes.Node_HUMAN_INPUT:
		return humaninput.New()
	}
	mlog.Errorf("unsurported node_type=%v", node_type)
	return nil
}
