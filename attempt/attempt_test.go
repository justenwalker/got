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

package attempt

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"testing"
	"testing/quick"
	"time"
)

func ExampleWithRetry_decorrelated_jitter() {
	var i int
	r, err := WithRetry(context.Background(), RetryStrategy{
		MaximumAttempts: 3,
		ShouldRetry:     RetryAlways,
		Delayer:         DecorrelatedJitter(rand.Float64, 15*time.Millisecond, time.Second, DefaultDecorrelatedScale),
	}, func(ctx context.Context) (int, error) {
		fmt.Println("called")
		i++
		if i > 2 {
			return 123, nil
		}
		return 0, errors.New("failed")
	})
	fmt.Println("result", r, "error", err)
	// Output:
	// called
	// called
	// called
	// result 123 error <nil>
}

func ExampleWithRetry_exponential_backoff_with_full_jitter() {
	var i int
	r, err := WithRetry(context.Background(), RetryStrategy{
		MaximumAttempts: 3,
		ShouldRetry:     RetryAlways,
		Delayer: FullJitter(rand.Float64, ExponentialBackoff{
			InitialDelay: 15 * time.Millisecond,
			MaxDelay:     1 * time.Second,
			Coefficient:  2.0,
		}.Delay),
	}, func(ctx context.Context) (int, error) {
		fmt.Println("called")
		i++
		if i > 2 {
			return 123, nil
		}
		return 0, errors.New("failed")
	})
	fmt.Println("result", r, "error", err)
	// Output:
	// called
	// called
	// called
	// result 123 error <nil>
}

func ExampleWithTimeout_success() {
	r, err := WithTimeout(context.Background(), time.Second, func(ctx context.Context) (int, error) {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		default:
			return 123, nil
		}
	})
	fmt.Println("result", r, "error", err)
	// Output:
	// result 123 error <nil>
}

func ExampleWithTimeout_deadline_exceeded() {
	r, err := WithTimeout(context.Background(), time.Second, func(ctx context.Context) (int, error) {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		case <-time.After(5 * time.Second):
			return 123, nil
		}
	})
	fmt.Println("result", r, "error", err)
	// Output:
	// result 0 error context deadline exceeded
}

func TestWithRetry(t *testing.T) {
	var retryErr = errors.New("some error")
	type result struct {
		value int
		err   error
	}
	var tests = []struct {
		name          string
		strategy      RetryStrategy
		returnValues  []result
		isContextDone bool
		expectedError error
	}{
		{
			name: "happy_path",
			strategy: RetryStrategy{
				MaximumAttempts: 1,
				ShouldRetry:     RetryAlways,
				Delayer:         Duration(0),
			},
			returnValues: []result{
				{0, nil},
			},
		},
		{
			name: "retry_then_err",
			strategy: RetryStrategy{
				MaximumAttempts: 2,
				ShouldRetry:     RetryAlways,
				Delayer:         Duration(0),
			},
			returnValues: []result{
				{0, retryErr},
				{0, retryErr},
			},
			expectedError: &RetryExhaustedError{Attempt: 2, Err: retryErr},
		},
		{
			name: "ctx_cancelled",
			strategy: RetryStrategy{
				MaximumAttempts: 0,
				ShouldRetry:     RetryAlways,
				Delayer:         Duration(0),
			},
			isContextDone: true,
			expectedError: context.Canceled,
		},
		{
			name: "retry_backoff",
			strategy: RetryStrategy{
				MaximumAttempts: 3,
				ShouldRetry:     RetryAlways,
				Delayer:         Duration(time.Second),
			},
			returnValues: []result{
				{0, retryErr},
				{0, retryErr},
				{0, retryErr},
			},
			expectedError: &RetryExhaustedError{Attempt: 3, Err: retryErr},
		},
		{
			name: "retry_no_delay",
			strategy: RetryStrategy{
				MaximumAttempts: 1,
				ShouldRetry:     RetryNever,
				Delayer:         Duration(0),
			},
			returnValues: []result{
				{0, nil},
			},
		},
		{
			name: "retry_no_backoff",
			strategy: RetryStrategy{
				MaximumAttempts: 1,
				ShouldRetry:     RetryNever,
			},
			returnValues: []result{
				{0, nil},
			},
		},
		{
			name: "no_retry",
			strategy: RetryStrategy{
				MaximumAttempts: 1,
			},
			returnValues: []result{
				{0, nil},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			// Simulate context.Done signal
			if tt.isContextDone {
				var cancel context.CancelFunc
				ctx, cancel = context.WithCancel(ctx)
				cancel()
			}

			var attemptCount int
			_, err := WithRetry(ctx, tt.strategy, func(ctx context.Context) (int, error) {
				if attemptCount < len(tt.returnValues) {
					rv := tt.returnValues[attemptCount]
					attemptCount++
					return rv.value, rv.err
				}
				t.Fatalf("unexpected call to retry %d", attemptCount)
				return 0, nil
			})
			if attemptCount < len(tt.returnValues) {
				t.Fatalf("expected %d calls, got %d", len(tt.returnValues), attemptCount)
			}

			if tt.expectedError != nil && err == nil || tt.expectedError == nil && err != nil {
				t.Errorf("WithRetry error = %v, want %v", err, tt.expectedError)
			}

			if tt.expectedError != nil && err != nil &&
				tt.expectedError.Error() != err.Error() {
				t.Errorf("WithRetry error = %v, want %v", err, tt.expectedError)
			}
		})
	}
}

type mockRetryEvaluator struct {
	shouldRetry bool
}

func (mr *mockRetryEvaluator) ShouldRetry(err error) bool {
	return mr.shouldRetry
}

type mockDelayer struct {
	delayDuration time.Duration
}

func (md *mockDelayer) Delay(int) time.Duration {
	return md.delayDuration
}

func TestWithTimeout(t *testing.T) {
	mockErr := errors.New("mock error")
	cancelledContext, cancel := context.WithCancel(context.Background())
	cancel()
	tests := []struct {
		name          string
		ctx           context.Context
		timeout       time.Duration
		fn            func(ctx context.Context) (int, error)
		expectedValue int
		expectedErr   error
	}{
		{
			name:          "successfully ran within timeout",
			ctx:           context.Background(),
			timeout:       1 * time.Second,
			fn:            mockTimeoutFunc(500*time.Millisecond, 1, nil),
			expectedValue: 1,
			expectedErr:   nil,
		},
		{
			name:          "context already cancelled",
			ctx:           cancelledContext,
			timeout:       1 * time.Second,
			fn:            mockTimeoutFunc(500*time.Millisecond, 1, mockErr),
			expectedValue: 0,
			expectedErr:   context.Canceled,
		},
		{
			name:          "execution timed out",
			ctx:           context.Background(),
			timeout:       1 * time.Second,
			fn:            mockTimeoutFunc(2*time.Second, 1, nil),
			expectedValue: 0,
			expectedErr:   context.DeadlineExceeded,
		},
		{
			name:          "returns error",
			ctx:           context.Background(),
			timeout:       1 * time.Second,
			fn:            mockTimeoutFunc(500*time.Millisecond, 1, mockErr),
			expectedValue: 0,
			expectedErr:   mockErr,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := WithTimeout[int](test.ctx, test.timeout, test.fn)

			if got != test.expectedValue {
				t.Errorf("Got value %v, expected %v", got, test.expectedValue)
			}
			if !errors.Is(err, test.expectedErr) {
				t.Errorf("Got error %v, expected error %v", err, test.expectedErr)
			}
		})
	}
}

func mockTimeoutFunc(delay time.Duration, num int, err error) func(context.Context) (int, error) {
	return func(ctx context.Context) (int, error) {
		time.Sleep(delay)
		return num, err
	}
}

func Test_FullJitter(t *testing.T) {
	type args struct {
		rnd     Rand
		delayer Delayer
	}
	tests := []struct {
		name string
		args args
		want time.Duration
	}{
		{
			name: "test_1",
			args: args{
				rnd: func() float64 {
					return 0.5
				},
				delayer: Duration(time.Second),
			},
			want: time.Second / 2,
		},
		{
			name: "test_2",
			args: args{
				rnd: func() float64 {
					return 0.2
				},
				delayer: Duration(time.Minute),
			},
			want: time.Minute / 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FullJitter(tt.args.rnd, tt.args.delayer)(1); got != tt.want {
				t.Errorf("FullJitter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_EqualJitter(t *testing.T) {
	type args struct {
		rnd     Rand
		delayer Delayer
	}
	tests := []struct {
		name string
		args args
		want time.Duration
	}{
		{
			name: "test_1",
			args: args{
				rnd: func() float64 {
					return 0.5
				},
				delayer: Duration(time.Second),
			},
			want: 3 * time.Second / 4,
		},
		{
			name: "test_2",
			args: args{
				rnd: func() float64 {
					return 0.2
				},
				delayer: Duration(time.Minute),
			},
			want: 3 * time.Minute / 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EqualJitter(tt.args.rnd, tt.args.delayer)(1); got != tt.want {
				t.Errorf("EqualJitter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecorrelatedJitter(t *testing.T) {
	tests := []struct {
		name  string
		rnd   Rand
		base  time.Duration
		cap   time.Duration
		scale float64
		want  time.Duration
	}{
		{
			name:  "min_base_cap",
			rnd:   func() float64 { return 0 },
			base:  time.Duration(1),
			cap:   time.Duration(1),
			scale: 2,
			want:  time.Duration(1),
		},
		{
			name:  "rnd_with_cap",
			rnd:   func() float64 { return 0.5 },
			base:  time.Duration(1),
			cap:   time.Duration(3),
			scale: 2,
			want:  time.Duration(1),
		},
		{
			name:  "exceed_cap",
			rnd:   func() float64 { return 1 },
			base:  time.Duration(1),
			cap:   time.Duration(2),
			scale: 3,
			want:  time.Duration(2),
		},
		{
			name:  "scale_below_one",
			rnd:   func() float64 { return 0.5 },
			base:  time.Duration(1),
			cap:   time.Duration(5),
			scale: 0.5,
			want:  time.Duration(1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			delayer := DecorrelatedJitter(tt.rnd, tt.base, tt.cap, tt.scale)
			got := delayer(1)
			if got != tt.want {
				t.Errorf("DecorrelatedJitter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExponentialBackoff_Delay(t *testing.T) {
	tests := []struct {
		name               string
		ExponentialBackoff ExponentialBackoff
		attempt            int
		expectedDuration   time.Duration
	}{
		{
			"ZeroAttempt",
			ExponentialBackoff{
				InitialDelay: 1 * time.Second,
				MaxDelay:     10 * time.Second,
				Coefficient:  2,
			},
			0,
			0,
		},
		{
			"DelayWithinBounds",
			ExponentialBackoff{
				InitialDelay: 1 * time.Second,
				MaxDelay:     10 * time.Second,
				Coefficient:  2,
			},
			2,
			4 * time.Second,
		},
		{
			"DelayExceedsMaximum",
			ExponentialBackoff{
				InitialDelay: 1 * time.Second,
				MaxDelay:     10 * time.Second,
				Coefficient:  2,
			},
			4,
			10 * time.Second, // this would be 16 seconds without the max cap
		},
		{
			"DifferentCoefficient",
			ExponentialBackoff{
				InitialDelay: 1 * time.Second,
				MaxDelay:     100 * time.Second, // Increase maximum delay
				Coefficient:  3,                 // Change coefficient from 2 to 3
			},
			2,
			9 * time.Second, // With coefficient of 3, the delay for n = 2 should be 9 seconds
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ExponentialBackoff.Delay(tt.attempt); got != tt.expectedDuration {
				t.Errorf("ExponentialBackoff.Delay() = %v, want = %v", got, tt.expectedDuration)
			}
		})
	}
}

func TestDuration_Delay(t *testing.T) {
	dur := Duration(time.Hour)
	err := quick.Check(func(i int) bool {
		return dur(i) == time.Hour
	}, nil)
	if err != nil {
		t.Fatal("ERROR", err)
	}
}
