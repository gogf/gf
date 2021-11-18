// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcfg

import (
	"context"

	"github.com/gogf/gf/v2/internal/intlog"
)

// SetContent sets customized configuration content for specified `file`.
// The `file` is unnecessary param, default is DefaultConfigFile.
func (c *AdapterFile) SetContent(content string, file ...string) {
	name := DefaultConfigFile
	if len(file) > 0 {
		name = file[0]
	}
	// Clear file cache for instances which cached `name`.
	localInstances.LockFunc(func(m map[string]interface{}) {
		if customConfigContentMap.Contains(name) {
			for _, v := range m {
				if configInstance, ok := v.(*Config); ok {
					if fileConfig, ok := configInstance.GetAdapter().(*AdapterFile); ok {
						fileConfig.jsonMap.Remove(name)
					}
				}
			}
		}
		customConfigContentMap.Set(name, content)
	})
}

// GetContent returns customized configuration content for specified `file`.
// The `file` is unnecessary param, default is DefaultConfigFile.
func (c *AdapterFile) GetContent(file ...string) string {
	name := DefaultConfigFile
	if len(file) > 0 {
		name = file[0]
	}
	return customConfigContentMap.Get(name)
}

// RemoveContent removes the global configuration with specified `file`.
// If `name` is not passed, it removes configuration of the default group name.
func (c *AdapterFile) RemoveContent(file ...string) {
	name := DefaultConfigFile
	if len(file) > 0 {
		name = file[0]
	}
	// Clear file cache for instances which cached `name`.
	localInstances.LockFunc(func(m map[string]interface{}) {
		if customConfigContentMap.Contains(name) {
			for _, v := range m {
				if configInstance, ok := v.(*Config); ok {
					if fileConfig, ok := configInstance.GetAdapter().(*AdapterFile); ok {
						fileConfig.jsonMap.Remove(name)
					}
				}
			}
			customConfigContentMap.Remove(name)
		}
	})

	intlog.Printf(context.TODO(), `RemoveContent: %s`, name)
}

// ClearContent removes all global configuration contents.
func (c *AdapterFile) ClearContent() {
	customConfigContentMap.Clear()
	// Clear cache for all instances.
	localInstances.LockFunc(func(m map[string]interface{}) {
		for _, v := range m {
			if configInstance, ok := v.(*Config); ok {
				if fileConfig, ok := configInstance.GetAdapter().(*AdapterFile); ok {
					fileConfig.jsonMap.Clear()
				}
			}
		}
	})
	intlog.Print(context.TODO(), `RemoveConfig`)
}
