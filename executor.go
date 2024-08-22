package async

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
)

// ExecutorStatus represents the status of an [ExecutorService].
type ExecutorStatus uint32

const (
	ExecutorStatusRunning ExecutorStatus = iota
	ExecutorStatusTerminating
	ExecutorStatusShutdown
)

var (
	ErrExecutorQueueFull = errors.New("async: executor queue is full")
	ErrExecutorShutdown  = errors.New("async: executor is shut down")
)

// ExecutorService is an interface that defines a task executor.
type ExecutorService[T any] interface {
	// Submit submits a function to the executor service.
	// The function will be executed asynchronously and the result will be
	// available via the returned future.
	Submit(func(context.Context) (T, error)) (Future[T], error)

	// Shutdown shuts down the executor service.
	// Once the executor service is shut down, no new tasks can be submitted
	// and any pending tasks will be cancelled.
	Shutdown() error

	// Status returns the current status of the executor service.
	Status() ExecutorStatus
}

// ExecutorConfig represents the Executor configuration.
type ExecutorConfig struct {
	WorkerPoolSize int
	QueueSize      int
}

// NewExecutorConfig returns a new [ExecutorConfig].
func NewExecutorConfig(workerPoolSize, queueSize int) *ExecutorConfig {
	return &ExecutorConfig{
		WorkerPoolSize: workerPoolSize,
		QueueSize:      queueSize,
	}
}

// Executor implements the [ExecutorService] interface.
type Executor[T any] struct {
	cancel context.CancelFunc
	queue  chan job[T]
	status atomic.Uint32
}

var _ ExecutorService[any] = (*Executor[any])(nil)

type job[T any] struct {
	promise Promise[T]
	task    func(context.Context) (T, error)
}

// NewExecutor returns a new [Executor].
func NewExecutor[T any](ctx context.Context, config *ExecutorConfig) *Executor[T] {
	ctx, cancel := context.WithCancel(ctx)
	executor := &Executor[T]{
		cancel: cancel,
		queue:  make(chan job[T], config.QueueSize),
	}
	// init the workers pool
	go executor.startWorkers(ctx, config.WorkerPoolSize)

	// set status to terminating when ctx is done
	go executor.monitorCtx(ctx)

	// set the executor status to running
	executor.status.Store(uint32(ExecutorStatusRunning))

	return executor
}

func (e *Executor[T]) monitorCtx(ctx context.Context) {
	<-ctx.Done()
	e.status.Store(uint32(ExecutorStatusTerminating))
}

func (e *Executor[T]) startWorkers(ctx context.Context, poolSize int) {
	var wg sync.WaitGroup
	for i := 0; i < poolSize; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
		loop:
			for ExecutorStatus(e.status.Load()) == ExecutorStatusRunning {
				select {
				case job := <-e.queue:
					result, err := job.task(ctx)
					if err != nil {
						job.promise.Failure(err)
					} else {
						job.promise.Success(result)
					}
				case <-ctx.Done():
					break loop
				}
			}
		}()
	}

	// wait for all workers to exit
	wg.Wait()
	// close the queue and cancel all pending tasks
	close(e.queue)
	for job := range e.queue {
		job.promise.Failure(ErrExecutorShutdown)
	}
	// mark the executor as shut down
	e.status.Store(uint32(ExecutorStatusShutdown))
}

// Submit submits a function to the executor.
// The function will be executed asynchronously and the result will be
// available via the returned future.
func (e *Executor[T]) Submit(f func(context.Context) (T, error)) (Future[T], error) {
	promise := NewPromise[T]()
	if ExecutorStatus(e.status.Load()) == ExecutorStatusRunning {
		select {
		case e.queue <- job[T]{promise, f}:
		default:
			return nil, ErrExecutorQueueFull
		}
	} else {
		return nil, ErrExecutorShutdown
	}
	return promise.Future(), nil
}

// Shutdown shuts down the executor.
// Once the executor service is shut down, no new tasks can be submitted
// and any pending tasks will be cancelled.
func (e *Executor[T]) Shutdown() error {
	e.cancel()
	return nil
}

// Status returns the current status of the executor.
func (e *Executor[T]) Status() ExecutorStatus {
	return ExecutorStatus(e.status.Load())
}
