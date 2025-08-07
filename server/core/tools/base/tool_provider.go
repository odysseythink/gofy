package base

import (
	"fmt"
	"slices"

	"mlib.com/gofy/server/core/exceptions"
	providerentities "mlib.com/gofy/server/entities/provider"
	toolsentities "mlib.com/gofy/server/entities/tools"
	providerenumtypes "mlib.com/gofy/server/enum_types/provider"
	toolsenumtypes "mlib.com/gofy/server/enum_types/tools"
)

type ToolProviderController interface {
	GetTool(tool_name string) *Tooler
	GetCredentialsSchema() []*providerentities.ProviderConfig
	/*
		validate the credentials of the provider

		:param tool_name: the name of the tool, defined in `get_tools`
		:param credentials: the credentials of the tool
	*/
	ValidateCredentialsFormat(credentials map[string]any)
	ProviderType() toolsenumtypes.ToolProviderType
}

type BaseToolProviderController struct {
	Entity toolsentities.ToolProviderEntity `json:"entity"`
}

func (controller *BaseToolProviderController) GetCredentialsSchema() []*providerentities.ProviderConfig {
	return controller.Entity.CredentialsSchema
}

func (controller *BaseToolProviderController) ProviderType() toolsenumtypes.ToolProviderType {
	return toolsenumtypes.ToolProvider_BUILT_IN
}

// func (controller *BaseToolProviderController) ValidateParameters(real_controller ToolProviderController, tool_id int, tool_name string, tool_parameters map[string]any) {
// 	tool_parameters_schema := controller.GetParameters(real_controller, tool_name)

// 	tool_parameters_need_to_validate := map[string]*toolsentities.ToolParameter{}
// 	for _, parameter := range tool_parameters_schema {
// 		tool_parameters_need_to_validate[parameter.Name] = parameter
// 	}
// 	for tool_parameter := range tool_parameters {
// 		if _, ok := tool_parameters_need_to_validate[tool_parameter]; !ok {
// 			panic(exceptions.NewToolParameterValidationError(fmt.Sprintf("parameter %s not found in tool {tool_name}", tool_parameter)))
// 		}
// 		// check type
// 		parameter_schema := tool_parameters_need_to_validate[tool_parameter]
// 		switch parameter_schema.Type {
// 		case toolsenumtypes.ToolParameter_STRING:
// 			if _, ok := tool_parameters[tool_parameter].(string); !ok {
// 				panic(exceptions.NewToolParameterValidationError(fmt.Sprintf("parameter %s should be string", tool_parameter)))
// 			}
// 		case toolsenumtypes.ToolParameter_NUMBER:
// 			if real_val, ok := tool_parameters[tool_parameter].(int); ok {
// 				if float64(real_val) < parameter_schema.Min {
// 					panic(exceptions.NewToolParameterValidationError(fmt.Sprintf("tool_parameters[%s]=%d should be greater than %f", tool_parameter, real_val, parameter_schema.Min)))
// 				}
// 				if float64(real_val) > parameter_schema.Max {
// 					panic(exceptions.NewToolParameterValidationError(fmt.Sprintf("tool_parameters[%s]=%d should be less than %f", tool_parameter, real_val, parameter_schema.Max)))
// 				}
// 			} else if real_val, ok := tool_parameters[tool_parameter].(float64); ok {
// 				if real_val < parameter_schema.Min {
// 					panic(exceptions.NewToolParameterValidationError(fmt.Sprintf("tool_parameters[%s]=%f should be greater than %f", tool_parameter, real_val, parameter_schema.Min)))
// 				}
// 				if real_val > parameter_schema.Max {
// 					panic(exceptions.NewToolParameterValidationError(fmt.Sprintf("tool_parameters[%s]=%f should be less than %f", tool_parameter, real_val, parameter_schema.Max)))
// 				}
// 			} else {
// 				panic(exceptions.NewToolParameterValidationError(fmt.Sprintf("parameter %s should be number", tool_parameter)))
// 			}
// 		case toolsenumtypes.ToolParameter_BOOLEAN:
// 			if _, ok := tool_parameters[tool_parameter].(bool); !ok {
// 				panic(exceptions.NewToolParameterValidationError(fmt.Sprintf("parameter %s should be boolean", tool_parameter)))
// 			}
// 		case toolsenumtypes.ToolParameter_SELECT:
// 			if real_val, ok := tool_parameters[tool_parameter].(string); !ok {
// 				panic(exceptions.NewToolParameterValidationError(fmt.Sprintf("parameter %s should be string", tool_parameter)))
// 			} else {
// 				if !slices.ContainsFunc(parameter_schema.Options, func(option *toolsentities.ToolParameterOption) bool {
// 					return option.Value == real_val
// 				}) {
// 					panic(exceptions.NewToolParameterValidationError(fmt.Sprintf("parameter %s should be one of %#v", tool_parameter, parameter_schema.Options)))
// 				}
// 			}
// 		}
// 		delete(tool_parameters_need_to_validate, tool_parameter)
// 	}
// 	for tool_parameter_validate := range tool_parameters_need_to_validate {
// 		parameter_schema := tool_parameters_need_to_validate[tool_parameter_validate]
// 		if parameter_schema.Required {
// 			panic(exceptions.NewToolParameterValidationError(fmt.Sprintf("parameter %s is required", tool_parameter_validate)))
// 		}
// 		// the parameter is not set currently, set the default value if needed
// 		if parameter_schema.Default != nil {
// 			tool_parameters[tool_parameter_validate] = parameter_schema.Default
// 		}
// 	}
// }

func (controller *BaseToolProviderController) ValidateCredentialsFormat(credentials map[string]any) {
	credentials_schema := map[string]*providerentities.ProviderConfig{}
	if credentials_schema == nil {
		return
	}
	for _, credential := range controller.Entity.CredentialsSchema {
		credentials_schema[credential.Name] = credential
	}
	credentials_need_to_validate := map[string]*providerentities.ProviderConfig{}
	for credential_name := range credentials_schema {
		credentials_need_to_validate[credential_name] = credentials_schema[credential_name]
	}
	for credential_name := range credentials {
		if _, ok := credentials_need_to_validate[credential_name]; !ok {
			panic(exceptions.NewToolProviderCredentialValidationError(fmt.Sprintf("credential {%s} not found in provider {%s}", credential_name, controller.Entity.Identity.Name)))
		}
		// check type
		credential_schema := credentials_need_to_validate[credential_name]
		if !credential_schema.Required && credentials[credential_name] == nil {
			continue
		}
		if slices.Contains([]providerenumtypes.BasicProviderConfigType{providerenumtypes.BasicProviderConfig_SECRET_INPUT, providerenumtypes.BasicProviderConfig_TEXT_INPUT}, credential_schema.Type) {
			if _, ok := credentials[credential_name].(string); !ok {
				panic(exceptions.NewToolProviderCredentialValidationError(fmt.Sprintf("credential {%s} should be string", credential_name)))
			}
		} else if credential_schema.Type == providerenumtypes.BasicProviderConfig_SELECT {
			if _, ok := credentials[credential_name].(string); !ok {
				panic(exceptions.NewToolProviderCredentialValidationError(fmt.Sprintf("credential {%s} should be string", credential_name)))
			}
			options := credential_schema.Options
			// if not isinstance(options, list){
			//     panic(exceptions.NewToolProviderCredentialValidationError(fmt.Sprintf("credential {%s} options should be list", credential_name)))
			// }
			values := []string{}
			for _, x := range options {
				values = append(values, x.Value)
			}
			if !slices.Contains(values, credentials[credential_name].(string)) {
				panic(exceptions.NewToolProviderCredentialValidationError(fmt.Sprintf("credential {%s} should be one of {%#v}", credential_name, options)))
			}
		}
		delete(credentials_need_to_validate, credential_name)
	}
	for credential_name := range credentials_need_to_validate {
		credential_schema := credentials_need_to_validate[credential_name]
		if credential_schema.Required {
			panic(exceptions.NewToolProviderCredentialValidationError(fmt.Sprintf("credential {%s} is required", credential_name)))
		}
		// the credential is not set currently, set the default value if needed
		if credential_schema.Default != nil {
			default_value := credential_schema.Default
			// parse default value into the correct type
			if slices.Contains([]providerenumtypes.BasicProviderConfigType{
				providerenumtypes.BasicProviderConfig_SECRET_INPUT,
				providerenumtypes.BasicProviderConfig_TEXT_INPUT,
				providerenumtypes.BasicProviderConfig_SELECT,
			}, credential_schema.Type) {
				if _, ok := default_value.(string); !ok {
					panic(exceptions.NewToolProviderCredentialValidationError(fmt.Sprintf("credential {%s} default_value should be string", credential_name)))
				}
			}
			credentials[credential_name] = default_value
		}
	}
}
