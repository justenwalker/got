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

package semaphore_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/justenwalker/got/semaphore"
	"testing"
	"time"
)

func ExampleSemaphore() {
	// This example simulates concurrently working on an array of 6 elements
	// Each goroutine is responsible for setting the index value to the square of its position.
	// The semaphore only allows 3 goroutines to work on the slice concurrently.
	sem := semaphore.New(3)
	result := make([]int, 5)
	for i := range result {
		// Acquire semaphore before starting goroutine
		_ = sem.Acquire(context.Background())
		go func(i int) {
			defer sem.Release() // free up a slot after goroutine exits
			time.Sleep(1 * time.Millisecond)
			result[i] = i * i
		}(i)
	}
	// wait for all goroutines to finish
	sem.Wait()
	for i, d := range result {
		fmt.Println(i, d)
	}
	// Output:
	// 0 0
	// 1 1
	// 2 4
	// 3 9
	// 4 16
}

func TestSemaphore_Acquire(t *testing.T) {
	tests := []struct {
		name      string
		size      int
		contextFn func(t *testing.T) context.Context
		want      error
	}{
		{
			name: "normal",
			size: 1,
			contextFn: func(t *testing.T) context.Context {
				return context.Background()
			},
			want: nil,
		},
		{
			name: "context_canceled",
			size: 0,
			contextFn: func(t *testing.T) context.Context {
				// A context that is already cancelled
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			want: context.Canceled,
		},
		{
			name: "deadline_exceeded",
			size: 0,
			contextFn: func(t *testing.T) context.Context {
				ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
				t.Cleanup(cancel)
				return ctx
			},
			want: context.DeadlineExceeded,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := semaphore.New(tt.size)

			got := s.Acquire(tt.contextFn(t))

			if !errors.Is(got, tt.want) {
				t.Fatalf("TestAcquire() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSemaphore_TryAcquire(t *testing.T) {
	tests := []struct {
		name string
		size int
		fill bool
		want bool
	}{
		{"success", 1, false, true},
		{"fail_full", 1, true, false},
		{"fail_zero", 0, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := semaphore.New(tt.size)

			// if the fill flag is true, we fill the semaphore
			if tt.fill {
				s <- struct{}{}
			}

			got := s.TryAcquire()
			if got != tt.want {
				t.Errorf("TryAcquire() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSemaphore_Release(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		sem := semaphore.New(1)
		if ok := sem.TryAcquire(); !ok {
			t.Fatalf("expected acquire to succeed, but it failed")
		}

		// Release after acquire should succeed
		tryBlockingOp(t, sem.Release, false)
	})
	t.Run("failed", func(t *testing.T) {
		sem := semaphore.New(1)
		// Release before acquire should fail
		tryBlockingOp(t, sem.Release, true)
	})
}

func TestSemaphore_Wait(t *testing.T) {
	t.Run("fail_outstanding_acquire", func(t *testing.T) {
		sem := semaphore.New(1)
		if ok := sem.TryAcquire(); !ok {
			t.Fatalf("expected acquire to succeed, but it failed")
		}
		// Wait should fail
		tryBlockingOp(t, sem.Wait, true)
	})
	t.Run("success_empty_semaphore", func(t *testing.T) {
		sem := semaphore.New(1)
		// Wait before acquire should fail
		tryBlockingOp(t, sem.Wait, false)
	})
	t.Run("success_acquire_release_semaphore", func(t *testing.T) {
		sem := semaphore.New(1)
		if ok := sem.TryAcquire(); !ok {
			t.Fatalf("expected acquire to succeed, but it failed")
		}
		// should release
		tryBlockingOp(t, sem.Release, false)

		// Release before acquire should fail
		tryBlockingOp(t, sem.Wait, false)
	})
}

func tryBlockingOp(t *testing.T, op func(), shouldBlock bool) {
	t.Helper()
	doneCh := make(chan struct{})
	go func() {
		defer close(doneCh)
		op()
	}()
	select {
	case <-doneCh:
		if shouldBlock {
			t.Fatalf("expected to block indefinitely")
		}
	case <-time.After(1 * time.Second):
		if !shouldBlock {
			t.Fatalf("expected operation to return")
		}
	}
}
