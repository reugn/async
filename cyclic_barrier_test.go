package async

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/reugn/async/internal/assert"
)

func TestCyclicBarrier_Basic(t *testing.T) {
	t.Parallel()

	barrier := NewCyclicBarrier(3)
	var (
		wg    sync.WaitGroup
		count atomic.Int32
	)

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			count.Add(1)
			assert.IsNil(t, barrier.Await())
			// all parties should be released together
			assert.Equal(t, int32(3), count.Load())
			assert.Equal(t, 3, barrier.Parties())
			assert.Equal(t, 0, barrier.Waiting())
		}()
	}

	wg.Wait()
	assert.Equal(t, int32(3), count.Load())
	assert.Equal(t, 3, barrier.Parties())
	assert.Equal(t, 0, barrier.Waiting())
}

func TestCyclicBarrier_AwaitContext(t *testing.T) {
	t.Parallel()

	barrier := NewCyclicBarrier(3)
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			assert.ErrorIs(t, barrier.AwaitContext(ctx), ErrBrokenBarrier)
		}()
	}

	// sleep to ensure the goroutine is waiting
	time.Sleep(10 * time.Millisecond)

	assert.Equal(t, 2, barrier.Waiting())
	// cancel the context to release the goroutine
	cancel()

	wg.Wait()
	assert.Equal(t, 0, barrier.Waiting())
}

func TestCyclicBarrier_ContextCancelRace(t *testing.T) {
	t.Parallel()

	barrier := NewCyclicBarrier(3)
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(3)
		go func() {
			defer wg.Done()
			// first cycle
			assert.ErrorIs(t, barrier.AwaitContext(ctx), ErrBrokenBarrier)
			// second cycle
			assert.IsNil(t, barrier.Await())
		}()
		go func() {
			defer wg.Done()
			// first cycle
			assert.ErrorIs(t, barrier.AwaitContext(ctx), ErrBrokenBarrier)
			// second cycle
			assert.IsNil(t, barrier.Await())
		}()
		go func() {
			defer wg.Done()
			// skip	the first cycle to ensure the barrier is reset
			time.Sleep(10 * time.Millisecond)
			// second cycle
			assert.IsNil(t, barrier.Await())
		}()

		// wait before canceling the context to ensure the first cycle is registered
		time.Sleep(2 * time.Millisecond)
		cancel()
		wg.Wait()
		ctx, cancel = context.WithCancel(context.Background())
	}

	cancel() // cancel the final context created in the loop

	// all goroutines should be released
	assert.Equal(t, 0, barrier.Waiting())
}

func TestCyclicBarrier_SingleParty(t *testing.T) {
	t.Parallel()

	barrier := NewCyclicBarrier(1)
	var (
		wg       sync.WaitGroup
		released atomic.Bool
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		assert.IsNil(t, barrier.Await())
		released.Store(true)
	}()

	wg.Wait()
	assert.Equal(t, true, released.Load())
}

func TestCyclicBarrier_Validation(t *testing.T) {
	t.Parallel()

	assert.PanicMsgContains(t, func() { NewCyclicBarrier(0) }, "async: number of parties must be at least 1")
	assert.PanicMsgContains(t, func() { NewCyclicBarrier(-1) }, "async: number of parties must be at least 1")
}

func TestCyclicBarrier_Reusability(t *testing.T) {
	t.Parallel()

	barrier := NewCyclicBarrier(3)
	var (
		wg    sync.WaitGroup
		phase atomic.Int32
	)

	// phase 1
	wg.Add(3)
	for i := 0; i < 3; i++ {
		go func() {
			defer wg.Done()
			assert.IsNil(t, barrier.Await())
			phase.Add(1)
		}()
	}
	wg.Wait()
	assert.Equal(t, int32(3), phase.Load())

	// phase 2 - reuse the same barrier
	phase.Store(0)
	wg.Add(3)
	for i := 0; i < 3; i++ {
		go func() {
			defer wg.Done()
			assert.IsNil(t, barrier.Await())
			phase.Add(1)
		}()
	}
	wg.Wait()
	assert.Equal(t, int32(3), phase.Load())
}

func TestCyclicBarrier_ConcurrentCycles(t *testing.T) {
	t.Parallel()

	barrier := NewCyclicBarrier(20)
	var (
		wg        sync.WaitGroup
		completed atomic.Int32
	)

	// many goroutines calling Await concurrently
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			assert.IsNil(t, barrier.Await())
			completed.Add(1)
		}()
	}

	wg.Wait()
	// should complete in groups of 20
	assert.Equal(t, int32(100), completed.Load())
}

func TestCyclicBarrier_LastPartyReleasesAll(t *testing.T) {
	t.Parallel()

	barrier := NewCyclicBarrier(4)
	var (
		wg       sync.WaitGroup
		released atomic.Int32
	)

	// first 3 parties arrive quickly
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			assert.IsNil(t, barrier.Await())
			released.Add(1)
		}()
	}

	// last party arrives after delay
	time.Sleep(50 * time.Millisecond)
	wg.Add(1)
	go func() {
		defer wg.Done()
		assert.IsNil(t, barrier.Await())
		released.Add(1)
	}()

	wg.Wait()
	// all 4 should be released
	assert.Equal(t, int32(4), released.Load())
}

func TestCyclicBarrier_MixedTiming(t *testing.T) {
	t.Parallel()

	// different arrival times in milliseconds
	delays := []int{100, 50, 200, 10, 150, 30}
	barrier := NewCyclicBarrier(len(delays))
	var (
		wg    sync.WaitGroup
		count atomic.Int32
	)

	for _, delay := range delays {
		wg.Add(1)
		go func(d int) {
			defer wg.Done()

			time.Sleep(time.Duration(d) * time.Millisecond)
			count.Add(1)
			assert.IsNil(t, barrier.Await())

			// after barrier, all parties should have arrived
			assert.Equal(t, int32(len(delays)), count.Load())
		}(delay)
	}

	wg.Wait()
	// all parties should have arrived
	assert.Equal(t, int32(len(delays)), count.Load())
}

func TestCyclicBarrier_Reset(t *testing.T) {
	t.Parallel()

	barrier := NewCyclicBarrier(5)
	var (
		wg       sync.WaitGroup
		released atomic.Int32
	)

	// start 3 goroutines waiting at the barrier
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			assert.ErrorIs(t, barrier.Await(), ErrBrokenBarrier)
			released.Add(1)
		}()
	}

	// wait a bit to ensure they're waiting
	time.Sleep(50 * time.Millisecond)
	assert.Equal(t, 3, barrier.Waiting())

	// reset should release all waiting goroutines
	barrier.Reset()

	wg.Wait()
	// all 3 should be released
	assert.Equal(t, int32(3), released.Load())
	assert.Equal(t, 0, barrier.Waiting())
}

func TestCyclicBarrier_ResetWhenEmpty(t *testing.T) {
	t.Parallel()

	barrier := NewCyclicBarrier(3)

	// reset when no one is waiting should not panic
	barrier.Reset()
	assert.Equal(t, 0, barrier.Waiting())

	// barrier should still be usable after reset
	var wg sync.WaitGroup
	var released atomic.Int32

	wg.Add(3)
	for i := 0; i < 3; i++ {
		go func() {
			defer wg.Done()
			assert.IsNil(t, barrier.Await())
			released.Add(1)
		}()
	}

	wg.Wait()
	assert.Equal(t, int32(3), released.Load())
}

func TestCyclicBarrier_ResetAndReuse(t *testing.T) {
	t.Parallel()

	barrier := NewCyclicBarrier(4)
	var (
		wg       sync.WaitGroup
		released atomic.Int32
	)

	// first cycle - start 2 goroutines, then reset
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			assert.ErrorIs(t, barrier.Await(), ErrBrokenBarrier)
			released.Add(1)
		}()
	}

	time.Sleep(50 * time.Millisecond)
	barrier.Reset()

	// second cycle - use barrier again after reset
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			assert.IsNil(t, barrier.Await())
			released.Add(1)
		}()
	}

	wg.Wait()
	// 2 from first cycle + 4 from second cycle
	assert.Equal(t, int32(6), released.Load())
	assert.Equal(t, 0, barrier.Waiting())
}

func TestCyclicBarrier_MultipleResets(t *testing.T) {
	t.Parallel()

	barrier := NewCyclicBarrier(5)
	var (
		wg       sync.WaitGroup
		released atomic.Int32
	)

	// multiple reset cycles
	for cycle := 0; cycle < 3; cycle++ {
		// start some goroutines
		for i := 0; i < 3; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				assert.ErrorIs(t, barrier.Await(), ErrBrokenBarrier)
				released.Add(1)
			}()
		}

		time.Sleep(20 * time.Millisecond)
		assert.Equal(t, 3, barrier.Waiting())

		// reset
		barrier.Reset()
		assert.Equal(t, 0, barrier.Waiting())
	}

	wg.Wait()
	// 3 cycles * 3 goroutines
	assert.Equal(t, int32(9), released.Load())
}
