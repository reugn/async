package async

import (
	"fmt"
	"time"
)

// FutureSeq reduces many Futures into a single Future.
// The resulting array may contain both T values and errors.
func FutureSeq[T any](futures []Future[T]) Future[[]any] {
	next := newFuture[[]any]()
	go func() {
		seq := make([]any, len(futures))
		for i, future := range futures {
			result, err := future.Join()
			if err != nil {
				seq[i] = err
			} else {
				seq[i] = result
			}
		}
		next.complete(seq, nil)
	}()
	return next
}

// FutureFirstCompletedOf asynchronously returns a new Future to the result
// of the first Future in the list that is completed.
// This means no matter if it is completed as a success or as a failure.
func FutureFirstCompletedOf[T any](futures ...Future[T]) Future[T] {
	next := newFuture[T]()
	go func() {
		for _, f := range futures {
			go func(future Future[T]) {
				next.complete(future.Join())
			}(f)
		}
	}()
	return next
}

// FutureTimer returns Future that will have been resolved after given duration;
// useful for FutureFirstCompletedOf for timeout purposes.
func FutureTimer[T any](d time.Duration) Future[T] {
	next := newFuture[T]()
	go func() {
		<-time.After(d)
		var zero T
		next.(*futureImpl[T]).
			complete(zero, fmt.Errorf("future timeout after %s", d))
	}()
	return next
}
