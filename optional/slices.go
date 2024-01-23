package optional

// Coalesce takes multiple Value[T] as input and returns the first valid Value[T].
// If all the Value[T] are invalid or there are no Value[T] provided, it returns Nothing[T]().
func Coalesce[T any](vals ...Value[T]) Value[T] {
	for _, v := range vals {
		if v.IsValid() {
			return v
		}
	}
	return Nothing[T]()
}

// MapSlice takes a slice of Value[A] and a mapper function as input and returns a new slice of Value[B].
// The mapper function, mapFn, is applied to each valid element in the input slice, and the result is used to create a new Value[B] in the output slice.
// If an element in the input slice is invalid, the corresponding element in the output slice will be invalid as well.
// The input slice values is not modified.
func MapSlice[A any, B any](values []Value[A], mapFn func(a A) B) []Value[B] {
	result := make([]Value[B], len(values))
	for i, a := range values {
		if a.IsValid() {
			result[i] = Value[B]{
				Wrapped: mapFn(a.Wrapped),
				Valid:   true,
			}
		}
	}
	return result
}
