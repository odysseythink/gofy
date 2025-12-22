package llm

import (
	"encoding/json"
	"fmt"
	"iter"
	"slices"
	"strings"

	"gorm.io/gorm"
	"mlib.com/confy/cast"
	"mlib.com/gofy/server/constants"
	"mlib.com/gofy/server/core/exceptions"
	llmnodesexceptions "mlib.com/gofy/server/core/exceptions/nodes/llm"
	modelmanager "mlib.com/gofy/server/core/manageres/model_manager"
	providermanager "mlib.com/gofy/server/core/manageres/provider_manager"
	"mlib.com/gofy/server/core/memory"
	promptutils "mlib.com/gofy/server/core/prompt/utils"
	"mlib.com/gofy/server/core/variables"
	"mlib.com/gofy/server/core/workflow/nodes/base"
	variabletemplateparserutils "mlib.com/gofy/server/core/workflow/utils/variable_template_parser"
	dbengine "mlib.com/gofy/server/db_engine"
	appconfigentities "mlib.com/gofy/server/entities/app/config"
	graphengineentities "mlib.com/gofy/server/entities/graph_engine"
	modelruntimeentities "mlib.com/gofy/server/entities/model_runtime"
	nodesevententities "mlib.com/gofy/server/entities/nodes/event"
	llmnodesentities "mlib.com/gofy/server/entities/nodes/llm"
	promptentities "mlib.com/gofy/server/entities/prompt"
	ragentities "mlib.com/gofy/server/entities/rag"
	workflowentities "mlib.com/gofy/server/entities/workflow"
	modelenumtypes "mlib.com/gofy/server/enum_types/model"
	modelruntimeenumtypes "mlib.com/gofy/server/enum_types/model_runtime"
	nodesenumtypes "mlib.com/gofy/server/enum_types/nodes"
	providerenumtypes "mlib.com/gofy/server/enum_types/provider"
	variableenumtypes "mlib.com/gofy/server/enum_types/variable"
	workflowenumtypes "mlib.com/gofy/server/enum_types/workflow"
	"mlib.com/gofy/server/models"
	"mlib.com/gofy/server/utils"
	"mlib.com/mlog"
)

type LLMNode struct {
	*base.BaseNode[*llmnodesentities.LLMNodeData]
}

func (n *LLMNode) Type() nodesenumtypes.NodeType {
	return nodesenumtypes.Node_LLM
}

func DeductLLMQuota(tenant_id string, model_instance *modelmanager.ModelInstance, usage *modelruntimeentities.LLMUsage) {
	provider_model_bundle := model_instance.ProviderModelBundle
	provider_configuration := provider_model_bundle.Configuration

	if provider_configuration.UsingProviderType != providerenumtypes.Provider_SYSTEM {
		return
	}
	system_configuration := provider_configuration.SystemConfiguration

	var quota_unit providerenumtypes.QuotaUnitType
	for _, quota_configuration := range system_configuration.QuotaConfigurations {
		if quota_configuration.QuotaType == system_configuration.CurrentQuotaType {
			quota_unit = quota_configuration.QuotaUnit

			if quota_configuration.QuotaLimit == -1 {
				return
			}
			break
		}
	}
	var used_quota int
	if string(quota_unit) != "" {
		if quota_unit == providerenumtypes.QuotaUnit_TOKENS {
			used_quota = usage.TotalTokens
		} else if quota_unit == providerenumtypes.QuotaUnit_CREDITS {
			used_quota = 1
		} else {
			used_quota = 1
		}
	}
	if used_quota > 0 && string(system_configuration.CurrentQuotaType) != "" {
		dbengine.Instance().DB.Debug().Model(&models.Provider{}).Update("quota_used", gorm.Expr("quota_used + ?", used_quota)).Where("tenant_id = ? and provider_name = ? and provider_type=? and quota_type = ? and quota_limit > quota_used", tenant_id, model_instance.Provider, providerenumtypes.Provider_SYSTEM, system_configuration.CurrentQuotaType)
	}
}
func (n *LLMNode) Run() (*workflowentities.NodeRunResult, iter.Seq[any]) {
	return nil, func(yield func(any) bool) {
		node_inputs := map[string]any{}
		process_data := map[string]any{}
		result_text := ""
		finish_reason := ""
		var usage *modelruntimeentities.LLMUsage
		func() {
			defer func() {
				if r := recover(); r != nil {
					mlog.Errorf("panic recover:%v\n%s", r, utils.GetCurrentGoroutineStack())
					if real_exp, ok := r.(*llmnodesexceptions.LLMNodeError); ok {
						yield(&nodesevententities.RunCompletedEvent{
							RunResult: &workflowentities.NodeRunResult{
								Status:      models.WorkflowNodeExecutionStatus_FAILED,
								Error:       real_exp.Error(),
								Inputs:      node_inputs,
								ProcessData: process_data,
							},
						})
					} else if real_exp, ok := r.(error); ok {
						yield(&nodesevententities.RunCompletedEvent{
							RunResult: &workflowentities.NodeRunResult{
								Status:      models.WorkflowNodeExecutionStatus_FAILED,
								Error:       real_exp.Error(),
								Inputs:      node_inputs,
								ProcessData: process_data,
							},
						})
					} else {
						panic(r)
					}
				}
			}()
			// init messages template
			// n.NodeData.PromptTemplate = n._transform_chat_messages(n.NodeData.PromptTemplate)

			// fetch variables and fetch values from variable pool
			inputs, err := _fetch_inputs(n, n.NodeData)
			if err != nil {
				yield(&nodesevententities.RunCompletedEvent{
					RunResult: &workflowentities.NodeRunResult{
						Status:      models.WorkflowNodeExecutionStatus_FAILED,
						Error:       err.Error(),
						Inputs:      node_inputs,
						ProcessData: process_data,
					},
				})
				return
			}
			// fetch jinja2 inputs
			// jinja_inputs = self._fetch_jinja_inputs(node_data=n.NodeData)

			// merge inputs
			// inputs.update(jinja_inputs)
			mlog.Debugf("------inputs=%#v", inputs)
			mlog.Debugf("------node_inputs=%#v", node_inputs)
			// fetch files
			// files = (
			// 	self._fetch_files(selector=n.NodeData.vision.configs.variable_selector)
			// 	if n.NodeData.vision.enabled
			// 	else []
			// )

			// if files:
			// 	node_inputs["#files#"] = [file.to_dict() for file in files]

			// fetch context value
			generator := _fetch_context(n, n.NodeData)
			context := ""
			for event, err := range generator {
				if err != nil {
					yield(&nodesevententities.RunCompletedEvent{
						RunResult: &workflowentities.NodeRunResult{
							Status:      models.WorkflowNodeExecutionStatus_FAILED,
							Error:       err.Error(),
							Inputs:      node_inputs,
							ProcessData: process_data,
						},
					})
					return

				}
				context = event.Context
				if !yield(event) {
					return
				}
			}
			if context != "" {
				node_inputs["#context#"] = context
			}
			// fetch model config
			model_instance, model_config := FetchModelConfig(n.GetTenantID(), n.NodeData.Model)
			// fetch memory
			mem := FetchMemory(n.GetGraphRuntimeState(), n.GetAppID(), n.NodeData.Memory, model_instance)

			query := ""
			if n.NodeData.Memory != nil {
				query = n.NodeData.Memory.QueryPromptTemplate
				query_variable := n.GetGraphRuntimeState().VariablePool.Get([]string{constants.SYSTEM_VARIABLE_NODE_ID, string(workflowenumtypes.SystemVariableKey_QUERY)})
				if query == "" && query_variable != nil {
					query = query_variable.Text()
				}
			}
			var prompt_messages []modelruntimeentities.PromptMessager
			var stop []string
			if real_prompt_template, ok := n.NodeData.PromptTemplate.([]*llmnodesentities.LLMNodeChatModelMessage); ok {
				prompt_messages, stop = FetchPromptMessages(
					query,
					context,
					mem,
					model_config,
					real_prompt_template,
					n.NodeData.Memory,
					n.NodeData.Vision.Enabled,
					n.NodeData.Vision.Configs.Detail,
					n.GetGraphRuntimeState().VariablePool,
				)

			} else if real_prompt_template, ok := n.NodeData.PromptTemplate.(*llmnodesentities.LLMNodeCompletionModelPromptTemplate); ok {
				prompt_messages, stop = FetchPromptMessages(
					query,
					context,
					mem,
					model_config,
					real_prompt_template,
					n.NodeData.Memory,
					n.NodeData.Vision.Enabled,
					n.NodeData.Vision.Configs.Detail,
					n.GetGraphRuntimeState().VariablePool,
				)

			} else {
				mlog.Errorf("unsuported PromptTemplate=%#v", n.NodeData.PromptTemplate)
				yield(&nodesevententities.RunCompletedEvent{
					RunResult: &workflowentities.NodeRunResult{
						Status:      models.WorkflowNodeExecutionStatus_FAILED,
						Error:       fmt.Sprintf("unsuported PromptTemplate=%#v", n.NodeData.PromptTemplate),
						Inputs:      node_inputs,
						ProcessData: process_data,
					},
				})
				return
			}

			process_data = map[string]any{
				"model_mode": model_config.Mode,
				"prompts": promptutils.PromptMessagesToPromptForSaving(
					model_config.Mode, prompt_messages,
				),
				"model_provider": model_config.Provider,
				"model_name":     model_config.Model,
			}

			// handle invoke result
			llm_generator := InvokeLLM(
				"",
				n.GetID(),
				n.NodeData.Model,
				model_instance,
				prompt_messages,
				stop,
			)
			if err != nil {
				yield(&nodesevententities.RunCompletedEvent{
					RunResult: &workflowentities.NodeRunResult{
						Status:      models.WorkflowNodeExecutionStatus_FAILED,
						Error:       err.Error(),
						Inputs:      node_inputs,
						ProcessData: process_data,
					},
				})
				return
			}

			usage = modelruntimeentities.NewLLMUsage()
			for event := range llm_generator {
				if real_event, ok := any(event).(*nodesevententities.RunStreamChunkEvent); ok {
					if !yield(real_event) {
						return
					}
				} else if real_event, ok := any(event).(*nodesevententities.ModelInvokeCompletedEvent); ok {
					result_text = real_event.Text
					usage = real_event.Usage
					finish_reason = real_event.FinishReason
					// deduct quota
					DeductLLMQuota(n.GetTenantID(), model_instance, usage)
					break
				}
			}
		}()

		bindata, _ := json.Marshal(usage)
		outputs := map[string]any{"text": result_text, "usage": string(bindata), "finish_reason": finish_reason}

		yield(&nodesevententities.RunCompletedEvent{
			RunResult: &workflowentities.NodeRunResult{
				Status:      models.WorkflowNodeExecutionStatus_SUCCEEDED,
				Inputs:      node_inputs,
				ProcessData: process_data,
				Outputs:     outputs,
				Metadata: map[workflowenumtypes.NodeRunMetadataKey]any{
					workflowenumtypes.NodeRunMetadataKey_TOTAL_TOKENS: usage.TotalTokens,
					workflowenumtypes.NodeRunMetadataKey_TOTAL_PRICE:  usage.TotalPrice,
					workflowenumtypes.NodeRunMetadataKey_CURRENCY:     usage.Currency,
				},
				LLmUsage: usage,
			},
		})
	}
	// return &workflowentities.NodeRunResult{
	// 	Status: models.WorkflowNodeExecutionStatus_SUCCEEDED,
	// }, nil, nil
}

func (n *LLMNode) ExtractVariableSelectorToVariableMapping(graph_config map[string]any, node_id string, node_data *llmnodesentities.LLMNodeData) map[string][]string {
	prompt_template := node_data.PromptTemplate
	mlog.Debugf("------node_data=%#v", node_data)
	mlog.Debugf("------prompt_template=%#v", prompt_template)
	variable_selectors := []*workflowentities.VariableSelector{}
	switch real_prompt_template := prompt_template.(type) {
	case []*llmnodesentities.LLMNodeChatModelMessage:
		for _, prompt := range real_prompt_template {
			variable_template_parser := variabletemplateparserutils.NewVariableTemplateParser(prompt.Text)
			variable_selectors = append(variable_selectors, variable_template_parser.ExtractVariableSelectors()...)
		}
	case *llmnodesentities.LLMNodeCompletionModelPromptTemplate:
		variable_template_parser := variabletemplateparserutils.NewVariableTemplateParser(real_prompt_template.Text)
		variable_selectors = variable_template_parser.ExtractVariableSelectors()
	default:
		panic(llmnodesexceptions.NewInvalidVariableTypeError(fmt.Sprintf("Invalid prompt template type: %#v", prompt_template)))
	}

	variable_mapping := map[string][]string{}
	for _, variable_selector := range variable_selectors {
		variable_mapping[variable_selector.Variable] = variable_selector.ValueSelector
	}
	mem := node_data.Memory
	if mem != nil && mem.QueryPromptTemplate != "" {
		query_variable_selectors := variabletemplateparserutils.NewVariableTemplateParser(mem.QueryPromptTemplate).ExtractVariableSelectors()
		for _, variable_selector := range query_variable_selectors {
			variable_mapping[variable_selector.Variable] = variable_selector.ValueSelector
		}
	}
	if node_data.Context.Enabled {
		variable_mapping["#context#"] = node_data.Context.VariableSelector
	}
	if node_data.Vision.Enabled {
		variable_mapping["#files#"] = []string{"sys", string(workflowenumtypes.SystemVariableKey_FILES)}
	}
	if node_data.Memory != nil {
		variable_mapping["#sys.query#"] = []string{"sys", string(workflowenumtypes.SystemVariableKey_QUERY)}
	}
	// if node_data.PromptConfig{
	//     enable_jinja = False

	//     if isinstance(prompt_template, list){
	//         for prompt in prompt_template{
	//             if prompt.edition_type == "jinja2"{
	//                 enable_jinja = True
	//                 break
	//             }
	//         }
	//     }else{
	//         if prompt_template.edition_type == "jinja2"{
	//             enable_jinja = True
	//         }
	//     }
	//     if enable_jinja{
	//         for variable_selector in node_data.prompt_config.jinja2_variables or []{
	//             variable_mapping[variable_selector.variable] = variable_selector.value_selector
	//         }
	//     }
	// }
	new_variable_mapping := map[string][]string{}
	for key, value := range variable_mapping {
		new_variable_mapping[node_id+"."+key] = value
	}

	return new_variable_mapping
}

func (n *LLMNode) GetDefaultConfig(filters map[string]any) map[string]any {
	return map[string]any{
		"type": "llm",
		"config": map[string]any{
			"prompt_templates": map[string]any{
				"chat_model": map[string]any{
					"prompts": []map[string]any{
						{"role": "system", "text": "You are a helpful AI assistant.", "edition_type": "basic"},
					},
				},
				"completion_model": map[string]any{
					"conversation_histories_role": map[string]any{"user_prefix": "Human", "assistant_prefix": "Assistant"},
					"prompt": map[string]any{
						"text":         "Here are the chat histories between human and assistant, inside <histories></histories> XML tags.\n\n<histories>\n{{#histories#}}\n</histories>\n\n\nHuman: {{#sys.query#}}\n\nAssistant:",
						"edition_type": "basic",
					},
					"stop": []string{"Human:"},
				},
			},
		},
	}
}

func _fetch_inputs(n *LLMNode, node_data *llmnodesentities.LLMNodeData) (map[string]any, error) {
	inputs := map[string]any{}
	prompt_template := node_data.PromptTemplate

	variable_selectors := []*workflowentities.VariableSelector{}
	switch real_prompt_template := prompt_template.(type) {
	case []*llmnodesentities.LLMNodeChatModelMessage:
		for _, prompt := range real_prompt_template {
			variable_template_parser := variabletemplateparserutils.NewVariableTemplateParser(prompt.Text)
			variable_selectors = append(variable_selectors, variable_template_parser.ExtractVariableSelectors()...)
		}
	case *llmnodesentities.LLMNodeCompletionModelPromptTemplate:
		variable_template_parser := variabletemplateparserutils.NewVariableTemplateParser(real_prompt_template.Text)
		variable_selectors = variable_template_parser.ExtractVariableSelectors()
	default:
		return nil, llmnodesexceptions.NewInvalidVariableTypeError(fmt.Sprintf("Invalid prompt template type: %#v", prompt_template))
	}

	for _, variable_selector := range variable_selectors {
		variable := n.GetGraphRuntimeState().VariablePool.Get(variable_selector.ValueSelector)
		if variable == nil {
			return nil, llmnodesexceptions.NewVariableNotFoundError(fmt.Sprintf("Variable %s not found", variable_selector.Variable))
		}
		if variable.ValueType() == variableenumtypes.Variable_NONE {
			inputs[variable_selector.Variable] = ""
		}

		inputs[variable_selector.Variable] = variable.ToObject()
	}
	mem := node_data.Memory
	if mem != nil && mem.QueryPromptTemplate != "" {
		query_variable_selectors := variabletemplateparserutils.NewVariableTemplateParser(mem.QueryPromptTemplate).ExtractVariableSelectors()
		for _, variable_selector := range query_variable_selectors {
			variable := n.GetGraphRuntimeState().VariablePool.Get(variable_selector.ValueSelector)
			if variable == nil {
				return nil, llmnodesexceptions.NewVariableNotFoundError(fmt.Sprintf("Variable %s not found", variable_selector.Variable))
			}
			if variable.ValueType() == variableenumtypes.Variable_NONE {
				continue
			}
			inputs[variable_selector.Variable] = variable.ToObject()
		}
	}
	return inputs, nil
}

func _handle_invoke_result[T1 *modelruntimeentities.LLMResult | iter.Seq[*modelruntimeentities.LLMResultChunk]](node_id string, invoke_result T1) iter.Seq[nodesevententities.NodeEventer] {
	return func(yield func(nodesevententities.NodeEventer) bool) {
		if _, ok := any(invoke_result).(*modelruntimeentities.LLMResult); ok {
			return
		} else if real_invoke_result, ok := any(invoke_result).(iter.Seq[*modelruntimeentities.LLMResultChunk]); ok {
			model := ""
			prompt_messages := []modelruntimeentities.PromptMessager{}
			full_text := ""
			var usage *modelruntimeentities.LLMUsage
			finish_reason := ""
			for result := range real_invoke_result {
				text := result.Delta.Message.Content
				full_text += text

				if !yield(&nodesevententities.RunStreamChunkEvent{ChunkContent: text, FromVariableSelector: []string{node_id, "text"}}) {
					return
				}

				if model == "" {
					model = result.Model
				}
				if len(prompt_messages) == 0 {
					prompt_messages = result.PromptMessages
				}
				if usage == nil && result.Delta.Usage != nil {
					usage = result.Delta.Usage
				}
				if finish_reason == "" && result.Delta.FinishReason != "" {
					finish_reason = result.Delta.FinishReason
				}
			}
			if usage == nil {
				usage = modelruntimeentities.NewLLMUsage()
			}
			if !yield(&nodesevententities.ModelInvokeCompletedEvent{Text: full_text, Usage: usage, FinishReason: finish_reason}) {
				return
			}
		}
	}
}

func InvokeLLM(
	user_id string,
	node_id string,
	node_data_model *llmnodesentities.ModelConfig,
	model_instance *modelmanager.ModelInstance,
	prompt_messages []modelruntimeentities.PromptMessager,
	stop []string,
) iter.Seq[nodesevententities.NodeEventer] {
	invoke_result := model_instance.InvokeLLMStream(
		prompt_messages,
		node_data_model.CompletionParams,
		nil,
		stop,
		user_id,
		nil,
	)

	return _handle_invoke_result(node_id, invoke_result)
}

func _convert_to_original_retriever_resource(context_dict map[string]any) *ragentities.RetrievalSourceMetadata {
	if _, ok := context_dict["metadata"]; ok {
		if metadata, ok := context_dict["metadata"].(map[string]any); ok {
			if _, ok := metadata["_source"]; ok {
				if _, ok := metadata["_source"].(string); ok && metadata["_source"].(string) == "knowledge" {
					rsm := &ragentities.RetrievalSourceMetadata{}
					bindata, _ := json.Marshal(metadata)
					err := json.Unmarshal(bindata, rsm)
					if err != nil {
						mlog.Error("json unmarshal failed:%v", err)
						return nil
					}

					rsm.DataSourceType = metadata["document_data_source_type"].(string)
					rsm.HitCount = metadata["segment_hit_count"].(int)
					rsm.WordCount = metadata["segment_word_count"].(int)
					rsm.IndexNodeHash = metadata["segment_index_node_hash"].(string)
					return rsm

				}
			}
		}
	}
	return nil
}

func _fetch_context(n *LLMNode, node_data *llmnodesentities.LLMNodeData) iter.Seq2[*nodesevententities.RunRetrieverResourceEvent, error] {
	return func(yield func(*nodesevententities.RunRetrieverResourceEvent, error) bool) {
		if !node_data.Context.Enabled {
			return
		}
		if len(node_data.Context.VariableSelector) == 0 {
			return
		}
		context_value_variable := n.GetGraphRuntimeState().VariablePool.Get(node_data.Context.VariableSelector)
		if context_value_variable != nil {
			switch real_segment := any(context_value_variable).(type) {
			case *variables.StringVariable:
				if !yield(&nodesevententities.RunRetrieverResourceEvent{Context: real_segment.Value}, nil) {
					return
				}
			case *variables.ArrayStringVariable:
				context_str := ""
				original_retriever_resource := []*ragentities.RetrievalSourceMetadata{}
				for _, item := range real_segment.Value {
					context_str += item + "\n"
				}
				if !yield(&nodesevententities.RunRetrieverResourceEvent{RetrieverResources: original_retriever_resource, Context: strings.TrimSpace(context_str)}, nil) {
					return
				}
			case *variables.ArrayObjectVariable:
				context_str := ""
				original_retriever_resource := []*ragentities.RetrievalSourceMetadata{}
				for _, item := range real_segment.Value {
					if _, ok := item["content"]; !ok {
						mlog.Errorf("Invalid context structure: %#v", item)
						yield(nil, llmnodesexceptions.NewInvalidContextStructureError(fmt.Sprintf("Invalid context structure: %#v", item)))
						return
					}
					if _, ok := item["content"].(string); !ok {
						mlog.Errorf("Invalid context structure: %#v", item)
						yield(nil, llmnodesexceptions.NewInvalidContextStructureError(fmt.Sprintf("Invalid context structure: %#v", item)))
						return
					}
					context_str += item["content"].(string) + "\n"

					retriever_resource := _convert_to_original_retriever_resource(item)
					if retriever_resource != nil {
						original_retriever_resource = append(original_retriever_resource, retriever_resource)
					}

				}
				if !yield(&nodesevententities.RunRetrieverResourceEvent{RetrieverResources: original_retriever_resource, Context: strings.TrimSpace(context_str)}, nil) {
					return
				}
			default:
				mlog.Errorf("unsuported context_value_variable Segmenter: %#v", real_segment)
				yield(nil, llmnodesexceptions.NewInvalidContextStructureError(fmt.Sprintf("unsuported context_value_variable Segmenter: %#v", real_segment)))
				return
			}
		}
	}
}

func FetchModelConfig(tenant_id string, node_data_model *llmnodesentities.ModelConfig) (*modelmanager.ModelInstance, *appconfigentities.ModelConfigWithCredentialsEntity) {
	model_name := node_data_model.Name
	provider_name := node_data_model.Provider

	model_instance := (&modelmanager.ModelManager{}).GetModelInstance(
		tenant_id, provider_name, modelruntimeenumtypes.Model_LLM, model_name,
	)

	provider_model_bundle := model_instance.ProviderModelBundle
	model_type_instance := any(model_instance.ModelTypeInstance).(modelruntimeentities.LargeLanguageModeler)

	model_credentials := model_instance.Credentials

	// check model
	provider_model := (&providermanager.ProviderConfigurationManager{}).GetProviderModel(
		provider_model_bundle.Configuration, modelruntimeenumtypes.Model_LLM, model_name, false,
	)
	if provider_model == nil {
		panic(llmnodesexceptions.NewModelNotExistError(fmt.Sprintf("Model %s not exist.", model_name)))
	}

	if provider_model.Status == modelenumtypes.ModelStatus_NO_CONFIGURE {
		panic(exceptions.NewProviderTokenNotInitError(fmt.Sprintf("Model %s credentials is not initialized.", model_name)))
	} else if provider_model.Status == modelenumtypes.ModelStatus_NO_PERMISSION {
		panic(exceptions.NewModelCurrentlyNotSupportError(fmt.Sprintf("Gofy Hosted OpenAI %s currently not support.", model_name)))
	} else if provider_model.Status == modelenumtypes.ModelStatus_QUOTA_EXCEEDED {
		panic(exceptions.NewQuotaExceededError(fmt.Sprintf("Model provider %s quota exceeded.", provider_name)))
	}
	// model config
	completion_params := node_data_model.CompletionParams
	stop := []string{}
	if _, ok := completion_params["stop"]; ok {
		if _, ok := completion_params["stop"].([]string); ok {
			stop = completion_params["stop"].([]string)
			delete(completion_params, "stop")
		} else if _, ok := completion_params["stop"].([]any); ok {
			is_all_string := true
			new_stop := []string{}
			for _, v := range completion_params["stop"].([]any) {
				if _, ok := v.(string); !ok {
					is_all_string = false
					break
				} else {
					new_stop = append(new_stop, v.(string))
				}
			}
			if is_all_string {
				stop = new_stop
				delete(completion_params, "stop")
			}
		}

	}
	// get model mode
	model_mode := node_data_model.Mode
	if string(model_mode) == "" {
		panic(llmnodesexceptions.NewLLMModeRequiredError("LLM mode is required."))
	}
	model_schema := model_type_instance.GetModelSchema(model_type_instance, model_name, model_credentials)

	if model_schema == nil {
		panic(llmnodesexceptions.NewModelNotExistError(fmt.Sprintf("Model %s not exist.", model_name)))
	}
	return model_instance, &appconfigentities.ModelConfigWithCredentialsEntity{
		Provider:            provider_name,
		Model:               model_name,
		ModelSchema:         model_schema,
		Mode:                model_mode,
		ProviderModelBundle: provider_model_bundle,
		Credentials:         model_credentials,
		Parameters:          completion_params,
		Stop:                stop,
	}
}

func FetchMemory(
	graph_runtime_state *graphengineentities.GraphRuntimeState, app_id string, node_data_memory *promptentities.MemoryConfig, model_instance *modelmanager.ModelInstance,
) *memory.TokenBufferMemory {
	if node_data_memory == nil {
		return nil
	}
	// get conversation id
	conversation_id_variable := graph_runtime_state.VariablePool.Get([]string{"sys", string(workflowenumtypes.SystemVariableKey_CONVERSATION_ID)})
	if conversation_id_variable == nil {
		return nil
	}
	if conversation_id_variable.ValueType() != variableenumtypes.Variable_STRING {
		mlog.Warningf("conversation_id_variable must be string variable")
		return nil
	}

	conversation_id := conversation_id_variable.(*variables.StringVariable).Value

	// get conversation
	conversation := new(models.Conversation)
	err := dbengine.Instance().DB.Debug().Model(&models.Conversation{}).Where("app_id = ? and id = ? ", app_id, conversation_id).First(conversation).Error
	if err != nil {
		mlog.Warningf("get Conversation failded:%v", err)
		conversation = nil
	}

	if conversation == nil {
		return nil
	}

	return memory.NewTokenBufferMemory(conversation, model_instance)
}

func _handle_list_messages(
	messages []*llmnodesentities.LLMNodeChatModelMessage,
	context string,
	variable_pool *workflowentities.VariablePool,
	vision_detail_config modelruntimeentities.ImagePromptMessageContentDETAIL,
) []modelruntimeentities.PromptMessager {
	prompt_messages := []modelruntimeentities.PromptMessager{}
	for _, message := range messages {
		template := message.Text
		// Get segment group from basic message
		if context != "" {
			template = strings.ReplaceAll(message.Text, "{#context#}", context)
		}
		segment_group := variable_pool.ConvertTemplate(template)

		// Create message with text from all segments
		plain_text := segment_group.Text()
		if plain_text != "" {
			prompt_message := _combine_message_content_with_role(
				[]modelruntimeentities.PromptMessageContenter{modelruntimeentities.NewTextPromptMessageContent(plain_text)}, message.Role,
			)

			prompt_messages = append(prompt_messages, prompt_message)
		}
	}
	return prompt_messages
}

func FetchPromptMessages[T []*llmnodesentities.LLMNodeChatModelMessage | *llmnodesentities.LLMNodeCompletionModelPromptTemplate](
	sys_query string,
	// sys_files: Sequence["File"],
	context string,
	mem *memory.TokenBufferMemory,
	model_config *appconfigentities.ModelConfigWithCredentialsEntity,
	prompt_template T,
	memory_config *promptentities.MemoryConfig,
	vision_enabled bool,
	vision_detail modelruntimeentities.ImagePromptMessageContentDETAIL,
	variable_pool *workflowentities.VariablePool,
) ([]modelruntimeentities.PromptMessager, []string) {
	// FIXME: fix the type error cause prompt_messages is type quick a few times
	prompt_messages := []modelruntimeentities.PromptMessager{}

	if real_prompt_template, ok := any(prompt_template).([]*llmnodesentities.LLMNodeChatModelMessage); ok {
		// For chat model
		tmp_prompt_messages := _handle_list_messages(
			real_prompt_template,
			context,
			variable_pool,
			vision_detail,
		)
		if len(tmp_prompt_messages) > 0 {
			prompt_messages = append(prompt_messages, tmp_prompt_messages...)
		}

		// Get memory messages for chat mode
		memory_messages := _handle_memory_chat_mode(
			mem,
			memory_config,
			model_config,
		)

		// Extend prompt_messages with memory messages
		prompt_messages = append(prompt_messages, memory_messages...)

		// Add current query to the prompt messages
		if sys_query != "" {
			message := &llmnodesentities.LLMNodeChatModelMessage{
				ChatModelMessage: &promptentities.ChatModelMessage{
					Text:        sys_query,
					Role:        modelruntimeentities.PromptMessageRole_USER,
					EditionType: promptentities.Edition_Basic,
				},
			}
			tmp_prompt_messages := _handle_list_messages(
				[]*llmnodesentities.LLMNodeChatModelMessage{message},
				"",
				variable_pool,
				vision_detail,
			)
			if len(tmp_prompt_messages) > 0 {
				prompt_messages = append(prompt_messages, tmp_prompt_messages...)
			}
		}
	} else if real_prompt_template, ok := any(prompt_template).(*llmnodesentities.LLMNodeCompletionModelPromptTemplate); ok {
		// For completion model
		tmp_prompt_messages := _handle_completion_template(
			real_prompt_template,
			context,
			variable_pool,
		)

		if len(tmp_prompt_messages) > 0 {
			prompt_messages = append(prompt_messages, tmp_prompt_messages...)
		}

		// Get memory text for completion model
		memory_text := _handle_memory_completion_mode(
			mem,
			memory_config,
			model_config,
		)

		// Insert histories into the prompt
		prompt_content := prompt_messages[0].GetContent()
		// For issue #11247 - Check if prompt content is a string or a list
		// prompt_content_type = type(prompt_content)
		if real_prompt_content, ok := prompt_content.(string); ok {
			if strings.Contains(real_prompt_content, "#histories#") {
				real_prompt_content = strings.ReplaceAll(real_prompt_content, "#histories#", memory_text)
			} else {
				real_prompt_content = memory_text + "\n" + real_prompt_content
			}
			prompt_messages[0].SetContent(real_prompt_content)
		} else if real_prompt_content, ok := prompt_content.([]modelruntimeentities.PromptMessageContenter); ok {
			for idx, content_item := range real_prompt_content {
				if real_content_item, ok := any(content_item).(*modelruntimeentities.TextPromptMessageContent); ok {
					if strings.Contains(real_content_item.Data(), "#histories#") {
						real_content_item.SetData(strings.ReplaceAll(real_content_item.Data(), "#histories#", memory_text))
					} else {
						real_content_item.SetData(memory_text + "\n" + real_content_item.Data())
					}
					real_prompt_content[idx] = real_content_item
				}
			}
			prompt_messages[0].SetContent(real_prompt_content)
		} else {
			mlog.Errorf("Invalid prompt content=%#v type", prompt_content)
			panic(exceptions.NewValueError("Invalid prompt content type"))
		}
		// Add current query to the prompt message
		if sys_query != "" {
			if real_prompt_content, ok := prompt_content.(string); ok {
				real_prompt_content = strings.ReplaceAll(real_prompt_content, "#sys.query#", sys_query)
				prompt_messages[0].SetContent(real_prompt_content)
			} else if real_prompt_content, ok := prompt_content.([]modelruntimeentities.PromptMessageContenter); ok {
				for idx, content_item := range real_prompt_content {
					if real_content_item, ok := any(content_item).(*modelruntimeentities.TextPromptMessageContent); ok {
						real_content_item.SetData(sys_query + "\n" + real_content_item.Data())
						real_prompt_content[idx] = real_content_item
					}
				}
				prompt_messages[0].SetContent(real_prompt_content)
			} else {
				mlog.Errorf("Invalid prompt content=%#v type", prompt_content)
				panic(exceptions.NewValueError("Invalid prompt content type"))
			}
		}
	}

	// Remove empty messages and filter unsupported content
	filtered_prompt_messages := []modelruntimeentities.PromptMessager{}
	for idx, prompt_message := range prompt_messages {
		if real_prompt_message, ok := prompt_message.(*modelruntimeentities.UserPromptMessage[[]modelruntimeentities.PromptMessageContenter]); ok {
			prompt_message_content := []modelruntimeentities.PromptMessageContenter{}
			for _, content_item := range real_prompt_message.Content {
				// Skip content if features are not defined
				if len(model_config.ModelSchema.Features) == 0 {
					if content_item.Type() != modelruntimeentities.PromptMessageContent_TEXT {
						continue
					}
					prompt_message_content = append(prompt_message_content, content_item)
					continue
				}
				// Skip content if corresponding feature is not supported

				if ((content_item.Type() == modelruntimeentities.PromptMessageContent_IMAGE &&
					!slices.Contains(model_config.ModelSchema.Features, modelruntimeenumtypes.ModelFeature_VISION)) ||
					(content_item.Type() == modelruntimeentities.PromptMessageContent_DOCUMENT &&
						!slices.Contains(model_config.ModelSchema.Features, modelruntimeenumtypes.ModelFeature_DOCUMENT))) ||
					(content_item.Type() == modelruntimeentities.PromptMessageContent_VIDEO &&
						!slices.Contains(model_config.ModelSchema.Features, modelruntimeenumtypes.ModelFeature_VIDEO)) ||
					(content_item.Type() == modelruntimeentities.PromptMessageContent_AUDIO &&
						!slices.Contains(model_config.ModelSchema.Features, modelruntimeenumtypes.ModelFeature_AUDIO)) {
					continue
				}
				prompt_message_content = append(prompt_message_content, content_item)
			}
			real_prompt_message.Content = prompt_message_content
			prompt_messages[idx] = real_prompt_message
		}

		filtered_prompt_messages = append(filtered_prompt_messages, prompt_message)
	}
	if len(filtered_prompt_messages) == 0 {
		panic(llmnodesexceptions.NewNoPromptFoundError("No prompt found in the LLM configuration. Please ensure a prompt is properly configured before proceeding."))
	}
	stop := model_config.Stop
	return filtered_prompt_messages, stop

}

func New() *LLMNode {
	return &LLMNode{
		BaseNode: &base.BaseNode[*llmnodesentities.LLMNodeData]{},
	}
}

// func _render_jinja2_message(
//     template string,
//     jinjia2_variables []*workflowentities.VariableSelector,
//     variable_pool *workflowentities.VariablePool,
// ) string{
//     if template == ""{
//         return ""
//     }
//     jinjia2_inputs := map[string]any{}
//     for _, jinja2_variable := range jinjia2_variables{
//         variable := variable_pool.Get(jinja2_variable.ValueSelector)
//         if variable != nil {
//             jinjia2_inputs[jinja2_variable.Variable] = variable.ToObject()
//         } else {
//             jinjia2_inputs[jinja2_variable.Variable] = ""
//         }
//     }
//     code_execute_resp = codeexecutor.ExecuteWorkflowCodeTemplate(
//         language=CodeLanguage.JINJA2,
//         code=template,
//         inputs=jinjia2_inputs,
//     )
//     result_text = code_execute_resp["result"]
//     return result_text

// }
func _calculate_rest_token(
	prompt_messages []modelruntimeentities.PromptMessager, model_config *appconfigentities.ModelConfigWithCredentialsEntity,
) int {
	rest_tokens := 2000

	if _, ok := model_config.ModelSchema.ModelProperties[modelruntimeenumtypes.ModelPropertyKey_CONTEXT_SIZE]; ok {
		if _, ok := model_config.ModelSchema.ModelProperties[modelruntimeenumtypes.ModelPropertyKey_CONTEXT_SIZE].(int); !ok {
			mlog.Errorf("ModelProperties[%s]=%#v must be int", modelruntimeenumtypes.ModelPropertyKey_CONTEXT_SIZE, model_config.ModelSchema.ModelProperties[modelruntimeenumtypes.ModelPropertyKey_CONTEXT_SIZE])
			panic(exceptions.NewValueError(fmt.Sprintf("modelProperties[%s]=%#v must be int", modelruntimeenumtypes.ModelPropertyKey_CONTEXT_SIZE, model_config.ModelSchema.ModelProperties[modelruntimeenumtypes.ModelPropertyKey_CONTEXT_SIZE])))
		}
		model_context_tokens := model_config.ModelSchema.ModelProperties[modelruntimeenumtypes.ModelPropertyKey_CONTEXT_SIZE].(int)
		model_instance := modelmanager.NewModelInstance(
			model_config.ProviderModelBundle, model_config.Model,
		)

		curr_message_tokens := model_instance.GetLLMNumTokens(prompt_messages, nil)
		max_tokens := 0
		for _, parameter_rule := range model_config.ModelSchema.ParameterRules {
			if parameter_rule.Name == "max_tokens" || (parameter_rule.UseTemplate != "" && parameter_rule.UseTemplate == "max_tokens") {
				if _, ok := model_config.Parameters[parameter_rule.Name]; ok {
					val, err := cast.ToE[int](model_config.Parameters[parameter_rule.Name])
					if err != nil {
						mlog.Errorf("cast(%#v) to int failed:%v", model_config.Parameters[parameter_rule.Name], err)
						if _, ok := model_config.Parameters[parameter_rule.UseTemplate]; ok {
							val, err := cast.ToE[int](model_config.Parameters[parameter_rule.UseTemplate])
							if err != nil {
								mlog.Errorf("cast(%#v) to int failed:%v", model_config.Parameters[parameter_rule.UseTemplate], err)
							} else {
								max_tokens = val
							}
						}
					} else {
						max_tokens = val
					}
				}
			}
		}

		rest_tokens = model_context_tokens - max_tokens - curr_message_tokens
		rest_tokens = max(rest_tokens, 0)
	}
	return rest_tokens

}

func _handle_memory_chat_mode(
	mem *memory.TokenBufferMemory,
	memory_config *promptentities.MemoryConfig,
	model_config *appconfigentities.ModelConfigWithCredentialsEntity,
) []modelruntimeentities.PromptMessager {
	memory_messages := []modelruntimeentities.PromptMessager{}
	// Get messages from memory for chat model
	if mem != nil && memory_config != nil {
		rest_tokens := _calculate_rest_token(nil, model_config)
		if memory_config.Window.Enabled {
			memory_messages = mem.GetHistoryPromptMessages(rest_tokens, memory_config.Window.Size)
		} else {
			memory_messages = mem.GetHistoryPromptMessages(rest_tokens, 0)
		}
	}
	return memory_messages

}

func _handle_memory_completion_mode(
	mem *memory.TokenBufferMemory,
	memory_config *promptentities.MemoryConfig,
	model_config *appconfigentities.ModelConfigWithCredentialsEntity,
) string {
	memory_text := ""
	// Get history text from memory for completion model
	if mem != nil && memory_config != nil {
		rest_tokens := _calculate_rest_token(nil, model_config)

		if memory_config.RolePrefix == nil {
			panic(llmnodesexceptions.NewMemoryRolePrefixRequiredError("Memory role prefix is required for completion model."))
		}
		if memory_config.Window.Enabled {
			memory_text = mem.GetHistoryPromptText(memory_config.RolePrefix.User, memory_config.RolePrefix.Assistant, rest_tokens, memory_config.Window.Size)
		} else {
			memory_text = mem.GetHistoryPromptText(memory_config.RolePrefix.User, memory_config.RolePrefix.Assistant, rest_tokens, 0)
		}
	}
	return memory_text

}
func _combine_message_content_with_role(content []modelruntimeentities.PromptMessageContenter, role modelruntimeentities.PromptMessageRole) modelruntimeentities.PromptMessager {
	switch role {
	case modelruntimeentities.PromptMessageRole_USER:
		return modelruntimeentities.NewUserPromptMessage(content, "")
	case modelruntimeentities.PromptMessageRole_ASSISTANT:
		if len(content) != 1 {
			mlog.Errorf("ASSISTANT PromptMessage only surport string content")
			panic(exceptions.NewNotImplementedError("ASSISTANT PromptMessage only surport string content"))
		}
		if real_content, ok := any(content[0]).(*modelruntimeentities.TextPromptMessageContent); ok {
			return modelruntimeentities.NewAssistantPromptMessage(real_content.Data(), "", nil)
		} else {
			panic(exceptions.NewNotImplementedError("ASSISTANT PromptMessage only surport string content"))
		}

	case modelruntimeentities.PromptMessageRole_SYSTEM:
		if len(content) != 1 {
			mlog.Errorf("SYSTEM PromptMessage only surport string content")
			panic(exceptions.NewNotImplementedError("SYSTEM PromptMessage only surport string content"))
		}
		if real_content, ok := any(content[0]).(*modelruntimeentities.TextPromptMessageContent); ok {
			return modelruntimeentities.NewSystemPromptMessage(real_content.Data(), "")
		} else {
			panic(exceptions.NewNotImplementedError("SYSTEM PromptMessage only surport string content"))
		}
	}
	panic(exceptions.NewNotImplementedError(fmt.Sprintf("Role %s is not supported", role)))
}
func _handle_completion_template(
	template *llmnodesentities.LLMNodeCompletionModelPromptTemplate,
	context string,
	variable_pool *workflowentities.VariablePool,
) []modelruntimeentities.PromptMessager {
	/*Handle completion template processing outside of LLMNode class.

	  Args:
	      template: The completion model prompt template
	      context: Optional context string
	      jinja2_variables: Variables for jinja2 template rendering
	      variable_pool: Variable pool for template conversion

	  Returns:
	      Sequence of prompt messages
	*/
	template_text := ""
	prompt_messages := []modelruntimeentities.PromptMessager{}
	if context != "" {
		template_text = strings.ReplaceAll(template.Text, "{#context#}", context)
	} else {
		template_text = template.Text
	}
	result_text := variable_pool.ConvertTemplate(template_text).Text()
	prompt_message := _combine_message_content_with_role(
		[]modelruntimeentities.PromptMessageContenter{modelruntimeentities.NewTextPromptMessageContent(result_text)}, modelruntimeentities.PromptMessageRole_USER,
	)

	prompt_messages = append(prompt_messages, prompt_message)
	return prompt_messages
}
