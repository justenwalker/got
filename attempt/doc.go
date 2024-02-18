// Package attempt provides tools for retrying and adding timeouts to actions.
//
// It contains two main functions: WithRetry and WithTimeout.
//
// WithRetry allows you to repeatedly attempt an action under a defined RetryStrategy.
// This is useful when dealing with network requests or other uncertain actions
// that might fail temporarily. By using WithRetry, you can ensure that temporary
// failures don't cause your program to terminate or progress in an unwanted state.
//
// WithTimeout allows you to run an action, but abort if it takes longer than expected.
// This is useful when you're dealing with actions that might get stuck, such as network
// requests, file reads, etc. Instead of your program hanging indefinitely, WithTimeout
// allows you to decide how long you're willing to wait and abort if necessary.
package attempt
