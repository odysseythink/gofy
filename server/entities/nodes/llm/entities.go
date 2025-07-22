package llm

import (
	"encoding/json"

	"mlib.com/gofy/server/core/exceptions"
	modelruntimeentities "mlib.com/gofy/server/entities/model_runtime"
	basenodesentities "mlib.com/gofy/server/entities/nodes/base"
	promptentities "mlib.com/gofy/server/entities/prompt"
	workflowentities "mlib.com/gofy/server/entities/workflow"
	"mlib.com/mlog"
)

type ModelConfig struct {
	Provider         string                       `json:"provider"`
	Name             string                       `json:"name"`
	Mode             modelruntimeentities.LLMMode `json:"mode"`
	CompletionParams map[string]any               `json:"completion_params"`
}
type ContextConfig struct {
	Enabled          bool     `json:"enabled"`
	VariableSelector []string `json:"variable_selector"`
}
type VisionConfigOptions struct {
	VariableSelector []string                                             `json:"variable_selector"` //= Field(default_factory=lambda: ["sys", "files"])
	Detail           modelruntimeentities.ImagePromptMessageContentDETAIL `json:"detail"`            //: ImagePromptMessageContent.DETAIL = ImagePromptMessageContent.DETAIL.HIGH
}

func NewVisionConfigOptions() *VisionConfigOptions {
	return &VisionConfigOptions{
		VariableSelector: []string{"sys", "files"},
		Detail:           modelruntimeentities.ImagePromptMessageContentDETAIL_HIGH,
	}
}

type VisionConfig struct {
	Enabled bool                `json:"enabled"`
	Configs VisionConfigOptions `json:"configs"`
}

func NewVisionConfig() *VisionConfig {
	return &VisionConfig{
		Configs: VisionConfigOptions{
			VariableSelector: []string{"sys", "files"},
			Detail:           modelruntimeentities.ImagePromptMessageContentDETAIL_HIGH,
		},
	}
}

type PromptConfig struct {
	Jinja2Variables []*workflowentities.VariableSelector `json:"jinja2_variables"`
}
type LLMNodeChatModelMessage struct {
	*promptentities.ChatModelMessage
	Jinja2Text string `json:"jinja2_text"`
}
type LLMNodeCompletionModelPromptTemplate struct {
	*promptentities.CompletionModelPromptTemplate
	Jinja2Text string `json:"jinja2_text"`
}

// LLMNodeData represents answer node data
type LLMNodeData struct {
	*basenodesentities.BaseNodeData
	Model          *ModelConfig `json:"model"`
	PromptTemplate any/*[]*LLMNodeChatModelMessage | *LLMNodeCompletionModelPromptTemplate*/ `json:"prompt_template"`
	PromptConfig   *PromptConfig                `json:"prompt_config"`
	Memory         *promptentities.MemoryConfig `json:"memory"`
	Context        *ContextConfig               `json:"context"`
	Vision         *VisionConfig                `json:"vision"`
}

func New() *LLMNodeData {
	return &LLMNodeData{
		Vision: NewVisionConfig(),
	}
}

func (data *LLMNodeData) Marshal(config map[string]any, dest any) error {
	mlog.Debugf("---llm config=%#v", config)
	if config == nil {
		return exceptions.NewValueError("ivnalid config")
	}
	bindata, _ := json.Marshal(config)
	if dest != data {
		err := json.Unmarshal(bindata, dest)
		if err != nil {
			return err
		}
		return nil
	}
	err := json.Unmarshal(bindata, data)
	if err != nil {
		return err
	}
	if _, ok := data.PromptTemplate.([]any); ok {
		var tmplist []*LLMNodeChatModelMessage
		bindata, _ := json.Marshal(data.PromptTemplate)
		err := json.Unmarshal(bindata, &tmplist)
		if err != nil {
			mlog.Warningf("data.PromptTemplate=%#v must be LLMNodeChatModelMessage list", data.PromptTemplate)
		} else {
			data.PromptTemplate = tmplist
		}
	} else if _, ok := data.PromptTemplate.(map[string]any); ok {
		tmp := new(LLMNodeCompletionModelPromptTemplate)
		bindata, _ := json.Marshal(data.PromptTemplate)
		err := json.Unmarshal(bindata, tmp)
		if err != nil {
			mlog.Warningf("data.PromptTemplate=%#v must be LLMNodeCompletionModelPromptTemplate", data.PromptTemplate)
		} else {
			data.PromptTemplate = tmp
		}
	}
	mlog.Debugf("---data=%#v", data)
	return nil
}
