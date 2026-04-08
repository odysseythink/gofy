package mcp

import (
	"encoding/json"
	"errors"

	"mlib.com/mlog"
)

type ReferenceType string

const (
	Reference_RESOURCE ReferenceType = "ref/resource"
	Reference_PROMPT   ReferenceType = "ref/prompt"
)

type Referencer interface {
	Type() ReferenceType
	ToDict() map[string]any
}

type ResourceReference struct {
	/*A reference to a resource or resource template definition.*/

	// Type string `json:"type"` //"ref/resource"]
	URI string `json:"uri"`
	/*The URI or URI template of the resource.*/
	ModelConfig map[string]any `json:"model_config"`
}

func (ref *ResourceReference) Type() ReferenceType {
	return Reference_RESOURCE
}
func (ref *ResourceReference) ToDict() map[string]any {
	return map[string]any{
		"type":         ref.Type(),
		"uri":          ref.URI,
		"model_config": ref.ModelConfig,
	}
}
func (ref ResourceReference) MarshalJSON() ([]byte, error) {
	return json.Marshal((&ref).ToDict())
}

func (ref *ResourceReference) UnmarshalJSON(data []byte) error {
	var basedata struct {
		Type ReferenceType `json:"type"`
	}
	if err := json.Unmarshal(data, &basedata); err != nil {
		mlog.Errorf("data=%s must contain type field:%v", string(data), err)
		return err
	}
	if basedata.Type != Reference_RESOURCE {
		mlog.Errorf("data=%s type field must be ref/resource", string(data))
		errors.New("data type field must be ref/resource")
	}

	type alias ResourceReference
	aux := &struct {
		*alias
	}{
		alias: (*alias)(ref), // 把原始对象嵌进去，避免递归
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	return nil
}

type PromptReference struct {
	/*Identifies a prompt.*/

	// Type string `json:"type"` //["ref/prompt"]
	Name string `json:"name"`
	/*The name of the prompt or prompt template*/
	ModelConfig map[string]any `json:"model_config"`
}

func (ref *PromptReference) Type() ReferenceType {
	return Reference_RESOURCE
}
func (ref *PromptReference) ToDict() map[string]any {
	return map[string]any{
		"type":         ref.Type(),
		"name":         ref.Name,
		"model_config": ref.ModelConfig,
	}
}
func (ref PromptReference) MarshalJSON() ([]byte, error) {
	return json.Marshal((&ref).ToDict())
}

func (ref *PromptReference) UnmarshalJSON(data []byte) error {
	var basedata struct {
		Type ReferenceType `json:"type"`
	}
	if err := json.Unmarshal(data, &basedata); err != nil {
		mlog.Errorf("data=%s must contain type field:%v", string(data), err)
		return err
	}
	if basedata.Type != Reference_PROMPT {
		mlog.Errorf("data=%s type field must be ref/prompt", string(data))
		errors.New("data type field must be ref/prompt")
	}

	type alias PromptReference
	aux := &struct {
		*alias
	}{
		alias: (*alias)(ref), // 把原始对象嵌进去，避免递归
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	return nil
}
