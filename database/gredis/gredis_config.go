// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gredis

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/util/gconv"
)

// Config is redis configuration.
type Config struct {
	// Address It supports single and cluster redis server. Multiple addresses joined with char ','. Eg: 192.168.1.1:6379, 192.168.1.2:6379.
	Address         string        `json:"address"         v:"required" dc:"Redis server address|i18n:config.redis.address"`
	Db              int           `json:"db"              d:"0" dc:"Redis database index|i18n:config.redis.db"`              // Redis db.
	User            string        `json:"user"            dc:"Username for AUTH|i18n:config.redis.user"`            // Username for AUTH.
	Pass            string        `json:"pass"            dc:"Password for AUTH|i18n:config.redis.pass"`            // Password for AUTH.
	SentinelUser    string        `json:"sentinel_user"   dc:"Username for sentinel AUTH|i18n:config.redis.sentinelUser"`   // Username for sentinel AUTH.
	SentinelPass    string        `json:"sentinel_pass"   dc:"Password for sentinel AUTH|i18n:config.redis.sentinelPass"`   // Password for sentinel AUTH.
	MinIdle         int           `json:"minIdle"         d:"0" dc:"Min idle connections|i18n:config.redis.minIdle"`         // Minimum number of connections allowed to be idle (default is 0)
	MaxIdle         int           `json:"maxIdle"         d:"10" dc:"Max idle connections|i18n:config.redis.maxIdle"`         // Maximum number of connections allowed to be idle (default is 10)
	MaxActive       int           `json:"maxActive"       d:"0" dc:"Max active connections (0=unlimited)|i18n:config.redis.maxActive"`       // Maximum number of connections limit (default is 0 means no limit).
	MaxConnLifetime time.Duration `json:"maxConnLifetime" d:"30s" dc:"Max connection lifetime|i18n:config.redis.maxConnLifetime"` // Maximum lifetime of the connection (default is 30 seconds, not allowed to be set to 0)
	IdleTimeout     time.Duration `json:"idleTimeout"     d:"10s" dc:"Idle connection timeout|i18n:config.redis.idleTimeout"`     // Maximum idle time for connection (default is 10 seconds, not allowed to be set to 0)
	WaitTimeout     time.Duration `json:"waitTimeout"     dc:"Wait timeout for connection pool|i18n:config.redis.waitTimeout"`     // Timed out duration waiting to get a connection from the connection pool.
	DialTimeout     time.Duration `json:"dialTimeout"     dc:"Dial connection timeout|i18n:config.redis.dialTimeout"`     // Dial connection timeout for TCP.
	ReadTimeout     time.Duration `json:"readTimeout"     dc:"Read timeout|i18n:config.redis.readTimeout"`     // Read timeout for TCP. DO NOT set it if not necessary.
	WriteTimeout    time.Duration `json:"writeTimeout"    dc:"Write timeout|i18n:config.redis.writeTimeout"`    // Write timeout for TCP.
	MasterName      string        `json:"masterName"      dc:"Master name for Sentinel mode|i18n:config.redis.masterName"`      // Used in Redis Sentinel mode.
	TLS             bool          `json:"tls"             d:"false" dc:"Enable TLS connection|i18n:config.redis.tls"`             // Specifies whether TLS should be used when connecting to the server.
	TLSSkipVerify   bool          `json:"tlsSkipVerify"   d:"false" dc:"Skip TLS server name verification|i18n:config.redis.tlsSkipVerify"`   // Disables server name verification when connecting over TLS.
	TLSConfig       *tls.Config   `json:"-"`               // TLS Config to use. When set TLS will be negotiated.
	SlaveOnly       bool          `json:"slaveOnly"       d:"false" dc:"Route commands to slave nodes only|i18n:config.redis.slaveOnly"`       // Route all commands to slave read-only nodes.
	Cluster         bool          `json:"cluster"         d:"false" dc:"Enable cluster mode|i18n:config.redis.cluster"`         // Specifies whether cluster mode be used.
	Protocol        int           `json:"protocol"        d:"3" dc:"RESP protocol version (2 or 3)|i18n:config.redis.protocol"`        // Specifies the RESP version (Protocol 2 or 3.)
}

const (
	DefaultGroupName = "default" // Default configuration group name.
)

var (
	// configChecker checks whether the *Config is nil.
	configChecker = func(v *Config) bool { return v == nil }
	// Configuration groups.
	localConfigMap = gmap.NewKVMapWithChecker[string, *Config](configChecker, true)
)

// SetConfig sets the global configuration for specified group.
// If `name` is not passed, it sets configuration for the default group name.
func SetConfig(config *Config, name ...string) {
	group := DefaultGroupName
	if len(name) > 0 {
		group = name[0]
	}
	localConfigMap.Set(group, config)

	intlog.Printf(context.TODO(), `SetConfig for group "%s": %+v`, group, config)
}

// SetConfigByMap sets the global configuration for specified group with map.
// If `name` is not passed, it sets configuration for the default group name.
func SetConfigByMap(m map[string]any, name ...string) error {
	group := DefaultGroupName
	if len(name) > 0 {
		group = name[0]
	}
	config, err := ConfigFromMap(m)
	if err != nil {
		return err
	}
	localConfigMap.Set(group, config)
	return nil
}

// ConfigFromMap parses and returns config from given map.
func ConfigFromMap(m map[string]any) (config *Config, err error) {
	config = &Config{}
	if err = gconv.Scan(m, config); err != nil {
		err = gerror.NewCodef(gcode.CodeInvalidConfiguration, `invalid redis configuration: %#v`, m)
	}
	if config.DialTimeout < time.Second {
		config.DialTimeout = config.DialTimeout * time.Second
	}
	if config.WaitTimeout < time.Second {
		config.WaitTimeout = config.WaitTimeout * time.Second
	}
	if config.WriteTimeout < time.Second {
		config.WriteTimeout = config.WriteTimeout * time.Second
	}
	if config.ReadTimeout < time.Second {
		config.ReadTimeout = config.ReadTimeout * time.Second
	}
	if config.IdleTimeout < time.Second {
		config.IdleTimeout = config.IdleTimeout * time.Second
	}
	if config.MaxConnLifetime < time.Second {
		config.MaxConnLifetime = config.MaxConnLifetime * time.Second
	}
	if config.Protocol != 2 && config.Protocol != 3 {
		config.Protocol = 3
	}
	return
}

// GetConfig returns the global configuration with specified group name.
// If `name` is not passed, it returns configuration of the default group name.
func GetConfig(name ...string) (config *Config, ok bool) {
	group := DefaultGroupName
	if len(name) > 0 {
		group = name[0]
	}
	if v := localConfigMap.Get(group); v != nil {
		return v, true
	}
	return &Config{}, false
}

// RemoveConfig removes the global configuration with specified group.
// If `name` is not passed, it removes configuration of the default group name.
func RemoveConfig(name ...string) {
	group := DefaultGroupName
	if len(name) > 0 {
		group = name[0]
	}
	localConfigMap.Remove(group)

	intlog.Printf(context.TODO(), `RemoveConfig: %s`, group)
}

// ClearConfig removes all configurations of redis.
func ClearConfig() {
	localConfigMap.Clear()
}
