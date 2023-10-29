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
	ctx     context.Context
	done    chan struct{}
	counter atomic.Int32
	state   atomic.Int32
}

// NewWaitGroupContext returns a new WaitGroupContext with Context ctx.
func NewWaitGroupContext(ctx context.Context) *WaitGroupContext {
	return &WaitGroupContext{
		ctx:  ctx,
		done: make(chan struct{}),
	}
}

// Add adds delta, which may be negative, to the WaitGroupContext counter.
// If the counter becomes zero, all goroutines blocked on Wait are released.
// If the counter goes negative, Add panics.
func (wgc *WaitGroupContext) Add(delta int) {
	counter := wgc.counter.Add(int32(delta))
	if counter == 0 && wgc.state.CompareAndSwap(0, 1) {
		wgc.release()
	} else if counter < 0 && wgc.state.Load() == 0 {
		panic("async: negative WaitGroupContext counter")
	}
}

// Done decrements the WaitGroupContext counter by one.
func (wgc *WaitGroupContext) Done() {
	wgc.Add(-1)
}

// Wait blocks until the wait group counter is zero or ctx is done.
func (wgc *WaitGroupContext) Wait() {
	select {
	case <-wgc.ctx.Done():
	case <-wgc.done:
	}
}

func (wgc *WaitGroupContext) release() {
	close(wgc.done)
}
