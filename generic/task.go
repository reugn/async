//go:build go1.18
// +build go1.18

package generic

import (
	"github.com/reugn/async"
)

// asyncTask is a data type for controlling possibly lazy & asynchronous computations.
type asyncTask[T any] struct {
	taskFunc func() (T, error)
}

func newAsyncTask[T any](taskFunc func() (T, error)) *asyncTask[T] {
	return &asyncTask[T]{
		taskFunc: taskFunc,
	}
}

func (task *asyncTask[T]) call() async.Future {
	promise := async.NewPromise()
	go func() {
		res, err := task.taskFunc()
		if err == nil {
			promise.Success(res)
		} else {
			promise.Failure(err)
		}
	}()
	return promise.Future()
}
