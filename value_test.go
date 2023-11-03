package async

import (
	"testing"

	"github.com/reugn/async/internal/assert"
	"github.com/reugn/async/internal/util"
)

//nolint:funlen
func TestValueCompareAndSwap(t *testing.T) {
	var value Value
	swapped := value.CompareAndSwap(1, 2)
	assert.Equal(t, swapped, false)
	assert.Equal(t, value.Load(), nil)

	swapped = value.CompareAndSwap(1, nil)
	assert.Equal(t, swapped, false)
	assert.Equal(t, value.Load(), nil)

	swapped = value.CompareAndSwap(nil, 1)
	assert.Equal(t, swapped, false)
	assert.Equal(t, value.Load(), nil)

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

	stringPointer := util.Ptr("b")
	swapped = value.CompareAndSwap("a", stringPointer)
	assert.Equal(t, swapped, true)
	if value.Load() != stringPointer {
		t.Fail()
	}

	swapped = value.CompareAndSwap(util.Ptr("b"), "c")
	assert.Equal(t, swapped, false)
	if value.Load() != stringPointer {
		t.Fail()
	}

	swapped = value.CompareAndSwap(stringPointer, "c")
	assert.Equal(t, swapped, true)
	assert.Equal(t, value.Load(), "c")
}

func TestValueLoad(t *testing.T) {
	var value Value
	assert.Equal(t, value.Load(), nil)

	value.Store(1)
	assert.Equal(t, value.Load(), 1)
}

func TestValueStore(t *testing.T) {
	var value Value
	value.Store(1)
	assert.Equal(t, value.Load(), 1)

	assert.Panic(t, func() { value.Store(nil) })

	value.Store("a")
	assert.Equal(t, value.Load(), "a")

	stringPointer := util.Ptr("b")
	value.Store(stringPointer)
	if value.Load() != stringPointer {
		t.Fail()
	}
}

func TestValueSwap(t *testing.T) {
	var value Value
	old := value.Swap(1)
	assert.Equal(t, old, nil)

	assert.Panic(t, func() { _ = value.Swap(nil) })

	old = value.Swap("a")
	assert.Equal(t, old, 1)

	stringPointer := util.Ptr("b")
	old = value.Swap(stringPointer)
	assert.Equal(t, old, "a")
	if value.Load() != stringPointer {
		t.Fail()
	}
}
