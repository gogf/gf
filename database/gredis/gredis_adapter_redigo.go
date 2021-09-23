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
	"fmt"
	"github.com/gogf/gf/errors/gcode"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/internal/intlog"
	"time"

	"github.com/gogf/gf/container/gmap"
	"github.com/gomodule/redigo/redis"
)

type AdapterRedigo struct {
	pool   *redis.Pool // Underlying connection pool.
	config *Config     // Configuration.
}

const (
	defaultPoolIdleTimeout = 10 * time.Second
	defaultPoolConnTimeout = 10 * time.Second
	defaultPoolMaxIdle     = 10
	defaultPoolMaxActive   = 100
	defaultPoolMaxLifeTime = 30 * time.Second
)

var (
	localAdapterRedigoPools = gmap.NewStrAnyMap(true)
)

func NewAdapterRedigo(config *Config) *AdapterRedigo {
	// The MaxIdle is the most important attribute of the connection pool.
	// Only if this attribute is set, the created connections from client
	// can not exceed the limit of the server.
	if config.MaxIdle == 0 {
		config.MaxIdle = defaultPoolMaxIdle
	}
	// This value SHOULD NOT exceed the connection limit of redis server.
	if config.MaxActive == 0 {
		config.MaxActive = defaultPoolMaxActive
	}
	if config.IdleTimeout == 0 {
		config.IdleTimeout = defaultPoolIdleTimeout
	}
	if config.ConnectTimeout == 0 {
		config.ConnectTimeout = defaultPoolConnTimeout
	}
	if config.MaxConnLifetime == 0 {
		config.MaxConnLifetime = defaultPoolMaxLifeTime
	}
	return &AdapterRedigo{
		config: config,
		pool: localAdapterRedigoPools.GetOrSetFuncLock(fmt.Sprintf("%v", config), func() interface{} {
			return &redis.Pool{
				Wait:            true,
				IdleTimeout:     config.IdleTimeout,
				MaxActive:       config.MaxActive,
				MaxIdle:         config.MaxIdle,
				MaxConnLifetime: config.MaxConnLifetime,
				Dial: func() (redis.Conn, error) {
					c, err := redis.Dial(
						"tcp",
						fmt.Sprintf("%s:%d", config.Host, config.Port),
						redis.DialConnectTimeout(config.ConnectTimeout),
						redis.DialUseTLS(config.TLS),
						redis.DialTLSSkipVerify(config.TLSSkipVerify),
					)
					if err != nil {
						return nil, err
					}
					intlog.Printf(context.TODO(), `open new connection, config:%+v`, config)
					// AUTH
					if len(config.Pass) > 0 {
						if _, err = c.Do("AUTH", config.Pass); err != nil {
							return nil, err
						}
					}
					// DB
					if _, err = c.Do("SELECT", config.Db); err != nil {
						return nil, err
					}
					return c, nil
				},
				// After conn is taken from the connection pool, to test if the connection is available,
				// If error is returned then it closes the connection object and recreate a new connection.
				TestOnBorrow: func(c redis.Conn, t time.Time) error {
					_, err := c.Do("PING")
					return err
				},
			}
		}).(*redis.Pool),
	}
}

// Close closes the redis connection pool,
// it will release all connections reserved by this pool.
// It is not necessary to call Close manually.
func (r *AdapterRedigo) Close(ctx context.Context) error {
	localAdapterRedigoPools.Remove(fmt.Sprintf("%v", r.config))
	return r.pool.Close()
}

// Conn returns a raw underlying connection object,
// which expose more methods to communicate with server.
// **You should call Close function manually if you do not use this connection any further.**
func (r *AdapterRedigo) Conn(ctx context.Context) (Conn, error) {
	conn := r.pool.Get()
	if conn == nil {
		return nil, gerror.NewCode(gcode.CodeInternalError, `retrieving connection from pool failed`)
	}
	connTimeout, ok := conn.(redis.ConnWithTimeout)
	if !ok {
		return nil, gerror.NewCode(gcode.CodeNotSupported, `current connection does not support "ConnWithTimeout"`)
	}
	return &localAdapterRedisConn{ConnWithTimeout: connTimeout}, nil
}

// Stats returns pool's statistics.
func (r *AdapterRedigo) Stats(ctx context.Context) (Stats, error) {
	return &localAdapterRedigoPoolStats{r.pool.Stats()}, nil
}
