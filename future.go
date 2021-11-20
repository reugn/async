package async

import "sync"

// Future represents a value which may or may not currently be available,
// but will be available at some point, or an error if that value could not be made available.
type Future interface {

	// Map creates a new Future by applying a function to the successful result of this Future.
	Map(func(interface{}) (interface{}, error)) Future

	// FlatMap creates a new Future by applying a function to the successful result of
	// this Future.
	FlatMap(func(interface{}) (Future, error)) Future

	// Get blocks until the Future is completed and returns either a result or an error.
	Get() (interface{}, error)

	// Recover handles any error that this Future might contain using a resolver function.
	Recover(func() (interface{}, error)) Future

	// RecoverWith handles any error that this Future might contain using another Future.
	RecoverWith(Future) Future

	// complete completes the Future with either a value or an error.
	// Is used by Promise internally.
	complete(interface{}, error)
}

// FutureImpl implements the Future interface.
type FutureImpl struct {
	acc   sync.Once
	compl sync.Once
	done  chan interface{}
	value interface{}
	err   error
}

// NewFuture returns a new Future.
func NewFuture() Future {
	return &FutureImpl{
		done: make(chan interface{}),
	}
}

// accept blocks once, until the Future result is available.
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

// Map creates a new Future by applying a function to the successful result of this Future
// and returns the result of the function as a new Future.
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

// FlatMap creates a new Future by applying a function to the successful result of
// this Future and returns the result of the function as a new Future.
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

// Get blocks until the Future is completed and returns either a result or an error.
func (fut *FutureImpl) Get() (interface{}, error) {
	fut.accept()
	return fut.value, fut.err
}

// Recover handles any error that this Future might contain using a given resolver function.
// Returns the result as a new Future.
func (fut *FutureImpl) Recover(f func() (interface{}, error)) Future {
	next := NewFuture()
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

// RecoverWith handles any error that this Future might contain using another Future.
// Returns the result as a new Future.
func (fut *FutureImpl) RecoverWith(rf Future) Future {
	next := NewFuture()
	go func() {
		fut.accept()
		if fut.err != nil {
			next.complete(rf.Get())
		} else {
			next.complete(fut.value, nil)
		}
	}()
	return next
}

// complete completes the Future with either a value or an error.
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
