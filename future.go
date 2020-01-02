package async

import "sync"

// Future represents a value which may or may not currently be available,
// but will be available at some point, or an error if that value could not be made available.
type Future interface {

	// Creates a new future by applying a function to the successful result of this future.
	Map(func(interface{}) (interface{}, error)) Future

	// Creates a new future by applying a function to the successful result of
	// this future, and returns the result of the function as the new future.
	FlatMap(func(interface{}) (Future, error)) Future

	// Blocks until Future completed and return either result or error.
	Get() (interface{}, error)

	// Creates a new future that will handle any error that this
	// future might contain. If this future contains
	// a valid result then the new future will contain the same.
	Recover(func() (interface{}, error)) Future

	// Creates a new future that will handle any error that this
	// future might contain by assigning it a value of another future.
	// If this future contains a valid result then the new future will contain the same result.
	RecoverWith(Future) Future

	// Complete future with either result or error
	// For Promise use internally
	complete(interface{}, error)
}

// FutureImpl Future implementation
type FutureImpl struct {
	acc   sync.Once
	compl sync.Once
	done  chan interface{}
	value interface{}
	err   error
}

// NewFuture returns new Future
func NewFuture() Future {
	return &FutureImpl{
		done: make(chan interface{}),
	}
}

// accept blocks once until result is available
func (fut *FutureImpl) accept() {
	fut.acc.Do(func() {
		sig := <-fut.done
		switch v := sig.(type) {
		case error:
			fut.err = v
		default:
			fut.value = v
		}
	})
}

// Map default implementation
func (fut *FutureImpl) Map(f func(interface{}) (interface{}, error)) Future {
	next := NewFuture()
	go func() {
		fut.accept()
		if fut.err != nil {
			next.complete(nil, fut.err)
		} else {
			next.complete(f(fut.value))
		}
	}()
	return next
}

// FlatMap default implementation
func (fut *FutureImpl) FlatMap(f func(interface{}) (Future, error)) Future {
	next := NewFuture()
	go func() {
		fut.accept()
		if fut.err != nil {
			next.complete(nil, fut.err)
		} else {
			tfut, terr := f(fut.value)
			if terr != nil {
				next.complete(nil, terr)
			} else {
				next.complete(tfut.Get())
			}
		}
	}()
	return next
}

// Get default implementation
func (fut *FutureImpl) Get() (interface{}, error) {
	fut.accept()
	return fut.value, fut.err
}

// Recover default implementation
func (fut *FutureImpl) Recover(f func() (interface{}, error)) Future {
	fut.accept()
	if fut.err != nil {
		next := NewFuture()
		next.complete(f())
		return next
	}
	return fut
}

// RecoverWith default implementation
func (fut *FutureImpl) RecoverWith(rf Future) Future {
	fut.accept()
	if fut.err != nil {
		next := NewFuture()
		next.complete(rf.Get())
		return next
	}
	return fut
}

// complete future with either value or error
func (fut *FutureImpl) complete(v interface{}, e error) {
	fut.compl.Do(func() {
		go func() {
			if e != nil {
				fut.done <- e
			} else {
				fut.done <- v
			}
		}()
	})
}
