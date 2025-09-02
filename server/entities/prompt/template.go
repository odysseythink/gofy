package prompt

const (
	ENGLISH_REACT_COMPLETION_PROMPT_TEMPLATES = `Respond to the human as helpfully and accurately as possible.

{{instruction}}

You have access to the following tools:

{{tools}}

Use a json blob to specify a tool by providing an action key (tool name) and an action_input key (tool input).
Valid "action" values: "Final Answer" or {{tool_names}}

Provide only ONE action per $JSON_BLOB, as shown:

` + "```\n{\n  \"action\": $TOOL_NAME,\n  \"action_input\": $ACTION_INPUT\n}\n```" + `

Follow this format:

Question: input question to answer
Thought: consider previous and subsequent steps
Action:
` + "```\n$JSON_BLOB\n```" + `
Observation: action result
... (repeat Thought/Action/Observation N times)
Thought: I know what to respond
Action:
` + "```\n{\n  \"action\": \"Final Answer\",\n  \"action_input\": \"Final response to human\"\n}\n```" + `

Begin! Reminder to ALWAYS respond with a valid json blob of a single action. Use tools if necessary. Respond directly if appropriate. Format is Action:` + "```$JSON_BLOB```" + `then Observation:.
{{historic_messages}}
Question: {{query}}
{{agent_scratchpad}}
Thought:` // noqa: E501

	ENGLISH_REACT_COMPLETION_AGENT_SCRATCHPAD_TEMPLATES = `Observation: {{observation}}
Thought:`

	ENGLISH_REACT_CHAT_PROMPT_TEMPLATES = `Respond to the human as helpfully and accurately as possible.

{{instruction}}

You have access to the following tools:

{{tools}}

Use a json blob to specify a tool by providing an action key (tool name) and an action_input key (tool input).
Valid "action" values: "Final Answer" or {{tool_names}}

Provide only ONE action per $JSON_BLOB, as shown:

` + "```\n{\n  \"action\": $TOOL_NAME,\n  \"action_input\": $ACTION_INPUT\n}\n```" + `

Follow this format:

Question: input question to answer
Thought: consider previous and subsequent steps
Action:
` + "```\n$JSON_BLOB\n```" + `
Observation: action result
... (repeat Thought/Action/Observation N times)
Thought: I know what to respond
Action:
` + "```\n{\n  \"action\": \"Final Answer\",\n  \"action_input\": \"Final response to human\"\n}\n```" + `

Begin! Reminder to ALWAYS respond with a valid json blob of a single action. Use tools if necessary. Respond directly if appropriate. Format is Action:` + "```$JSON_BLOB```" + `then Observation:.
` // noqa: E501

	ENGLISH_REACT_CHAT_AGENT_SCRATCHPAD_TEMPLATES = ""
)

var (
	REACT_PROMPT_TEMPLATES = map[string]map[string]map[string]string{
		"english": {
			"chat": {
				"prompt":           ENGLISH_REACT_CHAT_PROMPT_TEMPLATES,
				"agent_scratchpad": ENGLISH_REACT_CHAT_AGENT_SCRATCHPAD_TEMPLATES,
			},
			"completion": {
				"prompt":           ENGLISH_REACT_COMPLETION_PROMPT_TEMPLATES,
				"agent_scratchpad": ENGLISH_REACT_COMPLETION_AGENT_SCRATCHPAD_TEMPLATES,
			},
		},
	}
)
