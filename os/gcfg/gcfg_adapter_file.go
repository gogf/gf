// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcfg

import (
	"context"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/command"
	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gfsnotify"
	"github.com/gogf/gf/v2/os/gres"
	"github.com/gogf/gf/v2/util/gmode"
)

type AdapterFile struct {
	defaultName   string           // Default configuration file name.
	searchPaths   *garray.StrArray // Searching path array.
	jsonMap       *gmap.StrAnyMap  // The pared JSON objects for configuration files.
	violenceCheck bool             // Whether it does violence check in value index searching. It affects the performance when set true(false in default).
}

const (
	commandEnvKeyForFile = "gf.gcfg.file" // commandEnvKeyForFile is the configuration key for command argument or environment configuring file name.
	commandEnvKeyForPath = "gf.gcfg.path" // commandEnvKeyForPath is the configuration key for command argument or environment configuring directory path.
)

var (
	supportedFileTypes     = []string{"toml", "yaml", "yml", "json", "ini", "xml"} // All supported file types suffixes.
	localInstances         = gmap.NewStrAnyMap(true)                               // Instances map containing configuration instances.
	customConfigContentMap = gmap.NewStrStrMap(true)                               // Customized configuration content.

	// Prefix array for trying searching in resource manager.
	resourceTryFolders = []string{
		"", "/", "config/", "config", "/config", "/config/",
		"manifest/config/", "manifest/config", "/manifest/config", "/manifest/config/",
	}

	// Prefix array for trying searching in local system.
	localSystemTryFolders = []string{"", "config/", "manifest/config"}
)

// NewAdapterFile returns a new configuration management object.
// The parameter `file` specifies the default configuration file name for reading.
func NewAdapterFile(file ...string) (*AdapterFile, error) {
	var (
		err  error
		name = DefaultConfigFileName
	)
	if len(file) > 0 {
		name = file[0]
	} else {
		// Custom default configuration file name from command line or environment.
		if customFile := command.GetOptWithEnv(commandEnvKeyForFile); customFile != "" {
			name = customFile
		}
	}
	config := &AdapterFile{
		defaultName: name,
		searchPaths: garray.NewStrArray(true),
		jsonMap:     gmap.NewStrAnyMap(true),
	}
	// Customized dir path from env/cmd.
	if customPath := command.GetOptWithEnv(commandEnvKeyForPath); customPath != "" {
		if gfile.Exists(customPath) {
			if err = config.SetPath(customPath); err != nil {
				return nil, err
			}
		} else {
			return nil, gerror.Newf(`configuration directory path "%s" does not exist`, customPath)
		}
	} else {
		// ================================================================================
		// Automatic searching directories.
		// It does not affect adapter object cresting if these directories do not exist.
		// ================================================================================

		// Dir path of working dir.
		if err = config.AddPath(gfile.Pwd()); err != nil {
			intlog.Errorf(context.TODO(), `%+v`, err)
		}

		// Dir path of main package.
		if mainPath := gfile.MainPkgPath(); mainPath != "" && gfile.Exists(mainPath) {
			if err = config.AddPath(mainPath); err != nil {
				intlog.Errorf(context.TODO(), `%+v`, err)
			}
		}

		// Dir path of binary.
		if selfPath := gfile.SelfDir(); selfPath != "" && gfile.Exists(selfPath) {
			if err = config.AddPath(selfPath); err != nil {
				intlog.Errorf(context.TODO(), `%+v`, err)
			}
		}
	}
	return config, nil
}

// SetViolenceCheck sets whether to perform hierarchical conflict checking.
// This feature needs to be enabled when there is a level symbol in the key name.
// It is off in default.
//
// Note that, turning on this feature is quite expensive, and it is not recommended
// allowing separators in the key names. It is best to avoid this on the application side.
func (c *AdapterFile) SetViolenceCheck(check bool) {
	c.violenceCheck = check
	c.Clear()
}

// SetFileName sets the default configuration file name.
func (c *AdapterFile) SetFileName(name string) {
	c.defaultName = name
}

// GetFileName returns the default configuration file name.
func (c *AdapterFile) GetFileName() string {
	return c.defaultName
}

// Get retrieves and returns value by specified `pattern`.
// It returns all values of current Json object if `pattern` is given empty or string ".".
// It returns nil if no value found by `pattern`.
//
// We can also access slice item by its index number in `pattern` like:
// "list.10", "array.0.name", "array.0.1.id".
//
// It returns a default value specified by `def` if value for `pattern` is not found.
func (c *AdapterFile) Get(ctx context.Context, pattern string) (value interface{}, err error) {
	j, err := c.getJson()
	if err != nil {
		return nil, err
	}
	if j != nil {
		return j.Get(pattern).Val(), nil
	}
	return nil, nil
}

// Set sets value with specified `pattern`.
// It supports hierarchical data access by char separator, which is '.' in default.
// It is commonly used for updates certain configuration value in runtime.
// Note that, it is not recommended using `Set` configuration at runtime as the configuration would be
// automatically refreshed if underlying configuration file changed.
func (c *AdapterFile) Set(pattern string, value interface{}) error {
	j, err := c.getJson()
	if err != nil {
		return err
	}
	if j != nil {
		return j.Set(pattern, value)
	}
	return nil
}

// Data retrieves and returns all configuration data as map type.
func (c *AdapterFile) Data(ctx context.Context) (data map[string]interface{}, err error) {
	j, err := c.getJson()
	if err != nil {
		return nil, err
	}
	if j != nil {
		return j.Var().Map(), nil
	}
	return nil, nil
}

// MustGet acts as function Get, but it panics if error occurs.
func (c *AdapterFile) MustGet(ctx context.Context, pattern string) *gvar.Var {
	v, err := c.Get(ctx, pattern)
	if err != nil {
		panic(err)
	}
	return gvar.New(v)
}

// Clear removes all parsed configuration files content cache,
// which will force reload configuration content from file.
func (c *AdapterFile) Clear() {
	c.jsonMap.Clear()
}

// Dump prints current Json object with more manually readable.
func (c *AdapterFile) Dump() {
	if j, _ := c.getJson(); j != nil {
		j.Dump()
	}
}

// Available checks and returns whether configuration of given `file` is available.
func (c *AdapterFile) Available(ctx context.Context, fileName ...string) bool {
	var (
		usedFileName string
	)
	if len(fileName) > 0 && fileName[0] != "" {
		usedFileName = fileName[0]
	} else {
		usedFileName = c.defaultName
	}
	// Custom configuration content exists.
	if c.GetContent(usedFileName) != "" {
		return true
	}
	// Configuration file exists in system path.
	if path, _ := c.GetFilePath(usedFileName); path != "" {
		return true
	}
	return false
}

// autoCheckAndAddMainPkgPathToSearchPaths automatically checks and adds directory path of package main
// to the searching path list if it's currently in development environment.
func (c *AdapterFile) autoCheckAndAddMainPkgPathToSearchPaths() {
	if gmode.IsDevelop() {
		mainPkgPath := gfile.MainPkgPath()
		if mainPkgPath != "" {
			if !c.searchPaths.Contains(mainPkgPath) {
				c.searchPaths.Append(mainPkgPath)
			}
		}
	}
}

// getJson returns a *gjson.Json object for the specified `file` content.
// It would print error if file reading fails. It returns nil if any error occurs.
func (c *AdapterFile) getJson(fileName ...string) (configJson *gjson.Json, err error) {
	var (
		usedFileName = c.defaultName
	)
	if len(fileName) > 0 && fileName[0] != "" {
		usedFileName = fileName[0]
	} else {
		usedFileName = c.defaultName
	}
	// It uses json map to cache specified configuration file content.
	result := c.jsonMap.GetOrSetFuncLock(usedFileName, func() interface{} {
		var (
			content  string
			filePath string
		)
		// The configured content can be any kind of data type different from its file type.
		isFromConfigContent := true
		if content = c.GetContent(usedFileName); content == "" {
			isFromConfigContent = false
			filePath, err = c.GetFilePath(usedFileName)
			if err != nil {
				return nil
			}
			if filePath == "" {
				return nil
			}
			if file := gres.Get(filePath); file != nil {
				content = string(file.Content())
			} else {
				content = gfile.GetContents(filePath)
			}
		}
		// Note that the underlying configuration json object operations are concurrent safe.
		dataType := gfile.ExtName(usedFileName)
		if gjson.IsValidDataType(dataType) && !isFromConfigContent {
			configJson, err = gjson.LoadContentType(dataType, content, true)
		} else {
			configJson, err = gjson.LoadContent(content, true)
		}
		if err != nil {
			if filePath != "" {
				err = gerror.Wrapf(err, `load config file "%s" failed`, filePath)
			} else {
				err = gerror.Wrap(err, `load configuration failed`)
			}
			return nil
		}
		configJson.SetViolenceCheck(c.violenceCheck)
		// Add monitor for this configuration file,
		// any changes of this file will refresh its cache in Config object.
		if filePath != "" && !gres.Contains(filePath) {
			_, err = gfsnotify.Add(filePath, func(event *gfsnotify.Event) {
				c.jsonMap.Remove(usedFileName)
			}, ``)
			if err != nil {
				return nil
			}
		}
		return configJson
	})
	if result != nil {
		return result.(*gjson.Json), err
	}
	return
}
