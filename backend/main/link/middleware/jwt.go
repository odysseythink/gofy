package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/odysseythink/confy"
	"github.com/odysseythink/gofy/backend/core/exceptions"
	httpexceptions "github.com/odysseythink/gofy/backend/core/exceptions/http"
	dbengine "github.com/odysseythink/gofy/backend/db_engine"
	enumtypes "github.com/odysseythink/gofy/backend/enum_types"
	"github.com/odysseythink/gofy/backend/models"
	jwtutils "github.com/odysseythink/gofy/backend/utils/jwt"
	"github.com/odysseythink/mlog"
)

func LoadUser(account_id string) (*models.Account, error) {
	account := new(models.Account)
	err := dbengine.Instance().DB.Where("id = ?", account_id).First(account).Error
	if err != nil {
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
		dbengine.Instance().DB.Updates(&models.Account{Model: models.Model{ID: account_id}, LastActiveAt: &now})
	}
	return account, nil
}

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 我们这里jwt鉴权取头部信息 x-token 登录时回返回token信息 这里前端需要把token存储到cookie或者本地localStorage中 不过需要跟后端协商过期时间 可以约定刷新令牌或者重新登录
		token := jwtutils.GetToken(c)
		if token == "" {
			mlog.Errorf("未登录或非法访问")
			c.JSON(401, gin.H{"code": "unauthorized", "message": "Unauthorized."})
			c.Abort()
			return
		}

		// parseToken 解析token包含的信息
		claims, err := jwtutils.ParseToken(token, confy.Get[string]("SECRET_KEY"))
		if err != nil || claims == nil {
			mlog.Errorf("parse token failed:%v", err)
			c.JSON(401, gin.H{"code": "unauthorized", "message": "Unauthorized."})
			c.Abort()
			return
		}
		current_user, exp := LoadUser(claims.UserID)
		if exp != nil {
			mlog.Errorf("load user(%s) failed:%v", claims.UserID, exp.Error())
			c.JSON(401, gin.H{"code": "unauthorized", "message": "Unauthorized."})
			c.Abort()
			return
		}
		if current_user.Status == "uninitialized" {
			(&httpexceptions.AccountNotInitializedError{}).Response(c)
			c.Abort()
		}
		c.Set("current_user", current_user)

		c.Next()

		if newToken, exists := c.Get("new-token"); exists {
			c.Header("new-token", newToken.(string))
		}
		if newExpiresAt, exists := c.Get("new-expires-at"); exists {
			c.Header("new-expires-at", newExpiresAt.(string))
		}
	}
}
