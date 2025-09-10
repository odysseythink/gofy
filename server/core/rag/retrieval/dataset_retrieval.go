package retrieval

import (
	"encoding/json"
	"strings"

	"mlib.com/gofy/server/core/exceptions"
	"mlib.com/gofy/server/core/manageres"
	modelmanager "mlib.com/gofy/server/core/manageres/model_manager"
	providermanager "mlib.com/gofy/server/core/manageres/provider_manager"
	"mlib.com/gofy/server/core/prompt"
	dbengine "mlib.com/gofy/server/db_engine"
	ragentities "mlib.com/gofy/server/entities/rag"
	appconfigentities "mlib.com/gofy/server/entities/app/config"
	modelruntimeentities "mlib.com/gofy/server/entities/model_runtime"
	promptentities "mlib.com/gofy/server/entities/prompt"
	promptenumttypes "mlib.com/gofy/server/enum_types/prompt"
	ragretrievalenumtypes "mlib.com/gofy/server/enum_types/rag/retrieval"
	"mlib.com/gofy/server/models"
	"mlib.com/gofy/server/utils/mapstruct"
)

var (
	default_retrieval_model = map[string]any{
		"search_method":           ragretrievalenumtypes.RetrievalMethod_SEMANTIC_SEARCH,
		"reranking_enable":        false,
		"reranking_model":         map[string]any{"reranking_provider_name": "", "reranking_model_name": ""},
		"top_k":                   2,
		"score_threshold_enabled": false,
	}
)

type DatasetRetrieval struct {
}


func _fetch_model_config(
        tenant_id  string, model  *appconfigentities.ModelConfig,
    ) (*modelmanager.ModelInstance, *appconfigentities.ModelConfigWithCredentialsEntity){
        if model == nil {
            panic(exceptions.NewValueError("single_retrieval_config is required"))
		}
        model_name := model.Name
        provider_name := model.Provider

				model_instance := manageres.Instance.Model.GetModelInstance(
			tenant_id,
			provider_name,
			modelruntimeenumtypes.Model_LLM,
			model_name,
		)

        provider_model_bundle := model_instance.ProviderModelBundle
        model_type_instance := model_instance.ModelTypeInstance
        // model_type_instance = cast(LargeLanguageModel, model_type_instance)

        model_credentials := model_instance.Credentials

        // check model
	provider_model := (&providermanager.ProviderConfigurationManager{}).GetProviderModel(
		provider_model_bundle.Configuration,
		modelruntimeenumtypes.Model_LLM,
		model_name,
		false,
	)
        if provider_model == nil{
            panic(exceptions.NewValueError(fmt.Sprintf("Model %s not exist.", model_name)))
		}
        if provider_model.status == ModelStatus.NO_CONFIGURE{
            panic(exceptions.NewValueError(fmt.Sprintf("Model {%s} credentials is not initialized.", model_name)))
        } else if provider_model.status == ModelStatus.NO_PERMISSION{
            panic(exceptions.NewValueError(fmt.Sprintf("Dify Hosted OpenAI {%s} currently not support.", model_name)))
        } else if provider_model.status == ModelStatus.QUOTA_EXCEEDED{
            panic(exceptions.NewValueError(fmt.Sprintf("Model provider {%s} quota exceeded.", provider_name)))
		}
        // model config
        completion_params := model.CompletionParams
        stop := []string{}
        if  _, ok := completion_params["stop"]; ok {
			stop = mapstruct.Get(completion_params, "stop", []string{})
            delete(completion_params,"stop")
		}
        // get model mode
        model_mode := model.Mode
        if string(model_mode) == ""{
            panic(exceptions.NewValueError("LLM mode is required."))
		}
        model_schema := model_type_instance.GetModelSchema(model_type_instance, model_name, model_credentials)

        if model_schema == nil{
            panic(exceptions.NewValueError(fmt.Sprintf("Model {%s} not exist.", model_name)))
		}
        return model_instance, &appconfigentities.ModelConfigWithCredentialsEntity{
            Provider:provider_name,
            Model:model_name,
            ModelSchema:model_schema,
            Mode:model_mode,
            ProviderModelBundle:provider_model_bundle,
            Credentials:model_credentials,
            Parameters:completion_params,
            Stop:stop,
        }
}
func _get_prompt_template(
        model_config *appconfigentities.ModelConfigWithCredentialsEntity, mode  string, metadata_fields []any, query  string,
    ){
        model_mode := promptenumttypes.ModelMode(mode)
        input_text := query

        var prompt_messages []modelruntimeentities.PromptMessager
        if model_mode == promptenumttypes.ModelMode_CHAT{
            prompt_template := []*promptentities.ChatModelMessage{}
            system_prompt_messages := &promptentities.ChatModelMessage{Role:modelruntimeentities.PromptMessageRole_SYSTEM, Text: METADATA_FILTER_SYSTEM_PROMPT}
            prompt_template=append(prompt_template,system_prompt_messages)
            user_prompt_message_1 := &promptentities.ChatModelMessage{Role:modelruntimeentities.PromptMessageRole_USER, Text: METADATA_FILTER_USER_PROMPT_1}
            prompt_template=append(prompt_template,user_prompt_message_1)
            assistant_prompt_message_1 := &promptentities.ChatModelMessage{
                Role:modelruntimeentities.PromptMessageRole_ASSISTANT, Text: METADATA_FILTER_ASSISTANT_PROMPT_1,
            }
            prompt_template=append(prompt_template,assistant_prompt_message_1)
            user_prompt_message_2 := &promptentities.ChatModelMessage{Role:modelruntimeentities.PromptMessageRole_USER, Text: METADATA_FILTER_USER_PROMPT_2}
            prompt_template=append(prompt_template,user_prompt_message_2)
            assistant_prompt_message_2 := &promptentities.ChatModelMessage{
                Role:modelruntimeentities.PromptMessageRole_ASSISTANT, Text: METADATA_FILTER_ASSISTANT_PROMPT_2,
            }
            prompt_template=append(prompt_template,assistant_prompt_message_2)
			bindata, _ := json.Marshal(metadata_fields)
            user_prompt_message_3 := &promptentities.ChatModelMessage{
                Role:modelruntimeentities.PromptMessageRole_USER,
                Text: strings.ReplaceAll(strings.ReplaceAll(METADATA_FILTER_USER_PROMPT_3,"{input_text}", input_text),"{metadata_fields}", string(bindata)),
            }
            prompt_template=append(prompt_template,user_prompt_message_3)
			prompt_messages = prompt.NewAdvancedPromptTransform[[]*promptentities.ChatModelMessage](false, modelruntimeentities.ImagePromptMessageContentDETAIL("")).GetPrompt(
            prompt_template,
            map[string]any{},
            query,
            nil,
            "",
            nil,
            nil,
            model_config,
        	)
        } else if model_mode == promptenumttypes.ModelMode_COMPLETION{
			bindata, _ := json.Marshal(metadata_fields)
            prompt_template := &promptentities.CompletionModelPromptTemplate{
                Text: strings.ReplaceAll(strings.ReplaceAll(METADATA_FILTER_COMPLETION_PROMPT,"{input_text}", input_text),"{metadata_fields}", string(bindata)),
            }
			prompt_messages = prompt.NewAdvancedPromptTransform[*promptentities.CompletionModelPromptTemplate](false, modelruntimeentities.ImagePromptMessageContentDETAIL("")).GetPrompt(
            prompt_template,
            map[string]any{},
            query,
            nil,
            "",
            nil,
            nil,
            model_config,
        	)
        }else{
            panic(exceptions.NewValueError(fmt.Sprintf("Model mode {%v} not support.", model_mode)))
		}
        stop = model_config.Stop

        return prompt_messages, stop

}
func get_metadata_filter_condition(
        dataset_ids []string,
        query  string,
        tenant_id  string,
        user_id  string,
        metadata_filtering_mode  string,
        metadata_model_config  *appconfigentities.ModelConfig,
        metadata_filtering_conditions *appconfigentities.MetadataFilteringCondition,
        inputs map[string]any,
    ) (map[string][]string, *ragentities.MetadataCondition){
        document_query := dbengine.Instance().DB.Model(&models.Document).Where("dataset_id in ? and indexing_status = ? and enabled = ? and archived = ?", dataset_ids,"completed", true, false)
        filters = []  // type: ignore
        var metadata_condition *MetadataCondition
        if metadata_filtering_mode == "disabled"{
            return None, None
        } else if metadata_filtering_mode == "automatic"{
            automatic_metadata_filters = self._automatic_metadata_filter_func(
                dataset_ids, query, tenant_id, user_id, metadata_model_config
            )
            if automatic_metadata_filters{
                conditions = []
                for sequence, filter in enumerate(automatic_metadata_filters){
                    self._process_metadata_filter_func(
                        sequence,
                        filter.get("condition"),  // type: ignore
                        filter.get("metadata_name"),  // type: ignore
                        filter.get("value"),
                        filters,  // type: ignore
                    )
                    conditions.append(
                        Condition(
                            name=filter.get("metadata_name"),  // type: ignore
                            comparison_operator=filter.get("condition"),  // type: ignore
                            value=filter.get("value"),
                        )
                    )
				}
                metadata_condition = MetadataCondition(
                    logical_operator=metadata_filtering_conditions.logical_operator
                    if metadata_filtering_conditions
                    else "or",  // type: ignore
                    conditions=conditions,
                )
			}
        } else if metadata_filtering_mode == "manual"{
            if metadata_filtering_conditions{
                conditions = []
                for sequence, condition in enumerate(metadata_filtering_conditions.conditions):  // type: ignore
                    metadata_name = condition.name
                    expected_value = condition.value
                    if expected_value is not None and condition.comparison_operator not in ("empty", "not empty"){
                        if isinstance(expected_value, str){
                            expected_value = self._replace_metadata_filter_value(expected_value, inputs)
						}
					}
                    conditions.append(
                        Condition(
                            name=metadata_name,
                            comparison_operator=condition.comparison_operator,
                            value=expected_value,
                        )
                    )
                    filters = self._process_metadata_filter_func(
                        sequence,
                        condition.comparison_operator,
                        metadata_name,
                        expected_value,
                        filters,
                    )
				}
                metadata_condition = MetadataCondition(
                    logical_operator=metadata_filtering_conditions.logical_operator,
                    conditions=conditions,
                )
			}
        }else{
            panic(exceptions.NewValueError("Invalid metadata filtering mode"))
		}
        if filters{
            if metadata_filtering_conditions and metadata_filtering_conditions.logical_operator == "and"{  // type: ignore
                document_query = document_query.where(and_(*filters))
            }else{
                document_query = document_query.where(or_(*filters))
			}
		}
        documents = document_query.all()
        // group by dataset_id
        metadata_filter_document_ids = defaultdict(list) if documents else None  // type: ignore
        for document in documents{
            metadata_filter_document_ids[document.dataset_id].append(document.id)  // type: ignore
        return metadata_filter_document_ids, metadata_condition

}
func _automatic_metadata_filter_func(
        dataset_ids: list, query  string, tenant_id  string, user_id  string, metadata_model_config  *appconfigentities.ModelConfig
    ) -> Optional[list[dict[str, Any]]]{
        // get all metadata field
        metadata_fields = db.session.query(DatasetMetadata).where(DatasetMetadata.dataset_id.in_(dataset_ids)).all()
        all_metadata_fields = [metadata_field.name for metadata_field in metadata_fields]
        // get metadata model config
        if metadata_model_config is None{
            panic(exceptions.NewValueError("metadata_model_config is required")
        // get metadata model instance
        // fetch model config
        model_instance, model_config = self._fetch_model_config(tenant_id, metadata_model_config)

        // fetch prompt messages
        prompt_messages, stop = self._get_prompt_template(
            model_config=model_config,
            mode=metadata_model_config.mode,
            metadata_fields=all_metadata_fields,
            query=query or "",
        )

        result_text = ""
        try{
            // handle invoke result
            invoke_result = cast(
                Generator[LLMResult, None, None],
                model_instance.invoke_llm(
                    prompt_messages=prompt_messages,
                    model_parameters=model_config.parameters,
                    stop=stop,
                    stream=True,
                    user=user_id,
                ),
            )

            // handle invoke result
            result_text, usage = self._handle_invoke_result(invoke_result=invoke_result)

            result_text_json = parse_and_check_json_markdown(result_text, [])
            automatic_metadata_filters = []
            if "metadata_map" in result_text_json{
                metadata_map = result_text_json["metadata_map"]
                for item in metadata_map{
                    if item.get("metadata_field_name") in all_metadata_fields{
                        automatic_metadata_filters.append(
                            {
                                "metadata_name": item.get("metadata_field_name"),
                                "value": item.get("metadata_field_value"),
                                "condition": item.get("comparison_operator"),
                            }
                        )
        except Exception as e{
            return None
        return automatic_metadata_filters		

}
func _handle_invoke_result(invoke_result: Generator) -> tuple[str, LLMUsage]{
        model = None
        prompt_messages: list[PromptMessage] = []
        full_text = ""
        usage = None
        for result in invoke_result{
            text = result.delta.message.content
            full_text += text

            if not model{
                model = result.model

            if not prompt_messages{
                prompt_messages = result.prompt_messages

            if not usage and result.delta.usage{
                usage = result.delta.usage

        if not usage{
            usage = LLMUsage.empty_usage()

        return full_text, usage
}