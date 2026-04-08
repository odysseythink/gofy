package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type WebApi struct{}

// GetConversations returns conversations for an app user.
func (a *WebApi) GetConversations(c *gin.Context) {
	// TODO: Implement - query conversations by app_id and end_user
	c.JSON(http.StatusOK, gin.H{"data": []any{}, "has_more": false})
}

// DeleteConversation deletes a conversation.
func (a *WebApi) DeleteConversation(c *gin.Context) {
	// TODO: Implement
	c.JSON(http.StatusNoContent, nil)
}

// RenameConversation renames a conversation.
func (a *WebApi) RenameConversation(c *gin.Context) {
	// TODO: Implement
	c.JSON(http.StatusOK, gin.H{"result": "success"})
}

// GetMessages returns messages in a conversation.
func (a *WebApi) GetMessages(c *gin.Context) {
	// TODO: Implement
	c.JSON(http.StatusOK, gin.H{"data": []any{}, "has_more": false})
}

// MessageFeedback submits feedback for a message.
func (a *WebApi) MessageFeedback(c *gin.Context) {
	// TODO: Implement
	c.JSON(http.StatusOK, gin.H{"result": "success"})
}

// GetSuggestedQuestions returns suggested follow-up questions.
func (a *WebApi) GetSuggestedQuestions(c *gin.Context) {
	// TODO: Implement
	c.JSON(http.StatusOK, gin.H{"data": []string{}})
}

// UploadFile handles file uploads.
func (a *WebApi) UploadFile(c *gin.Context) {
	// TODO: Implement
	c.JSON(http.StatusCreated, gin.H{"id": "", "name": ""})
}

// GetAppParameters returns app configuration parameters for the web client.
func (a *WebApi) GetAppParameters(c *gin.Context) {
	// TODO: Implement - return app config visible to end users
	c.JSON(http.StatusOK, gin.H{
		"opening_statement":   "",
		"suggested_questions": []string{},
		"speech_to_text":      map[string]any{"enabled": false},
		"text_to_speech":      map[string]any{"enabled": false},
		"retriever_resource":  map[string]any{"enabled": false},
		"annotation_reply":    map[string]any{"enabled": false},
		"file_upload":         map[string]any{"enabled": false},
		"system_parameters":   map[string]any{},
	})
}

// GetAppMeta returns app metadata.
func (a *WebApi) GetAppMeta(c *gin.Context) {
	// TODO: Implement
	c.JSON(http.StatusOK, gin.H{
		"tool_icons": map[string]any{},
	})
}

// GetWorkflowRunDetail returns details of a workflow run.
func (a *WebApi) GetWorkflowRunDetail(c *gin.Context) {
	// TODO: Implement
	c.JSON(http.StatusOK, gin.H{})
}
