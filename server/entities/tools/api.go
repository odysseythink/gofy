package tools

import (
	"encoding/json"
	"fmt"

	"mlib.com/gofy/server/core/exceptions"
	pluginparameter "mlib.com/gofy/server/entities/plugin/parameter"
	toolsenumtypes "mlib.com/gofy/server/enum_types/tools"
	commontypes "mlib.com/gofy/server/types/common"
	"mlib.com/mlog"
)

type ToolApiEntity struct{
    Author string `json:"author"`
    Name string `json:"name"`  // identifier
    Label commontypes.I18nObject `json:"label"`  // label
    Description commontypes.I18nObject `json:"description"`
    Parameters []*ToolParameter `json:"parameters"`
    Labels []string `json:"labels"`
    OutputSchema map[string]any `json:"output_schema"`
}
var (
ToolProviderTypeApiLiteral = []string{"builtin", "api", "workflow", "mcp"}
)


type ToolProviderApiEntity[T string|map[string]any] struct{
    ID string `json:"id"`
    Author string `json:"author"`
    Name string `json:"name"`  // identifier
    Description commontypes.I18nObject `json:"description"`
    Icon T `json:"icon"` 
    IconDark T `json:"icon_dark"`//description="The dark icon of the tool"
    Label commontypes.I18nObject `json:"label"`  // label
    Type toolsenumtypes.ToolProviderType `json:"type"`
    MaskedCredentials map[string]any `json:"masked_credentials"`
    OriginalCredentials map[string]any `json:"original_credentials"`
    IsTeamAuthorization bool  `json:"is_team_authorization"`
    AllowDelete bool  `json:"allow_delete"` //default true
    PluginID *string `json:"plugin_id"` //description="The plugin id of the tool"
    PluginUniqueIdentifier *string `json:"plugin_unique_identifier"` //description="The unique identifier of the tool"
    Tools []*ToolApiEntity `json:"tools"`
    Labels []string`json:"labels"`
    // MCP
    ServerURL *string `json:"server_url"` //description="The server url of the tool"
    UpdatedAt int `json:"updated_at"`//default_factory=lambda: int(datetime.now().timestamp()))
    ServerIdentifier *string `json:"server_identifier"` //description="The server identifier of the MCP tool"
}

    func (entity *ToolProviderApiEntity) ToDict() map[string]any{
        // -------------
        // overwrite tool parameter types for temp fix
        for idx1, tool := range  entity.Tools{
			for idx2, parameter := range tool.Parameters{
				if parameter.Type == toolsenumtypes.ToolParameter_SYSTEM_FILES{
					parameter.Type = toolsenumtypes.ToolParameter_FILES
				}
				if parameter.InputSchema == nil{
					parameter.pop("input_schema", None)
				}
			}
		}
        // -------------
        optional_fields = entity.optional_field("server_url", entity.server_url)
        if entity.type == ToolProviderType.MCP.value:
            optional_fields.update(entity.optional_field("updated_at", entity.updated_at))
            optional_fields.update(entity.optional_field("server_identifier", entity.server_identifier))
        return {
            "id": entity.id,
            "author": entity.author,
            "name": entity.name,
            "plugin_id": entity.plugin_id,
            "plugin_unique_identifier": entity.plugin_unique_identifier,
            "description": entity.description.to_dict(),
            "icon": entity.icon,
            "icon_dark": entity.icon_dark,
            "label": entity.label.to_dict(),
            "type": entity.type.value,
            "team_credentials": entity.masked_credentials,
            "is_team_authorization": entity.is_team_authorization,
            "allow_delete": entity.allow_delete,
            "tools": tools,
            "labels": entity.labels,
            **optional_fields,
        }
}
        func (entity *ToolProviderApiEntity) optional_field(key: str, value: Any) -> dict:
        """Return dict with key-value if value is truthy, empty dict otherwise."""
        return {key: value} if value else {}
}

type ToolProviderCredentialApiEntity struct{
    id: str = Field(description="The unique id of the credential")
    name: str = Field(description="The name of the credential")
    provider: str = Field(description="The provider of the credential")
    credential_type: CredentialType = Field(description="The type of the credential")
    is_default: bool = Field(
        default=False, description="Whether the credential is the default credential for the provider in the workspace"
    )
    credentials: dict = Field(description="The credentials of the provider")


}
type ToolProviderCredentialInfoApiEntity struct{
    supported_credential_types: list[str] = Field(description="The supported credential types of the provider")
    is_oauth_custom_client_enabled: bool = Field(
        default=False, description="Whether the OAuth custom client is enabled for the provider"
    )
    credentials: list[ToolProviderCredentialApiEntity] = Field(description="The credentials of the provider")
}
type UserToolProviderCredentials struct {
	Credentials map[string]*ToolProviderCredentials
}
