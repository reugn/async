//go:build go1.18
// +build go1.18

package generic

import (
	"github.com/reugn/async"
)

// AsyncTask is a data type for controlling possibly lazy & asynchronous computations.
type AsyncTask[T any] struct {
	taskFunc func() (T, error)
}

// NewAsyncTask returns a new AsyncTask.
func NewAsyncTask[T any](taskFunc func() (T, error)) *AsyncTask[T] {
	return &AsyncTask[T]{
		taskFunc: taskFunc,
	}
}

// Call executes the AsyncTask and returns a Future.
func (task *AsyncTask[T]) Call() async.Future {
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
