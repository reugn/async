package internal

import (
	"reflect"
	"strings"
	"testing"
)

// AssertEqual verifies equality of two objects.
func AssertEqual[T any](t *testing.T, a, b T) {
	if !reflect.DeepEqual(a, b) {
		t.Fatalf("%v != %v", a, b)
	}
}

// AssertElementsMatch checks whether the given slices contain the same elements.
func AssertElementsMatch[T any](t *testing.T, a, b []T) {
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

// AssertErrorContains checks whether the given error contains the specified string.
func AssertErrorContains(t *testing.T, err error, str string) {
	if err == nil {
		t.Fatalf("Error is nil")
	} else if !strings.Contains(err.Error(), str) {
		t.Fatalf("Error doen't contain string: %s", str)
	}
}

// AssertPanic checks whether the given function panics.
func AssertPanic(t *testing.T, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The function did not panic")
		}
	}()
	f()
}
