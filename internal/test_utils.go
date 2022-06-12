package internal

import (
	"reflect"
	"testing"
)

// AssertEqual verifies equality of two objects.
func AssertEqual[T any](t *testing.T, a T, b T) {
	if !reflect.DeepEqual(a, b) {
		t.Fatalf("%v != %v", a, b)
	}
}
