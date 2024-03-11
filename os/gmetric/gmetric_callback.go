// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmetric

// GlobalCallbackItem is the global callback item registered.
type GlobalCallbackItem struct {
	Callback GlobalCallback        // Global callback.
	Metrics  []CanBeCallbackMetric // Callback on certain metrics.
}

var (
	// Registered callbacks.
	globalCallbackItems = make([]GlobalCallbackItem, 0)
)

// RegisterCallback registers callback on certain metrics.
// A callback is bound to certain component and version, it is called when the associated metrics are read.
// Multiple callbacks on the same component and version will be called by their registered sequence.
func RegisterCallback(callback GlobalCallback, canBeCallbackMetrics ...CanBeCallbackMetric) error {
	if globalProvider != nil {
		return globalProvider.RegisterCallback(callback, canBeCallbackMetrics...)
	}

	if len(canBeCallbackMetrics) == 0 {
		return nil
	}
	globalCallbackItems = append(globalCallbackItems, GlobalCallbackItem{
		Callback: callback,
		Metrics:  canBeCallbackMetrics,
	})
	return nil
}

// MustRegisterCallback performs as RegisterCallback, but it panics if any error occurs.
func MustRegisterCallback(callback GlobalCallback, canBeCallbackMetrics ...CanBeCallbackMetric) {
	err := RegisterCallback(callback, canBeCallbackMetrics...)
	if err != nil {
		panic(err)
	}
}

// GetRegisteredCallbacks retrieves and returns the registered global callbacks.
// It truncates the callback slice is the callbacks are returned.
func GetRegisteredCallbacks() []GlobalCallbackItem {
	items := globalCallbackItems
	globalCallbackItems = globalCallbackItems[:0]
	return items
}
