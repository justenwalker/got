package optional

import "testing"

func TestValue(t *testing.T) {
	ni := New(123)
	if v, ok := ni.Get(); !ok || v != 123 {
		t.Errorf("Expected nil.Get() = (123,true); got (%v,%t)", v, ok)
	}
	if !ni.IsValid() {
		t.Errorf("Expected ni.IsValid() to be true")
	}
	var pass bool
	ni.WithValue(func(val int) {
		if val != 123 {
			t.Errorf("expected val=123, but got=%v", val)
		}
		pass = true
	})
	if !pass {
		t.Errorf("expected nil.WithValue to execute function, but it didn't")
	}
	if 123 != ni.GetWithDefault(456) {
		t.Errorf("Expected GetWithDefault to return the value, but it returned the default value instead")
	}
	nb := Map(ni, func(a int) int64 {
		return int64(a)
	})
	if v, ok := nb.Get(); !ok || v != 123 {
		t.Errorf("Expected nil.Get() = (true,true); got (%v,%t)", v, ok)
	}
}

func TestNothing(t *testing.T) {
	ni := Nothing[int]()
	if v, ok := ni.Get(); ok {
		t.Errorf("Expected Get() to be invalid, but it is valid: got (%v,%t)", v, ok)
	}
	if ni.IsValid() {
		t.Errorf("Expected Nothing() to be invalid, but it is valid")
	}
	ni.WithValue(func(val int) {
		t.Errorf("WithValue should not be called on Nothing()")
	})
	if 123 != ni.GetWithDefault(123) {
		t.Errorf("Expected GetWithDefault to return default value, but it returned wrapped value")
	}
	nb := Map[int, bool](ni, func(a int) bool {
		t.Errorf("Map should not be called on Nothing()")
		return true
	})
	if nb.IsValid() {
		t.Errorf("Expected nb.IsValue() to be false")
	}
}
