package assert

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

// Equal verifies equality of two objects.
func Equal[T any](t *testing.T, a, b T) {
	if !reflect.DeepEqual(a, b) {
		t.Helper()
		t.Fatalf("%v != %v", a, b)
	}
}

// NotEqual asserts that the given objects are not equal.
func NotEqual[T any](t *testing.T, a, b T) {
	if reflect.DeepEqual(a, b) {
		t.Helper()
		t.Fatalf("%v == %v", a, b)
	}
}

// Same asserts that the given pointers point to the same object.
func Same[T any](t *testing.T, a, b *T) {
	if a != b {
		t.Helper()
		t.Fatalf("%p != %p", a, b)
	}
}

// IsNil verifies that the object is nil.
func IsNil(t *testing.T, obj any) {
	if obj != nil {
		value := reflect.ValueOf(obj)
		switch value.Kind() {
		case reflect.Ptr, reflect.Map, reflect.Slice,
			reflect.Interface, reflect.Func, reflect.Chan:
			if value.IsNil() {
				return
			}
		}
		t.Helper()
		t.Fatalf("%v is not nil", obj)
	}
}

// ElementsMatch checks whether the given slices contain the same elements.
func ElementsMatch[S ~[]E, E any](t *testing.T, a, b S) {
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
	t.Helper()
	t.Fatalf("Slice elements are not equal: %v != %v", a, b)
}

// ErrorContains checks whether the given error contains the specified string.
func ErrorContains(t *testing.T, err error, str string) {
	if err == nil {
		t.Helper()
		t.Fatalf("Error is nil")
	} else if !strings.Contains(err.Error(), str) {
		t.Helper()
		t.Fatalf("Error does not contain string: %s", str)
	}
}

// ErrorIs checks whether any error in err's tree matches target.
func ErrorIs(t *testing.T, err error, target error) {
	if !errors.Is(err, target) {
		t.Helper()
		t.Fatalf("Error type mismatch: %v != %v", err, target)
	}
}

// Panics checks whether the given function panics.
func Panics(t *testing.T, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Helper()
			t.Errorf("Function did not panic")
		}
	}()
	f()
}

// PanicMsgContains checks whether a function panics with a message containing
// the given string.
func PanicMsgContains(t *testing.T, f func(), str string) {
	defer func() {
		if r := recover(); r != nil {
			panicMessage := fmt.Sprintf("%v", r)
			if !strings.Contains(panicMessage, str) {
				t.Helper()
				t.Errorf("Panic message does not contain: %s", str)
			}
		} else {
			t.Helper()
			t.Errorf("Function did not panic")
		}
	}()
	f()
}
