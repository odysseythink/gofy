package services

import (
	crand "crypto/rand"
	"encoding/base64"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"mlib.com/confy"
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
	"mlib.com/mlog"
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
		ID:                uuid.NewV4().String(),
		Name:              name,
		Email:             email,
		InterfaceLanguage: interface_language,
		InterfaceTheme:    interface_theme,
		Timezone:          timezone,
		LastLoginAt:       &now,
		InitializedAt:     &now,
		CreatedAt:         &now,
		UpdatedAt:         &now,
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
			dbengine.Instance().DB.Updates(&models.TenantAccountJoin{ID: result.TaID, Current: true})
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
