package async

// Task is a data type for controlling possibly lazy & asynchronous computations.
type Task[T any] struct {
	taskFunc func() (T, error)
}

// NewTask returns a new Task.
func NewTask[T any](taskFunc func() (T, error)) *Task[T] {
	return &Task[T]{
		taskFunc: taskFunc,
	}
}

// Call executes the Task and returns a Future.
func (task *Task[T]) Call() Future[T] {
	promise := NewPromise[T]()
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
