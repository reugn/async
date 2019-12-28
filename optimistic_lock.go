package async

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/reugn/async/internal"
)

type lockStatus uint8

const (
	lockStatusSteady lockStatus = iota
	lockStatusWriting
)

// OptimisticLock allows optimistic reading
// Could be retried or switched to RLock in case of failure
type OptimisticLock struct {
	rw     *sync.RWMutex
	stamp  int64
	status lockStatus
}

// NewOptimisticLock returns new OptimisticLock
func NewOptimisticLock() *OptimisticLock {
	return &OptimisticLock{
		rw:     &sync.RWMutex{},
		stamp:  0,
		status: lockStatusSteady,
	}
}

// Lock locks resource for write
func (o *OptimisticLock) Lock() {
	o.rw.Lock()
	o.status = lockStatusWriting
}

// Unlock unlocks resource after write
func (o *OptimisticLock) Unlock() {
	atomic.StoreInt64(&o.stamp, internal.Cas())
	o.status = lockStatusSteady
	o.rw.Unlock()
}

// RLock locks resource for read
func (o *OptimisticLock) RLock() {
	o.rw.RLock()
}

// RUnlock unlocks resource after read
func (o *OptimisticLock) RUnlock() {
	o.rw.RUnlock()
}

// OptLock returns stamp to be checked on OptUnlock
func (o *OptimisticLock) OptLock() int64 {
	return o.stamp
}

// OptUnlock returns boolean result of optimistic unlock
// Retry or switch to read lock in case of negative outcome
func (o *OptimisticLock) OptUnlock(stamp int64) bool {
	if o.status == lockStatusSteady && stamp == o.stamp {
		return true
	}
	time.Sleep(time.Nanosecond) //switch context
	return false
}
