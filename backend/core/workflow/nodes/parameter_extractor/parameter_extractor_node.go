package parameterextractor

import (
	"encoding/json"
	"fmt"
	"iter"
	"maps"
	"math"
	"slices"
	"strconv"
	"strings"

	"github.com/odysseythink/mlog"
	uuid "github.com/satori/go.uuid"
	penodesexceptions "mlib.com/gofy/server/core/exceptions/nodes/parameter_extractor"
	"mlib.com/gofy/server/core/file"
	modelmanager "mlib.com/gofy/server/core/manageres/model_manager"
	"mlib.com/gofy/server/core/memory"
	"mlib.com/gofy/server/core/prompt"
	promptutils "mlib.com/gofy/server/core/prompt/utils"
	"mlib.com/gofy/server/core/workflow/nodes/base"
	"mlib.com/gofy/server/core/workflow/nodes/llm"
	variabletemplateparser "mlib.com/gofy/server/core/workflow/utils/variable_template_parser"
	appconfigentities "mlib.com/gofy/server/entities/app/config"
	modelruntimeentities "mlib.com/gofy/server/entities/model_runtime"
	llmnodesentities "mlib.com/gofy/server/entities/nodes/llm"
	penodesentities "mlib.com/gofy/server/entities/nodes/parameter_extractor"
	promptentities "mlib.com/gofy/server/entities/prompt"
	workflowentities "mlib.com/gofy/server/entities/workflow"
	modelruntimeenumtypes "mlib.com/gofy/server/enum_types/model_runtime"
	nodesenumtypes "mlib.com/gofy/server/enum_types/nodes"
	workflowenumtypes "mlib.com/gofy/server/enum_types/workflow"
	"mlib.com/gofy/server/models"
)

type ParameterExtractorNode struct {
	*base.BaseNode[*penodesentities.ParameterExtractorNodeData]
	_model_instance *modelmanager.ModelInstance
	_model_config   *appconfigentities.ModelConfigWithCredentialsEntity
}

func (n *ParameterExtractorNode) GetDefaultConfig(filters map[string]any) map[string]any {
	return map[string]any{
		"model": map[string]map[string]map[string]any{
			"prompt_templates": {
				"completion_model": {
					"conversation_histories_role": map[string]any{"user_prefix": "Human", "assistant_prefix": "Assistant"},
					"stop":                        []string{"Human:"},
				},
			},
		},
	}
}
func (n *ParameterExtractorNode) Type() nodesenumtypes.NodeType {
	return nodesenumtypes.Node_PARAMETER_EXTRACTOR
}

func (n *ParameterExtractorNode) Run() (*workflowentities.NodeRunResult, iter.Seq[any]) {
	node_data := n.NodeData
	variable := n.GetGraphRuntimeState().VariablePool.Get(node_data.Query)
	query := ""
	if variable != nil {
		query = variable.Text()
	}

	model_instance, model_config := n._fetch_model_config(&node_data.Model)
	if _, ok := any(model_instance.ModelTypeInstance).(modelruntimeentities.LargeLanguageModeler); !ok {
		panic(penodesexceptions.NewInvalidModelTypeError("Model is not a Large Language Model"))
	}
	llm_model := model_instance.ModelTypeInstance
	model_schema := llm_model.GetModelSchema(llm_model, model_config.Model, model_config.Credentials)
	if model_schema == nil {
		panic(penodesexceptions.NewModelSchemaNotFoundError("Model schema not found"))
	}

	// fetch memory
	mem := llm.FetchMemory(n.GetGraphRuntimeState(), n.GetAppID(), node_data.Memory, model_instance)
	var prompt_messages []modelruntimeentities.PromptMessager
	var prompt_message_tools []*modelruntimeentities.PromptMessageTool
	if (slices.Contains(model_schema.Features, modelruntimeenumtypes.ModelFeature_MULTI_TOOL_CALL) ||
		slices.Contains(model_schema.Features, modelruntimeenumtypes.ModelFeature_TOOL_CALL)) &&
		node_data.ReasoningMode == "function_call" {
		// use function call
		prompt_messages, prompt_message_tools = n._generate_function_call_prompt(
			node_data,
			query,
			n.GetGraphRuntimeState().VariablePool,
			model_config,
			mem,
			nil,
		)
	} else {
		// use prompt engineering
		prompt_messages = n._generate_prompt_engineering_prompt(
			node_data,
			query,
			n.GetGraphRuntimeState().VariablePool,
			model_config,
			mem,
			nil,
		)

		prompt_message_tools = []*modelruntimeentities.PromptMessageTool{}
	}
	inputs := map[string]any{
		"query":       query,
		"files":       nil,
		"parameters":  node_data.Parameters,
		"instruction": node_data.Instruction,
	}

	process_data := map[string]any{
		"model_mode": model_config.Mode,
		"prompts":    promptutils.PromptMessagesToPromptForSaving(model_config.Mode, prompt_messages),
		"usage":      nil,
		"tool_call":  nil,
	}
	if len(prompt_message_tools) > 0 {
		process_data["function"] = prompt_message_tools[0]
	}
	var ret *workflowentities.NodeRunResult
	var text string
	var usage *modelruntimeentities.LLMUsage
	var tool_call *modelruntimeentities.ToolCall
	func() {
		defer func() {
			if r := recover(); r != nil {
				if exp, ok := r.(*penodesexceptions.ParameterExtractorNodeError); ok {
					ret = &workflowentities.NodeRunResult{
						Status:      models.WorkflowNodeExecutionStatus_FAILED,
						Inputs:      inputs,
						ProcessData: process_data,
						Outputs:     map[string]any{"__is_success": 0, "__reason": exp.Error()},
						Error:       exp.Error(),
						Metadata:    nil,
					}
				} else if exp, ok := r.(error); ok {
					ret = &workflowentities.NodeRunResult{
						Status:      models.WorkflowNodeExecutionStatus_FAILED,
						Inputs:      inputs,
						ProcessData: process_data,
						Outputs:     map[string]any{"__is_success": 0, "__reason": "Failed to invoke model", "__error": exp.Error()},
						Error:       exp.Error(),
						Metadata:    nil,
					}
				} else {
					panic(r)
				}
			}
		}()

		text, usage, tool_call = n._invoke(
			&node_data.Model,
			model_instance,
			prompt_messages,
			prompt_message_tools,
			model_config.Stop,
		)
		process_data["usage"] = usage
		process_data["tool_call"] = tool_call
		process_data["llm_text"] = text
	}()
	if ret != nil {
		return ret, nil
	}

	errormsg := ""
	var result map[string]any
	if tool_call != nil {
		result = n._extract_json_from_tool_call(tool_call)
	} else {
		result = n._extract_complete_json_response(text)
		if result == nil {
			result = n._generate_default_result(node_data)
			errormsg = "Failed to extract result from function call or text response, using empty result."
		}
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				if exp, ok := r.(*penodesexceptions.ParameterExtractorNodeError); ok {
					errormsg = exp.Error()
				} else {
					panic(r)
				}
			}
		}()

		result = n._validate_result(node_data, result)
	}()

	// transform result into standard format
	result = n._transform_result(node_data, result)
	outputs := maps.Clone(result)
	if errormsg != "" {
		outputs["__is_success"] = 0
		outputs["__reason"] = errormsg
	} else {
		outputs["__is_success"] = 1
		outputs["__reason"] = errormsg
	}
	return &workflowentities.NodeRunResult{
		Status:      models.WorkflowNodeExecutionStatus_SUCCEEDED,
		Inputs:      inputs,
		ProcessData: process_data,
		Outputs:     outputs,
		Metadata: map[workflowenumtypes.NodeRunMetadataKey]any{
			workflowenumtypes.NodeRunMetadataKey_TOTAL_TOKENS: usage.TotalTokens,
			workflowenumtypes.NodeRunMetadataKey_TOTAL_PRICE:  usage.TotalPrice,
			workflowenumtypes.NodeRunMetadataKey_CURRENCY:     usage.Currency,
		},
		LLmUsage: usage,
	}, nil
}

func (n *ParameterExtractorNode) ExtractVariableSelectorToVariableMapping(graph_config map[string]any, node_id string, node_data *penodesentities.ParameterExtractorNodeData) map[string][]string {
	// FIXME: fix the type error later
	variable_mapping := map[string][]string{"query": node_data.Query}

	if node_data.Instruction != "" {
		selectors := variabletemplateparser.ExtractSelectorsFromTemplate(node_data.Instruction)
		for _, selector := range selectors {
			variable_mapping[selector.Variable] = selector.ValueSelector
		}
	}
	for key, value := range variable_mapping {
		variable_mapping[node_id+"."+key] = value
	}

	return variable_mapping
}

func (n *ParameterExtractorNode) _invoke(
	node_data_model *llmnodesentities.ModelConfig,
	model_instance *modelmanager.ModelInstance,
	prompt_messages []modelruntimeentities.PromptMessager,
	tools []*modelruntimeentities.PromptMessageTool,
	stop []string,
) (string, *modelruntimeentities.LLMUsage, *modelruntimeentities.ToolCall) {
	invoke_result := model_instance.InvokeLLM(
		prompt_messages,
		node_data_model.CompletionParams,
		tools,
		stop,
		n.GetUserID(),
		nil,
	)

	// handle invoke result
	text := invoke_result.Message.Content
	usage := invoke_result.Usage
	var tool_call *modelruntimeentities.ToolCall
	if len(invoke_result.Message.ToolCalls) > 0 {
		tool_call = invoke_result.Message.ToolCalls[0]
	}

	// deduct quota
	llm.DeductLLMQuota(n.GetTenantID(), model_instance, usage)

	return text, usage, tool_call
}

func (n *ParameterExtractorNode) _get_function_calling_prompt_template(
	node_data *penodesentities.ParameterExtractorNodeData,
	query string,
	variable_pool *workflowentities.VariablePool,
	mem *memory.TokenBufferMemory,
	max_token_limit int,
) []*promptentities.ChatModelMessage {
	if max_token_limit <= 0 {
		max_token_limit = 2000
	}
	model_mode := node_data.Model.Mode
	input_text := query
	memory_str := ""
	instruction := variable_pool.ConvertTemplate(node_data.Instruction).Text()

	if mem != nil && node_data.Memory != nil && node_data.Memory.Window.Enabled {
		memory_str = mem.GetHistoryPromptText("", "", max_token_limit, node_data.Memory.Window.Size)
	}
	if model_mode == modelruntimeentities.LLMMode_CHAT {
		system_prompt_messages := &promptentities.ChatModelMessage{
			Role: modelruntimeentities.PromptMessageRole_SYSTEM,
			Text: strings.ReplaceAll(strings.ReplaceAll(FUNCTION_CALLING_EXTRACTOR_SYSTEM_PROMPT, "{histories}", memory_str), "{instruction}", instruction),
		}
		user_prompt_message := &promptentities.ChatModelMessage{Role: modelruntimeentities.PromptMessageRole_USER, Text: input_text}
		return []*promptentities.ChatModelMessage{system_prompt_messages, user_prompt_message}
	} else {
		panic(penodesexceptions.NewInvalidModelModeError(fmt.Sprintf("Model mode {%s} not support.", model_mode)))
	}
}
func (n *ParameterExtractorNode) _get_prompt_engineering_prompt_template(
	node_data *penodesentities.ParameterExtractorNodeData,
	query string,
	variable_pool *workflowentities.VariablePool,
	mem *memory.TokenBufferMemory,
	max_token_limit int,
) any {
	if max_token_limit <= 0 {
		max_token_limit = 2000
	}
	model_mode := node_data.Model.Mode
	input_text := query
	memory_str := ""
	instruction := variable_pool.ConvertTemplate(node_data.Instruction).Text()

	if mem != nil && node_data.Memory != nil && node_data.Memory.Window.Enabled {
		memory_str = mem.GetHistoryPromptText("", "", max_token_limit, node_data.Memory.Window.Size)
	}
	if model_mode == modelruntimeentities.LLMMode_CHAT {
		system_prompt_messages := &promptentities.ChatModelMessage{
			Role: modelruntimeentities.PromptMessageRole_SYSTEM,
			Text: strings.ReplaceAll(strings.ReplaceAll(FUNCTION_CALLING_EXTRACTOR_SYSTEM_PROMPT, "{histories}", memory_str), "{instruction}", instruction),
		}
		user_prompt_message := &promptentities.ChatModelMessage{Role: modelruntimeentities.PromptMessageRole_USER, Text: input_text}
		return []*promptentities.ChatModelMessage{system_prompt_messages, user_prompt_message}
	} else if model_mode == modelruntimeentities.LLMMode_COMPLETION {
		return &promptentities.CompletionModelPromptTemplate{
			Text: strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(COMPLETION_GENERATE_JSON_PROMPT, "{histories}", memory_str), "{text}", input_text), "{instruction}", instruction), "{γγγ", ""), "}γγγ", ""),
		}
	} else {
		panic(penodesexceptions.NewInvalidModelModeError(fmt.Sprintf("Model mode {%s} not support.", model_mode)))
	}
}
func (n *ParameterExtractorNode) _fetch_model_config(
	node_data_model *llmnodesentities.ModelConfig,
) (*modelmanager.ModelInstance, *appconfigentities.ModelConfigWithCredentialsEntity) {
	/*
	   Fetch model config.
	*/
	if n._model_instance == nil || n._model_config == nil {
		n._model_instance, n._model_config = llm.FetchModelConfig(n.GetTenantID(), node_data_model)
	}
	return n._model_instance, n._model_config
}
func (n *ParameterExtractorNode) _calculate_rest_token(
	node_data *penodesentities.ParameterExtractorNodeData,
	query string,
	variable_pool *workflowentities.VariablePool,
	model_config *appconfigentities.ModelConfigWithCredentialsEntity,
	context string,
) int {

	model_instance, model_config := n._fetch_model_config(&node_data.Model)
	if _, ok := any(model_instance.ModelTypeInstance).(modelruntimeentities.LargeLanguageModeler); !ok {
		panic(penodesexceptions.NewInvalidModelTypeError("Model is not a Large Language Model"))
	}
	llm_model := model_instance.ModelTypeInstance
	model_schema := llm_model.GetModelSchema(llm_model, model_config.Model, model_config.Credentials)
	if model_schema == nil {
		panic(penodesexceptions.NewModelSchemaNotFoundError("Model schema not found"))
	}
	var prompt_template any
	if slices.Contains(model_schema.Features, modelruntimeenumtypes.ModelFeature_MULTI_TOOL_CALL) {
		prompt_template = n._get_function_calling_prompt_template(node_data, query, variable_pool, nil, 2000)
	} else {
		prompt_template = n._get_prompt_engineering_prompt_template(node_data, query, variable_pool, nil, 2000)
	}
	var prompt_messages []modelruntimeentities.PromptMessager
	if real_prompt_template, ok := prompt_template.([]*promptentities.ChatModelMessage); ok {
		prompt_transform := prompt.NewAdvancedPromptTransform[[]*promptentities.ChatModelMessage](true, modelruntimeentities.ImagePromptMessageContentDETAIL(""))
		prompt_messages = prompt_transform.GetPrompt(
			real_prompt_template,
			map[string]string{},
			"",
			nil,
			context,
			node_data.Memory,
			nil,
			model_config,
		)
	} else if real_prompt_template, ok := prompt_template.(*promptentities.CompletionModelPromptTemplate); ok {
		prompt_transform := prompt.NewAdvancedPromptTransform[*promptentities.CompletionModelPromptTemplate](true, modelruntimeentities.ImagePromptMessageContentDETAIL(""))
		prompt_messages = prompt_transform.GetPrompt(
			real_prompt_template,
			map[string]string{},
			"",
			nil,
			context,
			node_data.Memory,
			nil,
			model_config,
		)
	} else {
		panic(penodesexceptions.NewParameterExtractorNodeError("prompt_template must be []*ChatModelMessage or *CompletionModelPromptTemplat"))
	}

	rest_tokens := 2000

	if _, ok := model_config.ModelSchema.ModelProperties[modelruntimeenumtypes.ModelPropertyKey_CONTEXT_SIZE]; ok {

		if model_context_tokens, ok := model_config.ModelSchema.ModelProperties[modelruntimeenumtypes.ModelPropertyKey_CONTEXT_SIZE].(int); ok {
			model_type_instance := model_config.ProviderModelBundle.ModelTypeInstance.(modelruntimeentities.LargeLanguageModeler)

			curr_message_tokens := model_type_instance.GetNumTokens(model_config.Model, model_config.Credentials, prompt_messages, nil) + 1000 // add 1000 to ensure tool call messages

			max_tokens := 0
			for _, parameter_rule := range model_config.ModelSchema.ParameterRules {
				if parameter_rule.Name == "max_tokens" ||
					(parameter_rule.UseTemplate != "" && parameter_rule.UseTemplate == "max_tokens") {
					if _, ok := model_config.Parameters[parameter_rule.Name]; ok {
						if _, ok := model_config.Parameters[parameter_rule.Name].(int); ok {
							max_tokens = model_config.Parameters[parameter_rule.Name].(int)
						}
					}
					if max_tokens == 0 {
						if _, ok := model_config.Parameters[parameter_rule.UseTemplate]; ok {
							if _, ok := model_config.Parameters[parameter_rule.UseTemplate].(int); ok {
								max_tokens = model_config.Parameters[parameter_rule.UseTemplate].(int)
							}
						}
					}
				}
			}

			rest_tokens = model_context_tokens - max_tokens - curr_message_tokens
			rest_tokens = int(math.Max(float64(rest_tokens), float64(0)))
		}
	}
	return rest_tokens
}
func (n *ParameterExtractorNode) _generate_function_call_prompt(
	node_data *penodesentities.ParameterExtractorNodeData,
	query string,
	variable_pool *workflowentities.VariablePool,
	model_config *appconfigentities.ModelConfigWithCredentialsEntity,
	mem *memory.TokenBufferMemory,
	files []*file.File,
) ([]modelruntimeentities.PromptMessager, []*modelruntimeentities.PromptMessageTool) {
	/*
	   Generate function call prompt.
	*/
	bindata, _ := json.Marshal(node_data.GetParameterJsonSchema())
	query = strings.ReplaceAll(strings.ReplaceAll(FUNCTION_CALLING_EXTRACTOR_USER_TEMPLATE, "{content}", query), "{structure}", string(bindata))

	prompt_transform := prompt.NewAdvancedPromptTransform[[]*promptentities.ChatModelMessage](true, modelruntimeentities.ImagePromptMessageContentDETAIL(""))
	rest_token := n._calculate_rest_token(node_data, query, variable_pool, model_config, "")
	prompt_template := n._get_function_calling_prompt_template(
		node_data, query, variable_pool, mem, rest_token,
	)
	prompt_messages := prompt_transform.GetPrompt(
		prompt_template,
		nil,
		"",
		files,
		"",
		node_data.Memory,
		nil,
		model_config,
	)

	// find last user message
	last_user_message_idx := len(prompt_messages) - 1
	for i, prompt_message := range prompt_messages {
		if prompt_message.Role() == modelruntimeentities.PromptMessageRole_USER {
			last_user_message_idx = i
		}
	}
	// add function call messages before last user message
	example_messages := []modelruntimeentities.PromptMessager{}
	for _, example := range FUNCTION_CALLING_EXTRACTOR_EXAMPLE {
		id := uuid.NewV4().String()
		bindata, _ := json.Marshal(example["assistant"]["function_call"].(map[string]any)["parameters"])
		example_messages = append(example_messages, modelruntimeentities.NewUserPromptMessage(example["user"]["query"].(string), ""))
		example_messages = append(example_messages, modelruntimeentities.NewAssistantPromptMessage(example["assistant"]["text"].(string), "", []*modelruntimeentities.ToolCall{&modelruntimeentities.ToolCall{
			ID:   id,
			Type: "function",
			Function: modelruntimeentities.ToolCallFunction{
				Name:      example["assistant"]["function_call"].(map[string]any)["name"].(string),
				Arguments: string(bindata),
			},
		},
		}))
		example_messages = append(example_messages, modelruntimeentities.NewToolPromptMessage("Great! You have called the function with the correct parameters.", "", id))
		example_messages = append(example_messages, modelruntimeentities.NewAssistantPromptMessage("I have extracted the parameters, let's move on.", "", nil))
	}
	new_prompt_messages := []modelruntimeentities.PromptMessager{}
	new_prompt_messages = append(new_prompt_messages, prompt_messages[:last_user_message_idx]...)
	new_prompt_messages = append(new_prompt_messages, example_messages...)
	new_prompt_messages = append(new_prompt_messages, prompt_messages[last_user_message_idx:]...)

	// generate tool
	tool := &modelruntimeentities.PromptMessageTool{
		Name:        FUNCTION_CALLING_EXTRACTOR_NAME,
		Description: "Extract parameters from the natural language text",
		Parameters:  node_data.GetParameterJsonSchema(),
	}

	return new_prompt_messages, []*modelruntimeentities.PromptMessageTool{tool}
}
func (n *ParameterExtractorNode) _generate_prompt_engineering_completion_prompt(
	node_data *penodesentities.ParameterExtractorNodeData,
	query string,
	variable_pool *workflowentities.VariablePool,
	model_config *appconfigentities.ModelConfigWithCredentialsEntity,
	mem *memory.TokenBufferMemory,
	files *file.File,
) []modelruntimeentities.PromptMessager {
	/*
	   Generate completion prompt.
	*/
	bindata, _ := json.Marshal(node_data.GetParameterJsonSchema())
	rest_token := n._calculate_rest_token(node_data, query, variable_pool, model_config, "")
	prompt_template := n._get_prompt_engineering_prompt_template(node_data, query, variable_pool, mem, rest_token)
	var prompt_messages []modelruntimeentities.PromptMessager
	if real_prompt_template, ok := prompt_template.([]*promptentities.ChatModelMessage); ok {
		prompt_transform := prompt.NewAdvancedPromptTransform[[]*promptentities.ChatModelMessage](true, modelruntimeentities.ImagePromptMessageContentDETAIL(""))
		prompt_messages = prompt_transform.GetPrompt(
			real_prompt_template,
			map[string]string{"structure": string(bindata)},
			"",
			nil,
			"",
			node_data.Memory,
			mem,
			model_config,
		)
	} else if real_prompt_template, ok := prompt_template.(*promptentities.CompletionModelPromptTemplate); ok {
		prompt_transform := prompt.NewAdvancedPromptTransform[*promptentities.CompletionModelPromptTemplate](true, modelruntimeentities.ImagePromptMessageContentDETAIL(""))
		prompt_messages = prompt_transform.GetPrompt(
			real_prompt_template,
			map[string]string{"structure": string(bindata)},
			"",
			nil,
			"",
			node_data.Memory,
			mem,
			model_config,
		)
	} else {
		panic(penodesexceptions.NewParameterExtractorNodeError("prompt_template must be []*ChatModelMessage or *CompletionModelPromptTemplat"))
	}

	return prompt_messages
}
func (n *ParameterExtractorNode) _generate_prompt_engineering_chat_prompt(
	node_data *penodesentities.ParameterExtractorNodeData,
	query string,
	variable_pool *workflowentities.VariablePool,
	model_config *appconfigentities.ModelConfigWithCredentialsEntity,
	mem *memory.TokenBufferMemory,
	files *file.File,
) []modelruntimeentities.PromptMessager {
	/*
	   Generate chat prompt.
	*/

	bindata, _ := json.Marshal(node_data.GetParameterJsonSchema())
	rest_token := n._calculate_rest_token(node_data, query, variable_pool, model_config, "")
	prompt_template := n._get_prompt_engineering_prompt_template(node_data, strings.ReplaceAll(strings.ReplaceAll(CHAT_GENERATE_JSON_USER_MESSAGE_TEMPLATE, "{structure}", string(bindata)), "{text}", query), variable_pool, mem, rest_token)
	var prompt_messages []modelruntimeentities.PromptMessager
	if real_prompt_template, ok := prompt_template.([]*promptentities.ChatModelMessage); ok {
		prompt_transform := prompt.NewAdvancedPromptTransform[[]*promptentities.ChatModelMessage](true, modelruntimeentities.ImagePromptMessageContentDETAIL(""))
		prompt_messages = prompt_transform.GetPrompt(
			real_prompt_template,
			nil,
			"",
			nil,
			"",
			node_data.Memory,
			nil,
			model_config,
		)
	} else if real_prompt_template, ok := prompt_template.(*promptentities.CompletionModelPromptTemplate); ok {
		prompt_transform := prompt.NewAdvancedPromptTransform[*promptentities.CompletionModelPromptTemplate](true, modelruntimeentities.ImagePromptMessageContentDETAIL(""))
		prompt_messages = prompt_transform.GetPrompt(
			real_prompt_template,
			nil,
			"",
			nil,
			"",
			node_data.Memory,
			nil,
			model_config,
		)
	} else {
		panic(penodesexceptions.NewParameterExtractorNodeError("prompt_template must be []*ChatModelMessage or *CompletionModelPromptTemplat"))
	}

	// find last user message
	last_user_message_idx := len(prompt_messages) - 1
	for i, prompt_message := range prompt_messages {
		if prompt_message.Role() == modelruntimeentities.PromptMessageRole_USER {
			last_user_message_idx = i
		}
	}
	// add example messages before last user message
	example_messages := []modelruntimeentities.PromptMessager{}
	for _, example := range CHAT_EXAMPLE {
		bindata, _ := json.Marshal(example["user"]["json"])
		example_messages = append(example_messages, modelruntimeentities.NewUserPromptMessage(strings.ReplaceAll(strings.ReplaceAll(CHAT_GENERATE_JSON_USER_MESSAGE_TEMPLATE, "{structure}", string(bindata)), "{text}", example["user"]["query"].(string)), ""))
		bindata, _ = json.Marshal(example["assistant"]["json"])
		example_messages = append(example_messages, modelruntimeentities.NewAssistantPromptMessage(string(bindata), "", nil))
	}
	new_prompt_messages := []modelruntimeentities.PromptMessager{}
	new_prompt_messages = append(new_prompt_messages, prompt_messages[:last_user_message_idx]...)
	new_prompt_messages = append(new_prompt_messages, example_messages...)
	new_prompt_messages = append(new_prompt_messages, prompt_messages[last_user_message_idx:]...)

	return new_prompt_messages
}
func (n *ParameterExtractorNode) _generate_prompt_engineering_prompt(
	data *penodesentities.ParameterExtractorNodeData,
	query string,
	variable_pool *workflowentities.VariablePool,
	model_config *appconfigentities.ModelConfigWithCredentialsEntity,
	mem *memory.TokenBufferMemory,
	files *file.File,
) []modelruntimeentities.PromptMessager {
	/*
	   Generate prompt engineering prompt.
	*/
	model_mode := data.Model.Mode

	if model_mode == modelruntimeentities.LLMMode_COMPLETION {
		return n._generate_prompt_engineering_completion_prompt(
			data,
			query,
			variable_pool,
			model_config,
			mem,
			files,
		)
	} else if model_mode == modelruntimeentities.LLMMode_CHAT {
		return n._generate_prompt_engineering_chat_prompt(
			data,
			query,
			variable_pool,
			model_config,
			mem,
			files,
		)
	} else {
		panic(penodesexceptions.NewInvalidModelModeError(fmt.Sprintf("Invalid model mode: {%s}", model_mode)))
	}
}

func (n *ParameterExtractorNode) _validate_result(data *penodesentities.ParameterExtractorNodeData, result map[string]any) map[string]any {
	/*
	   Validate result.
	*/
	if len(data.Parameters) != len(result) {
		panic(penodesexceptions.NewInvalidNumberOfParametersError("Invalid number of parameters"))
	}
	for _, parameter := range data.Parameters {
		if _, ok := result[parameter.Name]; !ok && parameter.Required {
			panic(penodesexceptions.NewRequiredParameterMissingError(fmt.Sprintf("Parameter {%s} is required", parameter.Name)))
		}
		if parameter.Type == "select" && len(parameter.Options) > 0 {
			if _, ok := result[parameter.Name]; !ok {
				panic(penodesexceptions.NewInvalidSelectValueError(fmt.Sprintf("Invalid `select` value for parameter {%s}", parameter.Name)))
			} else {
				if _, ok := result[parameter.Name].(string); !ok {
					panic(penodesexceptions.NewInvalidSelectValueError(fmt.Sprintf("Invalid `select` value for parameter {%s}", parameter.Name)))
				} else {
					if !slices.Contains(parameter.Options, result[parameter.Name].(string)) {
						panic(penodesexceptions.NewInvalidSelectValueError(fmt.Sprintf("Invalid `select` value for parameter {%s}", parameter.Name)))
					}
				}
			}
		}
		if parameter.Type == "number" {
			if _, ok := result[parameter.Name]; !ok {
				panic(penodesexceptions.NewInvalidNumberValueError(fmt.Sprintf("Invalid `number` value for parameter {%s}", parameter.Name)))
			} else {
				switch result[parameter.Name].(type) {
				case int:
				case float32:
				case float64:
				default:
					panic(penodesexceptions.NewInvalidNumberValueError(fmt.Sprintf("Invalid `number` value for parameter {%s}", parameter.Name)))
				}
			}
		}
		if parameter.Type == "bool" {
			if _, ok := result[parameter.Name]; !ok {
				panic(penodesexceptions.NewInvalidBoolValueError(fmt.Sprintf("Invalid `bool` value for parameter {%s}", parameter.Name)))
			} else {
				switch result[parameter.Name].(type) {
				case bool:
				default:
					panic(penodesexceptions.NewInvalidBoolValueError(fmt.Sprintf("Invalid `bool` value for parameter {%s}", parameter.Name)))
				}
			}
		}
		if parameter.Type == "string" {
			if _, ok := result[parameter.Name]; !ok {
				panic(penodesexceptions.NewInvalidStringValueError(fmt.Sprintf("Invalid `string` value for parameter {%s}", parameter.Name)))
			} else {
				switch result[parameter.Name].(type) {
				case string:
				default:
					panic(penodesexceptions.NewInvalidStringValueError(fmt.Sprintf("Invalid `string` value for parameter {%s}", parameter.Name)))
				}
			}
		}
		if strings.HasPrefix(parameter.Type, "array") {
			if _, ok := result[parameter.Name]; !ok {
				panic(penodesexceptions.NewInvalidArrayValueError(fmt.Sprintf("Invalid `array` value for parameter {%s}", parameter.Name)))
			} else {
				nested_type := parameter.Type[6 : len(parameter.Type)-1]
				switch parameters := result[parameter.Name].(type) {
				case []int:
					if nested_type != "number" {
						panic(penodesexceptions.NewInvalidArrayValueError(fmt.Sprintf("Invalid `array[number]` value for parameter {%s}", parameter.Name)))
					}
				case []float32:
					if nested_type != "number" {
						panic(penodesexceptions.NewInvalidArrayValueError(fmt.Sprintf("Invalid `array[number]` value for parameter {%s}", parameter.Name)))
					}
				case []float64:
					if nested_type != "number" {
						panic(penodesexceptions.NewInvalidArrayValueError(fmt.Sprintf("Invalid `array[number]` value for parameter {%s}", parameter.Name)))
					}
				case []string:
					if nested_type != "string" {
						panic(penodesexceptions.NewInvalidArrayValueError(fmt.Sprintf("Invalid `array[string]` value for parameter {%s}", parameter.Name)))
					}
				case []map[string]any:
					if nested_type != "object" {
						panic(penodesexceptions.NewInvalidArrayValueError(fmt.Sprintf("Invalid `array[object]` value for parameter {%s}", parameter.Name)))
					}
				case []any:
					for _, item := range parameters {
						switch item.(type) {
						case int:
							if nested_type != "number" {
								panic(penodesexceptions.NewInvalidArrayValueError(fmt.Sprintf("Invalid `array[number]` value for parameter {%s}", parameter.Name)))
							}
						case float32:
							if nested_type != "number" {
								panic(penodesexceptions.NewInvalidArrayValueError(fmt.Sprintf("Invalid `array[number]` value for parameter {%s}", parameter.Name)))
							}
						case float64:
							if nested_type != "number" {
								panic(penodesexceptions.NewInvalidArrayValueError(fmt.Sprintf("Invalid `array[number]` value for parameter {%s}", parameter.Name)))
							}
						case string:
							if nested_type != "string" {
								panic(penodesexceptions.NewInvalidArrayValueError(fmt.Sprintf("Invalid `array[string]` value for parameter {%s}", parameter.Name)))
							}
						case map[string]any:
							if nested_type != "object" {
								panic(penodesexceptions.NewInvalidArrayValueError(fmt.Sprintf("Invalid `array[object]` value for parameter {%s}", parameter.Name)))
							}
						default:
							panic(penodesexceptions.NewInvalidArrayValueError(fmt.Sprintf("Invalid `array[%#v]` value for parameter {%s}", item, parameter.Name)))
						}
					}
				default:
					panic(penodesexceptions.NewInvalidArrayValueError(fmt.Sprintf("Invalid `array[%#v]` value for parameter {%s}", result[parameter.Name], parameter.Name)))
				}
			}
		}
	}
	return result
}
func (n *ParameterExtractorNode) _transform_result(data *penodesentities.ParameterExtractorNodeData, result map[string]any) map[string]any {
	/*
	   Transform result into standard format.
	*/
	transformed_result := map[string]any{}
	for _, parameter := range data.Parameters {
		if _, ok := result[parameter.Name]; ok {
			// transform value
			if parameter.Type == "number" {
				switch real_val := result[parameter.Name].(type) {
				case int:
					transformed_result[parameter.Name] = result[parameter.Name]
				case float32:
					transformed_result[parameter.Name] = result[parameter.Name]
				case float64:
					transformed_result[parameter.Name] = result[parameter.Name]
				case string:
					if strings.Contains(real_val, ".") {
						tmp_val, err := strconv.ParseFloat(real_val, 64)
						if err != nil {
							mlog.Errorf("string(%s) to float failed:%v", real_val, err)
							continue
						} else {
							result[parameter.Name] = tmp_val
							transformed_result[parameter.Name] = result[parameter.Name]
						}
					} else {
						tmp_val, err := strconv.Atoi(real_val)
						if err != nil {
							mlog.Errorf("string(%s) to int failed:%v", real_val, err)
							continue
						} else {
							result[parameter.Name] = tmp_val
							transformed_result[parameter.Name] = result[parameter.Name]
						}
					}
				default:
					continue
				}
				// TODO: bool is not supported in the current version
				// }else if parameter.Type == 'bool'{
				//     if isinstance(result[parameter.Name], bool){
				//         transformed_result[parameter.Name] = bool(result[parameter.Name])
				//     }else if isinstance(result[parameter.Name], string){
				//         if result[parameter.Name].lower() in ['true', 'false']{
				//             transformed_result[parameter.Name] = bool(result[parameter.Name].lower() == 'true')
				//     }else if isinstance(result[parameter.Name], int){
				//         transformed_result[parameter.Name] = bool(result[parameter.Name])
			} else if slices.Contains([]string{"string", "select"}, parameter.Type) {
				if _, ok := result[parameter.Name].(string); ok {
					transformed_result[parameter.Name] = result[parameter.Name]
				}
			} else if strings.HasPrefix(parameter.Type, "array") {
				nested_type := parameter.Type[6 : len(parameter.Type)-1]
				switch result[parameter.Name].(type) {
				case []int:
					if nested_type == "number" {
						transformed_result[parameter.Name] = result[parameter.Name]
					}
				case []float32:
					if nested_type == "number" {
						panic(penodesexceptions.NewInvalidArrayValueError(fmt.Sprintf("Invalid `array[number]` value for parameter {%s}", parameter.Name)))
					}
				case []float64:
					if nested_type == "number" {
						transformed_result[parameter.Name] = result[parameter.Name]
					}
				case []string:
					if nested_type == "string" {
						transformed_result[parameter.Name] = result[parameter.Name]
					}
				case []map[string]any:
					if nested_type == "object" {
						transformed_result[parameter.Name] = result[parameter.Name]
					}
				case []any:
					tmp_result := []any{}
					for _, item := range result[parameter.Name].([]any) {
						switch real_val := item.(type) {
						case int:
							if nested_type == "number" {
								tmp_result = append(tmp_result, real_val)
							}
						case float32:
							if nested_type == "number" {
								tmp_result = append(tmp_result, real_val)
							}
						case float64:
							if nested_type == "number" {
								tmp_result = append(tmp_result, real_val)
							}
						case string:
							if nested_type == "number" {
								if strings.Contains(real_val, ".") {
									tmp_val, err := strconv.ParseFloat(real_val, 64)
									if err != nil {
										mlog.Errorf("string(%s) to float failed:%v", real_val, err)
										continue
									} else {
										tmp_result = append(tmp_result, tmp_val)
									}
								} else {
									tmp_val, err := strconv.Atoi(real_val)
									if err != nil {
										mlog.Errorf("string(%s) to int failed:%v", real_val, err)
										continue
									} else {
										tmp_result = append(tmp_result, tmp_val)
									}
								}
							} else if nested_type == "string" {
								tmp_result = append(tmp_result, real_val)
							}
						case map[string]any:
							if nested_type == "object" {
								tmp_result = append(tmp_result, real_val)
							}
						default:
							continue
						}
					}
					transformed_result[parameter.Name] = tmp_result
				}
			}
		}
		if _, ok := transformed_result[parameter.Name]; !ok {
			if parameter.Type == "number" {
				transformed_result[parameter.Name] = 0
			} else if parameter.Type == "bool" {
				transformed_result[parameter.Name] = false
			} else if slices.Contains([]string{"string", "select"}, parameter.Type) {
				transformed_result[parameter.Name] = ""
			} else if strings.HasPrefix(parameter.Type, "array") {
				transformed_result[parameter.Name] = []any{}
			}
		}
	}
	return transformed_result
}
func extract_json(text string) string {
	/*
		From a given JSON started from '{' or '[' extract the complete JSON object.
	*/
	stack := []byte{}
	for i, c := range []byte(text) {
		if slices.Contains([]byte{'{', '['}, c) {
			stack = append(stack, c)
		} else if slices.Contains([]byte{'}', ']'}, c) {
			// check if stack is empty
			if len(stack) == 0 {
				return text[:i]
			}
			// check if the last element in stack is matching
			if (c == '}' && stack[len(stack)-1] == '{') || (c == ']' && stack[len(stack)-1] == '[') {
				stack = stack[:len(stack)-1]
				if len(stack) == 0 {
					return text[:i+1]
				}
			} else {
				return text[:i]
			}
		}
	}
	return ""
}
func (n *ParameterExtractorNode) _extract_complete_json_response(result string) map[string]any {
	/*
	   Extract complete json response.
	*/

	// extract json from the text
	for idx, c := range []byte(result) {
		if c == '{' || c == '[' {
			json_str := extract_json(result[idx:])
			if json_str != "" {
				tmpmap := map[string]any{}
				err := json.Unmarshal([]byte(json_str), &tmpmap)
				if err != nil {
					mlog.Errorf("json unmarshal(%s) failed:%v", json_str, err)
					continue
				} else {
					return tmpmap
				}
			}
		}
	}
	return nil
}
func (n *ParameterExtractorNode) _extract_json_from_tool_call(tool_call *modelruntimeentities.ToolCall) map[string]any {
	/*
	   Extract json from tool call.
	*/
	if tool_call == nil || tool_call.Function.Arguments == "" {
		return nil
	}
	tmpmap := map[string]any{}
	err := json.Unmarshal([]byte(tool_call.Function.Arguments), &tmpmap)
	if err != nil {
		mlog.Errorf("json unmarshal(%s) failed:%v", tool_call.Function.Arguments, err)
		return nil
	} else {
		return tmpmap
	}
}
func (n *ParameterExtractorNode) _generate_default_result(data *penodesentities.ParameterExtractorNodeData) map[string]any {
	/*
	   Generate default result.
	*/
	result := map[string]any{}
	for _, parameter := range data.Parameters {
		if parameter.Type == "number" {
			result[parameter.Name] = 0
		} else if parameter.Type == "bool" {
			result[parameter.Name] = false
		} else if slices.Contains([]string{"string", "select"}, parameter.Type) {
			result[parameter.Name] = ""
		}
	}
	return result
}

func New() *ParameterExtractorNode {
	return &ParameterExtractorNode{
		BaseNode: &base.BaseNode[*penodesentities.ParameterExtractorNodeData]{},
	}
}
