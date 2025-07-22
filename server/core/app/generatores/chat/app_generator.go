package chat

import (
	"mlib.com/gofy/server/core/app/generatores/base"
)

type ChatAppGenerator struct {
	*base.BaseAppGenerator
}

func NewChatAppGenerator() *ChatAppGenerator {
	return &ChatAppGenerator{
		BaseAppGenerator: &base.BaseAppGenerator{},
	}
}
