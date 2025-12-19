package async

import (
	"context"
	"errors"
	"sync"
)

// ErrBrokenBarrier is returned when the cyclic barrier is broken.
// This can happen if the barrier is reset or the await context is done.
var ErrBrokenBarrier = errors.New("async: cyclic barrier is broken")

// CyclicBarrier is a synchronization primitive that allows a group of goroutines
// to wait for each other to reach a common barrier point. It is a reusable barrier
// that can be used multiple times.
type CyclicBarrier struct {
	mu        sync.Mutex
	cycleChan chan error
	count     int // number of waiting parties in current cycle
	parties   int
}

// NewCyclicBarrier creates a new CyclicBarrier with the given number of parties.
// It panics if the number of parties is less than 1.
func NewCyclicBarrier(parties int) *CyclicBarrier {
	if parties < 1 {
		panic("async: number of parties must be at least 1")
	}
	return &CyclicBarrier{
		cycleChan: make(chan error, parties),
		parties:   parties,
	}
}

// Await waits for all parties to reach the barrier. If the current party is the
// last to arrive, it will release all parties and reset the barrier. Otherwise,
// it will wait for the other parties to arrive.
// It returns nil on success, or an error if the barrier is broken.
func (cb *CyclicBarrier) Await() error {
	return cb.AwaitContext(context.Background())
}

// AwaitContext waits for all parties to reach the barrier. If the current party
// is the last to arrive, it will release all parties and reset the barrier.
// Otherwise, it will wait for the other parties to arrive.
// It returns nil on success, or an error if the barrier is broken. The error
// can be ErrBrokenBarrier if the barrier is reset or the context is done.
// It breaks the barrier if the context is done and the barrier is not already broken.
func (cb *CyclicBarrier) AwaitContext(ctx context.Context) error {
	cb.mu.Lock()

	cb.count++
	if cb.count == cb.parties {
		// last party, release all and reset the barrier for the next cycle
		cb.resetBarrier()
		cb.mu.Unlock()
		return nil
	}

	// capture the channel reference before unlocking
	ch := cb.cycleChan
	cb.mu.Unlock()

	// wait for release or context done
	select {
	case err := <-ch:
		return err // barrier was released or broken
	case <-ctx.Done():
		// break the barrier if still active
		cb.mu.Lock()
		// check if still waiting in the same cycle
		if cb.cycleChan == ch {
			cb.breakBarrier()
		}
		cb.mu.Unlock()

		// try to read the result from the channel
		select {
		case err := <-ch:
			return err // barrier was released or broken
		default:
			return errors.Join(ErrBrokenBarrier, ctx.Err())
		}
	}
}

// Waiting returns the number of parties currently waiting at the barrier.
// If the barrier is not in a waiting state, it returns 0.
func (cb *CyclicBarrier) Waiting() int {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	return cb.count
}

// Parties returns the number of parties that must reach the barrier before it
// is released.
func (cb *CyclicBarrier) Parties() int {
	return cb.parties
}

// Reset resets the barrier to its initial state.
// It breaks the barrier by notifying all waiting goroutines that the barrier
// is broken. If there are no waiting goroutines, it does nothing.
func (cb *CyclicBarrier) Reset() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.breakBarrier()
}

// breakBarrier breaks the barrier by notifying all waiting goroutines that the
// barrier is broken and resets the barrier for the next cycle.
func (cb *CyclicBarrier) breakBarrier() {
	// return immediately if there are no waiting goroutines
	if cb.count == 0 {
		return
	}

	// send error to all waiting goroutines
	for i := 0; i < cb.count; i++ {
		select {
		case cb.cycleChan <- ErrBrokenBarrier:
		default:
		}
	}

	// reset the barrier for the next cycle
	cb.resetBarrier()
}

// resetBarrier resets the barrier for the next cycle. It is called when the
// barrier is released and the next cycle is started.
func (cb *CyclicBarrier) resetBarrier() {
	// close the channel to release all waiting goroutines
	close(cb.cycleChan)

	// reset for next cycle
	cb.cycleChan = make(chan error, cb.parties)
	cb.count = 0
}
