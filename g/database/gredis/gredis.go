// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gredis provides client for redis server.
//
// Redis客户端.
// Redis中文手册文档请参考：http://redisdoc.com/ , Redis官方命令请参考：https://redis.io/commands
package gredis

import (
    "fmt"
    "github.com/gogf/gf/g/container/gmap"
    "github.com/gogf/gf/third/github.com/gomodule/redigo/redis"
    "time"
)

const (
    gDEFAULT_POOL_IDLE_TIMEOUT  = 60 * time.Second
    gDEFAULT_POOL_MAX_LIFE_TIME = 60 * time.Second
)

// Redis客户端(管理连接池)
type Redis struct {
    pool   *redis.Pool
    config Config
}

// Redis连接对象(连接池中的单个连接)
type Conn redis.Conn

// Redis服务端但节点连接配置信息
type Config struct {
    Host            string        // 地址
    Port            int           // 端口
    Db              int           // 数据库
    Pass            string        // 授权密码
    MaxIdle         int           // 最大允许空闲存在的连接数(默认为0表示不存在闲置连接)
    MaxActive       int           // 最大连接数量限制(默认为0表示不限制)
    IdleTimeout     time.Duration // 连接最大空闲时间(默认为60秒,不允许设置为0)
    MaxConnLifetime time.Duration // 连接最长存活时间(默认为60秒,不允许设置为0)
}

// Redis链接池统计信息
type PoolStats struct {
    redis.PoolStats
}

// 连接池map
var pools = gmap.NewStringInterfaceMap()

// New creates a redis client object with given configuration.
// Redis client maintains a connection pool automatically.
//
// 创建redis操作对象.
func New(config Config) *Redis {
    if config.IdleTimeout == 0 {
        config.IdleTimeout = gDEFAULT_POOL_IDLE_TIMEOUT
    }
    if config.MaxConnLifetime == 0 {
        config.MaxConnLifetime = gDEFAULT_POOL_MAX_LIFE_TIME
    }
    return &Redis{
        config : config,
        pool   : pools.GetOrSetFuncLock(fmt.Sprintf("%v", config), func() interface{} {
            return &redis.Pool {
                IdleTimeout     : config.IdleTimeout,
                MaxConnLifetime : config.MaxConnLifetime,
                Dial            : func() (redis.Conn, error) {
                    c, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", config.Host, config.Port))
                    if err != nil {
                        return nil, err
                    }
                    // 密码设置
                    if len(config.Pass) > 0 {
                        if _, err := c.Do("AUTH", config.Pass); err != nil {
                            return nil, err
                        }
                    }
                    // 数据库设置
                    if _, err := c.Do("SELECT", config.Db); err != nil {
                        return nil, err
                    }
                    return c, nil
                },
                // 在被应用从连接池中获取出来之后，用以测试连接是否可用，
                // 如果返回error那么关闭该连接对象重新创建新的连接。
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
//
// 关闭redis管理对象，将会关闭底层的连接池。
func (r *Redis) Close() error {
    pools.Remove(fmt.Sprintf("%v", r.config))
    return r.pool.Close()
}

// See GetConn.
func (r *Redis) Conn() Conn {
    return r.GetConn()
}

// GetConn returns a raw connection object,
// which expose more methods communication with server.
// **You should call Close function manually if you do not use this connection any further.**
//
// 获得一个原生的redis连接对象，用于自定义连接操作，
// 但是需要注意的是如果不再使用该连接对象时，需要手动Close连接，否则会造成连接数超限。
func (r *Redis) GetConn() Conn {
    return r.pool.Get().(Conn)
}

// SetMaxIdle sets the MaxIdle attribute of the connection pool.
//
// 设置属性 - MaxIdle
func (r *Redis) SetMaxIdle(value int) {
    r.pool.MaxIdle = value
}

// SetMaxIdle sets the MaxActive attribute of the connection pool.
//
// 设置属性 - MaxActive
func (r *Redis) SetMaxActive(value int) {
    r.pool.MaxActive = value
}

// SetMaxIdle sets the IdleTimeout attribute of the connection pool.
//
// 设置属性 - IdleTimeout
func (r *Redis) SetIdleTimeout(value time.Duration) {
    r.pool.IdleTimeout = value
}

// SetMaxIdle sets the MaxConnLifetime attribute of the connection pool.
//
// 设置属性 - MaxConnLifetime
func (r *Redis) SetMaxConnLifetime(value time.Duration) {
    r.pool.MaxConnLifetime = value
}

// Stats returns pool's statistics.
//
// 获取当前连接池统计信息。
func (r *Redis) Stats() *PoolStats {
    return &PoolStats{r.pool.Stats()}
}

// Do sends a command to the server and returns the received reply.
// Do automatically get a connection from pool, and close it when reply received.
//
// 执行同步命令，自动从连接池中获取连接，使用完毕后关闭连接（丢回连接池），开发者不用自行Close.
func (r *Redis) Do(command string, args ...interface{}) (interface{}, error) {
    conn := r.pool.Get()
    defer conn.Close()
    return conn.Do(command, args...)
}

// Deprecated.
// Send writes the command to the client's output buffer.
//
// 执行异步命令 - Send
func (r *Redis) Send(command string, args ...interface{}) error {
    conn := r.pool.Get()
    defer conn.Close()
    return conn.Send(command, args...)
}

