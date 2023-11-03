package async

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/reugn/async/internal/assert"
)

func TestWaitGroupContext(t *testing.T) {
	var result atomic.Int32
	wgc := NewWaitGroupContext(context.Background())
	wgc.Add(2)

	go func() {
		defer wgc.Done()
		time.Sleep(10 * time.Millisecond)
		result.Add(1)
	}()
	go func() {
		defer wgc.Done()
		time.Sleep(20 * time.Millisecond)
		result.Add(2)
	}()
	go func() {
		wgc.Wait()
		result.Add(3)
	}()

	wgc.Wait()
	time.Sleep(10 * time.Millisecond)

	assert.Equal(t, int(result.Load()), 6)
}

func TestWaitGroupContextCanceled(t *testing.T) {
	var result atomic.Int32
	ctx, cancelFunc := context.WithCancel(context.Background())
	go func() {
		time.Sleep(100 * time.Millisecond)
		result.Add(10)
		cancelFunc()
	}()
	wgc := NewWaitGroupContext(ctx)
	wgc.Add(2)

	go func() {
		defer wgc.Done()
		time.Sleep(10 * time.Millisecond)
		result.Add(1)
	}()
	go func() {
		defer wgc.Done()
		time.Sleep(300 * time.Millisecond)
		result.Add(2)
	}()
	go func() {
		wgc.Wait()
		result.Add(100)
	}()

	wgc.Wait()
	time.Sleep(10 * time.Millisecond)

	assert.Equal(t, int(result.Load()), 111)
}

func TestWaitGroupContextPanic(t *testing.T) {
	negativeCounter := func() {
		wgc := NewWaitGroupContext(context.Background())
		wgc.Add(-2)
	}
	assert.Panic(t, negativeCounter)
}
