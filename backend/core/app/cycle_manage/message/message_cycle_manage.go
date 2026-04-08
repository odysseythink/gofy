package message

import (
	llmgenerator "mlib.com/gofy/server/core/llm_generator"
	dbengine "mlib.com/gofy/server/db_engine"
	appgeneratorentities "mlib.com/gofy/server/entities/app/generator"
	appqueueentities "mlib.com/gofy/server/entities/app/queue"
	appresponseentities "mlib.com/gofy/server/entities/app/response"
	apptaskentities "mlib.com/gofy/server/entities/app/task"
	"mlib.com/gofy/server/models"
	"mlib.com/mlog"
)

type MessageCycleManage[T1 *appgeneratorentities.ChatAppGenerateEntity | *appgeneratorentities.CompletionAppGenerateEntity | *appgeneratorentities.AgentChatAppGenerateEntity | *appgeneratorentities.AdvancedChatAppGenerateEntity, T2 *apptaskentities.EasyUITaskState | *apptaskentities.WorkflowTaskState] struct {
	ApplicationGenerateEntity T1
	TaskState                 T2
}

func New[T1 *appgeneratorentities.ChatAppGenerateEntity | *appgeneratorentities.CompletionAppGenerateEntity | *appgeneratorentities.AgentChatAppGenerateEntity | *appgeneratorentities.AdvancedChatAppGenerateEntity, T2 *apptaskentities.EasyUITaskState | *apptaskentities.WorkflowTaskState](
	application_generate_entity T1,
	task_state T2,
) *MessageCycleManage[T1, T2] {
	return &MessageCycleManage[T1, T2]{
		ApplicationGenerateEntity: application_generate_entity,
		TaskState:                 task_state,
	}
}

func (mgr *MessageCycleManage[T1, T2]) GenerateConversationName(conversation_id string, query string) {
	/*
	   Generate conversation name.
	   :param conversation: conversation
	   :param query: query
	   :return: thread
	*/
	is_first_message := false
	auto_generate_conversation_name := true
	if _, ok := any(mgr.ApplicationGenerateEntity).(*appgeneratorentities.CompletionAppGenerateEntity); ok {
		return
	} else if real_entity, ok := any(mgr.ApplicationGenerateEntity).(*appgeneratorentities.ChatAppGenerateEntity); ok {
		is_first_message = real_entity.ConversationID == ""
		if len(real_entity.Extras) > 0 {
			if _, ok := real_entity.Extras["auto_generate_conversation_name"]; ok {
				if _, ok := real_entity.Extras["auto_generate_conversation_name"].(bool); ok {
					auto_generate_conversation_name = real_entity.Extras["auto_generate_conversation_name"].(bool)
				}
			}
		}
	} else if real_entity, ok := any(mgr.ApplicationGenerateEntity).(*appgeneratorentities.AgentChatAppGenerateEntity); ok {
		is_first_message = real_entity.ConversationID == ""

		if len(real_entity.Extras) > 0 {
			if _, ok := real_entity.Extras["auto_generate_conversation_name"]; ok {
				if _, ok := real_entity.Extras["auto_generate_conversation_name"].(bool); ok {
					auto_generate_conversation_name = real_entity.Extras["auto_generate_conversation_name"].(bool)
				}
			}
		}
	} else if real_entity, ok := any(mgr.ApplicationGenerateEntity).(*appgeneratorentities.AdvancedChatAppGenerateEntity); ok {
		is_first_message = real_entity.ConversationID == ""
		if len(real_entity.Extras) > 0 {
			if _, ok := real_entity.Extras["auto_generate_conversation_name"]; ok {
				if _, ok := real_entity.Extras["auto_generate_conversation_name"].(bool); ok {
					auto_generate_conversation_name = real_entity.Extras["auto_generate_conversation_name"].(bool)
				}
			}
		}
	}

	if auto_generate_conversation_name && is_first_message {
		// start generate thread
		go mgr._generate_conversation_name_worker(conversation_id, query)
	}
}
func (mgr *MessageCycleManage[T1, T2]) _generate_conversation_name_worker(conversation_id string, query string) {
	// with flask_app.app_context():
	// get conversation and message
	conversation := new(models.Conversation)
	err := dbengine.Instance().DB.Model(&models.Conversation{}).Where("id = ?", conversation_id).Preload("App").First(conversation).Error
	if err != nil {
		mlog.Errorf("get Conversation(%s) failed:%v", conversation_id, err)
		conversation = nil
	}

	if conversation == nil {
		return
	}
	if conversation.Mode != models.AppMode_COMPLETION {
		app_model := conversation.App
		if app_model == nil {
			return
		}
		// generate conversation name

		name := (&llmgenerator.LLMGenerator{}).GenerateConversationName(app_model.TenantID, query, "", "")
		conversation.Name = name
		dbengine.Instance().DB.Save(conversation)
	}
}

func (mgr *MessageCycleManage[T1, T2]) HandleAnnotationReply(event *appqueueentities.QueueAnnotationReplyEvent) *models.MessageAnnotation {
	/*
	   Handle annotation reply.
	   :param event: event
	   :return:
	*/
	annotation := new(models.MessageAnnotation)
	err := dbengine.Instance().DB.Model(&models.MessageAnnotation{}).Where("id = ?", event.MessageAnnotationID).Preload("Account").First(annotation).Error
	if err != nil {
		mlog.Errorf("get MessageAnnotation(%s) failed:%v", event.MessageAnnotationID, err)
		annotation = nil
	}

	if annotation != nil {
		account := annotation.Account()
		name := "Gofy user"
		if account != nil {
			name = account.Name
		}
		annotation_reply := map[string]any{
			"id":      annotation.ID,
			"account": map[string]any{"id": annotation.AccountID, "name": name},
		}
		if real_task_state, ok := any(mgr.TaskState).(*apptaskentities.EasyUITaskState); ok {
			if real_task_state.Metadata == nil {
				real_task_state.Metadata = make(map[string]any)
			}
			real_task_state.Metadata["annotation_reply"] = annotation_reply
		} else if real_task_state, ok := any(mgr.TaskState).(*apptaskentities.WorkflowTaskState); ok {
			if real_task_state.Metadata == nil {
				real_task_state.Metadata = make(map[string]any)
			}
			real_task_state.Metadata["annotation_reply"] = annotation_reply
		}

		return annotation
	}
	return nil
}

func (mgr *MessageCycleManage[T1, T2]) HandleRetrieverResources(event *appqueueentities.QueueRetrieverResourcesEvent) {
	/*
	   Handle retriever resources.
	   :param event: event
	   :return:
	*/
	show_retrieve_source := false
	if real_entity, ok := any(mgr.ApplicationGenerateEntity).(*appgeneratorentities.CompletionAppGenerateEntity); ok {
		if real_entity.AppConfig != nil {
			if real_entity.AppConfig.AdditionalFeatures != nil {
				show_retrieve_source = real_entity.AppConfig.AdditionalFeatures.ShowRetrieveSource
			}
		}
	} else if real_entity, ok := any(mgr.ApplicationGenerateEntity).(*appgeneratorentities.ChatAppGenerateEntity); ok {
		if real_entity.AppConfig != nil {
			if real_entity.AppConfig.AdditionalFeatures != nil {
				show_retrieve_source = real_entity.AppConfig.AdditionalFeatures.ShowRetrieveSource
			}
		}
	} else if real_entity, ok := any(mgr.ApplicationGenerateEntity).(*appgeneratorentities.AgentChatAppGenerateEntity); ok {
		if real_entity.AppConfig != nil {
			if real_entity.AppConfig.AdditionalFeatures != nil {
				show_retrieve_source = real_entity.AppConfig.AdditionalFeatures.ShowRetrieveSource
			}
		}
	} else if real_entity, ok := any(mgr.ApplicationGenerateEntity).(*appgeneratorentities.AdvancedChatAppGenerateEntity); ok {
		if real_entity.AppConfig != nil {
			if real_entity.AppConfig.AdditionalFeatures != nil {
				show_retrieve_source = real_entity.AppConfig.AdditionalFeatures.ShowRetrieveSource
			}
		}
	}
	if show_retrieve_source {
		if real_task_state, ok := any(mgr.TaskState).(*apptaskentities.EasyUITaskState); ok {
			if real_task_state.Metadata == nil {
				real_task_state.Metadata = make(map[string]any)
			}
			real_task_state.Metadata["retriever_resources"] = event.RetrieverResources
		} else if real_task_state, ok := any(mgr.TaskState).(*apptaskentities.WorkflowTaskState); ok {
			if real_task_state.Metadata == nil {
				real_task_state.Metadata = make(map[string]any)
			}
			real_task_state.Metadata["retriever_resources"] = event.RetrieverResources
		}
	}
}

// func(mgr *MessageCycleManage[T1, T2]) _message_file_to_stream_response(event: QueueMessageFileEvent) -> Optional[MessageFileStreamResponse]{
//         /*
//         Message file to stream response.
//         :param event: event
//         :return:
//         */
//         message_file = db.session.query(MessageFile).filter(MessageFile.id == event.message_file_id).first()

//         if message_file and message_file.url is not None:
//             // get tool file id
//             tool_file_id = message_file.url.split("/")[-1]
//             // trim extension
//             tool_file_id = tool_file_id.split(".")[0]

//             // get extension
//             if "." in message_file.url:
//                 extension = f".{message_file.url.split('.')[-1]}"
//                 if len(extension) > 10:
//                     extension = ".bin"
//             else:
//                 extension = ".bin"
//             // add sign url to local file
//             if message_file.url.startswith("http"):
//                 url = message_file.url
//             else:
//                 url = ToolFileManager.sign_file(tool_file_id=tool_file_id, extension=extension)

//             return MessageFileStreamResponse(
//                 task_id=self._application_generate_entity.task_id,
//                 id=message_file.id,
//                 type=message_file.type,
//                 belongs_to=message_file.belongs_to or "user",
//                 url=url,
//             )

//	        return None
//	}
func (mgr *MessageCycleManage[T1, T2]) MessageToStreamResponse(
	answer string, message_id string, from_variable_selector []string,
) *appresponseentities.MessageStreamResponse {
	/*
	   Message to stream response.
	   :param answer: answer
	   :param message_id: message id
	   :return:
	*/
	task_id := ""
	if real_entity, ok := any(mgr.ApplicationGenerateEntity).(*appgeneratorentities.CompletionAppGenerateEntity); ok {
		task_id = real_entity.TaskID
	} else if real_entity, ok := any(mgr.ApplicationGenerateEntity).(*appgeneratorentities.ChatAppGenerateEntity); ok {
		task_id = real_entity.TaskID
	} else if real_entity, ok := any(mgr.ApplicationGenerateEntity).(*appgeneratorentities.AgentChatAppGenerateEntity); ok {
		task_id = real_entity.TaskID
	} else if real_entity, ok := any(mgr.ApplicationGenerateEntity).(*appgeneratorentities.AdvancedChatAppGenerateEntity); ok {
		task_id = real_entity.TaskID
	}
	return appresponseentities.NewMessageStreamResponse(
		task_id,
		message_id,
		answer,
		from_variable_selector,
	)
}

func (mgr *MessageCycleManage[T1, T2]) MessageReplaceToStreamResponse(answer string) *appresponseentities.MessageReplaceStreamResponse {
	/*
	   Message replace to stream response.
	   :param answer: answer
	   :return:
	*/
	task_id := ""
	if real_entity, ok := any(mgr.ApplicationGenerateEntity).(*appgeneratorentities.CompletionAppGenerateEntity); ok {
		task_id = real_entity.TaskID
	} else if real_entity, ok := any(mgr.ApplicationGenerateEntity).(*appgeneratorentities.ChatAppGenerateEntity); ok {
		task_id = real_entity.TaskID
	} else if real_entity, ok := any(mgr.ApplicationGenerateEntity).(*appgeneratorentities.AgentChatAppGenerateEntity); ok {
		task_id = real_entity.TaskID
	} else if real_entity, ok := any(mgr.ApplicationGenerateEntity).(*appgeneratorentities.AdvancedChatAppGenerateEntity); ok {
		task_id = real_entity.TaskID
	}
	return appresponseentities.NewMessageReplaceStreamResponse(task_id, answer)
}
