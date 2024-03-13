// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmetric

// noopObservableMetric is an implementer for interface ObservableMetric with no truly operations.
type noopObservableMetric struct{}

// newNoopObservableMetric creates and returns a CounterPerformer with no truly operations.
func newNoopObservableMetric() ObservableMetric {
	return noopObservableMetric{}
}

func (m noopObservableMetric) observable() {}
