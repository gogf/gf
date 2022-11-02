// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gredis provides convenient client for redis server.
//
// Redis Client.
//
// Redis Commands Official: https://redis.io/commands
//
// Redis Chinese Documentation: http://redisdoc.com/
package gredis

import (
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

// AdapterFunc is the function creating redis adapter.
type AdapterFunc func(config *Config) Adapter

var (
	// defaultAdapterFunc is the default adapter function creating redis adapter.
	defaultAdapterFunc AdapterFunc = func(config *Config) Adapter {
		return nil
	}
)

// New creates and returns a redis client.
// It creates a default redis adapter of go-redis.
func New(config ...*Config) (*Redis, error) {
	var (
		usedConfig  *Config
		usedAdapter Adapter
	)
	if len(config) > 0 && config[0] != nil {
		// Redis client with go redis implements adapter from given configuration.
		usedConfig = config[0]
		usedAdapter = defaultAdapterFunc(config[0])
	} else if configFromGlobal, ok := GetConfig(); ok {
		// Redis client with go redis implements adapter from package configuration.
		usedConfig = configFromGlobal
		usedAdapter = defaultAdapterFunc(configFromGlobal)
	}
	if usedConfig == nil {
		return nil, gerror.NewCode(
			gcode.CodeInvalidConfiguration,
			`no configuration found for creating Redis client`,
		)
	}
	if usedAdapter == nil {
		return nil, gerror.NewCode(
			gcode.CodeNecessaryPackageNotImport,
			errorNilAdapter,
		)
	}
	redis := &Redis{
		config:       config[0],
		localAdapter: defaultAdapterFunc(config[0]),
	}
	return redis.initGroup(), nil
}

// NewWithAdapter creates and returns a redis client with given adapter.
func NewWithAdapter(adapter Adapter) (*Redis, error) {
	if adapter == nil {
		return nil, gerror.NewCodef(gcode.CodeInvalidParameter, `adapter cannot be nil`)
	}
	redis := &Redis{localAdapter: adapter}
	return redis.initGroup(), nil
}

// RegisterAdapterFunc registers default function creating redis adapter.
func RegisterAdapterFunc(adapterFunc AdapterFunc) {
	defaultAdapterFunc = adapterFunc
}
