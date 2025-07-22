package message

import (
	"mlib.com/gofy/server/core/app/queue_manager/base"
	"mlib.com/gofy/server/core/exceptions"
	appqueueentities "mlib.com/gofy/server/entities/app/queue"
	appenumtypes "mlib.com/gofy/server/enum_types/app"
)

type MessageAppQueueManager struct {
	*base.BaseAppQueueManager[*appqueueentities.MessageQueueMessage]
	appMode        string
	conversationID string
	messageID      string
}

func New(task_id string, user_id string, invoke_from appenumtypes.InvokeFrom, conversation_id string, app_mode string, message_id string) *MessageAppQueueManager {
	mgr := &MessageAppQueueManager{
		BaseAppQueueManager: base.New[*appqueueentities.MessageQueueMessage](task_id, user_id, invoke_from),
		appMode:             app_mode,
		conversationID:      conversation_id,
		messageID:           message_id,
	}
	return mgr
}

func (mgr *MessageAppQueueManager) ConstructQueueMessage(event appqueueentities.AppQueueEventer) *appqueueentities.MessageQueueMessage {
	return &appqueueentities.MessageQueueMessage{
		QueueMessage: &appqueueentities.QueueMessage{
			TaskID:  mgr.TaskID,
			AppMode: mgr.appMode,
			Eventer: event,
		},
		MessageID:      mgr.messageID,
		ConversationID: mgr.conversationID,
	}
}
func (mgr *MessageAppQueueManager) Publish(event appqueueentities.AppQueueEventer, pub_from appenumtypes.PublishFrom) {
	/*
	   Publish event to queue
	   :param event:
	   :param pub_from:
	   :return:
	*/
	message := &appqueueentities.MessageQueueMessage{
		QueueMessage: &appqueueentities.QueueMessage{
			TaskID:  mgr.TaskID,
			AppMode: mgr.appMode,
			Eventer: event,
		},
		MessageID:      mgr.messageID,
		ConversationID: mgr.conversationID,
	}
	mgr.MsgQueue <- message

	switch any(event).(type) {
	case *appqueueentities.QueueStopEvent, *appqueueentities.QueueErrorEvent, *appqueueentities.QueueMessageEndEvent, *appqueueentities.QueueAdvancedChatMessageEndEvent:
		mgr.StopListen()
	}

	if pub_from == appenumtypes.PublishFrom_APPLICATION_MANAGER && mgr.IsStopped() {
		panic(&exceptions.GenerateTaskStoppedError{})
	}
}
