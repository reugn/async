package async

import (
	"sync"
)

// PriorityLock is a non-reentrant mutex that allows for the specification
// of lock acquisition priority. It extends the standard sync.Locker
// interface with three additional locking methods based on priority levels:
//
// - LockH ensures the highest priority for acquiring the lock, making it
// suitable for critical sections requiring immediate access.
// - LockM has a moderate priority level, intended for non-urgent critical
// sections that can be delayed by locks acquired using LockH.
// - LockL acquires the lock only if no other lock requests with higher
// priority are pending.
//
// The current implementation may cause starvation for lower priority
// lock requests.
type PriorityLock struct {
	high     chan struct{}
	moderate chan struct{}
	low      chan struct{}
	idle     chan struct{}
}

var _ sync.Locker = (*PriorityLock)(nil)

// NewPriorityLock instantiates and returns a new PriorityLock.
func NewPriorityLock() *PriorityLock {
	idle := make(chan struct{}, 1)
	idle <- struct{}{}
	return &PriorityLock{
		high:     make(chan struct{}),
		moderate: make(chan struct{}),
		low:      make(chan struct{}),
		idle:     idle,
	}
}

// Lock will block the calling goroutine until LockH acquires the lock.
func (pl *PriorityLock) Lock() {
	pl.LockH()
}

// LockH ensures the highest priority for acquiring the lock, making it
// suitable for critical sections requiring immediate access.
func (pl *PriorityLock) LockH() {
	select {
	case <-pl.high:
	case <-pl.idle:
	}
}

// LockM has a moderate priority level, intended for non-urgent critical
// sections that can be delayed by locks acquired using LockH.
func (pl *PriorityLock) LockM() {
	select {
	case <-pl.moderate:
	case <-pl.idle:
	}
}

// LockL acquires the lock only if no other lock requests with higher
// priority are pending.
func (pl *PriorityLock) LockL() {
	select {
	case <-pl.low:
	case <-pl.idle:
	}
}

// Unlock releases the previously acquired lock.
// It will panic if the lock is already unlocked.
func (pl *PriorityLock) Unlock() {
	select {
	case pl.high <- struct{}{}:
		return
	default:
	}
	select {
	case pl.moderate <- struct{}{}:
		return
	default:
	}
	select {
	case pl.low <- struct{}{}:
		return
	default:
	}
	select {
	case pl.idle <- struct{}{}:
		return
	default:
		panic("async: unlock of unlocked PriorityLock")
	}
}
