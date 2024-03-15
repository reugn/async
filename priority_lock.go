package async

import (
	"fmt"
	"sync"
)

const priorityLimit = 1024

// PriorityLock is a non-reentrant mutex that allows specifying a priority
// level when acquiring the lock. It extends the standard sync.Locker interface
// with an additional locking method, LockP, which takes a priority level as an
// argument.
//
// The current implementation may cause starvation for lower priority
// lock requests.
type PriorityLock struct {
	sem []chan struct{}
	max int
}

var _ sync.Locker = (*PriorityLock)(nil)

// NewPriorityLock instantiates and returns a new PriorityLock, specifying the
// maximum priority level that can be used in the LockP method. It panics if
// the maximum priority level is non-positive or exceeds the hard limit.
func NewPriorityLock(maxPriority int) *PriorityLock {
	if maxPriority < 1 {
		panic(fmt.Errorf("nonpositive maximum priority: %d", maxPriority))
	}
	if maxPriority > priorityLimit {
		panic(fmt.Errorf("maximum priority %d exceeds hard limit of %d",
			maxPriority, priorityLimit))
	}
	sem := make([]chan struct{}, maxPriority+1)
	sem[0] = make(chan struct{}, 1)
	sem[0] <- struct{}{}
	for i := 1; i <= maxPriority; i++ {
		sem[i] = make(chan struct{})
	}
	return &PriorityLock{
		sem: sem,
		max: maxPriority,
	}
}

// Lock will block the calling goroutine until it acquires the lock, using
// the highest available priority.
func (pl *PriorityLock) Lock() {
	pl.LockP(pl.max)
}

// LockP blocks the calling goroutine until it acquires the lock. Requests with
// higher priorities acquire the lock first. If the provided priority is
// outside the valid range, it will be assigned the boundary value.
func (pl *PriorityLock) LockP(priority int) {
	switch {
	case priority < 1:
		priority = 1
	case priority > pl.max:
		priority = pl.max
	}
	select {
	case <-pl.sem[priority]:
	case <-pl.sem[0]:
	}
}

// Unlock releases the previously acquired lock.
// It will panic if the lock is already unlocked.
func (pl *PriorityLock) Unlock() {
	for i := pl.max; i >= 0; i-- {
		select {
		case pl.sem[i] <- struct{}{}:
			return
		default:
		}
	}
	panic("async: unlock of unlocked PriorityLock")
}
