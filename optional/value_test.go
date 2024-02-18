// Copyright (c) 2024 Justen Walker
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
//
// SPDX-License-Identifier: MIT

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
