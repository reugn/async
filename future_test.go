package async

import (
	"errors"
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/reugn/async/internal/assert"
)

func TestFuture(t *testing.T) {
	p := NewPromise[bool]()
	go func() {
		time.Sleep(time.Millisecond * 100)
		p.Success(true)
	}()
	res, err := p.Future().Join()

	assert.Equal(t, res, true)
	assert.Equal(t, err, nil)
}

func TestFutureUtils(t *testing.T) {
	p1 := NewPromise[int]()
	p2 := NewPromise[int]()
	p3 := NewPromise[int]()
	go func() {
		time.Sleep(time.Millisecond * 100)
		p1.Success(1)
		time.Sleep(time.Millisecond * 200)
		p2.Success(2)
		time.Sleep(time.Millisecond * 300)
		p3.Success(3)
	}()
	arr := []Future[int]{p1.Future(), p2.Future(), p3.Future()}
	res := []interface{}{1, 2, 3}
	futRes, _ := FutureSeq(arr).Join()

	assert.Equal(t, res, futRes)
}

func TestFutureFirstCompleted(t *testing.T) {
	p := NewPromise[bool]()
	go func() {
		time.Sleep(time.Millisecond * 1000)
		p.Success(true)
	}()
	timeout := FutureTimer[bool](time.Millisecond * 100)
	futRes, futErr := FutureFirstCompletedOf(p.Future(), timeout).Join()

	assert.Equal(t, false, futRes)
	if futErr == nil {
		t.Fatalf("futErr is nil")
	}
}

func TestFutureTransform(t *testing.T) {
	p1 := NewPromise[int]()
	go func() {
		time.Sleep(time.Millisecond * 100)
		p1.Success(1)
	}()
	future := p1.Future().Map(func(v int) (int, error) {
		return v + 1, nil
	}).FlatMap(func(v int) (Future[int], error) {
		nv := v + 1
		p2 := NewPromise[int]()
		p2.Success(nv)
		return p2.Future(), nil
	}).Recover(func() (int, error) {
		return 5, nil
	})

	res, _ := future.Get(time.Second * 5)
	assert.Equal(t, 3, res)

	res, _ = future.Join()
	assert.Equal(t, 3, res)
}

func TestFutureFailure(t *testing.T) {
	p1 := NewPromise[int]()
	p2 := NewPromise[int]()
	go func() {
		time.Sleep(time.Millisecond * 100)
		p1.Failure(errors.New("Future error"))
		time.Sleep(time.Millisecond * 200)
		p2.Success(2)
	}()
	res, _ := p1.Future().RecoverWith(p2.Future()).Join()

	assert.Equal(t, 2, res)
}

func TestFutureTimeout(t *testing.T) {
	p := NewPromise[bool]()
	go func() {
		time.Sleep(time.Millisecond * 200)
		p.Success(true)
	}()
	future := p.Future()

	_, err := future.Get(time.Millisecond * 50)
	assert.ErrorContains(t, err, "timeout")

	_, err = future.Join()
	assert.ErrorContains(t, err, "timeout")
}

func TestFutureGoroutineLeak(t *testing.T) {
	var wg sync.WaitGroup

	fmt.Println(runtime.NumGoroutine())

	numFuture := 100
	for i := 0; i < numFuture; i++ {
		promise := NewPromise[string]()
		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(time.Millisecond * 100)
			promise.Success("OK")
		}()
		wg.Add(1)
		go func() {
			defer wg.Done()
			fut := promise.Future()
			_, _ = fut.Get(time.Millisecond * 10)
			time.Sleep(time.Millisecond * 100)
			_, _ = fut.Join()
		}()
	}

	wg.Wait()
	time.Sleep(time.Millisecond * 10)
	numGoroutine := runtime.NumGoroutine()
	fmt.Println(numGoroutine)
	if numGoroutine > numFuture {
		t.Fatalf("numGoroutine is %d", numGoroutine)
	}
}
