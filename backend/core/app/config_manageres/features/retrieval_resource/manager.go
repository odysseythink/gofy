package retrievalresource

// 定义RetrievalResourceConfigManager结构体
type RetrievalResourceConfigManager struct{}

// Convert 方法
func (m *RetrievalResourceConfigManager) Convert(config map[string]any) bool {
	retrieverResourceDict, ok := config["retriever_resource"].(map[string]any)
	if !ok {
		return false
	}

	if enabled, ok := retrieverResourceDict["enabled"].(bool); ok {
		return enabled
	}

	return false
}
