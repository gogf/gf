// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"fmt"
	"sync"
	"time"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

// Config is the configuration management object.
type Config map[string]ConfigGroup

// ConfigGroup is a slice of configuration node for specified named group.
type ConfigGroup []ConfigNode

// ConfigNode is configuration for one node.
type ConfigNode struct {
	// Host specifies the server address, can be either IP address or domain name
	// Example: "127.0.0.1", "localhost"
	Host string `json:"host"`

	// Port specifies the server port number
	// Default is typically "3306" for MySQL
	Port string `json:"port"`

	// User specifies the authentication username for database connection
	User string `json:"user"`

	// Pass specifies the authentication password for database connection
	Pass string `json:"pass"`

	// Name specifies the default database name to be used
	Name string `json:"name"`

	// Type specifies the database type
	// Example: mysql, mariadb, sqlite, mssql, pgsql, oracle, clickhouse, dm.
	Type string `json:"type"`

	// Link provides custom connection string that combines all configuration in one string
	// Optional field
	Link string `json:"link"`

	// Extra provides additional configuration options for third-party database drivers
	// Optional field
	Extra string `json:"extra"`

	// Role specifies the node role in master-slave setup
	// Optional field, defaults to "master"
	// Available values: "master", "slave"
	Role Role `json:"role"`

	// Debug enables debug mode for logging and output
	// Optional field
	Debug bool `json:"debug"`

	// Prefix specifies the table name prefix
	// Optional field
	Prefix string `json:"prefix"`

	// DryRun enables simulation mode where SELECT statements are executed
	// but INSERT/UPDATE/DELETE statements are not
	// Optional field
	DryRun bool `json:"dryRun"`

	// Weight specifies the node weight for load balancing calculations
	// Optional field, only effective in multi-node setups
	Weight int `json:"weight"`

	// Charset specifies the character set for database operations
	// Optional field, defaults to "utf8"
	Charset string `json:"charset"`

	// Protocol specifies the network protocol for database connection
	// Optional field, defaults to "tcp"
	// See net.Dial for available network protocols
	Protocol string `json:"protocol"`

	// Timezone sets the time zone for timestamp interpretation and display
	// Optional field
	Timezone string `json:"timezone"`

	// Namespace specifies the schema namespace for certain databases
	// Optional field, e.g., in PostgreSQL, Name is the catalog and Namespace is the schema
	Namespace string `json:"namespace"`

	// MaxIdleConnCount specifies the maximum number of idle connections in the pool
	// Optional field
	MaxIdleConnCount int `json:"maxIdle"`

	// MaxOpenConnCount specifies the maximum number of open connections in the pool
	// Optional field
	MaxOpenConnCount int `json:"maxOpen"`

	// MaxConnLifeTime specifies the maximum lifetime of a connection
	// Optional field
	MaxConnLifeTime time.Duration `json:"maxLifeTime"`

	// QueryTimeout specifies the maximum execution time for DQL operations
	// Optional field
	QueryTimeout time.Duration `json:"queryTimeout"`

	// ExecTimeout specifies the maximum execution time for DML operations
	// Optional field
	ExecTimeout time.Duration `json:"execTimeout"`

	// TranTimeout specifies the maximum execution time for a transaction block
	// Optional field
	TranTimeout time.Duration `json:"tranTimeout"`

	// PrepareTimeout specifies the maximum execution time for prepare operations
	// Optional field
	PrepareTimeout time.Duration `json:"prepareTimeout"`

	// CreatedAt specifies the field name for automatic timestamp on record creation
	// Optional field
	CreatedAt string `json:"createdAt"`

	// UpdatedAt specifies the field name for automatic timestamp on record updates
	// Optional field
	UpdatedAt string `json:"updatedAt"`

	// DeletedAt specifies the field name for automatic timestamp on record deletion
	// Optional field
	DeletedAt string `json:"deletedAt"`

	// TimeMaintainDisabled controls whether automatic time maintenance is disabled
	// Optional field
	TimeMaintainDisabled bool `json:"timeMaintainDisabled"`
}

type Role string

const (
	RoleMaster Role = "master"
	RoleSlave  Role = "slave"
)

const (
	DefaultGroupName = "default" // Default group name.
)

// configs specifies internal used configuration object.
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
func SetConfig(config Config) error {
	defer instances.Clear()
	configs.Lock()
	defer configs.Unlock()

	for k, nodes := range config {
		for i, node := range nodes {
			parsedNode, err := parseConfigNode(node)
			if err != nil {
				return err
			}
			nodes[i] = parsedNode
		}
		config[k] = nodes
	}
	configs.config = config
	return nil
}

// SetConfigGroup sets the configuration for given group.
func SetConfigGroup(group string, nodes ConfigGroup) error {
	defer instances.Clear()
	configs.Lock()
	defer configs.Unlock()

	for i, node := range nodes {
		parsedNode, err := parseConfigNode(node)
		if err != nil {
			return err
		}
		nodes[i] = parsedNode
	}
	configs.config[group] = nodes
	return nil
}

// AddConfigNode adds one node configuration to configuration of given group.
func AddConfigNode(group string, node ConfigNode) error {
	defer instances.Clear()
	configs.Lock()
	defer configs.Unlock()

	parsedNode, err := parseConfigNode(node)
	if err != nil {
		return err
	}
	configs.config[group] = append(configs.config[group], parsedNode)
	return nil
}

// parseConfigNode parses `Link` configuration syntax.
func parseConfigNode(node ConfigNode) (ConfigNode, error) {
	if node.Link != "" {
		parsedLinkNode, err := parseConfigNodeLink(&node)
		if err != nil {
			return node, err
		}
		node = *parsedLinkNode
	}
	if node.Link != "" && node.Type == "" {
		match, _ := gregex.MatchString(`([a-z]+):(.+)`, node.Link)
		if len(match) == 3 {
			node.Type = gstr.Trim(match[1])
			node.Link = gstr.Trim(match[2])
		}
	}
	return node, nil
}

// AddDefaultConfigNode adds one node configuration to configuration of default group.
func AddDefaultConfigNode(node ConfigNode) error {
	return AddConfigNode(DefaultGroupName, node)
}

// AddDefaultConfigGroup adds multiple node configurations to configuration of default group.
func AddDefaultConfigGroup(nodes ConfigGroup) error {
	return SetConfigGroup(DefaultGroupName, nodes)
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
func (c *Core) SetLogger(logger glog.ILogger) {
	c.logger = logger
}

// GetLogger returns the (logger) of the orm.
func (c *Core) GetLogger() glog.ILogger {
	return c.logger
}

// SetMaxIdleConnCount sets the maximum number of connections in the idle
// connection pool.
//
// If MaxOpenConns is greater than 0 but less than the new MaxIdleConns,
// then the new MaxIdleConns will be reduced to match the MaxOpenConns limit.
//
// If n <= 0, no idle connections are retained.
//
// The default max idle connections is currently 2. This may change in
// a future release.
func (c *Core) SetMaxIdleConnCount(n int) {
	c.dynamicConfig.MaxIdleConnCount = n
}

// SetMaxOpenConnCount sets the maximum number of open connections to the database.
//
// If MaxIdleConns is greater than 0 and the new MaxOpenConns is less than
// MaxIdleConns, then MaxIdleConns will be reduced to match the new
// MaxOpenConns limit.
//
// If n <= 0, then there is no limit on the number of open connections.
// The default is 0 (unlimited).
func (c *Core) SetMaxOpenConnCount(n int) {
	c.dynamicConfig.MaxOpenConnCount = n
}

// SetMaxConnLifeTime sets the maximum amount of time a connection may be reused.
//
// Expired connections may be closed lazily before reuse.
//
// If d <= 0, connections are not closed due to a connection's age.
func (c *Core) SetMaxConnLifeTime(d time.Duration) {
	c.dynamicConfig.MaxConnLifeTime = d
}

// GetConfig returns the current used node configuration.
func (c *Core) GetConfig() *ConfigNode {
	var configNode = c.getConfigNodeFromCtx(c.db.GetCtx())
	if configNode != nil {
		// Note:
		// It so here checks and returns the config from current DB,
		// if different schemas between current DB and config.Name from context,
		// for example, in nested transaction scenario, the context is passed all through the logic procedure,
		// but the config.Name from context may be still the original one from the first transaction object.
		if c.config.Name == configNode.Name {
			return configNode
		}
	}
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
func (c *Core) SetDryRun(enabled bool) {
	c.config.DryRun = enabled
}

// GetDryRun returns the DryRun value.
func (c *Core) GetDryRun() bool {
	return c.config.DryRun || allDryRun
}

// GetPrefix returns the table prefix string configured.
func (c *Core) GetPrefix() string {
	return c.config.Prefix
}

// GetSchema returns the schema configured.
func (c *Core) GetSchema() string {
	schema := c.schema
	if schema == "" {
		schema = c.db.GetConfig().Name
	}
	return schema
}

func parseConfigNodeLink(node *ConfigNode) (*ConfigNode, error) {
	var (
		link  = node.Link
		match []string
	)
	if link != "" {
		// To be compatible with old configuration,
		// it checks and converts the link to new configuration.
		if node.Type != "" && !gstr.HasPrefix(link, node.Type+":") {
			link = fmt.Sprintf("%s:%s", node.Type, link)
		}
		match, _ = gregex.MatchString(linkPattern, link)
		if len(match) <= 5 {
			return nil, gerror.NewCodef(
				gcode.CodeInvalidParameter,
				`invalid link configuration: %s, shuold be pattern like: %s`,
				link, linkPatternDescription,
			)
		}
		node.Type = match[1]
		node.User = match[2]
		node.Pass = match[3]
		node.Protocol = match[4]
		array := gstr.Split(match[5], ":")
		if node.Protocol == "file" {
			node.Name = match[5]
		} else {
			if len(array) == 2 {
				// link with port.
				node.Host = array[0]
				node.Port = array[1]
			} else {
				// link without port.
				node.Host = array[0]
			}
			node.Name = match[6]
		}
		if len(match) > 6 && match[7] != "" {
			node.Extra = match[7]
		}
	}
	if node.Extra != "" {
		if m, _ := gstr.Parse(node.Extra); len(m) > 0 {
			_ = gconv.Struct(m, &node)
		}
	}
	// Default value checks.
	if node.Charset == "" {
		node.Charset = defaultCharset
	}
	if node.Protocol == "" {
		node.Protocol = defaultProtocol
	}
	return node, nil
}
