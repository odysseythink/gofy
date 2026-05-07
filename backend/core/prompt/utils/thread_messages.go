package utils

import (
	"slices"

	"github.com/odysseythink/gofy/backend/constants"
	dbengine "github.com/odysseythink/gofy/backend/db_engine"
	"github.com/odysseythink/gofy/backend/models"
	"github.com/odysseythink/mlog"
)

func ExtractThreadMessages(messages []*models.Message) []*models.Message {
	thread_messages := []*models.Message{}
	var next_message string

	for _, message := range messages {
		if message.ParentMessageID == "" {
			// If the message is regenerated and does not have a parent message, it is the start of a new thread
			thread_messages = append(thread_messages, message)
			break
		}
		if next_message == "" {
			thread_messages = append(thread_messages, message)
			next_message = message.ParentMessageID
		} else {
			if slices.Contains([]string{message.ID, constants.UUID_NIL}, next_message) {
				thread_messages = append(thread_messages, message)
				next_message = message.ParentMessageID
			}
		}
	}
	return thread_messages
}

func GetThreadMessagesLength(conversation_id string) int {
	/*
	   Get the number of thread messages based on the parent message id.
	*/
	// Fetch all messages related to the conversation
	var messages []*models.Message
	err := dbengine.Instance().DB.Model(&models.Message{}).Where("conversation_id = ?", conversation_id).Order("created_at DESC").Find(&messages).Error
	if err != nil {
		mlog.Errorf("get Message failed:%v", err)
	}

	// Extract thread messages
	thread_messages := ExtractThreadMessages(messages)

	// Exclude the newly created message with an empty answer
	if len(thread_messages) > 0 && thread_messages[0].Answer == "" {
		thread_messages = thread_messages[1:]
	}
	return len(thread_messages)
}
