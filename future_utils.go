package async

import (
	"fmt"
	"time"
)

// FutureSeq reduces many Futures into a single Future.
func FutureSeq[T any](futures []Future[T]) Future[[]any] {
	next := NewFuture[[]any]()
	go func() {
		seq := make([]any, len(futures))
		for i, future := range futures {
			res, err := future.Join()
			if err != nil {
				seq[i] = err
			} else {
				seq[i] = res
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
	next := NewFuture[T]()
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
	next := NewFuture[T]()
	go func() {
		timer := time.NewTimer(d)
		<-timer.C
		var nilT T
		next.complete(nilT, fmt.Errorf("FutureTimer %s timeout", d))
	}()
	return next
}
