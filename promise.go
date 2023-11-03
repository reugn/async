package async

import "sync"

// Promise represents a writable, single-assignment container,
// which completes a Future.
type Promise[T any] interface {

	// Success completes the underlying Future with a value.
	Success(*T)

	// Failure fails the underlying Future with an error.
	Failure(error)

	// Future returns the underlying Future.
	Future() Future[T]
}

type promiseStatus uint8

const (
	ready promiseStatus = iota
	completed
)

// PromiseImpl implements the Promise interface.
type PromiseImpl[T any] struct {
	sync.Mutex
	future Future[T]
	status promiseStatus
}

// Verify PromiseImpl satisfies the Promise interface.
var _ Promise[any] = (*PromiseImpl[any])(nil)

// NewPromise returns a new PromiseImpl.
func NewPromise[T any]() Promise[T] {
	return &PromiseImpl[T]{
		future: NewFuture[T](),
		status: ready,
	}
}

// Success completes the underlying Future with a given value.
func (p *PromiseImpl[T]) Success(value *T) {
	p.Lock()
	defer p.Unlock()

	if p.status != completed {
		p.future.complete(value, nil)
		p.status = completed
	}
}

// Failure fails the underlying Future with a given error.
func (p *PromiseImpl[T]) Failure(err error) {
	p.Lock()
	defer p.Unlock()

	if p.status != completed {
		p.future.complete(nil, err)
		p.status = completed
	}
}

// Future returns the underlying Future.
func (p *PromiseImpl[T]) Future() Future[T] {
	return p.future
}
