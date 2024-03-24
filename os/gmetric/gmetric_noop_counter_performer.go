// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmetric

import "context"

// noopCounterPerformer is an implementer for interface CounterPerformer with no truly operations.
type noopCounterPerformer struct{}

// newNoopCounterPerformer creates and returns a CounterPerformer with no truly operations.
func newNoopCounterPerformer() CounterPerformer {
	return noopCounterPerformer{}
}

// Inc increments the counter by 1.
func (noopCounterPerformer) Inc(ctx context.Context, option ...Option) {}

// Add adds the given value to the counter. It panics if the value is < 0.
func (noopCounterPerformer) Add(ctx context.Context, increment float64, option ...Option) {}
