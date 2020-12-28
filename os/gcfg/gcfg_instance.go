// Copyright GoFrame Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcfg

import (
	"fmt"
	"github.com/gogf/gf/container/gmap"
)

const (
	// Default group name for instance usage.
	DefaultName = "default"
)

var (
	// Instances map containing configuration instances.
	instances = gmap.NewStrAnyMap(true)
)

// Instance returns an instance of Config with default settings.
// The parameter <name> is the name for the instance. But very note that, if the file "name.toml"
// exists in the configuration directory, it then sets it as the default configuration file. The
// toml file type is the default configuration file type.
func Instance(name ...string) *Config {
	key := DefaultName
	if len(name) > 0 && name[0] != "" {
		key = name[0]
	}
	return instances.GetOrSetFuncLock(key, func() interface{} {
		c := New()
		for _, fileType := range supportedFileTypes {
			if file := fmt.Sprintf(`%s.%s`, key, fileType); c.Available(file) {
				c.SetFileName(file)
				break
			}
		}
		return c
	}).(*Config)
}
