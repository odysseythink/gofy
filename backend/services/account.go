package services

import (
	"context"
	crand "crypto/rand"
	"encoding/base64"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/odysseythink/confy"
	"github.com/odysseythink/mlog"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"mlib.com/gofy/server/cache"
	"mlib.com/gofy/server/core/exceptions"
	httpexceptions "mlib.com/gofy/server/core/exceptions/http"
	tokenmanager "mlib.com/gofy/server/core/manageres/token_manager"
	dbengine "mlib.com/gofy/server/db_engine"
	enumtypes "mlib.com/gofy/server/enum_types"
	"mlib.com/gofy/server/global"
	"mlib.com/gofy/server/libs/helper"
	"mlib.com/gofy/server/libs/password"
	"mlib.com/gofy/server/models"
	"mlib.com/gofy/server/models/request"
	jwtutils "mlib.com/gofy/server/utils/jwt"
)

const (
	LOGIN_MAX_ERROR_LIMITS       = 5
	REFRESH_TOKEN_PREFIX         = "refresh_token:"
	ACCOUNT_REFRESH_TOKEN_PREFIX = "account_refresh_token:"
)

type AccountService struct {
}

var (
	reset_password_rate_limiter              = helper.NewRateLimiter("reset_password_rate_limit", 1, 60*1)
	email_code_login_rate_limiter            = helper.NewRateLimiter("email_code_login_rate_limit", 1, 60*1)
	email_code_account_deletion_rate_limiter = helper.NewRateLimiter("email_code_account_deletion_rate_limit", 1, 60*1)
)

func (s *AccountService) CreateAccount(
	email string,
	name string,
	interface_language string,
	pwd string,
	interface_theme string,
	is_setup bool,
) *models.Account {
	if interface_theme == "" {
		interface_theme = "light"
	}
	if !ServiceGroupApp.Feature.GetSystemFeatures().IsAllowRegister && !is_setup {
		panic(httpexceptions.NewAccountNotFound())
	}

	if confy.Get[bool]("billing_enabled") && ServiceGroupApp.Billing.IsEmailInFreeze(email) {
		panic(exceptions.NewAccountRegisterError("This email account has been deleted within the past 30 days and is temporarily unavailable for new account registration"))
	}
	now := time.Now()
	timezone := global.LANGUAGE_TIMEZONE_MAPPING[interface_language]
	if timezone == "" {
		timezone = "UTC"
	}
	account := &models.Account{
		Model: models.Model{
			ID:        uuid.NewV4().String(),
			CreatedAt: &now,
			UpdatedAt: &now,
		},
		Name:              name,
		Email:             email,
		InterfaceLanguage: interface_language,
		InterfaceTheme:    interface_theme,
		Timezone:          timezone,
		LastLoginAt:       &now,
		InitializedAt:     &now,
		LastActiveAt:      &now,
	}

	if pwd != "" {
		// generate password salt
		salt := make([]byte, 16)
		_, err := crand.Read(salt)
		if err != nil {
			mlog.Errorf("create salt failed:%v", err)
			panic(exceptions.NewValueError(fmt.Sprintf("create salt failed:%v", err)))
		}

		base64_salt := base64.StdEncoding.EncodeToString(salt)
		// encrypt password with salt
		password_hashed := password.HashPassword(pwd, salt)
		base64_password_hashed := base64.StdEncoding.EncodeToString([]byte(password_hashed))

		account.Password = base64_password_hashed
		account.PasswordSalt = base64_salt
	}
	dbengine.Instance().DB.Create(account)
	return account
}

// Create 创建Account记录
func (s *AccountService) Create(acc models.Account) (err error) {
	err = dbengine.Instance().DB.Create(&acc).Error
	return err
}

// Delete 删除Account记录
func (s *AccountService) Delete(acc models.Account) (err error) {
	err = dbengine.Instance().DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Account{}).Where("id = ?", acc.ID).Error; err != nil {
			return err
		}
		if err = tx.Delete(&acc).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// DeleteByIds 批量删除Account记录
func (s *AccountService) DeleteByIds(ids []string) (err error) {
	if err := dbengine.Instance().DB.Model(&models.Account{}).Delete("id in ?", ids).Error; err != nil {
		return err
	}

	return nil
}

// Update 更新Account记录
func (s *AccountService) Update(acc *models.Account) (err error) {
	err = dbengine.Instance().DB.Save(acc).Error
	return err
}

// Get 根据id获取Account记录
func (s *AccountService) Get(id string) (acc *models.Account, err error) {
	acc = &models.Account{}
	err = dbengine.Instance().DB.Where("id = ?", id).First(acc).Error
	if err != nil {
		acc = nil
	}
	return
}

// GetByEmail 根据email获取Account记录
func (s *AccountService) GetByEmail(email string) (acc *models.Account, err error) {
	acc = &models.Account{}
	err = dbengine.Instance().DB.Debug().Where("email = ?", email).First(&acc).Error
	if err != nil {
		acc = nil
	}
	return
}

// GetList 分页获取Account记录
func (s *AccountService) GetList(info request.PageInfo) (list []models.Account, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := dbengine.Instance().DB.Model(&models.Account{})
	var endpoints []models.Account

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	err = db.Limit(limit).Offset(offset).Find(&endpoints).Error
	return endpoints, total, err
}

func (s *AccountService) Authenticate(email, pwd, invite_token string) *models.Account {
	// """authenticate account with email and password"""
	now := time.Now()

	account, err := s.GetByEmail(email)
	if err != nil {
		mlog.Errorf("get account by email(%s) failed:%v", email, err)
		panic(exceptions.NewAccountNotFoundError("account not found"))
	}
	if account.Status == enumtypes.AccountStatusBANNED {
		panic(exceptions.NewAccountLoginError("Account is banned."))
	}
	if pwd != "" && /*invite_token &&*/ account.Password == "" {
		// # if invite_token is valid, set password and password_salt
		// salt = secrets.token_bytes(16)
		// base64_salt = base64.b64encode(salt).decode()
		// password_hashed = hash_password(password, salt)
		// base64_password_hashed = base64.b64encode(password_hashed).decode()
		// account.password = base64_password_hashed
		// account.password_salt = base64_salt
		salt := make([]byte, 16)
		crand.Read(salt)

		base64_salt := base64.StdEncoding.EncodeToString(salt)
		password_hashed := password.HashPassword(pwd, salt)
		base64_password_hashed := base64.StdEncoding.EncodeToString([]byte(password_hashed))
		account.Password = base64_password_hashed
		account.PasswordSalt = base64_salt
	}
	if account.Password == "" || !password.ComparePassword(pwd, account.Password, account.PasswordSalt) {
		panic(exceptions.NewAccountPasswordError("Invalid email or password."))
	}

	if account.Status == enumtypes.AccountStatusPENDING {
		account.Status = enumtypes.AccountStatusACTIVE
		account.InitializedAt = &now
	}
	s.Update(account)
	mlog.Debugf("----account=%#v", account)
	return account
}

func (s *AccountService) UpdateLastLogin(account *models.Account, ip_address string) {
	// """Update last login time and ip"""
	now := time.Now()
	account.LastLoginAt = &now
	account.LastLoginIP = ip_address
	s.Update(account)
}

func (s *AccountService) GetAccountJWTToken(account *models.Account, exp time.Duration) string {
	jtokenstr, _, _ := jwtutils.GenToken(account.ID, confy.Get[string]("SECRET_KEY"), confy.Get[string]("EDITION"), 12*3600*time.Second)
	return jtokenstr
}

func (s *AccountService) _get_login_cache_key(account_id, token string) string {
	return fmt.Sprintf("account_login:%s:%s", account_id, token)
}
func (s *AccountService) _get_refresh_token_key(refresh_token string) string {
	return REFRESH_TOKEN_PREFIX + refresh_token
}

func (s *AccountService) _get_account_refresh_token_key(account_id string) string {
	return ACCOUNT_REFRESH_TOKEN_PREFIX + account_id
}
func (s *AccountService) _store_refresh_token(refresh_token, account_id string) {
	cache.Instance().SetEx(s._get_refresh_token_key(refresh_token), account_id, time.Duration(confy.GetWithDefault[int]("refresh_token_expire_days", 30))*time.Second*3600*24)
	cache.Instance().SetEx(s._get_account_refresh_token_key(account_id), refresh_token, time.Duration(confy.GetWithDefault[int]("refresh_token_expire_days", 30))*time.Second*3600*24)
}
func (s *AccountService) _delete_refresh_token(refresh_token string, account_id string) {
	cache.Instance().DelKey(s._get_refresh_token_key(refresh_token))
	cache.Instance().DelKey(s._get_account_refresh_token_key(account_id))
}

func (s *AccountService) Login(account *models.Account, ip_address string) *TokenPair {
	if ip_address != "" {
		s.UpdateLastLogin(account, ip_address)
	}
	if account.Status == enumtypes.AccountStatusPENDING {
		account.Status = enumtypes.AccountStatusACTIVE
		s.Update(account)
	}
	exp := 30 * time.Second * 3600 * 24
	access_token := s.GetAccountJWTToken(account, exp)
	refresh_token := password.GenerateRefreshToken(0)
	mlog.Debugf("----access_token=%s,refresh_token=%s ", access_token, refresh_token)
	s._store_refresh_token(refresh_token, account.ID)

	return &TokenPair{
		AccessToken:  access_token,
		RefreshToken: refresh_token,
	}
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (s *AccountService) Logout(account *models.Account) {
	refresh_token := cache.Instance().GetString(s._get_account_refresh_token_key(account.ID))
	if refresh_token != "" {
		s._delete_refresh_token(refresh_token, account.ID)
	}
}

func (s *AccountService) SetCurrentTenant(acc *models.Account, value *models.Tenant) {
	if acc == nil || value == nil {
		mlog.Errorf("invalid args(acc=%#v, tenant=%#v)", acc, value)
		return
	}
	tenant := value
	ta := &models.TenantAccountJoin{}
	err := dbengine.Instance().DB.Model(&models.TenantAccountJoin{}).Where("tenant_id = ? and account_id = ?", tenant.ID, acc.ID).First(ta).Error
	if err != nil {
		mlog.Errorf("get TenantAccountJoin by(tenant_id = %s and account_id = %s) failed:%v", tenant.ID, acc.ID, err)
		tenant = nil
	} else {
		tenant.CurrentRole = ta.Role
	}
	acc.SetCurrentTenant(tenant)
}

func (s *AccountService) SetCurrentTenantID(acc *models.Account, value string) {
	if acc == nil || value == "" {
		mlog.Errorf("invalid args(acc=%#v, tenant_id=%#v)", acc, value)
		return
	}
	result := struct {
		models.Tenant
		TaRole string `gorm:"column:ta_role;" json:"ta_role"`
	}{
		Tenant: models.Tenant{},
	}
	err := dbengine.Instance().DB.Raw("SELECT t.*, ta.role as ta_role FROM tenants t join tenant_account_joins ta on t.id = ta.tenant_id WHERE ta.account_id = ? AND t.id = ?;", acc.ID, value).First(&result).Error
	if err != nil {
		acc.SetCurrentTenant(nil)
		mlog.Errorf("exec(SELECT t.*, ta.role as ta_role FROM tenants t join tenant_account_joins ta on t.id = ta.tenant_id WHERE ta.account_id = '%s' AND t.id = '%s';) failed:%v", acc.ID, value, err)
		return
	} else {
		result.Tenant.CurrentRole = result.TaRole
		acc.SetCurrentTenant(&result.Tenant)
	}
}

func (s *AccountService) LoadUser(account_id string) (*models.Account, error) {
	account, err := s.Get(account_id)
	if err != nil || account == nil {
		mlog.Errorf("get account by user_id=%s failed:%v", account_id, err)
		return nil, exceptions.NewAccountNotFoundError("")
	}
	if account.Status == enumtypes.AccountStatusBANNED {
		mlog.Error("account is banned or closed")
		return nil, httpexceptions.NewUnauthorizedAndForceLogout("Account is banned.")
	}
	result := struct {
		models.Tenant
		TaID string `gorm:"column:ta_id;" json:"ta_id"`
		Role string `gorm:"column:role;" json:"ta_role"`
	}{
		Tenant: models.Tenant{},
	}
	err = dbengine.Instance().DB.Debug().Raw("SELECT t.*, ta.role, ta.id as ta_id FROM tenants t join tenant_account_joins ta on t.id = ta.tenant_id WHERE ta.account_id = ? AND ta.current = ?;", account.ID, true).Scan(&result).Error
	if err != nil {
		mlog.Warningf("get current TenantAccountJoin by account_id=%s failed:%v", account.ID, err)
		err = dbengine.Instance().DB.Raw("SELECT t.*, ta.role, ta.id as ta_id FROM tenants t join tenant_account_joins ta on t.id = ta.tenant_id WHERE ta.account_id = ?;", account.ID).First(&result).Error
		if err != nil {
			mlog.Errorf("get TenantAccountJoin by account_id=%s failed:%v", account.ID, err)
			return nil, exceptions.NewAccountNotFoundError("")
		} else {
			result.Tenant.CurrentRole = result.Role
			account.SetCurrentTenant(&result.Tenant)
			dbengine.Instance().DB.Updates(&models.TenantAccountJoin{Model: models.Model{ID: result.TaID}, Current: true})
		}
	} else {
		result.Tenant.CurrentRole = result.Role
		account.SetCurrentTenant(&result.Tenant)
	}

	now := time.Now()

	if now.Sub(*account.LastActiveAt) > 10*time.Minute {
		account.LastActiveAt = &now
		s.Update(account)
	}
	return account, nil
}

func (s *AccountService) IsLoginErrorRateLimit(email string) bool {
	key := fmt.Sprintf("login_error_rate_limit:%s", email)
	var countstr string
	if err := cache.Instance().Get(key, &countstr); err != nil {
		mlog.Errorf("get redis(%s) failed:%v", key, err)
		return false
	}
	if count, err := strconv.Atoi(countstr); err != nil {
		mlog.Errorf("get redis(%s)=%s convert to int failed:%v", key, countstr, err)
		return false
	} else {
		if count > LOGIN_MAX_ERROR_LIMITS {
			return true
		}
	}
	return false

}

func (s *AccountService) AddLoginErrorRateLimit(email string) {
	key := fmt.Sprintf("login_error_rate_limit:%s", email)
	var countstr string
	err := cache.Instance().Get(key, &countstr)
	if err != nil {
		mlog.Warningf("get redis(%s) failed:%v", key, err)
		countstr = "0"
	}
	count := 0
	count, err = strconv.Atoi(countstr)
	if err != nil {
		mlog.Errorf("get redis(%s)=%s convert to int failed:%v", key, countstr, err)
		count = 0
	} else {
		count += 1
	}
	cache.Instance().SetEx(key, count, time.Duration(confy.Get[int]("login_lockout_duration"))*time.Second)
}
func (cls *AccountService) SendResetPasswordEmail(account *models.Account, email string, language string /*"en-US"*/) string {
	account_email := email
	if account != nil {
		account_email = account.Email
	}

	if account_email == "" {
		panic(exceptions.NewValueError("Email must be provided."))
	}
	if reset_password_rate_limiter.IsRateLimited(account_email) {
		panic(httpexceptions.NewPasswordResetRateLimitExceededError())
	}
	code := ""
	{
		for i := 0; i < 6; i++ {
			code += strconv.Itoa(rand.Intn(9))
		}
	}
	token, _ := (&tokenmanager.TokenManager{}).GenerateToken("reset_password", account, email, map[string]any{"code": code})
	// send_reset_password_mail_task.delay(
	// 	language=language,
	// 	to=account_email,
	// 	code=code,
	// )
	// cls.reset_password_rate_limiter.increment_rate_limit(account_email)
	// return token
	return token
}

func (cls *AccountService) ResetLoginErrorRateLimit(email string) {
	key := fmt.Sprintf("login_error_rate_limit:%s", email)
	cache.Instance().DelKey(key)
}
func (cls *AccountService) LoadLoggedInAccount(account_id string) (*models.Account, error) {
	return cls.LoadUser(account_id)
}

// === Registration & Invitation ===

// CreateAccountWithEmail creates a new account with email/password after checking uniqueness.
func (s *AccountService) CreateAccountWithEmail(email, name, pwd, language string) (*models.Account, error) {
	var count int64
	dbengine.Instance().DB.Model(&models.Account{}).Where("email = ?", email).Count(&count)
	if count > 0 {
		return nil, fmt.Errorf("email already registered")
	}

	now := time.Now()
	timezone := global.LANGUAGE_TIMEZONE_MAPPING[language]
	if timezone == "" {
		timezone = "UTC"
	}
	account := &models.Account{
		Model: models.Model{
			ID:        uuid.NewV4().String(),
			CreatedAt: &now,
			UpdatedAt: &now,
		},
		Email:             email,
		Name:              name,
		InterfaceLanguage: language,
		Status:            enumtypes.AccountStatusACTIVE,
		Timezone:          timezone,
		LastLoginAt:       &now,
		InitializedAt:     &now,
		LastActiveAt:      &now,
	}
	if pwd != "" {
		salt := make([]byte, 16)
		crand.Read(salt)
		base64_salt := base64.StdEncoding.EncodeToString(salt)
		password_hashed := password.HashPassword(pwd, salt)
		base64_password_hashed := base64.StdEncoding.EncodeToString([]byte(password_hashed))
		account.Password = base64_password_hashed
		account.PasswordSalt = base64_salt
	}
	if err := dbengine.Instance().DB.Create(account).Error; err != nil {
		return nil, err
	}
	return account, nil
}

// CheckEmailUnique returns true if the email is not already registered.
func (s *AccountService) CheckEmailUnique(email string) bool {
	var count int64
	dbengine.Instance().DB.Model(&models.Account{}).Where("email = ?", email).Count(&count)
	return count == 0
}

// === Tenant/Member Management ===

// CreateTenant creates a new tenant and adds the creator as owner.
func (s *AccountService) CreateTenant(name string, creatorID string) (*models.Tenant, error) {
	now := time.Now()
	tenant := &models.Tenant{
		Model: models.Model{
			ID:        uuid.NewV4().String(),
			CreatedAt: &now,
			UpdatedAt: &now,
		},
		Name: name,
	}
	if err := dbengine.Instance().DB.Create(tenant).Error; err != nil {
		return nil, err
	}
	join := &models.TenantAccountJoin{
		Model: models.Model{
			ID:        uuid.NewV4().String(),
			CreatedAt: &now,
			UpdatedAt: &now,
		},
		TenantID:  tenant.ID,
		AccountID: creatorID,
		Role:      string(enumtypes.TenantAccountRole_OWNER),
		Current:   true,
	}
	if err := dbengine.Instance().DB.Create(join).Error; err != nil {
		return nil, err
	}
	return tenant, nil
}

// InviteMember invites an existing account to a tenant by email.
func (s *AccountService) InviteMember(tenantID, email, role, inviterID string) error {
	account, err := s.GetByEmail(email)
	if err != nil || account == nil {
		return fmt.Errorf("account not found for email: %s", email)
	}
	var count int64
	dbengine.Instance().DB.Model(&models.TenantAccountJoin{}).Where("tenant_id = ? AND account_id = ?", tenantID, account.ID).Count(&count)
	if count > 0 {
		return fmt.Errorf("user is already a member")
	}
	now := time.Now()
	join := &models.TenantAccountJoin{
		Model: models.Model{
			ID:        uuid.NewV4().String(),
			CreatedAt: &now,
			UpdatedAt: &now,
		},
		TenantID:  tenantID,
		AccountID: account.ID,
		Role:      role,
		InvitedBy: inviterID,
	}
	return dbengine.Instance().DB.Create(join).Error
}

// RemoveMember removes a non-owner member from a tenant.
func (s *AccountService) RemoveMember(tenantID, accountID string) error {
	var join models.TenantAccountJoin
	if err := dbengine.Instance().DB.Where("tenant_id = ? AND account_id = ?", tenantID, accountID).First(&join).Error; err != nil {
		return fmt.Errorf("member not found")
	}
	if join.Role == string(enumtypes.TenantAccountRole_OWNER) {
		return fmt.Errorf("cannot remove the owner")
	}
	return dbengine.Instance().DB.Delete(&join).Error
}

// UpdateMemberRole updates the role of a tenant member.
func (s *AccountService) UpdateMemberRole(tenantID, accountID, newRole string) error {
	return dbengine.Instance().DB.Model(&models.TenantAccountJoin{}).
		Where("tenant_id = ? AND account_id = ?", tenantID, accountID).
		Update("role", newRole).Error
}

// GetUserRole returns the role of a user in a tenant, or empty string if not a member.
func (s *AccountService) GetUserRole(tenantID, accountID string) string {
	var join models.TenantAccountJoin
	if err := dbengine.Instance().DB.Where("tenant_id = ? AND account_id = ?", tenantID, accountID).First(&join).Error; err != nil {
		return ""
	}
	return join.Role
}

// IsMember checks whether an account is a member of a tenant.
func (s *AccountService) IsMember(tenantID, accountID string) bool {
	var count int64
	dbengine.Instance().DB.Model(&models.TenantAccountJoin{}).Where("tenant_id = ? AND account_id = ?", tenantID, accountID).Count(&count)
	return count > 0
}

// GetTenantMembers returns all TenantAccountJoin records for a tenant.
func (s *AccountService) GetTenantMembers(tenantID string) []*models.TenantAccountJoin {
	var joins []*models.TenantAccountJoin
	dbengine.Instance().DB.Where("tenant_id = ?", tenantID).Find(&joins)
	return joins
}

// GetTenantMemberAccounts returns all Account objects for members of a tenant.
func (s *AccountService) GetTenantMemberAccounts(tenantID string) []*models.Account {
	var accounts []*models.Account
	dbengine.Instance().DB.Raw("SELECT acc.* FROM accounts acc JOIN tenant_account_joins ta ON acc.id = ta.account_id WHERE ta.tenant_id = ?", tenantID).Find(&accounts)
	return accounts
}

// GetAccountTenants returns all tenants an account belongs to.
func (s *AccountService) GetAccountTenants(accountID string) []*models.Tenant {
	var tenants []*models.Tenant
	dbengine.Instance().DB.Raw("SELECT t.* FROM tenants t JOIN tenant_account_joins ta ON t.id = ta.tenant_id WHERE ta.account_id = ?", accountID).Find(&tenants)
	return tenants
}

// GetTenantByID retrieves a tenant by ID.
func (s *AccountService) GetTenantByID(tenantID string) (*models.Tenant, error) {
	var tenant models.Tenant
	if err := dbengine.Instance().DB.Where("id = ?", tenantID).First(&tenant).Error; err != nil {
		return nil, err
	}
	return &tenant, nil
}

// UpdateTenant updates a tenant's name.
func (s *AccountService) UpdateTenant(tenantID, name string) error {
	return dbengine.Instance().DB.Model(&models.Tenant{}).Where("id = ?", tenantID).Update("name", name).Error
}

// === Password Management ===

// ChangePassword verifies the old password and sets a new one.
func (s *AccountService) ChangePassword(accountID, oldPassword, newPassword string) error {
	account, err := s.Get(accountID)
	if err != nil || account == nil {
		return fmt.Errorf("account not found")
	}
	if !password.ComparePassword(oldPassword, account.Password, account.PasswordSalt) {
		return fmt.Errorf("incorrect password")
	}
	salt := make([]byte, 16)
	crand.Read(salt)
	base64_salt := base64.StdEncoding.EncodeToString(salt)
	password_hashed := password.HashPassword(newPassword, salt)
	base64_password_hashed := base64.StdEncoding.EncodeToString([]byte(password_hashed))
	account.Password = base64_password_hashed
	account.PasswordSalt = base64_salt
	return s.Update(account)
}

// ResetPassword sets a new password for an account (no old password check).
func (s *AccountService) ResetPassword(accountID, newPassword string) error {
	account, err := s.Get(accountID)
	if err != nil || account == nil {
		return fmt.Errorf("account not found")
	}
	salt := make([]byte, 16)
	crand.Read(salt)
	base64_salt := base64.StdEncoding.EncodeToString(salt)
	password_hashed := password.HashPassword(newPassword, salt)
	base64_password_hashed := base64.StdEncoding.EncodeToString([]byte(password_hashed))
	account.Password = base64_password_hashed
	account.PasswordSalt = base64_salt
	return s.Update(account)
}

// === Account Integration ===

// LinkAccountIntegrate links a third-party provider account to an account.
func (s *AccountService) LinkAccountIntegrate(accountID, provider, openID string) error {
	now := time.Now()
	integrate := &models.AccountIntegrate{
		Model: models.Model{
			ID:        uuid.NewV4().String(),
			CreatedAt: &now,
			UpdatedAt: &now,
		},
		AccountID: accountID,
		Provider:  provider,
		OpenID:    openID,
	}
	return dbengine.Instance().DB.Create(integrate).Error
}

// GetAccountIntegrates returns all integrations for an account.
func (s *AccountService) GetAccountIntegrates(accountID string) []*models.AccountIntegrate {
	var integrates []*models.AccountIntegrate
	dbengine.Instance().DB.Where("account_id = ?", accountID).Find(&integrates)
	return integrates
}

// GetAccountByIntegrate finds an account by third-party provider and open ID.
func (s *AccountService) GetAccountByIntegrate(provider, openID string) (*models.Account, error) {
	var integrate models.AccountIntegrate
	if err := dbengine.Instance().DB.Where("provider = ? AND open_id = ?", provider, openID).First(&integrate).Error; err != nil {
		return nil, fmt.Errorf("integration not found")
	}
	return s.Get(integrate.AccountID)
}

// UnlinkAccountIntegrate removes a third-party integration from an account.
func (s *AccountService) UnlinkAccountIntegrate(accountID, provider string) error {
	return dbengine.Instance().DB.Where("account_id = ? AND provider = ?", accountID, provider).Delete(&models.AccountIntegrate{}).Error
}

// === Account Status ===

// CloseAccount sets an account's status to banned/closed.
func (s *AccountService) CloseAccount(accountID string) error {
	return dbengine.Instance().DB.Model(&models.Account{}).Where("id = ?", accountID).Update("status", enumtypes.AccountStatusBANNED).Error
}

// IsAccountActive checks if an account exists and is active.
func (s *AccountService) IsAccountActive(accountID string) bool {
	var count int64
	dbengine.Instance().DB.Model(&models.Account{}).Where("id = ? AND status = ?", accountID, enumtypes.AccountStatusACTIVE).Count(&count)
	return count > 0
}

// GetAccountByEmail returns an account by email, or nil if not found.
func (s *AccountService) GetAccountByEmail(email string) *models.Account {
	acc, err := s.GetByEmail(email)
	if err != nil {
		return nil
	}
	return acc
}

// GetAccountByID returns an account by ID, or nil if not found.
func (s *AccountService) GetAccountByID(accountID string) *models.Account {
	acc, err := s.Get(accountID)
	if err != nil {
		return nil
	}
	return acc
}

// === Profile Updates ===

// UpdateAccountName updates only the name field.
func (s *AccountService) UpdateAccountName(accountID, name string) error {
	return dbengine.Instance().DB.Model(&models.Account{}).Where("id = ?", accountID).Update("name", name).Error
}

// UpdateAccountAvatar updates only the avatar field.
func (s *AccountService) UpdateAccountAvatar(accountID, avatar string) error {
	return dbengine.Instance().DB.Model(&models.Account{}).Where("id = ?", accountID).Update("avatar", avatar).Error
}

// UpdateAccountLanguage updates the interface language.
func (s *AccountService) UpdateAccountLanguage(accountID, language string) error {
	return dbengine.Instance().DB.Model(&models.Account{}).Where("id = ?", accountID).Update("interface_language", language).Error
}

// UpdateAccountTimezone updates the timezone.
func (s *AccountService) UpdateAccountTimezone(accountID, timezone string) error {
	return dbengine.Instance().DB.Model(&models.Account{}).Where("id = ?", accountID).Update("timezone", timezone).Error
}

// UpdateAccountTheme updates the interface theme.
func (s *AccountService) UpdateAccountTheme(accountID, theme string) error {
	return dbengine.Instance().DB.Model(&models.Account{}).Where("id = ?", accountID).Update("interface_theme", theme).Error
}

// === Email Sending (stubs — need email service integration) ===

func (s *AccountService) SendRegistrationEmail(email, verificationCode string) error {
	// TODO: Integrate with email service (config/email.go)
	mlog.Infof("sending registration email to %s with code %s", email, verificationCode)
	return nil
}

func (s *AccountService) SendResetPasswordNotification(email, resetToken string) error {
	mlog.Infof("sending reset password email to %s", email)
	return nil
}

func (s *AccountService) SendInvitationEmail(inviterName, email, tenantName, inviteCode string) error {
	mlog.Infof("sending invitation email to %s for tenant %s", email, tenantName)
	return nil
}

func (s *AccountService) SendChangeEmailCode(accountID, newEmail, code string) error {
	mlog.Infof("sending email change code to %s", newEmail)
	return nil
}

func (s *AccountService) SendAccountDeletionCode(email, code string) error {
	mlog.Infof("sending account deletion code to %s", email)
	return nil
}

// === Rate Limiting ===

func (s *AccountService) CheckRateLimit(key string, maxAttempts int, windowSeconds int) error {
	cacheKey := fmt.Sprintf("rate_limit:%s", key)
	// Use Redis INCR + EXPIRE for rate limiting
	client := cache.Instance().Client()
	ctx := context.Background()

	count, err := client.Incr(ctx, cacheKey).Result()
	if err != nil {
		return nil // fail open
	}
	if count == 1 {
		client.Expire(ctx, cacheKey, time.Duration(windowSeconds)*time.Second)
	}
	if count > int64(maxAttempts) {
		return fmt.Errorf("rate limit exceeded, please try again later")
	}
	return nil
}

func (s *AccountService) CheckEmailSendRateLimit(email string) error {
	return s.CheckRateLimit(fmt.Sprintf("email_send:%s", email), 5, 300) // 5 per 5 minutes
}

func (s *AccountService) CheckLoginRateLimit(email string) error {
	return s.CheckRateLimit(fmt.Sprintf("login:%s", email), 10, 600) // 10 per 10 minutes
}

func (s *AccountService) CheckPasswordResetRateLimit(email string) error {
	return s.CheckRateLimit(fmt.Sprintf("pwd_reset:%s", email), 3, 3600) // 3 per hour
}

// === Verification Codes ===

func (s *AccountService) GenerateVerificationCode(purpose, identifier string) (string, error) {
	code := fmt.Sprintf("%06d", rand.Intn(1000000))
	cacheKey := fmt.Sprintf("verify:%s:%s", purpose, identifier)
	client := cache.Instance().Client()
	client.Set(context.Background(), cacheKey, code, 10*time.Minute)
	return code, nil
}

func (s *AccountService) VerifyCode(purpose, identifier, code string) bool {
	cacheKey := fmt.Sprintf("verify:%s:%s", purpose, identifier)
	client := cache.Instance().Client()
	stored, err := client.Get(context.Background(), cacheKey).Result()
	if err != nil || stored != code {
		return false
	}
	client.Del(context.Background(), cacheKey)
	return true
}

// === Owner Transfer ===

func (s *AccountService) TransferOwnership(tenantID, currentOwnerID, newOwnerID string) error {
	// Verify current owner
	role := s.GetUserRole(tenantID, currentOwnerID)
	if role != "owner" {
		return fmt.Errorf("only the owner can transfer ownership")
	}
	// Update roles
	if err := s.UpdateMemberRole(tenantID, currentOwnerID, "admin"); err != nil {
		return err
	}
	return s.UpdateMemberRole(tenantID, newOwnerID, "owner")
}

// === Dataset Operator Members ===

// GetDatasetOperatorMembers returns accounts that can operate on datasets in a tenant.
func (s *AccountService) GetDatasetOperatorMembers(tenantID string) []*models.Account {
	var accounts []*models.Account
	dbengine.Instance().DB.Raw(
		"SELECT acc.* FROM accounts acc JOIN tenant_account_joins ta ON acc.id = ta.account_id WHERE ta.tenant_id = ? AND ta.role IN (?, ?, ?, ?)",
		tenantID,
		string(enumtypes.TenantAccountRole_OWNER),
		string(enumtypes.TenantAccountRole_ADMIN),
		string(enumtypes.TenantAccountRole_EDITOR),
		string(enumtypes.TenantAccountRole_DATASET_OPERATOR),
	).Find(&accounts)
	return accounts
}

// === Invitation Code ===

// GetInvitationByCode looks up an invitation code record by code string.
func (s *AccountService) GetInvitationByCode(code string) *models.InvitationCode {
	var invitation models.InvitationCode
	if err := dbengine.Instance().DB.Where("code = ? AND status = ?", code, "unused").First(&invitation).Error; err != nil {
		return nil
	}
	return &invitation
}

// UseInvitationCode marks an invitation code as used by a tenant and account.
func (s *AccountService) UseInvitationCode(code string, tenantID, accountID string) error {
	invitation := s.GetInvitationByCode(code)
	if invitation == nil {
		return fmt.Errorf("invitation code not found or already used")
	}
	now := time.Now()
	return dbengine.Instance().DB.Model(invitation).Updates(map[string]interface{}{
		"status":             "used",
		"used_at":            &now,
		"used_by_tenant_id":  tenantID,
		"used_by_account_id": accountID,
	}).Error
}

// === Token Refresh ===

// RefreshAccessToken validates a refresh token and issues a new token pair.
func (s *AccountService) RefreshAccessToken(refreshToken string) (*TokenPair, error) {
	accountID := cache.Instance().GetString(s._get_refresh_token_key(refreshToken))
	if accountID == "" {
		return nil, fmt.Errorf("invalid or expired refresh token")
	}
	account, err := s.Get(accountID)
	if err != nil || account == nil {
		return nil, fmt.Errorf("account not found")
	}
	// Delete old refresh token
	s._delete_refresh_token(refreshToken, accountID)
	// Issue new pair
	exp := 30 * time.Second * 3600 * 24
	accessToken := s.GetAccountJWTToken(account, exp)
	newRefreshToken := password.GenerateRefreshToken(0)
	s._store_refresh_token(newRefreshToken, accountID)
	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

// === Batch Operations ===

// GetAccountsByIDs returns multiple accounts by their IDs.
func (s *AccountService) GetAccountsByIDs(ids []string) []*models.Account {
	var accounts []*models.Account
	if len(ids) == 0 {
		return accounts
	}
	dbengine.Instance().DB.Where("id IN ?", ids).Find(&accounts)
	return accounts
}

// CountTenantMembers returns the number of members in a tenant.
func (s *AccountService) CountTenantMembers(tenantID string) int64 {
	var count int64
	dbengine.Instance().DB.Model(&models.TenantAccountJoin{}).Where("tenant_id = ?", tenantID).Count(&count)
	return count
}

// SwitchCurrentTenant sets the current flag for a user's tenant.
func (s *AccountService) SwitchCurrentTenant(accountID, tenantID string) error {
	// Unset all current flags for this account
	dbengine.Instance().DB.Model(&models.TenantAccountJoin{}).
		Where("account_id = ?", accountID).
		Update("current", false)
	// Set the target tenant as current
	result := dbengine.Instance().DB.Model(&models.TenantAccountJoin{}).
		Where("account_id = ? AND tenant_id = ?", accountID, tenantID).
		Update("current", true)
	if result.RowsAffected == 0 {
		return fmt.Errorf("tenant membership not found")
	}
	return result.Error
}
