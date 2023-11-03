package async

import (
	"runtime"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/reugn/async/internal/assert"
	"github.com/reugn/async/internal/util"
)

func TestClear(t *testing.T) {
	m := prepareConcurrentMap()
	m.Clear()
	assert.Equal(t, m.Size(), 0)
	m.Put(1, util.Ptr("a"))
	assert.Equal(t, m.Size(), 1)
}

func TestComputeIfAbsent(t *testing.T) {
	m := prepareConcurrentMap()
	assert.Equal(
		t,
		m.ComputeIfAbsent(4, func(_ int) *string { return util.Ptr("d") }),
		util.Ptr("d"),
	)
	assert.Equal(t, m.Size(), 4)
	assert.Equal(
		t,
		m.ComputeIfAbsent(4, func(_ int) *string { return util.Ptr("e") }),
		util.Ptr("d"),
	)
	assert.Equal(t, m.Size(), 4)
}

func TestContainsKey(t *testing.T) {
	m := prepareConcurrentMap()
	assert.Equal(t, m.ContainsKey(3), true)
	assert.Equal(t, m.ContainsKey(4), false)
}

func TestGet(t *testing.T) {
	m := prepareConcurrentMap()
	assert.Equal(t, m.Get(1), util.Ptr("a"))
	assert.Equal(t, m.Get(4), nil)
}

func TestGetOrDefault(t *testing.T) {
	m := prepareConcurrentMap()
	assert.Equal(t, m.GetOrDefault(1, util.Ptr("e")), util.Ptr("a"))
	assert.Equal(t, m.GetOrDefault(5, util.Ptr("e")), util.Ptr("e"))
}

func TestIsEmpty(t *testing.T) {
	m := prepareConcurrentMap()
	assert.Equal(t, m.IsEmpty(), false)
	m.Clear()
	assert.Equal(t, m.IsEmpty(), true)
}

func TestKeySet(t *testing.T) {
	m := prepareConcurrentMap()
	assert.ElementsMatch(t, m.KeySet(), []int{1, 2, 3})
	m.Put(4, util.Ptr("d"))
	assert.ElementsMatch(t, m.KeySet(), []int{1, 2, 3, 4})
}

func TestPut(t *testing.T) {
	m := prepareConcurrentMap()
	assert.Equal(t, m.Size(), 3)
	m.Put(4, util.Ptr("d"))
	assert.Equal(t, m.Size(), 4)
	assert.Equal(t, m.Get(4), util.Ptr("d"))
	m.Put(4, util.Ptr("e"))
	assert.Equal(t, m.Size(), 4)
	assert.Equal(t, m.Get(4), util.Ptr("e"))
}

func TestRemove(t *testing.T) {
	m := prepareConcurrentMap()
	assert.Equal(t, m.Remove(3), util.Ptr("c"))
	assert.Equal(t, m.Size(), 2)
	assert.Equal(t, m.Remove(5), nil)
	assert.Equal(t, m.Size(), 2)
}

func TestSize(t *testing.T) {
	m := prepareConcurrentMap()
	assert.Equal(t, m.Size(), 3)
}

func TestValues(t *testing.T) {
	m := prepareConcurrentMap()
	assert.ElementsMatch(
		t,
		m.Values(),
		[]*string{util.Ptr("a"), util.Ptr("b"), util.Ptr("c")},
	)
	m.Put(4, util.Ptr("d"))
	assert.ElementsMatch(
		t,
		m.Values(),
		[]*string{util.Ptr("a"), util.Ptr("b"), util.Ptr("c"), util.Ptr("d")},
	)
}

func TestMemoryLeaks(t *testing.T) {
	var statsBefore runtime.MemStats
	runtime.ReadMemStats(&statsBefore)

	m := NewConcurrentMap[int, string]()

	var wg sync.WaitGroup
	wg.Add(4)
	go func() {
		defer wg.Done()
		for i := 0; i < 1000000; i++ {
			m.Put(i, util.Ptr(strconv.Itoa(i)))
			time.Sleep(time.Nanosecond)
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			m.Clear()
			time.Sleep(time.Millisecond)
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			m.KeySet()
			time.Sleep(10 * time.Millisecond)
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 80; i++ {
			m.Values()
			time.Sleep(12 * time.Millisecond)
		}
	}()

	wg.Wait()
	m.Clear()
	runtime.GC()

	var statsAfter runtime.MemStats
	runtime.ReadMemStats(&statsAfter)

	assert.Equal(t, m.IsEmpty(), true)
	if statsAfter.HeapObjects > statsBefore.HeapObjects+50 {
		t.Error("HeapObjects leak")
	}
}

func prepareConcurrentMap() *ConcurrentMap[int, string] {
	syncMap := NewConcurrentMap[int, string]()
	syncMap.Put(1, util.Ptr("a"))
	syncMap.Put(2, util.Ptr("b"))
	syncMap.Put(3, util.Ptr("c"))
	return syncMap
}
