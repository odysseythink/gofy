package mcp

import (
	"encoding/json"
)

type ResourceContentType string

const (
	ResourceContent_TEXT ResourceContentType = "text"
	ResourceContent_BLOB ResourceContentType = "blob"
)

type ResourceContenter interface {
	Type() ResourceContentType
	GetURL() string
	SetURL(string)
	GetMimeType() string
	SetMimeType(string)
	GetModelConfig() map[string]any
	SetModelConfig(map[string]any)
	ToDict() map[string]any
}

/*
	URI
	MimeType
	ModelConfig
*/

type BaseResourceContent struct {
	/*The contents of a specific resource or sub-resource.*/

	URI string `json:"uri"` //[AnyUrl, UrlConstraints(host_required=False)]
	/*The URI of this resource.*/
	MimeType string `json:"mimeType"`
	/*The MIME type of this resource, if known.*/
	ModelConfig map[string]any `json:"model_config"`
}

func (content *BaseResourceContent) GetURI() string {
	return content.URI
}
func (content *BaseResourceContent) SetURI(val string) {
	content.URI = val
}
func (content *BaseResourceContent) GetMimeType() string {
	return content.MimeType
}
func (content *BaseResourceContent) SetMimeType(val string) {
	content.MimeType = val
}
func (content *BaseResourceContent) GetModelConfig() map[string]any {
	return content.ModelConfig
}
func (content *BaseResourceContent) SetModelConfig(val map[string]any) {
	content.ModelConfig = val
}

type TextResourceContent struct {
	*BaseResourceContent
	/*Text contents of a resource.*/

	Text string `json:"text"`
	/*
	   The text of the item. This must only be set if the item can actually be represented
	   as text (not binary data).
	*/
}

func (content *TextResourceContent) Type() ResourceContentType {
	return ResourceContent_TEXT
}
func (content *TextResourceContent) ToDict() map[string]any {
	bindata, _ := json.Marshal(content)
	var dict map[string]any
	json.Unmarshal(bindata, &dict)
	dict["type"] = content.Type()
	return dict
}

type BlobResourceContents struct {
	*BaseResourceContent
	/*Binary contents of a resource.*/

	Blob string `json:"blob"`
	/*A base64-encoded string representing the binary data of the item.*/
}

func (content *BlobResourceContents) Type() ResourceContentType {
	return ResourceContent_TEXT
}
func (content *BlobResourceContents) ToDict() map[string]any {
	bindata, _ := json.Marshal(content)
	var dict map[string]any
	json.Unmarshal(bindata, &dict)
	dict["type"] = content.Type()
	return dict
}
