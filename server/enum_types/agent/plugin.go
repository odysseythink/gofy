package agent

import (
	parameterenumtypes "mlib.com/gofy/server/enum_types/parameter"
	pluginenumtypes "mlib.com/gofy/server/enum_types/plugin"
)

type AgentStrategyParameterType string

const (
	AgentStrategyParameter_STRING         = AgentStrategyParameterType(parameterenumtypes.CommonParameter_STRING)
	AgentStrategyParameter_NUMBER         = AgentStrategyParameterType(parameterenumtypes.CommonParameter_NUMBER)
	AgentStrategyParameter_BOOLEAN        = AgentStrategyParameterType(parameterenumtypes.CommonParameter_BOOLEAN)
	AgentStrategyParameter_SELECT         = AgentStrategyParameterType(parameterenumtypes.CommonParameter_SELECT)
	AgentStrategyParameter_SECRET_INPUT   = AgentStrategyParameterType(parameterenumtypes.CommonParameter_SECRET_INPUT)
	AgentStrategyParameter_FILE           = AgentStrategyParameterType(parameterenumtypes.CommonParameter_FILE)
	AgentStrategyParameter_FILES          = AgentStrategyParameterType(parameterenumtypes.CommonParameter_FILES)
	AgentStrategyParameter_APP_SELECTOR   = AgentStrategyParameterType(parameterenumtypes.CommonParameter_APP_SELECTOR)
	AgentStrategyParameter_MODEL_SELECTOR = AgentStrategyParameterType(parameterenumtypes.CommonParameter_MODEL_SELECTOR)
	AgentStrategyParameter_TOOLS_SELECTOR = AgentStrategyParameterType(parameterenumtypes.CommonParameter_TOOLS_SELECTOR)
	AgentStrategyParameter_ANY            = AgentStrategyParameterType(parameterenumtypes.CommonParameter_ANY)

	// deprecated, should not use.
	AgentStrategyParameter_SYSTEM_FILES = AgentStrategyParameterType(parameterenumtypes.CommonParameter_SYSTEM_FILES)
)

func (param AgentStrategyParameterType) AsNormalType() string {
	return pluginenumtypes.AsNormalType(string(param))
}
func (param AgentStrategyParameterType) CastValue(value any) any {
	return pluginenumtypes.CastParameterValue(string(param), value)
}

type AgentFeatureType string

const (
	AgentFeature_HISTORY_MESSAGES AgentFeatureType = "history-messages"
)
