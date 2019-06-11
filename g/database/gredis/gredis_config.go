// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gredis

import "github.com/gogf/gf/g/container/gmap"

const (
    // Default configuration group name.
    DEFAULT_GROUP_NAME = "default"
)
var (
    // Configuration groups.
    configs = gmap.NewStrAnyMap()
)

// SetConfig sets the global configuration for specified group.
// If <name> is not passed, it sets configuration for the default group name.
func SetConfig(config Config, name...string) {
    group := DEFAULT_GROUP_NAME
    if len(name) > 0 {
        group = name[0]
    }
    configs.Set(group, config)
    instances.Remove(group)
}

// GetConfig returns the global configuration with specified group name.
// If <name> is not passed, it returns configuration of the default group name.
func GetConfig(name...string) (config Config, ok bool) {
    group := DEFAULT_GROUP_NAME
    if len(name) > 0 {
        group = name[0]
    }
    if v := configs.Get(group); v != nil {
        return v.(Config), true
    }
    return Config{}, false
}

// RemoveConfig removes the global configuration with specified group.
// If <name> is not passed, it removes configuration of the default group name.
func RemoveConfig(name...string) {
    group := DEFAULT_GROUP_NAME
    if len(name) > 0 {
        group = name[0]
    }
    configs.Remove(group)
    instances.Remove(group)
}

// ClearConfig removes all configurations and instances of redis.
func ClearConfig() {
    configs.Clear()
    instances.Clear()
}


