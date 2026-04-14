package callbacks

import (
	"fmt"

	"github.com/odysseythink/mlog"
	modelruntimeentities "mlib.com/gofy/server/entities/model_runtime"
)

type LoggingCallback struct {
	*BaseCallback
}

func NewLoggingCallback() *LoggingCallback {
	return &LoggingCallback{
		BaseCallback: &BaseCallback{},
	}
}

func (cb *LoggingCallback) OnBeforeInvoke(
	llm_instance modelruntimeentities.AIModeler,
	model string,
	credentials map[string]any,
	prompt_messages []modelruntimeentities.PromptMessager,
	model_parameters map[string]any,
	tools []*modelruntimeentities.PromptMessageTool,
	stop []string,
	stream bool, /*true*/
	user string,
) {
	/*
		Before invoke callback

		:param llm_instance: LLM instance
		:param model: model name
		:param credentials: model credentials
		:param prompt_messages: prompt messages
		:param model_parameters: model parameters
		:param tools: tools for tool calling
		:param stop: stop words
		:param stream: is stream response
		:param user: unique user id
	*/
	cb.PrintText("\n[on_llm_before_invoke]\n", "blue")
	cb.PrintText(fmt.Sprintf("Model: %s\n", model), "blue")
	cb.PrintText("Parameters:\n", "blue")
	for key, value := range model_parameters {
		cb.PrintText(fmt.Sprintf("\t%s: %v\n", key, value), "blue")
	}
	if len(stop) > 0 {
		cb.PrintText(fmt.Sprintf("\tstop: %v\n", stop), "blue")
	}
	if len(tools) > 0 {
		cb.PrintText("\tTools:\n", "blue")
		for _, tool := range tools {
			cb.PrintText(fmt.Sprintf("\t\t%s\n", tool.Name), "blue")
		}
	}
	cb.PrintText(fmt.Sprintf("Stream: %v\n", stream), "blue")

	if user != "" {
		cb.PrintText(fmt.Sprintf("User: %s\n", user), "blue")
	}
	cb.PrintText("Prompt messages:\n", "blue")
	for _, prompt_message := range prompt_messages {
		if prompt_message.Name() != "" {
			cb.PrintText(fmt.Sprintf("\tname: %s\n", prompt_message.Name()), "blue")
		}
		cb.PrintText(fmt.Sprintf("\trole: %s\n", prompt_message.Role()), "blue")
		cb.PrintText(fmt.Sprintf("\tcontent: %v\n", prompt_message.GetContent()), "blue")
	}
	if stream {
		cb.PrintText("\n[on_llm_new_chunk]", "blue")
	}
}
func (cb *LoggingCallback) OnNewChunk(
	llm_instance modelruntimeentities.AIModeler,
	chunk *modelruntimeentities.LLMResultChunk,
	model string,
	credentials map[string]any,
	prompt_messages []modelruntimeentities.PromptMessager,
	model_parameters map[string]any,
	tools []*modelruntimeentities.PromptMessageTool,
	stop []string,
	stream bool, /*true*/
	user string,
) {
	/*
		On new chunk callback

		:param llm_instance: LLM instance
		:param chunk: chunk
		:param model: model name
		:param credentials: model credentials
		:param prompt_messages: prompt messages
		:param model_parameters: model parameters
		:param tools: tools for tool calling
		:param stop: stop words
		:param stream: is stream response
		:param user: unique user id
	*/
	fmt.Println(chunk.Delta.Message.Content)
}
func (cb *LoggingCallback) OnAfterInvoke(
	llm_instance modelruntimeentities.AIModeler,
	result *modelruntimeentities.LLMResult,
	model string,
	credentials map[string]any,
	prompt_messages []modelruntimeentities.PromptMessager,
	model_parameters map[string]any,
	tools []*modelruntimeentities.PromptMessageTool,
	stop []string,
	stream bool, /*true*/
	user string,
) {
	/*
		After invoke callback

		:param llm_instance: LLM instance
		:param result: result
		:param model: model name
		:param credentials: model credentials
		:param prompt_messages: prompt messages
		:param model_parameters: model parameters
		:param tools: tools for tool calling
		:param stop: stop words
		:param stream: is stream response
		:param user: unique user id
	*/
	cb.PrintText("\n[on_llm_after_invoke]\n", "yellow")
	cb.PrintText(fmt.Sprintf("Content: %s\n", result.Message.Content), "yellow")

	if len(result.Message.ToolCalls) > 0 {
		cb.PrintText("Tool calls:\n", "yellow")
		for _, tool_call := range result.Message.ToolCalls {
			cb.PrintText(fmt.Sprintf("\t%s\n", tool_call.ID), "yellow")
			cb.PrintText(fmt.Sprintf("\t%s\n", tool_call.Function.Name), "yellow")
			cb.PrintText(fmt.Sprintf("\t%s}\n", tool_call.Function.Arguments), "yellow")
		}
	}
	cb.PrintText(fmt.Sprintf("Model: %s\n", result.Model), "yellow")
	cb.PrintText(fmt.Sprintf("Usage: %s\n", result.Usage), "yellow")
	cb.PrintText(fmt.Sprintf("System Fingerprint: %s\n", result.SystemFingerprint), "yellow")
}
func (cb *LoggingCallback) OnInvokeError(
	llm_instance modelruntimeentities.AIModeler,
	err error,
	model string,
	credentials map[string]any,
	prompt_messages []modelruntimeentities.PromptMessager,
	model_parameters map[string]any,
	tools []*modelruntimeentities.PromptMessageTool,
	stop []string,
	stream bool, /*true*/
	user string,
) {
	/*
		Invoke error callback

		:param llm_instance: LLM instance
		:param ex: exception
		:param model: model name
		:param credentials: model credentials
		:param prompt_messages: prompt messages
		:param model_parameters: model parameters
		:param tools: tools for tool calling
		:param stop: stop words
		:param stream: is stream response
		:param user: unique user id
	*/
	cb.PrintText("\n[on_llm_invoke_error]\n", "red")
	mlog.Error(err)
}
