// Copyright GoFrame Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"fmt"
	"github.com/gogf/gf/os/gcache"
	"sync"
	"time"

	"github.com/gogf/gf/os/glog"
)

const (
	DEFAULT_GROUP_NAME = "default" // Deprecated, use DefaultGroupName instead.
	DefaultGroupName   = "default" // Default group name.
)

// Config is the configuration management object.
type Config map[string]ConfigGroup

// ConfigGroup is a slice of configuration node for specified named group.
type ConfigGroup []ConfigNode

// ConfigNode is configuration for one node.
type ConfigNode struct {
	Host                 string        // Host of server, ip or domain like: 127.0.0.1, localhost
	Port                 string        // Port, it's commonly 3306.
	User                 string        // Authentication username.
	Pass                 string        // Authentication password.
	Name                 string        // Default used database name.
	Type                 string        // Database type: mysql, sqlite, mssql, pgsql, oracle.
	Role                 string        // (Optional, "master" in default) Node role, used for master-slave mode: master, slave.
	Debug                bool          // (Optional) Debug mode enables debug information logging and output.
	Prefix               string        // (Optional) Table prefix.
	DryRun               bool          // (Optional) Dry run, which does SELECT but no INSERT/UPDATE/DELETE statements.
	Weight               int           // (Optional) Weight for load balance calculating, it's useless if there's just one node.
	Charset              string        // (Optional, "utf8mb4" in default) Custom charset when operating on database.
	LinkInfo             string        `json:"link"`        // (Optional) Custom link information, when it is used, configuration Host/Port/User/Pass/Name are ignored.
	MaxIdleConnCount     int           `json:"maxidle"`     // (Optional) Max idle connection configuration for underlying connection pool.
	MaxOpenConnCount     int           `json:"maxopen"`     // (Optional) Max open connection configuration for underlying connection pool.
	MaxConnLifetime      time.Duration `json:"maxlifetime"` // (Optional) Max connection TTL configuration for underlying connection pool.
	CreatedAt            string        // (Optional) The filed name of table for automatic-filled created datetime.
	UpdatedAt            string        // (Optional) The filed name of table for automatic-filled updated datetime.
	DeletedAt            string        // (Optional) The filed name of table for automatic-filled updated datetime.
	TimeMaintainDisabled bool          // (Optional) Disable the automatic time maintaining feature.
}

// configs is internal used configuration object.
var configs struct {
	sync.RWMutex
	config Config // All configurations.
	group  string // Default configuration group.
}

func init() {
	configs.config = make(Config)
	configs.group = DefaultGroupName
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
	AddConfigNode(DefaultGroupName, node)
}

// AddDefaultConfigGroup adds multiple node configurations to configuration of default group.
func AddDefaultConfigGroup(nodes ConfigGroup) {
	SetConfigGroup(DefaultGroupName, nodes)
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

// IsConfigured checks and returns whether the database configured.
// It returns true if any configuration exists.
func IsConfigured() bool {
	configs.RLock()
	defer configs.RUnlock()
	return len(configs.config) > 0
}

// SetLogger sets the logger for orm.
func (c *Core) SetLogger(logger *glog.Logger) {
	c.logger = logger
}

// GetLogger returns the logger of the orm.
func (c *Core) GetLogger() *glog.Logger {
	return c.logger
}

// SetMaxIdleConnCount sets the max idle connection count for underlying connection pool.
func (c *Core) SetMaxIdleConnCount(n int) {
	c.config.MaxIdleConnCount = n
}

// SetMaxOpenConnCount sets the max open connection count for underlying connection pool.
func (c *Core) SetMaxOpenConnCount(n int) {
	c.config.MaxOpenConnCount = n
}

// SetMaxConnLifetime sets the connection TTL for underlying connection pool.
// If parameter <d> <= 0, it means the connection never expires.
func (c *Core) SetMaxConnLifetime(d time.Duration) {
	c.config.MaxConnLifetime = d
}

// String returns the node as string.
func (node *ConfigNode) String() string {
	return fmt.Sprintf(
		`%s@%s:%s,%s,%s,%s,%s,%v,%d-%d-%d#%s`,
		node.User, node.Host, node.Port,
		node.Name, node.Type, node.Role, node.Charset, node.Debug,
		node.MaxIdleConnCount,
		node.MaxOpenConnCount,
		node.MaxConnLifetime,
		node.LinkInfo,
	)
}

// GetConfig returns the current used node configuration.
func (c *Core) GetConfig() *ConfigNode {
	return c.config
}

// SetDebug enables/disables the debug mode.
func (c *Core) SetDebug(debug bool) {
	c.debug.Set(debug)
}

// GetDebug returns the debug value.
func (c *Core) GetDebug() bool {
	return c.debug.Val()
}

// GetCache returns the internal cache object.
func (c *Core) GetCache() *gcache.Cache {
	return c.cache
}

// GetGroup returns the group string configured.
func (c *Core) GetGroup() string {
	return c.group
}

// SetDryRun enables/disables the DryRun feature.
// Deprecated, use GetConfig instead.
func (c *Core) SetDryRun(enabled bool) {
	c.config.DryRun = enabled
}

// GetDryRun returns the DryRun value.
// Deprecated, use GetConfig instead.
func (c *Core) GetDryRun() bool {
	return c.config.DryRun
}

// GetPrefix returns the table prefix string configured.
// Deprecated, use GetConfig instead.
func (c *Core) GetPrefix() string {
	return c.config.Prefix
}

// SetSchema changes the schema for this database connection object.
// Importantly note that when schema configuration changed for the database,
// it affects all operations on the database object in the future.
func (c *Core) SetSchema(schema string) {
	c.schema.Set(schema)
}

// GetSchema returns the schema configured.
func (c *Core) GetSchema() string {
	return c.schema.Val()
}
