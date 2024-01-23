package optional

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
)

func ExampleValue_MarshalJSON() {
	type myStruct struct {
		Value1 *Value[int] `json:"value1,omitempty"`
		Value2 Value[int]  `json:"value2"`
	}
	mv := myStruct{
		Value1: New(123).Ptr(),
		Value2: Value[int]{Wrapped: 456, Valid: true},
	}
	data, _ := json.Marshal(mv)
	fmt.Println(string(data))
	// Output:
	// {"value1":123,"value2":456}
}

func ExampleValue_UnmarshalJSON() {
	type myStruct struct {
		Value1 *Value[int] `json:"value1,omitempty"`
		Value2 Value[int]  `json:"value2"`
	}
	var mv myStruct
	_ = json.Unmarshal([]byte(`{"value1":123,"value2":456}`), &mv)
	fmt.Println(mv.Value1.Get())
	fmt.Println(mv.Value2.Get())
	// Output:
	// 123 true
	// 456 true
}

func TestValue_MarshalJSON(t *testing.T) {
	type myStruct struct {
		Value1 *Value[int] `json:"value1,omitempty"`
		Value2 Value[int]  `json:"value2"`
	}
	tests := []struct {
		name   string
		obj    any
		expect []byte
	}{
		{
			name: "nothing-nil",
			obj: &myStruct{
				Value1: Nothing[int]().Ptr(),
			},
			expect: []byte(`{"value2":null}`),
		},
		{
			name: "nothing-value",
			obj: &myStruct{
				Value1: Nothing[int]().Ptr(),
				Value2: Value[int]{Valid: true, Wrapped: 123},
			},
			expect: []byte(`{"value2":123}`),
		},
		{
			name: "value",
			obj: &myStruct{
				Value1: New(123).Ptr(),
				Value2: New(456),
			},
			expect: []byte(`{"value1":123,"value2":456}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.obj)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !bytes.Equal(data, tt.expect) {
				t.Fatalf("unexpected data:\n%s\n\nwanted:\n%s", string(data), string(tt.expect))
			}
		})
	}
}

func TestValue_UnmarshalJSON(t *testing.T) {
	type myStruct struct {
		Value1 *Value[int] `json:"value1,omitempty"`
		Value2 Value[int]  `json:"value2"`
	}
	tests := []struct {
		name   string
		data   []byte
		expect myStruct
	}{
		{
			name: "null",
			data: []byte(`null`),
			expect: myStruct{
				Value1: Nothing[int]().Ptr(),
			},
		},
		{
			name: "empty",
			data: []byte(`{}`),
			expect: myStruct{
				Value1: Nothing[int]().Ptr(),
			},
		},
		{
			name: "value1",
			data: []byte(`{"value1":123}`),
			expect: myStruct{
				Value1: New[int](123).Ptr(),
			},
		},
		{
			name: "value2",
			data: []byte(`{"value2":123}`),
			expect: myStruct{
				Value2: New[int](123),
			},
		},
		{
			name: "null-value2",
			data: []byte(`{"value1":null,"value2":456}`),
			expect: myStruct{
				Value1: Nothing[int]().Ptr(),
				Value2: New[int](456),
			},
		},
		{
			name: "value1-null",
			data: []byte(`{"value1":123,"value2":null}`),
			expect: myStruct{
				Value1: New[int](123).Ptr(),
			},
		},
		{
			name: "value1-value2",
			data: []byte(`{"value1":123,"value2":456}`),
			expect: myStruct{
				Value1: New[int](123).Ptr(),
				Value2: New[int](456),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var actual myStruct
			err := json.Unmarshal(tt.data, &actual)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			a1 := actual.Value1
			a2 := actual.Value2
			a1v, a1ok := a1.Get()
			a2v, a2ok := a2.Get()
			b1v, b1ok := tt.expect.Value1.Get()
			b2v, b2ok := tt.expect.Value2.Get()
			if a1ok != b1ok || a1v != b1v {
				t.Errorf("myStruct.Value1: (%v,%t), want (%v,%t)", a1v, a1ok, b1v, b1ok)
			}
			if a2ok != b2ok || a2v != b2v {
				t.Errorf("myStruct.Value2: (%v,%t), want (%v,%t)", a2v, a2ok, b2v, b2ok)
			}
		})
	}
}

func TestValue_UnmarshalJSON_error(t *testing.T) {
	type myStruct struct {
		Value1 *Value[int] `json:"value1,omitempty"`
		Value2 Value[int]  `json:"value2"`
	}
	var out myStruct
	err := json.Unmarshal([]byte(`{"value1":true,"value2":false}`), &out)
	if err == nil {
		t.Fatal("expected json unmarshal error")
	}
}
