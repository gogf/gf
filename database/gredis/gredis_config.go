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
	Address         string        `json:"address"`         // It supports single and cluster redis server. Multiple addresses joined with char ','. Eg: 192.168.1.1:6379, 192.168.1.2:6379.
	Db              int           `json:"db"`              // Redis db.
	Pass            string        `json:"pass"`            // Password for AUTH.
	MinIdle         int           `json:"minIdle"`         // Minimum number of connections allowed to be idle (default is 0)
	MaxIdle         int           `json:"maxIdle"`         // Maximum number of connections allowed to be idle (default is 10)
	MaxActive       int           `json:"maxActive"`       // Maximum number of connections limit (default is 0 means no limit).
	MaxConnLifetime time.Duration `json:"maxConnLifetime"` // Maximum lifetime of the connection (default is 30 seconds, not allowed to be set to 0)
	IdleTimeout     time.Duration `json:"idleTimeout"`     // Maximum idle time for connection (default is 10 seconds, not allowed to be set to 0)
	WaitTimeout     time.Duration `json:"waitTimeout"`     // Timed out duration waiting to get a connection from the connection pool.
	DialTimeout     time.Duration `json:"dialTimeout"`     // Dial connection timeout for TCP.
	ReadTimeout     time.Duration `json:"readTimeout"`     // Read timeout for TCP. DO NOT set it if not necessary.
	WriteTimeout    time.Duration `json:"writeTimeout"`    // Write timeout for TCP.
	MasterName      string        `json:"masterName"`      // Used in Redis Sentinel mode.
	TLS             bool          `json:"tls"`             // Specifies whether TLS should be used when connecting to the server.
	TLSSkipVerify   bool          `json:"tlsSkipVerify"`   // Disables server name verification when connecting over TLS.
	TLSConfig       *tls.Config   `json:"-"`               // TLS Config to use. When set TLS will be negotiated.
	SlaveOnly       bool          `json:"slaveOnly"`       // Route all commands to slave read-only nodes.
}

const (
	DefaultGroupName = "default" // Default configuration group name.
)

var (
	// Configuration groups.
	localConfigMap = gmap.NewStrAnyMap(true)
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
func SetConfigByMap(m map[string]interface{}, name ...string) error {
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
func ConfigFromMap(m map[string]interface{}) (config *Config, err error) {
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
		return v.(*Config), true
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
