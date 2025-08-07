package agent

import (
	"strings"

	toolsentities "mlib.com/gofy/server/entities/tools"
)

type AgentToolEntity struct {
	/*
	   Agent Tool Entity.
	*/

	ProviderType   string         `json:"provider_type"` //Literal["builtin", "api", "workflow"]
	ProviderID     string         `json:"provider_id"`
	ToolName       string         `json:"tool_name"`
	ToolParameters map[string]any `json:"tool_parameters"`
}

type AgentPromptEntity struct {
	/*
	   Agent Prompt Entity.
	*/

	FirstPrompt   string `json:"first_prompt"`
	NextIteration string `json:"next_iteration"`
}

type Action struct {
	/*
		Action Entity.
	*/

	ActionName  string `json:"action_name"`
	ActionInput any/*: Union[dict, str]*/ `json:"action_input"`
}

func (a *Action) ToDict() map[string]any {
	/*
		Convert to dictionary.
	*/
	return map[string]any{
		"action":       a.ActionName,
		"action_input": a.ActionInput,
	}
}

type AgentScratchpadUnit struct {
	/*
	   Agent First Prompt Entity.
	*/

	AgentResponse *string `json:"agent_response"`
	Thought       *string `json:"thought"`
	ActionStr     *string `json:"action_str"`
	Observation   *string `json:"observation"`
	Action        *Action `json:"action"`
}

func (a *AgentScratchpadUnit) IsFinal() bool {
	/*
	   Check if the scratchpad unit is final.
	*/
	return a.Action == nil || (strings.Contains(strings.ToLower(a.Action.ActionName), "final") && strings.Contains(strings.ToLower(a.Action.ActionName), "answer"))
}

type AgentEntityStrategy string

const (
	/*
	   Agent Strategy.
	*/

	AgentEntityStrategy_CHAIN_OF_THOUGHT AgentEntityStrategy = "chain-of-thought"
	AgentEntityStrategy_FUNCTION_CALLING AgentEntityStrategy = "function-calling"
)

type AgentEntity struct {
	/*
	   Agent Entity.
	*/

	Provider     string              `json:"provider"`
	Model        string              `json:"model"`
	Strategy     AgentEntityStrategy `json:"strategy"`
	Prompt       *AgentPromptEntity  `json:"prompt"`
	Tools        []*AgentToolEntity  `json:"tools"`
	MaxIteration int                 `json:"max_iteration"`
}

func NewAgentEntity() *AgentEntity {
	return &AgentEntity{
		MaxIteration: 5,
	}
}

type AgentInvokeMessage struct {
	*toolsentities.ToolInvokeMessage
}
