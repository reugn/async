package async

import (
	"sync"
	"testing"
)

type synchronizedAdder struct {
	value int
	lock  *ReentrantLock
}

func newSynchronizedAdder() *synchronizedAdder {
	return &synchronizedAdder{
		lock: NewReentrantLock(),
	}
}

func (sa *synchronizedAdder) Value() int {
	return sa.value
}

func (sa *synchronizedAdder) addOne() {
	sa.lock.Lock()
	defer sa.lock.Unlock()
	sa.value++
}

func (sa *synchronizedAdder) addTwo() {
	sa.lock.Lock()
	defer sa.lock.Unlock()
	sa.value += 2
}

func (sa *synchronizedAdder) addThree() {
	sa.lock.Lock()
	defer sa.lock.Unlock()
	sa.value += 3
}

func (sa *synchronizedAdder) addFour() {
	sa.lock.Lock()
	defer sa.lock.Unlock()
	sa.value += 4
	sa.addThree()
}

func (sa *synchronizedAdder) addFive() {
	sa.lock.Lock()
	defer sa.lock.Unlock()
	sa.value += 5
	sa.addFour()
}

func TestPacker1(t *testing.T) {
	adder := newSynchronizedAdder()
	var wg sync.WaitGroup
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go func() {
			adder.addOne()
			adder.addTwo()
			adder.addThree()
			wg.Done()
		}()
	}
	wg.Wait()
	assertEqual(t, 30, adder.Value())
}

func TestPacker2(t *testing.T) {
	adder := newSynchronizedAdder()
	var wg sync.WaitGroup
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go func() {
			adder.addFive()
			wg.Done()
		}()
	}
	wg.Wait()
	assertEqual(t, 60, adder.Value())
}
