package assert

import (
	"reflect"
	"strings"
	"testing"
)

// Equal verifies equality of two objects.
func Equal[T any](t *testing.T, a, b T) {
	if !reflect.DeepEqual(a, b) {
		t.Fatalf("%v != %v", a, b)
	}
}

// ElementsMatch checks whether the given slices contain the same elements.
func ElementsMatch[T any](t *testing.T, a, b []T) {
	if len(a) == len(b) {
		count := 0
		for _, va := range a {
			for _, vb := range b {
				if reflect.DeepEqual(va, vb) {
					count++
					break
				}
			}
		}
		if count == len(a) {
			return
		}
	}
	t.Fatalf("Slice elements are not equal: %v != %v", a, b)
}

// ErrorContains checks whether the given error contains the specified string.
func ErrorContains(t *testing.T, err error, str string) {
	if err == nil {
		t.Fatalf("Error is nil")
	} else if !strings.Contains(err.Error(), str) {
		t.Fatalf("Error doen't contain string: %s", str)
	}
}

// Panic checks whether the given function panics.
func Panic(t *testing.T, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The function did not panic")
		}
	}()
	f()
}
