package async

import "sync"

// Promise represents a writable, single-assignment container,
// which completes a Future.
type Promise[T any] interface {

	// Success completes the underlying Future with a value.
	Success(T)

	// Failure fails the underlying Future with an error.
	Failure(error)

	// Future returns the underlying Future.
	Future() Future[T]
}

// promiseImpl implements the Promise interface.
type promiseImpl[T any] struct {
	once   sync.Once
	future Future[T]
}

// Verify promiseImpl satisfies the Promise interface.
var _ Promise[any] = (*promiseImpl[any])(nil)

// NewPromise returns a new Promise.
func NewPromise[T any]() Promise[T] {
	return &promiseImpl[T]{
		future: newFuture[T](),
	}
}

// Success completes the underlying Future with a given value.
func (p *promiseImpl[T]) Success(value T) {
	p.once.Do(func() {
		p.future.complete(value, nil)
	})
}

// Failure fails the underlying Future with a given error.
func (p *promiseImpl[T]) Failure(err error) {
	p.once.Do(func() {
		var zero T
		p.future.complete(zero, err)
	})
}

// Future returns the underlying Future.
func (p *promiseImpl[T]) Future() Future[T] {
	return p.future
}
