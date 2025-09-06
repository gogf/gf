// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gi18n

import "github.com/gogf/gf/v2/container/gmap"

const (
	// DefaultName is the default group name for instance usage.
	DefaultName = "default"
)

var (
	// instances is the instances map for management
	// for multiple i18n instance by name.
	instances = gmap.NewStrAnyMap(true)
)

// Instance returns an instance of Resource.
// The parameter `name` is the name for the instance.
func Instance(name ...string) *Manager {
	key := DefaultName
	if len(name) > 0 && name[0] != "" {
		key = name[0]
	}
	return instances.GetOrSetFuncLock(key, func() any {
		return New()
	}).(*Manager)
}
