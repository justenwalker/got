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
