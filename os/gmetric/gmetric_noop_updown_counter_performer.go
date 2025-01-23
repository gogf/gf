// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmetric

import "context"

// noopUpDownCounterPerformer is an implementer for interface CounterPerformer with no truly operations.
type noopUpDownCounterPerformer struct{}

// newNoopUpDownCounterPerformer creates and returns a CounterPerformer with no truly operations.
func newNoopUpDownCounterPerformer() UpDownCounterPerformer {
	return noopUpDownCounterPerformer{}
}

// Inc increments the counter by 1.
func (noopUpDownCounterPerformer) Inc(ctx context.Context, option ...Option) {}

// Dec decrements the counter by 1.
func (noopUpDownCounterPerformer) Dec(ctx context.Context, option ...Option) {}

// Add adds the given value to the counter.
func (noopUpDownCounterPerformer) Add(ctx context.Context, increment float64, option ...Option) {}
