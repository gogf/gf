// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmetric

// noopHistogramPerformer is an implementer for interface HistogramPerformer with no truly operations.
type noopHistogramPerformer struct{}

// newNoopHistogramPerformer creates and returns a HistogramPerformer with no truly operations.
func newNoopHistogramPerformer() HistogramPerformer {
	return noopHistogramPerformer{}
}

// Record adds a single value to the histogram. The value is usually positive or zero.
func (noopHistogramPerformer) Record(increment float64, option ...Option) {}
