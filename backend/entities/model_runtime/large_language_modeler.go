package modelruntime

import "iter"

type LargeLanguageModeler interface {
	AIModeler
	GetNumTokens(
		model string,
		credentials map[string]any,
		prompt_messages []PromptMessager,
		tools []*PromptMessageTool,
	) int
	Invoke(
		model string,
		credentials map[string]any,
		prompt_messages []PromptMessager,
		model_parameters map[string]any,
		tools []*PromptMessageTool,
		stop []string,
		user string,
	) *LLMResult /*-> Union[LLMResult, Generator]*/
	InvokeStream(
		model string,
		credentials map[string]any,
		prompt_messages []PromptMessager,
		model_parameters map[string]any,
		tools []*PromptMessageTool,
		stop []string,
		user string,
	) iter.Seq[*LLMResultChunk] /*-> Union[LLMResult, Generator]*/
	GetParameterRules(modeler LargeLanguageModeler, model string, credentials map[string]any) []*ParameterRule
	CalcResponseUsage(LargeLanguageModeler, string, map[string]any, int, int) *LLMUsage
}
