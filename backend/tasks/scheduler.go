package tasks

import (
	"sync"
	"time"

	"github.com/odysseythink/mlog"
)

// ScheduledTask represents a recurring task.
type ScheduledTask struct {
	Name     string
	Interval time.Duration
	Handler  func()
	ticker   *time.Ticker
	quit     chan struct{}
}

// Scheduler manages recurring tasks.
type Scheduler struct {
	tasks []*ScheduledTask
	mu    sync.Mutex
}

// NewScheduler creates a new scheduler.
func NewScheduler() *Scheduler {
	return &Scheduler{}
}

// AddTask adds a recurring task.
func (s *Scheduler) AddTask(name string, interval time.Duration, handler func()) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tasks = append(s.tasks, &ScheduledTask{
		Name:     name,
		Interval: interval,
		Handler:  handler,
		quit:     make(chan struct{}),
	})
}

// Start begins all scheduled tasks.
func (s *Scheduler) Start() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, task := range s.tasks {
		t := task
		t.ticker = time.NewTicker(t.Interval)
		go func() {
			mlog.Info("scheduler: starting task '%s' every %v", t.Name, t.Interval)
			for {
				select {
				case <-t.quit:
					t.ticker.Stop()
					return
				case <-t.ticker.C:
					func() {
						defer func() {
							if r := recover(); r != nil {
								mlog.Errorf("scheduler: task '%s' panicked: %v", t.Name, r)
							}
						}()
						t.Handler()
					}()
				}
			}
		}()
	}
}

// Stop stops all scheduled tasks.
func (s *Scheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, task := range s.tasks {
		close(task.quit)
	}
}
