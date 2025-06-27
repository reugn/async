package async

import (
	"fmt"
	"hash/fnv"
	"sync"
)

// ShardedMap implements the async.Map interface in a thread-safe manner,
// delegating load/store operations to one of the underlying async.SynchronizedMaps
// (shards), using a key hash to calculate the shard number.
// A ShardedMap must not be copied.
type ShardedMap[K comparable, V any] struct {
	shards   uint64
	shardMap []*SynchronizedMap[K, V]
	hashFunc func(K) uint64
}

var _ Map[int, any] = (*ShardedMap[int, any])(nil)

// NewShardedMap returns a new ShardedMap, where shards is the number of partitions for this
// map. It uses the 64-bit FNV-1a hash function to calculate the shard number for a key.
// If the shards argument is not positive, NewShardedMap will panic.
func NewShardedMap[K comparable, V any](shards int) *ShardedMap[K, V] {
	return NewShardedMapWithHash[K, V](shards, func(key K) uint64 {
		h := fnv.New64a()
		_, _ = fmt.Fprint(h, key)
		return h.Sum64()
	})
}

// NewShardedMapWithHash returns a new ShardedMap, where shards is the number of partitions
// for this map, and hashFunc is a hash function to calculate the shard number for a key.
// If shards is not positive or hashFunc is nil, NewShardedMapWithHash will panic.
func NewShardedMapWithHash[K comparable, V any](shards int, hashFunc func(K) uint64) *ShardedMap[K, V] {
	if shards < 1 {
		panic(fmt.Sprintf("nonpositive shards: %d", shards))
	}
	if hashFunc == nil {
		panic("hashFunc is nil")
	}
	shardMap := make([]*SynchronizedMap[K, V], shards)
	for i := 0; i < shards; i++ {
		shardMap[i] = NewSynchronizedMap[K, V]()
	}
	return &ShardedMap[K, V]{
		shards:   uint64(shards),
		shardMap: shardMap,
		hashFunc: hashFunc,
	}
}

// Clear removes all of the mappings from this map.
func (sm *ShardedMap[K, V]) Clear() {
	for _, shard := range sm.shardMap {
		shard.Clear()
	}
}

// ComputeIfAbsent attempts to compute a value using the given mapping
// function and enters it into the map, if the specified key is not
// already associated with a value.
func (sm *ShardedMap[K, V]) ComputeIfAbsent(key K, mappingFunction func(K) *V) *V {
	return sm.shard(key).ComputeIfAbsent(key, mappingFunction)
}

// ContainsKey returns true if this map contains a mapping for the
// specified key.
func (sm *ShardedMap[K, V]) ContainsKey(key K) bool {
	return sm.shard(key).ContainsKey(key)
}

// Get returns the value to which the specified key is mapped, or nil if
// this map contains no mapping for the key.
func (sm *ShardedMap[K, V]) Get(key K) *V {
	return sm.shard(key).Get(key)
}

// GetOrDefault returns the value to which the specified key is mapped, or
// defaultValue if this map contains no mapping for the key.
func (sm *ShardedMap[K, V]) GetOrDefault(key K, defaultValue *V) *V {
	return sm.shard(key).GetOrDefault(key, defaultValue)
}

// IsEmpty returns true if this map contains no key-value mappings.
func (sm *ShardedMap[K, V]) IsEmpty() bool {
	for _, shard := range sm.shardMap {
		if !shard.IsEmpty() {
			return false
		}
	}
	return true
}

// KeySet returns a slice of the keys contained in this map.
func (sm *ShardedMap[K, V]) KeySet() []K {
	var keys []K
	for _, shard := range sm.shardMap {
		keys = append(keys, shard.KeySet()...)
	}
	return keys
}

// Put associates the specified value with the specified key in this map.
func (sm *ShardedMap[K, V]) Put(key K, value *V) {
	sm.shard(key).Put(key, value)
}

// Remove removes the mapping for a key from this map if it is present,
// returning the previous value or nil if none.
func (sm *ShardedMap[K, V]) Remove(key K) *V {
	return sm.shard(key).Remove(key)
}

// Size returns the number of key-value mappings in this map.
func (sm *ShardedMap[K, V]) Size() int {
	var size int
	for _, shard := range sm.shardMap {
		size += shard.Size()
	}
	return size
}

// Values returns a slice of the values contained in this map.
func (sm *ShardedMap[K, V]) Values() []*V {
	var values []*V
	for _, shard := range sm.shardMap {
		values = append(values, shard.Values()...)
	}
	return values
}

// shard returns an underlying synchronized map for the key.
func (sm *ShardedMap[K, V]) shard(key K) Map[K, V] {
	return sm.shardMap[sm.hashFunc(key)%sm.shards]
}

// SynchronizedMap implements the async.Map interface in a thread-safe manner,
// delegating load/store operations to a Go map and using a sync.RWMutex
// for synchronization.
type SynchronizedMap[K comparable, V any] struct {
	sync.RWMutex
	store map[K]*V
}

var _ Map[int, any] = (*SynchronizedMap[int, any])(nil)

// NewSynchronizedMap returns a new SynchronizedMap.
func NewSynchronizedMap[K comparable, V any]() *SynchronizedMap[K, V] {
	return &SynchronizedMap[K, V]{
		store: make(map[K]*V),
	}
}

// Clear removes all of the mappings from this map.
func (sync *SynchronizedMap[K, V]) Clear() {
	sync.Lock()
	defer sync.Unlock()
	sync.store = make(map[K]*V)
}

// ComputeIfAbsent attempts to compute a value using the given mapping
// function and enters it into the map, if the specified key is not
// already associated with a value.
func (sync *SynchronizedMap[K, V]) ComputeIfAbsent(key K, mappingFunction func(K) *V) *V {
	sync.Lock()
	defer sync.Unlock()
	value, ok := sync.store[key]
	if !ok {
		value = mappingFunction(key)
		sync.store[key] = value
	}
	return value
}

// ContainsKey returns true if this map contains a mapping for the
// specified key.
func (sync *SynchronizedMap[K, V]) ContainsKey(key K) bool {
	sync.RLock()
	defer sync.RUnlock()
	_, ok := sync.store[key]
	return ok
}

// Get returns the value to which the specified key is mapped, or nil if
// this map contains no mapping for the key.
func (sync *SynchronizedMap[K, V]) Get(key K) *V {
	sync.RLock()
	defer sync.RUnlock()
	return sync.store[key]
}

// GetOrDefault returns the value to which the specified key is mapped, or
// defaultValue if this map contains no mapping for the key.
func (sync *SynchronizedMap[K, V]) GetOrDefault(key K, defaultValue *V) *V {
	sync.RLock()
	defer sync.RUnlock()
	value, ok := sync.store[key]
	if ok {
		return value
	}
	return defaultValue
}

// IsEmpty returns true if this map contains no key-value mappings.
func (sync *SynchronizedMap[K, V]) IsEmpty() bool {
	sync.RLock()
	defer sync.RUnlock()
	return len(sync.store) == 0
}

// KeySet returns a slice of the keys contained in this map.
func (sync *SynchronizedMap[K, V]) KeySet() []K {
	sync.RLock()
	defer sync.RUnlock()
	keys := make([]K, 0, len(sync.store))
	for key := range sync.store {
		keys = append(keys, key)
	}
	return keys
}

// Put associates the specified value with the specified key in this map.
func (sync *SynchronizedMap[K, V]) Put(key K, value *V) {
	sync.Lock()
	defer sync.Unlock()
	sync.store[key] = value
}

// Remove removes the mapping for a key from this map if it is present,
// returning the previous value or nil if none.
func (sync *SynchronizedMap[K, V]) Remove(key K) *V {
	sync.Lock()
	defer sync.Unlock()
	value, ok := sync.store[key]
	if ok {
		delete(sync.store, key)
	}
	return value
}

// Size returns the number of key-value mappings in this map.
func (sync *SynchronizedMap[K, V]) Size() int {
	sync.RLock()
	defer sync.RUnlock()
	return len(sync.store)
}

// Values returns a slice of the values contained in this map.
func (sync *SynchronizedMap[K, V]) Values() []*V {
	sync.RLock()
	defer sync.RUnlock()
	values := make([]*V, 0, len(sync.store))
	for _, value := range sync.store {
		values = append(values, value)
	}
	return values
}
