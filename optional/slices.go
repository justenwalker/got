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
