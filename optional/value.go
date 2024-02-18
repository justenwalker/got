package optional

// New creates a Value wrapping type T with the given concrete value t.
func New[T any](t T) Value[T] {
	return Value[T]{
		Wrapped: t,
		Valid:   true,
	}
}

// Nothing creates an unset/invalid Value.
// A nil pointer to any Value is also Nothing.
func Nothing[T any]() Value[T] {
	return Value[T]{}
}

// Value is a generic type that wraps a value of any type T.
//
// A Value has several method to support interacting with values (set or unset) in a way that doesn't panic.
//
// The Value type also supports JSON marshaling and unmarshalling. A valid Value serializes to its contained value.
// An invalid or unset Value serializes to null. You can support omitempty by using a pointer to a Value,
// as the Value type:
//
//	type MyJSONStruct struct {
//	    Int *Value[int] `json:"int,omitempty"`
//	}
type Value[T any] struct {
	// Value is the wrapped value of type T
	Wrapped T
	// Set indicates if Wrapped is valid
	Valid bool
}

// Ptr returns a pointer to the current Value instance.
func (v Value[T]) Ptr() *Value[T] {
	if v.IsValid() {
		return &v
	}
	return nil
}

// Get returns the wrapped value and a boolean indicating if it is valid.
func (v *Value[T]) Get() (T, bool) {
	if v == nil {
		var z T
		return z, false
	}
	return v.Wrapped, v.Valid
}

// Dereference returns a new Value[T] that is a dereferenced copy of the receiver, or an empty Value[T] if the receiver is nil.
func (v *Value[T]) Dereference() Value[T] {
	if v == nil {
		return Value[T]{}
	}
	return *v
}

// GetWithDefault returns the wrapped value if it is valid, otherwise it returns the default value passed.
func (v *Value[T]) GetWithDefault(def T) T {
	if v.IsValid() {
		return v.Wrapped
	}
	return def
}

// IsValid checks if the Value is valid.
func (v *Value[T]) IsValid() bool {
	if v == nil {
		return false
	}
	return v.Valid
}

// WithValue calls the provided function `fn` if the `Value` is valid.
// The function takes the wrapped value of type `T` as a parameter.
func (v *Value[T]) WithValue(fn func(val T)) {
	if v.IsValid() {
		fn(v.Wrapped)
	}
}

// Map applies the given map function which maps type A -> B.
// The function takes a wrapped value of type A and returns a new wrapped value of type B.
// If a is not valid, it returns Nothing[B]()
func Map[A any, B any](a Value[A], mapFn func(a A) B) Value[B] {
	if a.IsValid() {
		return Value[B]{
			Wrapped: mapFn(a.Wrapped),
			Valid:   true,
		}
	}
	return Nothing[B]()
}
