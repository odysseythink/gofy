package models

import (
	"encoding/json"
	"time"

	uuid "github.com/satori/go.uuid"
	toolsentities "mlib.com/gofy/server/entities/tools"
	toolsenumtypes "mlib.com/gofy/server/enum_types/tools"
	commontypes "mlib.com/gofy/server/types/common"
	"mlib.com/mlog"
)

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

func (self *BuiltinToolProvider) Credentials() map[string]any {
	tmp := map[string]any{}
	err := json.Unmarshal([]byte(self.EncryptedCredentials), &tmp)
	if err != nil {
		mlog.Errorf("json unmarshal failed:%v", err)
		return nil
	} else {
		return tmp
	}
}

// PublishedAppTool [...]
type PublishedAppTool struct {
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
func (PublishedAppTool) TableName() string {
	return "tool_published_apps"
}
func (self *PublishedAppTool) Description2I18n() *commontypes.I18nObject {
	e := new(commontypes.I18nObject)
	err := json.Unmarshal([]byte(self.Description), e)
	if err != nil {
		mlog.Errorf("json unmarshal failed:%v", err)
		return nil
	} else {
		return e
	}
}

// ApiToolProvider [...]
type ApiToolProvider struct {
	ID               string                               `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	Name             string                               `gorm:"column:name;type:varchar(40);not null" json:"name"`
	Schema           string                               `gorm:"column:schema;type:text;not null" json:"schema"`
	SchemaType       toolsenumtypes.ApiProviderSchemaType `gorm:"column:schema_type_str;type:varchar(40);not null" json:"schema_type_str"`
	UserID           string                               `gorm:"column:user_id;type:varchar(36);not null" json:"user_id"`
	User             *Account                             `json:"user" form:"user" gorm:"foreignKey:UserID;references:ID;"`
	TenantID         string                               `gorm:"column:tenant_id;type:varchar(36);not null" json:"tenant_id"`
	Tenant           *Tenant                              `json:"tenant" form:"tenant" gorm:"foreignKey:TenantID;references:ID;"`
	ToolsStr         string                               `gorm:"column:tools_str;type:text;not null" json:"tools_str"`
	Icon             string                               `gorm:"column:icon;type:varchar(255);not null" json:"icon"`
	CredentialsStr   string                               `gorm:"column:credentials_str;type:text;not null" json:"credentials_str"`
	Description      string                               `gorm:"column:description;type:text;not null" json:"description"`
	CreatedAt        *time.Time                           `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        *time.Time                           `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	PrivacyPolicy    string                               `gorm:"column:privacy_policy;type:varchar(255)" json:"privacy_policy"`
	CustomDisclaimer string                               `gorm:"column:custom_disclaimer;type:text;not null" json:"custom_disclaimer"`
}

// TableName get sql table name.获取数据库表名
func (ApiToolProvider) TableName() string {
	return "tool_api_providers"
}

func (self *ApiToolProvider) Tools() []*toolsentities.ApiToolBundle {
	if self.ToolsStr != "" {
		var tools []*toolsentities.ApiToolBundle
		err := json.Unmarshal([]byte(self.ToolsStr), &tools)
		if err != nil {
			mlog.Errorf("json unmarshal(%s) failed:%v", self.ToolsStr, err)
			return nil
		}
		return tools
	} else {
		return nil
	}
}
func (self *ApiToolProvider) Credentials() map[string]any {
	if self.CredentialsStr == "" {
		return nil
	} else {
		var tmp map[string]any
		err := json.Unmarshal([]byte(self.CredentialsStr), &tmp)
		if err != nil {
			mlog.Errorf("json unmarshal(%s) failed:%v", self.CredentialsStr, err)
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
	ID                        string     `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	Name                      string     `gorm:"column:name;type:varchar(40);not null" json:"name"`
	Icon                      string     `gorm:"column:icon;type:varchar(255);not null" json:"icon"`
	AppID                     string     `gorm:"column:app_id;type:varchar(36);not null" json:"app_id"`
	App                       *App       `json:"app" form:"app" gorm:"foreignKey:AppID;references:ID;"`
	UserID                    string     `gorm:"column:user_id;type:varchar(36);not null" json:"user_id"`
	User                      *Account   `json:"user" form:"user" gorm:"foreignKey:UserID;references:ID;"`
	TenantID                  string     `gorm:"column:tenant_id;type:varchar(36);not null" json:"tenant_id"`
	Tenant                    *Tenant    `json:"tenant" form:"tenant" gorm:"foreignKey:TenantID;references:ID;"`
	Description               string     `gorm:"column:description;type:text;not null" json:"description"`
	ParameterConfigurationStr string     `gorm:"column:parameter_configuration;type:varchar(255);default:[]" json:"parameter_configuration"`
	CreatedAt                 *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt                 *time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	PrivacyPolicy             string     `gorm:"column:privacy_policy;type:varchar(255);default:''" json:"privacy_policy"`
	Version                   string     `gorm:"column:version;type:varchar(255);default:''" json:"version"`
	Label                     string     `gorm:"column:label;type:varchar(255);default:''" json:"label"`
}

// TableName get sql table name.获取数据库表名
func (WorkflowToolProvider) TableName() string {
	return "tool_workflow_providers"
}
func (self *WorkflowToolProvider) ParameterConfigurations() []*toolsentities.WorkflowToolParameterConfiguration {
	if self.ParameterConfigurationStr == "" {
		return nil
	} else {
		var tmp []*toolsentities.WorkflowToolParameterConfiguration
		err := json.Unmarshal([]byte(self.ParameterConfigurationStr), &tmp)
		if err != nil {
			mlog.Errorf("json unmarshal(%s) failed:%v", self.ParameterConfigurationStr, err)
			return nil
		}
		return tmp
	}
}

// ToolModelInvoke [...]
type ToolModelInvoke struct {
	ID                      string     `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	UserID                  string     `gorm:"column:user_id;type:varchar(36);not null" json:"user_id"`
	TenantID                string     `gorm:"column:tenant_id;type:varchar(36);not null" json:"tenant_id"`
	Provider                string     `gorm:"column:provider;type:varchar(40);not null" json:"provider"`
	ToolType                string     `gorm:"column:tool_type;type:varchar(40);not null" json:"tool_type"`
	ToolName                string     `gorm:"column:tool_name;type:varchar(40);not null" json:"tool_name"`
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

func (self *ToolConversationVariable) Variables() map[string]any {
	if self.VariablesStr == "" {
		return nil
	} else {
		var tmp map[string]any
		err := json.Unmarshal([]byte(self.VariablesStr), &tmp)
		if err != nil {
			mlog.Errorf("json unmarshal(%s) failed:%v", self.VariablesStr, err)
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
	OriginalURL    string `gorm:"column:original_url;type:text" json:"original_url"`
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
