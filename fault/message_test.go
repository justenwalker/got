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
	"testing"
)

const (
	testErr1 = Message("error: 1")
	testErr2 = Message("error: 2")
)

func TestMessage_Error(t *testing.T) {
	var err1 error = testErr1
	var err2 error = testErr2
	if err1.Error() != string(testErr1) {
		t.Errorf("err1.Error() = %s, want %s", err1.Error(), string(testErr1))
	}
	testExpectTrueHelper(t, errors.Is(err1, testErr1), "errors.Is(err1, testErr1)")
	testExpectTrueHelper(t, errors.Is(err2, testErr2), "errors.Is(err2, testErr2)")
	testExpectTrueHelper(t, errors.Is(err1, Message("error: 1")), `errors.Is(err1, Message("error: 1")`)
	testExpectTrueHelper(t, !errors.Is(err1, err2), "!errors.Is(err1, err2)")
}

func testExpectTrueHelper(t *testing.T, b bool, msg string) {
	t.Helper()
	if !b {
		t.Error("expectation failed:", msg)
	}
}
