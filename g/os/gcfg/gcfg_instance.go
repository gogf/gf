// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcfg

import (
	"github.com/gogf/gf/g/container/gmap"
)

var (
    // Instances map.
    instances = gmap.NewStringInterfaceMap()
)

// Instance returns an instance of Config.
func Instance(file...string) *Config {
	configFile := DEFAULT_CONFIG_FILE
	if len(file) > 0 {
		configFile = file[0]
	}
	return instances.GetOrSetFuncLock(configFile, func() interface{} {
		return New(configFile)
	}).(*Config)
}
