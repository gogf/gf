// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gredis

import (
	"context"
	"github.com/gogf/gf/errors/gcode"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/internal/intlog"
	"time"

	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/text/gregex"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
)

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

const (
	DefaultGroupName = "default" // Default configuration group name.
	DefaultRedisPort = 6379      // Default redis port configuration if not passed.
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

// SetConfigByStr sets the global configuration for specified group with string.
// If `name` is not passed, it sets configuration for the default group name.
func SetConfigByStr(str string, name ...string) error {
	group := DefaultGroupName
	if len(name) > 0 {
		group = name[0]
	}
	config, err := ConfigFromStr(str)
	if err != nil {
		return err
	}
	localConfigMap.Set(group, config)
	return nil
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

// ConfigFromStr parses and returns config from given str.
// Eg: host:port[,db,pass?maxIdle=x&maxActive=x&idleTimeout=x&maxConnLifetime=x]
func ConfigFromStr(str string) (config *Config, err error) {
	array, _ := gregex.MatchString(`^([^:]+):*(\d*),{0,1}(\d*),{0,1}(.*)\?(.+)$`, str)
	if len(array) == 6 {
		parse, _ := gstr.Parse(array[5])
		config = &Config{
			Host: array[1],
			Port: gconv.Int(array[2]),
			Db:   gconv.Int(array[3]),
			Pass: array[4],
		}
		if config.Port == 0 {
			config.Port = DefaultRedisPort
		}
		if err = gconv.Struct(parse, config); err != nil {
			return nil, err
		}
		return
	}
	array, _ = gregex.MatchString(`([^:]+):*(\d*),{0,1}(\d*),{0,1}(.*)`, str)
	if len(array) == 5 {
		config = &Config{
			Host: array[1],
			Port: gconv.Int(array[2]),
			Db:   gconv.Int(array[3]),
			Pass: array[4],
		}
		if config.Port == 0 {
			config.Port = DefaultRedisPort
		}
	} else {
		err = gerror.NewCodef(gcode.CodeInvalidConfiguration, `invalid redis configuration: "%s"`, str)
	}
	return
}

// ClearConfig removes all configurations of redis.
func ClearConfig() {
	localConfigMap.Clear()
}
