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
func (a *AdapterFile) SetContent(content string, fileNameOrPath ...string) {
	var usedFileNameOrPath = DefaultConfigFileName
	if len(fileNameOrPath) > 0 {
		usedFileNameOrPath = fileNameOrPath[0]
	}
	// Clear file cache for instances which cached `name`.
	localInstances.LockFunc(func(m map[string]interface{}) {
		if customConfigContentMap.Contains(usedFileNameOrPath) {
			for _, v := range m {
				if configInstance, ok := v.(*Config); ok {
					if fileConfig, ok := configInstance.GetAdapter().(*AdapterFile); ok {
						fileConfig.jsonMap.Remove(usedFileNameOrPath)
					}
				}
			}
		}
		customConfigContentMap.Set(usedFileNameOrPath, content)
	})
}

// GetContent returns customized configuration content for specified `file`.
// The `file` is unnecessary param, default is DefaultConfigFile.
func (a *AdapterFile) GetContent(fileNameOrPath ...string) string {
	var usedFileNameOrPath = DefaultConfigFileName
	if len(fileNameOrPath) > 0 {
		usedFileNameOrPath = fileNameOrPath[0]
	}
	return customConfigContentMap.Get(usedFileNameOrPath)
}

// RemoveContent removes the global configuration with specified `file`.
// If `name` is not passed, it removes configuration of the default group name.
func (a *AdapterFile) RemoveContent(fileNameOrPath ...string) {
	var usedFileNameOrPath = DefaultConfigFileName
	if len(fileNameOrPath) > 0 {
		usedFileNameOrPath = fileNameOrPath[0]
	}
	// Clear file cache for instances which cached `name`.
	localInstances.LockFunc(func(m map[string]interface{}) {
		if customConfigContentMap.Contains(usedFileNameOrPath) {
			for _, v := range m {
				if configInstance, ok := v.(*Config); ok {
					if fileConfig, ok := configInstance.GetAdapter().(*AdapterFile); ok {
						fileConfig.jsonMap.Remove(usedFileNameOrPath)
					}
				}
			}
			customConfigContentMap.Remove(usedFileNameOrPath)
		}
	})

	intlog.Printf(context.TODO(), `RemoveContent: %s`, usedFileNameOrPath)
}

// ClearContent removes all global configuration contents.
func (a *AdapterFile) ClearContent() {
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
