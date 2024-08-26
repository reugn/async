package async

import (
	"context"
	"errors"
	"runtime"
	"testing"
	"time"

	"github.com/reugn/async/internal/assert"
)

func TestExecutor(t *testing.T) {
	ctx := context.Background()
	executor := NewExecutor[int](ctx, NewExecutorConfig(2, 2))

	job := func(_ context.Context) (int, error) {
		time.Sleep(time.Millisecond)
		return 1, nil
	}
	jobLong := func(_ context.Context) (int, error) {
		time.Sleep(10 * time.Millisecond)
		return 1, nil
	}

	future1 := submitJob[int](t, executor, job)
	future2 := submitJob[int](t, executor, job)

	// wait for the first two jobs to complete
	time.Sleep(3 * time.Millisecond)

	// submit four more jobs
	future3 := submitJob[int](t, executor, jobLong)
	future4 := submitJob[int](t, executor, jobLong)
	future5 := submitJob[int](t, executor, jobLong)
	future6 := submitJob[int](t, executor, jobLong)

	// the queue has reached its maximum capacity
	future7, err := executor.Submit(job)
	assert.ErrorIs(t, err, ErrExecutorQueueFull)
	assert.IsNil(t, future7)

	assert.Equal(t, executor.Status(), ExecutorStatusRunning)

	routines := runtime.NumGoroutine()

	// shut down the executor
	_ = executor.Shutdown()
	time.Sleep(time.Millisecond)

	// verify that submit fails after the executor was shut down
	_, err = executor.Submit(job)
	assert.ErrorIs(t, err, ErrExecutorShutdown)

	// validate the executor status
	assert.Equal(t, executor.Status(), ExecutorStatusTerminating)
	time.Sleep(10 * time.Millisecond)
	assert.Equal(t, executor.Status(), ExecutorStatusShutdown)

	assert.Equal(t, routines, runtime.NumGoroutine()+4)

	assertFutureResult(t, 1, future1, future2, future3, future4)
	assertFutureError(t, ErrExecutorShutdown, future5, future6)
}

func TestExecutor_context(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	executor := NewExecutor[int](ctx, NewExecutorConfig(2, 2))

	job := func(_ context.Context) (int, error) {
		return 0, errors.New("error")
	}

	future, err := executor.Submit(job)
	assert.IsNil(t, err)

	result, err := future.Join()
	assert.Equal(t, result, 0)
	assert.ErrorContains(t, err, "error")

	cancel()
	time.Sleep(5 * time.Millisecond)
	assert.Equal(t, executor.Status(), ExecutorStatusShutdown)
}

func submitJob[T any](t *testing.T, executor ExecutorService[T],
	f func(context.Context) (T, error)) Future[T] {
	future, err := executor.Submit(f)
	assert.IsNil(t, err)

	time.Sleep(time.Millisecond) // switch context
	return future
}

func assertFutureResult[T any](t *testing.T, expected T, futures ...Future[T]) {
	for _, future := range futures {
		result, err := future.Join()
		assert.IsNil(t, err)
		assert.Equal(t, expected, result)
	}
}

func assertFutureError[T any](t *testing.T, expected error, futures ...Future[T]) {
	for _, future := range futures {
		result, err := future.Join()
		var zero T
		assert.Equal(t, zero, result)
		assert.ErrorIs(t, err, expected)
	}
}
