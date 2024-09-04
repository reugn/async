package async

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
)

// ExecutorStatus represents the status of an [ExecutorService].
type ExecutorStatus uint32

const (
	ExecutorStatusRunning ExecutorStatus = iota
	ExecutorStatusTerminating
	ExecutorStatusShutDown
)

var (
	ErrExecutorQueueFull = errors.New("async: executor queue is full")
	ErrExecutorShutDown  = errors.New("async: executor is shut down")
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
// workerPoolSize must be positive and queueSize non-negative.
func NewExecutorConfig(workerPoolSize, queueSize int) *ExecutorConfig {
	return &ExecutorConfig{
		WorkerPoolSize: workerPoolSize,
		QueueSize:      queueSize,
	}
}

// Executor implements the [ExecutorService] interface.
type Executor[T any] struct {
	mtx    sync.RWMutex
	cancel context.CancelFunc
	queue  chan executorJob[T]
	status atomic.Uint32
}

var _ ExecutorService[any] = (*Executor[any])(nil)

type executorJob[T any] struct {
	promise Promise[T]
	task    func(context.Context) (T, error)
}

// run executes the task, handling possible panics.
func (job *executorJob[T]) run(ctx context.Context) (result T, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered: %v", r)
		}
	}()
	return job.task(ctx)
}

// NewExecutor returns a new [Executor].
func NewExecutor[T any](ctx context.Context, config *ExecutorConfig) *Executor[T] {
	ctx, cancel := context.WithCancel(ctx)
	executor := &Executor[T]{
		cancel: cancel,
		queue:  make(chan executorJob[T], config.QueueSize),
	}
	// set the executor status to running explicitly
	executor.status.Store(uint32(ExecutorStatusRunning))

	// init the workers pool
	go executor.startWorkers(ctx, config.WorkerPoolSize)

	// set status to terminating when ctx is done
	go executor.monitorCtx(ctx)

	return executor
}

func (e *Executor[T]) monitorCtx(ctx context.Context) {
	<-ctx.Done()
	_ = e.status.CompareAndSwap(uint32(ExecutorStatusRunning),
		uint32(ExecutorStatusTerminating))
}

func (e *Executor[T]) startWorkers(ctx context.Context, poolSize int) {
	var wg sync.WaitGroup
	for i := 0; i < poolSize; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
		loop:
			// check the status to break the loop even if the queue is not empty
			for ExecutorStatus(e.status.Load()) == ExecutorStatusRunning {
				select {
				case job := <-e.queue:
					result, err := job.run(ctx)
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
	// mark the executor as terminating
	e.status.Store(uint32(ExecutorStatusTerminating))

	// avoid submissions while draining the queue
	e.mtx.Lock()
	defer e.mtx.Unlock()

	// close the queue and cancel all pending tasks
	close(e.queue)
	for job := range e.queue {
		job.promise.Failure(ErrExecutorShutDown)
	}
	// mark the executor as shut down
	e.status.Store(uint32(ExecutorStatusShutDown))
}

// Submit submits a function to the executor.
// The function will be executed asynchronously and the result will be
// available via the returned future.
func (e *Executor[T]) Submit(f func(context.Context) (T, error)) (Future[T], error) {
	e.mtx.RLock()
	defer e.mtx.RUnlock()

	if ExecutorStatus(e.status.Load()) == ExecutorStatusRunning {
		promise := NewPromise[T]()
		select {
		case e.queue <- executorJob[T]{promise, f}:
			return promise.Future(), nil
		default:
			return nil, ErrExecutorQueueFull
		}
	}
	return nil, ErrExecutorShutDown
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
