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
	"fmt"
	"math"
	"time"
)

// RetryExhaustedError is an error that is returned by WithRetry when the maximum attempts have been exhausted.
type RetryExhaustedError struct {
	// Attempt is the attempt that failed.
	Attempt int
	// Err is the last error returned by the retried function.
	Err error
}

func (e *RetryExhaustedError) Error() string {
	return fmt.Sprintf("attempt: retry exhausted after %d attempts. last error: %v", e.Attempt, e.Err)
}

func (e *RetryExhaustedError) Unwrap() error {
	return e.Err
}

// WithRetry retries the Call using the RetryStrategy provided
func WithRetry[T any](ctx context.Context, rs RetryStrategy, fn func(ctx context.Context) (T, error)) (T, error) {
	var zero T
	if rs.ShouldRetry == nil {
		return fn(ctx)
	}
	// don't run if context is already finished
	select {
	case <-ctx.Done():
		return zero, ctx.Err()
	default:
	}
	var attempt int
	for {
		attempt++
		t, err := fn(ctx)
		if err == nil {
			return t, nil
		}
		if !rs.ShouldRetry(err) {
			return zero, err
		}
		if rs.MaximumAttempts != 0 && attempt >= rs.MaximumAttempts {
			return zero, &RetryExhaustedError{
				Attempt: attempt,
				Err:     err,
			}
		}
		var delay time.Duration
		if rs.Delayer != nil {
			delay = rs.Delayer(attempt)
		}
		if delay == 0 {
			select {
			case <-ctx.Done():
				return zero, ctx.Err()
			default:
			}
			continue
		}
		ticker := time.NewTicker(delay)
		select {
		case <-ctx.Done():
			ticker.Stop()
			return zero, ctx.Err()
		case <-ticker.C:

		}
	}
}

type result[T any] struct {
	value T
	err   error
}

// WithTimeout calls the given function and returns early if the function takes longer than the timeout provided.
//
// Note: The function is called with a context that is cancelled after the timeout duration.
// The function provided should therefore support cancellation via context, otherwise this may leak resources.
func WithTimeout[T any](ctx context.Context, timeout time.Duration, fn func(ctx context.Context) (T, error)) (T, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	resultCh := make(chan result[T], 1)
	defer cancel()
	go func() {
		t, err := fn(ctx)
		resultCh <- result[T]{value: t, err: err}
	}()
	var zero T
	select {
	case <-ctx.Done():
		return zero, ctx.Err()
	case r := <-resultCh:
		if r.err != nil {
			return zero, r.err
		}
		return r.value, nil
	}
}

// RetryStrategy represents a strategy for retrying a specific operation in WithRetry.
type RetryStrategy struct {
	// MaximumAttempts limits the number of attempts.
	// If MaximumAttempts is 0, retries will be performed indefinitely.
	MaximumAttempts int
	// ShouldRetry is responsible for evaluating whether to retry an operation.
	// if it is not set, no retries will be performed
	ShouldRetry func(err error) bool
	// Delayer is responsible for determining the delay duration before the next retry attempt.
	// If it is not set, there will be no delays between retries.
	Delayer func(attempt int) time.Duration
}

// RetryAlways always returns true, allowing a retry for any error.
func RetryAlways(_ error) bool {
	return true
}

// RetryNever always returns false, never allowing a retry for any error.
func RetryNever(_ error) bool {
	return false
}

// Delayer represents a policy for determining the delay duration
// before the next retry attempt.
type Delayer = func(attempt int) time.Duration

// ExponentialBackoff implements Delayer using Exponential Back-off strategy.
//
// Given an InitialDelay 'i', a MaxDelay 'M', a Coefficient 'c', and an attempt 'n'
// the delay is calculated as follows:
//
//	delay = min(M, (c^n) * i)
//
// ## Preconditions
// 1. Coefficient > 0
// 2. InitialDelay < MaxDelay
//
// If the preconditions are not met, behavior is undefined.
type ExponentialBackoff struct {
	InitialDelay time.Duration
	MaxDelay     time.Duration
	Coefficient  float64
}

// DefaultBackoffCoefficient is the default exponential back-off coefficient used
const DefaultBackoffCoefficient = 2.0

// Delay implements Delayer for ExponentialBackoff
func (e ExponentialBackoff) Delay(attempt int) time.Duration {
	if attempt == 0 {
		return 0
	}
	if e.Coefficient == 0 {
		e.Coefficient = DefaultBackoffCoefficient
	}
	dur := time.Duration(float64(e.InitialDelay) * math.Pow(e.Coefficient, float64(attempt)))
	if dur > e.MaxDelay {
		return e.MaxDelay
	}
	return dur
}

// Rand is a function that returns a value between [0.0,1.0) with uniformly distributed probability
type Rand = func() float64

func randomDurationBetween(rnd Rand, start, end time.Duration) time.Duration {
	r := rnd()
	return start + time.Duration(r*float64(end-start))
}

// FullJitter augments the given Delayer by adding full jitter to the returned duration.
// Full Jitter is defined as:
//
//	random_between( 0 , delay )
func FullJitter(rnd Rand, delayer Delayer) Delayer {
	return func(attempt int) time.Duration {
		return randomDurationBetween(rnd, 0, delayer(attempt))
	}
}

// EqualJitter augments the given Delayer by adding jitter to 1/2 the delay.
// EqualJitter is defined as:
//
//	delay/2 + random_between( 0 , delay/2 )
func EqualJitter(rnd Rand, delayer Delayer) Delayer {
	return func(attempt int) time.Duration {
		d := delayer(attempt)
		return d/2 + randomDurationBetween(rnd, 0, d/2)
	}
}

// DefaultDecorrelatedScale is the default scale factor for DecorrelatedJitter
const DefaultDecorrelatedScale = 3.0

// DecorrelatedJitter returns a Delayer using a Decorrelated Jitter algorithm.
// See: https://www.awsarchitectureblog.com/2015/03/backoff.html
//
// This is implemented as:
//
//	sleep = base
//	sleep = min(cap, random_between(0,sleep * scale))
//
// **NOTE**: This delayer has internal state, and is therefore not safe to be called concurrently.
// A new DecorrelatedJitter should be used for each goroutine.
//
// ## Preconditions
// 1. scale should be > 0
// 2. base < cap
//
// If the preconditions are not met, behavior is undefined.
func DecorrelatedJitter(rnd Rand, base time.Duration, cap time.Duration, scale float64) Delayer {
	sleep := base
	if scale <= 0 {
		scale = DefaultDecorrelatedScale
	}
	return func(attempt int) time.Duration {
		sleep = randomDurationBetween(rnd, base, time.Duration(float64(sleep)*scale))
		if sleep > cap {
			return cap
		}
		return sleep
	}
}

// Duration returns  Delayer that always returns the same Duration.
func Duration(dur time.Duration) Delayer {
	return func(attempt int) time.Duration {
		return dur
	}
}
