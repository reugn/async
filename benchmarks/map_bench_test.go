package benchmarks

import (
	"math"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/reugn/async"
	"github.com/reugn/async/internal/ptr"
)

var (
	iter   = 5
	shards = 64
)

// go test -bench=. -benchmem -v
func benchmarkMixedConcurrentLoad(m async.Map[mkey, int]) {
	var wg sync.WaitGroup
	for r := 0; r < iter; r++ {
		from := r * iter
		to := (r + 1) * iter
		wg.Add(4)
		go func() {
			defer wg.Done()
			for i := from; i < to; i++ {
				m.Put(mkey{uint64(i), strconv.Itoa(i)}, ptr.Of(i))
			}
		}()
		go func() {
			defer wg.Done()
			for i := from; i < to; i++ {
				_ = m.ComputeIfAbsent(mkey{uint64(i), strconv.Itoa(i)}, func(_ mkey) *int {
					time.Sleep(time.Nanosecond)
					return ptr.Of(i)
				})
			}
		}()
		go func() {
			defer wg.Done()
			for i := from; i < to; i++ {
				_ = m.Get(mkey{uint64(i), strconv.Itoa(i)})
			}
		}()
		go func() {
			defer wg.Done()
			for i := from; i < to; i++ {
				_ = m.GetOrDefault(mkey{uint64(i), strconv.Itoa(i)}, ptr.Of(i))
			}
		}()
	}
	wg.Wait()
}

func BenchmarkMapMixedLoad_ConcurrentMap(b *testing.B) {
	m := async.NewConcurrentMap[mkey, int]()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		benchmarkMixedConcurrentLoad(m)
	}
}

func BenchmarkMapMixedLoad_ShardedMap(b *testing.B) {
	m := async.NewShardedMap[mkey, int](shards)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		benchmarkMixedConcurrentLoad(m)
	}
}

func BenchmarkMapMixedLoad_ShardedMapWithHash(b *testing.B) {
	m := async.NewShardedMapWithHash[mkey, int](
		shards,
		func(k mkey) uint64 { return k.i },
	)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		benchmarkMixedConcurrentLoad(m)
	}
}

func BenchmarkMapMixedLoad_SynchronizedMap(b *testing.B) {
	m := async.NewSynchronizedMap[mkey, int]()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		benchmarkMixedConcurrentLoad(m)
	}
}

func benchmarkReadConcurrentLoad(m async.Map[mkey, int]) {
	var wg sync.WaitGroup
	for r := 0; r < iter; r++ {
		from := r * iter
		to := (r + 1) * iter
		wg.Add(5)
		go func() {
			defer wg.Done()
			for i := from; i < to; i++ {
				_ = m.Get(mkey{uint64(i), strconv.Itoa(i)})
				_ = m.Get(mkey{math.MaxUint64, strconv.Itoa(i)})
			}
		}()
		go func() {
			defer wg.Done()
			for i := from; i < to; i++ {
				_ = m.GetOrDefault(mkey{uint64(i), strconv.Itoa(i)}, ptr.Of(i))
				_ = m.GetOrDefault(mkey{math.MaxUint64, strconv.Itoa(i)}, ptr.Of(i))
			}
		}()
		go func() {
			defer wg.Done()
			for i := from; i < to; i++ {
				_ = m.ContainsKey(mkey{uint64(i), strconv.Itoa(i)})
				_ = m.ContainsKey(mkey{math.MaxUint64, strconv.Itoa(i)})
			}
		}()
		go func() {
			defer wg.Done()
			for i := from; i < to; i++ {
				_ = m.IsEmpty()
			}
		}()
		go func() {
			defer wg.Done()
			for i := from; i < to; i++ {
				_ = m.Size()
			}
		}()
	}
	wg.Wait()
}

func fillMap(m async.Map[mkey, int]) {
	for i := 0; i < 100; i++ {
		m.Put(mkey{uint64(i), strconv.Itoa(i)}, ptr.Of(i))
	}
}

func BenchmarkMapReadLoad_ConcurrentMap(b *testing.B) {
	m := async.NewConcurrentMap[mkey, int]()
	fillMap(m)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		benchmarkReadConcurrentLoad(m)
	}
}

func BenchmarkMapReadLoad_ShardedMap(b *testing.B) {
	m := async.NewShardedMap[mkey, int](shards)
	fillMap(m)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		benchmarkReadConcurrentLoad(m)
	}
}

func BenchmarkMapReadLoad_ShardedMapWithHash(b *testing.B) {
	m := async.NewShardedMapWithHash[mkey, int](
		shards,
		func(k mkey) uint64 { return k.i },
	)
	fillMap(m)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		benchmarkReadConcurrentLoad(m)
	}
}

func BenchmarkMapReadLoad_SynchronizedMap(b *testing.B) {
	m := async.NewSynchronizedMap[mkey, int]()
	fillMap(m)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		benchmarkReadConcurrentLoad(m)
	}
}

func benchmarkWriteConcurrentLoad(m async.Map[mkey, int]) {
	var wg sync.WaitGroup
	for r := 0; r < iter; r++ {
		from := r * iter
		to := (r + 1) * iter
		wg.Add(3)
		go func() {
			defer wg.Done()
			for i := from; i < to; i++ {
				m.Put(mkey{uint64(i), strconv.Itoa(i)}, ptr.Of(i))
			}
		}()
		go func() {
			defer wg.Done()
			for i := from; i < to; i++ {
				_ = m.ComputeIfAbsent(mkey{uint64(i), strconv.Itoa(i)}, func(_ mkey) *int {
					time.Sleep(time.Nanosecond)
					return ptr.Of(i)
				})
			}
		}()
		go func() {
			defer wg.Done()
			for i := to; i >= from; i-- {
				_ = m.Remove(mkey{uint64(i), strconv.Itoa(i)})
			}
		}()
	}
	wg.Wait()
}

func BenchmarkMapWriteLoad_ConcurrentMap(b *testing.B) {
	m := async.NewConcurrentMap[mkey, int]()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		benchmarkWriteConcurrentLoad(m)
	}
}

func BenchmarkMapWriteLoad_ShardedMap(b *testing.B) {
	m := async.NewShardedMap[mkey, int](shards)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		benchmarkWriteConcurrentLoad(m)
	}
}

func BenchmarkMapWriteLoad_ShardedMapWithHash(b *testing.B) {
	m := async.NewShardedMapWithHash[mkey, int](
		shards,
		func(k mkey) uint64 { return k.i },
	)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		benchmarkWriteConcurrentLoad(m)
	}
}

func BenchmarkMapWriteLoad_SynchronizedMap(b *testing.B) {
	m := async.NewSynchronizedMap[mkey, int]()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		benchmarkWriteConcurrentLoad(m)
	}
}

type mkey struct {
	i uint64
	s string
}
