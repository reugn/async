package internal

import (
	"reflect"
	"testing"
)

// AssertEqual verifies equality of two objects
func AssertEqual(t *testing.T, a interface{}, b interface{}) {
	if !reflect.DeepEqual(a, b) {
		t.Fatalf("%v != %v", a, b)
	}
}
