// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmetric

// noopCounterPerformer is an implementer for interface CounterPerformer with no truly operations.
type noopCounterPerformer struct{}

// newNoopCounterPerformer creates and returns a CounterPerformer with no truly operations.
func newNoopCounterPerformer() CounterPerformer {
	return noopCounterPerformer{}
}

// Inc increments the counter by 1. Use Add to increment it by arbitrary
// non-negative values.
func (noopCounterPerformer) Inc(option ...Option) {}

// Add adds the given value to the counter. It panics if the value is < 0.
func (noopCounterPerformer) Add(increment float64, option ...Option) {}

func (noopCounterPerformer) canBeCallback() {}
