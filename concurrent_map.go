package async

import (
	"sync"
	"sync/atomic"
	"time"
)

// ConcurrentMap implements the async.Map interface in a thread-safe manner
// by delegating load/store operations to the underlying sync.Map.
//
// The sync.Map type is optimized for two common use cases: (1) when the entry for a given
// key is only ever written once but read many times, as in caches that only grow,
// or (2) when multiple goroutines read, write, and overwrite entries for disjoint
// sets of keys. In these two cases, use of a sync.Map may significantly reduce lock
// contention compared to a Go map paired with a separate sync.Mutex or sync.RWMutex.
type ConcurrentMap[K comparable, V any] struct {
	m        atomic.Value
	size     int64
	clearing int32 // TODO: use atomic.Bool when upgrading to/past go1.19
}

var _ Map[int, any] = (*ConcurrentMap[int, any])(nil)

// NewConcurrentMap returns a new ConcurrentMap instance.
func NewConcurrentMap[K comparable, V any]() *ConcurrentMap[K, V] {
	var underlying atomic.Value
	underlying.Store(&sync.Map{})
	return &ConcurrentMap[K, V]{
		m: underlying,
	}
}

// Clear removes all of the mappings from this map.
func (cm *ConcurrentMap[K, V]) Clear() {
	atomic.StoreInt32(&cm.clearing, 1)
	defer atomic.StoreInt32(&cm.clearing, 0)
	_ = cm.m.Swap(&sync.Map{})
	atomic.StoreInt64(&cm.size, 0)
}

// ComputeIfAbsent attempts to compute a value using the given mapping
// function and enters it into the map, if the specified key is not
// already associated with a value.
func (cm *ConcurrentMap[K, V]) ComputeIfAbsent(key K, mappingFunction func(K) *V) *V {
	value := cm.Get(key)
	if value == nil {
		computed, loaded := cm.smap().LoadOrStore(key, mappingFunction(key))
		if !loaded {
			atomic.AddInt64(&cm.size, 1)
		}
		return computed.(*V)
	}
	return value
}

// ContainsKey returns true if this map contains a mapping for the
// specified key.
func (cm *ConcurrentMap[K, V]) ContainsKey(key K) bool {
	return cm.Get(key) != nil
}

// Get returns the value to which the specified key is mapped, or nil if
// this map contains no mapping for the key.
func (cm *ConcurrentMap[K, V]) Get(key K) *V {
	value, ok := cm.smap().Load(key)
	if !ok {
		return nil
	}
	return value.(*V)
}

// GetOrDefault returns the value to which the specified key is mapped, or
// defaultValue if this map contains no mapping for the key.
func (cm *ConcurrentMap[K, V]) GetOrDefault(key K, defaultValue *V) *V {
	value, ok := cm.smap().Load(key)
	if !ok {
		return defaultValue
	}
	return value.(*V)
}

// IsEmpty returns true if this map contains no key-value mappings.
func (cm *ConcurrentMap[K, V]) IsEmpty() bool {
	return cm.Size() == 0
}

// KeySet returns a slice of the keys contained in this map.
func (cm *ConcurrentMap[K, V]) KeySet() []K {
	keys := make([]K, 0, cm.Size())
	rangeKeysFunc := func(key any, _ any) bool {
		keys = append(keys, key.(K))
		return true
	}
	cm.smap().Range(rangeKeysFunc)
	return keys
}

// Put associates the specified value with the specified key in this map.
func (cm *ConcurrentMap[K, V]) Put(key K, value *V) {
	// TODO: use sync.Map.Swap when upgrading to/past go1.20
	_, loaded := cm.smap().LoadOrStore(key, value)
	if !loaded {
		atomic.AddInt64(&cm.size, 1)
	} else {
		cm.smap().Store(key, value)
	}
}

// Remove removes the mapping for a key from this map if it is present,
// returning the previous value or nil if none.
func (cm *ConcurrentMap[K, V]) Remove(key K) *V {
	value, loaded := cm.smap().LoadAndDelete(key)
	if !loaded {
		return nil
	}
	atomic.AddInt64(&cm.size, -1)
	return value.(*V)
}

// Size returns the number of key-value mappings in this map.
func (cm *ConcurrentMap[K, V]) Size() int {
	size := atomic.LoadInt64(&cm.size)
	if size > 0 {
		return int(size)
	}
	return 0
}

// Values returns a slice of the values contained in this map.
func (cm *ConcurrentMap[K, V]) Values() []*V {
	values := make([]*V, 0, cm.Size())
	rangeValuesFunc := func(_ any, value any) bool {
		values = append(values, value.(*V))
		return true
	}
	cm.smap().Range(rangeValuesFunc)
	return values
}

func (cm *ConcurrentMap[K, V]) smap() *sync.Map {
	for {
		c := atomic.LoadInt32(&cm.clearing)
		if c == 0 {
			break
		}
		time.Sleep(time.Nanosecond)
	}
	return cm.m.Load().(*sync.Map)
}
