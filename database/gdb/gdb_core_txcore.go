// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"context"
	"database/sql"
	"reflect"

	"github.com/gogf/gf/v3/errors/gcode"
	"github.com/gogf/gf/v3/errors/gerror"
	"github.com/gogf/gf/v3/internal/reflection"
	"github.com/gogf/gf/v3/text/gregex"
	"github.com/gogf/gf/v3/util/gconv"
)

// TXCore is the struct for transaction management.
type TXCore struct {
	// db is the database management interface that implements the DB interface,
	// providing access to database operations and configuration.
	db DB
	// tx is the underlying SQL transaction object from database/sql package,
	// which manages the actual transaction operations.
	tx *sql.Tx
	// master is the underlying master database connection pool,
	// used for direct database operations when needed.
	master *sql.DB
	// transactionId is a unique identifier for this transaction instance,
	// used for tracking and debugging purposes.
	transactionId string
	// transactionCount tracks the number of nested transaction begins,
	// used for managing transaction nesting depth.
	transactionCount int
	// isClosed indicates whether this transaction has been finalized
	// through either a commit or rollback operation.
	isClosed bool
	// cancelFunc is the context cancellation function associated with ctx,
	// used to cancel the transaction context when needed.
	cancelFunc context.CancelFunc
}

func (c *Core) newEmptyTX() TX {
	return &TXCore{
		db: c.db,
	}
}

// transactionKeyForNestedPoint forms and returns the transaction key at current save point.
func (tx *TXCore) transactionKeyForNestedPoint() string {
	return tx.db.GetCore().QuoteWord(
		transactionPointerPrefix + gconv.String(tx.transactionCount),
	)
}

// GetDB returns the DB for current transaction.
func (tx *TXCore) GetDB() DB {
	return tx.db
}

// GetSqlTX returns the underlying transaction object for current transaction.
func (tx *TXCore) GetSqlTX() *sql.Tx {
	return tx.tx
}

// Commit commits current transaction.
// Note that it releases previous saved transaction point if it's in a nested transaction procedure,
// or else it commits the hole transaction.
func (tx *TXCore) Commit(ctx context.Context) error {
	if tx.transactionCount > 0 {
		tx.transactionCount--
		_, err := tx.Exec(ctx, "RELEASE SAVEPOINT "+tx.transactionKeyForNestedPoint())
		return err
	}
	_, err := tx.db.DoCommit(ctx, DoCommitInput{
		Tx:            tx.tx,
		Sql:           "COMMIT",
		Type:          SqlTypeTXCommit,
		TxCancelFunc:  tx.cancelFunc,
		IsTransaction: true,
	})
	if err == nil {
		tx.isClosed = true
	}
	return err
}

// Rollback aborts current transaction.
// Note that it aborts current transaction if it's in a nested transaction procedure,
// or else it aborts the hole transaction.
func (tx *TXCore) Rollback(ctx context.Context) error {
	if tx.transactionCount > 0 {
		tx.transactionCount--
		_, err := tx.Exec(ctx, "ROLLBACK TO SAVEPOINT "+tx.transactionKeyForNestedPoint())
		return err
	}
	_, err := tx.db.DoCommit(ctx, DoCommitInput{
		Tx:            tx.tx,
		Sql:           "ROLLBACK",
		Type:          SqlTypeTXRollback,
		TxCancelFunc:  tx.cancelFunc,
		IsTransaction: true,
	})
	if err == nil {
		tx.isClosed = true
	}
	return err
}

// IsClosed checks and returns this transaction has already been committed or rolled back.
func (tx *TXCore) IsClosed() bool {
	return tx.isClosed
}

// Begin starts a nested transaction procedure.
func (tx *TXCore) Begin(ctx context.Context) error {
	_, err := tx.Exec(ctx, "SAVEPOINT "+tx.transactionKeyForNestedPoint())
	if err != nil {
		return err
	}
	tx.transactionCount++
	return nil
}

// SavePoint performs `SAVEPOINT xxx` SQL statement that saves transaction at current point.
// The parameter `point` specifies the point name that will be saved to server.
func (tx *TXCore) SavePoint(ctx context.Context, point string) error {
	_, err := tx.Exec(ctx, "SAVEPOINT "+tx.db.GetCore().QuoteWord(point))
	return err
}

// RollbackTo performs `ROLLBACK TO SAVEPOINT xxx` SQL statement that rollbacks to specified saved transaction.
// The parameter `point` specifies the point name that was saved previously.
func (tx *TXCore) RollbackTo(ctx context.Context, point string) error {
	_, err := tx.Exec(ctx, "ROLLBACK TO SAVEPOINT "+tx.db.GetCore().QuoteWord(point))
	return err
}

// Transaction wraps the transaction logic using function `f`.
// It rollbacks the transaction and returns the error from function `f` if
// it returns non-nil error. It commits the transaction and returns nil if
// function `f` returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function `f`
// as it is automatically handled by this function.
func (tx *TXCore) Transaction(ctx context.Context, f func(ctx context.Context, tx TX) error) (err error) {
	// Check transaction object from context.
	if TXFromCtx(ctx, tx.db.GetGroup()) == nil {
		// Inject transaction object into context.
		ctx = WithTX(ctx, tx)
	}
	if err = tx.Begin(ctx); err != nil {
		return err
	}
	err = callTxFunc(ctx, tx, f)
	return
}

// TransactionWithOptions wraps the transaction logic with propagation options using function `f`.
func (tx *TXCore) TransactionWithOptions(
	ctx context.Context, opts TxOptions, f func(ctx context.Context, tx TX) error,
) (err error) {
	return tx.db.TransactionWithOptions(ctx, opts, f)
}

// Query does query operation on transaction.
// See Core.Query.
func (tx *TXCore) Query(ctx context.Context, sql string, args ...interface{}) (result Result, err error) {
	return tx.db.DoQuery(ctx, &txLink{tx.tx}, sql, args...)
}

// Exec does none query operation on transaction.
// See Core.Exec.
func (tx *TXCore) Exec(ctx context.Context, sql string, args ...interface{}) (sql.Result, error) {
	return tx.db.DoExec(ctx, &txLink{tx.tx}, sql, args...)
}

// Prepare creates a prepared statement for later queries or executions.
// Multiple queries or executions may be run concurrently from the
// returned statement.
// The caller must call the statement's Close method
// when the statement is no longer needed.
func (tx *TXCore) Prepare(ctx context.Context, sql string) (*Stmt, error) {
	return tx.db.DoPrepare(ctx, &txLink{tx.tx}, sql)
}

// GetAll queries and returns data records from database.
func (tx *TXCore) GetAll(ctx context.Context, sql string, args ...interface{}) (Result, error) {
	return tx.Query(ctx, sql, args...)
}

// GetOne queries and returns one record from database.
func (tx *TXCore) GetOne(ctx context.Context, sql string, args ...interface{}) (Record, error) {
	list, err := tx.GetAll(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	if len(list) > 0 {
		return list[0], nil
	}
	return nil, nil
}

// GetStruct queries one record from database and converts it to given struct.
// The parameter `pointer` should be a pointer to struct.
func (tx *TXCore) GetStruct(ctx context.Context, obj interface{}, sql string, args ...interface{}) error {
	one, err := tx.GetOne(ctx, sql, args...)
	if err != nil {
		return err
	}
	return one.Struct(obj)
}

// GetStructs queries records from database and converts them to given struct.
// The parameter `pointer` should be type of struct slice: []struct/[]*struct.
func (tx *TXCore) GetStructs(ctx context.Context, objPointerSlice interface{}, sql string, args ...interface{}) error {
	all, err := tx.GetAll(ctx, sql, args...)
	if err != nil {
		return err
	}
	return all.Structs(objPointerSlice)
}

// GetScan queries one or more records from database and converts them to given struct or
// struct array.
//
// If parameter `pointer` is type of struct pointer, it calls GetStruct internally for
// the conversion. If parameter `pointer` is type of slice, it calls GetStructs internally
// for conversion.
func (tx *TXCore) GetScan(ctx context.Context, pointer interface{}, sql string, args ...interface{}) error {
	reflectInfo := reflection.OriginTypeAndKind(pointer)
	if reflectInfo.InputKind != reflect.Ptr {
		return gerror.NewCodef(
			gcode.CodeInvalidParameter,
			"params should be type of pointer, but got: %v",
			reflectInfo.InputKind,
		)
	}
	switch reflectInfo.OriginKind {
	case reflect.Array, reflect.Slice:
		return tx.GetStructs(ctx, pointer, sql, args...)

	case reflect.Struct:
		return tx.GetStruct(ctx, pointer, sql, args...)

	default:
	}
	return gerror.NewCodef(
		gcode.CodeInvalidParameter,
		`in valid parameter type "%v", of which element type should be type of struct/slice`,
		reflectInfo.InputType,
	)
}

// GetValue queries and returns the field value from database.
// The sql should query only one field from database, or else it returns only one
// field of the result.
func (tx *TXCore) GetValue(ctx context.Context, sql string, args ...interface{}) (Value, error) {
	one, err := tx.GetOne(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	for _, v := range one {
		return v, nil
	}
	return nil, nil
}

// GetCount queries and returns the count from database.
func (tx *TXCore) GetCount(ctx context.Context, sql string, args ...interface{}) (int64, error) {
	if !gregex.IsMatchString(`(?i)SELECT\s+COUNT\(.+\)\s+FROM`, sql) {
		sql, _ = gregex.ReplaceString(`(?i)(SELECT)\s+(.+)\s+(FROM)`, `$1 COUNT($2) $3`, sql)
	}
	value, err := tx.GetValue(ctx, sql, args...)
	if err != nil {
		return 0, err
	}
	return value.Int64(), nil
}

// Insert does "INSERT INTO ..." statement for the table.
// If there's already one unique record of the data in the table, it returns error.
//
// The parameter `data` can be type of map/gmap/struct/*struct/[]map/[]struct, etc.
// Example:
// Data(g.Map{"uid": 10000, "name":"john"})
// Data(g.Slice{g.Map{"uid": 10000, "name":"john"}, g.Map{"uid": 20000, "name":"smith"})
//
// The parameter `batch` specifies the batch operation count when given data is slice.
func (tx *TXCore) Insert(ctx context.Context, table string, data interface{}, batch ...int) (sql.Result, error) {
	if len(batch) > 0 {
		return tx.Model(table).Data(data).Batch(batch[0]).Insert(ctx)
	}
	return tx.Model(table).Data(data).Insert(ctx)
}

// InsertIgnore does "INSERT IGNORE INTO ..." statement for the table.
// If there's already one unique record of the data in the table, it ignores the inserting.
//
// The parameter `data` can be type of map/gmap/struct/*struct/[]map/[]struct, etc.
// Example:
// Data(g.Map{"uid": 10000, "name":"john"})
// Data(g.Slice{g.Map{"uid": 10000, "name":"john"}, g.Map{"uid": 20000, "name":"smith"})
//
// The parameter `batch` specifies the batch operation count when given data is slice.
func (tx *TXCore) InsertIgnore(ctx context.Context, table string, data interface{}, batch ...int) (sql.Result, error) {
	if len(batch) > 0 {
		return tx.Model(table).Data(data).Batch(batch[0]).InsertIgnore(ctx)
	}
	return tx.Model(table).Data(data).InsertIgnore(ctx)
}

// InsertAndGetId performs action Insert and returns the last insert id that automatically generated.
func (tx *TXCore) InsertAndGetId(ctx context.Context, table string, data interface{}, batch ...int) (int64, error) {
	if len(batch) > 0 {
		return tx.Model(table).Data(data).Batch(batch[0]).InsertAndGetId(ctx)
	}
	return tx.Model(table).Data(data).InsertAndGetId(ctx)
}

// Replace does "REPLACE INTO ..." statement for the table.
// If there's already one unique record of the data in the table, it deletes the record
// and inserts a new one.
//
// The parameter `data` can be type of map/gmap/struct/*struct/[]map/[]struct, etc.
// Example:
// Data(g.Map{"uid": 10000, "name":"john"})
// Data(g.Slice{g.Map{"uid": 10000, "name":"john"}, g.Map{"uid": 20000, "name":"smith"})
//
// The parameter `data` can be type of map/gmap/struct/*struct/[]map/[]struct, etc.
// If given data is type of slice, it then does batch replacing, and the optional parameter
// `batch` specifies the batch operation count.
func (tx *TXCore) Replace(ctx context.Context, table string, data interface{}, batch ...int) (sql.Result, error) {
	if len(batch) > 0 {
		return tx.Model(table).Data(data).Batch(batch[0]).Replace(ctx)
	}
	return tx.Model(table).Data(data).Replace(ctx)
}

// Save does "INSERT INTO ... ON DUPLICATE KEY UPDATE..." statement for the table.
// It updates the record if there's primary or unique index in the saving data,
// or else it inserts a new record into the table.
//
// The parameter `data` can be type of map/gmap/struct/*struct/[]map/[]struct, etc.
// Example:
// Data(g.Map{"uid": 10000, "name":"john"})
// Data(g.Slice{g.Map{"uid": 10000, "name":"john"}, g.Map{"uid": 20000, "name":"smith"})
//
// If given data is type of slice, it then does batch saving, and the optional parameter
// `batch` specifies the batch operation count.
func (tx *TXCore) Save(ctx context.Context, table string, data interface{}, batch ...int) (sql.Result, error) {
	if len(batch) > 0 {
		return tx.Model(table).Data(data).Batch(batch[0]).Save(ctx)
	}
	return tx.Model(table).Data(data).Save(ctx)
}

// Update does "UPDATE ... " statement for the table.
//
// The parameter `data` can be type of string/map/gmap/struct/*struct, etc.
// Example: "uid=10000", "uid", 10000, g.Map{"uid": 10000, "name":"john"}
//
// The parameter `condition` can be type of string/map/gmap/slice/struct/*struct, etc.
// It is commonly used with parameter `args`.
// Example:
// "uid=10000",
// "uid", 10000
// "money>? AND name like ?", 99999, "vip_%"
// "status IN (?)", g.Slice{1,2,3}
// "age IN(?,?)", 18, 50
// User{ Id : 1, UserName : "john"}.
func (tx *TXCore) Update(ctx context.Context, table string, data interface{}, condition interface{}, args ...interface{}) (sql.Result, error) {
	return tx.Model(table).Data(data).Where(condition, args...).Update(ctx)
}

// Delete does "DELETE FROM ... " statement for the table.
//
// The parameter `condition` can be type of string/map/gmap/slice/struct/*struct, etc.
// It is commonly used with parameter `args`.
// Example:
// "uid=10000",
// "uid", 10000
// "money>? AND name like ?", 99999, "vip_%"
// "status IN (?)", g.Slice{1,2,3}
// "age IN(?,?)", 18, 50
// User{ Id : 1, UserName : "john"}.
func (tx *TXCore) Delete(ctx context.Context, table string, condition interface{}, args ...interface{}) (sql.Result, error) {
	return tx.Model(table).Where(condition, args...).Delete(ctx)
}

// QueryContext implements interface function Link.QueryContext.
func (tx *TXCore) QueryContext(ctx context.Context, sql string, args ...interface{}) (*sql.Rows, error) {
	return tx.tx.QueryContext(ctx, sql, args...)
}

// ExecContext implements interface function Link.ExecContext.
func (tx *TXCore) ExecContext(ctx context.Context, sql string, args ...interface{}) (sql.Result, error) {
	return tx.tx.ExecContext(ctx, sql, args...)
}

// PrepareContext implements interface function Link.PrepareContext.
func (tx *TXCore) PrepareContext(ctx context.Context, sql string) (*sql.Stmt, error) {
	return tx.tx.PrepareContext(ctx, sql)
}

// IsOnMaster implements interface function Link.IsOnMaster.
func (tx *TXCore) IsOnMaster() bool {
	return true
}

// IsTransaction implements interface function Link.IsTransaction.
func (tx *TXCore) IsTransaction() bool {
	return tx != nil
}
