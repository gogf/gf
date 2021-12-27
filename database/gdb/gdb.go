// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gdb provides ORM features for popular relationship databases.
package gdb

import (
	"context"
	"database/sql"
	"time"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/container/gtype"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/util/grand"
)

// DB defines the interfaces for ORM operations.
type DB interface {
	// ===========================================================================
	// Model creation.
	// ===========================================================================

	// Model creates and returns a new ORM model from given schema.
	// The parameter `table` can be more than one table names, and also alias name, like:
	// 1. Model names:
	//    Model("user")
	//    Model("user u")
	//    Model("user, user_detail")
	//    Model("user u, user_detail ud")
	// 2. Model name with alias: Model("user", "u")
	// Also see Core.Model.
	Model(tableNameOrStruct ...interface{}) *Model

	// Raw creates and returns a model based on a raw sql not a table.
	Raw(rawSql string, args ...interface{}) *Model

	// Schema creates and returns a schema.
	// Also see Core.Schema.
	Schema(schema string) *Schema

	// With creates and returns an ORM model based on meta data of given object.
	// Also see Core.With.
	With(objects ...interface{}) *Model

	// Open creates a raw connection object for database with given node configuration.
	// Note that it is not recommended using the this function manually.
	// Also see DriverMysql.Open.
	Open(config *ConfigNode) (*sql.DB, error)

	// Ctx is a chaining function, which creates and returns a new DB that is a shallow copy
	// of current DB object and with given context in it.
	// Also see Core.Ctx.
	Ctx(ctx context.Context) DB

	// Close closes the database and prevents new queries from starting.
	// Close then waits for all queries that have started processing on the server
	// to finish.
	//
	// It is rare to Close a DB, as the DB handle is meant to be
	// long-lived and shared between many goroutines.
	Close(ctx context.Context) error

	// ===========================================================================
	// Query APIs.
	// ===========================================================================

	Query(ctx context.Context, sql string, args ...interface{}) (Result, error)    // See Core.Query.
	Exec(ctx context.Context, sql string, args ...interface{}) (sql.Result, error) // See Core.Exec.
	Prepare(ctx context.Context, sql string, execOnMaster ...bool) (*Stmt, error)  // See Core.Prepare.

	// ===========================================================================
	// Common APIs for CURD.
	// ===========================================================================

	Insert(ctx context.Context, table string, data interface{}, batch ...int) (sql.Result, error)                               // See Core.Insert.
	InsertIgnore(ctx context.Context, table string, data interface{}, batch ...int) (sql.Result, error)                         // See Core.InsertIgnore.
	InsertAndGetId(ctx context.Context, table string, data interface{}, batch ...int) (int64, error)                            // See Core.InsertAndGetId.
	Replace(ctx context.Context, table string, data interface{}, batch ...int) (sql.Result, error)                              // See Core.Replace.
	Save(ctx context.Context, table string, data interface{}, batch ...int) (sql.Result, error)                                 // See Core.Save.
	Update(ctx context.Context, table string, data interface{}, condition interface{}, args ...interface{}) (sql.Result, error) // See Core.Update.
	Delete(ctx context.Context, table string, condition interface{}, args ...interface{}) (sql.Result, error)                   // See Core.Delete.

	// ===========================================================================
	// Internal APIs for CURD, which can be overwritten by custom CURD implements.
	// ===========================================================================

	DoGetAll(ctx context.Context, link Link, sql string, args ...interface{}) (result Result, err error)                                           // See Core.DoGetAll.
	DoInsert(ctx context.Context, link Link, table string, data List, option DoInsertOption) (result sql.Result, err error)                        // See Core.DoInsert.
	DoUpdate(ctx context.Context, link Link, table string, data interface{}, condition string, args ...interface{}) (result sql.Result, err error) // See Core.DoUpdate.
	DoDelete(ctx context.Context, link Link, table string, condition string, args ...interface{}) (result sql.Result, err error)                   // See Core.DoDelete.
	DoQuery(ctx context.Context, link Link, sql string, args ...interface{}) (result Result, err error)                                            // See Core.DoQuery.
	DoExec(ctx context.Context, link Link, sql string, args ...interface{}) (result sql.Result, err error)                                         // See Core.DoExec.
	DoFilter(ctx context.Context, link Link, sql string, args []interface{}) (newSql string, newArgs []interface{}, err error)                     // See Core.DoFilter.
	DoCommit(ctx context.Context, in DoCommitInput) (out *DoCommitOutput, err error)                                                               // See Core.DoCommit.
	DoPrepare(ctx context.Context, link Link, sql string) (*Stmt, error)                                                                           // See Core.DoPrepare.

	// ===========================================================================
	// Query APIs for convenience purpose.
	// ===========================================================================

	GetAll(ctx context.Context, sql string, args ...interface{}) (Result, error)                // See Core.GetAll.
	GetOne(ctx context.Context, sql string, args ...interface{}) (Record, error)                // See Core.GetOne.
	GetValue(ctx context.Context, sql string, args ...interface{}) (Value, error)               // See Core.GetValue.
	GetArray(ctx context.Context, sql string, args ...interface{}) ([]Value, error)             // See Core.GetArray.
	GetCount(ctx context.Context, sql string, args ...interface{}) (int, error)                 // See Core.GetCount.
	GetScan(ctx context.Context, objPointer interface{}, sql string, args ...interface{}) error // See Core.GetScan.
	Union(unions ...*Model) *Model                                                              // See Core.Union.
	UnionAll(unions ...*Model) *Model                                                           // See Core.UnionAll.

	// ===========================================================================
	// Master/Slave specification support.
	// ===========================================================================

	Master(schema ...string) (*sql.DB, error) // See Core.Master.
	Slave(schema ...string) (*sql.DB, error)  // See Core.Slave.

	// ===========================================================================
	// Ping-Pong.
	// ===========================================================================

	PingMaster() error // See Core.PingMaster.
	PingSlave() error  // See Core.PingSlave.

	// ===========================================================================
	// Transaction.
	// ===========================================================================

	Begin(ctx context.Context) (*TX, error)                                           // See Core.Begin.
	Transaction(ctx context.Context, f func(ctx context.Context, tx *TX) error) error // See Core.Transaction.

	// ===========================================================================
	// Configuration methods.
	// ===========================================================================

	GetCache() *gcache.Cache            // See Core.GetCache.
	SetDebug(debug bool)                // See Core.SetDebug.
	GetDebug() bool                     // See Core.GetDebug.
	SetSchema(schema string)            // See Core.SetSchema.
	GetSchema() string                  // See Core.GetSchema.
	GetPrefix() string                  // See Core.GetPrefix.
	GetGroup() string                   // See Core.GetGroup.
	SetDryRun(enabled bool)             // See Core.SetDryRun.
	GetDryRun() bool                    // See Core.GetDryRun.
	SetLogger(logger *glog.Logger)      // See Core.SetLogger.
	GetLogger() *glog.Logger            // See Core.GetLogger.
	GetConfig() *ConfigNode             // See Core.GetConfig.
	SetMaxIdleConnCount(n int)          // See Core.SetMaxIdleConnCount.
	SetMaxOpenConnCount(n int)          // See Core.SetMaxOpenConnCount.
	SetMaxConnLifeTime(d time.Duration) // See Core.SetMaxConnLifeTime.

	// ===========================================================================
	// Utility methods.
	// ===========================================================================

	GetCtx() context.Context                                                                         // See Core.GetCtx.
	GetCore() *Core                                                                                  // See Core.GetCore
	GetChars() (charLeft string, charRight string)                                                   // See Core.GetChars.
	Tables(ctx context.Context, schema ...string) (tables []string, err error)                       // See Core.Tables.
	TableFields(ctx context.Context, table string, schema ...string) (map[string]*TableField, error) // See Core.TableFields.
	ConvertDataForRecord(ctx context.Context, data interface{}) map[string]interface{}               // See Core.ConvertDataForRecord
	FilteredLink() string                                                                            // FilteredLink is used for filtering sensitive information in `Link` configuration before output it to tracing server.
}

// Core is the base struct for database management.
type Core struct {
	db     DB              // DB interface object.
	ctx    context.Context // Context for chaining operation only. Do not set a default value in Core initialization.
	group  string          // Configuration group name.
	debug  *gtype.Bool     // Enable debug mode for the database, which can be changed in runtime.
	cache  *gcache.Cache   // Cache manager, SQL result cache only.
	links  *gmap.StrAnyMap // links caches all created links by node.
	schema *gtype.String   // Custom schema for this object.
	logger *glog.Logger    // Logger for logging functionality.
	config *ConfigNode     // Current config node.
}

// DoCommitInput is the input parameters for function DoCommit.
type DoCommitInput struct {
	Stmt *sql.Stmt
	Link Link
	Sql  string
	Args []interface{}
	Type string
}

// DoCommitOutput is the output parameters for function DoCommit.
type DoCommitOutput struct {
	Row    *sql.Row   // Row is the result of Stmt.QueryRowContext.
	Rows   *sql.Rows  // Rows is the result of query statement.
	Result sql.Result // Result is the result of exec statement.
}

// Driver is the interface for integrating sql drivers into package gdb.
type Driver interface {
	// New creates and returns a database object for specified database server.
	New(core *Core, node *ConfigNode) (DB, error)
}

// Link is a common database function wrapper interface.
type Link interface {
	Query(sql string, args ...interface{}) (*sql.Rows, error)
	Exec(sql string, args ...interface{}) (sql.Result, error)
	Prepare(sql string) (*sql.Stmt, error)
	QueryContext(ctx context.Context, sql string, args ...interface{}) (*sql.Rows, error)
	ExecContext(ctx context.Context, sql string, args ...interface{}) (sql.Result, error)
	PrepareContext(ctx context.Context, sql string) (*sql.Stmt, error)
	IsTransaction() bool
}

// Sql is the sql recording struct.
type Sql struct {
	Sql           string        // SQL string(may contain reserved char '?').
	Type          string        // SQL operation type.
	Args          []interface{} // Arguments for this sql.
	Format        string        // Formatted sql which contains arguments in the sql.
	Error         error         // Execution result.
	Start         int64         // Start execution timestamp in milliseconds.
	End           int64         // End execution timestamp in milliseconds.
	Group         string        // Group is the group name of the configuration that the sql is executed from.
	IsTransaction bool          // IsTransaction marks whether this sql is executed in transaction.
	RowsAffected  int64         // RowsAffected marks retrieved or affected number with current sql statement.
}

// DoInsertOption is the input struct for function DoInsert.
type DoInsertOption struct {
	OnDuplicateStr string
	OnDuplicateMap map[string]interface{}
	InsertOption   int // Insert operation.
	BatchCount     int // Batch count for batch inserting.
}

// TableField is the struct for table field.
type TableField struct {
	Index   int         // For ordering purpose as map is unordered.
	Name    string      // Field name.
	Type    string      // Field type.
	Null    bool        // Field can be null or not.
	Key     string      // The index information(empty if it's not an index).
	Default interface{} // Default value for the field.
	Extra   string      // Extra information.
	Comment string      // Field comment.
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
	defaultModelSafe        = false
	defaultCharset          = `utf8`
	queryTypeNormal         = 0
	queryTypeCount          = 1
	unionTypeNormal         = 0
	unionTypeAll            = 1
	insertOptionDefault     = 0
	insertOptionReplace     = 1
	insertOptionSave        = 2
	insertOptionIgnore      = 3
	defaultBatchNumber      = 10               // Per count for batch insert/replace/save.
	defaultMaxIdleConnCount = 10               // Max idle connection count in pool.
	defaultMaxOpenConnCount = 0                // Max open connection count in pool. Default is no limit.
	defaultMaxConnLifeTime  = 30 * time.Second // Max lifetime for per connection in pool in seconds.
	ctxTimeoutTypeExec      = iota
	ctxTimeoutTypeQuery
	ctxTimeoutTypePrepare
	commandEnvKeyForDryRun = "gf.gdb.dryrun"
	sqlTypeBegin           = `DB.Begin`
	sqlTypeTXCommit        = `TX.Commit`
	sqlTypeTXRollback      = `TX.Rollback`
	sqlTypeQueryContext    = `DB.QueryContext`
	sqlTypeExecContext     = `DB.ExecContext`
	sqlTypePrepareContext  = `DB.PrepareContext`
	modelForDaoSuffix      = `ForDao`
	dbRoleSlave            = `slave`
)

const (
	DoCommitTypeExecContext         = "ExecContext"
	DoCommitTypeQueryContext        = "QueryContext"
	DoCommitTypeStmtExecContext     = "Statement.ExecContext"
	DoCommitTypeStmtQueryContext    = "Statement.QueryContext"
	DoCommitTypeStmtQueryRowContext = "Statement.QueryRowContext"
)

var (
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
	regularFieldNameRegPattern = `^[\w\.\-]+$`

	// regularFieldNameWithoutDotRegPattern is similar to regularFieldNameRegPattern but not allows '.'.
	// Note that, although some databases allow char '.' in the field name, but it here does not allow '.'
	// in the field name as it conflicts with "db.table.field" pattern in SOME situations.
	regularFieldNameWithoutDotRegPattern = `^[\w\-]+$`

	// tableFieldsMap caches the table information retrieved from database.
	tableFieldsMap = gmap.New(true)

	// allDryRun sets dry-run feature for all database connections.
	// It is commonly used for command options for convenience.
	allDryRun = false
)

func init() {
	// allDryRun is initialized from environment or command options.
	allDryRun = gcmd.GetOptWithEnv(commandEnvKeyForDryRun, false).Bool()
}

// Register registers custom database driver to gdb.
func Register(name string, driver Driver) error {
	driverMap[name] = driver
	return nil
}

// New creates and returns an ORM object with global configurations.
// The parameter `name` specifies the configuration group name,
// which is DefaultGroupName in default.
func New(group ...string) (db DB, err error) {
	groupName := configs.group
	if len(group) > 0 && group[0] != "" {
		groupName = group[0]
	}
	configs.RLock()
	defer configs.RUnlock()

	if len(configs.config) < 1 {
		return nil, gerror.NewCode(
			gcode.CodeInvalidConfiguration,
			"database configuration is empty, please set the database configuration before using",
		)
	}
	if _, ok := configs.config[groupName]; ok {
		if node, err := getConfigNodeByGroup(groupName, true); err == nil {
			c := &Core{
				group:  groupName,
				debug:  gtype.NewBool(),
				cache:  gcache.New(),
				links:  gmap.NewStrAnyMap(true),
				schema: gtype.NewString(),
				logger: glog.New(),
				config: node,
			}
			if v, ok := driverMap[node.Type]; ok {
				c.db, err = v.New(c, node)
				if err != nil {
					return nil, err
				}
				return c.db, nil
			} else {
				return nil, gerror.NewCodef(
					gcode.CodeInvalidConfiguration,
					`cannot find database driver for specified database type "%s", did you misspell type name "%s" or forget importing the database driver?`,
					node.Type, node.Type,
				)
			}
		} else {
			return nil, err
		}
	} else {
		return nil, gerror.NewCodef(
			gcode.CodeInvalidConfiguration,
			`database configuration node "%s" is not found, did you misspell group name "%s" or miss the database configuration?`,
			groupName, groupName,
		)
	}
}

// Instance returns an instance for DB operations.
// The parameter `name` specifies the configuration group name,
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
// The parameter `master` specifies whether retrieving a master node, or else a slave node
// if master-slave configured.
func getConfigNodeByGroup(group string, master bool) (*ConfigNode, error) {
	if list, ok := configs.config[group]; ok {
		// Separates master and slave configuration nodes array.
		var (
			masterList = make(ConfigGroup, 0)
			slaveList  = make(ConfigGroup, 0)
		)
		for i := 0; i < len(list); i++ {
			if list[i].Role == dbRoleSlave {
				slaveList = append(slaveList, list[i])
			} else {
				masterList = append(masterList, list[i])
			}
		}
		if len(masterList) < 1 {
			return nil, gerror.NewCode(
				gcode.CodeInvalidConfiguration,
				"at least one master node configuration's need to make sense",
			)
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
		return nil, gerror.NewCodef(
			gcode.CodeInvalidConfiguration,
			"empty database configuration for item name '%s'",
			group,
		)
	}
}

// getConfigNodeByWeight calculates the configuration weights and randomly returns a node.
//
// Calculation algorithm brief:
// 1. If we have 2 nodes, and their weights are both 1, then the weight range is [0, 199];
// 2. Node1 weight range is [0, 99], and node2 weight range is [100, 199], ratio is 1:1;
// 3. If the random number is 99, it then chooses and returns node1;.
func getConfigNodeByWeight(cg ConfigGroup) *ConfigNode {
	if len(cg) < 2 {
		return &cg[0]
	}
	var total int
	for i := 0; i < len(cg); i++ {
		total += cg[i].Weight * 100
	}
	// If total is 0 means all the nodes have no weight attribute configured.
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
		// fmt.Printf("r: %d, min: %d, max: %d\n", r, min, max)
		if r >= min && r < max {
			return &cg[i]
		} else {
			min = max
		}
	}
	return nil
}

// getSqlDb retrieves and returns a underlying database connection object.
// The parameter `master` specifies whether retrieves master node connection if
// master-slave nodes are configured.
func (c *Core) getSqlDb(master bool, schema ...string) (sqlDb *sql.DB, err error) {
	// Load balance.
	node, err := getConfigNodeByGroup(c.group, master)
	if err != nil {
		return nil, err
	}
	// Default value checks.
	if node.Charset == "" {
		node.Charset = defaultCharset
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
	v := c.links.GetOrSetFuncLock(node.String(), func() interface{} {
		intlog.Printf(
			c.db.GetCtx(),
			`open new connection, master:%#v, config:%#v, node:%#v`,
			master, c.config, node,
		)
		defer func() {
			if err != nil {
				intlog.Printf(c.db.GetCtx(), `open new connection failed: %v, %#v`, err, node)
			} else {
				intlog.Printf(
					c.db.GetCtx(),
					`open new connection success, master:%#v, config:%#v, node:%#v`,
					master, c.config, node,
				)
			}
		}()

		if sqlDb, err = c.db.Open(node); err != nil {
			return nil
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
		if c.config.MaxConnLifeTime > 0 {
			sqlDb.SetConnMaxLifetime(c.config.MaxConnLifeTime)
		} else {
			sqlDb.SetConnMaxLifetime(defaultMaxConnLifeTime)
		}
		return sqlDb
	})
	if v != nil && sqlDb == nil {
		sqlDb = v.(*sql.DB)
	}
	if node.Debug {
		c.db.SetDebug(node.Debug)
	}
	if node.DryRun {
		c.db.SetDryRun(node.DryRun)
	}
	return
}
