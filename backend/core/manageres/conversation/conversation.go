package conversation

import (
	"github.com/odysseythink/mlog"
	"mlib.com/gofy/server/core/exceptions"
	httpexceptions "mlib.com/gofy/server/core/exceptions/http"
	dbengine "mlib.com/gofy/server/db_engine"
	"mlib.com/gofy/server/models"
)

type ConversationManager struct {
}

func New() *ConversationManager {
	return &ConversationManager{}
}

func GetConversationByUser[T *models.Account | *models.EndUser](
	app_model *models.App, conversation_id string, user T,
) *models.Conversation {
	db := dbengine.Instance().DB.Model(&models.Conversation{}).Where("id = ? and app_id = ? and status = ? and is_deleted = ?", conversation_id, app_model.ID, "normal", false)

	if real_user, ok := any(user).(*models.Account); ok {
		db.Where("from_account_id = ?", real_user.ID)
	} else if real_user, ok := any(user).(*models.EndUser); ok && real_user != nil {
		db.Where("from_end_user_id = ?", real_user.ID)
	}
	conversation := new(models.Conversation)
	err := db.First(conversation).Error
	if err != nil {
		mlog.Errorf("get conversation failed:%v", err)
		conversation = nil
	}

	if conversation == nil {
		panic(exceptions.NewConversationNotExistsError(""))
	}
	if conversation.Status != "normal" {
		panic(httpexceptions.NewConversationCompletedError(""))
	}
	return conversation
}

func GetAppModelConfig(app_model *models.App, conversation *models.Conversation) *models.AppModelConfig {
	var app_model_config *models.AppModelConfig
	if conversation != nil {
		app_model_config = new(models.AppModelConfig)
		err := dbengine.Instance().DB.Model(&models.AppModelConfig{}).Where("id = ? and app_id = ?", conversation.AppModelConfigID, app_model.ID).First(app_model_config).Error
		if err != nil {
			mlog.Errorf("get AppModelConfig failed:%v", err)
			app_model_config = nil
		}

		if app_model_config == nil {
			panic(exceptions.NewAppModelConfigBrokenError(""))
		}
	} else {
		if app_model.AppModelConfigID == "" {
			panic(exceptions.NewAppModelConfigBrokenError(""))
		}
		app_model_config := app_model.AppModelConfig

		if app_model_config == nil {
			panic(exceptions.NewAppModelConfigBrokenError(""))
		}
	}
	return app_model_config
}
