package ptr

// Of returns a pointer to the given value.
func Of[T any](value T) *T {
	return &value
}
