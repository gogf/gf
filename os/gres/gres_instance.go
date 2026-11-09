// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gres

import "github.com/gogf/gf/v2/container/gmap"

const (
	// DefaultName default group name for instance usage.
	DefaultName = "default"
)

var (
	// checker checks whether the value is nil.
	checker = func(v *Resource) bool { return v == nil }
	// Instances map.
	instances = gmap.NewKVMapWithChecker[string, *Resource](checker, true)
)

// Instance returns an instance of Resource.
// The parameter `name` is the name for the instance.
func Instance(name ...string) *Resource {
	key := DefaultName
	if len(name) > 0 && name[0] != "" {
		key = name[0]
	}
	return instances.GetOrSetFuncLock(key, New)
}
