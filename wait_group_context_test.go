package async

import (
	"context"
	"testing"
	"time"

	"github.com/reugn/async/internal/assert"
)

func TestWaitGroupContext(t *testing.T) {
	result := 0
	wgc := NewWaitGroupContext(context.Background())
	wgc.Add(2)

	go func() {
		defer wgc.Done()
		time.Sleep(time.Millisecond * 10)
		result++
	}()
	go func() {
		defer wgc.Done()
		time.Sleep(time.Millisecond * 20)
		result += 2
	}()
	go func() {
		wgc.Wait()
		result += 3
	}()

	wgc.Wait()
	time.Sleep(time.Millisecond * 10)

	assert.Equal(t, result, 6)
}

func TestWaitGroupContextCanceled(t *testing.T) {
	result := 0
	ctx, cancelFunc := context.WithCancel(context.Background())
	go func() {
		time.Sleep(time.Millisecond * 100)
		result += 10
		cancelFunc()
	}()
	wgc := NewWaitGroupContext(ctx)
	wgc.Add(2)

	go func() {
		defer wgc.Done()
		time.Sleep(time.Millisecond * 10)
		result++
	}()
	go func() {
		defer wgc.Done()
		time.Sleep(time.Millisecond * 300)
		result += 2
	}()
	go func() {
		wgc.Wait()
		result += 100
	}()

	wgc.Wait()
	time.Sleep(time.Millisecond * 10)

	assert.Equal(t, result, 111)
}

func TestWaitGroupContextPanic(t *testing.T) {
	negativeCounter := func() {
		wgc := NewWaitGroupContext(context.Background())
		wgc.Add(-2)
	}
	assert.Panic(t, negativeCounter)
}
