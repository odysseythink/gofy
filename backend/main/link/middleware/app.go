package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	dbengine "github.com/odysseythink/gofy/backend/db_engine"
	enumtypes "github.com/odysseythink/gofy/backend/enum_types"
	"github.com/odysseythink/gofy/backend/models"
	jwtutils "github.com/odysseythink/gofy/backend/utils/jwt"
	"github.com/odysseythink/mlog"
)

func AppAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := jwtutils.GetToken(c)
		if token == "" {
			mlog.Errorf("非法访问")
			c.JSON(401, gin.H{"code": "unauthorized", "message": "Unauthorized."})
			c.Abort()
			return
		}
		current_time := time.Now()
		api_token := new(models.ApiToken)
		err := dbengine.Instance().DB.Model(&models.ApiToken{}).Where("token = ? and type = ?", token, "app").First(api_token).Error
		if err != nil {
			mlog.Errorf("get api token=%s failed:%v", token, err)
			c.JSON(401, gin.H{"code": "unauthorized", "message": "Access token is invalid."})
			c.Abort()
			return
		}

		app_model := new(models.App)
		err = dbengine.Instance().DB.Model(&models.App{}).Where("id = ?", api_token.AppID).First(app_model).Error
		if err != nil {
			mlog.Errorf("get App failed:%v", err)
			c.JSON(http.StatusForbidden, gin.H{"code": "forbidden", "message": "You don't have the permission to access the requested resource. The app no longer exists."})
			c.Abort()
			return
		}

		if app_model.Status != "normal" {
			mlog.Errorf("The app's status is abnormal.")
			c.JSON(http.StatusForbidden, gin.H{"code": "forbidden", "message": "You don't have the permission to access the requested resource. The app's status is abnormal."})
			c.Abort()
			return
		}
		if !app_model.EnableAPI {
			mlog.Errorf("The app's API service has been disabled.")
			c.JSON(http.StatusForbidden, gin.H{"code": "forbidden", "message": "You don't have the permission to access the requested resource. The app's API service has been disabled."})
			c.Abort()
			return
		}
		tenant := new(models.Tenant)
		err = dbengine.Instance().DB.Model(&models.Tenant{}).Where("id = ?", app_model.TenantID).First(tenant).Error
		if err != nil {
			mlog.Errorf("get tenant failed:%v", err)
			c.JSON(http.StatusForbidden, gin.H{"code": "forbidden", "message": "You don't have the permission to access the requested resource. Tenant does not exist."})
			c.Abort()
			return
		}

		if tenant.Status == enumtypes.TenantStatus_ARCHIVE {
			mlog.Errorf("The workspace's status is archived.")
			c.JSON(http.StatusForbidden, gin.H{"code": "forbidden", "message": "You don't have the permission to access the requested resource. The workspace's status is archived."})
			c.Abort()
			return
		}
		c.Set("app_model", app_model)

		api_token.LastUsedAt = &current_time
		dbengine.Instance().DB.Updates(&models.ApiToken{ID: api_token.ID, LastUsedAt: api_token.LastUsedAt})
		c.Next()
	}
}
