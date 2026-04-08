package services

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"mlib.com/confy"
	"mlib.com/gofy/server/cache"
	"mlib.com/gofy/server/core/exceptions"
	dbengine "mlib.com/gofy/server/db_engine"
	"mlib.com/gofy/server/global"
	"mlib.com/gofy/server/models"
	"mlib.com/gofy/server/utils/validate"
	"mlib.com/mlog"
)

type RegisterService struct {
}

func (s *RegisterService) getInvitationTokenKey(token string) string {
	return fmt.Sprintf("member_invite:token:", token)
}

func (s *RegisterService) getInvitationByToken(token, workspace_id, email string) map[string]string {
	if workspace_id != "" && email != "" {
		hasher := sha256.New()
		hasher.Write([]byte(email))
		email_hash := hex.EncodeToString(hasher.Sum(nil))
		// email_hash = sha256(email.encode()).hexdigest()
		cache_key := fmt.Sprintf("member_invite_token:%s, %s:%s", workspace_id, email_hash, token)
		account_id := cache.Instance().GetString(cache_key)
		// account_id = redis_client.get(cache_key)

		if account_id == "" {
			return nil
		}

		return map[string]string{
			"account_id":   account_id,
			"email":        email,
			"workspace_id": workspace_id,
		}
	} else {
		data := cache.Instance().GetString(s.getInvitationTokenKey(token))
		if data == "" {
			return nil
		}
		invitation := map[string]string{}
		if err := json.Unmarshal([]byte(data), &invitation); err != nil {
			mlog.Errorf("json unmarshal(%s) failed:%v", data, err)
			return nil
		} else {
			return invitation
		}
	}
}
func (s *RegisterService) GetInvitationIfTokenValid(workspace_id, email, token string) map[string]any {
	invitation_data := s.getInvitationByToken(token, workspace_id, email)
	if invitation_data == nil {
		return nil
	}
	err := validate.StringMapTypeVerify(invitation_data, validate.Rules{
		"workspace_id": {validate.RuleTypeOfField(reflect.String), validate.NotEmpty()},
		"email":        {validate.RuleTypeOfField(reflect.String), validate.NotEmpty()},
		"account_id":   {validate.RuleTypeOfField(reflect.String), validate.NotEmpty()},
	})
	if err != nil {
		mlog.Errorf("verify failed:%v", err)
		return nil
	}
	var tenant models.Tenant
	if err := dbengine.Instance().DB.Where("id = ? and status = ? ", invitation_data["workspace_id"], "normal").First(&tenant); err != nil {
		mlog.Errorf("get Tenant failed:%v", err)
		return nil
	}

	var tenant_account struct {
		*models.Account
		Role string `gorm:"column:role;" json:"role"`
	}
	err = dbengine.Instance().DB.Raw("SELECT acc.*, ta.role FROM accounts acc join tenant_account_joins ta on ta.account_id = acc.id WHERE acc.email = ? AND ta.tenant_id = ?;", invitation_data["email"], tenant.ID).First(&tenant_account).Error
	if err != nil {
		mlog.Errorf("get account failed:%v", err)
		return nil
	}
	if invitation_data["account_id"] != tenant_account.ID {
		return nil
	}

	return map[string]any{
		"account": tenant_account,
		"data":    invitation_data,
		"tenant":  tenant,
	}
}

func (cls *RegisterService) Setup(email string, name string, password string, ip_address string) {
	/*
	   Setup dify

	   :param email: email
	   :param name: username
	   :param password: password
	   :param ip_address: ip address
	*/
	defer func() {
		if r := recover(); r != nil {
			// db.session.query(DifySetup).delete()
			// db.session.query(TenantAccountJoin).delete()
			// db.session.query(Account).delete()
			// db.session.query(Tenant).delete()
			// db.session.commit()

			mlog.Errorf("Setup account failed, email: {%s}, name: {%s}", email, name)
			panic(exceptions.NewValueError(fmt.Sprintf("Setup failed: %v", r)))
		}
		// Register
		account := ServiceGroupApp.Account.CreateAccount(
			email,
			name,
			global.LANGUAGES[0],
			password,
			"",
			true,
		)
		now := time.Now()
		account.LastLoginIP = ip_address
		account.InitializedAt = &now

		ServiceGroupApp.Tenant.CreateOwnerTenantIfNotExist(account, "", true)

		dify_setup := &models.DifySetup{
			Version: confy.Get[string]("CURRENT_VERSION"),
			SetupAt: &now,
		}
		dbengine.Instance().DB.Create(dify_setup)
	}()
}
