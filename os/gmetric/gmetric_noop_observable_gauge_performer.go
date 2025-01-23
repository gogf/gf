// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmetric

// noopObservableGaugePerformer is an implementer for interface ObservableGaugePerformer with no truly operations.
type noopObservableGaugePerformer struct{}

// newNoopObservableGaugePerformer creates and returns a ObservableGaugePerformer with no truly operations.
func newNoopObservableGaugePerformer() ObservableGaugePerformer {
	return noopObservableGaugePerformer{}
}

func (noopObservableGaugePerformer) observable() {}
