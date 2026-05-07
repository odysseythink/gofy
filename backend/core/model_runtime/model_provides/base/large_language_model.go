package base

import (
	"fmt"
	"iter"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/odysseythink/gofy/backend/core/exceptions"
	modelruntimeentities "github.com/odysseythink/gofy/backend/entities/model_runtime"
	modelruntimeenumtypes "github.com/odysseythink/gofy/backend/enum_types/model_runtime"
	"github.com/odysseythink/mlog"
)

const (
	HTML_THINKING_TAG = `<details style="color:gray;background-color: #f8f8f8;padding: 8px;border-radius: 4px;" open> 
<summary> Thinking... </summary>`
)

type LargeLanguageModel struct {
	*BaseAIModel
}

func (m *LargeLanguageModel) ModelType() modelruntimeenumtypes.ModelType {
	return modelruntimeenumtypes.Model_LLM
}

// func (m *LargeLanguageModel) CodeBlockModeWrapper(
// 	modeler modelruntimeentities.LargeLanguageModeler,
// 	model string,
// 	credentials map[string]any,
// 	prompt_messages []modelruntimeentities.PromptMessager,
// 	model_parameters map[string]any,
// 	tools []*modelruntimeentities.PromptMessageTool,
// 	stop []string,
// 	stream bool, /* = True*/
// 	user string,
// 	callbacks []modelruntimeentities.Callbacker,
// ) (any, error) /*-> Union[LLMResult, Generator]*/ {
// 	/*
// 		Code block mode wrapper, ensure the response is a code block with output markdown quote

// 		:param model: model name
// 		:param credentials: model credentials
// 		:param prompt_messages: prompt messages
// 		:param model_parameters: model parameters
// 		:param tools: tools for tool calling
// 		:param stop: stop words
// 		:param stream: is stream response
// 		:param user: unique user id
// 		:param callbacks: callbacks
// 		:return: full response or stream response chunk generator result
// 	*/

// 	block_prompts := `You should always follow the instructions and output a valid {{block}} object.
// The structure of the {{block}} object you can found in the instructions, use {"answer": "$your_answer"} as the default structure
// if you are not sure about the structure.

// <instructions>
// {{instructions}}
// </instructions>
// ` // noqa: E501
// 	code_block := ""
// 	if _, ok := model_parameters["response_format"]; ok {
// 		if _, ok := model_parameters["response_format"].(string); ok {
// 			code_block = model_parameters["response_format"].(string)
// 		}
// 	}

// 	if code_block == "" {
// 		return modeler.Invoke(
// 			model,
// 			credentials,
// 			prompt_messages,
// 			model_parameters,
// 			tools,
// 			stop,
// 			stream,
// 			user,
// 		)
// 	}
// 	delete(model_parameters, "response_format")

// 	if len(stop) == 0 {
// 		stop = make([]string, 0)
// 	}
// 	stop = append(stop, []string{"\n```", "```\n"}...)
// 	block_prompts = strings.ReplaceAll(block_prompts, "{{block}}", code_block)

// 	// check if there is a system message
// 	if len(prompt_messages) > 0 && prompt_messages[0].Role() == modelruntimeentities.PromptMessageRole_SYSTEM {
// 		// override the system message
// 		realpromptmsg := prompt_messages[0].(*modelruntimeentities.SystemPromptMessage)
// 		prompt_messages[0] = modelruntimeentities.NewSystemPromptMessage(
// 			strings.ReplaceAll(block_prompts, "{{instructions}}", realpromptmsg.Content), "",
// 		)
// 	} else {
// 		// insert the system message
// 		prompt_messages = append([]modelruntimeentities.PromptMessager{modelruntimeentities.NewSystemPromptMessage(strings.ReplaceAll(block_prompts, "{{instructions}}", fmt.Sprintf("Please output a valid %s object.", code_block)), "")}, prompt_messages...)
// 	}

// 	if len(prompt_messages) > 0 && prompt_messages[len(prompt_messages)-1].Role() == modelruntimeentities.PromptMessageRole_USER {
// 		// add ```JSON\n to the last text message

// 		last_prompt_message := prompt_messages[len(prompt_messages)-1]
// 		if realpromptmsg, ok := any(last_prompt_message).(*modelruntimeentities.UserPromptMessage[string]); ok {
// 			realpromptmsg.Content += fmt.Sprintf("\n```%s\n", code_block)
// 			prompt_messages[len(prompt_messages)-1] = realpromptmsg
// 		} else if realpromptmsg, ok := any(last_prompt_message).(*modelruntimeentities.UserPromptMessage[[]modelruntimeentities.PromptMessageContenter]); ok {
// 			for i := len(realpromptmsg.Content); i >= 0; i-- {
// 				if realpromptmsg.Content[i].Type() == modelruntimeentities.PromptMessageContent_TEXT {
// 					realval := any(realpromptmsg.Content[i]).(*modelruntimeentities.TextPromptMessageContent)
// 					realval.Data += fmt.Sprintf("\n```%s\n", code_block)
// 					realpromptmsg.Content[i] = realval
// 					break
// 				}
// 			}
// 			prompt_messages[len(prompt_messages)-1] = realpromptmsg
// 		}
// 	} else {
// 		// append a user message
// 		prompt_messages = append(prompt_messages, modelruntimeentities.NewUserPromptMessage(fmt.Sprintf("```%s\n", code_block), ""))
// 	}
// 	response, err := modeler.Invoke(
// 		model,
// 		credentials,
// 		prompt_messages,
// 		model_parameters,
// 		tools,
// 		stop,
// 		stream,
// 		user,
// 	)
// 	if err != nil {
// 		mlog.Errorf("invoke failed:%v", err)
// 		return nil, err
// 	}

// 	if _, ok := response.(iter.Seq[*modelruntimeentities.LLMResultChunk]); ok {
// 		return func(yield func(*modelruntimeentities.LLMResultChunk) bool) {
// 			var first_chunk *modelruntimeentities.LLMResultChunk
// 			state := ""
// 			backtick_count := 0
// 			content_piece := ""
// 			for piece := range response.(iter.Seq[*modelruntimeentities.LLMResultChunk]) {
// 				if first_chunk == nil {
// 					first_chunk = piece
// 					if first_chunk.Delta.Message.Content != "" && strings.HasPrefix(first_chunk.Delta.Message.Content, "`") {
// 						state = "search_start"
// 					} else {
// 						state = "normal"
// 					}
// 				}
// 				if first_chunk.Delta.Message.Content != "" && strings.HasPrefix(first_chunk.Delta.Message.Content, "`") {
// 					if piece.Delta.Message.Content != "" {
// 						content := piece.Delta.Message.Content
// 						// Reset content to ensure we're only processing and yielding the relevant parts
// 						piece.Delta.Message.Content = ""
// 						// Yield a piece with cleared content before processing it to maintain the generator structure
// 						yield(piece)
// 						content_piece = content
// 					} else {
// 						// Yield pieces without content directly
// 						yield(piece)
// 						continue
// 					}
// 					if state == "done" {
// 						continue
// 					}
// 					new_piece := ""
// 					for _, char := range []byte(content_piece) {
// 						if state == "search_start" {
// 							if string(char) == "`" {
// 								backtick_count += 1
// 								if backtick_count == 3 {
// 									state = "skip_language"
// 									backtick_count = 0
// 								}
// 							} else {
// 								backtick_count = 0
// 							}
// 						} else if state == "skip_language" {
// 							// Skip everything until the first newline, marking the end of the language identifier
// 							if string(char) == "\n" {
// 								state = "in_code_block"
// 							}
// 						} else if state == "in_code_block" {
// 							if string(char) == "`" {
// 								backtick_count += 1
// 								if backtick_count == 3 {
// 									state = "done"
// 									break
// 								}
// 							} else {
// 								if backtick_count > 0 {
// 									// If backticks were counted but we're still collecting content, it was a false start
// 									new_piece += strings.Repeat("`", backtick_count)
// 									backtick_count = 0
// 								}
// 								new_piece += string(char)
// 							}
// 						} else if state == "done" {
// 							break
// 						}
// 						if new_piece != "" {
// 							// Only yield content collected within the code block
// 							yield(&modelruntimeentities.LLMResultChunk{
// 								Model:          model,
// 								PromptMessages: prompt_messages,
// 								Delta: &modelruntimeentities.LLMResultChunkDelta{
// 									Index:   0,
// 									Message: modelruntimeentities.NewAssistantPromptMessage(new_piece, "", nil),
// 								},
// 							})
// 						}
// 					}
// 				} else {
// 					if piece.Delta.Message.Content != "" {
// 						content := piece.Delta.Message.Content
// 						piece.Delta.Message.Content = ""
// 						if !yield(piece) {
// 							return
// 						}
// 						content_piece = content
// 					} else {
// 						if !yield(piece) {
// 							return
// 						}
// 						continue
// 					}
// 					new_piece := ""
// 					for _, c := range []byte(content_piece) {
// 						char := string(c)
// 						if state == "normal" {
// 							if char == "`" {
// 								state = "in_backticks"
// 								backtick_count = 1
// 							} else {
// 								new_piece += char
// 							}
// 						} else if state == "in_backticks" {
// 							if char == "`" {
// 								backtick_count += 1
// 								if backtick_count == 3 {
// 									state = "skip_content"
// 									backtick_count = 0
// 								}
// 							} else {
// 								new_piece += strings.Repeat("`", backtick_count) + char
// 								state = "normal"
// 								backtick_count = 0
// 							}
// 						} else if state == "skip_content" {
// 							if char == " " {
// 								state = "normal"
// 							}
// 						}
// 						if new_piece != "" {
// 							yield(&modelruntimeentities.LLMResultChunk{
// 								Model:          model,
// 								PromptMessages: prompt_messages,
// 								Delta: &modelruntimeentities.LLMResultChunkDelta{
// 									Index:   0,
// 									Message: modelruntimeentities.NewAssistantPromptMessage(new_piece, "", nil),
// 								},
// 							})
// 						}
// 					}
// 				}
// 			}

//			}, nil
//		} else if _, ok := response.(*modelruntimeentities.LLMResult); ok {
//			return response, nil
//		} else {
//			return nil, exceptions.NewValueError(fmt.Sprintf("unsurport invoke response(%#v)", response))
//		}
//	}
func (m *LargeLanguageModel) CodeBlockModeStreamProcessor(
	model string, prompt_messages []modelruntimeentities.PromptMessager, input_generator iter.Seq[*modelruntimeentities.LLMResultChunk],
) iter.Seq[*modelruntimeentities.LLMResultChunk] {
	/*
		Code block mode stream processor, ensure the response is a code block with output markdown quote

		:param model: model name
		:param prompt_messages: prompt messages
		:param input_generator: input generator
		:return: output generator
	*/
	return func(yield func(*modelruntimeentities.LLMResultChunk) bool) {
		state := "normal"
		backtick_count := 0
		content_piece := ""
		for piece := range input_generator {
			if piece.Delta.Message.Content != "" {
				content := piece.Delta.Message.Content
				piece.Delta.Message.Content = ""
				if !yield(piece) {
					return
				}
				content_piece = content
			} else {
				if !yield(piece) {
					return
				}
				continue
			}
			new_piece := ""
			for _, c := range []byte(content_piece) {
				char := string(c)
				if state == "normal" {
					if char == "`" {
						state = "in_backticks"
						backtick_count = 1
					} else {
						new_piece += char
					}
				} else if state == "in_backticks" {
					if char == "`" {
						backtick_count += 1
						if backtick_count == 3 {
							state = "skip_content"
							backtick_count = 0
						}
					} else {
						new_piece += strings.Repeat("`", backtick_count) + char
						state = "normal"
						backtick_count = 0
					}
				} else if state == "skip_content" {
					if char == " " {
						state = "normal"
					}
				}
				if new_piece != "" {
					yield(&modelruntimeentities.LLMResultChunk{
						Model:          model,
						PromptMessages: prompt_messages,
						Delta: &modelruntimeentities.LLMResultChunkDelta{
							Index:   0,
							Message: modelruntimeentities.NewAssistantPromptMessage(new_piece, "", nil),
						},
					})
				}
			}
		}
	}
}
func (m *LargeLanguageModel) CodeBlockModeStreamProcessorWithBacktick(
	model string, prompt_messages []modelruntimeentities.PromptMessager, input_generator iter.Seq[*modelruntimeentities.LLMResultChunk],
) iter.Seq[*modelruntimeentities.LLMResultChunk] {
	/*
		Code block mode stream processor, ensure the response is a code block with output markdown quote.
		This version skips the language identifier that follows the opening triple backticks.

		:param model: model name
		:param prompt_messages: prompt messages
		:param input_generator: input generator
		:return: output generator
	*/
	return func(yield func(*modelruntimeentities.LLMResultChunk) bool) {
		state := "search_start"
		backtick_count := 0
		content_piece := ""

		for piece := range input_generator {
			if piece.Delta.Message.Content != "" {
				content := piece.Delta.Message.Content
				// Reset content to ensure we're only processing and yielding the relevant parts
				piece.Delta.Message.Content = ""
				// Yield a piece with cleared content before processing it to maintain the generator structure
				yield(piece)
				content_piece = content
			} else {
				// Yield pieces without content directly
				yield(piece)
				continue
			}
			if state == "done" {
				continue
			}
			new_piece := ""
			for _, char := range []byte(content_piece) {
				if state == "search_start" {
					if string(char) == "`" {
						backtick_count += 1
						if backtick_count == 3 {
							state = "skip_language"
							backtick_count = 0
						}
					} else {
						backtick_count = 0
					}
				} else if state == "skip_language" {
					// Skip everything until the first newline, marking the end of the language identifier
					if string(char) == "\n" {
						state = "in_code_block"
					}
				} else if state == "in_code_block" {
					if string(char) == "`" {
						backtick_count += 1
						if backtick_count == 3 {
							state = "done"
							break
						}
					} else {
						if backtick_count > 0 {
							// If backticks were counted but we're still collecting content, it was a false start
							new_piece += strings.Repeat("`", backtick_count)
							backtick_count = 0
						}
						new_piece += string(char)
					}
				} else if state == "done" {
					break
				}
				if new_piece != "" {
					// Only yield content collected within the code block
					yield(&modelruntimeentities.LLMResultChunk{
						Model:          model,
						PromptMessages: prompt_messages,
						Delta: &modelruntimeentities.LLMResultChunkDelta{
							Index:   0,
							Message: modelruntimeentities.NewAssistantPromptMessage(new_piece, "", nil),
						},
					})
				}
			}
		}
	}
}
func (m *LargeLanguageModel) WrapThinkingByReasoningContent(delta map[string]any, is_reasoning bool) (string, bool) {
	/*
		If the reasoning response is from delta.get("reasoning_content"), we wrap
		it with HTML details tag.

		:param delta: delta dictionary from LLM streaming response
		:param is_reasoning: is reasoning
		:return: tuple of (processed_content, is_reasoning)
	*/
	content := ""
	if _, ok := delta["content"]; ok {
		if _, ok := delta["content"].(string); ok {
			content = delta["content"].(string)
		}
	}
	reasoning_content := ""
	if _, ok := delta["reasoning_content"]; ok {
		if _, ok := delta["reasoning_content"].(string); ok {
			reasoning_content = delta["reasoning_content"].(string)
		}
	}

	if reasoning_content != "" {
		if !is_reasoning {
			content = HTML_THINKING_TAG + reasoning_content
			is_reasoning = true
		} else {
			content = reasoning_content
		}
	} else if is_reasoning {
		content = "</details>" + content
		is_reasoning = false
	}
	return content, is_reasoning
}
func (m *LargeLanguageModel) WrapThinkingByTag(content string) string {
	/*
		if the reasoning response is a <think>...</think> block from delta.get("content"),
		we replace <think> to <detail>.

			:param content: delta.get("content")
			:return: processed_content
	*/
	return strings.ReplaceAll(strings.ReplaceAll(content, "<think>", HTML_THINKING_TAG), "</think>", "</details>")
}
func (m *LargeLanguageModel) InvokeResultGenerator(
	modeler modelruntimeentities.LargeLanguageModeler,
	model string,
	result iter.Seq[*modelruntimeentities.LLMResultChunk],
	credentials map[string]any,
	prompt_messages []modelruntimeentities.PromptMessager,
	model_parameters map[string]any,
	tools []*modelruntimeentities.PromptMessageTool,
	stop []string,
	stream bool, /* = True*/
	user string,
	callbacks []modelruntimeentities.Callbacker,
) iter.Seq2[*modelruntimeentities.LLMResultChunk, error] {
	/*
		Invoke result generator

		:param result: result generator
		:return: result generator
	*/
	return func(yield func(*modelruntimeentities.LLMResultChunk, error) bool) {
		prompt_message := modelruntimeentities.NewAssistantPromptMessage("", "", nil)
		var usage *modelruntimeentities.LLMUsage
		var system_fingerprint string
		real_model := model

		// try{
		for chunk := range result {
			yield(chunk, nil)

			m.TriggerNewChunkCallbacks(
				modeler,
				chunk,
				model,
				credentials,
				prompt_messages,
				model_parameters,
				tools,
				stop,
				stream,
				user,
				callbacks,
			)

			prompt_message.Content += chunk.Delta.Message.Content
			real_model = chunk.Model
			if chunk.Delta.Usage != nil {
				usage = chunk.Delta.Usage
			}
			if chunk.SystemFingerprint != "" {
				system_fingerprint = chunk.SystemFingerprint
			}
		}
		// except Exception as e{
		// 	raise m.TransformInvokeError(e)
		if usage == nil {
			usage = modelruntimeentities.NewLLMUsage()
		}
		m.TriggerAfterInvokeCallbacks(
			modeler,
			model,
			&modelruntimeentities.LLMResult{
				Model:             real_model,
				PromptMessages:    prompt_messages,
				Message:           prompt_message,
				Usage:             usage,
				SystemFingerprint: system_fingerprint,
			},
			credentials,
			prompt_messages,
			model_parameters,
			tools,
			stop,
			stream,
			user,
			callbacks,
		)
	}
}
func (m *LargeLanguageModel) EnforceStopTokens(text string, stop []string) string {
	/*Cut off the text as soon as any stop words occur.*/
	re, err := regexp.Compile(strings.Join(stop, "|"))
	if err != nil {
		mlog.Errorf("regxp compile(%s) failed:%v", strings.Join(stop, "|"), err)
		return text
	}

	return re.Split(text, 1)[0]
}

func (m *LargeLanguageModel) CalcResponseUsage(
	modeler modelruntimeentities.LargeLanguageModeler, model string, credentials map[string]any, prompt_tokens int, completion_tokens int,
) *modelruntimeentities.LLMUsage {
	/*
		Calculate response usage

		:param model: model name
		:param credentials: model credentials
		:param prompt_tokens: prompt tokens
		:param completion_tokens: completion tokens
		:return: usage
	*/
	// get prompt price info
	prompt_price_info, _ := m.GetPrice(
		modeler,
		model,
		credentials,
		modelruntimeentities.PriceType_INPUT,
		prompt_tokens,
	)

	// get completion price info
	completion_price_info, _ := m.GetPrice(
		modeler, model, credentials, modelruntimeentities.PriceType_OUTPUT, completion_tokens,
	)

	// transform usage
	usage := &modelruntimeentities.LLMUsage{
		PromptTokens:        prompt_tokens,
		PromptUnitPrice:     prompt_price_info.UnitPrice,
		PromptPriceUnit:     prompt_price_info.Unit,
		PromptPrice:         prompt_price_info.TotalAmount,
		CompletionTokens:    completion_tokens,
		CompletionUnitPrice: completion_price_info.UnitPrice,
		CompletionPriceUnit: completion_price_info.Unit,
		CompletionPrice:     completion_price_info.TotalAmount,
		TotalTokens:         prompt_tokens + completion_tokens,
		TotalPrice:          prompt_price_info.TotalAmount + completion_price_info.TotalAmount,
		Currency:            prompt_price_info.Currency,
		Latency:             time.Since(m.StartedAt).Seconds(),
	}

	return usage
}
func (m *LargeLanguageModel) TriggerBeforeInvokeCallbacks(
	modeler modelruntimeentities.LargeLanguageModeler,
	model string,
	credentials map[string]any,
	prompt_messages []modelruntimeentities.PromptMessager,
	model_parameters map[string]any,
	tools []*modelruntimeentities.PromptMessageTool,
	stop []string,
	stream bool, /* = True*/
	user string,
	callbacks []modelruntimeentities.Callbacker,
) {
	/*
		Trigger before invoke callbacks

			:param model: model name
			:param credentials: model credentials
			:param prompt_messages: prompt messages
			:param model_parameters: model parameters
			:param tools: tools for tool calling
			:param stop: stop words
			:param stream: is stream response
			:param user: unique user id
			:param callbacks: callbacks
	*/
	if len(callbacks) > 0 {
		for _, callback := range callbacks {
			// try{
			callback.OnBeforeInvoke(
				modeler,
				model,
				credentials,
				prompt_messages,
				model_parameters,
				tools,
				stop,
				stream,
				user,
			)
		}
	}
	// except Exception as e{
	// 	if callback.raise_error{
	// 		raise e
	// 	else{
	// 		logger.warning(f"Callback {callback.__class__.__name__} on_before_invoke failed with error {e}")
}
func (m *LargeLanguageModel) TriggerNewChunkCallbacks(
	modeler modelruntimeentities.LargeLanguageModeler,
	chunk *modelruntimeentities.LLMResultChunk,
	model string,
	credentials map[string]any,
	prompt_messages []modelruntimeentities.PromptMessager,
	model_parameters map[string]any,
	tools []*modelruntimeentities.PromptMessageTool,
	stop []string,
	stream bool, /* = True*/
	user string,
	callbacks []modelruntimeentities.Callbacker,
) {
	/*
		Trigger new chunk callbacks

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
	if len(callbacks) > 0 {
		for _, callback := range callbacks {
			// try{
			callback.OnNewChunk(
				modeler,
				chunk,
				model,
				credentials,
				prompt_messages,
				model_parameters,
				tools,
				stop,
				stream,
				user,
			)
		}
	}
	// except Exception as e{
	// 	if callback.raise_error{
	// 		raise e
	// 	else{
	// 		logger.warning(f"Callback {callback.__class__.__name__} on_new_chunk failed with error {e}")
}
func (m *LargeLanguageModel) TriggerAfterInvokeCallbacks(
	modeler modelruntimeentities.LargeLanguageModeler,
	model string,
	result *modelruntimeentities.LLMResult,
	credentials map[string]any,
	prompt_messages []modelruntimeentities.PromptMessager,
	model_parameters map[string]any,
	tools []*modelruntimeentities.PromptMessageTool,
	stop []string,
	stream bool, /* = True*/
	user string,
	callbacks []modelruntimeentities.Callbacker,
) {
	/*
		Trigger after invoke callbacks

			:param model: model name
			:param result: result
			:param credentials: model credentials
			:param prompt_messages: prompt messages
			:param model_parameters: model parameters
			:param tools: tools for tool calling
			:param stop: stop words
			:param stream: is stream response
			:param user: unique user id
			:param callbacks: callbacks
	*/
	if len(callbacks) > 0 {
		for _, callback := range callbacks {
			// try{
			callback.OnAfterInvoke(
				modeler,
				result,
				model,
				credentials,
				prompt_messages,
				model_parameters,
				tools,
				stop,
				stream,
				user,
			)
		}
	}
	// except Exception as e{
	// 	if callback.raise_error{
	// 		raise e
	// 	else{
	// 		logger.warning(f"Callback {callback.__class__.__name__} on_after_invoke failed with error {e}")
}
func (m *LargeLanguageModel) Trigger_invoke_error_callbacks(
	modeler modelruntimeentities.LargeLanguageModeler,
	model string,
	err error,
	credentials map[string]any,
	prompt_messages []modelruntimeentities.PromptMessager,
	model_parameters map[string]any,
	tools []*modelruntimeentities.PromptMessageTool,
	stop []string,
	stream bool, /* = True*/
	user string,
	callbacks []modelruntimeentities.Callbacker,
) {
	/*
		Trigger invoke error callbacks

		:param model: model name
		:param ex: exception
		:param credentials: model credentials
		:param prompt_messages: prompt messages
		:param model_parameters: model parameters
		:param tools: tools for tool calling
		:param stop: stop words
		:param stream: is stream response
		:param user: unique user id
		:param callbacks: callbacks
	*/
	if len(callbacks) > 0 {
		for _, callback := range callbacks {
			// try{
			callback.OnInvokeError(
				modeler,
				err,
				model,
				credentials,
				prompt_messages,
				model_parameters,
				tools,
				stop,
				stream,
				user,
			)
		}
	}
	// except Exception as e{
	// 	if callback.raise_error{
	// 		raise e
	// 	else{
	// 		logger.warning(f"Callback {callback.__class__.__name__} on_invoke_error failed with error {e}")
}

func (m *LargeLanguageModel) GetParameterRules(modeler modelruntimeentities.LargeLanguageModeler, model string, credentials map[string]any) []*modelruntimeentities.ParameterRule {
	/*
		Get parameter rules

		:param model: model name
		:param credentials: model credentials
		:return: parameter rules
	*/
	model_schema := m.GetModelSchema(modeler, model, credentials)
	if model_schema != nil {
		for idx, v := range model_schema.ParameterRules {
			if v.Label.ZhHans == "" {
				v.Label.ZhHans = v.Name
			}
			if v.Label.EnUS == "" {
				v.Label.EnUS = v.Name
			}
			model_schema.ParameterRules[idx] = v
		}
		return model_schema.ParameterRules
	}
	return nil
}
func (m *LargeLanguageModel) GetModelMode(modeler modelruntimeentities.LargeLanguageModeler, model string, credentials map[string]any) modelruntimeentities.LLMMode {
	/*
		Get model mode

		:param model: model name
		:param credentials: model credentials
		:return: model mode
	*/
	model_schema := m.GetModelSchema(modeler, model, credentials)

	mode := modelruntimeentities.LLMMode_CHAT
	if model_schema != nil {
		if _, ok := model_schema.ModelProperties[modelruntimeenumtypes.ModelPropertyKey_MODE]; ok {
			mode = modelruntimeentities.LLMMode(model_schema.ModelProperties[modelruntimeenumtypes.ModelPropertyKey_MODE].(string))
		}
	}
	return mode
}

func (m *LargeLanguageModel) ValidateAndFilterModelParameters(modeler modelruntimeentities.LargeLanguageModeler, model string, model_parameters map[string]any, credentials map[string]any) (map[string]any, error) {
	/*
		Validate model parameters

		:param model: model name
		:param model_parameters: model parameters
		:param credentials: model credentials
		:return
	*/
	parameter_rules := m.GetParameterRules(modeler, model, credentials)

	// validate model parameters
	filtered_model_parameters := map[string]any{}
	for _, parameter_rule := range parameter_rules {
		parameter_name := parameter_rule.Name
		var parameter_value any
		var ok bool
		if parameter_value, ok = model_parameters[parameter_name]; !ok || parameter_value == nil {
			if _, ok1 := model_parameters[parameter_rule.UseTemplate]; ok1 && parameter_rule.UseTemplate != "" {
				// if parameter value is None, use template value variable name instead
				parameter_value = model_parameters[parameter_rule.UseTemplate]
			} else {
				if parameter_rule.Required {
					if parameter_rule.Default != nil {
						filtered_model_parameters[parameter_name] = parameter_rule.Default
						continue
					} else {
						return nil, exceptions.NewValueError(fmt.Sprintf("Model Parameter %s is required.", parameter_name))
					}
				} else {
					continue
				}
			}
		}
		// validate parameter value type
		if parameter_rule.Type == modelruntimeenumtypes.ParameterType_INT {
			if _, ok := parameter_value.(int); !ok {
				return nil, exceptions.NewValueError(fmt.Sprintf("Model Parameter %v should be int.", parameter_name))
			}
			// validate parameter value range
			if float64(parameter_value.(int)) < parameter_rule.Min {
				return nil, exceptions.NewValueError(fmt.Sprintf("Model Parameter %v should be greater than or equal to %v.", parameter_name, parameter_rule.Min))
			}
			if float64(parameter_value.(int)) > parameter_rule.Max {
				return nil, exceptions.NewValueError(fmt.Sprintf("Model Parameter %v should be less than or equal to %v.", parameter_name, parameter_rule.Max))
			}
		} else if parameter_rule.Type == modelruntimeenumtypes.ParameterType_FLOAT {
			var real_parameter_value float64
			if _, ok := parameter_value.(int); !ok {
				if _, ok := parameter_value.(float64); !ok {
					return nil, exceptions.NewValueError(fmt.Sprintf("Model Parameter %v should be float.", parameter_name))
				} else {
					real_parameter_value = parameter_value.(float64)
				}
			} else {
				real_parameter_value = float64(parameter_value.(int))
			}

			// validate parameter value precision
			if parameter_rule.Precision == 0 {
				if parameter_value != int(real_parameter_value) {
					return nil, exceptions.NewValueError(fmt.Sprintf("Model Parameter %v should be int.", parameter_name))
				}
			} else {
				// math.Round
				// if parameter_value != round(parameter_value, parameter_rule.Precision){
				// 	raise ValueError(
				// 		f"Model Parameter {parameter_name} should be round to {parameter_rule.Precision}"
				// 		f" decimal places."
				// 	)
				// }
			}

			// validate parameter value range
			if real_parameter_value < parameter_rule.Min {
				return nil, exceptions.NewValueError(fmt.Sprintf("Model Parameter %v should be greater than or equal to %v.", parameter_name, parameter_rule.Min))
			}
			if real_parameter_value > parameter_rule.Max {
				return nil, exceptions.NewValueError(fmt.Sprintf("Model Parameter %v should be less than or equal to %v.", parameter_name, parameter_rule.Max))
			}
		} else if parameter_rule.Type == modelruntimeenumtypes.ParameterType_BOOLEAN {
			if _, ok := parameter_value.(bool); !ok {
				exceptions.NewValueError(fmt.Sprintf("Model Parameter %v should be bool.", parameter_name))
			}
		} else if parameter_rule.Type == modelruntimeenumtypes.ParameterType_STRING {
			if _, ok := parameter_value.(string); !ok {
				exceptions.NewValueError(fmt.Sprintf("Model Parameter %v should be string.", parameter_name))
			}

			// validate options
			if len(parameter_rule.Options) > 0 && !slices.Contains(parameter_rule.Options, parameter_value.(string)) {
				exceptions.NewValueError(fmt.Sprintf("Model Parameter %v should be one of %v.", parameter_name, parameter_rule.Options))
			}
		} else if parameter_rule.Type == modelruntimeenumtypes.ParameterType_TEXT {
			if _, ok := parameter_value.(string); !ok {
				exceptions.NewValueError(fmt.Sprintf("Model Parameter %v should be text.", parameter_name))
			}

			// validate options
			if len(parameter_rule.Options) > 0 && !slices.Contains(parameter_rule.Options, parameter_value.(string)) {
				exceptions.NewValueError(fmt.Sprintf("Model Parameter %v should be one of %v.", parameter_name, parameter_rule.Options))
			}
		} else {
			exceptions.NewValueError(fmt.Sprintf("Model Parameter %v type %v is not supported.", parameter_name, parameter_rule.Type))
		}
		filtered_model_parameters[parameter_name] = parameter_value
	}
	return filtered_model_parameters, nil
}
