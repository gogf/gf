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

// New creates and returns a redis client.
// It creates a default redis adapter of go-redis.
func New(config ...*Config) (*Redis, error) {
	if len(config) > 0 {
		return &Redis{adapter: NewAdapterGoRedis(config[0])}, nil
	}
	configFromGlobal, ok := GetConfig()
	if !ok {
		return nil, gerror.NewCode(
			gcode.CodeMissingConfiguration,
			`configuration not found for creating Redis client`,
		)
	}
	return &Redis{adapter: NewAdapterGoRedis(configFromGlobal)}, nil
}

// NewWithAdapter creates and returns a redis client with given adapter.
func NewWithAdapter(adapter Adapter) *Redis {
	return &Redis{adapter: adapter}
}
