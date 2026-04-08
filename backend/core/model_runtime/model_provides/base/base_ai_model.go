package base

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
	"mlib.com/confy"
	"mlib.com/gofy/server/core/exceptions"
	modelruntimeentities "mlib.com/gofy/server/entities/model_runtime"
	modelruntimeenumtypes "mlib.com/gofy/server/enum_types/model_runtime"
	"mlib.com/mlog"
)

type BaseAIModel struct {
	ModeType     modelruntimeenumtypes.ModelType
	ModelSchemas []*modelruntimeentities.AIModelEntity
	StartedAt    time.Time

	// pydantic configs
	ModelConfig *confy.Confy
}

func (m *BaseAIModel) GetCustomizableModelSchemaFromCredentials(modeler modelruntimeentities.AIModeler, model string, credentials map[string]any) *modelruntimeentities.AIModelEntity {
	/*
		Get customizable model schema from credentials

		:param model: model name
		:param credentials: model credentials
		:return: model schema
	*/
	return m.getCustomizableModelSchema(modeler, model, credentials)
}
func (m *BaseAIModel) getCustomizableModelSchema(modeler modelruntimeentities.AIModeler, model string, credentials map[string]any) *modelruntimeentities.AIModelEntity {
	/*
		Get customizable model schema and fill in the template
	*/
	schema := modeler.GetCustomizableModelSchema(model, credentials)

	if schema == nil {
		return nil
	}
	// fill in the template
	new_parameter_rules := []*modelruntimeentities.ParameterRule{}
	for _, parameter_rule := range schema.ParameterRules {
		if parameter_rule.UseTemplate != "" {
			// try:
			default_parameter_name := modelruntimeenumtypes.DefaultParameterNameType(parameter_rule.UseTemplate)
			default_parameter_rule, err := m.GetDefaultParameterRuleVariableMap(default_parameter_name)
			if err != nil {
				mlog.Error("get default parameter rule variable map failed:", err)
				continue
			}
			if parameter_rule.Max == 0.0 {
				if _, ok := default_parameter_rule["max"]; ok {
					if _, ok := default_parameter_rule["max"].(float64); ok {
						parameter_rule.Max = default_parameter_rule["max"].(float64)
					}
				}
			}
			if parameter_rule.Min == 0.0 {
				if _, ok := default_parameter_rule["min"]; ok {
					if _, ok := default_parameter_rule["min"].(float64); ok {
						parameter_rule.Min = default_parameter_rule["min"].(float64)
					}
				}
			}
			if parameter_rule.Precision == 0 {
				if _, ok := default_parameter_rule["precision"]; ok {
					if _, ok := default_parameter_rule["precision"].(int); ok {
						parameter_rule.Precision = default_parameter_rule["precision"].(int)
					}
				}
			}

			if _, ok := default_parameter_rule["Required"]; ok {
				if _, ok := default_parameter_rule["Required"].(bool); ok {
					parameter_rule.Required = default_parameter_rule["Required"].(bool)
				}
			}

			if parameter_rule.Default == nil {
				if _, ok := default_parameter_rule["default"]; ok {

					parameter_rule.Default = default_parameter_rule["default"]

				}
			}
			if parameter_rule.Help.EnUS == "" {
				if _, ok := default_parameter_rule["help"]; ok {
					if _, ok := default_parameter_rule["help"].(map[string]any); ok {
						if _, ok := default_parameter_rule["help"].(map[string]any)["en_US"]; ok {
							if _, ok := default_parameter_rule["help"].(map[string]any)["en_US"].(string); ok {
								parameter_rule.Help.EnUS = default_parameter_rule["help"].(map[string]any)["en_US"].(string)
							}
						}
					}
				}
			}
			if parameter_rule.Help.ZhHans == "" {
				if _, ok := default_parameter_rule["help"]; ok {
					if _, ok := default_parameter_rule["help"].(map[string]any); ok {
						if _, ok := default_parameter_rule["help"].(map[string]any)["zh_Hans"]; ok {
							if _, ok := default_parameter_rule["help"].(map[string]any)["zh_Hans"].(string); ok {
								parameter_rule.Help.ZhHans = default_parameter_rule["help"].(map[string]any)["zh_Hans"].(string)
							}
						}
					}
				}
			}

			// except ValueError:
			// 	pass
		}
		new_parameter_rules = append(new_parameter_rules, parameter_rule)
	}
	schema.ParameterRules = new_parameter_rules

	return schema
}

func (m *BaseAIModel) GetDefaultParameterRuleVariableMap(name modelruntimeenumtypes.DefaultParameterNameType) (map[string]any, error) {
	/*
		Get default parameter rule for given name

		:param name: parameter name
		:return: parameter rule
	*/
	default_parameter_rule := map[string]any{}
	if _, ok := modelruntimeentities.PARAMETER_RULE_TEMPLATE[name]; ok {
		default_parameter_rule = modelruntimeentities.PARAMETER_RULE_TEMPLATE[name]
	}

	if len(default_parameter_rule) == 0 {
		return nil, fmt.Errorf("invalid model parameter rule name %s", name)
	}
	return default_parameter_rule, nil
}

func (m *BaseAIModel) PredefinedModels(modeler modelruntimeentities.AIModeler) []*modelruntimeentities.AIModelEntity {
	/*
		Get all predefined models for given provider.

		:return:
	*/
	if m.ModelSchemas != nil {
		return m.ModelSchemas
	}
	model_schemas := []*modelruntimeentities.AIModelEntity{}

	// get module name
	model_type := string(modeler.ModelType())

	// get provider name
	provider_name := modeler.ProviderName()

	// get parent path of the current path
	provider_model_type_path := filepath.Join("model_runtime", "model_provides", provider_name, model_type)

	rd, err := os.ReadDir(provider_model_type_path)
	if err != nil {
		mlog.Error("read dir failed:", err)
		return nil
	}
	model_schema_yaml_paths := []string{}
	for _, fi := range rd {
		if !fi.IsDir() {
			if !strings.HasPrefix(fi.Name(), "__") &&
				!strings.HasPrefix(fi.Name(), "_") &&
				strings.HasSuffix(fi.Name(), ".yaml") {
				model_schema_yaml_paths = append(model_schema_yaml_paths, filepath.Join(provider_model_type_path, fi.Name()))
			}
		}
	}

	// get _position.yaml file path
	// position_map := positionhelper.GetPositionMap(provider_model_type_path, "")

	// // traverse all model_schema_yaml_paths
	for _, model_schema_yaml_path := range model_schema_yaml_paths {
		// read yaml data from yaml file
		filecontent, err := os.ReadFile(model_schema_yaml_path)
		if err != nil {
			mlog.Errorf("read file(%s) failed:%v", model_schema_yaml_path, err)
			continue
		}
		yaml_data := map[string]any{}
		err = yaml.Unmarshal(filecontent, &yaml_data)
		if err != nil {
			mlog.Errorf("yaml.Unmarshal(%s) failed:%v", string(filecontent), err)
			continue
		}

		new_parameter_rules := []map[string]any{}
		parameter_rules := []any{}
		if _, ok := yaml_data["parameter_rules"]; ok {
			if _, ok := yaml_data["parameter_rules"].([]any); ok {
				parameter_rules = yaml_data["parameter_rules"].([]any)
			}
		}
		for _, parameter_rule := range parameter_rules {
			if _, ok := parameter_rule.(map[string]any); !ok {
				if _, ok := parameter_rule.(map[any]any); !ok {
					mlog.Errorf("parameter_rule(%#v) must be dict", parameter_rule)
					continue
				} else {
					parameter_rule_dict := map[string]any{}
					for k, v := range parameter_rule.(map[any]any) {
						if _, ok := k.(string); !ok {
							mlog.Errorf("parameter_rule(%#v) must be dict", parameter_rule)
							continue
						}
						parameter_rule_dict[k.(string)] = v
					}
					parameter_rule = parameter_rule_dict
				}
			}
			parameter_rule_dict := parameter_rule.(map[string]any)
			if _, ok := parameter_rule_dict["use_template"]; ok {
				if _, ok := parameter_rule_dict["use_template"].(string); ok {
					default_parameter_name := modelruntimeenumtypes.DefaultParameterNameType(parameter_rule_dict["use_template"].(string))
					default_parameter_rule, err := m.GetDefaultParameterRuleVariableMap(default_parameter_name)
					if err != nil {
						mlog.Error("get default parameter rule variable map failed:", err)
						continue
					}
					for k, v := range default_parameter_rule {
						parameter_rule_dict[k] = v
					}
				}
			}
			if _, ok := parameter_rule_dict["label"]; !ok {
				parameter_rule_dict["label"] = map[string]any{"zh_Hans": parameter_rule_dict["name"], "en_US": parameter_rule_dict["name"]}
			}
			new_parameter_rules = append(new_parameter_rules, parameter_rule_dict)
		}
		yaml_data["parameter_rules"] = new_parameter_rules

		if _, ok := yaml_data["label"]; !ok {
			yaml_data["label"] = map[string]any{"zh_Hans": yaml_data["model"], "en_US": yaml_data["model"]}
		}
		yaml_data["fetch_from"] = modelruntimeenumtypes.FetchFrom_PREDEFINED_MODEL

		// yaml_data to entity
		model_schema := new(modelruntimeentities.AIModelEntity)
		err = yaml.Unmarshal(filecontent, model_schema)
		if err != nil {
			mlog.Errorf("invalid model schema for %s.%s.%s :%v", provider_name, model_type, model_schema_yaml_path, err)
			continue
		}

		// cache model schema
		model_schemas = append(model_schemas, model_schema)
	}
	// resort model schemas by position
	// model_schemas = sort_by_position_map(position_map, model_schemas, lambda x: x.model)

	// cache model schemas
	m.ModelSchemas = model_schemas

	return model_schemas
}
func (m *BaseAIModel) GetModelSchema(modeler modelruntimeentities.AIModeler, model string, credentials map[string]any) *modelruntimeentities.AIModelEntity {
	/*
		Get model schema by model name and credentials

		:param model: model name
		:param credentials: model credentials
		:return: model schema
	*/
	//Try to get model schema from predefined models
	for _, predefined_model := range m.PredefinedModels(modeler) {
		if model == predefined_model.Model {
			return predefined_model
		}
	}
	// Try to get model schema from credentials
	if len(credentials) > 0 {
		model_schema := m.GetCustomizableModelSchemaFromCredentials(modeler, model, credentials)
		if model_schema != nil {
			return model_schema
		}
	}
	return nil
}
func (m *BaseAIModel) GetPrice(modeler modelruntimeentities.AIModeler, model string, credentials map[string]any, price_type modelruntimeentities.PriceType, tokens int) (*modelruntimeentities.PriceInfo, error) {
	/*
		Get price for given model and tokens

		:param model: model name
		:param credentials: model credentials
		:param price_type: price type
		:param tokens: number of tokens
		:return: price info
	*/
	// get model schema
	model_schema := m.GetModelSchema(modeler, model, credentials)

	// get price info from predefined model schema
	var price_config *modelruntimeentities.PriceConfig
	if model_schema != nil && model_schema.Pricing != nil {
		price_config = model_schema.Pricing
	}
	// get unit price
	var unit_price float64
	if price_config != nil {
		if price_type == modelruntimeentities.PriceType_INPUT {
			unit_price = price_config.Input
		} else if price_type == modelruntimeentities.PriceType_OUTPUT {
			unit_price = price_config.Output
		}
	}
	if unit_price == 0.0 {
		return &modelruntimeentities.PriceInfo{
			UnitPrice:   0.0,
			Unit:        0.0,
			TotalAmount: 0.0,
			Currency:    "USD",
		}, nil
	}
	// calculate total amount
	if price_config == nil {
		return nil, exceptions.NewValueError("Price config not found for model " + model)
	}

	return &modelruntimeentities.PriceInfo{
		UnitPrice:   unit_price,
		Unit:        price_config.Unit,
		TotalAmount: float64(tokens) * unit_price * price_config.Unit,
		Currency:    price_config.Currency,
	}, nil
}

// func (m *BaseAIModel) GetNumTokensByGPT2(text string) int {
/*
	Get number of tokens for given prompt messages by gpt2
	Some provider models do not provide an interface for obtaining the number of tokens.
	Here, the gpt2 tokenizer is used to calculate the number of tokens.
	This method can be executed offline, and the gpt2 tokenizer has been cached in the project.

	:param text: plain text of prompt. You need to convert the original message to plain text
	:return: number of tokens
*/
// return GPT2Tokenizer.get_num_tokens(text)
// }
