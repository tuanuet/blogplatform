package service

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestTaskRunner_Submit(t *testing.T) {
	runner := NewTaskRunner(1 * time.Second)

	var wg sync.WaitGroup
	wg.Add(1)

	taskExecuted := false
	runner.Submit(func(ctx context.Context) {
		taskExecuted = true
		wg.Done()
	})

	wg.Wait()

	if !taskExecuted {
		t.Error("Task was not executed")
	}
}

func TestTaskRunner_Timeout(t *testing.T) {
	runner := NewTaskRunner(10 * time.Millisecond)

	var wg sync.WaitGroup
	wg.Add(1)

	ctxCancelled := false
	runner.Submit(func(ctx context.Context) {
		select {
		case <-time.After(100 * time.Millisecond):
		case <-ctx.Done():
			ctxCancelled = true
		}
		wg.Done()
	})

	wg.Wait()

	if !ctxCancelled {
		t.Error("Context was not cancelled on timeout")
	}
}
