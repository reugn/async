package async

import (
	"strings"
	"testing"
	"time"

	"github.com/reugn/async/internal/assert"
)

func TestPriorityLock(t *testing.T) {
	p := NewPriorityLock()
	var b strings.Builder

	p.LockH() // acquire first to make the result predictable
	go func() {
		time.Sleep(time.Millisecond)
		p.Unlock()
	}()
	for i := 0; i < 10; i++ {
		go func() {
			p.LockH()
			time.Sleep(time.Microsecond)
			b.WriteRune('h')
			p.Unlock()
		}()
		go func() {
			p.LockM()
			time.Sleep(time.Microsecond)
			b.WriteRune('m')
			p.Unlock()
		}()
		go func() {
			p.LockL()
			time.Sleep(time.Microsecond)
			b.WriteRune('l')
			p.Unlock()
		}()
	}
	time.Sleep(5 * time.Millisecond)
	p.LockL()
	expected := strings.Repeat("h", 10) + strings.Repeat("m", 10) + strings.Repeat("l", 10)
	p.Unlock()
	assert.Equal(t, b.String(), expected)
}

func TestPriorityLock_IdleLock(t *testing.T) {
	p := NewPriorityLock()
	var b strings.Builder
	p.LockH()
	b.WriteRune('h')
	p.Unlock()
	p.LockM()
	b.WriteRune('m')
	p.Unlock()
	p.LockL()
	b.WriteRune('l')
	p.Unlock()
	assert.Equal(t, b.String(), "hml")
}

func TestPriorityLock_Panic(t *testing.T) {
	p := NewPriorityLock()
	p.Lock()
	time.Sleep(time.Nanosecond) // to silence empty critical section warning
	p.Unlock()
	assert.PanicMsgContains(t, func() { p.Unlock() }, "unlock of unlocked PriorityLock")
}
