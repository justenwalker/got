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

// Package semaphore provides a basic implementation of a semaphore.
// A semaphore is a synchronization primitive that limits the number of concurrent accesses to a shared resource.
package semaphore

import "context"

// Semaphore represents a synchronization primitive that limits the number of goroutines
// that can access a certain resource or a section of code simultaneously.
//
// Example usage:
//
// Creating a Semaphore:
// sem := New(size)
//
// Acquiring a Semaphore:
// err := sem.Acquire(ctx)
//
// Trying to acquire a Semaphore without blocking:
// acquired := sem.TryAcquire()
//
// Releasing a Semaphore:
// sem.Release()
//
// Waiting until all resources are released:
// sem.Wait()
type Semaphore chan struct{}

// New creates a new Semaphore with the specified size.
func New(size int) Semaphore {
	return make(chan struct{}, size)
}

// Acquire acquires the semaphore by blocking until it is available.
// If the context is cancelled, then the semaphore is not acquired, and an error is returned.
func (s Semaphore) Acquire(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case s <- struct{}{}:
		return nil
	}
}

// TryAcquire tries to acquire a semaphore without blocking; returns true if successfully acquired.
func (s Semaphore) TryAcquire() bool {
	select {
	case s <- struct{}{}:
		return true
	default:
		return false
	}
}

// Release releases the semaphore.
// This MUST be called after a successful call to Acquire or TryAcquire to release the resource
// Failing to call this may lead to deadlocks.
func (s Semaphore) Release() {
	<-s
}

// Wait waits for all semaphore acquires to be released back to the pool.
// After the call to wait, the semaphore should not be re-used.
func (s Semaphore) Wait() {
	for i := 0; i < cap(s); i++ {
		s <- struct{}{}
	}
}
