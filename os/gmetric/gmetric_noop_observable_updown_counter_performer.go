// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmetric

// noopObservableUpDownCounterPerformer is an implementer for interface ObservableUpDownCounterPerformer
// with no truly operations.
type noopObservableUpDownCounterPerformer struct{}

// newNoopObservableUpDownCounterPerformer creates and returns a ObservableUpDownCounterPerformer
// with no truly operations.
func newNoopObservableUpDownCounterPerformer() ObservableUpDownCounterPerformer {
	return noopObservableUpDownCounterPerformer{}
}

func (noopObservableUpDownCounterPerformer) observable() {}
