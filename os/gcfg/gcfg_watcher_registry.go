// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcfg

import (
	"context"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/internal/intlog"
)

// WatcherRegistry is a helper type for managing configuration watchers.
// It provides a unified implementation of watcher management to avoid code duplication
// across different adapter implementations.
type WatcherRegistry struct {
	watchers *gmap.KVMap[string, WatcherFunc] // Watchers map storing watcher callbacks.
}

// NewWatcherRegistry creates and returns a new WatcherRegistry instance.
func NewWatcherRegistry() *WatcherRegistry {
	return &WatcherRegistry{
		watchers: gmap.NewKVMap[string, WatcherFunc](true),
	}
}

// IsWatching checks whether the watcher with the specified name is registered.
func (r *WatcherRegistry) IsWatching(name string) bool {
	return r.watchers.Contains(name)
}

// Add adds a watcher with the specified name and callback function.
func (r *WatcherRegistry) Add(name string, fn WatcherFunc) {
	r.watchers.Set(name, fn)
}

// Remove removes the watcher with the specified name.
func (r *WatcherRegistry) Remove(name string) {
	r.watchers.Remove(name)
}

// GetNames returns all watcher names.
func (r *WatcherRegistry) GetNames() []string {
	return r.watchers.Keys()
}

// Notify notifies all registered watchers by calling their callback functions.
// Each callback is executed in a separate goroutine with panic recovery to prevent
// one watcher's panic from affecting others.
func (r *WatcherRegistry) Notify(ctx context.Context) {
	r.watchers.Iterator(func(k string, fn WatcherFunc) bool {
		go func(k string, fn WatcherFunc, ctx context.Context) {
			defer func() {
				if r := recover(); r != nil {
					intlog.Errorf(ctx, "watcher %s panic: %v", k, r)
				}
			}()
			fn(ctx)
		}(k, fn, ctx)
		return true
	})
}
