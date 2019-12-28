package async

import (
	"fmt"
	"time"
)

// FutureSeq reduces many Futures into a single Future
func FutureSeq(futures []Future) Future {
	next := NewFuture()
	go func() {
		seq := make([]interface{}, len(futures))
		for i, f := range futures {
			v, e := f.Get()
			if e != nil {
				seq[i] = e
			} else {
				seq[i] = v
			}
		}
		next.complete(seq, nil)
	}()
	return next
}

// FutureFirstCompletedOf asynchronously returns a new Future to the result of the first Future
// in the list that is completed. This means no matter if it is completed as a success or as a failure.
func FutureFirstCompletedOf(futures ...Future) Future {
	next := NewFuture()
	go func() {
		for _, f := range futures {
			go func(future Future) {
				next.complete(future.Get())
			}(f)
		}
	}()
	return next
}

// FutureTimer returns Future that will have been resolved after given duration
// useful for FutureFirstCompletedOf for timeout purposes
func FutureTimer(d time.Duration) Future {
	next := NewFuture()
	go func() {
		timer := time.NewTimer(d)
		<-timer.C
		next.complete(nil, fmt.Errorf("FutureTimer %v timeout", d))
	}()
	return next
}
