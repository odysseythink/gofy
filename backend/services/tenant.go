package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/odysseythink/confy"
	"github.com/odysseythink/mlog"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"mlib.com/gofy/server/core/exceptions"
	httpexceptions "mlib.com/gofy/server/core/exceptions/http"
	dbengine "mlib.com/gofy/server/db_engine"
	enumtypes "mlib.com/gofy/server/enum_types"
	"mlib.com/gofy/server/events"
	"mlib.com/gofy/server/models"
	"mlib.com/gofy/server/models/request"
	cryptutils "mlib.com/gofy/server/utils/crypt"
)

type TenantService struct {
}

// Create 创建Tenant记录
func (s *TenantService) Create(te models.Tenant) (err error) {
	err = dbengine.Instance().DB.Create(&te).Error
	return err
}

// Delete 删除Tenant记录
func (s *TenantService) Delete(te models.Tenant) (err error) {
	err = dbengine.Instance().DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Tenant{}).Where("id = ?", te.ID).Error; err != nil {
			return err
		}
		if err = tx.Delete(&te).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// DeleteByIds 批量删除Tenant记录
func (s *TenantService) DeleteByIds(ids []string) (err error) {
	if err := dbengine.Instance().DB.Model(&models.Tenant{}).Delete("id in ?", ids).Error; err != nil {
		return err
	}

	return nil
}

// Update 更新Tenant记录
func (s *TenantService) Update(te *models.Tenant) (err error) {
	err = dbengine.Instance().DB.Save(te).Error
	return err
}

// Get 根据id获取Tenant记录
func (s *TenantService) Get(id string) (te *models.Tenant, err error) {
	te = &models.Tenant{}
	err = dbengine.Instance().DB.Where("id = ?", id).First(te).Error
	if err != nil {
		te = nil
	}
	return
}

// GetList 分页获取Tenant记录
func (s *TenantService) GetList(info request.BasePageReq) (list []models.Tenant, total int64, err error) {
	limit := info.Page
	offset := info.Limit * (info.Page - 1)
	// 创建db
	db := dbengine.Instance().DB.Model(&models.Tenant{})
	var ts []models.Tenant

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	err = db.Limit(limit).Offset(offset).Find(&ts).Error
	return ts, total, err
}

func (s *TenantService) GetJoinTenants(account *models.Account) []*models.Tenant {
	if account == nil {
		mlog.Errorf("invalid arg")
		return nil
	}
	mlog.Debugf("--account=%#v", account)
	// """Get acc join tenants"""
	var ts []*models.Tenant
	err := dbengine.Instance().DB.Raw("SELECT t.*, ta.role, ta.id as ta_id FROM tenants t join tenant_account_joins ta on ta.tenant_id = t.id WHERE ta.account_id =? AND t.status = ?;", account.ID, enumtypes.TenantStatus_NORMAL).Find(&ts).Error
	if err != nil {
		mlog.Errorf("get Tenant from mysql failed:%v", err)
		return nil
	}

	return ts
}

func (s *TenantService) GetCurrentTenantByAccount(acc *models.Account) (*models.Tenant, error) {
	// """Get tenant by account and add the role"""
	if acc == nil {
		mlog.Errorf("invalid arg")
		return nil, errors.New("invalid arg")
	}

	if acc.CurrentTenant == nil {
		return nil, errors.New("tenant not found")
	}
	ta := &models.TenantAccountJoin{}
	err := dbengine.Instance().DB.Model(&models.TenantAccountJoin{}).Where("tenant_id = ? and account_id = ?", acc.CurrentTenant().ID, acc.ID).First(ta).Error
	if err != nil {
		acc.CurrentTenant().CurrentRole = ta.Role
	} else {
		return nil, errors.New("tenant not found for the account")
	}
	return acc.CurrentTenant(), nil
}

func (s *TenantService) RestCount(t *models.Tenant) int64 {
	if t == nil {
		mlog.Errorf("invalid arg")
		return 0
	}

	var count int64 = 0
	err := dbengine.Instance().DB.Model(&models.Tenant{}).Where("created_at < ? and id != ?", t.CreatedAt, t.ID).Count(&count).Error
	if err != nil {
		mlog.Errorf("get rest failed:%v", err)
		return 0
	} else {
		return count
	}
}

func (s *TenantService) HasRoles(tenant *models.Tenant, roles []enumtypes.TenantAccountJoinRole) bool {
	/*Check if user has any of the given roles for a tenant*/
	str_roles := []string{}
	for _, role := range roles {
		if !role.Valid() {
			panic(exceptions.NewValueError("all roles must be TenantAccountJoinRole"))
		}
		str_roles = append(str_roles, string(role))
	}
	ta := new(models.TenantAccountJoin)
	err := dbengine.Instance().DB.Model(&models.TenantAccountJoin{}).Where("tenant_id=? and role in ?", tenant.ID, str_roles).First(ta).Error
	if err != nil {
		mlog.Errorf("get TenantAccountJoin from mysql failed:%v", err)
		ta = nil
	}
	return ta != nil
}

func (s *TenantService) GetTenantInfo(tenant *models.Tenant, current_user *models.Account) map[string]any {
	if tenant == nil {
		mlog.Errorf("invalid arg")
		return nil
	}
	tenant_info := map[string]any{
		"id":               tenant.ID,
		"name":             tenant.Name,
		"plan":             tenant.Plan,
		"status":           tenant.Status,
		"created_at":       tenant.CreatedAt,
		"in_trail":         true,
		"trial_end_reason": "",
		"role":             "normal",
	}

	// # Get role of user
	ta := &models.TenantAccountJoin{}
	err := dbengine.Instance().DB.Model(&models.TenantAccountJoin{}).Where("tenant_id = ? and account_id = ?", tenant.ID, current_user.ID).First(ta).Error
	if err != nil {
		mlog.Errorf("get TenantAccountJoin by tenant=%#v, current_user=%#v failed:%v", tenant, current_user, err)
		return nil
	}

	tenant_info["role"] = ta.Role

	if confy.Get[bool]("CAN_REPLACE_LOGO") && s.HasRoles(tenant, []enumtypes.TenantAccountJoinRole{enumtypes.TenantAccountJoinRole_OWNER, enumtypes.TenantAccountJoinRole_ADMIN}) {
		base_url := confy.Get[string]("FILES_URL")
		replace_webapp_logo := ""
		remove_webapp_brand := false
		custom_config_dict := make(map[string]any)
		if tenant.CustomConfig != "" {
			err := json.Unmarshal([]byte(tenant.CustomConfig), &custom_config_dict)
			if err != nil {
				mlog.Warningf("tenant(%#v) CustomConfig must be dict", tenant)
			} else {
				if _, ok := custom_config_dict["replace_webapp_logo"]; ok {
					if tmp, ok := custom_config_dict["replace_webapp_logo"].(bool); ok && tmp {
						replace_webapp_logo = fmt.Sprintf("%s/files/workspaces/%s/webapp-logo", base_url, tenant.ID)
					}
				}
				if _, ok := custom_config_dict["remove_webapp_brand"]; ok {
					if tmp, ok := custom_config_dict["remove_webapp_brand"].(bool); ok {
						remove_webapp_brand = tmp
					}
				}
			}
		}

		tenant_info["custom_config"] = map[string]any{
			"remove_webapp_brand": remove_webapp_brand,
			"replace_webapp_logo": replace_webapp_logo,
		}
	}
	return tenant_info
}

func (s *TenantService) GetTenantMembers(tenant_id string) []*models.Account {
	type Result struct {
		models.Account
		Role enumtypes.TenantAccountJoinRole `gorm:"column:role;type:varchar(16);default:normal" json:"role"`
	}
	var accs []*Result
	sql := fmt.Sprintf("SELECT accounts.*,tenant_account_joins.role  FROM accounts JOIN tenant_account_joins ON tenant_account_joins.account_id = account.id AND tenant_account_joins.tenant_id = '%s';", tenant_id)
	err := dbengine.Instance().Raw(sql).Scan(&accs).Error
	if err != nil {
		mlog.Errorf("exce sql(%s) failed:%v", sql, err)
		return nil
	}

	// Initialize an empty list to store the updated accounts
	updated_accounts := []*models.Account{}

	for _, v := range accs {
		v.Account.SetRole(string(v.Role))
	}
	return updated_accounts
}
func (s *TenantService) GetTenantCount() int64 {
	var count int64 = 0
	err := dbengine.Instance().DB.Model(&models.Tenant{}).Count(&count).Error
	if err != nil {
		mlog.Errorf("count failed:%v", err)
		return 0
	} else {
		return count
	}
}

func (s *TenantService) CreateTenant(name string, is_setup bool, is_from_dashboard bool) *models.Tenant {
	/*Create tenant*/
	if !ServiceGroupApp.Feature.GetSystemFeatures().IsAllowCreateWorkspace && !is_setup && !is_from_dashboard {
		panic(httpexceptions.NewNotAllowedCreateWorkspace())
	}
	now := time.Now()
	tenant := &models.Tenant{
		Model: models.Model{
			ID:        uuid.NewV4().String(),
			CreatedAt: &now,
			UpdatedAt: &now,
		},
		Name: name,
	}

	tenant.EncryptPublicKey = cryptutils.GenerateKeyPair(tenant.ID)
	dbengine.Instance().Create(tenant)
	return tenant
}

func (s *TenantService) CreateTenantMember(tenant *models.Tenant, account *models.Account, role string) *models.TenantAccountJoin {
	/*Create tenant member*/
	if role == "" {
		role = "normal"
	}
	if enumtypes.TenantAccountJoinRole(role) == enumtypes.TenantAccountJoinRole_OWNER {
		if s.HasRoles(tenant, []enumtypes.TenantAccountJoinRole{enumtypes.TenantAccountJoinRole_OWNER}) {
			mlog.Errorf("Tenant {%s} has already an owner.", tenant.ID)
			panic(errors.New("tenant already has an owner."))
		}
	}
	ta := new(models.TenantAccountJoin)
	err := dbengine.Instance().DB.Model(&models.TenantAccountJoin{}).Where("tenant_id=? and account_id = ?", tenant.ID, account.ID).First(ta).Error
	if err != nil {
		mlog.Errorf("get TenantAccountJoin from mysql failed:%v", err)
		ta = nil
	}

	if ta != nil {
		ta.Role = role
		dbengine.Instance().DB.Updates(&models.TenantAccountJoin{Model: models.Model{ID: ta.ID}, Role: role})
	} else {
		now := time.Now()
		ta = &models.TenantAccountJoin{
			Model: models.Model{
				ID:        uuid.NewV4().String(),
				CreatedAt: &now,
				UpdatedAt: &now,
			},
			TenantID:  tenant.ID,
			AccountID: account.ID,
			Role:      role,
		}
		dbengine.Instance().DB.Create(ta)
	}

	return ta
}
func (s *TenantService) CreateOwnerTenantIfNotExist(
	account *models.Account, name string, is_setup bool,
) {
	/*Check if user have a workspace or not*/
	available_ta := new(models.TenantAccountJoin)
	err := dbengine.Instance().DB.Model(&models.TenantAccountJoin{}).Where("account_id = ?", account.ID).Order("id ASC").First(available_ta).Error
	if err != nil {
		mlog.Errorf("get TenantAccountJoin from mysql failed:%v", err)
		available_ta = nil
	}

	if available_ta != nil {
		return
	}
	/*Create owner tenant if not exist*/
	if !ServiceGroupApp.Feature.GetSystemFeatures().IsAllowCreateWorkspace && !is_setup {
		panic(exceptions.NewWorkSpaceNotAllowedCreateError())
	}
	var tenant *models.Tenant
	if name != "" {
		tenant = s.CreateTenant(name, is_setup, false)
	} else {
		tenant = s.CreateTenant(fmt.Sprintf("%s's Workspace", account.Name), is_setup, false)
	}
	s.CreateTenantMember(tenant, account, "owner")
	account.SetCurrentTenant(tenant)
	events.Instance.TenantWasCreatedSig.Emit(tenant)
}

func (s *TenantService) SwitchTenant(account *models.Account, tenant_id string) {
	// """Switch the current workspace for the account"""

	// # Ensure tenant_id is provided
	if account == nil {
		mlog.Errorf("account must be provided")
		panic(exceptions.NewValueError("account must be provided"))
	}
	if tenant_id == "" {
		mlog.Errorf("tenant ID must be provided")
		panic(exceptions.NewValueError("tenant ID must be provided"))
	}

	tenant_account_join := new(models.TenantAccountJoin)
	err := dbengine.Instance().DB.Raw("SELECT ta.* FROM tenant_account_joins ta join tenants t on ta.tenant_id = t.id WHERE ta.account_id = ? and ta.tenant_id = ? and t.status = ?;", account.ID, tenant_id, enumtypes.TenantStatus_NORMAL).First(tenant_account_join).Error
	if err != nil {
		mlog.Errorf("get TenantAccountJoin by account=%#v, tenant_id=%s failed:%v", account, tenant_id, err)
		tenant_account_join = nil
	}
	if tenant_account_join == nil {
		panic(httpexceptions.NewAccountNotLinkTenantError("Tenant not found or account is not a member of the tenant."))
	} else {
		dbengine.Instance().DB.Model(&models.TenantAccountJoin{}).Update("current", false).Where("account_id = ? and tenant_id=?", account.ID, tenant_id)
		tenant_account_join.Current = true
		dbengine.Instance().DB.Updates(&models.TenantAccountJoin{Model: models.Model{ID: tenant_account_join.ID}, Current: true})

		// Set the current tenant for the account
		account.SetCurrentTenantID(tenant_account_join.TenantID)
	}
}
