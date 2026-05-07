package memory

import (
	"fmt"
	"slices"
	"strings"

	modelmanager "github.com/odysseythink/gofy/backend/core/manageres/model_manager"
	promptutils "github.com/odysseythink/gofy/backend/core/prompt/utils"
	dbengine "github.com/odysseythink/gofy/backend/db_engine"
	modelruntimeentities "github.com/odysseythink/gofy/backend/entities/model_runtime"
	"github.com/odysseythink/gofy/backend/models"
	"github.com/odysseythink/mlog"
)

type TokenBufferMemory struct {
	Conversation  *models.Conversation
	ModelInstance *modelmanager.ModelInstance
}

func NewTokenBufferMemory(conversation *models.Conversation, modelInstance *modelmanager.ModelInstance) *TokenBufferMemory {
	return &TokenBufferMemory{
		Conversation:  conversation,
		ModelInstance: modelInstance,
	}
}

func (t *TokenBufferMemory) GetHistoryPromptMessages(max_token_limit int, message_limit int) []modelruntimeentities.PromptMessager {
	if max_token_limit <= 0 {
		max_token_limit = 2000
	}
	// app_record := t.Conversation.App
	var messages []*models.Message
	db := dbengine.Instance().DB.Debug().Model(&models.Message{}).Where("conversation_id = ?", t.Conversation.ID).Order("created_at DESC")
	if message_limit > 0 {
		message_limit = min(message_limit, 500)
	} else {
		message_limit = 500
	}
	err := db.Limit(message_limit).Find(&messages).Error
	if err != nil {
		mlog.Warningf("find Message failed:%v", err)
	}

	// instead of all messages from the conversation, we only need to extract messages
	// that belong to the thread of last message
	thread_messages := promptutils.ExtractThreadMessages(messages)

	// for newly created message, its answer is temporarily empty, we don't need to add it to memory
	if len(thread_messages) > 0 && thread_messages[0].Answer == "" {
		thread_messages = thread_messages[1:]
	}
	slices.Reverse(thread_messages)
	messages = thread_messages

	prompt_messages := make([]modelruntimeentities.PromptMessager, 0)
	for _, message := range messages {
		var files []*models.MessageFile
		err := dbengine.Instance().DB.Debug().Model(&models.MessageFile{}).Where("message_id = ?", message.ID).Find(&files).Error
		if err != nil {
			mlog.Warningf("get MessageFile by message_id(%s) failed:%v", message.ID, err)
		}
		if len(files) > 0 {
			// var file_extra_config *file.FileUploadConfig
			if !slices.Contains([]models.AppMode{models.AppMode_ADVANCED_CHAT, models.AppMode_WORKFLOW}, models.AppMode(t.Conversation.Mode)) {
				// file_extra_config = (&fileupload.FileUploadConfigManager{}).Convert(t.Conversation.ModelConfig, true)
			} else {
				if message.WorkflowRunID != "" {
					workflow_run := new(models.WorkflowRun)
					err = dbengine.Instance().DB.Debug().Model(&models.WorkflowRun{}).Where("id = ?", message.WorkflowRunID).Preload("Workflow").First(workflow_run).Error
					if err != nil {
						mlog.Warningf("get WorkflowRun by id(%s) failed:%v", message.WorkflowRunID, err)
						workflow_run = nil
					}

					if workflow_run != nil && workflow_run.Workflow != nil {
						// file_extra_config = (&fileupload.FileUploadConfigManager{}).Convert(workflow_run.Workflow.FeaturesDict(), false)
					}
				}
			}
			// detail := modelruntimeentities.ImagePromptMessageContentDETAIL_LOW
			// if file_extra_config != nil && app_record != nil{
			//     file_objs = filefactory.BuildFromMessageFiles(
			//         message_files=files, tenant_id=app_record.tenant_id, config=file_extra_config
			//     )
			//     if file_extra_config.image_config and file_extra_config.image_config.detail{
			//         detail = file_extra_config.image_config.detail
			// 	}
			// }else{
			//     file_objs = []
			// }
			// if not file_objs{
			prompt_messages = append(prompt_messages, modelruntimeentities.NewUserPromptMessage(message.Query, ""))
			// }else{
			//     prompt_message_contents: list[PromptMessageContent] = []
			//     prompt_message_contents.append(TextPromptMessageContent(data=message.query))
			//     for file := range file_objs{
			//         prompt_message = file_manager.to_prompt_message_content(
			//             file,
			//             image_detail_config=detail,
			//         )
			//         prompt_message_contents.append(prompt_message)

			//     prompt_messages.append(UserPromptMessage(content=prompt_message_contents))
			// }
		} else {
			prompt_messages = append(prompt_messages, modelruntimeentities.NewUserPromptMessage(message.Query, ""))
		}
		prompt_messages = append(prompt_messages, modelruntimeentities.NewAssistantPromptMessage(message.Answer, "", nil))
	}

	if len(prompt_messages) == 0 {
		return nil
	}

	// prune the chat message if it exceeds the max token limit
	curr_message_tokens := t.ModelInstance.GetLLMNumTokens(prompt_messages, nil)
	if err != nil {
		mlog.Warningf("GetLLMNumTokens failed:%v", err)
	}

	if curr_message_tokens > max_token_limit {
		pruned_memory := []modelruntimeentities.PromptMessager{}
		for curr_message_tokens > max_token_limit && len(prompt_messages) > 1 {
			pruned_memory = append(pruned_memory, prompt_messages[0])
			prompt_messages = prompt_messages[1:]
			curr_message_tokens = t.ModelInstance.GetLLMNumTokens(prompt_messages, nil)
			if err != nil {
				mlog.Warningf("GetLLMNumTokens failed:%v", err)
			}
		}
	}
	return prompt_messages
}

func (t *TokenBufferMemory) GetHistoryPromptText(
	human_prefix string, /* "Human"*/
	ai_prefix string, /* "Assistant"*/
	max_token_limit int, /* 2000*/
	message_limit int,
) string {
	/*
		Get history prompt text.
		:param human_prefix: human prefix
		:param ai_prefix: ai prefix
		:param max_token_limit: max token limit
		:param message_limit: message limit
		:return:
	*/
	if human_prefix == "" {
		human_prefix = "Human"
	}
	if ai_prefix == "" {
		ai_prefix = "Assistant"
	}
	if max_token_limit <= 0 {
		max_token_limit = 2000
	}
	prompt_messages := t.GetHistoryPromptMessages(max_token_limit, message_limit)

	string_messages := []string{}
	for _, m := range prompt_messages {
		role := ""
		if m.Role() == modelruntimeentities.PromptMessageRole_USER {
			role = human_prefix
		} else if m.Role() == modelruntimeentities.PromptMessageRole_ASSISTANT {
			role = human_prefix
		} else {
			continue
		}

		switch real_content := any(m.GetContent).(type) {
		case []modelruntimeentities.PromptMessageContenter:
			inner_msg := ""
			for _, content := range real_content {
				if sub_real_content, ok := any(content).(*modelruntimeentities.TextPromptMessageContent); ok {
					inner_msg += fmt.Sprintf("%s\n", sub_real_content.Data)
				} else if _, ok := any(content).(*modelruntimeentities.ImagePromptMessageContent); ok {
					inner_msg += "[image]\n"
				}
			}
			string_messages = append(string_messages, fmt.Sprintf("%s: %s", role, strings.TrimSpace(inner_msg)))
		case string:
			message := fmt.Sprintf("%s: %s", role, real_content)
			string_messages = append(string_messages, message)
		}
	}
	return strings.Join(string_messages, "\n")
}
