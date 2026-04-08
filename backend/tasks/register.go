package tasks

import (
	"fmt"
	"time"

	"mlib.com/mlog"
)

var (
	defaultQueue     *Queue
	defaultScheduler *Scheduler
)

// InitTasks initializes and starts the task queue and scheduler.
func InitTasks(workers int) {
	// Initialize queue
	defaultQueue = NewQueue(workers)

	// Register task handlers
	defaultQueue.Register(TaskDocumentIndexing, HandleDocumentIndexing)
	defaultQueue.Register(TaskDocumentIndexingUpdate, HandleDocumentIndexingUpdate)
	defaultQueue.Register(TaskDocumentClean, HandleDocumentClean)
	defaultQueue.Register(TaskSegmentCreate, HandleSegmentCreate)
	defaultQueue.Register(TaskSegmentDelete, HandleSegmentDelete)
	defaultQueue.Register(TaskSegmentEnable, HandleSegmentEnable)
	defaultQueue.Register(TaskSegmentDisable, HandleSegmentDisable)

	// Workflow tasks
	defaultQueue.Register(TaskWorkflowExecution, HandleWorkflowExecution)
	defaultQueue.Register(TaskWorkflowNodeExecution, HandleWorkflowNodeExecution)

	// Mail tasks
	defaultQueue.Register(TaskMailRegistration, HandleMailRegistration)
	defaultQueue.Register(TaskMailResetPassword, HandleMailResetPassword)
	defaultQueue.Register(TaskMailInviteMember, HandleMailInviteMember)

	// Start queue workers
	defaultQueue.Start()

	// Initialize scheduler
	defaultScheduler = NewScheduler()

	// Register scheduled tasks
	defaultScheduler.AddTask("clean_messages", 24*time.Hour, CleanOldMessages)
	defaultScheduler.AddTask("clean_workflow_runs", 24*time.Hour, CleanOldWorkflowRuns)
	defaultScheduler.AddTask("clean_unused_datasets", 7*24*time.Hour, CleanUnusedDatasets)

	// Start scheduler
	defaultScheduler.Start()

	mlog.Info("task queue and scheduler initialized")
}

// StopTasks gracefully stops the task queue and scheduler.
func StopTasks() {
	if defaultQueue != nil {
		defaultQueue.Stop()
	}
	if defaultScheduler != nil {
		defaultScheduler.Stop()
	}
}

// EnqueueTask enqueues a task to the default queue.
func EnqueueTask(taskType string, payload any) error {
	if defaultQueue == nil {
		return fmt.Errorf("task queue not initialized")
	}
	return defaultQueue.Enqueue(taskType, payload)
}
