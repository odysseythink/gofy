package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"mlib.com/mlog"
)

// Task represents a background task.
type Task struct {
	Type      string          `json:"type"`
	Payload   json.RawMessage `json:"payload"`
	ID        string          `json:"id"`
	Retries   int             `json:"retries"`
	MaxRetry  int             `json:"max_retry"`
	CreatedAt time.Time       `json:"created_at"`
}

// HandlerFunc is a function that processes a task.
type HandlerFunc func(ctx context.Context, task *Task) error

// Queue manages task enqueueing and processing.
type Queue struct {
	handlers map[string]HandlerFunc
	workers  int
	quit     chan struct{}
}

// NewQueue creates a new task queue.
func NewQueue(workers int) *Queue {
	return &Queue{
		handlers: make(map[string]HandlerFunc),
		workers:  workers,
		quit:     make(chan struct{}),
	}
}

// Register registers a handler for a task type.
func (q *Queue) Register(taskType string, handler HandlerFunc) {
	q.handlers[taskType] = handler
}

// Enqueue adds a task to the queue.
func (q *Queue) Enqueue(taskType string, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	task := &Task{
		Type:      taskType,
		Payload:   data,
		ID:        fmt.Sprintf("%s-%d", taskType, time.Now().UnixNano()),
		MaxRetry:  3,
		CreatedAt: time.Now(),
	}

	taskBytes, err := json.Marshal(task)
	if err != nil {
		return err
	}

	// Push to Redis list using the existing cache/redis client
	return pushToRedis("gofy:tasks:queue", string(taskBytes))
}

// EnqueueDelayed enqueues a task to be processed after a delay.
func (q *Queue) EnqueueDelayed(taskType string, payload any, delay time.Duration) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	task := &Task{
		Type:      taskType,
		Payload:   data,
		ID:        fmt.Sprintf("%s-%d", taskType, time.Now().UnixNano()),
		MaxRetry:  3,
		CreatedAt: time.Now(),
	}

	taskBytes, err := json.Marshal(task)
	if err != nil {
		return err
	}

	score := float64(time.Now().Add(delay).Unix())
	return pushToRedisDelayed("gofy:tasks:delayed", string(taskBytes), score)
}

// Start begins processing tasks from the queue.
func (q *Queue) Start() {
	mlog.Info("task queue starting with %d workers", q.workers)
	for i := 0; i < q.workers; i++ {
		go q.worker(i)
	}
	go q.delayedPoller()
}

// Stop gracefully stops the queue.
func (q *Queue) Stop() {
	close(q.quit)
}

func (q *Queue) worker(id int) {
	for {
		select {
		case <-q.quit:
			mlog.Info("worker %d stopping", id)
			return
		default:
			taskStr, err := popFromRedis("gofy:tasks:queue", 2*time.Second)
			if err != nil || taskStr == "" {
				continue
			}

			var task Task
			if err := json.Unmarshal([]byte(taskStr), &task); err != nil {
				mlog.Errorf("worker %d: failed to unmarshal task: %v", id, err)
				continue
			}

			handler, ok := q.handlers[task.Type]
			if !ok {
				mlog.Errorf("worker %d: no handler for task type: %s", id, task.Type)
				continue
			}

			ctx := context.Background()
			if err := handler(ctx, &task); err != nil {
				mlog.Errorf("worker %d: task %s failed: %v", id, task.ID, err)
				task.Retries++
				if task.Retries < task.MaxRetry {
					q.retryTask(&task)
				} else {
					mlog.Errorf("worker %d: task %s exceeded max retries", id, task.ID)
				}
			}
		}
	}
}

func (q *Queue) retryTask(task *Task) {
	taskBytes, _ := json.Marshal(task)
	delay := time.Duration(task.Retries*task.Retries) * time.Second // exponential backoff
	score := float64(time.Now().Add(delay).Unix())
	if err := pushToRedisDelayed("gofy:tasks:delayed", string(taskBytes), score); err != nil {
		mlog.Errorf("retryTask: failed to push delayed task %s: %v", task.ID, err)
	}
}

func (q *Queue) delayedPoller() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-q.quit:
			return
		case <-ticker.C:
			moveDelayedToReady("gofy:tasks:delayed", "gofy:tasks:queue")
		}
	}
}
