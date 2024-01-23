package ptr

// To is a generic function that takes a non-pointer type T and returns a ptr to it.
// This is useful for taking a ptr to literals; such as Of("string")
func To[T any](t T) *T {
	return &t
}

// Value is a generic function that takes a pointer to T and dereferences the value safely.
// If the pointer is null, then the zero-value of that type is returned.
func Value[T any](pt *T) T {
	if pt == nil {
		var zero T
		return zero
	}
	return *pt
}
