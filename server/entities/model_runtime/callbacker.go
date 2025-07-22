package modelruntime

type Callbacker interface {
	OnBeforeInvoke(
		llm_instance AIModeler,
		model string,
		credentials map[string]any,
		prompt_messages []PromptMessager,
		model_parameters map[string]any,
		tools []*PromptMessageTool,
		stop []string,
		stream bool, /*true*/
		user string,
	)
	OnNewChunk(
		llm_instance AIModeler,
		chunk *LLMResultChunk,
		model string,
		credentials map[string]any,
		prompt_messages []PromptMessager,
		model_parameters map[string]any,
		tools []*PromptMessageTool,
		stop []string,
		stream bool, /*true*/
		user string,
	)

	OnAfterInvoke(
		llm_instance AIModeler,
		result *LLMResult,
		model string,
		credentials map[string]any,
		prompt_messages []PromptMessager,
		model_parameters map[string]any,
		tools []*PromptMessageTool,
		stop []string,
		stream bool, /*true*/
		user string,
	)
	OnInvokeError(
		llm_instance AIModeler,
		err error,
		model string,
		credentials map[string]any,
		prompt_messages []PromptMessager,
		model_parameters map[string]any,
		tools []*PromptMessageTool,
		stop []string,
		stream bool, /*true*/
		user string,
	)
}
