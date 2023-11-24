package async

import (
	"context"
	"sync/atomic"
)

// A WaitGroupContext waits for a collection of goroutines to finish.
// The main goroutine calls Add to set the number of goroutines to wait for.
// Then each of the goroutines runs and calls Done when finished. At the same
// time, Wait can be used to block until all goroutines have finished or the
// given context is done.
type WaitGroupContext struct {
	ctx   context.Context
	sem   chan struct{}
	state atomic.Uint64 // high 32 bits are counter, low 32 bits are waiter count.
}

// NewWaitGroupContext returns a new WaitGroupContext with Context ctx.
func NewWaitGroupContext(ctx context.Context) *WaitGroupContext {
	return &WaitGroupContext{
		ctx: ctx,
		sem: make(chan struct{}),
	}
}

// Add adds delta, which may be negative, to the WaitGroupContext counter.
// If the counter becomes zero, all goroutines blocked on Wait are released.
// If the counter goes negative, Add panics.
func (wgc *WaitGroupContext) Add(delta int) {
	state := wgc.state.Add(uint64(delta) << 32)
	counter := int32(state >> 32)
	if counter == 0 {
		wgc.notifyAll()
	} else if counter < 0 {
		panic("async: negative WaitGroupContext counter")
	}
}

// Done decrements the WaitGroupContext counter by one.
func (wgc *WaitGroupContext) Done() {
	wgc.Add(-1)
}

// Wait blocks until the wait group counter is zero or ctx is done.
func (wgc *WaitGroupContext) Wait() {
	for {
		state := wgc.state.Load()
		counter := int32(state >> 32)
		if counter == 0 {
			return
		}
		if wgc.state.CompareAndSwap(state, state+1) {
			select {
			case <-wgc.sem:
				if wgc.state.Load() != 0 {
					panic("async: WaitGroupContext is reused before " +
						"previous Wait has returned")
				}
			case <-wgc.ctx.Done():
			}
			return
		}
	}
}

// notifyAll releases all goroutines blocked in Wait and resets
// the wait group state.
func (wgc *WaitGroupContext) notifyAll() {
	state := wgc.state.Load()
	waiting := uint32(state)
	wgc.state.Store(0)
	for ; waiting != 0; waiting-- {
		select {
		case wgc.sem <- struct{}{}:
		case <-wgc.ctx.Done():
		}
	}
}
