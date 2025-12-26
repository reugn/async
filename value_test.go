package async

import (
	"testing"

	"github.com/reugn/async/internal/assert"
	"github.com/reugn/async/internal/ptr"
)

func TestValueCompareAndSwap(t *testing.T) {
	var value Value
	swapped := value.CompareAndSwap(1, 2)
	assert.Equal(t, swapped, false)
	assert.IsNil(t, value.Load())

	swapped = value.CompareAndSwap(1, nil)
	assert.Equal(t, swapped, false)
	assert.IsNil(t, value.Load())

	swapped = value.CompareAndSwap(nil, 1)
	assert.Equal(t, swapped, false)
	assert.IsNil(t, value.Load())

	value.Store(1)

	swapped = value.CompareAndSwap("a", nil)
	assert.Equal(t, swapped, false)
	assert.Equal(t, value.Load(), 1)

	swapped = value.CompareAndSwap(nil, nil)
	assert.Equal(t, swapped, false)
	assert.Equal(t, value.Load(), 1)

	swapped = value.CompareAndSwap("a", 2)
	assert.Equal(t, swapped, false)
	assert.Equal(t, value.Load(), 1)

	swapped = value.CompareAndSwap(-1, 2)
	assert.Equal(t, swapped, false)
	assert.Equal(t, value.Load(), 1)

	swapped = value.CompareAndSwap(1, 2)
	assert.Equal(t, swapped, true)
	assert.Equal(t, value.Load(), 2)

	swapped = value.CompareAndSwap(2, "a")
	assert.Equal(t, swapped, true)
	assert.Equal(t, value.Load(), "a")

	stringPointer := ptr.Of("b")
	swapped = value.CompareAndSwap("a", stringPointer)
	assert.Equal(t, swapped, true)
	assert.Same(t, value.Load().(*string), stringPointer)

	swapped = value.CompareAndSwap(ptr.Of("b"), "c")
	assert.Equal(t, swapped, false)
	assert.Same(t, value.Load().(*string), stringPointer)

	swapped = value.CompareAndSwap(stringPointer, "c")
	assert.Equal(t, swapped, true)
	assert.Equal(t, value.Load(), "c")
}

func TestValueLoad(t *testing.T) {
	var value Value
	assert.IsNil(t, value.Load())

	value.Store(1)
	assert.Equal(t, value.Load(), 1)
}

func TestValueStore(t *testing.T) {
	var value Value
	value.Store(1)
	assert.Equal(t, value.Load(), 1)

	assert.Panics(t, func() { value.Store(nil) })

	value.Store("a")
	assert.Equal(t, value.Load(), "a")

	stringPointer := ptr.Of("b")
	value.Store(stringPointer)
	assert.Same(t, value.Load().(*string), stringPointer)
}

func TestValueSwap(t *testing.T) {
	var value Value
	old := value.Swap(1)
	assert.IsNil(t, old)

	assert.Panics(t, func() { _ = value.Swap(nil) })

	old = value.Swap("a")
	assert.Equal(t, old, 1)

	stringPointer := ptr.Of("b")
	old = value.Swap(stringPointer)
	assert.Equal(t, old, "a")
	assert.Same(t, value.Load().(*string), stringPointer)
}
