// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gcfg provides reading, caching and managing for configuration.
package gcfg

import (
	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/internal/intlog"
	"github.com/gogf/gf/os/gcmd"
)

// Config is the configuration manager.
type Config struct {
	defaultName   string           // Default configuration file name.
	searchPaths   *garray.StrArray // Searching path array.
	jsonMap       *gmap.StrAnyMap  // The pared JSON objects for configuration files.
	violenceCheck bool             // Whether do violence check in value index searching. It affects the performance when set true(false in default).
}

const (
	DefaultName                = "config"             // DefaultName is the default group name for instance usage.
	DefaultConfigFile          = "config.toml"        // DefaultConfigFile is the default configuration file name.
	commandEnvKeyForFile       = "gf.gcfg.file"       // commandEnvKeyForFile is the configuration key for command argument or environment configuring file name.
	commandEnvKeyForPath       = "gf.gcfg.path"       // commandEnvKeyForPath is the configuration key for command argument or environment configuring directory path.
	commandEnvKeyForErrorPrint = "gf.gcfg.errorprint" // commandEnvKeyForErrorPrint is used to specify the key controlling error printing to stdout.
)

var (
	supportedFileTypes     = []string{"toml", "yaml", "yml", "json", "ini", "xml"}         // All supported file types suffixes.
	resourceTryFiles       = []string{"", "/", "config/", "config", "/config", "/config/"} // Prefix array for trying searching in resource manager.
	instances              = gmap.NewStrAnyMap(true)                                       // Instances map containing configuration instances.
	customConfigContentMap = gmap.NewStrStrMap(true)                                       // Customized configuration content.
)

// SetContent sets customized configuration content for specified `file`.
// The `file` is unnecessary param, default is DefaultConfigFile.
func SetContent(content string, file ...string) {
	name := DefaultConfigFile
	if len(file) > 0 {
		name = file[0]
	}
	// Clear file cache for instances which cached `name`.
	instances.LockFunc(func(m map[string]interface{}) {
		if customConfigContentMap.Contains(name) {
			for _, v := range m {
				v.(*Config).jsonMap.Remove(name)
			}
		}
		customConfigContentMap.Set(name, content)
	})
}

// GetContent returns customized configuration content for specified `file`.
// The `file` is unnecessary param, default is DefaultConfigFile.
func GetContent(file ...string) string {
	name := DefaultConfigFile
	if len(file) > 0 {
		name = file[0]
	}
	return customConfigContentMap.Get(name)
}

// RemoveContent removes the global configuration with specified `file`.
// If `name` is not passed, it removes configuration of the default group name.
func RemoveContent(file ...string) {
	name := DefaultConfigFile
	if len(file) > 0 {
		name = file[0]
	}
	// Clear file cache for instances which cached `name`.
	instances.LockFunc(func(m map[string]interface{}) {
		if customConfigContentMap.Contains(name) {
			for _, v := range m {
				v.(*Config).jsonMap.Remove(name)
			}
			customConfigContentMap.Remove(name)
		}
	})

	intlog.Printf(`RemoveContent: %s`, name)
}

// ClearContent removes all global configuration contents.
func ClearContent() {
	customConfigContentMap.Clear()
	// Clear cache for all instances.
	instances.LockFunc(func(m map[string]interface{}) {
		for _, v := range m {
			v.(*Config).jsonMap.Clear()
		}
	})

	intlog.Print(`RemoveConfig`)
}

// errorPrint checks whether printing error to stdout.
func errorPrint() bool {
	return gcmd.GetOptWithEnv(commandEnvKeyForErrorPrint, true).Bool()
}
