package models

import (
	"time"

	modelruntimeenumtypes "github.com/odysseythink/gofy/backend/enum_types/model_runtime"
	providerenumtypes "github.com/odysseythink/gofy/backend/enum_types/provider"
)

// Provider [...]
type Provider struct {
	ID                             string  `gorm:"column:id;type:varchar(36);not null" json:"id"`
	TenantID                       string  `gorm:"column:tenant_id;type:varchar(36);not null" json:"tenant_id"`
	Tenant                         *Tenant `json:"tenant" form:"tenant" gorm:"foreignKey:TenantID;references:ID;"`
	ProviderName                   string  `gorm:"column:provider_name;type:varchar(255);not null" json:"provider_name"`
	providerenumtypes.ProviderType `gorm:"column:provider_type;type:varchar(40);default:custom" json:"provider_type"`
	EncryptedConfig                string     `gorm:"column:encrypted_config;type:text" json:"encrypted_config"`
	IsValid                        bool       `gorm:"column:is_valid;type:tinyint(1);not null;default:0" json:"is_valid"`
	LastUsed                       *time.Time `gorm:"column:last_used;type:timestamp" json:"last_used"`
	QuotaType                      string     `gorm:"column:quota_type;type:varchar(40);default:''" json:"quota_type"`
	QuotaLimit                     int64      `gorm:"column:quota_limit;type:bigint" json:"quota_limit"`
	QuotaUsed                      int64      `gorm:"column:quota_used;type:bigint" json:"quota_used"`
	CreatedAt                      *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt                      *time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName get sql table name.获取数据库表名
func (Provider) TableName() string {
	return "providers"
}

func (p *Provider) TokenIsSet() bool {
	// """
	// Returns True if the encrypted_config is not None, indicating that the token is set.
	// """
	return p.EncryptedConfig != ""
}
func (p *Provider) IsEnabled() bool {
	// """
	// Returns True if the provider is enabled.
	// """
	if p.ProviderType == providerenumtypes.Provider_SYSTEM {
		return p.IsValid
	} else {
		return p.IsValid && p.TokenIsSet()
	}
}

// ProviderModel [...]
type ProviderModel struct {
	ID              string     `gorm:"column:id;type:varchar(36);not null" json:"id"`
	TenantID        string     `gorm:"column:tenant_id;type:varchar(36);not null" json:"tenant_id"`
	ProviderName    string     `gorm:"column:provider_name;type:varchar(255);not null" json:"provider_name"`
	ModelName       string     `gorm:"column:model_name;type:varchar(255);not null" json:"model_name"`
	ModelType       string     `gorm:"column:model_type;type:varchar(40);not null" json:"model_type"`
	EncryptedConfig string     `gorm:"column:encrypted_config;type:text" json:"encrypted_config"`
	IsValid         bool       `gorm:"column:is_valid;type:tinyint(1);not null;default:0" json:"is_valid"`
	CreatedAt       *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       *time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName get sql table name.获取数据库表名
func (ProviderModel) TableName() string {
	return "provider_models"
}

// TenantDefaultModel [...]
type TenantDefaultModel struct {
	ID           string                          `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	TenantID     string                          `gorm:"column:tenant_id;type:varchar(36);not null" json:"tenant_id"`
	Tenant       *Tenant                         `json:"tenant" form:"tenant" gorm:"foreignKey:TenantID;references:ID;"`
	ProviderName string                          `gorm:"column:provider_name;type:varchar(255);not null" json:"provider_name"`
	ModelName    string                          `gorm:"column:model_name;type:varchar(255);not null" json:"model_name"`
	ModelType    modelruntimeenumtypes.ModelType `gorm:"column:model_type;type:varchar(40);not null" json:"model_type"`
	CreatedAt    *time.Time                      `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    *time.Time                      `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName get sql table name.获取数据库表名
func (TenantDefaultModel) TableName() string {
	return "tenant_default_models"
}

// TenantPreferredModelProvider [...]
type TenantPreferredModelProvider struct {
	ID                    string                         `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	TenantID              string                         `gorm:"column:tenant_id;type:varchar(36);not null" json:"tenant_id"`
	ProviderName          string                         `gorm:"column:provider_name;type:varchar(255);not null" json:"provider_name"`
	PreferredProviderType providerenumtypes.ProviderType `gorm:"column:preferred_provider_type;type:varchar(40);not null" json:"preferred_provider_type"`
	CreatedAt             *time.Time                     `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt             *time.Time                     `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName get sql table name.获取数据库表名
func (TenantPreferredModelProvider) TableName() string {
	return "tenant_preferred_model_providers"
}

// ProviderOrder [...]
type ProviderOrder struct {
	ID               string     `gorm:"column:id;type:varchar(36);not null" json:"id"`
	TenantID         string     `gorm:"column:tenant_id;type:varchar(36);not null" json:"tenant_id"`
	ProviderName     string     `gorm:"column:provider_name;type:varchar(255);not null" json:"provider_name"`
	AccountID        string     `gorm:"column:account_id;type:varchar(36);not null" json:"account_id"`
	PaymentProductID string     `gorm:"column:payment_product_id;type:varchar(191);not null" json:"payment_product_id"`
	PaymentID        string     `gorm:"column:payment_id;type:varchar(191)" json:"payment_id"`
	TransactionID    string     `gorm:"column:transaction_id;type:varchar(191)" json:"transaction_id"`
	Quantity         int        `gorm:"column:quantity;type:int;not null;default:1" json:"quantity"`
	Currency         string     `gorm:"column:currency;type:varchar(40)" json:"currency"`
	TotalAmount      int        `gorm:"column:total_amount;type:int" json:"total_amount"`
	PaymentStatus    string     `gorm:"column:payment_status;type:varchar(40);default:wait_pay" json:"payment_status"`
	PaidAt           *time.Time `gorm:"column:paid_at;type:timestamp" json:"paid_at"`
	PayFailedAt      *time.Time `gorm:"column:pay_failed_at;type:timestamp" json:"pay_failed_at"`
	RefundedAt       *time.Time `gorm:"column:refunded_at;type:timestamp" json:"refunded_at"`
	CreatedAt        *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        *time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName get sql table name.获取数据库表名
func (ProviderOrder) TableName() string {
	return "provider_orders"
}

// ProviderModelSetting [...]
type ProviderModelSetting struct {
	ID                   string     `gorm:"column:id;type:varchar(36);not null" json:"id"`
	TenantID             string     `gorm:"column:tenant_id;type:varchar(36);not null" json:"tenant_id"`
	ProviderName         string     `gorm:"column:provider_name;type:varchar(255);not null" json:"provider_name"`
	ModelName            string     `gorm:"column:model_name;type:varchar(255);not null" json:"model_name"`
	ModelType            string     `gorm:"column:model_type;type:varchar(40);not null" json:"model_type"`
	Enabled              bool       `gorm:"column:enabled;" json:"enabled"`
	LoadBalancingEnabled bool       `gorm:"column:load_balancing_enabled;" json:"load_balancing_enabled"`
	CreatedAt            *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt            *time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName get sql table name.获取数据库表名
func (ProviderModelSetting) TableName() string {
	return "provider_model_settings"
}

// LoadBalancingModelConfig [...]
type LoadBalancingModelConfig struct {
	ID              string     `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	TenantID        string     `gorm:"column:tenant_id;type:varchar(36);not null" json:"tenant_id"`
	ProviderName    string     `gorm:"column:provider_name;type:varchar(255);not null" json:"provider_name"`
	ModelName       string     `gorm:"column:model_name;type:varchar(255);not null" json:"model_name"`
	ModelType       string     `gorm:"column:model_type;type:varchar(40);not null" json:"model_type"`
	Name            string     `gorm:"column:name;type:varchar(255);not null" json:"name"`
	EncryptedConfig string     `gorm:"column:encrypted_config;type:text" json:"encrypted_config"`
	Enabled         bool       `gorm:"column:enabled;type:tinyint(1);not null;default:1" json:"enabled"`
	CreatedAt       *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       *time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName get sql table name.获取数据库表名
func (LoadBalancingModelConfig) TableName() string {
	return "load_balancing_model_configs"
}

// ProviderCredential [...]
type ProviderCredential struct {
	ID              string     `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	TenantID        string     `gorm:"column:tenant_id;type:varchar(36);not null" json:"tenant_id"`
	ProviderName    string     `gorm:"column:provider_name;type:varchar(255);not null" json:"provider_name"`
	CredentialName  string     `gorm:"column:credential_name;type:varchar(255);not null" json:"credential_name"`
	EncryptedConfig string     `gorm:"column:encrypted_config;type:text;not null" json:"encrypted_config"`
	CreatedAt       *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       *time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName get sql table name.获取数据库表名
func (ProviderCredential) TableName() string {
	return "provider_credentials"
}

// ProviderModelCredential [...]
type ProviderModelCredential struct {
	ID              string     `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	TenantID        string     `gorm:"column:tenant_id;type:varchar(36);not null" json:"tenant_id"`
	ProviderName    string     `gorm:"column:provider_name;type:varchar(255);not null" json:"provider_name"`
	ModelName       string     `gorm:"column:model_name;type:varchar(255);not null" json:"model_name"`
	ModelType       string     `gorm:"column:model_type;type:varchar(40);not null" json:"model_type"`
	CredentialName  string     `gorm:"column:credential_name;type:varchar(255);not null" json:"credential_name"`
	EncryptedConfig string     `gorm:"column:encrypted_config;type:text;not null" json:"encrypted_config"`
	CreatedAt       *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       *time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName get sql table name.获取数据库表名
func (ProviderModelCredential) TableName() string {
	return "provider_model_credentials"
}
