package async

import (
	"fmt"
	"sync"
)

// Once is an object that will execute the given function exactly once.
// Any subsequent call will return the previous result.
type Once[T any] struct {
	runOnce sync.Once
	result  T
	err     error
}

// Do calls the function f if and only if Do is being called for the
// first time for this instance of Once. In other words, given
//
//	var once Once[T]
//
// if once.Do(f) is called multiple times, only the first call will invoke f,
// even if f has a different value in each invocation. A new instance of
// Once is required for each function to execute.
//
// The return values for each subsequent call will be the result of the
// first execution.
//
// If f panics, Do considers it to have returned; future calls of Do return
// without calling f.
func (o *Once[T]) Do(f func() (T, error)) (T, error) {
	o.runOnce.Do(func() {
		defer func() {
			if err := recover(); err != nil {
				o.err = fmt.Errorf("recovered %v", err)
			}
		}()
		o.result, o.err = f()
	})
	return o.result, o.err
}
