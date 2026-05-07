package chat

import (
	"github.com/odysseythink/gofy/backend/core/app/generatores/base"
)

type ChatAppGenerator struct {
	*base.BaseAppGenerator
}

func NewChatAppGenerator() *ChatAppGenerator {
	return &ChatAppGenerator{
		BaseAppGenerator: &base.BaseAppGenerator{},
	}
}
