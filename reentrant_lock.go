package async

import (
	"sync"
)

// ReentrantLock allows goroutines to enter the lock more than once.
// Implements the sync.Locker interface.
//
// A ReentrantLock must not be copied after first use.
type ReentrantLock struct {
	outer     sync.Mutex
	inner     sync.Mutex
	goroutine uint64
	depth     int32
}

var _ sync.Locker = (*ReentrantLock)(nil)

// Lock locks the resource.
// Panics if the GoroutineID call returns an error.
func (r *ReentrantLock) Lock() {
	r.inner.Lock()

	current, err := GoroutineID()
	if err != nil {
		panic("async: Error on GoroutineID call")
	}

	switch r.goroutine {
	case current:
		// reentrant lock request
		r.depth++
		r.inner.Unlock()
	default:
		// initial or another goroutine lock request
		r.init(current)
	}
}

func (r *ReentrantLock) init(goroutine uint64) {
	r.inner.Unlock()
	r.outer.Lock()
	r.inner.Lock()
	r.goroutine = goroutine
	r.depth = 1
	r.inner.Unlock()
}

// Unlock unlocks the resource.
// Panics on trying to unlock the unlocked lock.
func (r *ReentrantLock) Unlock() {
	r.inner.Lock()
	defer r.inner.Unlock()

	r.depth--
	if r.depth == 0 {
		r.goroutine = 0
		r.outer.Unlock()
	}
}
