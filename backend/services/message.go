package services

import (
	achatcfgmanage "github.com/odysseythink/gofy/backend/core/app/config_manageres/advanced_chat"
	"github.com/odysseythink/gofy/backend/core/exceptions"
	httpexceptions "github.com/odysseythink/gofy/backend/core/exceptions/http"
	llmgenerator "github.com/odysseythink/gofy/backend/core/llm_generator"
	"github.com/odysseythink/gofy/backend/core/manageres"
	modelmanager "github.com/odysseythink/gofy/backend/core/manageres/model_manager"
	"github.com/odysseythink/gofy/backend/core/memory"
	dbengine "github.com/odysseythink/gofy/backend/db_engine"
	appenumtypes "github.com/odysseythink/gofy/backend/enum_types/app"
	modelruntimeenumtypes "github.com/odysseythink/gofy/backend/enum_types/model_runtime"
	"github.com/odysseythink/gofy/backend/models"
	"github.com/odysseythink/mlog"
)

type MessageService[T interface {
	*models.Account | *models.EndUser
}] struct {
}

func (service *MessageService[T]) GetMessage(app_model *models.App, user T, message_id string) *models.Message {
	message := new(models.Message)
	db := dbengine.Instance().DB.Model(&models.Message{}).Where("id =?", message_id).Where("app_id =?", app_model.ID)
	if real_user, ok := any(user).(*models.EndUser); ok {
		db = db.Where("from_source =?", "api").Where("from_end_user_id =?", real_user.ID).Where("from_account_id =?", "")
	} else if real_user, ok := any(user).(*models.Account); ok {
		db = db.Where("from_source =?", "console").Where("from_end_user_id =?", "").Where("from_account_id =?", real_user.ID)
	}
	err := db.First(message).Error
	if err != nil {
		mlog.Errorf("get message failed:%v", err)
		message = nil
	}
	if message == nil {
		panic(exceptions.NewMessageNotExistsError(""))
	}
	return message
}

func (service *MessageService[T]) GetSuggestedQuestionsAfterAnswer(
	app_model *models.App, user T, message_id string, invoke_from appenumtypes.InvokeFrom,
) []string {
	if user == nil {
		panic(exceptions.NewValueError("user cannot be None"))
	}
	message := service.GetMessage(app_model, user, message_id)
	var conversation *models.Conversation
	switch real_user := any(user).(type) {
	case *models.EndUser:
		conversation = ServiceGroupApp.EndUserConversation.GetConversation(app_model, message.ConversationID, real_user)
	case *models.Account:
		conversation = ServiceGroupApp.AccountConversation.GetConversation(app_model, message.ConversationID, real_user)
	default:
		mlog.Errorf("invalid user:%#v", user)
		panic(exceptions.NewValueError("user is invalid"))
	}

	if conversation == nil {
		panic(exceptions.NewConversationNotExistsError(""))
	}
	if conversation.Status != "normal" {
		panic(httpexceptions.NewConversationCompletedError(""))
	}
	var model_instance *modelmanager.ModelInstance
	if app_model.Mode == models.AppMode_ADVANCED_CHAT {
		var wf *models.Workflow
		if invoke_from == appenumtypes.InvokeFrom_DEBUGGER {
			wf = ServiceGroupApp.Workflow.GetDraftWorkflow(app_model)
		} else {
			wf = ServiceGroupApp.Workflow.GetPublishedWorkflow(app_model)
		}
		if wf == nil {
			return []string{}
		}
		app_config := (achatcfgmanage.New()).GetAppConfig(app_model, wf)

		if !app_config.AdditionalFeatures.SuggestedQuestionsAfterAnswer {
			panic(exceptions.NewSuggestedQuestionsAfterAnswerDisabledError(""))
		}
		model_instance = manageres.Instance.Model.GetDefaultModelInstance(app_model.TenantID, modelruntimeenumtypes.Model_LLM)
	} else {
		override_model_configs := conversation.OverrideModelConfigs()
		var app_model_config *models.AppModelConfig
		if len(override_model_configs) == 0 {
			app_model_config = new(models.AppModelConfig)
			err := dbengine.Instance().DB.Model(&models.AppModelConfig{}).Where("id =? and app_id = ?", conversation.AppModelConfigID, app_model.ID).First(app_model_config).Error
			if err != nil {
				mlog.Errorf("get AppModelConfig failed:%v", err)
				app_model_config = nil
			}
		} else {
			app_model_config = &models.AppModelConfig{
				ID:    conversation.AppModelConfigID,
				AppID: app_model.ID,
			}

			app_model_config.FromModelConfigDict(override_model_configs)
		}
		if app_model_config == nil {
			panic(exceptions.NewValueError("did not find app model config"))
		}
		suggested_questions_after_answer := app_model_config.SuggestedQuestionsAfterAnswerDict()
		suggested_questions_after_answer_enabled := false
		if _, ok := suggested_questions_after_answer["enabled"]; ok {
			if _, ok := suggested_questions_after_answer["enabled"].(bool); ok {
				suggested_questions_after_answer_enabled = suggested_questions_after_answer["enabled"].(bool)
			}
		}
		if !suggested_questions_after_answer_enabled {
			panic(exceptions.NewSuggestedQuestionsAfterAnswerDisabledError(""))
		}
		provider := ""
		model_name := ""
		model_dict := app_model_config.ModelDict()
		if _, ok := model_dict["provider"]; ok {
			if _, ok := model_dict["provider"].(string); ok {
				provider = model_dict["provider"].(string)
			}
		}
		if _, ok := model_dict["name"]; ok {
			if _, ok := model_dict["name"].(string); ok {
				model_name = model_dict["name"].(string)
			}
		}
		model_instance = manageres.Instance.Model.GetModelInstance(
			app_model.TenantID,
			provider,
			modelruntimeenumtypes.Model_LLM,
			model_name,
		)
	}
	// get memory of conversation (read-only)
	mem := memory.NewTokenBufferMemory(conversation, model_instance)

	histories := mem.GetHistoryPromptText("", "", 3000, 3)

	questions := (&llmgenerator.LLMGenerator{}).GenerateSuggestedQuestionsAfterAnswer(app_model.TenantID, histories)

	return questions
}
