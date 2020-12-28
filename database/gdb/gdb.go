// Copyright GoFrame Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gdb provides ORM features for popular relationship databases.
package gdb

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/os/gcmd"
	"time"

	"github.com/gogf/gf/container/gvar"
	"github.com/gogf/gf/internal/intlog"

	"github.com/gogf/gf/os/glog"

	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/container/gtype"
	"github.com/gogf/gf/os/gcache"
	"github.com/gogf/gf/util/grand"
)

// DB defines the interfaces for ORM operations.
type DB interface {
	// ===========================================================================
	// Model creation.
	// ===========================================================================

	// The DB interface is designed not only for
	// relational databases but also for NoSQL databases in the future. The name
	// "Table" is not proper for that purpose any more.
	Table(table ...string) *Model
	Model(table ...string) *Model
	Schema(schema string) *Schema

	// Open creates a raw connection object for database with given node configuration.
	// Note that it is not recommended using the this function manually.
	Open(config *ConfigNode) (*sql.DB, error)

	// Ctx is a chaining function, which creates and returns a new DB that is a shallow copy
	// of current DB object and with given context in it.
	// Note that this returned DB object can be used only once, so do not assign it to
	// a global or package variable for long using.
	Ctx(ctx context.Context) DB

	// ===========================================================================
	// Query APIs.
	// ===========================================================================

	Query(sql string, args ...interface{}) (*sql.Rows, error)
	Exec(sql string, args ...interface{}) (sql.Result, error)
	Prepare(sql string, execOnMaster ...bool) (*sql.Stmt, error)

	// ===========================================================================
	// Common APIs for CURD.
	// ===========================================================================

	Insert(table string, data interface{}, batch ...int) (sql.Result, error)
	InsertIgnore(table string, data interface{}, batch ...int) (sql.Result, error)
	Replace(table string, data interface{}, batch ...int) (sql.Result, error)
	Save(table string, data interface{}, batch ...int) (sql.Result, error)

	BatchInsert(table string, list interface{}, batch ...int) (sql.Result, error)
	BatchReplace(table string, list interface{}, batch ...int) (sql.Result, error)
	BatchSave(table string, list interface{}, batch ...int) (sql.Result, error)

	Update(table string, data interface{}, condition interface{}, args ...interface{}) (sql.Result, error)
	Delete(table string, condition interface{}, args ...interface{}) (sql.Result, error)

	// ===========================================================================
	// Internal APIs for CURD, which can be overwrote for custom CURD implements.
	// ===========================================================================

	DoQuery(link Link, sql string, args ...interface{}) (rows *sql.Rows, err error)
	DoGetAll(link Link, sql string, args ...interface{}) (result Result, err error)
	DoExec(link Link, sql string, args ...interface{}) (result sql.Result, err error)
	DoPrepare(link Link, sql string) (*sql.Stmt, error)
	DoInsert(link Link, table string, data interface{}, option int, batch ...int) (result sql.Result, err error)
	DoBatchInsert(link Link, table string, list interface{}, option int, batch ...int) (result sql.Result, err error)
	DoUpdate(link Link, table string, data interface{}, condition string, args ...interface{}) (result sql.Result, err error)
	DoDelete(link Link, table string, condition string, args ...interface{}) (result sql.Result, err error)

	// ===========================================================================
	// Query APIs for convenience purpose.
	// ===========================================================================

	GetAll(sql string, args ...interface{}) (Result, error)
	GetOne(sql string, args ...interface{}) (Record, error)
	GetValue(sql string, args ...interface{}) (Value, error)
	GetArray(sql string, args ...interface{}) ([]Value, error)
	GetCount(sql string, args ...interface{}) (int, error)
	GetStruct(objPointer interface{}, sql string, args ...interface{}) error
	GetStructs(objPointerSlice interface{}, sql string, args ...interface{}) error
	GetScan(objPointer interface{}, sql string, args ...interface{}) error

	// ===========================================================================
	// Master/Slave specification support.
	// ===========================================================================

	Master() (*sql.DB, error)
	Slave() (*sql.DB, error)

	// ===========================================================================
	// Ping-Pong.
	// ===========================================================================

	PingMaster() error
	PingSlave() error

	// ===========================================================================
	// Transaction.
	// ===========================================================================

	Begin() (*TX, error)
	Transaction(f func(tx *TX) error) (err error)

	// ===========================================================================
	// Configuration methods.
	// ===========================================================================

	GetCache() *gcache.Cache
	SetDebug(debug bool)
	GetDebug() bool
	SetSchema(schema string)
	GetSchema() string
	GetPrefix() string
	GetGroup() string
	SetDryRun(dryrun bool)
	GetDryRun() bool
	SetLogger(logger *glog.Logger)
	GetLogger() *glog.Logger
	GetConfig() *ConfigNode
	SetMaxIdleConnCount(n int)
	SetMaxOpenConnCount(n int)
	SetMaxConnLifetime(d time.Duration)

	// ===========================================================================
	// Utility methods.
	// ===========================================================================

	GetCtx() context.Context
	GetChars() (charLeft string, charRight string)
	GetMaster(schema ...string) (*sql.DB, error)
	GetSlave(schema ...string) (*sql.DB, error)
	QuoteWord(s string) string
	QuoteString(s string) string
	QuotePrefixTableName(table string) string
	Tables(schema ...string) (tables []string, err error)
	TableFields(table string, schema ...string) (map[string]*TableField, error)
	HasTable(name string) (bool, error)

	// HandleSqlBeforeCommit is a hook function, which deals with the sql string before
	// it's committed to underlying driver. The parameter <link> specifies the current
	// database connection operation object. You can modify the sql string <sql> and its
	// arguments <args> as you wish before they're committed to driver.
	HandleSqlBeforeCommit(link Link, sql string, args []interface{}) (string, []interface{})

	// ===========================================================================
	// Internal methods, for internal usage purpose, you do not need consider it.
	// ===========================================================================

	mappingAndFilterData(schema, table string, data map[string]interface{}, filter bool) (map[string]interface{}, error)
	convertFieldValueToLocalValue(fieldValue interface{}, fieldType string) interface{}
	convertRowsToResult(rows *sql.Rows) (Result, error)
}

// Core is the base struct for database management.
type Core struct {
	DB     DB              // DB interface object.
	group  string          // Configuration group name.
	debug  *gtype.Bool     // Enable debug mode for the database, which can be changed in runtime.
	cache  *gcache.Cache   // Cache manager, SQL result cache only.
	schema *gtype.String   // Custom schema for this object.
	logger *glog.Logger    // Logger.
	config *ConfigNode     // Current config node.
	ctx    context.Context // Context for chaining operation only.
}

// Driver is the interface for integrating sql drivers into package gdb.
type Driver interface {
	// New creates and returns a database object for specified database server.
	New(core *Core, node *ConfigNode) (DB, error)
}

// Sql is the sql recording struct.
type Sql struct {
	Sql    string        // SQL string(may contain reserved char '?').
	Args   []interface{} // Arguments for this sql.
	Format string        // Formatted sql which contains arguments in the sql.
	Error  error         // Execution result.
	Start  int64         // Start execution timestamp in milliseconds.
	End    int64         // End execution timestamp in milliseconds.
	Group  string        // Group is the group name of the configuration that the sql is executed from.
}

// TableField is the struct for table field.
type TableField struct {
	Index   int         // For ordering purpose as map is unordered.
	Name    string      // Field name.
	Type    string      // Field type.
	Null    bool        // Field can be null or not.
	Key     string      // The index information(empty if it's not a index).
	Default interface{} // Default value for the field.
	Extra   string      // Extra information.
	Comment string      // Comment.
}

// Link is a common database function wrapper interface.
type Link interface {
	Query(sql string, args ...interface{}) (*sql.Rows, error)
	Exec(sql string, args ...interface{}) (sql.Result, error)
	Prepare(sql string) (*sql.Stmt, error)
}

// Counter  is the type for update count.
type Counter struct {
	Field string
	Value float64
}

type (
	Raw    string                   // Raw is a raw sql that will not be treated as argument but as a direct sql part.
	Value  = *gvar.Var              // Value is the field value type.
	Record map[string]Value         // Record is the row record of the table.
	Result []Record                 // Result is the row record array.
	Map    = map[string]interface{} // Map is alias of map[string]interface{}, which is the most common usage map type.
	List   = []Map                  // List is type of map array.
)

const (
	insertOptionDefault     = 0
	insertOptionReplace     = 1
	insertOptionSave        = 2
	insertOptionIgnore      = 3
	defaultBatchNumber      = 10               // Per count for batch insert/replace/save.
	defaultMaxIdleConnCount = 10               // Max idle connection count in pool.
	defaultMaxOpenConnCount = 100              // Max open connection count in pool.
	defaultMaxConnLifeTime  = 30 * time.Second // Max life time for per connection in pool in seconds.
)

var (
	// ErrNoRows is alias of sql.ErrNoRows.
	ErrNoRows = sql.ErrNoRows

	// instances is the management map for instances.
	instances = gmap.NewStrAnyMap(true)

	// driverMap manages all custom registered driver.
	driverMap = map[string]Driver{
		"mysql":  &DriverMysql{},
		"mssql":  &DriverMssql{},
		"pgsql":  &DriverPgsql{},
		"oracle": &DriverOracle{},
		"sqlite": &DriverSqlite{},
	}

	// lastOperatorRegPattern is the regular expression pattern for a string
	// which has operator at its tail.
	lastOperatorRegPattern = `[<>=]+\s*$`

	// regularFieldNameRegPattern is the regular expression pattern for a string
	// which is a regular field name of table.
	regularFieldNameRegPattern = `^[\w\.\-\_]+$`

	// internalCache is the memory cache for internal usage.
	internalCache = gcache.New()

	// allDryRun sets dry-run feature for all database connections.
	// It is commonly used for command options for convenience.
	allDryRun = false
)

func init() {
	// allDryRun is initialized from environment or command options.
	allDryRun = gcmd.GetWithEnv("gf.gdb.dryrun", false).Bool()
}

// Register registers custom database driver to gdb.
func Register(name string, driver Driver) error {
	driverMap[name] = driver
	return nil
}

// New creates and returns an ORM object with global configurations.
// The parameter <name> specifies the configuration group name,
// which is DefaultGroupName in default.
func New(group ...string) (db DB, err error) {
	groupName := configs.group
	if len(group) > 0 && group[0] != "" {
		groupName = group[0]
	}
	configs.RLock()
	defer configs.RUnlock()

	if len(configs.config) < 1 {
		return nil, gerror.New("empty database configuration")
	}
	if _, ok := configs.config[groupName]; ok {
		if node, err := getConfigNodeByGroup(groupName, true); err == nil {
			c := &Core{
				group:  groupName,
				debug:  gtype.NewBool(),
				cache:  gcache.New(),
				schema: gtype.NewString(),
				logger: glog.New(),
				config: node,
			}
			if v, ok := driverMap[node.Type]; ok {
				c.DB, err = v.New(c, node)
				if err != nil {
					return nil, err
				}
				return c.DB, nil
			} else {
				return nil, gerror.New(fmt.Sprintf(`unsupported database type "%s"`, node.Type))
			}
		} else {
			return nil, err
		}
	} else {
		return nil, gerror.New(fmt.Sprintf(`database configuration node "%s" is not found`, groupName))
	}
}

// Instance returns an instance for DB operations.
// The parameter <name> specifies the configuration group name,
// which is DefaultGroupName in default.
func Instance(name ...string) (db DB, err error) {
	group := configs.group
	if len(name) > 0 && name[0] != "" {
		group = name[0]
	}
	v := instances.GetOrSetFuncLock(group, func() interface{} {
		db, err = New(group)
		return db
	})
	if v != nil {
		return v.(DB), nil
	}
	return
}

// getConfigNodeByGroup calculates and returns a configuration node of given group. It
// calculates the value internally using weight algorithm for load balance.
//
// The parameter <master> specifies whether retrieving a master node, or else a slave node
// if master-slave configured.
func getConfigNodeByGroup(group string, master bool) (*ConfigNode, error) {
	if list, ok := configs.config[group]; ok {
		// Separates master and slave configuration nodes array.
		masterList := make(ConfigGroup, 0)
		slaveList := make(ConfigGroup, 0)
		for i := 0; i < len(list); i++ {
			if list[i].Role == "slave" {
				slaveList = append(slaveList, list[i])
			} else {
				masterList = append(masterList, list[i])
			}
		}
		if len(masterList) < 1 {
			return nil, gerror.New("at least one master node configuration's need to make sense")
		}
		if len(slaveList) < 1 {
			slaveList = masterList
		}
		if master {
			return getConfigNodeByWeight(masterList), nil
		} else {
			return getConfigNodeByWeight(slaveList), nil
		}
	} else {
		return nil, gerror.New(fmt.Sprintf("empty database configuration for item name '%s'", group))
	}
}

// getConfigNodeByWeight calculates the configuration weights and randomly returns a node.
//
// Calculation algorithm brief:
// 1. If we have 2 nodes, and their weights are both 1, then the weight range is [0, 199];
// 2. Node1 weight range is [0, 99], and node2 weight range is [100, 199], ratio is 1:1;
// 3. If the random number is 99, it then chooses and returns node1;
func getConfigNodeByWeight(cg ConfigGroup) *ConfigNode {
	if len(cg) < 2 {
		return &cg[0]
	}
	var total int
	for i := 0; i < len(cg); i++ {
		total += cg[i].Weight * 100
	}
	// If total is 0 means all of the nodes have no weight attribute configured.
	// It then defaults each node's weight attribute to 1.
	if total == 0 {
		for i := 0; i < len(cg); i++ {
			cg[i].Weight = 1
			total += cg[i].Weight * 100
		}
	}
	// Exclude the right border value.
	r := grand.N(0, total-1)
	min := 0
	max := 0
	for i := 0; i < len(cg); i++ {
		max = min + cg[i].Weight*100
		//fmt.Printf("r: %d, min: %d, max: %d\n", r, min, max)
		if r >= min && r < max {
			return &cg[i]
		} else {
			min = max
		}
	}
	return nil
}

// getSqlDb retrieves and returns a underlying database connection object.
// The parameter <master> specifies whether retrieves master node connection if
// master-slave nodes are configured.
func (c *Core) getSqlDb(master bool, schema ...string) (sqlDb *sql.DB, err error) {
	// Load balance.
	node, err := getConfigNodeByGroup(c.group, master)
	if err != nil {
		return nil, err
	}
	// Default value checks.
	if node.Charset == "" {
		node.Charset = "utf8"
	}
	// Changes the schema.
	nodeSchema := c.schema.Val()
	if len(schema) > 0 && schema[0] != "" {
		nodeSchema = schema[0]
	}
	if nodeSchema != "" {
		// Value copy.
		n := *node
		n.Name = nodeSchema
		node = &n
	}
	// Cache the underlying connection pool object by node.
	v, _ := internalCache.GetOrSetFuncLock(node.String(), func() (interface{}, error) {
		sqlDb, err = c.DB.Open(node)
		if err != nil {
			intlog.Printf("DB open failed: %v, %+v", err, node)
			return nil, err
		}
		if c.config.MaxIdleConnCount > 0 {
			sqlDb.SetMaxIdleConns(c.config.MaxIdleConnCount)
		} else {
			sqlDb.SetMaxIdleConns(defaultMaxIdleConnCount)
		}
		if c.config.MaxOpenConnCount > 0 {
			sqlDb.SetMaxOpenConns(c.config.MaxOpenConnCount)
		} else {
			sqlDb.SetMaxOpenConns(defaultMaxOpenConnCount)
		}
		if c.config.MaxConnLifetime > 0 {
			// Automatically checks whether MaxConnLifetime is configured using string like: "30s", "60s", etc.
			// Or else it is configured just using number, which means value in seconds.
			if c.config.MaxConnLifetime > time.Second {
				sqlDb.SetConnMaxLifetime(c.config.MaxConnLifetime)
			} else {
				sqlDb.SetConnMaxLifetime(c.config.MaxConnLifetime * time.Second)
			}
		} else {
			sqlDb.SetConnMaxLifetime(defaultMaxConnLifeTime)
		}
		return sqlDb, nil
	}, 0)
	if v != nil && sqlDb == nil {
		sqlDb = v.(*sql.DB)
	}
	if node.Debug {
		c.DB.SetDebug(node.Debug)
	}
	if node.Debug {
		c.DB.SetDryRun(node.DryRun)
	}
	return
}
