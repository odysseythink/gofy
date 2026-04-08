package v1

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	dbengine "mlib.com/gofy/server/db_engine"
	"mlib.com/gofy/server/models"
	"mlib.com/gofy/server/storage"
)

type WebApi struct{}

// GetConversations returns conversations for an app user.
func (a *WebApi) GetConversations(c *gin.Context) {
	appCode := c.GetHeader("X-App-Code")
	if appCode == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "X-App-Code header required"})
		return
	}

	// Find app by app code (via Site model)
	var site models.Site
	if err := dbengine.Instance().DB.Where("code = ?", appCode).First(&site).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "app not found"})
		return
	}

	// Get end user from passport
	userID := c.GetString("end_user_id")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if limit > 100 {
		limit = 100
	}

	var conversations []*models.Conversation
	var total int64
	query := dbengine.Instance().DB.Where("app_id = ?", site.AppID)
	if userID != "" {
		query = query.Where("from_end_user_id = ?", userID)
	}
	query.Model(&models.Conversation{}).Count(&total)
	query.Order("updated_at DESC").Offset((page - 1) * limit).Limit(limit + 1).Find(&conversations)

	hasMore := len(conversations) > limit
	if hasMore {
		conversations = conversations[:limit]
	}

	c.JSON(http.StatusOK, gin.H{
		"data":     conversations,
		"has_more": hasMore,
		"limit":    limit,
	})
}

// DeleteConversation deletes a conversation.
func (a *WebApi) DeleteConversation(c *gin.Context) {
	conversationID := c.Param("conversation_id")
	userID := c.GetString("end_user_id")

	result := dbengine.Instance().DB.Where("id = ? AND from_end_user_id = ?", conversationID, userID).Delete(&models.Conversation{})
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "conversation not found"})
		return
	}
	c.Status(http.StatusNoContent)
}

// RenameConversation renames a conversation.
func (a *WebApi) RenameConversation(c *gin.Context) {
	conversationID := c.Param("conversation_id")
	var body struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := dbengine.Instance().DB.Model(&models.Conversation{}).Where("id = ?", conversationID).Update("name", body.Name)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "conversation not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": "success"})
}

// GetMessages returns messages in a conversation.
func (a *WebApi) GetMessages(c *gin.Context) {
	conversationID := c.Query("conversation_id")
	if conversationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "conversation_id is required"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if limit > 100 {
		limit = 100
	}

	var messages []*models.Message
	var total int64
	query := dbengine.Instance().DB.Where("conversation_id = ?", conversationID)
	query.Model(&models.Message{}).Count(&total)
	query.Order("created_at DESC").Offset((page - 1) * limit).Limit(limit + 1).Find(&messages)

	hasMore := len(messages) > limit
	if hasMore {
		messages = messages[:limit]
	}

	c.JSON(http.StatusOK, gin.H{
		"data":     messages,
		"has_more": hasMore,
		"limit":    limit,
	})
}

// MessageFeedback submits feedback for a message.
func (a *WebApi) MessageFeedback(c *gin.Context) {
	messageID := c.Param("message_id")
	var body struct {
		Rating string `json:"rating"` // like, dislike, null
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString("end_user_id")

	// Upsert feedback
	var existing models.MessageFeedback
	err := dbengine.Instance().DB.Where("message_id = ? AND from_end_user_id = ?", messageID, userID).First(&existing).Error
	if err != nil {
		// Create new
		feedback := models.MessageFeedback{
			ID:            uuid.NewV4().String(),
			MessageID:     messageID,
			Rating:        body.Rating,
			FromEndUserID: userID,
			FromSource:    "api",
		}
		dbengine.Instance().DB.Create(&feedback)
	} else {
		// Update existing
		dbengine.Instance().DB.Model(&existing).Update("rating", body.Rating)
	}

	c.JSON(http.StatusOK, gin.H{"result": "success"})
}

// GetSuggestedQuestions returns suggested follow-up questions.
func (a *WebApi) GetSuggestedQuestions(c *gin.Context) {
	messageID := c.Param("message_id")

	// Load message to get app context
	var message models.Message
	if err := dbengine.Instance().DB.Where("id = ?", messageID).First(&message).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "message not found"})
		return
	}

	// TODO: Generate suggested questions using LLM
	// For now return empty list
	c.JSON(http.StatusOK, gin.H{"data": []string{}})
}

// UploadFile handles file uploads.
func (a *WebApi) UploadFile(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}
	defer file.Close()

	// Read file content
	data, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file"})
		return
	}

	// Save to storage
	fileKey := fmt.Sprintf("uploads/%s/%s", time.Now().Format("20060102"), header.Filename)
	if err := storage.Save(fileKey, data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}

	// Create upload file record
	uploadFile := &models.UploadFile{
		ID:            uuid.NewV4().String(),
		TenantID:      c.GetString("tenant_id"),
		Name:          header.Filename,
		Size:          int(header.Size),
		Key:           fileKey,
		Extension:     filepath.Ext(header.Filename),
		MimeType:      header.Header.Get("Content-Type"),
		CreatedBy:     c.GetString("end_user_id"),
		CreatedByRole: "end_user",
	}
	dbengine.Instance().DB.Create(uploadFile)

	c.JSON(http.StatusCreated, gin.H{
		"id":         uploadFile.ID,
		"name":       uploadFile.Name,
		"size":       uploadFile.Size,
		"mime_type":  uploadFile.MimeType,
		"created_at": uploadFile.CreatedAt,
	})
}

// GetAppParameters returns app configuration parameters for the web client.
func (a *WebApi) GetAppParameters(c *gin.Context) {
	appCode := c.GetHeader("X-App-Code")
	var site models.Site
	if err := dbengine.Instance().DB.Where("code = ?", appCode).First(&site).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "app not found"})
		return
	}

	var app models.App
	if err := dbengine.Instance().DB.Where("id = ?", site.AppID).First(&app).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "app not found"})
		return
	}

	// Parse app model config for user-facing parameters
	var config map[string]any
	if app.AppModelConfigID != "" {
		var modelConfig models.AppModelConfig
		if err := dbengine.Instance().DB.Where("id = ?", app.AppModelConfigID).First(&modelConfig).Error; err == nil {
			json.Unmarshal([]byte(modelConfig.Configs), &config)
		}
	}

	params := map[string]any{
		"opening_statement":    "",
		"suggested_questions":  []string{},
		"speech_to_text":       map[string]any{"enabled": false},
		"text_to_speech":       map[string]any{"enabled": false},
		"retriever_resource":   map[string]any{"enabled": false},
		"annotation_reply":     map[string]any{"enabled": false},
		"file_upload":          map[string]any{"enabled": false},
		"system_parameters":    map[string]any{},
	}

	if config != nil {
		if os, ok := config["opening_statement"].(string); ok {
			params["opening_statement"] = os
		}
		if sq, ok := config["suggested_questions"].([]any); ok {
			params["suggested_questions"] = sq
		}
	}

	c.JSON(http.StatusOK, params)
}

// GetAppMeta returns app metadata.
func (a *WebApi) GetAppMeta(c *gin.Context) {
	appCode := c.GetHeader("X-App-Code")
	var site models.Site
	if err := dbengine.Instance().DB.Where("code = ?", appCode).First(&site).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"tool_icons": map[string]any{}})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tool_icons": map[string]any{},
	})
}

// GetWorkflowRunDetail returns details of a workflow run.
func (a *WebApi) GetWorkflowRunDetail(c *gin.Context) {
	// TODO: Implement
	c.JSON(http.StatusOK, gin.H{})
}
