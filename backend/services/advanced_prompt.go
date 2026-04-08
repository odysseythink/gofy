package services

type AdvancedPromptTemplateService struct{}

func (s *AdvancedPromptTemplateService) GetDefaultPromptConfig(appMode string, modelMode string) map[string]any {
	if appMode == "chat" {
		return map[string]any{
			"prompt": []map[string]any{
				{"role": "system", "text": ""},
			},
		}
	}
	return map[string]any{
		"prompt": map[string]any{"text": ""},
		"conversation_histories_role": map[string]any{
			"user_prefix":      "",
			"assistant_prefix": "",
		},
	}
}

func (s *AdvancedPromptTemplateService) GetDefaultAdvancedPromptConfig(appMode, modelMode string) map[string]any {
	return s.GetDefaultPromptConfig(appMode, modelMode)
}
