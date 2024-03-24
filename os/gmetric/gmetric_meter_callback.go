// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmetric

// CallbackItem is the global callback item registered.
type CallbackItem struct {
	Callback    Callback           // Global callback.
	Metrics     []ObservableMetric // Callback on certain metrics.
	MeterOption MeterOption        // MeterOption is the option that the meter holds.
	Provider    Provider           // Provider is the Provider that the callback item is bound to.
}

var (
	// Registered callbacks.
	globalCallbackItems = make([]CallbackItem, 0)
)

// RegisterCallback registers callback on certain metrics.
// A callback is bound to certain component and version, it is called when the associated metrics are read.
// Multiple callbacks on the same component and version will be called by their registered sequence.
func (meter *localMeter) RegisterCallback(callback Callback, observableMetrics ...ObservableMetric) error {
	if len(observableMetrics) == 0 {
		return nil
	}
	globalCallbackItems = append(globalCallbackItems, CallbackItem{
		Callback:    callback,
		Metrics:     observableMetrics,
		MeterOption: meter.MeterOption,
	})
	return nil
}

// MustRegisterCallback performs as RegisterCallback, but it panics if any error occurs.
func (meter *localMeter) MustRegisterCallback(callback Callback, observableMetrics ...ObservableMetric) {
	err := meter.RegisterCallback(callback, observableMetrics...)
	if err != nil {
		panic(err)
	}
}

// GetRegisteredCallbacks retrieves and returns the registered global callbacks.
// It truncates the callback slice is the callbacks are returned.
func GetRegisteredCallbacks() []CallbackItem {
	items := globalCallbackItems
	globalCallbackItems = globalCallbackItems[:0]
	return items
}
