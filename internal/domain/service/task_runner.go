package service

import (
	"context"
	"time"
)

// TaskRunner defines the interface for running asynchronous tasks
type TaskRunner interface {
	Submit(task func(ctx context.Context))
}

type taskRunner struct {
	timeout time.Duration
}

// NewTaskRunner creates a new TaskRunner with the given timeout for each task
func NewTaskRunner(timeout time.Duration) TaskRunner {
	return &taskRunner{
		timeout: timeout,
	}
}

// Submit runs the task in a background goroutine
func (r *taskRunner) Submit(task func(ctx context.Context)) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
		defer cancel()

		done := make(chan struct{})
		go func() {
			task(ctx)
			close(done)
		}()

		select {
		case <-done:
			// Task completed successfully
		case <-ctx.Done():
			// Task timed out
			// Note: It's up to the task to respect the context and stop.
		}
	}()
}
