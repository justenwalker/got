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
