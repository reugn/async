package async

import (
	"runtime"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/reugn/async/internal"
)

func TestClear(t *testing.T) {
	m := prepareConcurrentMap()
	m.Clear()
	internal.AssertEqual(t, m.Size(), 0)
	m.Put(1, ptr("a"))
	internal.AssertEqual(t, m.Size(), 1)
}

func TestComputeIfAbsent(t *testing.T) {
	m := prepareConcurrentMap()
	internal.AssertEqual(
		t,
		m.ComputeIfAbsent(4, func(_ int) *string { return ptr("d") }),
		ptr("d"),
	)
	internal.AssertEqual(t, m.Size(), 4)
	internal.AssertEqual(
		t,
		m.ComputeIfAbsent(4, func(_ int) *string { return ptr("e") }),
		ptr("d"),
	)
	internal.AssertEqual(t, m.Size(), 4)
}

func TestContainsKey(t *testing.T) {
	m := prepareConcurrentMap()
	internal.AssertEqual(t, m.ContainsKey(3), true)
	internal.AssertEqual(t, m.ContainsKey(4), false)
}

func TestGet(t *testing.T) {
	m := prepareConcurrentMap()
	internal.AssertEqual(t, m.Get(1), ptr("a"))
	internal.AssertEqual(t, m.Get(4), nil)
}

func TestGetOrDefault(t *testing.T) {
	m := prepareConcurrentMap()
	internal.AssertEqual(t, m.GetOrDefault(1, ptr("e")), ptr("a"))
	internal.AssertEqual(t, m.GetOrDefault(5, ptr("e")), ptr("e"))
}

func TestIsEmpty(t *testing.T) {
	m := prepareConcurrentMap()
	internal.AssertEqual(t, m.IsEmpty(), false)
	m.Clear()
	internal.AssertEqual(t, m.IsEmpty(), true)
}

func TestKeySet(t *testing.T) {
	m := prepareConcurrentMap()
	internal.AssertElementsMatch(t, m.KeySet(), []int{1, 2, 3})
	m.Put(4, ptr("d"))
	internal.AssertElementsMatch(t, m.KeySet(), []int{1, 2, 3, 4})
}

func TestPut(t *testing.T) {
	m := prepareConcurrentMap()
	internal.AssertEqual(t, m.Size(), 3)
	m.Put(4, ptr("d"))
	internal.AssertEqual(t, m.Size(), 4)
	internal.AssertEqual(t, m.Get(4), ptr("d"))
	m.Put(4, ptr("e"))
	internal.AssertEqual(t, m.Size(), 4)
	internal.AssertEqual(t, m.Get(4), ptr("e"))
}

func TestRemove(t *testing.T) {
	m := prepareConcurrentMap()
	internal.AssertEqual(t, m.Remove(3), ptr("c"))
	internal.AssertEqual(t, m.Size(), 2)
	internal.AssertEqual(t, m.Remove(5), nil)
	internal.AssertEqual(t, m.Size(), 2)
}

func TestSize(t *testing.T) {
	m := prepareConcurrentMap()
	internal.AssertEqual(t, m.Size(), 3)
}

func TestValues(t *testing.T) {
	m := prepareConcurrentMap()
	internal.AssertElementsMatch(
		t,
		m.Values(),
		[]*string{ptr("a"), ptr("b"), ptr("c")},
	)
	m.Put(4, ptr("d"))
	internal.AssertElementsMatch(
		t,
		m.Values(),
		[]*string{ptr("a"), ptr("b"), ptr("c"), ptr("d")},
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
			m.Put(i, ptr(strconv.Itoa(i)))
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
			time.Sleep(time.Millisecond * 10)
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 80; i++ {
			m.Values()
			time.Sleep(time.Millisecond * 12)
		}
	}()

	wg.Wait()
	m.Clear()
	runtime.GC()

	var statsAfter runtime.MemStats
	runtime.ReadMemStats(&statsAfter)

	internal.AssertEqual(t, m.IsEmpty(), true)
	if statsAfter.HeapObjects > statsBefore.HeapObjects+50 {
		t.Error("HeapObjects leak")
	}
}

func prepareConcurrentMap() *ConcurrentMap[int, string] {
	syncMap := NewConcurrentMap[int, string]()
	syncMap.Put(1, ptr("a"))
	syncMap.Put(2, ptr("b"))
	syncMap.Put(3, ptr("c"))
	return syncMap
}

func ptr(s string) *string {
	return &s
}
