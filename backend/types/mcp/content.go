package mcp

import (
	"encoding/json"
	"errors"

	"github.com/odysseythink/mlog"
)

type ContentType string

const (
	Content_TEXT     ContentType = "text"
	Content_IMAGE    ContentType = "image"
	Content_RESOURCE ContentType = "resource"
)

type Contenter interface {
	Type() ContentType
	ToDict() map[string]any
}
type TextContent struct {
	/*Text content for a message.*/

	Text string `json:"text"`
	/*The text content of the message.*/
	Annotations *Annotations   `json:"annotations"`
	ModelConfig map[string]any `json:"model_config"`
}

func (content *TextContent) Type() ContentType {
	return Content_TEXT
}
func (content *TextContent) ToDict() map[string]any {
	if content.Annotations != nil {
		return map[string]any{
			"type":         content.Type(),
			"text":         content.Text,
			"annotations":  content.Annotations.ToDict(),
			"model_config": content.ModelConfig,
		}
	} else {
		return map[string]any{
			"type":         content.Type(),
			"text":         content.Text,
			"annotations":  nil,
			"model_config": content.ModelConfig,
		}
	}
}
func (content TextContent) MarshalJSON() ([]byte, error) {
	return json.Marshal((&content).ToDict())
}

func (content *TextContent) UnmarshalJSON(data []byte) error {
	var basedata struct {
		Type ContentType `json:"type"`
	}
	if err := json.Unmarshal(data, &basedata); err != nil {
		mlog.Errorf("data=%s must contain type field:%v", string(data), err)
		return err
	}
	if basedata.Type != Content_TEXT {
		mlog.Errorf("data=%s type field must be text", string(data))
		errors.New("data type field must be text")
	}

	type alias TextContent
	aux := &struct {
		*alias
	}{
		alias: (*alias)(content), // 把原始对象嵌进去，避免递归
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	return nil
}

type ImageContent struct {
	/*Image content for a message.*/

	Data string `json:"data"`
	/*The base64-encoded image data.*/
	MimeType string `json:"mimeType"`
	/*
	   The MIME type of the image. Different providers may support different
	   image types.
	*/
	Annotations *Annotations   `json:"annotations"`
	ModelConfig map[string]any `json:"model_config"`
}

func (content *ImageContent) Type() ContentType {
	return Content_IMAGE
}
func (content *ImageContent) ToDict() map[string]any {
	if content.Annotations != nil {
		return map[string]any{
			"type":         content.Type(),
			"data":         content.Data,
			"mimeType":     content.MimeType,
			"annotations":  content.Annotations.ToDict(),
			"model_config": content.ModelConfig,
		}
	} else {
		return map[string]any{
			"type":         content.Type(),
			"data":         content.Data,
			"mimeType":     content.MimeType,
			"annotations":  nil,
			"model_config": content.ModelConfig,
		}
	}
}
func (content ImageContent) MarshalJSON() ([]byte, error) {
	return json.Marshal((&content).ToDict())
}

func (content *ImageContent) UnmarshalJSON(data []byte) error {
	var basedata struct {
		Type ContentType `json:"type"`
	}
	if err := json.Unmarshal(data, &basedata); err != nil {
		mlog.Errorf("data=%s must contain type field:%v", string(data), err)
		return err
	}
	if basedata.Type != Content_IMAGE {
		mlog.Errorf("data=%s type field must be image", string(data))
		errors.New("data type field must be image")
	}

	type alias ImageContent
	aux := &struct {
		*alias
	}{
		alias: (*alias)(content), // 把原始对象嵌进去，避免递归
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	return nil
}

type EmbeddedResourceContent struct {
	/*
	   The contents of a resource, embedded into a prompt or tool call result.

	   It is up to the client how best to render embedded resources for the benefit
	   of the LLM and/or the user.
	*/

	// Type string `json:"type"` //Literal["resource"]
	Resource    ResourceContenter `json:"resource"` //TextResourceContents | BlobResourceContents
	Annotations *Annotations      `json:"annotations"`
	ModelConfig map[string]any    `json:"model_config"`
}

func (content *EmbeddedResourceContent) Type() ContentType {
	return Content_TEXT
}
func (content *EmbeddedResourceContent) ToDict() map[string]any {
	ret := map[string]any{
		"type":         content.Type(),
		"model_config": content.ModelConfig,
	}
	if content.Annotations != nil {
		ret["annotations"] = content.Annotations.ToDict()
	} else {
		ret["annotations"] = nil
	}
	if content.Resource != nil {
		ret["resource"] = content.Resource.ToDict()
	} else {
		ret["resource"] = nil
	}
	return ret
}
func (content EmbeddedResourceContent) MarshalJSON() ([]byte, error) {
	return json.Marshal((&content).ToDict())
}

func (content *EmbeddedResourceContent) UnmarshalJSON(data []byte) error {
	var basedata struct {
		Type ContentType `json:"type"`
	}
	if err := json.Unmarshal(data, &basedata); err != nil {
		mlog.Errorf("data=%s must contain type field:%v", string(data), err)
		return err
	}
	if basedata.Type != Content_RESOURCE {
		mlog.Errorf("data=%s type field must be resource", string(data))
		errors.New("data type field must be resource")
	}

	type alias EmbeddedResourceContent
	aux := &struct {
		*alias
	}{
		alias: (*alias)(content), // 把原始对象嵌进去，避免递归
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	return nil
}
