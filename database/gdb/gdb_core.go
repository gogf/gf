// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//

package gdb

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/text/gstr"

	"github.com/gogf/gf/internal/utils"

	"github.com/gogf/gf/container/gvar"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/text/gregex"
	"github.com/gogf/gf/util/gconv"
)

// Ctx is a chaining function, which creates and returns a new DB that is a shallow copy
// of current DB object and with given context in it.
// Note that this returned DB object can be used only once, so do not assign it to
// a global or package variable for long using.
func (c *Core) Ctx(ctx context.Context) DB {
	if ctx == nil {
		return c.db
	}
	var (
		err        error
		newCore    = &Core{}
		configNode = c.db.GetConfig()
	)
	*newCore = *c
	newCore.ctx = ctx
	newCore.db, err = driverMap[configNode.Type].New(newCore, configNode)
	// Seldom error, just log it.
	if err != nil {
		c.db.GetLogger().Ctx(ctx).Error(err)
	}
	return newCore.db
}

// GetCtx returns the context for current DB.
// It returns `context.Background()` is there's no context previously set.
func (c *Core) GetCtx() context.Context {
	if c.ctx != nil {
		return c.ctx
	}
	return context.TODO()
}

// GetCtxTimeout returns the context and cancel function for specified timeout type.
func (c *Core) GetCtxTimeout(timeoutType int, ctx context.Context) (context.Context, context.CancelFunc) {
	if ctx == nil {
		ctx = c.db.GetCtx()
	} else {
		ctx = context.WithValue(ctx, "WrappedByGetCtxTimeout", nil)
	}
	switch timeoutType {
	case ctxTimeoutTypeExec:
		if c.db.GetConfig().ExecTimeout > 0 {
			return context.WithTimeout(ctx, c.db.GetConfig().ExecTimeout)
		}
	case ctxTimeoutTypeQuery:
		if c.db.GetConfig().QueryTimeout > 0 {
			return context.WithTimeout(ctx, c.db.GetConfig().QueryTimeout)
		}
	case ctxTimeoutTypePrepare:
		if c.db.GetConfig().PrepareTimeout > 0 {
			return context.WithTimeout(ctx, c.db.GetConfig().PrepareTimeout)
		}
	default:
		panic(gerror.Newf("invalid context timeout type: %d", timeoutType))
	}
	return ctx, func() {}
}

// Master creates and returns a connection from master node if master-slave configured.
// It returns the default connection if master-slave not configured.
func (c *Core) Master() (*sql.DB, error) {
	return c.getSqlDb(true, c.schema.Val())
}

// Slave creates and returns a connection from slave node if master-slave configured.
// It returns the default connection if master-slave not configured.
func (c *Core) Slave() (*sql.DB, error) {
	return c.getSqlDb(false, c.schema.Val())
}

// Query commits one query SQL to underlying driver and returns the execution result.
// It is most commonly used for data querying.
func (c *Core) Query(sql string, args ...interface{}) (rows *sql.Rows, err error) {
	link, err := c.db.Slave()
	if err != nil {
		return nil, err
	}
	return c.db.DoQuery(link, sql, args...)
}

// DoQuery commits the sql string and its arguments to underlying driver
// through given link object and returns the execution result.
func (c *Core) DoQuery(link Link, sql string, args ...interface{}) (rows *sql.Rows, err error) {
	sql, args = formatSql(sql, args)
	sql, args = c.db.HandleSqlBeforeCommit(link, sql, args)
	ctx := c.db.GetCtx()
	if c.GetConfig().QueryTimeout > 0 {
		ctx, _ = context.WithTimeout(ctx, c.GetConfig().QueryTimeout)
	}
	mTime1 := gtime.TimestampMilli()
	rows, err = link.QueryContext(ctx, sql, args...)
	mTime2 := gtime.TimestampMilli()
	sqlObj := &Sql{
		Sql:    sql,
		Type:   "DB.QueryContext",
		Args:   args,
		Format: FormatSqlWithArgs(sql, args),
		Error:  err,
		Start:  mTime1,
		End:    mTime2,
		Group:  c.db.GetGroup(),
	}
	c.addSqlToTracing(ctx, sqlObj)
	if c.db.GetDebug() {
		c.writeSqlToLogger(sqlObj)
	}
	if err == nil {
		return rows, nil
	} else {
		err = formatError(err, sql, args...)
	}
	return nil, err
}

// Exec commits one query SQL to underlying driver and returns the execution result.
// It is most commonly used for data inserting and updating.
func (c *Core) Exec(sql string, args ...interface{}) (result sql.Result, err error) {
	link, err := c.db.Master()
	if err != nil {
		return nil, err
	}
	return c.db.DoExec(link, sql, args...)
}

// DoExec commits the sql string and its arguments to underlying driver
// through given link object and returns the execution result.
func (c *Core) DoExec(link Link, sql string, args ...interface{}) (result sql.Result, err error) {
	sql, args = formatSql(sql, args)
	sql, args = c.db.HandleSqlBeforeCommit(link, sql, args)
	ctx := c.db.GetCtx()
	if c.GetConfig().ExecTimeout > 0 {
		var cancelFunc context.CancelFunc
		ctx, cancelFunc = context.WithTimeout(ctx, c.GetConfig().ExecTimeout)
		defer cancelFunc()
	}

	mTime1 := gtime.TimestampMilli()
	if !c.db.GetDryRun() {
		result, err = link.ExecContext(ctx, sql, args...)
	} else {
		result = new(SqlResult)
	}
	mTime2 := gtime.TimestampMilli()
	sqlObj := &Sql{
		Sql:    sql,
		Type:   "DB.ExecContext",
		Args:   args,
		Format: FormatSqlWithArgs(sql, args),
		Error:  err,
		Start:  mTime1,
		End:    mTime2,
		Group:  c.db.GetGroup(),
	}
	c.addSqlToTracing(ctx, sqlObj)
	if c.db.GetDebug() {
		c.writeSqlToLogger(sqlObj)
	}
	return result, formatError(err, sql, args...)
}

// Prepare creates a prepared statement for later queries or executions.
// Multiple queries or executions may be run concurrently from the
// returned statement.
// The caller must call the statement's Close method
// when the statement is no longer needed.
//
// The parameter `execOnMaster` specifies whether executing the sql on master node,
// or else it executes the sql on slave node if master-slave configured.
func (c *Core) Prepare(sql string, execOnMaster ...bool) (*Stmt, error) {
	var (
		err  error
		link Link
	)
	if len(execOnMaster) > 0 && execOnMaster[0] {
		if link, err = c.db.Master(); err != nil {
			return nil, err
		}
	} else {
		if link, err = c.db.Slave(); err != nil {
			return nil, err
		}
	}
	return c.db.DoPrepare(link, sql)
}

// doPrepare calls prepare function on given link object and returns the statement object.
func (c *Core) DoPrepare(link Link, sql string) (*Stmt, error) {
	ctx := c.db.GetCtx()
	if c.GetConfig().PrepareTimeout > 0 {
		// DO NOT USE cancel function in prepare statement.
		ctx, _ = context.WithTimeout(ctx, c.GetConfig().PrepareTimeout)
	}
	var (
		mTime1    = gtime.TimestampMilli()
		stmt, err = link.PrepareContext(ctx, sql)
		mTime2    = gtime.TimestampMilli()
		sqlObj    = &Sql{
			Sql:    sql,
			Type:   "DB.PrepareContext",
			Args:   nil,
			Format: FormatSqlWithArgs(sql, nil),
			Error:  err,
			Start:  mTime1,
			End:    mTime2,
			Group:  c.db.GetGroup(),
		}
	)
	c.addSqlToTracing(ctx, sqlObj)
	if c.db.GetDebug() {
		c.writeSqlToLogger(sqlObj)
	}
	return &Stmt{
		Stmt: stmt,
		core: c,
		sql:  sql,
	}, err
}

// GetAll queries and returns data records from database.
func (c *Core) GetAll(sql string, args ...interface{}) (Result, error) {
	return c.db.DoGetAll(nil, sql, args...)
}

// DoGetAll queries and returns data records from database.
func (c *Core) DoGetAll(link Link, sql string, args ...interface{}) (result Result, err error) {
	if link == nil {
		link, err = c.db.Slave()
		if err != nil {
			return nil, err
		}
	}
	rows, err := c.db.DoQuery(link, sql, args...)
	if err != nil || rows == nil {
		return nil, err
	}
	defer rows.Close()
	return c.db.convertRowsToResult(rows)
}

// GetOne queries and returns one record from database.
func (c *Core) GetOne(sql string, args ...interface{}) (Record, error) {
	list, err := c.db.GetAll(sql, args...)
	if err != nil {
		return nil, err
	}
	if len(list) > 0 {
		return list[0], nil
	}
	return nil, nil
}

// GetArray queries and returns data values as slice from database.
// Note that if there are multiple columns in the result, it returns just one column values randomly.
func (c *Core) GetArray(sql string, args ...interface{}) ([]Value, error) {
	all, err := c.db.DoGetAll(nil, sql, args...)
	if err != nil {
		return nil, err
	}
	return all.Array(), nil
}

// GetStruct queries one record from database and converts it to given struct.
// The parameter `pointer` should be a pointer to struct.
func (c *Core) GetStruct(pointer interface{}, sql string, args ...interface{}) error {
	one, err := c.db.GetOne(sql, args...)
	if err != nil {
		return err
	}
	return one.Struct(pointer)
}

// GetStructs queries records from database and converts them to given struct.
// The parameter `pointer` should be type of struct slice: []struct/[]*struct.
func (c *Core) GetStructs(pointer interface{}, sql string, args ...interface{}) error {
	all, err := c.db.GetAll(sql, args...)
	if err != nil {
		return err
	}
	return all.Structs(pointer)
}

// GetScan queries one or more records from database and converts them to given struct or
// struct array.
//
// If parameter `pointer` is type of struct pointer, it calls GetStruct internally for
// the conversion. If parameter `pointer` is type of slice, it calls GetStructs internally
// for conversion.
func (c *Core) GetScan(pointer interface{}, sql string, args ...interface{}) error {
	t := reflect.TypeOf(pointer)
	k := t.Kind()
	if k != reflect.Ptr {
		return fmt.Errorf("params should be type of pointer, but got: %v", k)
	}
	k = t.Elem().Kind()
	switch k {
	case reflect.Array, reflect.Slice:
		return c.db.GetStructs(pointer, sql, args...)
	case reflect.Struct:
		return c.db.GetStruct(pointer, sql, args...)
	}
	return fmt.Errorf("element type should be type of struct/slice, unsupported: %v", k)
}

// GetValue queries and returns the field value from database.
// The sql should queries only one field from database, or else it returns only one
// field of the result.
func (c *Core) GetValue(sql string, args ...interface{}) (Value, error) {
	one, err := c.db.GetOne(sql, args...)
	if err != nil {
		return gvar.New(nil), err
	}
	for _, v := range one {
		return v, nil
	}
	return gvar.New(nil), nil
}

// GetCount queries and returns the count from database.
func (c *Core) GetCount(sql string, args ...interface{}) (int, error) {
	// If the query fields do not contains function "COUNT",
	// it replaces the sql string and adds the "COUNT" function to the fields.
	if !gregex.IsMatchString(`(?i)SELECT\s+COUNT\(.+\)\s+FROM`, sql) {
		sql, _ = gregex.ReplaceString(`(?i)(SELECT)\s+(.+)\s+(FROM)`, `$1 COUNT($2) $3`, sql)
	}
	value, err := c.db.GetValue(sql, args...)
	if err != nil {
		return 0, err
	}
	return value.Int(), nil
}

// PingMaster pings the master node to check authentication or keeps the connection alive.
func (c *Core) PingMaster() error {
	if master, err := c.db.Master(); err != nil {
		return err
	} else {
		return master.Ping()
	}
}

// PingSlave pings the slave node to check authentication or keeps the connection alive.
func (c *Core) PingSlave() error {
	if slave, err := c.db.Slave(); err != nil {
		return err
	} else {
		return slave.Ping()
	}
}

// Begin starts and returns the transaction object.
// You should call Commit or Rollback functions of the transaction object
// if you no longer use the transaction. Commit or Rollback functions will also
// close the transaction automatically.
func (c *Core) Begin() (*TX, error) {
	if master, err := c.db.Master(); err != nil {
		return nil, err
	} else {
		//ctx := c.db.GetCtx()
		//if c.GetConfig().TranTimeout > 0 {
		//	var cancelFunc context.CancelFunc
		//	ctx, cancelFunc = context.WithTimeout(ctx, c.GetConfig().TranTimeout)
		//	defer cancelFunc()
		//}
		var (
			sqlStr     = "BEGIN"
			mTime1     = gtime.TimestampMilli()
			rawTx, err = master.Begin()
			mTime2     = gtime.TimestampMilli()
			sqlObj     = &Sql{
				Sql:    sqlStr,
				Type:   "DB.Begin",
				Args:   nil,
				Format: sqlStr,
				Error:  err,
				Start:  mTime1,
				End:    mTime2,
				Group:  c.db.GetGroup(),
			}
		)
		c.db.addSqlToTracing(c.db.GetCtx(), sqlObj)
		if c.db.GetDebug() {
			c.db.writeSqlToLogger(sqlObj)
		}
		if err == nil {
			return &TX{
				db:     c.db,
				tx:     rawTx,
				master: master,
			}, nil
		}
		return nil, err
	}
}

// Transaction wraps the transaction logic using function `f`.
// It rollbacks the transaction and returns the error from function `f` if
// it returns non-nil error. It commits the transaction and returns nil if
// function `f` returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function `f`
// as it is automatically handled by this function.
func (c *Core) Transaction(ctx context.Context, f func(ctx context.Context, tx *TX) error) (err error) {
	var tx *TX
	// Check transaction object from context.
	tx = TXFromCtx(ctx)
	if tx != nil {
		return tx.Transaction(ctx, f)
	}
	tx, err = c.db.Begin()
	if err != nil {
		return err
	}
	// Inject transaction object into context.
	ctx = WithTX(ctx, tx)
	defer func() {
		if err == nil {
			if e := recover(); e != nil {
				err = fmt.Errorf("%v", e)
			}
		}
		if err != nil {
			if e := tx.Rollback(); e != nil {
				err = e
			}
		} else {
			if e := tx.Commit(); e != nil {
				err = e
			}
		}
	}()
	err = f(ctx, tx)
	return
}

// Insert does "INSERT INTO ..." statement for the table.
// If there's already one unique record of the data in the table, it returns error.
//
// The parameter `data` can be type of map/gmap/struct/*struct/[]map/[]struct, etc.
// Eg:
// Data(g.Map{"uid": 10000, "name":"john"})
// Data(g.Slice{g.Map{"uid": 10000, "name":"john"}, g.Map{"uid": 20000, "name":"smith"})
//
// The parameter `batch` specifies the batch operation count when given data is slice.
func (c *Core) Insert(table string, data interface{}, batch ...int) (sql.Result, error) {
	if len(batch) > 0 {
		return c.Model(table).Data(data).Batch(batch[0]).Insert()
	}
	return c.Model(table).Data(data).Insert()
}

// InsertIgnore does "INSERT IGNORE INTO ..." statement for the table.
// If there's already one unique record of the data in the table, it ignores the inserting.
//
// The parameter `data` can be type of map/gmap/struct/*struct/[]map/[]struct, etc.
// Eg:
// Data(g.Map{"uid": 10000, "name":"john"})
// Data(g.Slice{g.Map{"uid": 10000, "name":"john"}, g.Map{"uid": 20000, "name":"smith"})
//
// The parameter `batch` specifies the batch operation count when given data is slice.
func (c *Core) InsertIgnore(table string, data interface{}, batch ...int) (sql.Result, error) {
	if len(batch) > 0 {
		return c.Model(table).Data(data).Batch(batch[0]).InsertIgnore()
	}
	return c.Model(table).Data(data).InsertIgnore()
}

// InsertAndGetId performs action Insert and returns the last insert id that automatically generated.
func (c *Core) InsertAndGetId(table string, data interface{}, batch ...int) (int64, error) {
	if len(batch) > 0 {
		return c.Model(table).Data(data).Batch(batch[0]).InsertAndGetId()
	}
	return c.Model(table).Data(data).InsertAndGetId()
}

// Replace does "REPLACE INTO ..." statement for the table.
// If there's already one unique record of the data in the table, it deletes the record
// and inserts a new one.
//
// The parameter `data` can be type of map/gmap/struct/*struct/[]map/[]struct, etc.
// Eg:
// Data(g.Map{"uid": 10000, "name":"john"})
// Data(g.Slice{g.Map{"uid": 10000, "name":"john"}, g.Map{"uid": 20000, "name":"smith"})
//
// The parameter `data` can be type of map/gmap/struct/*struct/[]map/[]struct, etc.
// If given data is type of slice, it then does batch replacing, and the optional parameter
// `batch` specifies the batch operation count.
func (c *Core) Replace(table string, data interface{}, batch ...int) (sql.Result, error) {
	if len(batch) > 0 {
		return c.Model(table).Data(data).Batch(batch[0]).Replace()
	}
	return c.Model(table).Data(data).Replace()
}

// Save does "INSERT INTO ... ON DUPLICATE KEY UPDATE..." statement for the table.
// It updates the record if there's primary or unique index in the saving data,
// or else it inserts a new record into the table.
//
// The parameter `data` can be type of map/gmap/struct/*struct/[]map/[]struct, etc.
// Eg:
// Data(g.Map{"uid": 10000, "name":"john"})
// Data(g.Slice{g.Map{"uid": 10000, "name":"john"}, g.Map{"uid": 20000, "name":"smith"})
//
// If given data is type of slice, it then does batch saving, and the optional parameter
// `batch` specifies the batch operation count.
func (c *Core) Save(table string, data interface{}, batch ...int) (sql.Result, error) {
	if len(batch) > 0 {
		return c.Model(table).Data(data).Batch(batch[0]).Save()
	}
	return c.Model(table).Data(data).Save()
}

// doInsert inserts or updates data for given table.
// This function is usually used for custom interface definition, you do not need call it manually.
// The parameter `data` can be type of map/gmap/struct/*struct/[]map/[]struct, etc.
// Eg:
// Data(g.Map{"uid": 10000, "name":"john"})
// Data(g.Slice{g.Map{"uid": 10000, "name":"john"}, g.Map{"uid": 20000, "name":"smith"})
//
// The parameter `option` values are as follows:
// 0: insert:  just insert, if there's unique/primary key in the data, it returns error;
// 1: replace: if there's unique/primary key in the data, it deletes it from table and inserts a new one;
// 2: save:    if there's unique/primary key in the data, it updates it or else inserts a new one;
// 3: ignore:  if there's unique/primary key in the data, it ignores the inserting;
func (c *Core) DoInsert(link Link, table string, data interface{}, option int, batch ...int) (result sql.Result, err error) {
	table = c.db.QuotePrefixTableName(table)
	var (
		fields       []string
		values       []string
		params       []interface{}
		dataMap      Map
		reflectValue = reflect.ValueOf(data)
		reflectKind  = reflectValue.Kind()
	)
	if reflectKind == reflect.Ptr {
		reflectValue = reflectValue.Elem()
		reflectKind = reflectValue.Kind()
	}
	switch reflectKind {
	case reflect.Slice, reflect.Array:
		return c.db.DoBatchInsert(link, table, data, option, batch...)
	case reflect.Struct:
		if _, ok := data.(apiInterfaces); ok {
			return c.db.DoBatchInsert(link, table, data, option, batch...)
		} else {
			dataMap = ConvertDataForTableRecord(data)
		}
	case reflect.Map:
		dataMap = ConvertDataForTableRecord(data)
	default:
		return result, gerror.New(fmt.Sprint("unsupported data type:", reflectKind))
	}
	if len(dataMap) == 0 {
		return nil, gerror.New("data cannot be empty")
	}
	var (
		charL, charR = c.db.GetChars()
		operation    = GetInsertOperationByOption(option)
		updateStr    = ""
	)
	for k, v := range dataMap {
		fields = append(fields, charL+k+charR)
		if s, ok := v.(Raw); ok {
			values = append(values, gconv.String(s))
		} else {
			values = append(values, "?")
			params = append(params, v)
		}
	}
	if option == insertOptionSave {
		for k, _ := range dataMap {
			// If it's SAVE operation,
			// do not automatically update the creating time.
			if c.isSoftCreatedFiledName(k) {
				continue
			}
			if len(updateStr) > 0 {
				updateStr += ","
			}
			updateStr += fmt.Sprintf(
				"%s%s%s=VALUES(%s%s%s)",
				charL, k, charR,
				charL, k, charR,
			)
		}
		updateStr = fmt.Sprintf("ON DUPLICATE KEY UPDATE %s", updateStr)
	}
	if link == nil {
		if link, err = c.db.Master(); err != nil {
			return nil, err
		}
	}
	return c.db.DoExec(
		link,
		fmt.Sprintf(
			"%s INTO %s(%s) VALUES(%s) %s",
			operation, table, strings.Join(fields, ","),
			strings.Join(values, ","), updateStr,
		),
		params...,
	)
}

// BatchInsert batch inserts data.
// The parameter `list` must be type of slice of map or struct.
func (c *Core) BatchInsert(table string, list interface{}, batch ...int) (sql.Result, error) {
	if len(batch) > 0 {
		return c.Model(table).Data(list).Batch(batch[0]).Insert()
	}
	return c.Model(table).Data(list).Insert()
}

// BatchInsertIgnore batch inserts data with ignore option.
// The parameter `list` must be type of slice of map or struct.
func (c *Core) BatchInsertIgnore(table string, list interface{}, batch ...int) (sql.Result, error) {
	if len(batch) > 0 {
		return c.Model(table).Data(list).Batch(batch[0]).InsertIgnore()
	}
	return c.Model(table).Data(list).InsertIgnore()
}

// BatchReplace batch replaces data.
// The parameter `list` must be type of slice of map or struct.
func (c *Core) BatchReplace(table string, list interface{}, batch ...int) (sql.Result, error) {
	if len(batch) > 0 {
		return c.Model(table).Data(list).Batch(batch[0]).Replace()
	}
	return c.Model(table).Data(list).Replace()
}

// BatchSave batch replaces data.
// The parameter `list` must be type of slice of map or struct.
func (c *Core) BatchSave(table string, list interface{}, batch ...int) (sql.Result, error) {
	if len(batch) > 0 {
		return c.Model(table).Data(list).Batch(batch[0]).Save()
	}
	return c.Model(table).Data(list).Save()
}

// DoBatchInsert batch inserts/replaces/saves data.
// This function is usually used for custom interface definition, you do not need call it manually.
func (c *Core) DoBatchInsert(link Link, table string, list interface{}, option int, batch ...int) (result sql.Result, err error) {
	table = c.db.QuotePrefixTableName(table)
	var (
		keys    []string      // Field names.
		values  []string      // Value holder string array, like: (?,?,?)
		params  []interface{} // Values that will be committed to underlying database driver.
		listMap List          // The data list that passed from caller.
	)
	switch value := list.(type) {
	case Result:
		listMap = value.List()
	case Record:
		listMap = List{value.Map()}
	case List:
		listMap = value
	case Map:
		listMap = List{value}
	default:
		var (
			rv   = reflect.ValueOf(list)
			kind = rv.Kind()
		)
		if kind == reflect.Ptr {
			rv = rv.Elem()
			kind = rv.Kind()
		}
		switch kind {
		// If it's slice type, it then converts it to List type.
		case reflect.Slice, reflect.Array:
			listMap = make(List, rv.Len())
			for i := 0; i < rv.Len(); i++ {
				listMap[i] = ConvertDataForTableRecord(rv.Index(i).Interface())
			}
		case reflect.Map:
			listMap = List{ConvertDataForTableRecord(value)}
		case reflect.Struct:
			if v, ok := value.(apiInterfaces); ok {
				var (
					array = v.Interfaces()
					list  = make(List, len(array))
				)
				for i := 0; i < len(array); i++ {
					list[i] = ConvertDataForTableRecord(array[i])
				}
				listMap = list
			} else {
				listMap = List{ConvertDataForTableRecord(value)}
			}
		default:
			return result, gerror.New(fmt.Sprint("unsupported list type:", kind))
		}
	}
	if len(listMap) < 1 {
		return result, gerror.New("data list cannot be empty")
	}
	if link == nil {
		if link, err = c.db.Master(); err != nil {
			return
		}
	}
	// Handle the field names and place holders.
	for k, _ := range listMap[0] {
		keys = append(keys, k)
	}
	// Prepare the batch result pointer.
	var (
		charL, charR = c.db.GetChars()
		batchResult  = new(SqlResult)
		keysStr      = charL + strings.Join(keys, charR+","+charL) + charR
		operation    = GetInsertOperationByOption(option)
		updateStr    = ""
	)
	if option == insertOptionSave {
		for _, k := range keys {
			// If it's SAVE operation,
			// do not automatically update the creating time.
			if c.isSoftCreatedFiledName(k) {
				continue
			}
			if len(updateStr) > 0 {
				updateStr += ","
			}
			updateStr += fmt.Sprintf(
				"%s%s%s=VALUES(%s%s%s)",
				charL, k, charR,
				charL, k, charR,
			)
		}
		updateStr = fmt.Sprintf("ON DUPLICATE KEY UPDATE %s", updateStr)
	}
	batchNum := defaultBatchNumber
	if len(batch) > 0 && batch[0] > 0 {
		batchNum = batch[0]
	}
	var (
		listMapLen  = len(listMap)
		valueHolder = make([]string, 0)
	)
	for i := 0; i < listMapLen; i++ {
		values = values[:0]
		// Note that the map type is unordered,
		// so it should use slice+key to retrieve the value.
		for _, k := range keys {
			if s, ok := listMap[i][k].(Raw); ok {
				values = append(values, gconv.String(s))
			} else {
				values = append(values, "?")
				params = append(params, listMap[i][k])
			}
		}
		valueHolder = append(valueHolder, "("+gstr.Join(values, ",")+")")
		if len(valueHolder) == batchNum || (i == listMapLen-1 && len(valueHolder) > 0) {
			r, err := c.db.DoExec(
				link,
				fmt.Sprintf(
					"%s INTO %s(%s) VALUES%s %s",
					operation, table, keysStr,
					gstr.Join(valueHolder, ","),
					updateStr,
				),
				params...,
			)
			if err != nil {
				return r, err
			}
			if n, err := r.RowsAffected(); err != nil {
				return r, err
			} else {
				batchResult.result = r
				batchResult.affected += n
			}
			params = params[:0]
			valueHolder = valueHolder[:0]
		}
	}
	return batchResult, nil
}

// Update does "UPDATE ... " statement for the table.
//
// The parameter `data` can be type of string/map/gmap/struct/*struct, etc.
// Eg: "uid=10000", "uid", 10000, g.Map{"uid": 10000, "name":"john"}
//
// The parameter `condition` can be type of string/map/gmap/slice/struct/*struct, etc.
// It is commonly used with parameter `args`.
// Eg:
// "uid=10000",
// "uid", 10000
// "money>? AND name like ?", 99999, "vip_%"
// "status IN (?)", g.Slice{1,2,3}
// "age IN(?,?)", 18, 50
// User{ Id : 1, UserName : "john"}
func (c *Core) Update(table string, data interface{}, condition interface{}, args ...interface{}) (sql.Result, error) {
	return c.Model(table).Data(data).Where(condition, args...).Update()
}

// doUpdate does "UPDATE ... " statement for the table.
// This function is usually used for custom interface definition, you do not need call it manually.
func (c *Core) DoUpdate(link Link, table string, data interface{}, condition string, args ...interface{}) (result sql.Result, err error) {
	table = c.db.QuotePrefixTableName(table)
	var (
		rv   = reflect.ValueOf(data)
		kind = rv.Kind()
	)
	if kind == reflect.Ptr {
		rv = rv.Elem()
		kind = rv.Kind()
	}
	var (
		params  []interface{}
		updates = ""
	)
	switch kind {
	case reflect.Map, reflect.Struct:
		var (
			fields  []string
			dataMap = ConvertDataForTableRecord(data)
		)
		for k, v := range dataMap {
			switch value := v.(type) {
			case *Counter:
				if value.Value != 0 {
					column := k
					if value.Field != "" {
						column = c.db.QuoteWord(value.Field)
					}
					fields = append(fields, fmt.Sprintf("%s=%s+?", column, column))
					params = append(params, value.Value)
				}
			case Counter:
				if value.Value != 0 {
					column := k
					if value.Field != "" {
						column = c.db.QuoteWord(value.Field)
					}
					fields = append(fields, fmt.Sprintf("%s=%s+?", column, column))
					params = append(params, value.Value)
				}
			default:
				if s, ok := v.(Raw); ok {
					fields = append(fields, c.db.QuoteWord(k)+"="+gconv.String(s))
				} else {
					fields = append(fields, c.db.QuoteWord(k)+"=?")
					params = append(params, v)
				}
			}
		}
		updates = strings.Join(fields, ",")
	default:
		updates = gconv.String(data)
	}
	if len(updates) == 0 {
		return nil, gerror.New("data cannot be empty")
	}
	if len(params) > 0 {
		args = append(params, args...)
	}
	// If no link passed, it then uses the master link.
	if link == nil {
		if link, err = c.db.Master(); err != nil {
			return nil, err
		}
	}
	return c.db.DoExec(
		link,
		fmt.Sprintf("UPDATE %s SET %s%s", table, updates, condition),
		args...,
	)
}

// Delete does "DELETE FROM ... " statement for the table.
//
// The parameter `condition` can be type of string/map/gmap/slice/struct/*struct, etc.
// It is commonly used with parameter `args`.
// Eg:
// "uid=10000",
// "uid", 10000
// "money>? AND name like ?", 99999, "vip_%"
// "status IN (?)", g.Slice{1,2,3}
// "age IN(?,?)", 18, 50
// User{ Id : 1, UserName : "john"}
func (c *Core) Delete(table string, condition interface{}, args ...interface{}) (result sql.Result, err error) {
	return c.Model(table).Where(condition, args...).Delete()
}

// DoDelete does "DELETE FROM ... " statement for the table.
// This function is usually used for custom interface definition, you do not need call it manually.
func (c *Core) DoDelete(link Link, table string, condition string, args ...interface{}) (result sql.Result, err error) {
	if link == nil {
		if link, err = c.db.Master(); err != nil {
			return nil, err
		}
	}
	table = c.db.QuotePrefixTableName(table)
	return c.db.DoExec(link, fmt.Sprintf("DELETE FROM %s%s", table, condition), args...)
}

// convertRowsToResult converts underlying data record type sql.Rows to Result type.
func (c *Core) convertRowsToResult(rows *sql.Rows) (Result, error) {
	if !rows.Next() {
		return nil, nil
	}
	// Column names and types.
	columns, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}
	columnTypes := make([]string, len(columns))
	columnNames := make([]string, len(columns))
	for k, v := range columns {
		columnTypes[k] = v.DatabaseTypeName()
		columnNames[k] = v.Name()
	}
	var (
		values   = make([]interface{}, len(columnNames))
		records  = make(Result, 0)
		scanArgs = make([]interface{}, len(values))
	)
	for i := range values {
		scanArgs[i] = &values[i]
	}
	for {
		if err := rows.Scan(scanArgs...); err != nil {
			return records, err
		}
		row := make(Record)
		for i, value := range values {
			if value == nil {
				row[columnNames[i]] = gvar.New(nil)
			} else {
				row[columnNames[i]] = gvar.New(c.db.convertFieldValueToLocalValue(value, columnTypes[i]))
			}
		}
		records = append(records, row)
		if !rows.Next() {
			break
		}
	}
	return records, nil
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
// It just returns the pointer address.
//
// Note that this interface implements mainly for workaround for a json infinite loop bug
// of Golang version < v1.14.
func (c *Core) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`%+v`, c)), nil
}

// writeSqlToLogger outputs the sql object to logger.
// It is enabled only if configuration "debug" is true.
func (c *Core) writeSqlToLogger(v *Sql) {
	s := fmt.Sprintf("[%3d ms] [%s] %s", v.End-v.Start, v.Group, v.Format)
	if v.Error != nil {
		s += "\nError: " + v.Error.Error()
		c.logger.Ctx(c.db.GetCtx()).Error(s)
	} else {
		c.logger.Ctx(c.db.GetCtx()).Debug(s)
	}
}

// HasTable determine whether the table name exists in the database.
func (c *Core) HasTable(name string) (bool, error) {
	tableList, err := c.db.Tables()
	if err != nil {
		return false, err
	}
	for _, table := range tableList {
		if table == name {
			return true, nil
		}
	}
	return false, nil
}

// isSoftCreatedFiledName checks and returns whether given filed name is an automatic-filled created time.
func (c *Core) isSoftCreatedFiledName(fieldName string) bool {
	if fieldName == "" {
		return false
	}
	if config := c.db.GetConfig(); config.CreatedAt != "" {
		if utils.EqualFoldWithoutChars(fieldName, config.CreatedAt) {
			return true
		}
		return gstr.InArray(append([]string{config.CreatedAt}, createdFiledNames...), fieldName)
	}
	for _, v := range createdFiledNames {
		if utils.EqualFoldWithoutChars(fieldName, v) {
			return true
		}
	}
	return false
}
