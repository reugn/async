package async

import "sync"

// ReentrantLock allows goroutines to enter the lock more than once.
type ReentrantLock struct {
	g           *sync.Mutex
	l           *sync.Mutex
	goroutineID uint
	counter     uint
	lockBalance int
}

// NewReentrantLock returns a new ReentrantLock.
func NewReentrantLock() *ReentrantLock {
	return &ReentrantLock{
		g:           &sync.Mutex{},
		l:           &sync.Mutex{},
		lockBalance: 1,
	}
}

// Lock locks the resource.
// Panics if the GoroutineID call returns an error.
func (r *ReentrantLock) Lock() {
	curr, err := GoroutineID()
	if err != nil {
		panic("async: Error on GoroutineID call")
	}
loop:
	for {
		r.l.Lock()
		switch r.goroutineID {
		case 0:
			// first time lock
			r.handleLock()
			r.goroutineID = curr
			r.counter++
			break loop
		case curr:
			// reentrant lock request
			r.counter++
			break loop
		default:
			// another goroutine lock request
			r.lockBalance--
			r.l.Unlock()
			r.g.Lock()
		}
	}
	r.l.Unlock()
}

func (r *ReentrantLock) handleLock() {
	if r.lockBalance > 0 {
		r.lockBalance--
		r.g.Lock()
	}
}

// Unlock unlocks the resource.
// Panics on trying to unlock the unlocked lock.
func (r *ReentrantLock) Unlock() {
	r.l.Lock()
	defer r.l.Unlock()
	if r.counter == 0 && r.goroutineID == 0 {
		panic("async: Unlock of unlocked ReentrantLock")
	}
	r.counter--
	if r.counter == 0 {
		r.goroutineID = 0
		r.lockBalance++
		r.g.Unlock()
	}
}
