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
	Model(tableNameOrStruct ...any) *Model

	// Raw creates and returns a model based on a raw sql not a table.
	Raw(rawSql string, args ...any) *Model

	// Schema switches to a specified schema.
	// Also see Core.Schema.
	Schema(schema string) *Schema

	// With creates and returns an ORM model based on metadata of given object.
	// Also see Core.With.
	With(objects ...any) *Model

	// Open creates a raw connection object for database with given node configuration.
	// Note that it is not recommended using the function manually.
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

	// Query executes a SQL query that returns rows using given SQL and arguments.
	// The args are for any placeholder parameters in the query.
	Query(ctx context.Context, sql string, args ...any) (Result, error)

	// Exec executes a SQL query that doesn't return rows (e.g., INSERT, UPDATE, DELETE).
	// It returns sql.Result for accessing LastInsertId or RowsAffected.
	Exec(ctx context.Context, sql string, args ...any) (sql.Result, error)

	// Prepare creates a prepared statement for later queries or executions.
	// The execOnMaster parameter determines whether the statement executes on master node.
	Prepare(ctx context.Context, sql string, execOnMaster ...bool) (*Stmt, error)

	// ===========================================================================
	// Common APIs for CRUD.
	// ===========================================================================

	// Insert inserts one or multiple records into table.
	// The data can be a map, struct, or slice of maps/structs.
	// The optional batch parameter specifies the batch size for bulk inserts.
	Insert(ctx context.Context, table string, data any, batch ...int) (sql.Result, error)

	// InsertIgnore inserts records but ignores duplicate key errors.
	// It works like Insert but adds IGNORE keyword to the SQL statement.
	InsertIgnore(ctx context.Context, table string, data any, batch ...int) (sql.Result, error)

	// InsertAndGetId inserts a record and returns the auto-generated ID.
	// It's a convenience method combining Insert with LastInsertId.
	InsertAndGetId(ctx context.Context, table string, data any, batch ...int) (int64, error)

	// Replace inserts or replaces records using REPLACE INTO syntax.
	// Existing records with same unique key will be deleted and re-inserted.
	Replace(ctx context.Context, table string, data any, batch ...int) (sql.Result, error)

	// Save inserts or updates records using INSERT ... ON DUPLICATE KEY UPDATE syntax.
	// It updates existing records instead of replacing them entirely.
	Save(ctx context.Context, table string, data any, batch ...int) (sql.Result, error)

	// Update updates records in table that match the condition.
	// The data can be a map or struct containing the new values.
	// The condition specifies the WHERE clause with optional placeholder args.
	Update(ctx context.Context, table string, data any, condition any, args ...any) (sql.Result, error)

	// Delete deletes records from table that match the condition.
	// The condition specifies the WHERE clause with optional placeholder args.
	Delete(ctx context.Context, table string, condition any, args ...any) (sql.Result, error)

	// ===========================================================================
	// Internal APIs for CRUD, which can be overwritten by custom CRUD implements.
	// ===========================================================================

	// DoSelect executes a SELECT query using the given link and returns the result.
	// This is an internal method that can be overridden by custom implementations.
	DoSelect(ctx context.Context, link Link, sql string, args ...any) (result Result, err error)

	// DoInsert performs the actual INSERT operation with given options.
	// This is an internal method that can be overridden by custom implementations.
	DoInsert(ctx context.Context, link Link, table string, data List, option DoInsertOption) (result sql.Result, err error)

	// DoUpdate performs the actual UPDATE operation.
	// This is an internal method that can be overridden by custom implementations.
	DoUpdate(ctx context.Context, link Link, table string, data any, condition string, args ...any) (result sql.Result, err error)

	// DoDelete performs the actual DELETE operation.
	// This is an internal method that can be overridden by custom implementations.
	DoDelete(ctx context.Context, link Link, table string, condition string, args ...any) (result sql.Result, err error)

	// DoQuery executes a query that returns rows.
	// This is an internal method that can be overridden by custom implementations.
	DoQuery(ctx context.Context, link Link, sql string, args ...any) (result Result, err error)

	// DoExec executes a query that doesn't return rows.
	// This is an internal method that can be overridden by custom implementations.
	DoExec(ctx context.Context, link Link, sql string, args ...any) (result sql.Result, err error)

	// DoFilter processes and filters SQL and args before execution.
	// This is an internal method that can be overridden to implement custom SQL filtering.
	DoFilter(ctx context.Context, link Link, sql string, args []any) (newSql string, newArgs []any, err error)

	// DoCommit handles the actual commit operation for transactions.
	// This is an internal method that can be overridden by custom implementations.
	DoCommit(ctx context.Context, in DoCommitInput) (out DoCommitOutput, err error)

	// DoPrepare creates a prepared statement on the given link.
	// This is an internal method that can be overridden by custom implementations.
	DoPrepare(ctx context.Context, link Link, sql string) (*Stmt, error)

	// ===========================================================================
	// Query APIs for convenience purpose.
	// ===========================================================================

	// GetAll executes a query and returns all rows as Result.
	// It's a convenience wrapper around Query.
	GetAll(ctx context.Context, sql string, args ...any) (Result, error)

	// GetOne executes a query and returns the first row as Record.
	// It's useful when you expect only one row to be returned.
	GetOne(ctx context.Context, sql string, args ...any) (Record, error)

	// GetValue executes a query and returns the first column of the first row.
	// It's useful for queries like SELECT COUNT(*) or getting a single value.
	GetValue(ctx context.Context, sql string, args ...any) (Value, error)

	// GetArray executes a query and returns the first column of all rows.
	// It's useful for queries like SELECT id FROM table.
	GetArray(ctx context.Context, sql string, args ...any) ([]Value, error)

	// GetCount executes a COUNT query and returns the result as an integer.
	// It's a convenience method for counting rows.
	GetCount(ctx context.Context, sql string, args ...any) (int, error)

	// GetScan executes a query and scans the result into the given object pointer.
	// It automatically maps database columns to struct fields or slice elements.
	GetScan(ctx context.Context, objPointer any, sql string, args ...any) error

	// Union combines multiple SELECT queries using UNION operator.
	// It returns a new Model that represents the combined query.
	Union(unions ...*Model) *Model

	// UnionAll combines multiple SELECT queries using UNION ALL operator.
	// Unlike Union, it keeps duplicate rows in the result.
	UnionAll(unions ...*Model) *Model

	// ===========================================================================
	// Master/Slave specification support.
	// ===========================================================================

	// Master returns a connection to the master database node.
	// The optional schema parameter specifies which database schema to use.
	Master(schema ...string) (*sql.DB, error)

	// Slave returns a connection to a slave database node.
	// The optional schema parameter specifies which database schema to use.
	Slave(schema ...string) (*sql.DB, error)

	// ===========================================================================
	// Ping-Pong.
	// ===========================================================================

	// PingMaster checks if the master database node is accessible.
	// It returns an error if the connection fails.
	PingMaster() error

	// PingSlave checks if any slave database node is accessible.
	// It returns an error if no slave connections are available.
	PingSlave() error

	// ===========================================================================
	// Transaction.
	// ===========================================================================

	// Begin starts a new transaction and returns a TX interface.
	// The returned TX must be committed or rolled back to release resources.
	Begin(ctx context.Context) (TX, error)

	// BeginWithOptions starts a new transaction with the given options and returns a TX interface.
	// The options allow specifying isolation level and read-only mode.
	// The returned TX must be committed or rolled back to release resources.
	BeginWithOptions(ctx context.Context, opts TxOptions) (TX, error)

	// Transaction executes a function within a transaction.
	// It automatically handles commit/rollback based on whether f returns an error.
	Transaction(ctx context.Context, f func(ctx context.Context, tx TX) error) error

	// TransactionWithOptions executes a function within a transaction with specific options.
	// It allows customizing transaction behavior like isolation level and timeout.
	TransactionWithOptions(ctx context.Context, opts TxOptions, f func(ctx context.Context, tx TX) error) error

	// ===========================================================================
	// Configuration methods.
	// ===========================================================================

	// GetCache returns the cache instance used by this database.
	// The cache is used for query results caching.
	GetCache() *gcache.Cache

	// SetDebug enables or disables debug mode for SQL logging.
	// When enabled, all SQL statements and their execution time are logged.
	SetDebug(debug bool)

	// GetDebug returns whether debug mode is enabled.
	GetDebug() bool

	// GetSchema returns the current database schema name.
	GetSchema() string

	// GetPrefix returns the table name prefix used by this database.
	GetPrefix() string

	// GetGroup returns the configuration group name of this database.
	GetGroup() string

	// SetDryRun enables or disables dry-run mode.
	// In dry-run mode, SQL statements are generated but not executed.
	SetDryRun(enabled bool)

	// GetDryRun returns whether dry-run mode is enabled.
	GetDryRun() bool

	// SetLogger sets a custom logger for database operations.
	// The logger must implement glog.ILogger interface.
	SetLogger(logger glog.ILogger)

	// GetLogger returns the current logger used by this database.
	GetLogger() glog.ILogger

	// GetConfig returns the configuration node used by this database.
	GetConfig() *ConfigNode

	// SetMaxIdleConnCount sets the maximum number of idle connections in the pool.
	SetMaxIdleConnCount(n int)

	// SetMaxOpenConnCount sets the maximum number of open connections to the database.
	SetMaxOpenConnCount(n int)

	// SetMaxConnLifeTime sets the maximum amount of time a connection may be reused.
	SetMaxConnLifeTime(d time.Duration)

	// ===========================================================================
	// Utility methods.
	// ===========================================================================

	// Stats returns statistics about the database connection pool.
	// It includes information like the number of active and idle connections.
	Stats(ctx context.Context) []StatsItem

	// GetCtx returns the context associated with this database instance.
	GetCtx() context.Context

	// GetCore returns the underlying Core instance of this database.
	GetCore() *Core

	// GetChars returns the left and right quote characters used for escaping identifiers.
	// For example, in MySQL these are backticks: ` and `.
	GetChars() (charLeft string, charRight string)

	// Tables returns a list of all table names in the specified schema.
	// If no schema is specified, it uses the default schema.
	Tables(ctx context.Context, schema ...string) (tables []string, err error)

	// TableFields returns detailed information about all fields in the specified table.
	// The returned map keys are field names and values contain field metadata.
	TableFields(ctx context.Context, table string, schema ...string) (map[string]*TableField, error)

	// ConvertValueForField converts a value to the appropriate type for a database field.
	// It handles type conversion from Go types to database-specific types.
	ConvertValueForField(ctx context.Context, fieldType string, fieldValue any) (any, error)

	// ConvertValueForLocal converts a database value to the appropriate Go type.
	// It handles type conversion from database-specific types to Go types.
	ConvertValueForLocal(ctx context.Context, fieldType string, fieldValue any) (any, error)

	// GetFormattedDBTypeNameForField returns the formatted database type name and pattern for a field type.
	GetFormattedDBTypeNameForField(fieldType string) (typeName, typePattern string)

	// CheckLocalTypeForField checks if a Go value is compatible with a database field type.
	// It returns the appropriate LocalType and any conversion errors.
	CheckLocalTypeForField(ctx context.Context, fieldType string, fieldValue any) (LocalType, error)

	// FormatUpsert formats an upsert (INSERT ... ON DUPLICATE KEY UPDATE) statement.
	// It generates the appropriate SQL based on the columns, values, and options provided.
	FormatUpsert(columns []string, list List, option DoInsertOption) (string, error)

	// OrderRandomFunction returns the SQL function for random ordering.
	// The implementation is database-specific (e.g., RAND() for MySQL).
	OrderRandomFunction() string
}

// TX defines the interfaces for ORM transaction operations.
type TX interface {
	Link

	// Ctx binds a context to current transaction.
	// The context is used for operations like timeout control.
	Ctx(ctx context.Context) TX

	// Raw creates and returns a model based on a raw SQL.
	// The rawSql can contain placeholders ? and corresponding args.
	Raw(rawSql string, args ...any) *Model

	// Model creates and returns a Model from given table name/struct.
	// The parameter can be table name as string, or struct/*struct type.
	Model(tableNameQueryOrStruct ...any) *Model

	// With creates and returns a Model from given object.
	// It automatically analyzes the object and generates corresponding SQL.
	With(object any) *Model

	// ===========================================================================
	// Nested transaction if necessary.
	// ===========================================================================

	// Begin starts a nested transaction.
	// It creates a new savepoint for current transaction.
	Begin() error

	// Commit commits current transaction/savepoint.
	// For nested transactions, it releases the current savepoint.
	Commit() error

	// Rollback rolls back current transaction/savepoint.
	// For nested transactions, it rolls back to the current savepoint.
	Rollback() error

	// Transaction executes given function in a nested transaction.
	// It automatically handles commit/rollback based on function's error return.
	Transaction(ctx context.Context, f func(ctx context.Context, tx TX) error) (err error)

	// TransactionWithOptions executes given function in a nested transaction with options.
	// It allows customizing transaction behavior like isolation level.
	TransactionWithOptions(ctx context.Context, opts TxOptions, f func(ctx context.Context, tx TX) error) error

	// ===========================================================================
	// Core method.
	// ===========================================================================

	// Query executes a query that returns rows using given SQL and arguments.
	// The args are for any placeholder parameters in the query.
	Query(sql string, args ...any) (result Result, err error)

	// Exec executes a query that doesn't return rows.
	// For example: INSERT, UPDATE, DELETE.
	Exec(sql string, args ...any) (sql.Result, error)

	// Prepare creates a prepared statement for later queries or executions.
	// Multiple queries or executions may be run concurrently from the statement.
	Prepare(sql string) (*Stmt, error)

	// ===========================================================================
	// Query.
	// ===========================================================================

	// GetAll executes a query and returns all rows as Result.
	// It's a convenient wrapper for Query.
	GetAll(sql string, args ...any) (Result, error)

	// GetOne executes a query and returns the first row as Record.
	// It's useful when you expect only one row to be returned.
	GetOne(sql string, args ...any) (Record, error)

	// GetStruct executes a query and scans the result into given struct.
	// The obj should be a pointer to struct.
	GetStruct(obj any, sql string, args ...any) error

	// GetStructs executes a query and scans all results into given struct slice.
	// The objPointerSlice should be a pointer to slice of struct.
	GetStructs(objPointerSlice any, sql string, args ...any) error

	// GetScan executes a query and scans the result into given variables.
	// The pointer can be type of struct/*struct/[]struct/[]*struct.
	GetScan(pointer any, sql string, args ...any) error

	// GetValue executes a query and returns the first column of first row.
	// It's useful for queries like SELECT COUNT(*).
	GetValue(sql string, args ...any) (Value, error)

	// GetCount executes a query that should return a count value.
	// It's a convenient wrapper for count queries.
	GetCount(sql string, args ...any) (int64, error)

	// ===========================================================================
	// CRUD.
	// ===========================================================================

	// Insert inserts one or multiple records into table.
	// The data can be map/struct/*struct/[]map/[]struct/[]*struct.
	Insert(table string, data any, batch ...int) (sql.Result, error)

	// InsertIgnore inserts one or multiple records with IGNORE option.
	// It ignores records that would cause duplicate key conflicts.
	InsertIgnore(table string, data any, batch ...int) (sql.Result, error)

	// InsertAndGetId inserts one record and returns its id value.
	// It's commonly used with auto-increment primary key.
	InsertAndGetId(table string, data any, batch ...int) (int64, error)

	// Replace inserts or replaces records using REPLACE INTO syntax.
	// Existing records with same unique key will be deleted and re-inserted.
	Replace(table string, data any, batch ...int) (sql.Result, error)

	// Save inserts or updates records using INSERT ... ON DUPLICATE KEY UPDATE syntax.
	// It updates existing records instead of replacing them entirely.
	Save(table string, data any, batch ...int) (sql.Result, error)

	// Update updates records in table that match given condition.
	// The data can be map/struct, and condition supports various formats.
	Update(table string, data any, condition any, args ...any) (sql.Result, error)

	// Delete deletes records from table that match given condition.
	// The condition supports various formats with optional arguments.
	Delete(table string, condition any, args ...any) (sql.Result, error)

	// ===========================================================================
	// Utility methods.
	// ===========================================================================

	// GetCtx returns the context that is bound to current transaction.
	GetCtx() context.Context

	// GetDB returns the underlying DB interface object.
	GetDB() DB

	// GetSqlTX returns the underlying *sql.Tx object.
	// Note: be very careful when using this method.
	GetSqlTX() *sql.Tx

	// IsClosed checks if current transaction is closed.
	// A transaction is closed after Commit or Rollback.
	IsClosed() bool

	// ===========================================================================
	// Save point feature.
	// ===========================================================================

	// SavePoint creates a save point with given name.
	// It's used in nested transactions to create rollback points.
	SavePoint(point string) error

	// RollbackTo rolls back transaction to previously created save point.
	// If the save point doesn't exist, it returns an error.
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
	localTypeMap  *gmap.StrAnyMap // Local type map for database field type conversion.
	dynamicConfig dynamicConfig   // Dynamic configurations, which can be changed in runtime.
	innerMemCache *gcache.Cache   // Internal memory cache for storing temporary data.
}

type dynamicConfig struct {
	MaxIdleConnCount int
	MaxOpenConnCount int
	MaxConnLifeTime  time.Duration
}

// DoCommitInput is the input parameters for function DoCommit.
type DoCommitInput struct {
	// Db is the underlying database connection object.
	Db *sql.DB

	// Tx is the underlying transaction object.
	Tx *sql.Tx

	// Stmt is the prepared statement object.
	Stmt *sql.Stmt

	// Link is the common database function wrapper interface.
	Link Link

	// Sql is the SQL string to be executed.
	Sql string

	// Args is the arguments for SQL placeholders.
	Args []any

	// Type indicates the type of SQL operation.
	Type SqlType

	// TxOptions specifies the transaction options.
	TxOptions sql.TxOptions

	// TxCancelFunc is the context cancel function for transaction.
	TxCancelFunc context.CancelFunc

	// IsTransaction indicates whether current operation is in transaction.
	IsTransaction bool
}

// DoCommitOutput is the output parameters for function DoCommit.
type DoCommitOutput struct {
	// Result is the result of exec statement.
	Result sql.Result

	// Records is the result of query statement.
	Records []Record

	// Stmt is the Statement object result for Prepare.
	Stmt *Stmt

	// Tx is the transaction object result for Begin.
	Tx TX

	// RawResult is the underlying result, which might be sql.Result/*sql.Rows/*sql.Row.
	RawResult any
}

// Driver is the interface for integrating sql drivers into package gdb.
type Driver interface {
	// New creates and returns a database object for specified database server.
	New(core *Core, node *ConfigNode) (DB, error)
}

// Link is a common database function wrapper interface.
// Note that, any operation using `Link` will have no SQL logging.
type Link interface {
	QueryContext(ctx context.Context, sql string, args ...any) (*sql.Rows, error)
	ExecContext(ctx context.Context, sql string, args ...any) (sql.Result, error)
	PrepareContext(ctx context.Context, sql string) (*sql.Stmt, error)
	IsOnMaster() bool
	IsTransaction() bool
}

// Sql is the sql recording struct.
type Sql struct {
	Sql           string  // SQL string(may contain reserved char '?').
	Type          SqlType // SQL operation type.
	Args          []any   // Arguments for this sql.
	Format        string  // Formatted sql which contains arguments in the sql.
	Error         error   // Execution result.
	Start         int64   // Start execution timestamp in milliseconds.
	End           int64   // End execution timestamp in milliseconds.
	Group         string  // Group is the group name of the configuration that the sql is executed from.
	Schema        string  // Schema is the schema name of the configuration that the sql is executed from.
	IsTransaction bool    // IsTransaction marks whether this sql is executed in transaction.
	RowsAffected  int64   // RowsAffected marks retrieved or affected number with current sql statement.
}

// DoInsertOption is the input struct for function DoInsert.
type DoInsertOption struct {
	// OnDuplicateStr is the custom string for `on duplicated` statement.
	OnDuplicateStr string

	// OnDuplicateMap is the custom key-value map from `OnDuplicateEx` function for `on duplicated` statement.
	OnDuplicateMap map[string]any

	// OnConflict is the custom conflict key of upsert clause, if the database needs it.
	OnConflict []string

	// InsertOption is the insert operation in constant value.
	InsertOption InsertOption

	// BatchCount is the batch count for batch inserting.
	BatchCount int
}

// TableField is the struct for table field.
type TableField struct {
	// Index is for ordering purpose as map is unordered.
	Index int

	// Name is the field name.
	Name string

	// Type is the field type. Eg: 'int(10) unsigned', 'varchar(64)'.
	Type string

	// Null is whether the field can be null or not.
	Null bool

	// Key is the index information(empty if it's not an index). Eg: PRI, MUL.
	Key string

	// Default is the default value for the field.
	Default any

	// Extra is the extra information. Eg: auto_increment.
	Extra string

	// Comment is the field comment.
	Comment string
}

// Counter is the type for update count.
type Counter struct {
	// Field is the field name.
	Field string

	// Value is the value.
	Value float64
}

type (
	// Raw is a raw sql that will not be treated as argument but as a direct sql part.
	Raw string

	// Value is the field value type.
	Value = *gvar.Var

	// Record is the row record of the table.
	Record map[string]Value

	// Result is the row record array.
	Result []Record

	// Map is alias of map[string]any, which is the most common usage map type.
	Map = map[string]any

	// List is type of map array.
	List = []Map
)

type CatchSQLManager struct {
	// SQLArray is the array of sql.
	SQLArray *garray.StrArray

	// DoCommit marks it will be committed to underlying driver or not.
	DoCommit bool
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
	cachePrefixTableFields                = `TableFields:`
	cachePrefixSelectCache                = `SelectCache:`
	commandEnvKeyForDryRun                = "gf.gdb.dryrun"
	modelForDaoSuffix                     = `ForDao`
	dbRoleSlave                           = `slave`
	ctxKeyForDB               gctx.StrKey = `CtxKeyForDB`
	ctxKeyCatchSQL            gctx.StrKey = `CtxKeyCatchSQL`
	ctxKeyInternalProducedSQL gctx.StrKey = `CtxKeyInternalProducedSQL`

	linkPattern            = `^(\w+):(.*?):(.*?)@(\w+?)\((.+?)\)/{0,1}([^\?]*)\?{0,1}(.*?)$`
	linkPatternDescription = `type:username:password@protocol(host:port)/dbname?param1=value1&...&paramN=valueN`
)

// Context key types to avoid collisions
type ctxKey string

const (
	ctxKeyWrappedByGetCtxTimeout ctxKey = "WrappedByGetCtxTimeout"
)

type ctxTimeoutType int

const (
	ctxTimeoutTypeExec ctxTimeoutType = iota
	ctxTimeoutTypeQuery
	ctxTimeoutTypePrepare
	ctxTimeoutTypeTrans
)

type SelectType int

const (
	SelectTypeDefault SelectType = iota
	SelectTypeCount
	SelectTypeValue
	SelectTypeArray
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

// LocalType is a type that defines the local storage type of a field value.
// It is used to specify how the field value should be processed locally.
type LocalType string

const (
	LocalTypeUndefined   LocalType = ""
	LocalTypeString      LocalType = "string"
	LocalTypeTime        LocalType = "time"
	LocalTypeDate        LocalType = "date"
	LocalTypeDatetime    LocalType = "datetime"
	LocalTypeInt         LocalType = "int"
	LocalTypeUint        LocalType = "uint"
	LocalTypeInt64       LocalType = "int64"
	LocalTypeUint64      LocalType = "uint64"
	LocalTypeBigInt      LocalType = "bigint"
	LocalTypeIntSlice    LocalType = "[]int"
	LocalTypeInt64Slice  LocalType = "[]int64"
	LocalTypeUint64Slice LocalType = "[]uint64"
	LocalTypeStringSlice LocalType = "[]string"
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
	fieldTypeInt128     = "int128"
	fieldTypeInt256     = "int256"
	fieldTypeUint128    = "uint128"
	fieldTypeUint256    = "uint256"
	fieldTypeReal       = "real"
	fieldTypeFloat      = "float"
	fieldTypeDouble     = "double"
	fieldTypeDecimal    = "decimal"
	fieldTypeMoney      = "money"
	fieldTypeNumeric    = "numeric"
	fieldTypeSmallmoney = "smallmoney"
	fieldTypeBool       = "bool"
	fieldTypeBit        = "bit"
	fieldTypeYear       = "year"      // YYYY
	fieldTypeDate       = "date"      // YYYY-MM-DD
	fieldTypeTime       = "time"      // HH:MM:SS
	fieldTypeDatetime   = "datetime"  // YYYY-MM-DD HH:MM:SS
	fieldTypeTimestamp  = "timestamp" // YYYYMMDD HHMMSS
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
		node, err = parseConfigNodeLink(node)
		if err != nil {
			return
		}
	}
	c := &Core{
		group:         group,
		debug:         gtype.NewBool(),
		cache:         gcache.New(),
		links:         gmap.New(true),
		logger:        glog.New(),
		config:        node,
		localTypeMap:  gmap.NewStrAnyMap(true),
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
	v := instances.GetOrSetFuncLock(group, func() any {
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
// The returned node is a clone of configuration node, which is safe for later modification.
//
// The parameter `master` specifies whether retrieving a master node, or else a slave node
// if master-slave nodes are configured.
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
// The returned node is a clone of configuration node, which is safe for later modification.
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
		minWeight = 0
		maxWeight = 0
		random    = grand.N(0, total-1)
	)
	for i := 0; i < len(cg); i++ {
		maxWeight = minWeight + cg[i].Weight*100
		if random >= minWeight && random < maxWeight {
			// ====================================================
			// Return a COPY of the ConfigNode.
			// ====================================================
			node := ConfigNode{}
			node = cg[i]
			return &node
		}
		minWeight = maxWeight
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
		// The returned node is a clone of configuration node, which is safe for later modification.
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
		instanceCacheFunc = func() any {
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
		// it here uses NODE VALUE not pointer as the cache key, in case of oracle ORA-12516 error.
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
