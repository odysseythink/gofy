package events

import (
	"log"

	"mlib.com/gofy/server/models"
	"mlib.com/mrun"
)

func init() {
	var err error
	Instance.MessageWasCreatedSig, err = mrun.NewSignal("message_was_created", func(*models.Message, any) {})
	if err != nil {
		log.Printf("[E]create signal failed:%v", err)
	}
}
