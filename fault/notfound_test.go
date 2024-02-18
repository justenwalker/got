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

package fault

import (
	"errors"
	"fmt"
	"testing"
)

type testNotFoundError struct {
	Value bool
}

func (e *testNotFoundError) Error() string {
	return "test"
}

func (e *testNotFoundError) NotFound() bool {
	return e.Value
}

func TestIsNotFound(t *testing.T) {
	tests := []struct {
		name   string
		errVal error
		expect bool
	}{
		{
			name:   "nil",
			errVal: nil,
			expect: false,
		},
		{
			name:   "standard-error",
			errVal: errors.New("standard error"),
			expect: false,
		},
		{
			name:   "not-found-true",
			errVal: &testNotFoundError{Value: true},
			expect: true,
		},
		{
			name:   "not-found-false",
			errVal: &testNotFoundError{Value: false},
			expect: false,
		},
		{
			name:   "wrapped-not-found-true",
			errVal: fmt.Errorf("test: %w", &testNotFoundError{Value: true}),
			expect: true,
		},
		{
			name:   "wrapped-not-found-false",
			errVal: fmt.Errorf("test: %w", &testNotFoundError{Value: false}),
			expect: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := IsNotFound(tt.errVal)
			if tt.expect != actual {
				t.Errorf("expected=%t, got=%t", tt.expect, actual)
			}
		})
	}
}
