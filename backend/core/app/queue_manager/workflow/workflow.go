package workflow

import (
	"github.com/odysseythink/mlog"
	"mlib.com/gofy/server/core/app/queue_manager/base"
	"mlib.com/gofy/server/core/exceptions"
	appqueueentities "mlib.com/gofy/server/entities/app/queue"
	appenumtypes "mlib.com/gofy/server/enum_types/app"
)

type WorkflowAppQueueManager struct {
	*base.BaseAppQueueManager[*appqueueentities.WorkflowQueueMessage]
	appMode string
}

func New(task_id string, user_id string, invoke_from appenumtypes.InvokeFrom, app_mode string) *WorkflowAppQueueManager {
	return &WorkflowAppQueueManager{
		BaseAppQueueManager: base.New[*appqueueentities.WorkflowQueueMessage](task_id, user_id, invoke_from),
		appMode:             app_mode,
	}
}
func (mgr *WorkflowAppQueueManager) Publish(event appqueueentities.AppQueueEventer, pub_from appenumtypes.PublishFrom) {
	/*
	   Publish event to queue
	   :param event:
	   :param pub_from:
	   :return:
	*/
	mlog.Debugf("------publish event=%#v", event)
	message := &appqueueentities.WorkflowQueueMessage{
		QueueMessage: &appqueueentities.QueueMessage{
			TaskID:  mgr.TaskID,
			AppMode: mgr.appMode,
			Eventer: event,
		},
	}
	mgr.MsgQueue <- message

	if event.Event() == appenumtypes.QueueEvent_STOP ||
		event.Event() == appenumtypes.QueueEvent_ERROR ||
		event.Event() == appenumtypes.QueueEvent_WORKFLOW_SUCCEEDED ||
		event.Event() == appenumtypes.QueueEvent_WORKFLOW_FAILED ||
		event.Event() == appenumtypes.QueueEvent_WORKFLOW_PARTIAL_SUCCEEDED {
		mgr.StopListen()
	}

	if pub_from == appenumtypes.PublishFrom_APPLICATION_MANAGER && mgr.IsStopped() {
		panic(&exceptions.GenerateTaskStoppedError{})
	}
}
