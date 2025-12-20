package async

import "iter"

// A Map is an object that maps keys to values.
type Map[K comparable, V any] interface {

	// Clear removes all of the mappings from this map.
	Clear()

	// ComputeIfAbsent attempts to compute a value using the given mapping
	// function and enters it into the map, if the specified key is not
	// already associated with a value.
	ComputeIfAbsent(key K, mappingFunction func(K) *V) *V

	// ContainsKey returns true if this map contains a mapping for the
	// specified key.
	ContainsKey(key K) bool

	// Get returns the value to which the specified key is mapped, or nil if
	// this map contains no mapping for the key.
	Get(key K) *V

	// GetOrDefault returns the value to which the specified key is mapped, or
	// defaultValue if this map contains no mapping for the key.
	GetOrDefault(key K, defaultValue *V) *V

	// IsEmpty returns true if this map contains no key-value mappings.
	IsEmpty() bool

	// KeySet returns a slice of the keys contained in this map.
	KeySet() []K

	// Put associates the specified value with the specified key in this map.
	Put(key K, value *V)

	// Remove removes the mapping for a key from this map if it is present,
	// returning the previous value or nil if none.
	Remove(key K) *V

	// Size returns the number of key-value mappings in this map.
	Size() int

	// Values returns a slice of the values contained in this map.
	Values() []*V

	// All returns an iterator of all key-value pairs in this map.
	// The order of the pairs is not specified.
	All() iter.Seq2[K, *V]
}
