// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gcfg provides reading, caching and managing for configuration.
package gcfg

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/internal/command"
	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/internal/utils"
	"github.com/gogf/gf/v2/os/genv"
)

// Config is the configuration management object.
type Config struct {
	adapter Adapter
	dataMap *gmap.StrAnyMap
}

const (
	DefaultName = "config" // DefaultName is the default group name for instance usage.
)

// New creates and returns a Config object with default adapter of AdapterFile.
func New() (*Config, error) {
	adapterFile, err := NewAdapterFile()
	if err != nil {
		return nil, err
	}
	return &Config{
		adapter: adapterFile,
		dataMap: gmap.NewStrAnyMap(true),
	}, nil
}

// NewWithAdapter creates and returns a Config object with given adapter.
func NewWithAdapter(adapter Adapter) *Config {
	return &Config{
		adapter: adapter,
		dataMap: gmap.NewStrAnyMap(true),
	}
}

// Instance returns an instance of Config with default settings.
// The parameter `name` is the name for the instance. But very note that, if the file "name.toml"
// exists in the configuration directory, it then sets it as the default configuration file. The
// toml file type is the default configuration file type.
func Instance(name ...string) *Config {
	var (
		ctx = context.TODO()
		key = DefaultName
	)
	if len(name) > 0 && name[0] != "" {
		key = name[0]
	}
	return localInstances.GetOrSetFuncLock(key, func() interface{} {
		adapter, err := NewAdapterFile()
		if err != nil {
			intlog.Error(context.Background(), err)
			return nil
		}
		// If it's not using default configuration or its configuration file is not available,
		// it searches the possible configuration file according to the name and all supported
		// file types.
		if key != DefaultName || !adapter.Available(ctx) {
			for _, fileType := range supportedFileTypes {
				if file := fmt.Sprintf(`%s.%s`, key, fileType); adapter.Available(ctx, file) {
					adapter.SetFileName(file)
					break
				}
			}
		}
		return NewWithAdapter(adapter)
	}).(*Config)
}

// SetAdapter sets the adapter of current Config object.
func (c *Config) SetAdapter(adapter Adapter) {
	c.adapter = adapter
}

// GetAdapter returns the adapter of current Config object.
func (c *Config) GetAdapter() Adapter {
	return c.adapter
}

// Available checks and returns the configuration service is available.
// The optional parameter `pattern` specifies certain configuration resource.
//
// It returns true if configuration file is present in default AdapterFile, or else false.
// Note that this function does not return error as it just does simply check for backend configuration service.
func (c *Config) Available(ctx context.Context, resource ...string) (ok bool) {
	return c.adapter.Available(ctx, resource...)
}

// Set sets value with specified `pattern`.
// It supports hierarchical data access by char separator, which is '.' in default.
// It is commonly used for updates certain configuration value in runtime.
func (c *Config) Set(ctx context.Context, pattern string, value interface{}) {
	c.dataMap.Set(pattern, value)
}

// Get retrieves and returns value by specified `pattern`.
// It returns all values of current Json object if `pattern` is given empty or string ".".
// It returns nil if no value found by `pattern`.
//
// It returns a default value specified by `def` if value for `pattern` is not found.
func (c *Config) Get(ctx context.Context, pattern string, def ...interface{}) (*gvar.Var, error) {
	var (
		err   error
		value interface{}
	)
	if value = c.dataMap.Get(pattern); value == nil {
		value, err = c.adapter.Get(ctx, pattern)
		if err != nil {
			return nil, err
		}
		if value == nil {
			if len(def) > 0 {
				return gvar.New(def[0]), nil
			}
			return nil, nil
		}
	}
	return gvar.New(value), nil
}

// GetWithEnv returns the configuration value specified by pattern `pattern`.
// If the configuration value does not exist, then it retrieves and returns the environment value specified by `key`.
// It returns the default value `def` if none of them exists.
//
// Fetching Rules: Environment arguments are in uppercase format, eg: GF_PACKAGE_VARIABLE.
func (c *Config) GetWithEnv(ctx context.Context, pattern string, def ...interface{}) (*gvar.Var, error) {
	value, err := c.Get(ctx, pattern)
	if err != nil {
		return nil, err
	}
	if value == nil {
		if v := genv.Get(utils.FormatEnvKey(pattern)); v != nil {
			return v, nil
		}
		if len(def) > 0 {
			return gvar.New(def[0]), nil
		}
		return nil, nil
	}
	return value, nil
}

// GetWithCmd returns the configuration value specified by pattern `pattern`.
// If the configuration value does not exist, then it retrieves and returns the command line option specified by `key`.
// It returns the default value `def` if none of them exists.
//
// Fetching Rules: Command line arguments are in lowercase format, eg: gf.package.variable.
func (c *Config) GetWithCmd(ctx context.Context, pattern string, def ...interface{}) (*gvar.Var, error) {
	value, err := c.Get(ctx, pattern)
	if err != nil {
		return nil, err
	}
	if value == nil {
		if v := command.GetOpt(utils.FormatCmdKey(pattern)); v != "" {
			return gvar.New(v), nil
		}
		if len(def) > 0 {
			return gvar.New(def[0]), nil
		}
		return nil, nil
	}
	return value, nil
}

// Data retrieves and returns all configuration data as map type.
func (c *Config) Data(ctx context.Context) (data map[string]interface{}, err error) {
	adapterData, err := c.adapter.Data(ctx)
	if err != nil {
		return nil, err
	}
	data = make(map[string]interface{})
	for k, v := range adapterData {
		data[k] = v
	}
	c.dataMap.Iterator(func(k string, v interface{}) bool {
		data[k] = v
		return true
	})
	return
}

// MustGet acts as function Get, but it panics if error occurs.
func (c *Config) MustGet(ctx context.Context, pattern string, def ...interface{}) *gvar.Var {
	v, err := c.Get(ctx, pattern, def...)
	if err != nil {
		panic(err)
	}
	if v == nil {
		return nil
	}
	return v
}

// MustGetWithEnv acts as function GetWithEnv, but it panics if error occurs.
func (c *Config) MustGetWithEnv(ctx context.Context, pattern string, def ...interface{}) *gvar.Var {
	v, err := c.GetWithEnv(ctx, pattern, def...)
	if err != nil {
		panic(err)
	}
	return v
}

// MustGetWithCmd acts as function GetWithCmd, but it panics if error occurs.
func (c *Config) MustGetWithCmd(ctx context.Context, pattern string, def ...interface{}) *gvar.Var {
	v, err := c.GetWithCmd(ctx, pattern, def...)
	if err != nil {
		panic(err)
	}
	return v
}

// MustData acts as function Data, but it panics if error occurs.
func (c *Config) MustData(ctx context.Context) map[string]interface{} {
	v, err := c.Data(ctx)
	if err != nil {
		panic(err)
	}
	return v
}
