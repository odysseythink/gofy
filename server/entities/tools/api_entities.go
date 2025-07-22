package tools

import (
	"encoding/json"

	toolsenumtypes "mlib.com/gofy/server/enum_types/tools"
	commontypes "mlib.com/gofy/server/types/common"
)

// 定义UserTool结构体
type UserTool struct {
	Author      string                 `json:"author"`
	Name        string                 `json:"name"`
	Label       commontypes.I18nObject `json:"label"`
	Description commontypes.I18nObject `json:"description"`
	Parameters  []*ToolParameter       `json:"parameters,omitempty"`
	Labels      []string               `json:"labels,omitempty"`
}

type UserToolProvider struct {
	ID                  string                          `json:"id"`
	Author              string                          `json:"author"`
	Name                string                          `json:"name"`
	Description         commontypes.I18nObject          `json:"description"`
	Icon                string                          `json:"icon"`
	Label               commontypes.I18nObject          `json:"label"`
	Type                toolsenumtypes.ToolProviderType `json:"type"`
	MaskedCredentials   map[string]string               `json:"masked_credentials,omitempty"`
	OriginalCredentials map[string]string               `json:"original_credentials,omitempty"`
	IsTeamAuthorization bool                            `json:"is_team_authorization"`
	AllowDelete         bool                            `json:"allow_delete"`
	Tools               []*UserTool                     `json:"tools"`
	Labels              []string                        `json:"labels"`
}

func NewUserToolProvider() *UserToolProvider {
	return &UserToolProvider{
		AllowDelete: true,
	}
}

// ToDict 将 UserToolProvider 转换为字典
func (utp *UserToolProvider) ToDict() map[string]any {
	dict := map[string]any{
		"id":                    utp.ID,
		"author":                utp.Author,
		"name":                  utp.Name,
		"description":           utp.Description, // 假设I18nObject有一个ToDict方法
		"icon":                  utp.Icon,
		"label":                 utp.Label, // 假设I18nObject有一个ToDict方法
		"type":                  utp.Type,  // 假设ToolProviderType有一个Value方法
		"team_credentials":      utp.MaskedCredentials,
		"is_team_authorization": utp.IsTeamAuthorization,
		"allow_delete":          utp.AllowDelete,
		"labels":                utp.Labels,
	}

	// 序列化Tools字段
	toolsBytes, _ := json.Marshal(utp.Tools)
	var tools []map[string]any
	json.Unmarshal(toolsBytes, &tools)

	// 这里需要对tools进行处理，类似于Python代码中的逻辑
	for idx, v := range tools {
		if _, ok := v["parameters"]; ok {
			parameters := v["parameters"].([]any)
			for idx1, v1 := range parameters {
				if v1.(map[string]any)["type"].(string) == string(toolsenumtypes.ToolParameter_FILE) {
					v1.(map[string]any)["type"] = "files"
					parameters[idx1] = v1
					v["parameters"] = v1
					tools[idx] = v
				}
			}
		}
	}

	dict["tools"] = tools
	return dict
}

type UserToolProviderCredentials struct {
	Credentials map[string]*ToolProviderCredentials
}
