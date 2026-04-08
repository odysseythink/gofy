package queue

import (
	"iter"

	appenumtypes "mlib.com/gofy/server/enum_types/app"
)

type AppQueueManager[T interface {
	*MessageQueueMessage | *WorkflowQueueMessage
}] interface {
	Publish(event AppQueueEventer, pub_from appenumtypes.PublishFrom)
	Listen(aqm AppQueueManager[T]) iter.Seq[T]
	StopListen()
	PublishError(AppQueueManager[T], error, appenumtypes.PublishFrom)
}
