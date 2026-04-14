package v1

import (
	"net/http"
	"os"

	dbengine "mlib.com/gofy/server/db_engine"
	"mlib.com/gofy/server/global"
	"mlib.com/gofy/server/models"
	"mlib.com/gofy/server/services"

	"github.com/gin-gonic/gin"
	"github.com/odysseythink/confy"
	"github.com/odysseythink/mlog"
)

type InitValidateApi struct {
}

// GetInitValidateStatus returns the initialization validation status.
// GET /console/api/init
func (api *InitValidateApi) GetInitValidateStatus(c *gin.Context) {
	status := getInitValidateStatus()
	if status {
		c.JSON(http.StatusOK, gin.H{"status": "finished"})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "not_started"})
	}
}

// InitValidate validates the init password.
// POST /console/api/init
func (api *InitValidateApi) InitValidate(c *gin.Context) {
	// Only for self-hosted edition
	if confy.Get[string]("EDITION") != "SELF_HOSTED" {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    "not_allowed",
			"message": "only self-hosted edition supports this endpoint",
		})
		return
	}

	tenant_count := services.ServiceGroupApp.Tenant.GetTenantCount()
	if tenant_count > 0 {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    "already_setup",
			"message": "The system has been set up. Please log in.",
		})
		return
	}

	var req struct {
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "invalid_request",
			"message": "missing password",
		})
		return
	}

	initPassword := os.Getenv("INIT_PASSWORD")
	if req.Password != initPassword {
		global.IsInitValidated = false
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    "init_validate_failed",
			"message": "Invalid init password.",
		})
		return
	}

	global.IsInitValidated = true
	c.JSON(http.StatusCreated, gin.H{"result": "success"})
}

func getInitValidateStatus() bool {
	if confy.Get[string]("EDITION") == "SELF_HOSTED" {
		if os.Getenv("INIT_PASSWORD") != "" {
			if global.IsInitValidated {
				return true
			}
			setup := new(models.GofySetup)
			err := dbengine.Instance().DB.First(setup).Error
			return err == nil
		}
	}
	return true
}

func init() {
	mlog.Debug("init_validate api registered")
}
