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
	"context"
	"github.com/gogf/gf/container/gvar"
	"github.com/gogf/gf/errors/gcode"
	"github.com/gogf/gf/errors/gerror"
	"time"
)

// Redis client.
type Redis struct {
	adapter Adapter
}

// Config is redis configuration.
type Config struct {
	Host            string        `json:"host"`
	Port            int           `json:"port"`
	Db              int           `json:"db"`
	Pass            string        `json:"pass"`            // Password for AUTH.
	MaxIdle         int           `json:"maxIdle"`         // Maximum number of connections allowed to be idle (default is 10)
	MaxActive       int           `json:"maxActive"`       // Maximum number of connections limit (default is 0 means no limit).
	IdleTimeout     time.Duration `json:"idleTimeout"`     // Maximum idle time for connection (default is 10 seconds, not allowed to be set to 0)
	MaxConnLifetime time.Duration `json:"maxConnLifetime"` // Maximum lifetime of the connection (default is 30 seconds, not allowed to be set to 0)
	ConnectTimeout  time.Duration `json:"connectTimeout"`  // Dial connection timeout.
	TLS             bool          `json:"tls"`             // Specifies the config to use when a TLS connection is dialed.
	TLSSkipVerify   bool          `json:"tlsSkipVerify"`   // Disables server name verification when connecting over TLS.
}

func New(config ...*Config) (*Redis, error) {
	if len(config) > 0 {
		return &Redis{adapter: NewAdapterRedigo(config[0])}, nil
	}
	configFromGlobal, ok := GetConfig()
	if !ok {
		return nil, gerror.NewCode(
			gcode.CodeMissingConfiguration,
			`configuration not found for creating Redis client`,
		)
	}
	return &Redis{adapter: NewAdapterRedigo(configFromGlobal)}, nil
}

func NewWithAdapter(adapter Adapter) *Redis {
	return &Redis{adapter: adapter}
}

func (r *Redis) SetAdapter(adapter Adapter) {
	r.adapter = adapter
}

// Conn returns a raw underlying connection object,
// which expose more methods to communicate with server.
// **You should call Close function manually if you do not use this connection any further.**
func (r *Redis) Conn(ctx context.Context) (Conn, error) {
	return r.adapter.Conn(ctx)
}

// Stats returns pool's statistics.
func (r *Redis) Stats(ctx context.Context) (Stats, error) {
	return r.adapter.Stats(ctx)
}

// Do sends a command to the server and returns the received reply.
// Do automatically get a connection from pool, and close it when the reply received.
// It does not really "close" the connection, but drops it back to the connection pool.
func (r *Redis) Do(ctx context.Context, command string, args ...interface{}) (*gvar.Var, error) {
	conn, err := r.Conn(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close(ctx)

	if len(args) > 0 {
		lastElement := args[len(args)-1]
		switch option := lastElement.(type) {
		case Option:
			return conn.Do(ctx, command, args, &option)

		case *Option:
			return conn.Do(ctx, command, args, option)

		default:
			return conn.Do(ctx, command, args, nil)
		}
	}
	return conn.Do(ctx, command, args, nil)
}
