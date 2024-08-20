package async

// Task is a data type for controlling possibly lazy and
// asynchronous computations.
type Task[T any] struct {
	taskFunc func() (T, error)
}

// NewTask returns a new Task associated with the specified function.
func NewTask[T any](taskFunc func() (T, error)) *Task[T] {
	return &Task[T]{
		taskFunc: taskFunc,
	}
}

// Call starts executing the task using a goroutine. It returns a
// Future which can be used to retrieve the result or error of the
// task when it is completed.
func (task *Task[T]) Call() Future[T] {
	promise := NewPromise[T]()
	go func() {
		result, err := task.taskFunc()
		if err == nil {
			promise.Success(result)
		} else {
			promise.Failure(err)
		}
	}()
	return promise.Future()
}
