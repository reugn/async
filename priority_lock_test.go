package async

import (
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/reugn/async/internal/assert"
)

func TestPriorityLock(t *testing.T) {
	p := NewPriorityLock(5)
	var b strings.Builder

	p.Lock() // acquire first to make the result predictable
	go func() {
		time.Sleep(time.Millisecond)
		p.Unlock()
	}()
	for i := 0; i < 10; i++ {
		for j := 5; j > 0; j-- {
			go func(n int) {
				p.LockP(n)
				time.Sleep(time.Microsecond)
				b.WriteString(strconv.Itoa(n))
				p.Unlock()
			}(j)
		}
	}
	time.Sleep(20 * time.Millisecond)

	p.Lock()
	result := b.String()
	p.Unlock()
	var expected strings.Builder
	for i := 5; i > 0; i-- {
		expected.WriteString(strings.Repeat(strconv.Itoa(i), 10))
	}
	assert.Equal(t, result, expected.String())
}

func TestPriorityLock_LockRange(t *testing.T) {
	p := NewPriorityLock(2)
	var b strings.Builder
	p.LockP(-1)
	b.WriteRune('1')
	p.Unlock()
	p.LockP(2048)
	b.WriteRune('1')
	p.Unlock()
	assert.Equal(t, b.String(), "11")
}

func TestPriorityLock_Panic(t *testing.T) {
	p := NewPriorityLock(2)
	p.Lock()
	time.Sleep(time.Nanosecond) // to silence empty critical section warning
	p.Unlock()
	assert.PanicMsgContains(t, func() { p.Unlock() }, "unlock of unlocked PriorityLock")
}

func TestPriorityLock_Validation(t *testing.T) {
	assert.PanicMsgContains(t, func() { NewPriorityLock(-1) }, "nonpositive maximum priority")
	assert.PanicMsgContains(t, func() { NewPriorityLock(2048) }, "exceeds hard limit")
}
