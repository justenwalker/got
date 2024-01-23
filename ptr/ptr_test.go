package ptr

import (
	"testing"
)

func TestTo(t *testing.T) {
	testPtrIsEqual[int](t, 123123, To[int](123123))
	testPtrIsEqual[int8](t, 127, To[int8](127))
	testPtrIsEqual[int16](t, 1<<14, To[int16](1<<14))
	testPtrIsEqual[int32](t, 1<<30, To[int32](1<<30))
	testPtrIsEqual[int64](t, 1<<62, To[int64](1<<62))
	testPtrIsEqual[uint](t, 123123, To[uint](123123))
	testPtrIsEqual[uint8](t, 255, To[uint8](255))
	testPtrIsEqual[uint16](t, 1<<15, To[uint16](1<<15))
	testPtrIsEqual[uint32](t, 1<<31, To[uint32](1<<31))
	testPtrIsEqual[uint64](t, 1<<63, To[uint64](1<<63))
	testPtrIsEqual[float32](t, 1.23, To[float32](1.23))
	testPtrIsEqual[float64](t, 1.64, To[float64](1.64))
	testPtrIsEqual[string](t, "test", To[string]("test"))
	testPtrIsEqual[bool](t, true, To[bool](true))
}

func TestValue(t *testing.T) {
	testIsEqual[int](t, 0, Value[int](nil))
	testIsEqual[uint](t, 123123, Value[uint](To[uint](123123)))
	testIsEqual[float32](t, 1.23, Value[float32](To[float32](1.23)))
	testIsEqual[float64](t, 1.64, Value[float64](To[float64](1.64)))
	testIsEqual[string](t, "test", Value[string](To[string]("test")))
	testIsEqual[bool](t, true, Value[bool](To[bool](true)))
}

func testPtrIsEqual[T comparable](t *testing.T, expected T, ptrIn *T) {
	t.Helper()
	if ptrIn == nil {
		t.Errorf("%T(%[1]v): expected non-nil pointer", ptrIn)
	}
	testIsEqual(t, expected, *ptrIn)
}

func testIsEqual[T comparable](t *testing.T, expected T, actual T) {
	t.Helper()
	if actual != expected {
		t.Errorf("expected=%[1]T(%[1]v), got=%[2]T(%[2]v)", expected, actual)
	}
}
