package optional

import (
	"bytes"
	"encoding/json"
)

var nullBytes = []byte(`null`)

// MarshalJSON marshals the wrapped value of type T to JSON.
// If the value is valid, it returns the JSON representation of the wrapped value.
// If the value is not valid, it returns a JSON 'null'
func (v Value[T]) MarshalJSON() ([]byte, error) {
	if v.IsValid() {
		return json.Marshal(v.Wrapped)
	}
	return nullBytes, nil
}

func (v *Value[T]) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, nullBytes) {
		*v = Value[T]{}
		return nil
	}
	var t T
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}
	*v = Value[T]{Wrapped: t, Valid: true}
	return nil
}