package agent

import (
	parameterentities "mlib.com/gofy/server/entities/plugin/parameter"
	toolsentities "mlib.com/gofy/server/entities/tools"
	agentenumtypes "mlib.com/gofy/server/enum_types/agent"
	commontypes "mlib.com/gofy/server/types/common"
)

type AgentStrategyProviderIdentity struct {
	*toolsentities.ToolProviderIdentity
}
type AgentStrategyParameter struct {
	*parameterentities.PluginParameter
	Type agentenumtypes.AgentStrategyParameterType `json:"type"` //description="The type of the parameter")
	Help *commontypes.I18nObject                   `json:"help"`
}

func (param *AgentStrategyParameter) InitFrontendParameter(value any) any {
	return parameterentities.InitFrontendParameter(param.PluginParameter, string(param.Type), value)
}

type AgentStrategyProviderEntity struct {
	Identity AgentStrategyProviderIdentity `json:"identity"`
	PluginID *string                       `json:"plugin_id"` //description="The id of the plugin")
}
type AgentStrategyIdentity struct {
	*toolsentities.ToolIdentity
}
type AgentStrategyEntity[T1 float64 | int | string, T2 float64 | int] struct {
	Identity     AgentStrategyIdentity             `json:"identity"`
	Parameters   []*AgentStrategyParameter         `json:"parameters"`
	Description  commontypes.I18nObject            `json:"description"` //description="The description of the agent strategy")
	OutputSchema map[string]any                    `json:"output_schema"`
	Features     []agentenumtypes.AgentFeatureType `json:"features"`
	MetaVersion  *string                           `json:"meta_version"`
	//pydantic configs
	ModelConfig map[string]any `json:"model_config"`
}
type AgentProviderEntityWithPlugin[T1 float64 | int | string, T2 float64 | int] struct {
	*AgentStrategyProviderEntity
	Strategies []*AgentStrategyEntity[T1, T2] `json:"strategies"`
}
