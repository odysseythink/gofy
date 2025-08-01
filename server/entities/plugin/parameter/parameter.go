package parameter

import (
	"fmt"
	"slices"

	"mlib.com/gofy/server/core/exceptions"
	pluginenumtypes "mlib.com/gofy/server/enum_types/plugin"
	commontypes "mlib.com/gofy/server/types/common"
)

type PluginParameterOption struct {
	Value string                 `json:"value"` //description="The value of the option")
	Label commontypes.I18nObject `json:"label"` //description="The label of the option")
	Icon  string                 `json:"icon"`  //description="The icon of the option, can be a url or a base64 encoded image"
}

func (option *PluginParameterOption) TransformIDToStr(value any) string {
	switch real_val := value.(type) {
	case string:
		return real_val
	default:
		return fmt.Sprintf("%v", real_val)
	}
}
func (option *PluginParameterOption) Copy() *PluginParameterOption {
	new_option := new(PluginParameterOption)
	new_option.Value = option.Value
	new_option.Label = commontypes.I18nObject{
		ZhHans: option.Label.ZhHans,
		EnUS:   option.Label.EnUS,
		PtBR:   option.Label.PtBR,
		JaJP:   option.Label.JaJP,
	}
	return new_option
}

type PluginParameterAutoGenerate struct {
	Type pluginenumtypes.PluginParameterAutoGenerateType `json:"type"`
}
type PluginParameterTemplate struct {
	Enabled bool `json:"enabled"` //description="Whether the parameter is jinja enabled")
}
type PluginParameter[T1 float64 | int | string, T2 float64 | int] struct {
	Name         string                       `json:"name"`        //description="The name of the parameter")
	Label        commontypes.I18nObject       `json:"label"`       //description="The label presented to the user")
	Placeholder  *commontypes.I18nObject      `json:"placeholder"` //description="The placeholder presented to the user")
	Scope        string                       `json:"scope"`
	AutoGenerate *PluginParameterAutoGenerate `json:"auto_generate"`
	Template     *PluginParameterTemplate     `json:"template"`
	Required     bool                         `json:"required"`
	Default      T1                           `json:"default"`
	Min          T2                           `json:"min"`
	Max          T2                           `json:"max"`
	Precision    int                          `json:"precision"`
	Options      []*PluginParameterOption     `json:"options"`
}

func InitFrontendParameter[T1 float64 | int | string, T2 float64 | int](rule *PluginParameter[T1, T2], typ string, value any) any {
	parameter_value := value
	if parameter_value == nil {
		// get default value
		parameter_value = rule.Default
		if parameter_value == nil && rule.Required {
			panic(exceptions.NewValueError(fmt.Sprintf("tool parameter {%s} not found in tool config", rule.Name)))
		}
	}
	if pluginenumtypes.PluginParameterType(typ) == pluginenumtypes.PluginParameter_SELECT {
		// check if tool_parameter_config in options
		options := []string{}
		for _, x := range rule.Options {
			options = append(options, x.Value)
		}
		if _, ok := parameter_value.(string); !ok {
			panic(exceptions.NewValueError(fmt.Sprintf("tool parameter {%s} value {%#v} must be string", rule.Name, parameter_value)))
		}
		if !slices.Contains(options, parameter_value.(string)) {
			panic(exceptions.NewValueError(fmt.Sprintf("tool parameter {%s} value {%#v} not in options {%#v}", rule.Name, parameter_value, options)))
		}
	}
	return pluginenumtypes.CastParameterValue(typ, parameter_value)
}
