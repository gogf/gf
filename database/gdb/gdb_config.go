// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"fmt"
	"sync"
	"time"

	"github.com/gogf/gf/os/glog"
)

const (
	DEFAULT_GROUP_NAME = "default" // Default group name.
)

// Config is the configuration management object.
type Config map[string]ConfigGroup

// ConfigGroup is a slice of configuration node for specified named group.
type ConfigGroup []ConfigNode

// ConfigNode is configuration for one node.
type ConfigNode struct {
	Host             string        // Host of server, ip or domain like: 127.0.0.1, localhost
	Port             string        // Port, it's commonly 3306.
	User             string        // Authentication username.
	Pass             string        // Authentication password.
	Name             string        // Default used database name.
	Type             string        // Database type: mysql, sqlite, mssql, pgsql, oracle.
	Role             string        // (Optional, "master" in default) Node role, used for master-slave mode: master, slave.
	Debug            bool          // (Optional) Debug mode enables debug information logging and output.
	Prefix           string        // (Optional) Table prefix.
	Weight           int           // (Optional) Weight for load balance calculating, it's useless if there's just one node.
	Charset          string        // (Optional, "utf8mb4" in default) Custom charset when operating on database.
	LinkInfo         string        // (Optional) Custom link information, when it is used, configuration Host/Port/User/Pass/Name are ignored.
	MaxIdleConnCount int           // (Optional) Max idle connection configuration for underlying connection pool.
	MaxOpenConnCount int           // (Optional) Max open connection configuration for underlying connection pool.
	MaxConnLifetime  time.Duration // (Optional) Max connection TTL configuration for underlying connection pool.
}

// configs is internal used configuration object.
var configs struct {
	sync.RWMutex
	config Config // All configurations.
	group  string // Default configuration group.
}

func init() {
	configs.config = make(Config)
	configs.group = DEFAULT_GROUP_NAME
}

// SetConfig sets the global configuration for package.
// It will overwrite the old configuration of package.
func SetConfig(config Config) {
	defer instances.Clear()
	configs.Lock()
	defer configs.Unlock()
	configs.config = config
}

// SetConfigGroup sets the configuration for given group.
func SetConfigGroup(group string, nodes ConfigGroup) {
	defer instances.Clear()
	configs.Lock()
	defer configs.Unlock()
	configs.config[group] = nodes
}

// AddConfigNode adds one node configuration to configuration of given group.
func AddConfigNode(group string, node ConfigNode) {
	defer instances.Clear()
	configs.Lock()
	defer configs.Unlock()
	configs.config[group] = append(configs.config[group], node)
}

// AddDefaultConfigNode adds one node configuration to configuration of default group.
func AddDefaultConfigNode(node ConfigNode) {
	AddConfigNode(DEFAULT_GROUP_NAME, node)
}

// AddDefaultConfigGroup adds multiple node configurations to configuration of default group.
func AddDefaultConfigGroup(nodes ConfigGroup) {
	SetConfigGroup(DEFAULT_GROUP_NAME, nodes)
}

// GetConfig retrieves and returns the configuration of given group.
func GetConfig(group string) ConfigGroup {
	configs.RLock()
	defer configs.RUnlock()
	return configs.config[group]
}

// SetDefaultGroup sets the group name for default configuration.
func SetDefaultGroup(name string) {
	defer instances.Clear()
	configs.Lock()
	defer configs.Unlock()
	configs.group = name
}

// GetDefaultGroup returns the { name of default configuration.
func GetDefaultGroup() string {
	defer instances.Clear()
	configs.RLock()
	defer configs.RUnlock()
	return configs.group
}

// SetLogger sets the logger for orm.
func (bs *dbBase) SetLogger(logger *glog.Logger) {
	bs.logger = logger
}

// GetLogger returns the logger of the orm.
func (bs *dbBase) GetLogger() *glog.Logger {
	return bs.logger
}

// SetMaxIdleConnCount sets the max idle connection count for underlying connection pool.
func (bs *dbBase) SetMaxIdleConnCount(n int) {
	bs.maxIdleConnCount = n
}

// SetMaxOpenConnCount sets the max open connection count for underlying connection pool.
func (bs *dbBase) SetMaxOpenConnCount(n int) {
	bs.maxOpenConnCount = n
}

// SetMaxConnLifetime sets the connection TTL for underlying connection pool.
// If parameter <d> <= 0, it means the connection never expires.
func (bs *dbBase) SetMaxConnLifetime(d time.Duration) {
	bs.maxConnLifetime = d
}

// String returns the node as string.
func (node *ConfigNode) String() string {
	if node.LinkInfo != "" {
		return node.LinkInfo
	}
	return fmt.Sprintf(
		`%s@%s:%s,%s,%s,%s,%s,%v,%d-%d-%d`,
		node.User, node.Host, node.Port,
		node.Name, node.Type, node.Role, node.Charset, node.Debug,
		node.MaxIdleConnCount,
		node.MaxOpenConnCount,
		node.MaxConnLifetime,
	)
}

// SetDebug enables/disables the debug mode.
func (bs *dbBase) SetDebug(debug bool) {
	if bs.debug.Val() == debug {
		return
	}
	bs.debug.Set(debug)
}

// getDebug returns the debug value.
func (bs *dbBase) getDebug() bool {
	return bs.debug.Val()
}
