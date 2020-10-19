// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
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
	"fmt"
	"time"

	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/container/gvar"
	"github.com/gomodule/redigo/redis"
)

// Redis client.
type Redis struct {
	pool   *redis.Pool // Underlying connection pool.
	group  string      // Configuration group.
	config Config      // Configuration.
}

// Redis connection.
type Conn struct {
	redis.Conn
}

// Redis configuration.
type Config struct {
	Host            string
	Port            int
	Db              int
	Pass            string        // Password for AUTH.
	MaxIdle         int           // Maximum number of connections allowed to be idle (default is 10)
	MaxActive       int           // Maximum number of connections limit (default is 0 means no limit).
	IdleTimeout     time.Duration // Maximum idle time for connection (default is 10 seconds, not allowed to be set to 0)
	MaxConnLifetime time.Duration // Maximum lifetime of the connection (default is 30 seconds, not allowed to be set to 0)
	ConnectTimeout  time.Duration // Dial connection timeout.
	TLS             bool          // Specifies the config to use when a TLS connection is dialed.
	TLSSkipVerify   bool          // Disables server name verification when connecting over TLS
}

// Pool statistics.
type PoolStats struct {
	redis.PoolStats
}

const (
	gDEFAULT_POOL_IDLE_TIMEOUT  = 10 * time.Second
	gDEFAULT_POOL_CONN_TIMEOUT  = 10 * time.Second
	gDEFAULT_POOL_MAX_IDLE      = 10
	gDEFAULT_POOL_MAX_ACTIVE    = 100
	gDEFAULT_POOL_MAX_LIFE_TIME = 30 * time.Second
)

var (
	// Pool map.
	pools = gmap.NewStrAnyMap(true)
)

// New creates a redis client object with given configuration.
// Redis client maintains a connection pool automatically.
func New(config Config) *Redis {
	// The MaxIdle is the most important attribute of the connection pool.
	// Only if this attribute is set, the created connections from client
	// can not exceed the limit of the server.
	if config.MaxIdle == 0 {
		config.MaxIdle = gDEFAULT_POOL_MAX_IDLE
	}
	// This value SHOULD NOT exceed the connection limit of redis server.
	if config.MaxActive == 0 {
		config.MaxActive = gDEFAULT_POOL_MAX_ACTIVE
	}
	if config.IdleTimeout == 0 {
		config.IdleTimeout = gDEFAULT_POOL_IDLE_TIMEOUT
	}
	if config.ConnectTimeout == 0 {
		config.ConnectTimeout = gDEFAULT_POOL_CONN_TIMEOUT
	}
	if config.MaxConnLifetime == 0 {
		config.MaxConnLifetime = gDEFAULT_POOL_MAX_LIFE_TIME
	}
	return &Redis{
		config: config,
		pool: pools.GetOrSetFuncLock(fmt.Sprintf("%v", config), func() interface{} {
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
					// AUTH
					if len(config.Pass) > 0 {
						if _, err := c.Do("AUTH", config.Pass); err != nil {
							return nil, err
						}
					}
					// DB
					if _, err := c.Do("SELECT", config.Db); err != nil {
						return nil, err
					}
					return c, nil
				},
				// After the conn is taken from the connection pool, to test if the connection is available,
				// If error is returned then it closes the connection object and recreate a new connection.
				TestOnBorrow: func(c redis.Conn, t time.Time) error {
					_, err := c.Do("PING")
					return err
				},
			}
		}).(*redis.Pool),
	}
}

// NewFromStr creates a redis client object with given configuration string.
// Redis client maintains a connection pool automatically.
// The parameter <str> like:
// 127.0.0.1:6379,0
// 127.0.0.1:6379,0,password
func NewFromStr(str string) (*Redis, error) {
	config, err := ConfigFromStr(str)
	if err != nil {
		return nil, err
	}
	return New(config), nil
}

// Close closes the redis connection pool,
// it will release all connections reserved by this pool.
// It is not necessary to call Close manually.
func (r *Redis) Close() error {
	if r.group != "" {
		// If it is an instance object,
		// it needs to remove it from the instance Map.
		instances.Remove(r.group)
	}
	pools.Remove(fmt.Sprintf("%v", r.config))
	return r.pool.Close()
}

// Conn returns a raw underlying connection object,
// which expose more methods to communicate with server.
// **You should call Close function manually if you do not use this connection any further.**
func (r *Redis) Conn() *Conn {
	return &Conn{r.pool.Get()}
}

// Alias of Conn, see Conn.
// Deprecated.
func (r *Redis) GetConn() *Conn {
	return r.Conn()
}

// SetMaxIdle sets the maximum number of idle connections in the pool.
func (r *Redis) SetMaxIdle(value int) {
	r.pool.MaxIdle = value
}

// SetMaxActive sets the maximum number of connections allocated by the pool at a given time.
// When zero, there is no limit on the number of connections in the pool.
//
// Note that if the pool is at the MaxActive limit, then all the operations will wait for
// a connection to be returned to the pool before returning.
func (r *Redis) SetMaxActive(value int) {
	r.pool.MaxActive = value
}

// SetIdleTimeout sets the IdleTimeout attribute of the connection pool.
// It closes connections after remaining idle for this duration. If the value
// is zero, then idle connections are not closed. Applications should set
// the timeout to a value less than the server's timeout.
func (r *Redis) SetIdleTimeout(value time.Duration) {
	r.pool.IdleTimeout = value
}

// SetMaxConnLifetime sets the MaxConnLifetime attribute of the connection pool.
// It closes connections older than this duration. If the value is zero, then
// the pool does not close connections based on age.
func (r *Redis) SetMaxConnLifetime(value time.Duration) {
	r.pool.MaxConnLifetime = value
}

// Stats returns pool's statistics.
func (r *Redis) Stats() *PoolStats {
	return &PoolStats{r.pool.Stats()}
}

// Do sends a command to the server and returns the received reply.
// Do automatically get a connection from pool, and close it when the reply received.
// It does not really "close" the connection, but drops it back to the connection pool.
func (r *Redis) Do(commandName string, args ...interface{}) (interface{}, error) {
	conn := &Conn{r.pool.Get()}
	defer conn.Close()
	return conn.Do(commandName, args...)
}

// DoWithTimeout sends a command to the server and returns the received reply.
// The timeout overrides the read timeout set when dialing the connection.
func (r *Redis) DoWithTimeout(timeout time.Duration, commandName string, args ...interface{}) (interface{}, error) {
	conn := &Conn{r.pool.Get()}
	defer conn.Close()
	return conn.DoWithTimeout(timeout, commandName, args...)
}

// DoVar returns value from Do as gvar.Var.
func (r *Redis) DoVar(commandName string, args ...interface{}) (*gvar.Var, error) {
	return resultToVar(r.Do(commandName, args...))
}

// DoVarWithTimeout returns value from Do as gvar.Var.
// The timeout overrides the read timeout set when dialing the connection.
func (r *Redis) DoVarWithTimeout(timeout time.Duration, commandName string, args ...interface{}) (*gvar.Var, error) {
	return resultToVar(r.DoWithTimeout(timeout, commandName, args...))
}
