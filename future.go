package async

import (
	"fmt"
	"sync"
	"time"
)

// Future represents a value which may or may not currently be available,
// but will be available at some point, or an error if that value could
// not be made available.
type Future[T any] interface {

	// Map creates a new Future by applying a function to the successful
	// result of this Future.
	Map(func(T) (T, error)) Future[T]

	// FlatMap creates a new Future by applying a function to the successful
	// result of this Future.
	FlatMap(func(T) (Future[T], error)) Future[T]

	// Join blocks until the Future is completed and returns either a result
	// or an error.
	Join() (T, error)

	// Get blocks for at most the given time duration for this Future to
	// complete and returns either a result or an error.
	Get(time.Duration) (T, error)

	// Recover handles any error that this Future might contain using a
	// resolver function.
	Recover(func() (T, error)) Future[T]

	// RecoverWith handles any error that this Future might contain using
	// another Future.
	RecoverWith(Future[T]) Future[T]

	// complete completes the Future with either a value or an error.
	// It is used by [Promise] internally.
	complete(T, error)
}

// futureImpl implements the Future interface.
type futureImpl[T any] struct {
	acceptOnce   sync.Once
	completeOnce sync.Once
	done         chan any
	value        T
	err          error
}

// Verify futureImpl satisfies the Future interface.
var _ Future[any] = (*futureImpl[any])(nil)

// newFuture returns a new Future.
func newFuture[T any]() Future[T] {
	return &futureImpl[T]{
		done: make(chan any, 1),
	}
}

// accept blocks once, until the Future result is available.
func (fut *futureImpl[T]) accept() {
	fut.acceptOnce.Do(func() {
		result := <-fut.done
		fut.setResult(result)
	})
}

// acceptTimeout blocks once, until the Future result is available or until
// the timeout expires.
func (fut *futureImpl[T]) acceptTimeout(timeout time.Duration) {
	fut.acceptOnce.Do(func() {
		timer := time.NewTimer(timeout)
		defer timer.Stop()
		select {
		case result := <-fut.done:
			fut.setResult(result)
		case <-timer.C:
			fut.setResult(fmt.Errorf("Future timeout after %s", timeout))
		}
	})
}

// setResult assigns a value to the Future instance.
func (fut *futureImpl[T]) setResult(result any) {
	switch value := result.(type) {
	case error:
		fut.err = value
	default:
		fut.value = value.(T)
	}
}

// Map creates a new Future by applying a function to the successful result
// of this Future and returns the result of the function as a new Future.
func (fut *futureImpl[T]) Map(f func(T) (T, error)) Future[T] {
	next := newFuture[T]()
	go func() {
		fut.accept()
		if fut.err != nil {
			var zero T
			next.complete(zero, fut.err)
		} else {
			next.complete(f(fut.value))
		}
	}()
	return next
}

// FlatMap creates a new Future by applying a function to the successful result
// of this Future and returns the result of the function as a new Future.
func (fut *futureImpl[T]) FlatMap(f func(T) (Future[T], error)) Future[T] {
	next := newFuture[T]()
	go func() {
		fut.accept()
		if fut.err != nil {
			var zero T
			next.complete(zero, fut.err)
		} else {
			tfut, terr := f(fut.value)
			if terr != nil {
				var zero T
				next.complete(zero, terr)
			} else {
				next.complete(tfut.Join())
			}
		}
	}()
	return next
}

// Join blocks until the Future is completed and returns either
// a result or an error.
func (fut *futureImpl[T]) Join() (T, error) {
	fut.accept()
	return fut.value, fut.err
}

// Get blocks for at most the given time duration for this Future to
// complete and returns either a result or an error.
func (fut *futureImpl[T]) Get(timeout time.Duration) (T, error) {
	fut.acceptTimeout(timeout)
	return fut.value, fut.err
}

// Recover handles any error that this Future might contain using
// a given resolver function.
// Returns the result as a new Future.
func (fut *futureImpl[T]) Recover(f func() (T, error)) Future[T] {
	next := newFuture[T]()
	go func() {
		fut.accept()
		if fut.err != nil {
			next.complete(f())
		} else {
			next.complete(fut.value, nil)
		}
	}()
	return next
}

// RecoverWith handles any error that this Future might contain using
// another Future.
// Returns the result as a new Future.
func (fut *futureImpl[T]) RecoverWith(rf Future[T]) Future[T] {
	next := newFuture[T]()
	go func() {
		fut.accept()
		if fut.err != nil {
			next.complete(rf.Join())
		} else {
			next.complete(fut.value, nil)
		}
	}()
	return next
}

// complete completes the Future with either a value or an error.
func (fut *futureImpl[T]) complete(value T, err error) {
	fut.completeOnce.Do(func() {
		if err != nil {
			fut.done <- err
		} else {
			fut.done <- value
		}
	})
}
