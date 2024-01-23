package fault

// test

import "errors"

// IsNotFound checks if the error or any of its wrapped errors is a NotFound error.
// a NotFound error implements the function `NotFound() bool` and returns true.
func IsNotFound(err error) bool {
	var asErr interface{ NotFound() bool }
	return errors.As(err, &asErr) && asErr.NotFound()
}
