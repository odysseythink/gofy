package retrieval

import (
	"encoding/json"
	"fmt"
	"iter"
	"regexp"
	"slices"
	"strings"

	"github.com/odysseythink/mlog"
	"gorm.io/gorm/clause"
	"mlib.com/confy/cast"
	"mlib.com/gofy/server/core/exceptions"
	"mlib.com/gofy/server/core/manageres"
	modelmanager "mlib.com/gofy/server/core/manageres/model_manager"
	providermanager "mlib.com/gofy/server/core/manageres/provider_manager"
	"mlib.com/gofy/server/core/prompt"
	dbengine "mlib.com/gofy/server/db_engine"
	appconfigentities "mlib.com/gofy/server/entities/app/config"
	modelruntimeentities "mlib.com/gofy/server/entities/model_runtime"
	promptentities "mlib.com/gofy/server/entities/prompt"
	ragentities "mlib.com/gofy/server/entities/rag"
	modelenumtypes "mlib.com/gofy/server/enum_types/model"
	modelruntimeenumtypes "mlib.com/gofy/server/enum_types/model_runtime"
	promptenumttypes "mlib.com/gofy/server/enum_types/prompt"
	ragretrievalenumtypes "mlib.com/gofy/server/enum_types/rag/retrieval"
	"mlib.com/gofy/server/models"
	"mlib.com/gofy/server/utils"
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
	tenant_id string, model *appconfigentities.ModelConfig,
) (*modelmanager.ModelInstance, *appconfigentities.ModelConfigWithCredentialsEntity) {
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
	if provider_model == nil {
		panic(exceptions.NewValueError(fmt.Sprintf("Model %s not exist.", model_name)))
	}
	if provider_model.Status == modelenumtypes.ModelStatus_NO_CONFIGURE {
		panic(exceptions.NewValueError(fmt.Sprintf("Model {%s} credentials is not initialized.", model_name)))
	} else if provider_model.Status == modelenumtypes.ModelStatus_NO_PERMISSION {
		panic(exceptions.NewValueError(fmt.Sprintf("Gofy Hosted OpenAI {%s} currently not support.", model_name)))
	} else if provider_model.Status == modelenumtypes.ModelStatus_QUOTA_EXCEEDED {
		panic(exceptions.NewValueError(fmt.Sprintf("Model provider {%s} quota exceeded.", provider_name)))
	}
	// model config
	completion_params := model.CompletionParams
	stop := []string{}
	if _, ok := completion_params["stop"]; ok {
		stop = mapstruct.Get(completion_params, "stop", []string{})
		delete(completion_params, "stop")
	}
	// get model mode
	model_mode := model.Mode
	if string(model_mode) == "" {
		panic(exceptions.NewValueError("LLM mode is required."))
	}
	model_schema := model_type_instance.GetModelSchema(model_type_instance, model_name, model_credentials)

	if model_schema == nil {
		panic(exceptions.NewValueError(fmt.Sprintf("Model {%s} not exist.", model_name)))
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
func _get_prompt_template(
	model_config *appconfigentities.ModelConfigWithCredentialsEntity, mode string, metadata_fields []string, query string,
) ([]modelruntimeentities.PromptMessager, []string) {
	model_mode := promptenumttypes.ModelModeType(mode)
	input_text := query

	var prompt_messages []modelruntimeentities.PromptMessager
	if model_mode == promptenumttypes.ModelMode_CHAT {
		prompt_template := []*promptentities.ChatModelMessage{}
		system_prompt_messages := &promptentities.ChatModelMessage{Role: modelruntimeentities.PromptMessageRole_SYSTEM, Text: METADATA_FILTER_SYSTEM_PROMPT}
		prompt_template = append(prompt_template, system_prompt_messages)
		user_prompt_message_1 := &promptentities.ChatModelMessage{Role: modelruntimeentities.PromptMessageRole_USER, Text: METADATA_FILTER_USER_PROMPT_1}
		prompt_template = append(prompt_template, user_prompt_message_1)
		assistant_prompt_message_1 := &promptentities.ChatModelMessage{
			Role: modelruntimeentities.PromptMessageRole_ASSISTANT, Text: METADATA_FILTER_ASSISTANT_PROMPT_1,
		}
		prompt_template = append(prompt_template, assistant_prompt_message_1)
		user_prompt_message_2 := &promptentities.ChatModelMessage{Role: modelruntimeentities.PromptMessageRole_USER, Text: METADATA_FILTER_USER_PROMPT_2}
		prompt_template = append(prompt_template, user_prompt_message_2)
		assistant_prompt_message_2 := &promptentities.ChatModelMessage{
			Role: modelruntimeentities.PromptMessageRole_ASSISTANT, Text: METADATA_FILTER_ASSISTANT_PROMPT_2,
		}
		prompt_template = append(prompt_template, assistant_prompt_message_2)
		bindata, _ := json.Marshal(metadata_fields)
		user_prompt_message_3 := &promptentities.ChatModelMessage{
			Role: modelruntimeentities.PromptMessageRole_USER,
			Text: strings.ReplaceAll(strings.ReplaceAll(METADATA_FILTER_USER_PROMPT_3, "{input_text}", input_text), "{metadata_fields}", string(bindata)),
		}
		prompt_template = append(prompt_template, user_prompt_message_3)
		prompt_messages = prompt.NewAdvancedPromptTransform[[]*promptentities.ChatModelMessage](false, modelruntimeentities.ImagePromptMessageContentDETAIL("")).GetPrompt(
			prompt_template,
			map[string]string{},
			query,
			nil,
			"",
			nil,
			nil,
			model_config,
		)
	} else if model_mode == promptenumttypes.ModelMode_COMPLETION {
		bindata, _ := json.Marshal(metadata_fields)
		prompt_template := &promptentities.CompletionModelPromptTemplate{
			Text: strings.ReplaceAll(strings.ReplaceAll(METADATA_FILTER_COMPLETION_PROMPT, "{input_text}", input_text), "{metadata_fields}", string(bindata)),
		}
		prompt_messages = prompt.NewAdvancedPromptTransform[*promptentities.CompletionModelPromptTemplate](false, modelruntimeentities.ImagePromptMessageContentDETAIL("")).GetPrompt(
			prompt_template,
			map[string]string{},
			query,
			nil,
			"",
			nil,
			nil,
			model_config,
		)
	} else {
		panic(exceptions.NewValueError(fmt.Sprintf("Model mode {%v} not support.", model_mode)))
	}
	stop := model_config.Stop

	return prompt_messages, stop

}

func _process_metadata_filter_func(
	sequence int, condition string, metadata_name string, value any, filters []clause.Expression,
) []clause.Expression {
	if value == nil {
		return nil
	}
	// key := metadata_name + "_" + strconv.Itoa(sequence)
	// key_value := metadata_name + "_" + strconv.Itoa(sequence) + "_value"
	switch strings.ToLower(condition) {
	case "contains":
		filters = append(filters, clause.Expr{
			SQL:  "documents.doc_metadata ->> ? LIKE ?",
			Vars: []interface{}{metadata_name, clause.NamedExpr{SQL: fmt.Sprintf("%%%v%%", value)}},
		})
	case "not contains":
		filters = append(filters, clause.Expr{
			SQL:  "documents.doc_metadata ->> ? NOT LIKE ?",
			Vars: []interface{}{metadata_name, clause.NamedExpr{SQL: fmt.Sprintf("%%%v%%", value)}},
		})

	case "start with":
		filters = append(filters, clause.Expr{
			SQL:  "documents.doc_metadata ->> ? LIKE ?",
			Vars: []interface{}{metadata_name, clause.NamedExpr{SQL: fmt.Sprintf("%v%%", value)}},
		})

	case "end with":
		filters = append(filters, clause.Expr{
			SQL:  "documents.doc_metadata ->> ? LIKE ?",
			Vars: []interface{}{metadata_name, clause.NamedExpr{SQL: fmt.Sprintf("%%%v", value)}},
		})

	case "is", "=":
		if str, ok := value.(string); ok {
			// 字符串：JSON 值带双引号
			filters = append(filters, clause.Expr{
				SQL:  "documents.doc_metadata ->> ? = ?",
				Vars: []interface{}{metadata_name, fmt.Sprintf("%s", str)},
			})
		} else {
			// 数值：CAST -> FLOAT
			filters = append(filters, clause.Expr{
				SQL:  "CAST(documents.doc_metadata ->> ? AS FLOAT) = ?",
				Vars: []interface{}{metadata_name, value},
			})
		}

	case "is not", "≠":
		if str, ok := value.(string); ok {
			filters = append(filters, clause.Expr{
				SQL:  "documents.doc_metadata ->> ? != ?",
				Vars: []interface{}{metadata_name, fmt.Sprintf("%s", str)},
			})
		} else {
			filters = append(filters, clause.Expr{
				SQL:  "CAST(documents.doc_metadata ->> ? AS FLOAT) != ?",
				Vars: []interface{}{metadata_name, value},
			})
		}
	case "empty":
		filters = append(filters, clause.Expr{
			SQL:  "documents.doc_metadata ->> ? IS NULL",
			Vars: []interface{}{metadata_name},
		})
	case "not empty":
		filters = append(filters, clause.Expr{
			SQL:  "documents.doc_metadata ->> ? IS NOT NULL",
			Vars: []interface{}{metadata_name},
		})
	case "before", "<":
		filters = append(filters, clause.Expr{
			SQL:  "CAST(documents.doc_metadata ->> ? AS FLOAT) < ?",
			Vars: []interface{}{metadata_name, cast.ToFloat64(value)},
		})
	case "after", ">":
		filters = append(filters, clause.Expr{
			SQL:  "CAST(documents.doc_metadata ->> ? AS FLOAT) > ?",
			Vars: []interface{}{metadata_name, cast.ToFloat64(value)},
		})
	case "≤", "<=":
		filters = append(filters, clause.Expr{
			SQL:  "CAST(documents.doc_metadata ->> ? AS FLOAT) <= ?",
			Vars: []interface{}{metadata_name, cast.ToFloat64(value)},
		})
	case "≥", ">=":
		filters = append(filters, clause.Expr{
			SQL:  "CAST(documents.doc_metadata ->> ? AS FLOAT) >= ?",
			Vars: []interface{}{metadata_name, cast.ToFloat64(value)},
		})
	default:

	}
	return filters
}

func _handle_invoke_result(invoke_result iter.Seq[*modelruntimeentities.LLMResultChunk]) (string, *modelruntimeentities.LLMUsage) {
	var model string
	var prompt_messages []modelruntimeentities.PromptMessager
	full_text := ""
	var usage *modelruntimeentities.LLMUsage
	for result := range invoke_result {
		text := result.Delta.Message.Content
		full_text += text

		if model == "" {
			model = result.Model
		}
		if prompt_messages == nil {
			prompt_messages = result.PromptMessages
		}
		if usage == nil && result.Delta.Usage != nil {
			usage = result.Delta.Usage
		}
	}
	if usage == nil {
		usage = modelruntimeentities.NewLLMUsage()
	}
	return full_text, usage
}

func _automatic_metadata_filter_func(
	dataset_ids []string, query string, tenant_id string, user_id string, metadata_model_config *appconfigentities.ModelConfig,
) []map[string]any {
	// get all metadata field
	var metadata_fields []*models.DatasetMetadata
	err := dbengine.Instance().DB.Model(&models.DatasetMetadata{}).Where("dataset_id in ?", dataset_ids).Find(&metadata_fields).Error
	if err != nil {
		mlog.Errorf("get DatasetMetadata failed:%v", err)
	}
	all_metadata_fields := []string{}
	for _, metadata_field := range metadata_fields {
		all_metadata_fields = append(all_metadata_fields, metadata_field.Name)
	}

	// get metadata model config
	if metadata_model_config == nil {
		panic(exceptions.NewValueError("metadata_model_config is required"))
	}
	// get metadata model instance
	// fetch model config
	model_instance, model_config := _fetch_model_config(tenant_id, metadata_model_config)

	// fetch prompt messages
	prompt_messages, stop := _get_prompt_template(
		model_config,
		string(metadata_model_config.Mode),
		all_metadata_fields,
		query,
	)

	automatic_metadata_filters := []map[string]any{}
	func() {
		defer func() {
			if r := recover(); r != nil {
				mlog.Errorf("panic recover:%v\n%s", r, utils.GetCurrentGoroutineStack())
				if real_exp, ok := r.(error); ok {
					mlog.Errorf("Unknown Error when generating:%v", real_exp)

				} else {
					panic(r)
				}
			}
		}()
		// handle invoke result
		invoke_result := model_instance.InvokeLLMStream(
			prompt_messages,
			model_config.Parameters,
			nil,
			stop,
			user_id,
			nil,
		)

		// handle invoke result
		result_text, _ := _handle_invoke_result(invoke_result)

		result_text_json := utils.ParseAndCheckJsonMarkdown(result_text, []string{})

		if _, ok := result_text_json["metadata_map"]; ok {
			metadata_map := mapstruct.Get(result_text_json, "metadata_map", []map[string]any{})
			for _, item := range metadata_map {
				if slices.Contains(all_metadata_fields, mapstruct.Get(item, "metadata_field_name", "")) {
					automatic_metadata_filters = append(automatic_metadata_filters, map[string]any{
						"metadata_name": mapstruct.Get(item, "metadata_field_name", ""),
						"value":         item["metadata_field_value"],
						"condition":     mapstruct.Get(item, "comparison_operator", ""),
					})
				}
			}
		}
	}()

	return automatic_metadata_filters
}

var (
	// 1. 匹配 {{word}}
	meta_var_pattern = regexp.MustCompile(`\{\{(\w+)\}\}`)
	// 2. 合并连续空白字符
	white_space_pattern = regexp.MustCompile(`[\r\n\t]+`)
)

func _replace_metadata_filter_value(text string, inputs map[string]any) string {
	if len(inputs) == 0 {
		return text
	}

	// 替换 {{key}}
	out := meta_var_pattern.ReplaceAllStringFunc(text, func(s string) string {
		key := meta_var_pattern.FindStringSubmatch(s)[1] // 捕获组1
		if v, ok := inputs[key]; ok {
			return cast.To[string](v)
		}
		return s // 未找到保留原模板
	})

	// 合并空白字符并 trim
	out = white_space_pattern.ReplaceAllString(out, " ")
	return strings.TrimSpace(out)
}

func GetMetadataFilterCondition(
	dataset_ids []string,
	query string,
	tenant_id string,
	user_id string,
	metadata_filtering_mode string,
	metadata_model_config *appconfigentities.ModelConfig,
	metadata_filtering_conditions *appconfigentities.MetadataFilteringCondition,
	inputs map[string]any,
) (map[string][]string, *ragentities.MetadataCondition) {
	document_query := dbengine.Instance().DB.Model(&models.Document{}).Where("dataset_id in ? and indexing_status = ? and enabled = ? and archived = ?", dataset_ids, "completed", true, false)
	filters := []clause.Expression{} // type: ignore
	var metadata_condition *ragentities.MetadataCondition
	if metadata_filtering_mode == "disabled" {
		return nil, nil
	} else if metadata_filtering_mode == "automatic" {
		automatic_metadata_filters := _automatic_metadata_filter_func(
			dataset_ids, query, tenant_id, user_id, metadata_model_config,
		)
		if len(automatic_metadata_filters) > 0 {
			conditions := []*ragentities.Condition{}
			for sequence, filter := range automatic_metadata_filters {
				filters = _process_metadata_filter_func(
					sequence,
					filter["condition"].(string),     // type: ignore
					filter["metadata_name"].(string), // type: ignore
					filter["value"],
					filters, // type: ignore
				)
				switch filter["value"].(type) {
				case string:
				case []string:
				case int:
				case float64:
				default:
					mlog.Errorf("Invalid value=%#v type", filter["value"])
					panic(exceptions.NewValueError("Invalid value type"))
				}
				conditions = append(conditions, &ragentities.Condition{
					Name:               filter["metadata_name"].(string), // type: ignore
					ComparisonOperator: filter["condition"].(string),     // type: ignore
					Value:              filter["value"],
				})
			}
			metadata_condition = &ragentities.MetadataCondition{
				LogicalOperator: "or", // type: ignore
				Conditions:      conditions,
			}
			if metadata_filtering_conditions != nil {
				metadata_condition.LogicalOperator = string(metadata_filtering_conditions.LogicalOperator)
			}
			if !slices.Contains([]string{"and", "or"}, metadata_condition.LogicalOperator) {
				mlog.Errorf("Invalid LogicalOperator=%#v type", metadata_condition.LogicalOperator)
				panic(exceptions.NewValueError("Invalid LogicalOperator type"))
			}
		}
	} else if metadata_filtering_mode == "manual" {
		if metadata_filtering_conditions != nil {
			conditions := []*ragentities.Condition{}
			for sequence, condition := range metadata_filtering_conditions.Conditions { // type: ignore
				metadata_name := condition.Name
				expected_value := condition.Value
				if expected_value != nil && !slices.Contains([]string{"empty", "not empty"}, condition.ComparisonOperator) {
					if _, ok := expected_value.(string); ok {
						expected_value = _replace_metadata_filter_value(expected_value.(string), inputs)
					}
				}
				switch expected_value.(type) {
				case string:
				case []string:
				case int:
				case float64:
				default:
					mlog.Errorf("Invalid expected_value=%#v type", expected_value)
					panic(exceptions.NewValueError("Invalid expected_value type"))
				}
				conditions = append(conditions, &ragentities.Condition{
					Name:               metadata_name,
					ComparisonOperator: condition.ComparisonOperator,
					Value:              expected_value,
				})
				filters = _process_metadata_filter_func(
					sequence,
					condition.ComparisonOperator,
					metadata_name,
					expected_value,
					filters,
				)
			}
			metadata_condition = &ragentities.MetadataCondition{
				LogicalOperator: string(metadata_filtering_conditions.LogicalOperator),
				Conditions:      conditions,
			}
		}
	} else {
		panic(exceptions.NewValueError("Invalid metadata filtering mode"))
	}
	if len(filters) > 0 {
		if metadata_filtering_conditions != nil && metadata_filtering_conditions.LogicalOperator == "and" { // type: ignore
			document_query = document_query.Clauses(clause.And(filters...))
		} else {
			document_query = document_query.Clauses(clause.Or(filters...))
		}
	}
	var documents []*models.Document
	err := document_query.Find(&documents).Error
	if err != nil {
		mlog.Errorf("get documents failed:%v", err)
	}
	// group by dataset_id
	var metadata_filter_document_ids map[string][]string // type: ignore
	for _, document := range documents {
		if metadata_filter_document_ids == nil {
			metadata_filter_document_ids = make(map[string][]string)
		}
		if _, ok := metadata_filter_document_ids[document.DatasetID]; !ok {
			metadata_filter_document_ids[document.DatasetID] = make([]string, 0)
		}
		metadata_filter_document_ids[document.DatasetID] = append(metadata_filter_document_ids[document.DatasetID], document.ID) // type: ignore
	}
	return metadata_filter_document_ids, metadata_condition

}
