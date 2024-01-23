package optional

import "testing"

func TestCoalesce(t *testing.T) {
	tests := []struct {
		name   string
		input  []Value[int]
		expect Value[int]
	}{
		{
			name:   "nil",
			input:  nil,
			expect: Nothing[int](),
		},
		{
			name:   "empty",
			input:  []Value[int]{},
			expect: Nothing[int](),
		},
		{
			name: "one",
			input: []Value[int]{
				New(1),
			},
			expect: New(1),
		},
		{
			name: "nil-one",
			input: []Value[int]{
				Nothing[int](),
				Nothing[int](),
				New(1),
			},
			expect: New(1),
		},
		{
			name: "nil-one-two",
			input: []Value[int]{
				Nothing[int](),
				Nothing[int](),
				New(1),
				New(2),
				New(3),
			},
			expect: New(1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Coalesce(tt.input...)
			a, ok1 := v.Get()
			b, ok2 := tt.expect.Get()
			if ok1 != ok2 || a != b {
				t.Errorf("Coalesce() = (%v,%t), want (%v,%t)", a, ok1, b, ok2)
			}
		})
	}
}

func TestMapSlice(t *testing.T) {
	tests := []struct {
		name   string
		input  []Value[int]
		expect []Value[int64]
	}{
		{
			name:   "nil",
			input:  nil,
			expect: nil,
		},
		{
			name:   "empty",
			input:  []Value[int]{},
			expect: []Value[int64]{},
		},
		{
			name: "one",
			input: []Value[int]{
				New(1),
			},
			expect: []Value[int64]{
				New(int64(1)),
			},
		},
		{
			name: "nil-one",
			input: []Value[int]{
				Nothing[int](),
				Nothing[int](),
				New(1),
			},
			expect: []Value[int64]{
				Nothing[int64](),
				Nothing[int64](),
				New(int64(1)),
			},
		},
		{
			name: "nil-one-two",
			input: []Value[int]{
				Nothing[int](),
				Nothing[int](),
				New(1),
				New(2),
				New(3),
			},
			expect: []Value[int64]{
				Nothing[int64](),
				Nothing[int64](),
				New(int64(1)),
				New(int64(2)),
				New(int64(3)),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vs := MapSlice[int, int64](tt.input, func(a int) int64 {
				return int64(a)
			})
			if len(tt.expect) != len(vs) {
				t.Fatalf("MapSlice() = length was %v, want %v", len(vs), len(tt.expect))
			}
			for i, v := range vs {
				if tt.expect[i] != v {
					t.Fatalf("MapSlice() = arr[%d]=%v, want %v", i, v, tt.expect[i])
				}
			}
		})
	}
}
