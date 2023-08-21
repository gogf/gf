// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gmutex inherits and extends sync.Mutex and sync.RWMutex with more futures.
//
// Note that, it is refracted using stdlib mutex of package sync from GoFrame version v2.5.2.
package gmutex

// New creates and returns a new mutex.
// Deprecated: use Mutex or RWMutex instead.
func New() *RWMutex {
	return &RWMutex{}
}
