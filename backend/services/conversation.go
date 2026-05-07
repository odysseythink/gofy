package services

import (
	"fmt"
	"time"

	"github.com/odysseythink/gofy/backend/core/exceptions"
	httpexceptions "github.com/odysseythink/gofy/backend/core/exceptions/http"
	dbengine "github.com/odysseythink/gofy/backend/db_engine"
	"github.com/odysseythink/gofy/backend/models"
	"github.com/odysseythink/mlog"
)

type ConversationService[T interface {
	*models.Account | *models.EndUser
}] struct {
}

func (service *ConversationService[T]) GetConversation(app_model *models.App, conversation_id string, user T) *models.Conversation {
	conversation := new(models.Conversation)
	db := dbengine.Instance().DB.Model(&models.Conversation{}).Where("id =? and app_id =?", conversation_id, app_model.ID)
	if real_user, ok := any(user).(*models.EndUser); ok {
		db = db.Where("from_source =?", "api").Where("from_end_user_id =?", real_user.ID).Where("from_account_id =?", "")
	} else if real_user, ok := any(user).(*models.Account); ok {
		db = db.Where("from_source =?", "console").Where("from_end_user_id =?", "").Where("from_account_id =?", real_user.ID)
	}
	db = db.Where("is_deleted = ?", false)
	err := db.First(conversation).Error
	if err != nil {
		mlog.Errorf("get conversation failed:%v", err)
		conversation = nil
	}
	if conversation == nil {
		panic(exceptions.NewConversationNotExistsError(""))
	}

	return conversation
}

func (service *ConversationService[T]) GetConversation1(current_user T, app_model *models.App, conversation_id string) *models.Conversation {
	if _, ok := any(current_user).(*models.Account); !ok || any(current_user).(*models.Account) == nil {
		mlog.Errorf("current_user must be account")
		panic(exceptions.NewValueError("current_user must be account."))
	}
	real_user := any(current_user).(*models.Account)
	conversation := new(models.Conversation)
	err := dbengine.Instance().DB.Model(&models.Conversation{}).Where("id = ? and app_id=?", conversation_id, app_model.ID).First(conversation).Error
	if err != nil {
		mlog.Errorf("get Conversation from mysql failed:%v", err)
		conversation = nil
	}

	if conversation == nil {
		panic(httpexceptions.NewNotFound("Conversation Not Exists."))
	}
	if conversation.ReadAt == nil {
		now := time.Now()
		conversation.ReadAt = &now
		conversation.ReadAccountID = real_user.ID
		dbengine.Instance().DB.Updates(&models.Conversation{ID: conversation.ID, ReadAt: conversation.ReadAt, ReadAccountID: conversation.ReadAccountID})
	}
	return conversation
}

func (service *ConversationService[T]) AutoGenerateName(conversationID string) (string, error) {
	var conv models.Conversation
	if err := dbengine.Instance().DB.Where("id = ?", conversationID).First(&conv).Error; err != nil {
		return "", fmt.Errorf("conversation not found")
	}
	// TODO: Use LLM to generate a name from conversation content
	// For now, use first message content as name
	var msg models.Message
	if err := dbengine.Instance().DB.Where("conversation_id = ?", conversationID).Order("created_at ASC").First(&msg).Error; err == nil {
		name := msg.Query
		if len([]rune(name)) > 50 {
			name = string([]rune(name)[:50]) + "..."
		}
		if name != "" {
			dbengine.Instance().DB.Model(&conv).Update("name", name)
			return name, nil
		}
	}
	return "New Conversation", nil
}

func (service *ConversationService[T]) DeleteConversation(conversationID, userID string) error {
	return dbengine.Instance().DB.Where("id = ? AND (from_account_id = ? OR from_end_user_id = ?)", conversationID, userID, userID).Delete(&models.Conversation{}).Error
}

func (service *ConversationService[T]) RenameConversation(conversationID, name string) error {
	return dbengine.Instance().DB.Model(&models.Conversation{}).Where("id = ?", conversationID).Update("name", name).Error
}

func (service *ConversationService[T]) GetConversationsByAppID(appID string, page, limit int, sortBy string) ([]*models.Conversation, int64) {
	var conversations []*models.Conversation
	var total int64
	query := dbengine.Instance().DB.Where("app_id = ?", appID)
	query.Model(&models.Conversation{}).Count(&total)
	if sortBy == "" {
		sortBy = "updated_at DESC"
	}
	query.Order(sortBy).Offset((page - 1) * limit).Limit(limit).Find(&conversations)
	return conversations, total
}

func (service *ConversationService[T]) GetConversationDetail(conversationID string) *models.Conversation {
	var conv models.Conversation
	if err := dbengine.Instance().DB.Where("id = ?", conversationID).First(&conv).Error; err != nil {
		return nil
	}
	return &conv
}

func (service *ConversationService[T]) PaginateByLastID(appID string, lastID string, limit int, userID string) ([]*models.Conversation, bool) {
	var conversations []*models.Conversation
	query := dbengine.Instance().DB.Where("app_id = ?", appID)
	if userID != "" {
		query = query.Where("from_end_user_id = ? OR from_account_id = ?", userID, userID)
	}
	if lastID != "" {
		var lastConv models.Conversation
		if err := dbengine.Instance().DB.Where("id = ?", lastID).First(&lastConv).Error; err == nil {
			query = query.Where("created_at < ?", lastConv.CreatedAt)
		}
	}
	query.Order("created_at DESC").Limit(limit + 1).Find(&conversations)
	hasMore := len(conversations) > limit
	if hasMore {
		conversations = conversations[:limit]
	}
	return conversations, hasMore
}
