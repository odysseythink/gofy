package base

import (
	"iter"
	"slices"
	"time"

	"github.com/spf13/viper"
	"mlib.com/gofy/server/cache"
	"mlib.com/gofy/server/core/exceptions"
	appqueueentities "mlib.com/gofy/server/entities/app/queue"
	appenumtypes "mlib.com/gofy/server/enum_types/app"
	"mlib.com/mlog"
)

type baseAppQueueMessageType interface {
	*appqueueentities.WorkflowQueueMessage | *appqueueentities.MessageQueueMessage
}

type BaseAppQueueManager[T baseAppQueueMessageType] struct {
	TaskID     string
	UserID     string
	InvokeFrom appenumtypes.InvokeFrom
	MsgQueue   chan T
}

func New[T *appqueueentities.WorkflowQueueMessage | *appqueueentities.MessageQueueMessage](task_id, user_id string, invoke_from appenumtypes.InvokeFrom) *BaseAppQueueManager[T] {
	if user_id == "" {
		panic(exceptions.NewValueError("user is required"))
	}
	mgr := &BaseAppQueueManager[T]{
		TaskID:     task_id,
		UserID:     user_id,
		InvokeFrom: invoke_from,
		// WorkflowMsgQueue: make(chan appentities.WorkflowQueueMessage, viper.GetIntWithDefault("app_config.workflow_msg_queue_capacity", 1024)),
		MsgQueue: make(chan T, viper.GetIntWithDefault("app_config.msg_queue_capacity", 1024)),
	}
	var user_prefix string
	if slices.Contains([]appenumtypes.InvokeFrom{appenumtypes.InvokeFrom_EXPLORE, appenumtypes.InvokeFrom_DEBUGGER}, invoke_from) {
		user_prefix = "account"
	} else {
		user_prefix = "end-user"
	}

	cache.Instance().SetEx(mgr.generateTaskBelongCacheKey(task_id), user_prefix+"-"+user_id, 1800*time.Second)

	return mgr
}

func (mgr *BaseAppQueueManager[T]) Listen(aqm appqueueentities.AppQueueManager[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		mlog.Debugf("****************Listening")
		refreshTicker := time.NewTicker(1000 * time.Millisecond)
		start_time := time.Now()
		last_ping_time := 0
		for {
			select {
			case message := <-mgr.MsgQueue:
				mlog.Debugf("------receive message:%#v", message)
				if message == nil {
					return
				}
				if !yield(message) {
					return
				}
			case <-refreshTicker.C:
				mlog.Debugf("------refresh")
				elapsed_time := time.Since(start_time).Seconds()
				if elapsed_time >= float64(viper.GetIntWithDefault("app_config.max_execution_time", 1200)) || mgr.IsStopped() {
					aqm.Publish(
						&appqueueentities.QueueStopEvent{StoppedBy: appenumtypes.QueueStopEvent_StopBy_USER_MANUAL}, appenumtypes.PublishFrom_TASK_PIPELINE,
					)
				}
				if int(elapsed_time)/10 > last_ping_time {
					aqm.Publish(&appqueueentities.QueuePingEvent{}, appenumtypes.PublishFrom_TASK_PIPELINE)
					last_ping_time = int(elapsed_time) / 10
				}
			}
		}
	}
	/*
	   Listen to queue
	   :return:
	*/
	// wait for APP_MAX_EXECUTION_TIME seconds to stop listen
	// listen_timeout := viper.GetInt("app_config.max_execution_time")
	// start_time := time.Now()
	// for {

	// }
	// last_ping_time: int | float = 0
	// for {
	//     try:
	//         message = self._q.get(timeout=1)
	//         if message is None:
	//             break

	//	yield message
	//
	// except queue.Empty:
	//
	//	continue
	//
	// finally:
	//
	//	    elapsed_time = time.time() - start_time
	//	    if elapsed_time >= listen_timeout or self._is_stopped():
	//	        // publish two messages to make sure the client can receive the stop signal
	//	        // and stop listening after the stop signal processed
	//	        self.publish(
	//	            QueueStopEvent(stopped_by=QueueStopEvent.StopBy.USER_MANUAL), PublishFrom.TASK_PIPELINE
	//	        )
	//	    }
	//	    if elapsed_time // 10 > last_ping_time:
	//	        self.publish(QueuePingEvent(), PublishFrom.TASK_PIPELINE)
	//	        last_ping_time = elapsed_time // 10
	//	    }
	//	}
}

func (mgr *BaseAppQueueManager[T]) StopListen() {
	/*
		Stop listen to queue
		:return:
	*/
	if mgr.MsgQueue != nil {
		mgr.MsgQueue <- nil
	}
}
func (mgr *BaseAppQueueManager[T]) generateTaskBelongCacheKey(task_id string) string {
	/*
	   Generate task belong cache key
	   :param task_id: task id
	   :return:
	*/
	return "generate_task_belong:" + task_id
}

func (mgr *BaseAppQueueManager[T]) generateStoppedCacheKey(task_id string) string {
	/*
	   Generate stopped cache key
	   :param task_id: task id
	   :return:
	*/
	return "generate_task_stopped:" + task_id
}

func (mgr *BaseAppQueueManager[T]) SetStopFlag(task_id string, invoke_from appenumtypes.InvokeFrom, user_id string) {
	/*
	   Set task stop flag
	   :return:
	*/
	result := cache.Instance().GetString(mgr.generateTaskBelongCacheKey(task_id))
	if result == "" {
		return
	}
	var user_prefix string
	if slices.Contains([]appenumtypes.InvokeFrom{appenumtypes.InvokeFrom_EXPLORE, appenumtypes.InvokeFrom_DEBUGGER}, invoke_from) {
		user_prefix = "account"
	} else {
		user_prefix = "end-user"
	}

	if result != user_prefix+"-"+user_id {
		return
	}
	stopped_cache_key := mgr.generateStoppedCacheKey(task_id)
	cache.Instance().SetEx(stopped_cache_key, 1, 600*time.Second)
}
func (mgr *BaseAppQueueManager[T]) IsStopped() bool {
	/*
	   Check if task is stopped
	   :return:
	*/
	stopped_cache_key := mgr.generateStoppedCacheKey(mgr.TaskID)
	result := cache.Instance().GetString(stopped_cache_key)
	if result != "" {
		return true
	}
	return false
}

func (mgr *BaseAppQueueManager[T]) PublishError(aqm appqueueentities.AppQueueManager[T], exp error, pub_from appenumtypes.PublishFrom) {
	/*
	   Publish error
	   :param e: error
	   :param pub_from: publish from
	   :return:
	*/
	aqm.Publish(&appqueueentities.QueueErrorEvent{Err: exp}, pub_from)
}
