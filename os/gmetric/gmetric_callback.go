// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmetric

// GlobalCallbackItem is the global callback item registered.
type GlobalCallbackItem struct {
	Callback GlobalCallback // Global callback.
	Metrics  []Metric       // Callback on certain metrics.
}

var (
	// Registered callbacks.
	globalCallbackItems = make([]GlobalCallbackItem, 0)
)

// RegisterCallback registers global callback on certain metrics.
// A global callback is called when these metrics are being read.
func RegisterCallback(callback GlobalCallback, metrics ...Metric) {
	globalCallbackItems = append(globalCallbackItems, GlobalCallbackItem{
		Callback: callback,
		Metrics:  metrics,
	})
}

// GetRegisteredCallbacks retrieves and returns the registered global callbacks.
// It truncates the callback slice is the callbacks are returned.
func GetRegisteredCallbacks() []GlobalCallbackItem {
	items := globalCallbackItems
	globalCallbackItems = globalCallbackItems[:0]
	return items
}
