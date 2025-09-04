package functioncall

import "mlib.com/gofy/server/core/app/runner/base"

type FunctionCallAgentRunner struct {
	*base.AppRunner[*appqueueentities.MessageQueueMessage]
}
