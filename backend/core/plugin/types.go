package plugin

// ToolResponseChunkType represents the type of a tool response chunk.
type ToolResponseChunkType string

const (
	ToolResponseChunkTypeText               ToolResponseChunkType = "text"
	ToolResponseChunkTypeFile               ToolResponseChunkType = "file"
	ToolResponseChunkTypeBlob               ToolResponseChunkType = "blob"
	ToolResponseChunkTypeBlobChunk          ToolResponseChunkType = "blob_chunk"
	ToolResponseChunkTypeJson               ToolResponseChunkType = "json"
	ToolResponseChunkTypeLink               ToolResponseChunkType = "link"
	ToolResponseChunkTypeImage              ToolResponseChunkType = "image"
	ToolResponseChunkTypeImageLink          ToolResponseChunkType = "image_link"
	ToolResponseChunkTypeVariable           ToolResponseChunkType = "variable"
	ToolResponseChunkTypeLog                ToolResponseChunkType = "log"
	ToolResponseChunkTypeRetrieverResources ToolResponseChunkType = "retriever_resources"
)

// IsValidToolResponseChunkType checks whether a string is a valid chunk type.
func IsValidToolResponseChunkType(t string) bool {
	switch ToolResponseChunkType(t) {
	case ToolResponseChunkTypeText,
		ToolResponseChunkTypeFile,
		ToolResponseChunkTypeBlob,
		ToolResponseChunkTypeBlobChunk,
		ToolResponseChunkTypeJson,
		ToolResponseChunkTypeLink,
		ToolResponseChunkTypeImage,
		ToolResponseChunkTypeImageLink,
		ToolResponseChunkTypeVariable,
		ToolResponseChunkTypeLog,
		ToolResponseChunkTypeRetrieverResources:
		return true
	default:
		return false
	}
}

// ToolResponseChunk represents a single chunk in a streaming tool response.
type ToolResponseChunk struct {
	Type    ToolResponseChunkType `json:"type"`
	Message map[string]any        `json:"message"`
	Meta    map[string]any        `json:"meta"`
}
