package async

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/reugn/async/internal"
)

const (
	lockStatusSteady int32 = iota
	lockStatusWriting
)

// OptimisticLock allows optimistic reading.
// Implements the Locker interface and is not reentrant.
type OptimisticLock struct {
	rw     *sync.RWMutex
	stamp  int64
	status int32
}

// NewOptimisticLock returns a new OptimisticLock.
func NewOptimisticLock() *OptimisticLock {
	return &OptimisticLock{
		rw:     &sync.RWMutex{},
		stamp:  0,
		status: lockStatusSteady,
	}
}

// Lock locks the resource for write.
func (o *OptimisticLock) Lock() {
	o.rw.Lock()
	atomic.StoreInt32(&o.status, lockStatusWriting)
}

// Unlock unlocks the resource after write.
func (o *OptimisticLock) Unlock() {
	atomic.StoreInt64(&o.stamp, internal.Cas())
	atomic.StoreInt32(&o.status, lockStatusSteady)
	o.rw.Unlock()
}

// RLock locks the resource for read.
func (o *OptimisticLock) RLock() {
	o.rw.RLock()
}

// RUnlock unlocks the resource after read.
func (o *OptimisticLock) RUnlock() {
	o.rw.RUnlock()
}

// OptLock returns the stamp to be verified on OptUnlock.
func (o *OptimisticLock) OptLock() int64 {
	return atomic.LoadInt64(&o.stamp)
}

// OptUnlock returns true if the lock has not been acquired in write mode since obtaining a given stamp.
// Retry or switch to RLock in case of failure.
func (o *OptimisticLock) OptUnlock(stamp int64) bool {
	if atomic.LoadInt32(&o.status) == lockStatusSteady && stamp == atomic.LoadInt64(&o.stamp) {
		return true
	}

	time.Sleep(time.Nanosecond) // switch context
	return false
}
