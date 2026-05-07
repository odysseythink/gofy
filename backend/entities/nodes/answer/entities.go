package answer

import (
	basenodesentities "github.com/odysseythink/gofy/backend/entities/nodes/base"

	answernodesenumtypes "github.com/odysseythink/gofy/backend/enum_types/nodes/answer"
)

// AnswerNodeData represents answer node data
type AnswerNodeData struct {
	*basenodesentities.BaseNodeData
	Answer string `json:"answer"`
}

// GenerateRouteChunk represents a chunk for generating routes
type GenerateRouteChunker interface {
	Type() answernodesenumtypes.GenerateRouteChunkType
}

// VarGenerateRouteChunk represents a variable chunk for generating routes
type VarGenerateRouteChunk struct {
	ValueSelector []string `json:"value_selector"`
}

func (rc *VarGenerateRouteChunk) Type() answernodesenumtypes.GenerateRouteChunkType {
	return answernodesenumtypes.GenerateRouteChunk_VAR
}

// TextGenerateRouteChunk represents a text chunk for generating routes
type TextGenerateRouteChunk struct {
	Text string `json:"text"`
}

func (rc *TextGenerateRouteChunk) Type() answernodesenumtypes.GenerateRouteChunkType {
	return answernodesenumtypes.GenerateRouteChunk_TEXT
}

// AnswerNodeDoubleLink represents a double link for answer nodes
type AnswerNodeDoubleLink struct {
	NodeID        string   `json:"node_id"`
	SourceNodeIDs []string `json:"source_node_ids"`
	TargetNodeIDs []string `json:"target_node_ids"`
}

// AnswerStreamGenerateRoute represents answer stream generate route
type AnswerStreamGenerateRoute struct {
	AnswerDependencies  map[string][]string               `json:"answer_dependencies"`
	AnswerGenerateRoute map[string][]GenerateRouteChunker `json:"answer_generate_route"`
}

func New() *AnswerNodeData {
	return &AnswerNodeData{
		BaseNodeData: &basenodesentities.BaseNodeData{},
	}
}
