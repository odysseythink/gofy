package prompt

import (
	modelruntimeentities "github.com/odysseythink/gofy/backend/entities/model_runtime"
)

type EditionType string

const (
	Edition_Basic  EditionType = "basic"
	Edition_Jinja2 EditionType = "jinja2"
)

type ChatModelMessage struct {
	/*
	   Chat Message.
	*/

	Text        string                                 `json:"text"`
	Role        modelruntimeentities.PromptMessageRole `json:"role"`
	EditionType `json:"edition_type"`                  // Optional[Literal["basic", "jinja2"]] = None

}
type CompletionModelPromptTemplate struct {
	/*
	   Completion Model Prompt Template.
	*/

	Text        string `json:"text"`
	EditionType `json:"edition_type"`
}
type RolePrefix struct {
	/*
	   Role Prefix.
	*/

	User      string `json:"user"`
	Assistant string `json:"assistant"`
}
type WindowConfig struct {
	/*
	   Window Config.
	*/

	Enabled bool `json:"enabled"`
	Size    int  `json:"size"`
}
type MemoryConfig struct {
	/*
	   Memory Config.
	*/

	RolePrefix          *RolePrefix  `json:"role_prefix"`
	Window              WindowConfig `json:"window"`
	QueryPromptTemplate string       `json:"query_prompt_template"`
}
