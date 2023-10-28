package async

import (
	"sync"
	"testing"

	"github.com/reugn/async/internal/assert"
)

type synchronizedAdder struct {
	ReentrantLock
	value int
}

func (sa *synchronizedAdder) getValue() int {
	sa.Lock()
	defer sa.Unlock()
	return sa.value
}

func (sa *synchronizedAdder) addOne() {
	sa.Lock()
	defer sa.Unlock()
	sa.value++
}

func (sa *synchronizedAdder) addTwo() {
	sa.Lock()
	defer sa.Unlock()
	sa.value++
	sa.addOne()
}

func (sa *synchronizedAdder) addThree() {
	sa.Lock()
	defer sa.Unlock()
	sa.value++
	sa.addTwo()
}

func (sa *synchronizedAdder) addFour() {
	sa.Lock()
	defer sa.Unlock()
	sa.value++
	sa.addThree()
}

func (sa *synchronizedAdder) addFive() {
	sa.Lock()
	defer sa.Unlock()
	sa.value++
	sa.addFour()
}

func Test_synchronizedAdder1(t *testing.T) {
	adder := &synchronizedAdder{}
	var wg sync.WaitGroup
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go func() {
			defer wg.Done()
			adder.addOne()
			adder.addTwo()
			adder.addThree()
			adder.addFour()
		}()
	}
	wg.Wait()
	assert.Equal(t, 50, adder.getValue())
}

func Test_synchronizedAdder2(t *testing.T) {
	adder := &synchronizedAdder{}
	var wg sync.WaitGroup
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go func() {
			defer wg.Done()
			adder.addFive()
		}()
	}
	wg.Wait()
	assert.Equal(t, 25, adder.getValue())
}
