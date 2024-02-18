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

// UnmarshalJSON unmarshals the JSON data into the Value of type T.
// If the JSON data is 'null', the Value is Nothing.
func (v *Value[T]) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, nullBytes) {
		*v = Nothing[T]()
		return nil
	}
	var t T
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}
	*v = Value[T]{Wrapped: t, Valid: true}
	return nil
}
