package agent

import (
	"slices"

	agententities "mlib.com/gofy/server/entities/agent"
	promptentities "mlib.com/gofy/server/entities/prompt"
	toolsenumtypes "mlib.com/gofy/server/enum_types/tools"
	"mlib.com/gofy/server/utils/mapstruct"
)

type AgentConfigManager struct {
}

func (mgr *AgentConfigManager) Convert(config map[string]any) *agententities.AgentEntity {
	if _, ok := config["agent_mode"]; !ok || config["agent_mode"] == nil {
		return nil
	}
	if _, ok := config["agent_mode"].(map[string]any); !ok {
		return nil
	}
	if _, ok := config["agent_mode"].(map[string]any)["enabled"]; !ok {
		return nil
	}
	agent_dict := mapstruct.Get(config, "agent_mode", map[string]any{})
	agent_strategy := mapstruct.Get(agent_dict, "strategy", "cot")
	var strategy agententities.AgentEntityStrategyType
	if agent_strategy == "function_call" {
		strategy = agententities.AgentEntityStrategy_FUNCTION_CALLING
	} else if slices.Contains([]string{"cot", "react"}, agent_strategy) {
		strategy = agententities.AgentEntityStrategy_CHAIN_OF_THOUGHT
	} else {
		// old configs, try to detect default strategy
		if mapstruct.Get(mapstruct.Get(config, "model", map[string]any{}), "provider", "") == "openai" {
			strategy = agententities.AgentEntityStrategy_FUNCTION_CALLING
		} else {
			strategy = agententities.AgentEntityStrategy_CHAIN_OF_THOUGHT
		}
	}
	agent_tools := []*agententities.AgentToolEntity{}
	for _, tool := range mapstruct.Get(agent_dict, "tools", []map[string]any{}) {
		// keys = tool.keys()
		if len(tool) >= 4 {
			if _, ok := tool["enabled"]; !ok || tool["enabled"] == nil {
				return nil
			}
			if !mapstruct.Get(tool, "enabled", false) {
				continue
			}

			agent_tools = append(agent_tools, &agententities.AgentToolEntity{
				ProviderType:   toolsenumtypes.ToolProviderType(mapstruct.Get(tool, "provider_type", "")),
				ProviderID:     mapstruct.Get(tool, "provider_id", ""),
				ToolName:       mapstruct.Get(tool, "tool_name", ""),
				ToolParameters: mapstruct.Get(tool, "tool_parameters", map[string]any{}),
				CredentialID:   mapstruct.Get(tool, "credential_id", ""),
			})
		}
	}
	if _, ok := agent_dict["strategy"]; ok {
		if _, ok := agent_dict["strategy"].(string); ok && !slices.Contains([]string{
			"react_router",
			"router",
		}, agent_dict["strategy"].(string)) {
			agent_prompt := mapstruct.Get(agent_dict, "prompt", map[string]any{})
			// check model mode
			model_mode := mapstruct.Get(mapstruct.Get(config, "model", map[string]any{}), "mode", "completion")
			var agent_prompt_entity *agententities.AgentPromptEntity
			if model_mode == "completion" {
				agent_prompt_entity = &agententities.AgentPromptEntity{
					FirstPrompt:   mapstruct.Get(agent_prompt, "first_prompt", promptentities.REACT_PROMPT_TEMPLATES["english"]["completion"]["prompt"]),
					NextIteration: mapstruct.Get(agent_prompt, "next_iteration", promptentities.REACT_PROMPT_TEMPLATES["english"]["completion"]["agent_scratchpad"]),
				}
			} else {
				agent_prompt_entity = &agententities.AgentPromptEntity{
					FirstPrompt:   mapstruct.Get(agent_prompt, "first_prompt", promptentities.REACT_PROMPT_TEMPLATES["english"]["chat"]["prompt"]),
					NextIteration: mapstruct.Get(agent_prompt, "next_iteration", promptentities.REACT_PROMPT_TEMPLATES["english"]["chat"]["agent_scratchpad"]),
				}
			}
			return &agententities.AgentEntity{
				Provider:     mapstruct.Get(mapstruct.Get(config, "model", map[string]any{}), "provider", ""),
				Model:        mapstruct.Get(mapstruct.Get(config, "model", map[string]any{}), "name", ""),
				Strategy:     strategy,
				Prompt:       agent_prompt_entity,
				Tools:        agent_tools,
				MaxIteration: mapstruct.Get(agent_dict, "max_iteration", 10),
			}
		}
	}

	return nil
}
