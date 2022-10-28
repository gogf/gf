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
	if len(config) > 0 && config[0] != nil {
		// Redis client with go redis implements adapter from given configuration.
		redis := &Redis{
			adapter: defaultAdapterFunc(config[0]),
			config:  config[0],
		}
		return redis.initGroup(), nil
	}
	// Redis client with go redis implements adapter from package configuration.
	if configFromGlobal, ok := GetConfig(); ok {
		redis := &Redis{
			adapter: defaultAdapterFunc(configFromGlobal),
			config:  configFromGlobal,
		}
		return redis.initGroup(), nil
	}
	// Redis client with empty adapter.
	redis := &Redis{}
	return redis.initGroup(), nil
}

// NewWithAdapter creates and returns a redis client with given adapter.
func NewWithAdapter(adapter Adapter) *Redis {
	redis := &Redis{adapter: adapter}
	return redis.initGroup()
}

// RegisterAdapterFunc registers default function creating redis adapter.
func RegisterAdapterFunc(adapterFunc AdapterFunc) {
	defaultAdapterFunc = adapterFunc
}
