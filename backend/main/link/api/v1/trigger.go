package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/odysseythink/mlog"
	"mlib.com/gofy/server/models"
	"mlib.com/gofy/server/services"
)

type TriggerApi struct{}

// GetAppTriggers lists triggers for an app
func (api *TriggerApi) GetAppTriggers(c *gin.Context) {
	appID := c.Param("app_id")
	rawUser, _ := c.Get("current_user")
	user, ok := rawUser.(*models.Account)
	if !ok || user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	triggers := services.ServiceGroupApp.Trigger.GetAppTriggers(user.CurrentTenantID(), appID)
	c.JSON(http.StatusOK, gin.H{"data": triggers})
}

// CreateAppTrigger creates a trigger for an app
func (api *TriggerApi) CreateAppTrigger(c *gin.Context) {
	appID := c.Param("app_id")
	rawUser, _ := c.Get("current_user")
	user, ok := rawUser.(*models.Account)
	if !ok || user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req struct {
		NodeID       string `json:"node_id" binding:"required"`
		TriggerType  string `json:"trigger_type" binding:"required"`
		Title        string `json:"title" binding:"required"`
		ProviderName string `json:"provider_name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	trigger, err := services.ServiceGroupApp.Trigger.CreateAppTrigger(
		user.CurrentTenantID(), appID, req.NodeID, req.TriggerType, req.Title, req.ProviderName,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": trigger})
}

// DeleteAppTrigger removes an app trigger
func (api *TriggerApi) DeleteAppTrigger(c *gin.Context) {
	triggerID := c.Param("trigger_id")
	if err := services.ServiceGroupApp.Trigger.DeleteAppTrigger(triggerID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": "success"})
}

// EnableAppTrigger enables or disables a trigger
func (api *TriggerApi) EnableAppTrigger(c *gin.Context) {
	appID := c.Param("app_id")
	_ = appID
	var req struct {
		TriggerID string `json:"trigger_id" binding:"required"`
		Enabled   bool   `json:"enabled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.ServiceGroupApp.Trigger.EnableAppTrigger(req.TriggerID, req.Enabled); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": "success"})
}

// GetTriggerLogs returns trigger execution logs
func (api *TriggerApi) GetTriggerLogs(c *gin.Context) {
	appID := c.Param("app_id")
	rawUser, _ := c.Get("current_user")
	user, ok := rawUser.(*models.Account)
	if !ok || user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	logs, total := services.ServiceGroupApp.Trigger.GetTriggerLogs(user.CurrentTenantID(), appID, page, pageSize)
	c.JSON(http.StatusOK, gin.H{
		"data":  logs,
		"total": total,
		"page":  page,
		"limit": pageSize,
	})
}

// GetWebhookTrigger returns webhook trigger info for a workflow node
func (api *TriggerApi) GetWebhookTrigger(c *gin.Context) {
	appID := c.Param("app_id")
	nodeID := c.Query("node_id")
	rawUser, _ := c.Get("current_user")
	user, ok := rawUser.(*models.Account)
	if !ok || user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	webhook := services.ServiceGroupApp.Trigger.GetWebhookTrigger(appID, nodeID, user.CurrentTenantID(), user.ID)
	if webhook == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "webhook trigger not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": webhook})
}

// CreateSchedulePlan creates a cron schedule for a trigger
func (api *TriggerApi) CreateSchedulePlan(c *gin.Context) {
	appID := c.Param("app_id")
	rawUser, _ := c.Get("current_user")
	user, ok := rawUser.(*models.Account)
	if !ok || user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req struct {
		NodeID         string `json:"node_id" binding:"required"`
		CronExpression string `json:"cron_expression" binding:"required"`
		Timezone       string `json:"timezone"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Timezone == "" {
		req.Timezone = "UTC"
	}

	plan, err := services.ServiceGroupApp.Trigger.CreateSchedulePlan(appID, req.NodeID, user.CurrentTenantID(), req.CronExpression, req.Timezone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": plan})
}

// GetSchedulePlans returns schedule plans for an app
func (api *TriggerApi) GetSchedulePlans(c *gin.Context) {
	appID := c.Param("app_id")
	rawUser, _ := c.Get("current_user")
	user, ok := rawUser.(*models.Account)
	if !ok || user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	plans := services.ServiceGroupApp.Trigger.GetSchedulePlans(appID, user.CurrentTenantID())
	c.JSON(http.StatusOK, gin.H{"data": plans})
}

// ListTriggerProviders lists available trigger providers for the workspace
func (api *TriggerApi) ListTriggerProviders(c *gin.Context) {
	// Return the built-in trigger types
	// Plugin providers would be fetched from the plugin daemon in a full implementation
	providers := []map[string]any{
		{
			"type":        "webhook",
			"name":        "Webhook",
			"description": "Trigger workflow via HTTP webhook",
		},
		{
			"type":        "schedule",
			"name":        "Schedule",
			"description": "Trigger workflow on a cron schedule",
		},
		{
			"type":        "plugin",
			"name":        "Plugin",
			"description": "Trigger workflow via plugin events",
		},
	}
	c.JSON(http.StatusOK, gin.H{"data": providers})
}

// ListSubscriptions lists trigger subscriptions for a provider
func (api *TriggerApi) ListSubscriptions(c *gin.Context) {
	provider := c.Param("provider")
	rawUser, _ := c.Get("current_user")
	user, ok := rawUser.(*models.Account)
	if !ok || user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	subs := services.ServiceGroupApp.Trigger.ListSubscriptions(user.CurrentTenantID(), provider)
	c.JSON(http.StatusOK, gin.H{"data": subs})
}

// CreateSubscription creates a new trigger subscription
func (api *TriggerApi) CreateSubscription(c *gin.Context) {
	provider := c.Param("provider")
	rawUser, _ := c.Get("current_user")
	user, ok := rawUser.(*models.Account)
	if !ok || user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req struct {
		Name           string         `json:"name" binding:"required"`
		Parameters     map[string]any `json:"parameters"`
		Properties     map[string]any `json:"properties"`
		Credentials    map[string]any `json:"credentials"`
		CredentialType string         `json:"credential_type"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sub, err := services.ServiceGroupApp.Trigger.CreateSubscription(
		user.CurrentTenantID(), user.ID, provider, req.Name,
		req.Parameters, req.Properties, req.Credentials, req.CredentialType,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": sub})
}

// DeleteSubscription removes a trigger subscription
func (api *TriggerApi) DeleteSubscription(c *gin.Context) {
	subscriptionID := c.Param("subscription_id")
	if err := services.ServiceGroupApp.Trigger.DeleteSubscription(subscriptionID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": "success"})
}

// HandlePluginIncoming is the public endpoint that receives plugin trigger events
func (api *TriggerApi) HandlePluginIncoming(c *gin.Context) {
	endpointID := c.Param("endpoint_id")
	if endpointID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "endpoint_id required"})
		return
	}

	data, err := extractWebhookData(c)
	if err != nil {
		mlog.Errorf("failed to extract plugin event data: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.ServiceGroupApp.Trigger.HandlePluginEvent(endpointID, data); err != nil {
		mlog.Errorf("failed to handle plugin event for endpoint %s: %v", endpointID, err)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "success"})
}

// HandleWebhookIncoming is the public endpoint that receives webhook calls
// This does NOT require authentication - it's called by external services
func (api *TriggerApi) HandleWebhookIncoming(c *gin.Context) {
	webhookID := c.Param("webhook_id")
	if webhookID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "webhook_id required"})
		return
	}

	// Extract webhook data from the request
	data, err := extractWebhookData(c)
	if err != nil {
		mlog.Errorf("failed to extract webhook data: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.ServiceGroupApp.Trigger.HandleWebhook(webhookID, data); err != nil {
		mlog.Errorf("failed to handle webhook %s: %v", webhookID, err)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "success"})
}
