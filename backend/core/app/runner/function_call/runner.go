package functioncall

import "github.com/odysseythink/gofy/backend/core/app/runner/base"

type FunctionCallAgentRunner struct {
	*base.AppRunner[*appqueueentities.MessageQueueMessage]
}
