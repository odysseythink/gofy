package questionclassifier

import (
	"encoding/json"
	"fmt"
	"iter"
	"math"
	"slices"
	"strings"

	"github.com/odysseythink/gofy/backend/core/exceptions"
	qcfnodesexceptions "github.com/odysseythink/gofy/backend/core/exceptions/nodes/question_classifier"
	modelmanager "github.com/odysseythink/gofy/backend/core/manageres/model_manager"
	"github.com/odysseythink/gofy/backend/core/memory"
	"github.com/odysseythink/gofy/backend/core/prompt"
	promptutils "github.com/odysseythink/gofy/backend/core/prompt/utils"
	"github.com/odysseythink/gofy/backend/core/variables"
	"github.com/odysseythink/gofy/backend/core/workflow/nodes/base"
	"github.com/odysseythink/gofy/backend/core/workflow/nodes/llm"
	variabletemplateparser "github.com/odysseythink/gofy/backend/core/workflow/utils/variable_template_parser"
	appconfigentities "github.com/odysseythink/gofy/backend/entities/app/config"
	modelruntimeentities "github.com/odysseythink/gofy/backend/entities/model_runtime"
	eventnodesentities "github.com/odysseythink/gofy/backend/entities/nodes/event"
	llmnodesentities "github.com/odysseythink/gofy/backend/entities/nodes/llm"
	qcfnodesentities "github.com/odysseythink/gofy/backend/entities/nodes/question_classifier"
	promptentities "github.com/odysseythink/gofy/backend/entities/prompt"
	workflowentities "github.com/odysseythink/gofy/backend/entities/workflow"
	modelruntimeenumtypes "github.com/odysseythink/gofy/backend/enum_types/model_runtime"
	nodesenumtypes "github.com/odysseythink/gofy/backend/enum_types/nodes"
	workflowenumtypes "github.com/odysseythink/gofy/backend/enum_types/workflow"
	"github.com/odysseythink/gofy/backend/models"
	"github.com/odysseythink/gofy/backend/utils"
	"github.com/odysseythink/mlog"
)

type QuestionClassifierNode struct {
	*base.BaseNode[*qcfnodesentities.QuestionClassifierNodeData]
}

func (n *QuestionClassifierNode) Type() nodesenumtypes.NodeType {
	return nodesenumtypes.Node_QUESTION_CLASSIFIER
}

func (n *QuestionClassifierNode) Run() (*workflowentities.NodeRunResult, iter.Seq[any]) {
	node_data := n.NodeData
	variable_pool := n.GetGraphRuntimeState().VariablePool
	// extract variables
	var variable variables.Variabler
	if node_data.QueryVariableSelector != nil {
		variable = variable_pool.Get(node_data.QueryVariableSelector)
	}
	query := ""
	if variable != nil {
		tmp := variable.GetValue()
		if _, ok := tmp.(string); ok {
			query = tmp.(string)
		}
	}

	inputs_variables := map[string]any{"query": query}
	// fetch model config
	model_instance, model_config := llm.FetchModelConfig(n.GetTenantID(), &node_data.Model)
	// fetch memory
	mem := llm.FetchMemory(n.GetGraphRuntimeState(), n.GetAppID(),
		node_data.Memory,
		model_instance,
	)
	// fetch instruction
	node_data.Instruction = variable_pool.ConvertTemplate(node_data.Instruction).Text()

	// fetch prompt messages
	rest_token := n._calculate_rest_token(
		node_data,
		query,
		model_config,
		"",
	)
	prompt_template := n._get_prompt_template(
		node_data,
		query,
		mem,
		rest_token,
	)
	var prompt_messages []modelruntimeentities.PromptMessager
	var stop []string
	result_text := ""
	usage := modelruntimeentities.NewLLMUsage()
	finish_reason := ""
	if real_prompt_template, ok := prompt_template.([]*llmnodesentities.LLMNodeChatModelMessage); ok {
		prompt_messages, stop = llm.FetchPromptMessages(
			query,
			"",
			mem,
			model_config,
			real_prompt_template,
			nil,
			node_data.Vision.Enabled,
			node_data.Vision.Configs.Detail,
			variable_pool,
		)

	} else if real_prompt_template, ok := prompt_template.(*llmnodesentities.LLMNodeCompletionModelPromptTemplate); ok {
		prompt_messages, stop = llm.FetchPromptMessages(
			query,
			"",
			mem,
			model_config,
			real_prompt_template,
			nil,
			node_data.Vision.Enabled,
			node_data.Vision.Configs.Detail,
			variable_pool,
		)

	} else {
		mlog.Errorf("unsuported PromptTemplate=%#v", prompt_template)
		return &workflowentities.NodeRunResult{
			Status: models.WorkflowNodeExecutionStatus_FAILED,
			Inputs: inputs_variables,
			Error:  "unsuported PromptTemplate",
			Metadata: map[workflowenumtypes.NodeRunMetadataKey]any{
				workflowenumtypes.NodeRunMetadataKey_TOTAL_TOKENS: usage.TotalTokens,
				workflowenumtypes.NodeRunMetadataKey_TOTAL_PRICE:  usage.TotalPrice,
				workflowenumtypes.NodeRunMetadataKey_CURRENCY:     usage.Currency,
			},
			LLmUsage: usage,
		}, nil
	}

	res := func() (res *workflowentities.NodeRunResult) {
		defer func() {
			if r := recover(); r != nil {
				if exp, ok := r.(error); ok {
					mlog.Errorf("accur error:%v", exp)
					res = &workflowentities.NodeRunResult{
						Status: models.WorkflowNodeExecutionStatus_FAILED,
						Inputs: inputs_variables,
						Error:  exp.Error(),
						Metadata: map[workflowenumtypes.NodeRunMetadataKey]any{
							workflowenumtypes.NodeRunMetadataKey_TOTAL_TOKENS: usage.TotalTokens,
							workflowenumtypes.NodeRunMetadataKey_TOTAL_PRICE:  usage.TotalPrice,
							workflowenumtypes.NodeRunMetadataKey_CURRENCY:     usage.Currency,
						},
						LLmUsage: usage,
					}
				} else {
					panic(r)
				}
			}
		}()
		// handle invoke result
		generator := llm.InvokeLLM(
			"",
			n.GetNodeID(),
			&node_data.Model,
			model_instance,
			prompt_messages,
			stop,
		)
		for event := range generator {
			if real_event, ok := any(event).(*eventnodesentities.ModelInvokeCompletedEvent); ok {
				result_text = real_event.Text
				usage = real_event.Usage
				finish_reason = real_event.FinishReason
				break
			}
		}
		category_name := node_data.Classes[0].Name
		category_id := node_data.Classes[0].ID
		result_text_json := utils.ParseAndCheckJsonMarkdown(result_text, nil)
		// result_text_json = json.loads(result_text.strip('```JSON\n'))
		_, ok1 := result_text_json["category_name"]
		_, ok2 := result_text_json["category_id"]
		if ok1 && ok2 {
			category_id_result := result_text_json["category_id"].(string)
			classes := node_data.Classes
			classes_map := map[string]string{}
			category_ids := []string{}
			for _, class_ := range classes {
				classes_map[class_.ID] = class_.Name
				category_ids = append(category_ids, class_.ID)
			}

			if slices.Contains(category_ids, category_id_result) {
				category_name = classes_map[category_id_result]
				category_id = category_id_result
			}
		}
		process_data := map[string]any{
			"model_mode": model_config.Mode,
			"prompts": promptutils.PromptMessagesToPromptForSaving(
				model_config.Mode, prompt_messages,
			),
			"usage":         usage,
			"finish_reason": finish_reason,
		}
		outputs := map[string]any{"class_name": category_name, "class_id": category_id}
		res = &workflowentities.NodeRunResult{
			Status:           models.WorkflowNodeExecutionStatus_SUCCEEDED,
			Inputs:           inputs_variables,
			ProcessData:      process_data,
			Outputs:          outputs,
			EdgeSourceHandle: category_id,
			Metadata: map[workflowenumtypes.NodeRunMetadataKey]any{
				workflowenumtypes.NodeRunMetadataKey_TOTAL_TOKENS: usage.TotalTokens,
				workflowenumtypes.NodeRunMetadataKey_TOTAL_PRICE:  usage.TotalPrice,
				workflowenumtypes.NodeRunMetadataKey_CURRENCY:     usage.Currency,
			},
			LLmUsage: usage,
		}
		return res
	}()
	mlog.Debugf("------res=%#v", res)
	return res, nil
}

func (n *QuestionClassifierNode) ExtractVariableSelectorToVariableMapping(graph_config map[string]any, node_id string, node_data *qcfnodesentities.QuestionClassifierNodeData) map[string][]string {
	variable_mapping := map[string][]string{"query": node_data.QueryVariableSelector}
	variable_selectors := []*workflowentities.VariableSelector{}
	if node_data.Instruction != "" {
		variable_template_parser := variabletemplateparser.NewVariableTemplateParser(node_data.Instruction)
		variable_selectors = append(variable_selectors, variable_template_parser.ExtractVariableSelectors()...)
	}
	for _, variable_selector := range variable_selectors {
		variable_mapping[variable_selector.Variable] = variable_selector.ValueSelector
	}
	new_variable_mapping := map[string][]string{}
	for key, value := range variable_mapping {
		new_variable_mapping[node_id+"."+key] = value
	}

	return new_variable_mapping
}

func (n *QuestionClassifierNode) GetDefaultConfig(filters map[string]any) map[string]any {
	return map[string]any{"type": "question-classifier", "config": map[string]any{"instructions": ""}}
}

func (n *QuestionClassifierNode) _calculate_rest_token(
	node_data *qcfnodesentities.QuestionClassifierNodeData,
	query string,
	model_config *appconfigentities.ModelConfigWithCredentialsEntity,
	context string,
) int {
	var prompt_messages []modelruntimeentities.PromptMessager
	prompt_template := n._get_prompt_template(node_data, query, nil, 2000)
	if real_prompt_template, ok := prompt_template.([]*llmnodesentities.LLMNodeChatModelMessage); ok {
		msges := []*promptentities.ChatModelMessage{}
		for _, v := range real_prompt_template {
			msges = append(msges, v.ChatModelMessage)
		}
		prompt_transform := prompt.NewAdvancedPromptTransform[[]*promptentities.ChatModelMessage](true, modelruntimeentities.ImagePromptMessageContentDETAIL(""))
		prompt_messages = prompt_transform.GetPrompt(
			msges,
			map[string]string{},
			"",
			nil,
			context,
			node_data.Memory,
			nil,
			model_config,
		)
	} else if real_prompt_template, ok := prompt_template.(*llmnodesentities.LLMNodeCompletionModelPromptTemplate); ok {
		prompt_transform := prompt.NewAdvancedPromptTransform[*promptentities.CompletionModelPromptTemplate](true, modelruntimeentities.ImagePromptMessageContentDETAIL(""))
		prompt_messages = prompt_transform.GetPrompt(
			real_prompt_template.CompletionModelPromptTemplate,
			map[string]string{},
			"",
			nil,
			context,
			node_data.Memory,
			nil,
			model_config,
		)
	} else {
		mlog.Errorf("unsupported prompt_template type:%#v", prompt_template)
		panic(exceptions.NewValueError("unsupported prompt_template type"))
	}

	rest_tokens := 2000

	if _, ok := model_config.ModelSchema.ModelProperties[modelruntimeenumtypes.ModelPropertyKey_CONTEXT_SIZE]; ok {
		if model_context_tokens, ok := model_config.ModelSchema.ModelProperties[modelruntimeenumtypes.ModelPropertyKey_CONTEXT_SIZE].(int); ok {
			model_instance := modelmanager.NewModelInstance(
				model_config.ProviderModelBundle, model_config.Model,
			)

			curr_message_tokens := model_instance.GetLLMNumTokens(prompt_messages, nil)

			max_tokens := 0
			for _, parameter_rule := range model_config.ModelSchema.ParameterRules {
				if parameter_rule.Name == "max_tokens" || (parameter_rule.UseTemplate != "" && parameter_rule.UseTemplate == "max_tokens") {
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
			rest_tokens = int(math.Max(float64(rest_tokens), 0.0))
		}
	}
	return rest_tokens
}

func (n *QuestionClassifierNode) _get_prompt_template(
	node_data *qcfnodesentities.QuestionClassifierNodeData,
	query string,
	mem *memory.TokenBufferMemory,
	max_token_limit int, /*= 2000*/
) any {
	if max_token_limit <= 0 {
		max_token_limit = 2000
	}
	model_mode := node_data.Model.Mode
	classes := node_data.Classes
	categories := []map[string]string{}
	for _, class := range classes {
		category := map[string]string{"category_id": class.ID, "category_name": class.Name}
		categories = append(categories, category)
	}
	instruction := node_data.Instruction
	input_text := query
	memory_str := ""
	if mem != nil {
		message_limit := 0
		if node_data.Memory != nil {
			message_limit = node_data.Memory.Window.Size
		}
		memory_str = mem.GetHistoryPromptText(
			"",
			"",
			max_token_limit,
			message_limit,
		)
	}

	if model_mode == modelruntimeentities.LLMMode_CHAT {
		prompt_messages := []*llmnodesentities.LLMNodeChatModelMessage{}
		system_prompt_messages := &llmnodesentities.LLMNodeChatModelMessage{
			ChatModelMessage: &promptentities.ChatModelMessage{
				Role: modelruntimeentities.PromptMessageRole_SYSTEM,
				Text: strings.ReplaceAll(QUESTION_CLASSIFIER_SYSTEM_PROMPT, "{histories}", memory_str),
			},
		}
		prompt_messages = append(prompt_messages, system_prompt_messages)
		user_prompt_message_1 := &llmnodesentities.LLMNodeChatModelMessage{
			ChatModelMessage: &promptentities.ChatModelMessage{
				Role: modelruntimeentities.PromptMessageRole_USER,
				Text: QUESTION_CLASSIFIER_USER_PROMPT_1,
			},
		}
		prompt_messages = append(prompt_messages, user_prompt_message_1)
		assistant_prompt_message_1 := &llmnodesentities.LLMNodeChatModelMessage{
			ChatModelMessage: &promptentities.ChatModelMessage{
				Role: modelruntimeentities.PromptMessageRole_ASSISTANT,
				Text: QUESTION_CLASSIFIER_ASSISTANT_PROMPT_1,
			},
		}
		prompt_messages = append(prompt_messages, assistant_prompt_message_1)
		user_prompt_message_2 := &llmnodesentities.LLMNodeChatModelMessage{
			ChatModelMessage: &promptentities.ChatModelMessage{
				Role: modelruntimeentities.PromptMessageRole_USER,
				Text: QUESTION_CLASSIFIER_USER_PROMPT_2,
			},
		}
		prompt_messages = append(prompt_messages, user_prompt_message_2)
		assistant_prompt_message_2 := &llmnodesentities.LLMNodeChatModelMessage{
			ChatModelMessage: &promptentities.ChatModelMessage{
				Role: modelruntimeentities.PromptMessageRole_ASSISTANT,
				Text: QUESTION_CLASSIFIER_ASSISTANT_PROMPT_2,
			},
		}
		prompt_messages = append(prompt_messages, assistant_prompt_message_2)
		bindata, _ := json.Marshal(categories)
		user_prompt_message_3 := &llmnodesentities.LLMNodeChatModelMessage{
			ChatModelMessage: &promptentities.ChatModelMessage{
				Role: modelruntimeentities.PromptMessageRole_USER,
				Text: strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(QUESTION_CLASSIFIER_USER_PROMPT_3, "{input_text}", input_text),
					"{categories}", string(bindata)),
					"{classification_instructions}", instruction),
			},
		}
		prompt_messages = append(prompt_messages, user_prompt_message_3)
		return prompt_messages
	} else if model_mode == modelruntimeentities.LLMMode_COMPLETION {
		bindata, _ := json.Marshal(categories)
		return &llmnodesentities.LLMNodeCompletionModelPromptTemplate{
			CompletionModelPromptTemplate: &promptentities.CompletionModelPromptTemplate{
				Text: strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(QUESTION_CLASSIFIER_COMPLETION_PROMPT,
					"{histories}", memory_str),
					"{input_text}", input_text),
					"{categories}", string(bindata)),
					"{classification_instructions}", instruction),
			},
		}

	} else {
		panic(qcfnodesexceptions.NewInvalidModelTypeError(fmt.Sprintf("Model mode %s not support.", model_mode)))
	}
}
func New() *QuestionClassifierNode {
	return &QuestionClassifierNode{
		BaseNode: &base.BaseNode[*qcfnodesentities.QuestionClassifierNodeData]{},
	}
}
