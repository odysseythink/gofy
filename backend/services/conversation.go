package services

import (
	"time"

	"mlib.com/gofy/server/core/exceptions"
	httpexceptions "mlib.com/gofy/server/core/exceptions/http"
	dbengine "mlib.com/gofy/server/db_engine"
	"mlib.com/gofy/server/models"
	"mlib.com/mlog"
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
