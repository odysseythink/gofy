package nodes

import (
	"fmt"
)

type NodeType string

const (
	Node_START                      NodeType = "start"
	Node_END                        NodeType = "end"
	Node_ANSWER                     NodeType = "answer"
	Node_LLM                        NodeType = "llm"
	Node_KNOWLEDGE_RETRIEVAL        NodeType = "knowledge-retrieval"
	Node_IF_ELSE                    NodeType = "if-else"
	Node_CODE                       NodeType = "code"
	Node_TEMPLATE_TRANSFORM         NodeType = "template-transform"
	Node_QUESTION_CLASSIFIER        NodeType = "question-classifier"
	Node_HTTP_REQUEST               NodeType = "http-request"
	Node_TOOL                       NodeType = "tool"
	Node_VARIABLE_AGGREGATOR        NodeType = "variable-aggregator"
	Node_LEGACY_VARIABLE_AGGREGATOR NodeType = "variable-assigner" // TODO: Merge this into VARIABLE_AGGREGATOR in the database.
	Node_LOOP                       NodeType = "loop"
	Node_LOOP_START                 NodeType = "loop-start" // Fake start node for loop.
	Node_LOOP_END                   NodeType = "loop-end"
	Node_ITERATION                  NodeType = "iteration"
	Node_ITERATION_START            NodeType = "iteration-start" // Fake start node for iteration.
	Node_PARAMETER_EXTRACTOR        NodeType = "parameter-extractor"
	Node_VARIABLE_ASSIGNER          NodeType = "assigner"
	Node_DOCUMENT_EXTRACTOR         NodeType = "document-extractor"
	Node_LIST_OPERATOR              NodeType = "list-operator"
)

func (n NodeType) Valid() bool {
	return n == Node_START ||
		n == Node_END ||
		n == Node_ANSWER ||
		n == Node_LLM ||
		n == Node_KNOWLEDGE_RETRIEVAL ||
		n == Node_IF_ELSE ||
		n == Node_CODE ||
		n == Node_TEMPLATE_TRANSFORM ||
		n == Node_QUESTION_CLASSIFIER ||
		n == Node_HTTP_REQUEST ||
		n == Node_TOOL ||
		n == Node_VARIABLE_AGGREGATOR ||
		n == Node_LEGACY_VARIABLE_AGGREGATOR ||
		n == Node_LOOP ||
		n == Node_LOOP_START ||
		n == Node_ITERATION ||
		n == Node_ITERATION_START ||
		n == Node_PARAMETER_EXTRACTOR ||
		n == Node_VARIABLE_ASSIGNER ||
		n == Node_DOCUMENT_EXTRACTOR ||
		n == Node_LIST_OPERATOR
}

// ParseNodeType returns the case-insensitive NodeType value for the given string.
func ParseNodeType(name string) (NodeType, error) {
	if !NodeType(name).Valid() {
		return NodeType(name), fmt.Errorf("logsink: invalid NodeType %q", name)
	}
	return NodeType(name), nil
}

func ParseNodeTypeWithoutError(name string) NodeType {

	return NodeType(name)
}

type ErrorStrategy string

const (
	ErrorStrategy_FAIL_BRANCH   ErrorStrategy = "fail-branch"
	ErrorStrategy_DEFAULT_VALUE ErrorStrategy = "default-value"
)

type FailBranchSourceHandle string

const (
	FailBranchSourceHandle_FAILED  FailBranchSourceHandle = "fail-branch"
	FailBranchSourceHandle_SUCCESS FailBranchSourceHandle = "success-branch"
)

var (
	CONTINUE_ON_ERROR_NODE_TYPE = []NodeType{Node_LLM, Node_CODE, Node_TOOL, Node_HTTP_REQUEST}
	RETRY_ON_ERROR_NODE_TYPE    = CONTINUE_ON_ERROR_NODE_TYPE
)
