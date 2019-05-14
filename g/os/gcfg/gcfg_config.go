// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcfg

import "github.com/gogf/gf/g/container/gmap"

var (
    // Customized configuration content.
    configs = gmap.NewStrStrMap()
)

// SetContent sets customized configuration content for specified <file>.
// The <file> is unnecessary param, default is DEFAULT_CONFIG_FILE.
func SetContent(content string, file ...string) {
    name := DEFAULT_CONFIG_FILE
    if len(file) > 0 {
        name = file[0]
    }
    // Clear file cache for instances which cached <name>.
    instances.LockFunc(func(m map[string]interface{}) {
    	if configs.Contains(name) {
		    for _, v := range m {
			    v.(*Config).jsons.Remove(name)
		    }
	    }
	    configs.Set(name, content)
    })
}

// GetContent returns customized configuration content for specified <file>.
// The <file> is unnecessary param, default is DEFAULT_CONFIG_FILE.
func GetContent(file ...string) string {
    name := DEFAULT_CONFIG_FILE
    if len(file) > 0 {
        name = file[0]
    }
    return configs.Get(name)
}

// RemoveConfig removes the global configuration with specified group.
// If <name> is not passed, it removes configuration of the default group name.
func RemoveConfig(file ...string) {
	name := DEFAULT_CONFIG_FILE
	if len(file) > 0 {
		name = file[0]
	}
	// Clear file cache for instances which cached <name>.
	instances.LockFunc(func(m map[string]interface{}) {
		if configs.Contains(name) {
			for _, v := range m {
				v.(*Config).jsons.Remove(name)
			}
			configs.Remove(name)
		}
	})
}

// ClearContent removes all global configuration contents.
func ClearContent() {
    configs.Clear()
	// Clear cache for all instances.
	instances.LockFunc(func(m map[string]interface{}) {
		for _, v := range m {
			v.(*Config).jsons.Clear()
		}
	})
}