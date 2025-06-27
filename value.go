package async

import (
	"sync/atomic"
)

// A Value provides an atomic load and store of a value of any type.
// The behavior is analogous to atomic.Value, except that
// the value is not required to be of the same specific type.
// Can be useful for storing different implementations of an interface.
type Value struct {
	p atomic.Pointer[atomic.Value]
}

// CompareAndSwap executes the compare-and-swap operation for the Value.
// The current implementation is not atomic.
//
//nolint:revive
func (v *Value) CompareAndSwap(old, new any) (swapped bool) {
	defer func() {
		if err := recover(); err != nil {
			swapped = false
		}
	}()
	delegate := v.p.Load()
	if delegate != nil {
		if old == delegate.Load() {
			v.p.Store(initValue(new))
			return true
		}
	}
	return false
}

// Load returns the value set by the most recent Store.
// It returns nil if there has been no call to Store for this Value.
func (v *Value) Load() (val any) {
	delegate := v.p.Load()
	if delegate != nil {
		return delegate.Load()
	}
	return nil
}

// Store sets the value of the Value v to val.
// Store(nil) panics.
func (v *Value) Store(val any) {
	v.p.Store(initValue(val))
}

// Swap stores new into Value and returns the previous value.
// It returns nil if the Value is empty.
// Swap(nil) panics.
//
//nolint:revive
func (v *Value) Swap(new any) (old any) {
	oldValue := v.p.Swap(initValue(new))
	if oldValue != nil {
		return oldValue.Load()
	}
	return nil
}

func initValue(val any) *atomic.Value {
	value := atomic.Value{}
	value.Store(val)
	return &value
}
