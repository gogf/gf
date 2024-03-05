// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package redis provides gredis.Adapter implements using go-redis.
package redis

import (
	"crypto/tls"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/text/gstr"
)

// Redis is an implement of Adapter using go-redis.
type Redis struct {
	gredis.AdapterOperation

	client redis.UniversalClient
	config *gredis.Config
}

const (
	defaultPoolMaxIdle     = 10
	defaultPoolMaxActive   = 100
	defaultPoolIdleTimeout = 10 * time.Second
	defaultPoolWaitTimeout = 10 * time.Second
	defaultPoolMaxLifeTime = 30 * time.Second
	defaultMaxRetries      = -1
)

func init() {
	gredis.RegisterAdapterFunc(func(config *gredis.Config) gredis.Adapter {
		return New(config)
	})
}

// New creates and returns a redis adapter using go-redis.
func New(config *gredis.Config) *Redis {
	fillWithDefaultConfiguration(config)
	opts := &redis.UniversalOptions{
		Addrs:            gstr.SplitAndTrim(config.Address, ","),
		Username:         config.User,
		Password:         config.Pass,
		SentinelUsername: config.SentinelUser,
		SentinelPassword: config.SentinelPass,
		DB:               config.Db,
		MaxRetries:       defaultMaxRetries,
		PoolSize:         config.MaxActive,
		MinIdleConns:     config.MinIdle,
		MaxIdleConns:     config.MaxIdle,
		ConnMaxLifetime:  config.MaxConnLifetime,
		ConnMaxIdleTime:  config.IdleTimeout,
		PoolTimeout:      config.WaitTimeout,
		DialTimeout:      config.DialTimeout,
		ReadTimeout:      config.ReadTimeout,
		WriteTimeout:     config.WriteTimeout,
		MasterName:       config.MasterName,
		TLSConfig:        config.TLSConfig,
		Protocol:         config.Protocol,
	}

	var client redis.UniversalClient
	if opts.MasterName != "" {
		redisSentinel := opts.Failover()
		redisSentinel.ReplicaOnly = config.SlaveOnly
		client = redis.NewFailoverClient(redisSentinel)
	} else if len(opts.Addrs) > 1 || config.Cluster {
		client = redis.NewClusterClient(opts.Cluster())
	} else {
		client = redis.NewClient(opts.Simple())
	}

	r := &Redis{
		client: client,
		config: config,
	}
	r.AdapterOperation = r
	return r
}

func fillWithDefaultConfiguration(config *gredis.Config) {
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
	if config.WaitTimeout == 0 {
		config.WaitTimeout = defaultPoolWaitTimeout
	}
	if config.MaxConnLifetime == 0 {
		config.MaxConnLifetime = defaultPoolMaxLifeTime
	}
	if config.WriteTimeout == 0 {
		config.WriteTimeout = -1
	}
	if config.ReadTimeout == 0 {
		config.ReadTimeout = -1
	}
	if config.TLSConfig == nil && config.TLS {
		config.TLSConfig = &tls.Config{
			InsecureSkipVerify: config.TLSSkipVerify,
		}
	}
}
