package models

import (
	"encoding/json"
	"net/url"
	"strings"
	"time"

	"github.com/odysseythink/gofy/backend/core/file"
	dbengine "github.com/odysseythink/gofy/backend/db_engine"
	toolsentities "github.com/odysseythink/gofy/backend/entities/tools"
	toolsenumtypes "github.com/odysseythink/gofy/backend/enum_types/tools"
	commontypes "github.com/odysseythink/gofy/backend/types/common"
	mcptypes "github.com/odysseythink/gofy/backend/types/mcp"
	"github.com/odysseythink/mlog"
	uuid "github.com/satori/go.uuid"
)

// system level tool oauth client params (client_id, client_secret, etc.)
type ToolOAuthSystemClient struct {
	ID                   string `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	PluginID             string `gorm:"column:plugin_id;type:varchar(512);not null" json:"plugin_id"`
	Provider             string `gorm:"column:provider;type:varchar(255);not null" json:"provider"`
	EncryptedOauthParams string `gorm:"column:encrypted_oauth_params;type:text" json:"encrypted_oauth_params"`
}

// TableName get sql table name.获取数据库表名
func (ToolOAuthSystemClient) TableName() string {
	return "tool_oauth_system_clients"
}

// tenant level tool oauth client params (client_id, client_secret, etc.)
type ToolOAuthTenantClient struct {
	ID                   string `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	TenantID             string `gorm:"column:tenant_id;type:varchar(36);not null" json:"tenant_id"`
	PluginID             string `gorm:"column:plugin_id;type:varchar(512);not null" json:"plugin_id"`
	Provider             string `gorm:"column:provider;type:varchar(255);not null" json:"provider"`
	Enabled              bool   `gorm:"column:enabled;type:tinyint(1);not null;default:1" json:"enabled"`
	EncryptedOauthParams string `gorm:"column:encrypted_oauth_params;type:text;not null" json:"encrypted_oauth_params"`
}

// TableName get sql table name.获取数据库表名
func (ToolOAuthTenantClient) TableName() string {
	return "tool_oauth_tenant_clients"
}
func (totc *ToolOAuthTenantClient) OauthParams() map[string]any {
	if totc.EncryptedOauthParams != "" {
		var dict map[string]any
		err := json.Unmarshal([]byte(totc.EncryptedOauthParams), &dict)
		if err != nil {
			mlog.Errorf("json unmarshal %s to dict failed:%v", totc.EncryptedOauthParams, err)
			return map[string]any{}
		}
		return dict
	} else {
		return map[string]any{}
	}
}

// BuiltinToolProvider [...]
type BuiltinToolProvider struct {
	ID                   string     `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	TenantID             string     `gorm:"column:tenant_id;type:varchar(36)" json:"tenant_id"`
	UserID               string     `gorm:"column:user_id;type:varchar(36);not null" json:"user_id"`
	Provider             string     `gorm:"column:provider;type:varchar(40);not null" json:"provider"`
	EncryptedCredentials string     `gorm:"column:encrypted_credentials;type:text" json:"encrypted_credentials"`
	CreatedAt            *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt            *time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName get sql table name.获取数据库表名
func (BuiltinToolProvider) TableName() string {
	return "tool_builtin_providers"
}

func (btp *BuiltinToolProvider) Credentials() map[string]any {
	tmp := map[string]any{}
	err := json.Unmarshal([]byte(btp.EncryptedCredentials), &tmp)
	if err != nil {
		mlog.Errorf("json unmarshal failed:%v", err)
		return nil
	} else {
		return tmp
	}
}

// ApiToolProvider [...]
type ApiToolProvider struct {
	ID               string     `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	Name             string     `gorm:"column:name;type:varchar(40);not null" json:"name"`
	Icon             string     `gorm:"column:icon;type:varchar(255);not null" json:"icon"`
	Schema           string     `gorm:"column:schema;type:text;not null" json:"schema"`
	SchemaTypeStr    string     `gorm:"column:schema_type_str;type:varchar(40);not null" json:"schema_type_str"`
	UserID           string     `gorm:"column:user_id;type:varchar(36);not null" json:"user_id"`
	TenantID         string     `gorm:"column:tenant_id;type:varchar(36);not null" json:"tenant_id"`
	Description      string     `gorm:"column:description;type:text;not null" json:"description"`
	ToolsStr         string     `gorm:"column:tools_str;type:text;not null" json:"tools_str"`
	CredentialsStr   string     `gorm:"column:credentials_str;type:text;not null" json:"credentials_str"`
	PrivacyPolicy    string     `gorm:"column:privacy_policy;type:varchar(255)" json:"privacy_policy"`
	CustomDisclaimer string     `gorm:"column:custom_disclaimer;type:text;not null" json:"custom_disclaimer"`
	CreatedAt        *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        *time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName get sql table name.获取数据库表名
func (ApiToolProvider) TableName() string {
	return "tool_api_providers"
}

func (atp *ApiToolProvider) SchemaType() toolsenumtypes.ApiProviderSchemaType {
	return toolsenumtypes.ApiProviderSchemaType(atp.SchemaTypeStr)
}

func (atp *ApiToolProvider) User() *Account {
	if atp.UserID == "" {
		return nil
	}
	user := new(Account)
	err := dbengine.Instance().DB.Model(&Account{}).Where("id = ?", atp.UserID).First(user).Error
	if err != nil {
		mlog.Errorf("get Account(%s) failed:%v", atp.UserID, err)
		return nil
	}
	return user
}
func (atp *ApiToolProvider) Tenant() *Tenant {
	if atp.TenantID == "" {
		return nil
	}
	tenant := new(Tenant)
	err := dbengine.Instance().DB.Model(&Tenant{}).Where("id = ?", atp.TenantID).First(tenant).Error
	if err != nil {
		mlog.Errorf("get Tenant(%s) failed:%v", atp.TenantID, err)
		return nil
	}
	return tenant
}

func (atp *ApiToolProvider) Tools() []*toolsentities.ApiToolBundle {
	if atp.ToolsStr != "" {
		var tools []*toolsentities.ApiToolBundle
		err := json.Unmarshal([]byte(atp.ToolsStr), &tools)
		if err != nil {
			mlog.Errorf("json unmarshal(%s) failed:%v", atp.ToolsStr, err)
			return nil
		}
		return tools
	} else {
		return nil
	}
}
func (atp *ApiToolProvider) Credentials() map[string]any {
	if atp.CredentialsStr == "" {
		return nil
	} else {
		var tmp map[string]any
		err := json.Unmarshal([]byte(atp.CredentialsStr), &tmp)
		if err != nil {
			mlog.Errorf("json unmarshal(%s) failed:%v", atp.CredentialsStr, err)
			return nil
		}
		return tmp
	}
}

// ToolLabelBinding [...]
type ToolLabelBinding struct {
	ID        string `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	ToolID    string `gorm:"column:tool_id;type:varchar(64);not null" json:"tool_id"`
	ToolType  string `gorm:"column:tool_type;type:varchar(40);not null" json:"tool_type"`
	LabelName string `gorm:"column:label_name;type:varchar(40);not null" json:"label_name"`
}

// TableName get sql table name.获取数据库表名
func (ToolLabelBinding) TableName() string {
	return "tool_label_bindings"
}

// WorkflowToolProvider [...]
type WorkflowToolProvider struct {
	ID                     string     `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	Name                   string     `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Label                  string     `gorm:"column:label;type:varchar(255);default:''" json:"label"`
	Icon                   string     `gorm:"column:icon;type:varchar(255);not null" json:"icon"`
	AppID                  string     `gorm:"column:app_id;type:varchar(36);not null" json:"app_id"`
	Version                string     `gorm:"column:version;type:varchar(255);default:''" json:"version"`
	UserID                 string     `gorm:"column:user_id;type:varchar(36);not null" json:"user_id"`
	TenantID               string     `gorm:"column:tenant_id;type:varchar(36);not null" json:"tenant_id"`
	Description            string     `gorm:"column:description;type:text;not null" json:"description"`
	ParameterConfiguration string     `gorm:"column:parameter_configuration;type:varchar(255);default:[]" json:"parameter_configuration"`
	PrivacyPolicy          string     `gorm:"column:privacy_policy;type:varchar(255);default:''" json:"privacy_policy"`
	CreatedAt              *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt              *time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName get sql table name.获取数据库表名
func (WorkflowToolProvider) TableName() string {
	return "tool_workflow_providers"
}
func (wftp *WorkflowToolProvider) ParameterConfigurations() []*toolsentities.WorkflowToolParameterConfiguration {
	if wftp.ParameterConfiguration == "" {
		return nil
	} else {
		var tmp []*toolsentities.WorkflowToolParameterConfiguration
		err := json.Unmarshal([]byte(wftp.ParameterConfiguration), &tmp)
		if err != nil {
			mlog.Errorf("json unmarshal(%s) failed:%v", wftp.ParameterConfiguration, err)
			return nil
		}
		return tmp
	}
}

func (wftp *WorkflowToolProvider) User() *Account {
	if wftp.UserID == "" {
		return nil
	}
	user := new(Account)
	err := dbengine.Instance().DB.Model(&Account{}).Where("id = ?", wftp.UserID).First(user).Error
	if err != nil {
		mlog.Errorf("get Account(%s) failed:%v", wftp.UserID, err)
		return nil
	}
	return user
}
func (wftp *WorkflowToolProvider) Tenant() *Tenant {
	if wftp.TenantID == "" {
		return nil
	}
	tenant := new(Tenant)
	err := dbengine.Instance().DB.Model(&Tenant{}).Where("id = ?", wftp.TenantID).First(tenant).Error
	if err != nil {
		mlog.Errorf("get Tenant(%s) failed:%v", wftp.TenantID, err)
		return nil
	}
	return tenant
}

func (wftp *WorkflowToolProvider) App() *App {
	if wftp.AppID == "" {
		return nil
	}
	app := new(App)
	err := dbengine.Instance().DB.Model(&App{}).Where("id = ?", wftp.AppID).First(app).Error
	if err != nil {
		mlog.Errorf("get App(%s) failed:%v", wftp.AppID, err)
		return nil
	}
	return app
}

type MCPToolProvider struct {
	ID                   string     `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	Name                 string     `gorm:"column:name;type:varchar(40);not null" json:"name"`
	ServerIdentifier     string     `gorm:"column:server_identifier;type:varchar(64);not null" json:"server_identifier"`
	ServerURL            string     `gorm:"column:server_url;not null" json:"server_url"`
	ServerURLHash        string     `gorm:"column:server_url_hash;type:varchar(64);not null" json:"server_url_hash"`
	Icon                 string     `gorm:"column:icon" json:"icon"`
	TenantID             string     `gorm:"column:tenant_id;type:varchar(36);not null" json:"tenant_id"`
	UserID               string     `gorm:"column:user_id;type:varchar(36);not null" json:"user_id"`
	EncryptedCredentials string     `gorm:"column:encrypted_credentials" json:"encrypted_credentials"`
	Authed               bool       `gorm:"column:authed;type:tinyint(0);not null;default:1" json:"authed"`
	Tools                string     `gorm:"column:tools;type:text;default:'[]'" json:"tools"`
	CreatedAt            *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt            *time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (MCPToolProvider) TableName() string {
	return "tool_mcp_providers"
}
func (mcptp *MCPToolProvider) LoadUser() *Account {
	if mcptp.UserID == "" {
		return nil
	}
	user := new(Account)
	err := dbengine.Instance().DB.Model(&Account{}).Where("id = ?", mcptp.UserID).First(user).Error
	if err != nil {
		mlog.Errorf("get Account(%s) failed:%v", mcptp.UserID, err)
		return nil
	}
	return user
}

func (mcptp *MCPToolProvider) Tenant() *Tenant {
	if mcptp.TenantID == "" {
		return nil
	}
	tenant := new(Tenant)
	err := dbengine.Instance().DB.Model(&Tenant{}).Where("id = ?", mcptp.TenantID).First(tenant).Error
	if err != nil {
		mlog.Errorf("get Tenant(%s) failed:%v", mcptp.TenantID, err)
		return nil
	}
	return tenant
}

func (mcptp *MCPToolProvider) Credentials() map[string]any {
	if mcptp.EncryptedCredentials == "" {
		return nil
	} else {
		var tmp map[string]any
		err := json.Unmarshal([]byte(mcptp.EncryptedCredentials), &tmp)
		if err != nil {
			mlog.Errorf("json unmarshal(%s) failed:%v", mcptp.EncryptedCredentials, err)
			return nil
		}
		return tmp
	}
}

func (mcptp *MCPToolProvider) MCPTools() []*mcptypes.Tool {

	if mcptp.Tools != "" {
		var datas []*mcptypes.Tool
		err := json.Unmarshal([]byte(mcptp.Tools), &datas)
		if err != nil {
			mlog.Errorf("json unmashal %s to tool list failed:%v", mcptp.Tools, err)
			return nil
		}
		return datas
	}
	return nil
}
func (mcptp *MCPToolProvider) ProviderIcon() any /*map[string]string | string*/ {

	if mcptp.Icon != "" {
		var datas map[string]string
		err := json.Unmarshal([]byte(mcptp.Icon), &datas)
		if err != nil {
			mlog.Warningf("json unmashal %s to map[string]string failed:%v", mcptp.Icon, err)
			return file.GetSignedFileURL(mcptp.Icon)
		}
		return datas
	}
	return nil
}

func (mcptp *MCPToolProvider) DecryptedServerURL() string {
	return mcptp.ServerURL
	// return cast(str, encrypter.DecryptToken(self.tenant_id, self.server_url))
}
func mask_url(rawURL string, mask_char string /* = "*"*/) string {
	/*
	   mask the url to a simple string
	*/
	if mask_char == "" {
		mask_char = "*"
	}
	u, err := url.Parse(rawURL)
	if err != nil {
		// 如果解析失败，直接返回原串
		return rawURL
	}
	base := u.Scheme + "://" + u.Host
	path := u.EscapedPath()
	if path == "" {
		path = "/"
	}
	if path != "/" {
		return base + "/" + strings.Repeat(mask_char, 6)
	}
	return base
}
func (mcptp *MCPToolProvider) MaskedServerURL() string {
	return mask_url(mcptp.DecryptedServerURL(), "")
}

// func (mcptp *MCPToolProvider)decrypted_credentials() -> dict:
//     from core.helper.provider_cache import NoOpProviderCredentialCache
//     from core.tools.mcp_tool.provider import MCPToolProviderController
//     from core.tools.utils.encryption import create_provider_encrypter

//     provider_controller = MCPToolProviderController._from_db(self)

//     encrypter, _ = create_provider_encrypter(
//         tenant_id=self.tenant_id,
//         config=[x.to_basic_provider_config() for x in provider_controller.get_credentials_schema()],
//         cache=NoOpProviderCredentialCache(),
//     )

//     return encrypter.decrypt(self.credentials)  # type: ignore

// ToolModelInvoke [...]
type ToolModelInvoke struct {
	ID                      string     `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	UserID                  string     `gorm:"column:user_id;type:varchar(36);not null" json:"user_id"`
	TenantID                string     `gorm:"column:tenant_id;type:varchar(36);not null" json:"tenant_id"`
	Provider                string     `gorm:"column:provider;type:varchar(255);not null" json:"provider"`
	ToolType                string     `gorm:"column:tool_type;type:varchar(40);not null" json:"tool_type"`
	ToolName                string     `gorm:"column:tool_name;type:varchar(128);not null" json:"tool_name"`
	ModelParameters         string     `gorm:"column:model_parameters;type:text;not null" json:"model_parameters"`
	PromptMessages          string     `gorm:"column:prompt_messages;type:text;not null" json:"prompt_messages"`
	ModelResponse           string     `gorm:"column:model_response;type:text;not null" json:"model_response"`
	PromptTokens            int        `gorm:"column:prompt_tokens;type:int;not null;default:0" json:"prompt_tokens"`
	AnswerTokens            int        `gorm:"column:answer_tokens;type:int;not null;default:0" json:"answer_tokens"`
	AnswerUnitPrice         float64    `gorm:"column:answer_unit_price;type:decimal(10,4);not null" json:"answer_unit_price"`
	AnswerPriceUnit         float64    `gorm:"column:answer_price_unit;type:decimal(10,7);not null;default:0.0010000" json:"answer_price_unit"`
	ProviderResponseLatency float64    `gorm:"column:provider_response_latency;type:double;not null;default:0" json:"provider_response_latency"`
	TotalPrice              float64    `gorm:"column:total_price;type:decimal(10,7)" json:"total_price"`
	Currency                string     `gorm:"column:currency;type:varchar(255);not null" json:"currency"`
	CreatedAt               *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt               *time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName get sql table name.获取数据库表名
func (ToolModelInvoke) TableName() string {
	return "tool_model_invokes"
}

// ToolConversationVariable [...]
type ToolConversationVariable struct {
	ID             string     `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	UserID         string     `gorm:"column:user_id;type:varchar(36);not null" json:"user_id"`
	TenantID       string     `gorm:"column:tenant_id;type:varchar(36);not null" json:"tenant_id"`
	ConversationID string     `gorm:"column:conversation_id;type:varchar(36);not null" json:"conversation_id"`
	VariablesStr   string     `gorm:"column:variables_str;type:text;not null" json:"variables_str"`
	CreatedAt      *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      *time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName get sql table name.获取数据库表名
func (ToolConversationVariable) TableName() string {
	return "tool_conversation_variables"
}

func (tcv *ToolConversationVariable) Variables() map[string]any {
	if tcv.VariablesStr == "" {
		return nil
	} else {
		var tmp map[string]any
		err := json.Unmarshal([]byte(tcv.VariablesStr), &tmp)
		if err != nil {
			mlog.Errorf("json unmarshal(%s) failed:%v", tcv.VariablesStr, err)
			return nil
		}
		return tmp
	}
}

// ToolFile [...]
type ToolFile struct {
	ID             string `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	UserID         string `gorm:"column:user_id;type:varchar(36);not null" json:"user_id"`
	TenantID       string `gorm:"column:tenant_id;type:varchar(36);not null" json:"tenant_id"`
	ConversationID string `gorm:"column:conversation_id;type:varchar(36)" json:"conversation_id"`
	FileKey        string `gorm:"column:file_key;type:varchar(255);not null" json:"file_key"`
	Mimetype       string `gorm:"column:mimetype;type:varchar(255);not null" json:"mimetype"`
	OriginalURL    string `gorm:"column:original_url;type:varchar(2048)" json:"original_url"`
	Name           string `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Size           int    `gorm:"column:size;type:int;not null" json:"size"`
}

func NewToolFile(
	user_id string,
	tenant_id string,
	conversation_id string,
	file_key string,
	mimetype string,
	original_url string,
	name string,
	size int,
) *ToolFile {
	return &ToolFile{
		ID:             uuid.NewV4().String(),
		UserID:         user_id,
		TenantID:       tenant_id,
		ConversationID: conversation_id,
		FileKey:        file_key,
		Mimetype:       mimetype,
		OriginalURL:    original_url,
		Name:           name,
		Size:           size,
	}
}

// TableName get sql table name.获取数据库表名
func (ToolFile) TableName() string {
	return "tool_files"
}

// DeprecatedPublishedAppTool [...]
type DeprecatedPublishedAppTool struct {
	ID               string     `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	AppID            string     `gorm:"column:app_id;type:varchar(36);not null" json:"app_id"`
	App              *App       `json:"app" form:"app" gorm:"foreignKey:AppID;references:ID;"`
	UserID           string     `gorm:"column:user_id;type:varchar(36);not null" json:"user_id"`
	Description      string     `gorm:"column:description;type:text;not null" json:"description"`
	LlmDescription   string     `gorm:"column:llm_description;type:text;not null" json:"llm_description"`
	QueryDescription string     `gorm:"column:query_description;type:text;not null" json:"query_description"`
	QueryName        string     `gorm:"column:query_name;type:varchar(40);not null" json:"query_name"`
	ToolName         string     `gorm:"column:tool_name;type:varchar(40);not null" json:"tool_name"`
	Author           string     `gorm:"column:author;type:varchar(40);not null" json:"author"`
	CreatedAt        *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        *time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName get sql table name.获取数据库表名
func (DeprecatedPublishedAppTool) TableName() string {
	return "tool_published_apps"
}
func (pat *DeprecatedPublishedAppTool) Description2I18n() *commontypes.I18nObject {
	e := new(commontypes.I18nObject)
	err := json.Unmarshal([]byte(pat.Description), e)
	if err != nil {
		mlog.Errorf("json unmarshal failed:%v", err)
		return nil
	} else {
		return e
	}
}
