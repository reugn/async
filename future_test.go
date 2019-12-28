package async

import (
	"testing"
	"time"
)

func TestFuture(t *testing.T) {
	p := NewPromise()
	go func() {
		time.Sleep(time.Millisecond * 100)
		p.Success(true)
	}()
	v, e := p.Future().Get()
	assertEqual(t, v.(bool), true)
	assertEqual(t, e, nil)
}

func TestFutureUtils(t *testing.T) {
	p1 := NewPromise()
	p2 := NewPromise()
	p3 := NewPromise()
	go func() {
		time.Sleep(time.Millisecond * 100)
		p1.Success(1)
		time.Sleep(time.Millisecond * 200)
		p2.Success(2)
		time.Sleep(time.Millisecond * 300)
		p3.Success(3)
	}()
	arr := []Future{p1.Future(), p2.Future(), p3.Future()}
	res := []interface{}{1, 2, 3}
	futRes, _ := FutureSeq(arr).Get()
	assertEqual(t, res, futRes)
}

func TestFutureFirstCompleted(t *testing.T) {
	p := NewPromise()
	go func() {
		time.Sleep(time.Millisecond * 1000)
		p.Success(true)
	}()
	timeout := FutureTimer(time.Millisecond * 100)
	futRes, futErr := FutureFirstCompletedOf(p.Future(), timeout).Get()
	assertEqual(t, nil, futRes)
	if futErr == nil {
		t.Fatalf("futErr is nil")
	}
}
