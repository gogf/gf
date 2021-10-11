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
	"github.com/gogf/gf/v2/internal/intlog"
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
	key := DefaultName
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
		if key != DefaultName || !adapter.Available() {
			for _, fileType := range supportedFileTypes {
				if file := fmt.Sprintf(`%s.%s`, key, fileType); adapter.Available(file) {
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
		if value == nil && len(def) > 0 {
			return gvar.New(def[0]), nil
		}
	}
	return gvar.New(value), nil
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
	return gvar.New(v)
}

// MustData acts as function Data, but it panics if error occurs.
func (c *Config) MustData(ctx context.Context) map[string]interface{} {
	v, err := c.Data(ctx)
	if err != nil {
		panic(err)
	}
	return v
}
