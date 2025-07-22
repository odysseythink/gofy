package models

import (
	"encoding/json"
	"time"

	dbengine "mlib.com/gofy/server/db_engine"
	enumtypes "mlib.com/gofy/server/enum_types"
	"mlib.com/mlog"
)

// Account [...]
type Account struct {
	ID                string                  `gorm:"column:id;type:varchar(36);not null" json:"id"`
	Name              string                  `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Email             string                  `gorm:"column:email;type:varchar(255);not null" json:"email"`
	Password          string                  `gorm:"column:password;type:varchar(255)" json:"password"`
	PasswordSalt      string                  `gorm:"column:password_salt;type:varchar(255)" json:"password_salt"`
	Avatar            string                  `gorm:"column:avatar;type:varchar(255)" json:"avatar"`
	InterfaceLanguage string                  `gorm:"column:interface_language;type:varchar(255)" json:"interface_language"`
	InterfaceTheme    string                  `gorm:"column:interface_theme;type:varchar(255)" json:"interface_theme"`
	Timezone          string                  `gorm:"column:timezone;type:varchar(255)" json:"timezone"`
	LastLoginAt       *time.Time              `gorm:"column:last_login_at;type:timestamp" json:"last_login_at"`
	LastLoginIP       string                  `gorm:"column:last_login_ip;type:varchar(255)" json:"last_login_ip"`
	Status            enumtypes.AccountStatus `gorm:"column:status;type:varchar(16);default:active" json:"status"`
	InitializedAt     *time.Time              `gorm:"column:initialized_at;type:timestamp" json:"initialized_at"`
	CreatedAt         *time.Time              `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt         *time.Time              `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	LastActiveAt      *time.Time              `gorm:"column:last_active_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"last_active_at"`
	RememberMe        bool                    `json:"remember_me" form:"-" gorm:"-:all"`
	InviteToken       string                  `json:"invite_token" form:"-" gorm:"-:all"`
	Language          string                  `json:"language" form:"-" gorm:"-:all"`
	_current_tenant   *Tenant                 `json:"-" form:"-" gorm:"-:all"`
	_role             string                  `json:"-" form:"-" gorm:"-:all"`
}

// TableName get sql table name.获取数据库表名
func (Account) TableName() string {
	return "accounts"
}

func (acc *Account) IsPasswordSet() bool {
	return acc.Password != ""
}

func (acc *Account) CurrentTenant() *Tenant {
	// FIXME: fix the type error later, because the type is important maybe cause some bugs
	return acc._current_tenant // type: ignore
}
func (acc *Account) SetCurrentTenant(value *Tenant) {
	tenant := value
	ta := new(TenantAccountJoin)
	err := dbengine.Instance().DB.Model(&TenantAccountJoin{}).Where("tenant_id=? and account_id = ?", tenant.ID, acc.ID).First(ta).Error
	if err != nil {
		mlog.Errorf("get TenantAccountJoin from mysql failed:%v", err)
		ta = nil
	}
	if ta != nil {
		tenant.CurrentRole = ta.Role
	} else {
		// FIXME: fix the type error later, because the type is important maybe cause some bugs
		tenant = nil // type: ignore
	}
	acc._current_tenant = tenant
}

func (acc *Account) CurrentTenantID() string {
	if acc._current_tenant != nil {
		return acc._current_tenant.ID
	} else {
		return ""
	}
}
func (acc *Account) SetCurrentTenantID(value string) {
	var tenant_account_join struct {
		*Tenant
		Role string `gorm:"column:role" json:"role"`
	}
	var tenant *Tenant
	err := dbengine.Instance().DB.Raw("SELECT t.*, ta.role FROM tenants t join tenant_account_joins ta on t.id = ta.tenant_id WHERE ta.account_id = ? AND t.id = ?;", acc.ID, value).First(&tenant_account_join).Error
	if err != nil {
		mlog.Errorf("get TenantAccountJoin from mysql failed:%v", err)
	} else {
		tenant = tenant_account_join.Tenant
		tenant.CurrentRole = tenant_account_join.Role
	}

	acc._current_tenant = tenant
}
func (acc *Account) SetRole(v string) {
	acc._role = v
}
func (acc *Account) Role() string {
	return acc._role
}
func (acc *Account) CurrentRole() string {
	return acc._current_tenant.CurrentRole
}
func (acc *Account) IsEditor() bool {
	if acc._current_tenant == nil {
		return false
	}
	return enumtypes.TenantAccountRole(acc._current_tenant.CurrentRole).IsEditingRole()
}

func (acc *Account) IsAdminOrOwner() bool {
	if acc._current_tenant == nil {
		return false
	}
	return enumtypes.TenantAccountRole(acc._current_tenant.CurrentRole).IsPrivilegedRole()
}
func (acc *Account) IsAdmin() bool {
	if acc._current_tenant == nil {
		return false
	}
	return enumtypes.TenantAccountRole(acc._current_tenant.CurrentRole).IsAdminRole()
}

func (acc *Account) IsDatasetEditor() bool {
	if acc._current_tenant == nil {
		return false
	}
	return enumtypes.TenantAccountRole(acc._current_tenant.CurrentRole).IsDatasetEditRole()
}
func (acc *Account) IsDatasetOperator() bool {
	if acc._current_tenant == nil {
		return false
	}
	return enumtypes.TenantAccountRole(acc._current_tenant.CurrentRole) == enumtypes.TenantAccountRole_DATASET_OPERATOR
}

// Tenant [...]
type Tenant struct {
	ID               string                 `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	Name             string                 `gorm:"column:name;type:varchar(255);not null" json:"name"`
	EncryptPublicKey string                 `gorm:"column:encrypt_public_key;type:text" json:"encrypt_public_key"`
	Plan             string                 `gorm:"column:plan;type:varchar(255);default:basic" json:"plan"`
	Status           enumtypes.TenantStatus `gorm:"column:status;type:varchar(255);default:normal" json:"status"`
	CreatedAt        *time.Time             `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        *time.Time             `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	CustomConfig     string                 `gorm:"column:custom_config;type:text" json:"custom_config"`
	CurrentRole      string                 `gorm:"-:all" json:"role"`
	Current          bool                   `json:"current" form:"-" gorm:"-:all"`
}

// TableName get sql table name.获取数据库表名
func (Tenant) TableName() string {
	return "tenants"
}

// func (t *Tenant) IsEditingRole() bool {
// 	return (t.Role == enumtypes.TenantAccountRole_OWNER || t.Role == enumtypes.TenantAccountRole_ADMIN || t.Role == enumtypes.TenantAccountRole_EDITOR)
// }

func (t *Tenant) GetAccounts() []*Account {
	var tmp_list []*Account
	err := dbengine.Instance().DB.Raw("SELECT acc.* FROM accounts acc join tenant_account_joins ta on acc.id = ta.account_id WHERE ta.tenant_id = ?;", t.ID).Find(&tmp_list).Error
	if err != nil {
		mlog.Errorf("get Account failed:%v", err)
	}
	return tmp_list
}
func (t *Tenant) CustomConfigDict() map[string]any {
	if t.CustomConfig != "" {
		var tmp map[string]any
		err := json.Unmarshal([]byte(t.CustomConfig), &tmp)
		if err != nil {
			mlog.Errorf("json unmarshal(%s) to dict failed:%v", t.CustomConfig, err)
			return map[string]any{}
		}
		return tmp
	}
	return map[string]any{}
}
func (t *Tenant) SetCustomConfigDict(value map[string]any) {
	bindata, _ := json.Marshal(value)
	t.CustomConfig = string(bindata)
}

// TenantAccountJoin [...]
type TenantAccountJoin struct {
	ID       string `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	TenantID string `gorm:"column:tenant_id;type:varchar(36);not null" json:"tenant_id"`
	// JoinTenant *Tenant `json:"tenant" form:"tenant" gorm:"foreignKey:TenantID;references:ID;"`
	AccountID string `gorm:"column:account_id;type:varchar(36);not null" json:"account_id"`
	// JoinAccount *Account   `json:"account" form:"account" gorm:"foreignKey:AccountID;references:ID;"`
	Role      string     `gorm:"column:role;type:varchar(16);default:normal" json:"role"`
	InvitedBy string     `gorm:"column:invited_by;type:varchar(36)" json:"invited_by"`
	CreatedAt *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	Current   bool       `gorm:"column:current;type:tinyint(1);not null;default:0" json:"current"`
}

// TableName get sql table name.获取数据库表名
func (TenantAccountJoin) TableName() string {
	return "tenant_account_joins"
}

// AccountIntegrate [...]
type AccountIntegrate struct {
	ID             string     `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	AccountID      string     `gorm:"column:account_id;type:varchar(36);not null" json:"account_id"`
	Provider       string     `gorm:"column:provider;type:varchar(16);not null" json:"provider"`
	OpenID         string     `gorm:"column:open_id;type:varchar(255);not null" json:"open_id"`
	EncryptedToken string     `gorm:"column:encrypted_token;type:varchar(255);not null" json:"encrypted_token"`
	CreatedAt      *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      *time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName get sql table name.获取数据库表名
func (AccountIntegrate) TableName() string {
	return "account_integrates"
}

// InvitationCode [...]
type InvitationCode struct {
	ID              int        `gorm:"primaryKey;column:id;type:int;not null" json:"id"`
	Batch           string     `gorm:"column:batch;type:varchar(255);not null" json:"batch"`
	Code            string     `gorm:"column:code;type:varchar(32);not null" json:"code"`
	Status          string     `gorm:"column:status;type:varchar(16);default:unused" json:"status"`
	UsedAt          *time.Time `gorm:"column:used_at;type:timestamp" json:"used_at"`
	UsedByTenantID  string     `gorm:"column:used_by_tenant_id;type:varchar(36)" json:"used_by_tenant_id"`
	UsedByAccountID string     `gorm:"column:used_by_account_id;type:varchar(36)" json:"used_by_account_id"`
	DeprecatedAt    *time.Time `gorm:"column:deprecated_at;type:timestamp" json:"deprecated_at"`
	CreatedAt       *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
}

// TableName get sql table name.获取数据库表名
func (InvitationCode) TableName() string {
	return "invitation_codes"
}
