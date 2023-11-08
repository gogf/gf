// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmetric

// noopGaugePerformer is an implementer for interface GaugePerformer with no truly operations.
type noopGaugePerformer struct{}

// newNoopGaugePerformer creates and returns a GaugePerformer with no truly operations.
func newNoopGaugePerformer() GaugePerformer {
	return noopGaugePerformer{}
}

// Set sets the Gauge to an arbitrary value.
func (noopGaugePerformer) Set(value float64, option ...Option) {}

// Inc increments the Gauge by 1. Use Add to increment it by arbitrary values.
func (noopGaugePerformer) Inc(option ...Option) {}

// Dec decrements the Gauge by 1. Use Sub to decrement it by arbitrary values.
func (noopGaugePerformer) Dec(option ...Option) {}

// Add adds the given value to the Gauge. (The value can be negative,
// resulting in a decrease of the Gauge.)
func (noopGaugePerformer) Add(increment float64, option ...Option) {}

// Sub subtracts the given value from the Gauge. (The value can be
// negative, resulting in an increase of the Gauge.)
func (noopGaugePerformer) Sub(decrement float64, option ...Option) {}
