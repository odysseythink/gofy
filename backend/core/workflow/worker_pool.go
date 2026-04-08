package workflow

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"mlib.com/mlog"
)

// WorkerPool manages goroutine workers for parallel node execution.
type WorkerPool struct {
	mu          sync.Mutex
	maxWorkers  int
	activeCount int32
	taskCh      chan *WorkerTask
	quit        chan struct{}
	wg          sync.WaitGroup
}

// WorkerTask represents a unit of work for the pool.
type WorkerTask struct {
	ID      string
	Execute func(ctx context.Context) error
	OnDone  func(err error)
	Ctx     context.Context
}

// NewWorkerPool creates a worker pool with the specified max workers.
func NewWorkerPool(maxWorkers int) *WorkerPool {
	if maxWorkers <= 0 {
		maxWorkers = 10
	}
	return &WorkerPool{
		maxWorkers: maxWorkers,
		taskCh:     make(chan *WorkerTask, maxWorkers*2),
		quit:       make(chan struct{}),
	}
}

// Start launches all worker goroutines.
func (wp *WorkerPool) Start() {
	for i := 0; i < wp.maxWorkers; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
	mlog.Infof("worker pool started with %d workers", wp.maxWorkers)
}

// Submit adds a task to the pool.
func (wp *WorkerPool) Submit(task *WorkerTask) {
	select {
	case wp.taskCh <- task:
	case <-wp.quit:
		if task.OnDone != nil {
			task.OnDone(context.Canceled)
		}
	}
}

// SubmitAndWait submits a task and blocks until it completes.
func (wp *WorkerPool) SubmitAndWait(ctx context.Context, fn func(ctx context.Context) error) error {
	done := make(chan error, 1)
	wp.Submit(&WorkerTask{
		Ctx:     ctx,
		Execute: fn,
		OnDone:  func(err error) { done <- err },
	})
	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

// RunParallel executes multiple functions concurrently and waits for all to complete.
func (wp *WorkerPool) RunParallel(ctx context.Context, fns []func(ctx context.Context) error) []error {
	errs := make([]error, len(fns))
	var wg sync.WaitGroup
	wg.Add(len(fns))

	for i, fn := range fns {
		idx := i
		f := fn
		wp.Submit(&WorkerTask{
			Ctx:     ctx,
			Execute: f,
			OnDone: func(err error) {
				errs[idx] = err
				wg.Done()
			},
		})
	}

	wg.Wait()
	return errs
}

// ActiveWorkers returns the number of currently active workers.
func (wp *WorkerPool) ActiveWorkers() int {
	return int(atomic.LoadInt32(&wp.activeCount))
}

// Stop gracefully shuts down the pool.
func (wp *WorkerPool) Stop() {
	close(wp.quit)
	wp.wg.Wait()
	mlog.Info("worker pool stopped")
}

// StopWithTimeout stops the pool with a timeout.
func (wp *WorkerPool) StopWithTimeout(timeout time.Duration) {
	close(wp.quit)
	done := make(chan struct{})
	go func() {
		wp.wg.Wait()
		close(done)
	}()
	select {
	case <-done:
		mlog.Info("worker pool stopped gracefully")
	case <-time.After(timeout):
		mlog.Errorf("worker pool stop timed out after %v", timeout)
	}
}

func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()
	for {
		select {
		case <-wp.quit:
			return
		case task := <-wp.taskCh:
			atomic.AddInt32(&wp.activeCount, 1)
			wp.executeTask(id, task)
			atomic.AddInt32(&wp.activeCount, -1)
		}
	}
}

func (wp *WorkerPool) executeTask(workerID int, task *WorkerTask) {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = context.Canceled
			mlog.Errorf("worker %d: task panicked: %v", workerID, r)
		}
		if task.OnDone != nil {
			task.OnDone(err)
		}
	}()

	ctx := task.Ctx
	if ctx == nil {
		ctx = context.Background()
	}
	err = task.Execute(ctx)
}
