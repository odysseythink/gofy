package tools

import (
	"encoding/json"

	toolsenumtypes "mlib.com/gofy/server/enum_types/tools"
	commontypes "mlib.com/gofy/server/types/common"
)

type ToolApiEntity[T1 float64 | int | string, T2 float64 | int] struct {
	Author       string                   `json:"author"`
	Name         string                   `json:"name"`  // identifier
	Label        commontypes.I18nObject   `json:"label"` // label
	Description  commontypes.I18nObject   `json:"description"`
	Parameters   []*ToolParameter[T1, T2] `json:"parameters"`
	Labels       []string                 `json:"labels"`
	OutputSchema map[string]any           `json:"output_schema"`
}

var (
	ToolProviderTypeApiLiteral = []string{"builtin", "api", "workflow", "mcp"}
)

type ToolProviderApiEntity[T1 string | map[string]any, T2 float64 | int | string, T3 float64 | int] struct {
	ID                     string                          `json:"id"`
	Author                 string                          `json:"author"`
	Name                   string                          `json:"name"` // identifier
	Description            commontypes.I18nObject          `json:"description"`
	Icon                   T1                              `json:"icon"`
	IconDark               T1                              `json:"icon_dark"` //description="The dark icon of the tool"
	Label                  commontypes.I18nObject          `json:"label"`     // label
	Type                   toolsenumtypes.ToolProviderType `json:"type"`
	MaskedCredentials      map[string]any                  `json:"masked_credentials"`
	OriginalCredentials    map[string]any                  `json:"original_credentials"`
	IsTeamAuthorization    bool                            `json:"is_team_authorization"`
	AllowDelete            bool                            `json:"allow_delete"`             //default true
	PluginID               *string                         `json:"plugin_id"`                //description="The plugin id of the tool"
	PluginUniqueIdentifier *string                         `json:"plugin_unique_identifier"` //description="The unique identifier of the tool"
	Tools                  []*ToolApiEntity[T2, T3]        `json:"tools"`
	Labels                 []string                        `json:"labels"`
	// MCP
	ServerURL        *string `json:"server_url"`        //description="The server url of the tool"
	UpdatedAt        int     `json:"updated_at"`        //default_factory=lambda: int(datetime.now().timestamp()))
	ServerIdentifier *string `json:"server_identifier"` //description="The server identifier of the MCP tool"
}

func (entity *ToolProviderApiEntity[T1, T2, T3]) ToDict() map[string]any {
	// -------------
	// overwrite tool parameter types for temp fix
	for idx1, tool := range entity.Tools {
		for idx2, parameter := range tool.Parameters {
			if parameter.Type == toolsenumtypes.ToolParameter_SYSTEM_FILES {
				parameter.Type = toolsenumtypes.ToolParameter_FILES
			}
			tool.Parameters[idx2] = parameter
		}
		entity.Tools[idx1] = tool
	}
	bindata, _ := json.Marshal(entity)
	ret := map[string]any{}
	json.Unmarshal(bindata, &ret)
	ret["description"] = entity.Description.ToDict()
	ret["label"] = entity.Label.ToDict()
	// -------------
	return ret
}

type ToolProviderCredentialApiEntity struct {
	ID             string                        `json:"id"`              //"The unique id of the credential")
	Name           string                        `json:"name"`            //"The name of the credential")
	Provider       string                        `json:"provider"`        //"The provider of the credential")
	CredentialType toolsenumtypes.CredentialType `json:"credential_type"` //"The type of the credential")
	IsDefault      bool                          `json:"is_default"`      //"Whether the credential is the default credential for the provider in the workspace"
	Credentials    map[string]any                `json:"credentials"`     //"The credentials of the provider")
}
type ToolProviderCredentialInfoApiEntity struct {
	SupportedCredentialTypes   []string                           `json:"supported_credential_types"`     //"The supported credential types of the provider")
	IsOauthCustomClientEnabled bool                               `json:"is_oauth_custom_client_enabled"` //description="Whether the OAuth custom client is enabled for the provider"
	Credentials                []*ToolProviderCredentialApiEntity `json:"credentials"`                    //"The credentials of the provider")
}

// type UserToolProviderCredentials struct {
// 	Credentials map[string]*ToolProviderCredentials
// }
