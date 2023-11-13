// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmetric

import "sync"

type GlobalCallbackItem struct {
	Callback GlobalCallback
	Metrics  []Metric
}

var (
	globalCallbackMu    sync.Mutex
	globalCallbackItems = make([]GlobalCallbackItem, 0)
)

func RegisterCallback(callback GlobalCallback, metrics ...Metric) error {
	globalCallbackMu.Lock()
	defer globalCallbackMu.Unlock()
	globalCallbackItems = append(globalCallbackItems, GlobalCallbackItem{
		Callback: callback,
		Metrics:  metrics,
	})
	// It removes the callback of metric if global callback defined on it.
	for _, m := range metrics {
		if callbackRemover, ok := m.(CallbackRemover); ok {
			callbackRemover.RemoveCallback()
		}
	}
	return nil
}

func GetRegisteredCallbacks() []GlobalCallbackItem {
	globalCallbackMu.Lock()
	defer globalCallbackMu.Unlock()
	items := globalCallbackItems
	globalCallbackItems = globalCallbackItems[:0]
	return items
}
