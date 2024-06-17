// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gdb provides ORM features for popular relationship databases.
//
// TODO use context.Context as required parameter for all DB operations.
package gdb

import (
	"context"
	"database/sql"
	"time"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/container/gtype"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/util/grand"
	"github.com/gogf/gf/v2/util/gutil"
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

	// With creates and returns an ORM model based on metadata of given object.
	// Also see Core.With.
	With(objects ...interface{}) *Model

	// Open creates a raw connection object for database with given node configuration.
	// Note that it is not recommended using the function manually.
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

	DoSelect(ctx context.Context, link Link, sql string, args ...interface{}) (result Result, err error)                                           // See Core.DoSelect.
	DoInsert(ctx context.Context, link Link, table string, data List, option DoInsertOption) (result sql.Result, err error)                        // See Core.DoInsert.
	DoUpdate(ctx context.Context, link Link, table string, data interface{}, condition string, args ...interface{}) (result sql.Result, err error) // See Core.DoUpdate.
	DoDelete(ctx context.Context, link Link, table string, condition string, args ...interface{}) (result sql.Result, err error)                   // See Core.DoDelete.

	DoQuery(ctx context.Context, link Link, sql string, args ...interface{}) (result Result, err error)    // See Core.DoQuery.
	DoExec(ctx context.Context, link Link, sql string, args ...interface{}) (result sql.Result, err error) // See Core.DoExec.

	DoFilter(ctx context.Context, link Link, sql string, args []interface{}) (newSql string, newArgs []interface{}, err error) // See Core.DoFilter.
	DoCommit(ctx context.Context, in DoCommitInput) (out DoCommitOutput, err error)                                            // See Core.DoCommit.

	DoPrepare(ctx context.Context, link Link, sql string) (*Stmt, error) // See Core.DoPrepare.

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

	Begin(ctx context.Context) (TX, error)                                           // See Core.Begin.
	Transaction(ctx context.Context, f func(ctx context.Context, tx TX) error) error // See Core.Transaction.

	// ===========================================================================
	// Configuration methods.
	// ===========================================================================

	GetCache() *gcache.Cache            // See Core.GetCache.
	SetDebug(debug bool)                // See Core.SetDebug.
	GetDebug() bool                     // See Core.GetDebug.
	GetSchema() string                  // See Core.GetSchema.
	GetPrefix() string                  // See Core.GetPrefix.
	GetGroup() string                   // See Core.GetGroup.
	SetDryRun(enabled bool)             // See Core.SetDryRun.
	GetDryRun() bool                    // See Core.GetDryRun.
	SetLogger(logger glog.ILogger)      // See Core.SetLogger.
	GetLogger() glog.ILogger            // See Core.GetLogger.
	GetConfig() *ConfigNode             // See Core.GetConfig.
	SetMaxIdleConnCount(n int)          // See Core.SetMaxIdleConnCount.
	SetMaxOpenConnCount(n int)          // See Core.SetMaxOpenConnCount.
	SetMaxConnLifeTime(d time.Duration) // See Core.SetMaxConnLifeTime.

	// ===========================================================================
	// Utility methods.
	// ===========================================================================

	Stats(ctx context.Context) []StatsItem                                                                   // See Core.Stats.
	GetCtx() context.Context                                                                                 // See Core.GetCtx.
	GetCore() *Core                                                                                          // See Core.GetCore
	GetChars() (charLeft string, charRight string)                                                           // See Core.GetChars.
	Tables(ctx context.Context, schema ...string) (tables []string, err error)                               // See Core.Tables. The driver must implement this function.
	TableFields(ctx context.Context, table string, schema ...string) (map[string]*TableField, error)         // See Core.TableFields. The driver must implement this function.
	ConvertValueForField(ctx context.Context, fieldType string, fieldValue interface{}) (interface{}, error) // See Core.ConvertValueForField
	ConvertValueForLocal(ctx context.Context, fieldType string, fieldValue interface{}) (interface{}, error) // See Core.ConvertValueForLocal
	CheckLocalTypeForField(ctx context.Context, fieldType string, fieldValue interface{}) (LocalType, error) // See Core.CheckLocalTypeForField
	FormatUpsert(columns []string, list List, option DoInsertOption) (string, error)                         // See Core.DoFormatUpsert
}

// TX defines the interfaces for ORM transaction operations.
type TX interface {
	Link

	Ctx(ctx context.Context) TX
	Raw(rawSql string, args ...interface{}) *Model
	Model(tableNameQueryOrStruct ...interface{}) *Model
	With(object interface{}) *Model

	// ===========================================================================
	// Nested transaction if necessary.
	// ===========================================================================

	Begin() error
	Commit() error
	Rollback() error
	Transaction(ctx context.Context, f func(ctx context.Context, tx TX) error) (err error)

	// ===========================================================================
	// Core method.
	// ===========================================================================

	Query(sql string, args ...interface{}) (result Result, err error)
	Exec(sql string, args ...interface{}) (sql.Result, error)
	Prepare(sql string) (*Stmt, error)

	// ===========================================================================
	// Query.
	// ===========================================================================

	GetAll(sql string, args ...interface{}) (Result, error)
	GetOne(sql string, args ...interface{}) (Record, error)
	GetStruct(obj interface{}, sql string, args ...interface{}) error
	GetStructs(objPointerSlice interface{}, sql string, args ...interface{}) error
	GetScan(pointer interface{}, sql string, args ...interface{}) error
	GetValue(sql string, args ...interface{}) (Value, error)
	GetCount(sql string, args ...interface{}) (int64, error)

	// ===========================================================================
	// CURD.
	// ===========================================================================

	Insert(table string, data interface{}, batch ...int) (sql.Result, error)
	InsertIgnore(table string, data interface{}, batch ...int) (sql.Result, error)
	InsertAndGetId(table string, data interface{}, batch ...int) (int64, error)
	Replace(table string, data interface{}, batch ...int) (sql.Result, error)
	Save(table string, data interface{}, batch ...int) (sql.Result, error)
	Update(table string, data interface{}, condition interface{}, args ...interface{}) (sql.Result, error)
	Delete(table string, condition interface{}, args ...interface{}) (sql.Result, error)

	// ===========================================================================
	// Utility methods.
	// ===========================================================================

	GetCtx() context.Context
	GetDB() DB
	GetSqlTX() *sql.Tx
	IsClosed() bool

	// ===========================================================================
	// Save point feature.
	// ===========================================================================

	SavePoint(point string) error
	RollbackTo(point string) error
}

// StatsItem defines the stats information for a configuration node.
type StatsItem interface {
	// Node returns the configuration node info.
	Node() ConfigNode

	// Stats returns the connection stat for current node.
	Stats() sql.DBStats
}

// Core is the base struct for database management.
type Core struct {
	db            DB              // DB interface object.
	ctx           context.Context // Context for chaining operation only. Do not set a default value in Core initialization.
	group         string          // Configuration group name.
	schema        string          // Custom schema for this object.
	debug         *gtype.Bool     // Enable debug mode for the database, which can be changed in runtime.
	cache         *gcache.Cache   // Cache manager, SQL result cache only.
	links         *gmap.Map       // links caches all created links by node.
	logger        glog.ILogger    // Logger for logging functionality.
	config        *ConfigNode     // Current config node.
	dynamicConfig dynamicConfig   // Dynamic configurations, which can be changed in runtime.
	innerMemCache *gcache.Cache
}

type dynamicConfig struct {
	MaxIdleConnCount int
	MaxOpenConnCount int
	MaxConnLifeTime  time.Duration
}

// DoCommitInput is the input parameters for function DoCommit.
type DoCommitInput struct {
	Db            *sql.DB
	Tx            *sql.Tx
	Stmt          *sql.Stmt
	Link          Link
	Sql           string
	Args          []interface{}
	Type          SqlType
	IsTransaction bool
}

// DoCommitOutput is the output parameters for function DoCommit.
type DoCommitOutput struct {
	Result    sql.Result  // Result is the result of exec statement.
	Records   []Record    // Records is the result of query statement.
	Stmt      *Stmt       // Stmt is the Statement object result for Prepare.
	Tx        TX          // Tx is the transaction object result for Begin.
	RawResult interface{} // RawResult is the underlying result, which might be sql.Result/*sql.Rows/*sql.Row.
}

// Driver is the interface for integrating sql drivers into package gdb.
type Driver interface {
	// New creates and returns a database object for specified database server.
	New(core *Core, node *ConfigNode) (DB, error)
}

// Link is a common database function wrapper interface.
// Note that, any operation using `Link` will have no SQL logging.
type Link interface {
	QueryContext(ctx context.Context, sql string, args ...interface{}) (*sql.Rows, error)
	ExecContext(ctx context.Context, sql string, args ...interface{}) (sql.Result, error)
	PrepareContext(ctx context.Context, sql string) (*sql.Stmt, error)
	IsOnMaster() bool
	IsTransaction() bool
}

// Sql is the sql recording struct.
type Sql struct {
	Sql           string        // SQL string(may contain reserved char '?').
	Type          SqlType       // SQL operation type.
	Args          []interface{} // Arguments for this sql.
	Format        string        // Formatted sql which contains arguments in the sql.
	Error         error         // Execution result.
	Start         int64         // Start execution timestamp in milliseconds.
	End           int64         // End execution timestamp in milliseconds.
	Group         string        // Group is the group name of the configuration that the sql is executed from.
	Schema        string        // Schema is the schema name of the configuration that the sql is executed from.
	IsTransaction bool          // IsTransaction marks whether this sql is executed in transaction.
	RowsAffected  int64         // RowsAffected marks retrieved or affected number with current sql statement.
}

// DoInsertOption is the input struct for function DoInsert.
type DoInsertOption struct {
	OnDuplicateStr string                 // Custom string for `on duplicated` statement.
	OnDuplicateMap map[string]interface{} // Custom key-value map from `OnDuplicateEx` function for `on duplicated` statement.
	OnConflict     []string               // Custom conflict key of upsert clause, if the database needs it.
	InsertOption   InsertOption           // Insert operation in constant value.
	BatchCount     int                    // Batch count for batch inserting.
}

// TableField is the struct for table field.
type TableField struct {
	Index   int         // For ordering purpose as map is unordered.
	Name    string      // Field name.
	Type    string      // Field type. Eg: 'int(10) unsigned', 'varchar(64)'.
	Null    bool        // Field can be null or not.
	Key     string      // The index information(empty if it's not an index). Eg: PRI, MUL.
	Default interface{} // Default value for the field.
	Extra   string      // Extra information. Eg: auto_increment.
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

type CatchSQLManager struct {
	SQLArray *garray.StrArray
	DoCommit bool // DoCommit marks it will be committed to underlying driver or not.
}

const (
	defaultModelSafe                      = false
	defaultCharset                        = `utf8`
	defaultProtocol                       = `tcp`
	unionTypeNormal                       = 0
	unionTypeAll                          = 1
	defaultMaxIdleConnCount               = 10               // Max idle connection count in pool.
	defaultMaxOpenConnCount               = 0                // Max open connection count in pool. Default is no limit.
	defaultMaxConnLifeTime                = 30 * time.Second // Max lifetime for per connection in pool in seconds.
	ctxTimeoutTypeExec                    = 0
	ctxTimeoutTypeQuery                   = 1
	ctxTimeoutTypePrepare                 = 2
	cachePrefixTableFields                = `TableFields:`
	cachePrefixSelectCache                = `SelectCache:`
	commandEnvKeyForDryRun                = "gf.gdb.dryrun"
	modelForDaoSuffix                     = `ForDao`
	dbRoleSlave                           = `slave`
	ctxKeyForDB               gctx.StrKey = `CtxKeyForDB`
	ctxKeyCatchSQL            gctx.StrKey = `CtxKeyCatchSQL`
	ctxKeyInternalProducedSQL gctx.StrKey = `CtxKeyInternalProducedSQL`

	// type:[username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
	linkPattern = `(\w+):([\w\-\$]*):(.*?)@(\w+?)\((.+?)\)/{0,1}([^\?]*)\?{0,1}(.*)`
)

type queryType int

const (
	queryTypeNormal queryType = iota
	queryTypeCount
	queryTypeValue
)

type joinOperator string

const (
	joinOperatorLeft  joinOperator = "LEFT"
	joinOperatorRight joinOperator = "RIGHT"
	joinOperatorInner joinOperator = "INNER"
)

type InsertOption int

const (
	InsertOptionDefault InsertOption = iota
	InsertOptionReplace
	InsertOptionSave
	InsertOptionIgnore
)

const (
	InsertOperationInsert      = "INSERT"
	InsertOperationReplace     = "REPLACE"
	InsertOperationIgnore      = "INSERT IGNORE"
	InsertOnDuplicateKeyUpdate = "ON DUPLICATE KEY UPDATE"
)

type SqlType string

const (
	SqlTypeBegin               SqlType = "DB.Begin"
	SqlTypeTXCommit            SqlType = "TX.Commit"
	SqlTypeTXRollback          SqlType = "TX.Rollback"
	SqlTypeExecContext         SqlType = "DB.ExecContext"
	SqlTypeQueryContext        SqlType = "DB.QueryContext"
	SqlTypePrepareContext      SqlType = "DB.PrepareContext"
	SqlTypeStmtExecContext     SqlType = "DB.Statement.ExecContext"
	SqlTypeStmtQueryContext    SqlType = "DB.Statement.QueryContext"
	SqlTypeStmtQueryRowContext SqlType = "DB.Statement.QueryRowContext"
)

type LocalType string

const (
	LocalTypeUndefined   LocalType = ""
	LocalTypeString      LocalType = "string"
	LocalTypeDate        LocalType = "date"
	LocalTypeDatetime    LocalType = "datetime"
	LocalTypeInt         LocalType = "int"
	LocalTypeUint        LocalType = "uint"
	LocalTypeInt64       LocalType = "int64"
	LocalTypeUint64      LocalType = "uint64"
	LocalTypeIntSlice    LocalType = "[]int"
	LocalTypeInt64Slice  LocalType = "[]int64"
	LocalTypeUint64Slice LocalType = "[]uint64"
	LocalTypeInt64Bytes  LocalType = "int64-bytes"
	LocalTypeUint64Bytes LocalType = "uint64-bytes"
	LocalTypeFloat32     LocalType = "float32"
	LocalTypeFloat64     LocalType = "float64"
	LocalTypeBytes       LocalType = "[]byte"
	LocalTypeBool        LocalType = "bool"
	LocalTypeJson        LocalType = "json"
	LocalTypeJsonb       LocalType = "jsonb"
)

const (
	fieldTypeBinary     = "binary"
	fieldTypeVarbinary  = "varbinary"
	fieldTypeBlob       = "blob"
	fieldTypeTinyblob   = "tinyblob"
	fieldTypeMediumblob = "mediumblob"
	fieldTypeLongblob   = "longblob"
	fieldTypeInt        = "int"
	fieldTypeTinyint    = "tinyint"
	fieldTypeSmallInt   = "small_int"
	fieldTypeSmallint   = "smallint"
	fieldTypeMediumInt  = "medium_int"
	fieldTypeMediumint  = "mediumint"
	fieldTypeSerial     = "serial"
	fieldTypeBigInt     = "big_int"
	fieldTypeBigint     = "bigint"
	fieldTypeBigserial  = "bigserial"
	fieldTypeReal       = "real"
	fieldTypeFloat      = "float"
	fieldTypeDouble     = "double"
	fieldTypeDecimal    = "decimal"
	fieldTypeMoney      = "money"
	fieldTypeNumeric    = "numeric"
	fieldTypeSmallmoney = "smallmoney"
	fieldTypeBool       = "bool"
	fieldTypeBit        = "bit"
	fieldTypeDate       = "date"
	fieldTypeDatetime   = "datetime"
	fieldTypeTimestamp  = "timestamp"
	fieldTypeTimestampz = "timestamptz"
	fieldTypeJson       = "json"
	fieldTypeJsonb      = "jsonb"
)

var (
	// instances is the management map for instances.
	instances = gmap.NewStrAnyMap(true)

	// driverMap manages all custom registered driver.
	driverMap = map[string]Driver{}

	// lastOperatorRegPattern is the regular expression pattern for a string
	// which has operator at its tail.
	lastOperatorRegPattern = `[<>=]+\s*$`

	// regularFieldNameRegPattern is the regular expression pattern for a string
	// which is a regular field name of table.
	regularFieldNameRegPattern = `^[\w\.\-]+$`

	// regularFieldNameWithCommaRegPattern is the regular expression pattern for one or more strings
	// which are regular field names of table, multiple field names joined with char ','.
	regularFieldNameWithCommaRegPattern = `^[\w\.\-,\s]+$`

	// regularFieldNameWithoutDotRegPattern is similar to regularFieldNameRegPattern but not allows '.'.
	// Note that, although some databases allow char '.' in the field name, but it here does not allow '.'
	// in the field name as it conflicts with "db.table.field" pattern in SOME situations.
	regularFieldNameWithoutDotRegPattern = `^[\w\-]+$`

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
	driverMap[name] = newDriverWrapper(driver)
	return nil
}

// New creates and returns an ORM object with given configuration node.
func New(node ConfigNode) (db DB, err error) {
	return newDBByConfigNode(&node, "")
}

// NewByGroup creates and returns an ORM object with global configurations.
// The parameter `name` specifies the configuration group name,
// which is DefaultGroupName in default.
func NewByGroup(group ...string) (db DB, err error) {
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
		var node *ConfigNode
		if node, err = getConfigNodeByGroup(groupName, true); err == nil {
			return newDBByConfigNode(node, groupName)
		}
		return nil, err
	}
	return nil, gerror.NewCodef(
		gcode.CodeInvalidConfiguration,
		`database configuration node "%s" is not found, did you misspell group name "%s" or miss the database configuration?`,
		groupName, groupName,
	)
}

// newDBByConfigNode creates and returns an ORM object with given configuration node and group name.
//
// Very Note:
// The parameter `node` is used for DB creation, not for underlying connection creation.
// So all db type configurations in the same group should be the same.
func newDBByConfigNode(node *ConfigNode, group string) (db DB, err error) {
	if node.Link != "" {
		node = parseConfigNodeLink(node)
	}
	c := &Core{
		group:         group,
		debug:         gtype.NewBool(),
		cache:         gcache.New(),
		links:         gmap.New(true),
		logger:        glog.New(),
		config:        node,
		innerMemCache: gcache.New(),
		dynamicConfig: dynamicConfig{
			MaxIdleConnCount: node.MaxIdleConnCount,
			MaxOpenConnCount: node.MaxOpenConnCount,
			MaxConnLifeTime:  node.MaxConnLifeTime,
		},
	}
	if v, ok := driverMap[node.Type]; ok {
		if c.db, err = v.New(c, node); err != nil {
			return nil, err
		}
		return c.db, nil
	}
	errorMsg := `cannot find database driver for specified database type "%s"`
	errorMsg += `, did you misspell type name "%s" or forget importing the database driver? `
	errorMsg += `possible reference: https://github.com/gogf/gf/tree/master/contrib/drivers`
	return nil, gerror.NewCodef(gcode.CodeInvalidConfiguration, errorMsg, node.Type, node.Type)
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
		db, err = NewByGroup(group)
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
	}
	return nil, gerror.NewCodef(
		gcode.CodeInvalidConfiguration,
		"empty database configuration for item name '%s'",
		group,
	)
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
	var (
		min    = 0
		max    = 0
		random = grand.N(0, total-1)
	)
	for i := 0; i < len(cg); i++ {
		max = min + cg[i].Weight*100
		if random >= min && random < max {
			// ====================================================
			// Return a COPY of the ConfigNode.
			// ====================================================
			node := ConfigNode{}
			node = cg[i]
			return &node
		}
		min = max
	}
	return nil
}

// getSqlDb retrieves and returns an underlying database connection object.
// The parameter `master` specifies whether retrieves master node connection if
// master-slave nodes are configured.
func (c *Core) getSqlDb(master bool, schema ...string) (sqlDb *sql.DB, err error) {
	var (
		node *ConfigNode
		ctx  = c.db.GetCtx()
	)
	if c.group != "" {
		// Load balance.
		configs.RLock()
		defer configs.RUnlock()
		// Value COPY for node.
		node, err = getConfigNodeByGroup(c.group, master)
		if err != nil {
			return nil, err
		}
	} else {
		// Value COPY for node.
		n := *c.db.GetConfig()
		node = &n
	}
	if node.Charset == "" {
		node.Charset = defaultCharset
	}
	// Changes the schema.
	nodeSchema := gutil.GetOrDefaultStr(c.schema, schema...)
	if nodeSchema != "" {
		node.Name = nodeSchema
	}
	// Update the configuration object in internal data.
	if err = c.setConfigNodeToCtx(ctx, node); err != nil {
		return
	}

	// Cache the underlying connection pool object by node.
	var (
		instanceCacheFunc = func() interface{} {
			if sqlDb, err = c.db.Open(node); err != nil {
				return nil
			}
			if sqlDb == nil {
				return nil
			}
			if c.dynamicConfig.MaxIdleConnCount > 0 {
				sqlDb.SetMaxIdleConns(c.dynamicConfig.MaxIdleConnCount)
			} else {
				sqlDb.SetMaxIdleConns(defaultMaxIdleConnCount)
			}
			if c.dynamicConfig.MaxOpenConnCount > 0 {
				sqlDb.SetMaxOpenConns(c.dynamicConfig.MaxOpenConnCount)
			} else {
				sqlDb.SetMaxOpenConns(defaultMaxOpenConnCount)
			}
			if c.dynamicConfig.MaxConnLifeTime > 0 {
				sqlDb.SetConnMaxLifetime(c.dynamicConfig.MaxConnLifeTime)
			} else {
				sqlDb.SetConnMaxLifetime(defaultMaxConnLifeTime)
			}
			return sqlDb
		}
		// it here uses node value not pointer as the cache key, in case of oracle ORA-12516 error.
		instanceValue = c.links.GetOrSetFuncLock(*node, instanceCacheFunc)
	)
	if instanceValue != nil && sqlDb == nil {
		// It reads from instance map.
		sqlDb = instanceValue.(*sql.DB)
	}
	if node.Debug {
		c.db.SetDebug(node.Debug)
	}
	if node.DryRun {
		c.db.SetDryRun(node.DryRun)
	}
	return
}
