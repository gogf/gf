// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmetric

// noopObservableCounterPerformer is an implementer for interface ObservableCounterPerformer with no truly operations.
type noopObservableCounterPerformer struct{}

// newNoopObservableCounterPerformer creates and returns a ObservableCounterPerformer with no truly operations.
func newNoopObservableCounterPerformer() ObservableCounterPerformer {
	return noopObservableCounterPerformer{}
}

func (noopObservableCounterPerformer) observable() {}
