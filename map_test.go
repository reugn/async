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

func TestMap_Clear(t *testing.T) {
	tests := prepareTestMaps()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.Clear()
			assert.Equal(t, tt.m.Size(), 0)
			tt.m.Put(1, util.Ptr("a"))
			assert.Equal(t, tt.m.Size(), 1)
		})
	}
}

func TestMap_ComputeIfAbsent(t *testing.T) {
	tests := prepareTestMaps()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(
				t,
				tt.m.ComputeIfAbsent(4, func(_ int) *string { return util.Ptr("d") }),
				util.Ptr("d"),
			)
			assert.Equal(t, tt.m.Size(), 4)
			assert.Equal(
				t,
				tt.m.ComputeIfAbsent(4, func(_ int) *string { return util.Ptr("e") }),
				util.Ptr("d"),
			)
			assert.Equal(t, tt.m.Size(), 4)
		})
	}
}

func TestMap_ContainsKey(t *testing.T) {
	tests := prepareTestMaps()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.m.ContainsKey(3), true)
			assert.Equal(t, tt.m.ContainsKey(4), false)
		})
	}
}

func TestMap_Get(t *testing.T) {
	tests := prepareTestMaps()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.m.Get(1), util.Ptr("a"))
			assert.IsNil(t, tt.m.Get(4))
		})
	}
}

func TestMap_GetOrDefault(t *testing.T) {
	tests := prepareTestMaps()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.m.GetOrDefault(1, util.Ptr("e")), util.Ptr("a"))
			assert.Equal(t, tt.m.GetOrDefault(5, util.Ptr("e")), util.Ptr("e"))
		})
	}
}

func TestMap_IsEmpty(t *testing.T) {
	tests := prepareTestMaps()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.m.IsEmpty(), false)
			tt.m.Clear()
			assert.Equal(t, tt.m.IsEmpty(), true)
		})
	}
}

func TestMap_KeySet(t *testing.T) {
	tests := prepareTestMaps()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.ElementsMatch(t, tt.m.KeySet(), []int{1, 2, 3})
			tt.m.Put(4, util.Ptr("d"))
			assert.ElementsMatch(t, tt.m.KeySet(), []int{1, 2, 3, 4})
		})
	}
}

func TestMap_Put(t *testing.T) {
	tests := prepareTestMaps()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.m.Size(), 3)
			tt.m.Put(4, util.Ptr("d"))
			assert.Equal(t, tt.m.Size(), 4)
			assert.Equal(t, tt.m.Get(4), util.Ptr("d"))
			tt.m.Put(4, util.Ptr("e"))
			assert.Equal(t, tt.m.Size(), 4)
			assert.Equal(t, tt.m.Get(4), util.Ptr("e"))
		})
	}
}

func TestMap_Remove(t *testing.T) {
	tests := prepareTestMaps()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.m.Remove(3), util.Ptr("c"))
			assert.Equal(t, tt.m.Size(), 2)
			assert.IsNil(t, tt.m.Remove(5))
			assert.Equal(t, tt.m.Size(), 2)
		})
	}
}

func TestMap_Size(t *testing.T) {
	tests := prepareTestMaps()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.m.Size(), 3)
		})
	}
}

func TestMap_Values(t *testing.T) {
	tests := prepareTestMaps()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.ElementsMatch(
				t,
				tt.m.Values(),
				[]*string{util.Ptr("a"), util.Ptr("b"), util.Ptr("c")},
			)
			tt.m.Put(4, util.Ptr("d"))
			assert.ElementsMatch(
				t,
				tt.m.Values(),
				[]*string{util.Ptr("a"), util.Ptr("b"), util.Ptr("c"), util.Ptr("d")},
			)
		})
	}
}

func TestShardedMap_ConstructorArguments(t *testing.T) {
	assert.PanicMsgContains(t, func() {
		NewShardedMap[int, string](0)
	}, "nonpositive shards")

	assert.PanicMsgContains(t, func() {
		NewShardedMapWithHash[int, string](0, func(_ int) uint64 { return 1 })
	}, "nonpositive shards")

	assert.PanicMsgContains(t, func() {
		NewShardedMapWithHash[int, string](2, nil)
	}, "hashFunc is nil")

	NewShardedMapWithHash[int, string](2, func(_ int) uint64 { return 1 })
}

func TestConcurrentMap_MemoryLeaks(t *testing.T) {
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

func prepareTestMaps() []testMap {
	tests := make([]testMap, 0, 2)
	concurrentMap := NewConcurrentMap[int, string]()
	putValues(concurrentMap)
	tests = append(tests, testMap{"concurrentMap", concurrentMap})
	shardedMap := NewShardedMap[int, string](2)
	putValues(shardedMap)
	tests = append(tests, testMap{"shardedMap", shardedMap})
	return tests
}

func putValues(m Map[int, string]) {
	m.Put(1, util.Ptr("a"))
	m.Put(2, util.Ptr("b"))
	m.Put(3, util.Ptr("c"))
}

type testMap struct {
	name string
	m    Map[int, string]
}
