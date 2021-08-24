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
	"github.com/gogf/gf/errors/gcode"
	"reflect"
	"strings"

	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/internal/utils"
	"github.com/gogf/gf/text/gstr"

	"github.com/gogf/gf/container/gvar"
	"github.com/gogf/gf/text/gregex"
	"github.com/gogf/gf/util/gconv"
)

// GetCore returns the underlying *Core object.
func (c *Core) GetCore() *Core {
	return c
}

// Ctx is a chaining function, which creates and returns a new DB that is a shallow copy
// of current DB object and with given context in it.
// Note that this returned DB object can be used only once, so do not assign it to
// a global or package variable for long using.
func (c *Core) Ctx(ctx context.Context) DB {
	if ctx == nil {
		return c.db
	}
	// It is already set context in previous chaining operation.
	if c.ctx != nil {
		return c.db
	}
	ctx = context.WithValue(ctx, ctxStrictKeyName, 1)
	// It makes a shallow copy of current db and changes its context for next chaining operation.
	var (
		err        error
		newCore    = &Core{}
		configNode = c.db.GetConfig()
	)
	*newCore = *c
	newCore.ctx = ctx
	// It creates a new DB object, which is commonly a wrapper for object `Core`.
	newCore.db, err = driverMap[configNode.Type].New(newCore, configNode)
	if err != nil {
		// It is really a serious error here.
		// Do not let it continue.
		panic(err)
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
		ctx = c.GetCtx()
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
		panic(gerror.NewCodef(gcode.CodeInvalidParameter, "invalid context timeout type: %d", timeoutType))
	}
	return ctx, func() {}
}

// Master creates and returns a connection from master node if master-slave configured.
// It returns the default connection if master-slave not configured.
func (c *Core) Master(schema ...string) (*sql.DB, error) {
	useSchema := ""
	if len(schema) > 0 && schema[0] != "" {
		useSchema = schema[0]
	} else {
		useSchema = c.schema.Val()
	}
	return c.getSqlDb(true, useSchema)
}

// Slave creates and returns a connection from slave node if master-slave configured.
// It returns the default connection if master-slave not configured.
func (c *Core) Slave(schema ...string) (*sql.DB, error) {
	useSchema := ""
	if len(schema) > 0 && schema[0] != "" {
		useSchema = schema[0]
	} else {
		useSchema = c.schema.Val()
	}
	return c.getSqlDb(false, useSchema)
}

// GetAll queries and returns data records from database.
func (c *Core) GetAll(sql string, args ...interface{}) (Result, error) {
	return c.db.DoGetAll(c.GetCtx(), nil, sql, args...)
}

// DoGetAll queries and returns data records from database.
func (c *Core) DoGetAll(ctx context.Context, link Link, sql string, args ...interface{}) (result Result, err error) {
	rows, err := c.db.DoQuery(ctx, link, sql, args...)
	if err != nil || rows == nil {
		return nil, err
	}
	defer rows.Close()
	return c.convertRowsToResult(rows)
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
	all, err := c.db.DoGetAll(c.GetCtx(), nil, sql, args...)
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
		return c.db.GetCore().GetStructs(pointer, sql, args...)
	case reflect.Struct:
		return c.db.GetCore().GetStruct(pointer, sql, args...)
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

// Union does "(SELECT xxx FROM xxx) UNION (SELECT xxx FROM xxx) ..." statement.
func (c *Core) Union(unions ...*Model) *Model {
	return c.doUnion(unionTypeNormal, unions...)
}

// UnionAll does "(SELECT xxx FROM xxx) UNION ALL (SELECT xxx FROM xxx) ..." statement.
func (c *Core) UnionAll(unions ...*Model) *Model {
	return c.doUnion(unionTypeAll, unions...)
}

func (c *Core) doUnion(unionType int, unions ...*Model) *Model {
	var (
		unionTypeStr   string
		composedSqlStr string
		composedArgs   = make([]interface{}, 0)
	)
	if unionType == unionTypeAll {
		unionTypeStr = "UNION ALL"
	} else {
		unionTypeStr = "UNION"
	}
	for _, v := range unions {
		sqlWithHolder, holderArgs := v.getFormattedSqlAndArgs(queryTypeNormal, false)
		if composedSqlStr == "" {
			composedSqlStr += fmt.Sprintf(`(%s)`, sqlWithHolder)
		} else {
			composedSqlStr += fmt.Sprintf(` %s (%s)`, unionTypeStr, sqlWithHolder)
		}
		composedArgs = append(composedArgs, holderArgs...)
	}
	return c.db.Raw(composedSqlStr, composedArgs...)
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

// DoInsert inserts or updates data forF given table.
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
func (c *Core) DoInsert(ctx context.Context, link Link, table string, list List, option DoInsertOption) (result sql.Result, err error) {
	var (
		keys           []string      // Field names.
		values         []string      // Value holder string array, like: (?,?,?)
		params         []interface{} // Values that will be committed to underlying database driver.
		onDuplicateStr string        // onDuplicateStr is used in "ON DUPLICATE KEY UPDATE" statement.
	)
	// Handle the field names and place holders.
	for k, _ := range list[0] {
		keys = append(keys, k)
	}
	// Prepare the batch result pointer.
	var (
		charL, charR = c.db.GetChars()
		batchResult  = new(SqlResult)
		keysStr      = charL + strings.Join(keys, charR+","+charL) + charR
		operation    = GetInsertOperationByOption(option.InsertOption)
	)
	if option.InsertOption == insertOptionSave {
		onDuplicateStr = c.formatOnDuplicate(keys, option)
	}
	var (
		listLength  = len(list)
		valueHolder = make([]string, 0)
	)
	for i := 0; i < listLength; i++ {
		values = values[:0]
		// Note that the map type is unordered,
		// so it should use slice+key to retrieve the value.
		for _, k := range keys {
			if s, ok := list[i][k].(Raw); ok {
				values = append(values, gconv.String(s))
			} else {
				values = append(values, "?")
				params = append(params, list[i][k])
			}
		}
		valueHolder = append(valueHolder, "("+gstr.Join(values, ",")+")")
		// Batch package checks: It meets the batch number or it is the last element.
		if len(valueHolder) == option.BatchCount || (i == listLength-1 && len(valueHolder) > 0) {
			r, err := c.db.DoExec(ctx, link, fmt.Sprintf(
				"%s INTO %s(%s) VALUES%s %s",
				operation, c.QuotePrefixTableName(table), keysStr,
				gstr.Join(valueHolder, ","),
				onDuplicateStr,
			), params...)
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

func (c *Core) formatOnDuplicate(columns []string, option DoInsertOption) string {
	var (
		onDuplicateStr string
	)
	if option.OnDuplicateStr != "" {
		onDuplicateStr = option.OnDuplicateStr
	} else if len(option.OnDuplicateMap) > 0 {
		for k, v := range option.OnDuplicateMap {
			if len(onDuplicateStr) > 0 {
				onDuplicateStr += ","
			}
			switch v.(type) {
			case Raw, *Raw:
				onDuplicateStr += fmt.Sprintf(
					"%s=%s",
					c.QuoteWord(k),
					v,
				)
			default:
				onDuplicateStr += fmt.Sprintf(
					"%s=VALUES(%s)",
					c.QuoteWord(k),
					c.QuoteWord(gconv.String(v)),
				)
			}
		}
	} else {
		for _, column := range columns {
			// If it's SAVE operation, do not automatically update the creating time.
			if c.isSoftCreatedFieldName(column) {
				continue
			}
			if len(onDuplicateStr) > 0 {
				onDuplicateStr += ","
			}
			onDuplicateStr += fmt.Sprintf(
				"%s=VALUES(%s)",
				c.QuoteWord(column),
				c.QuoteWord(column),
			)
		}
	}
	return fmt.Sprintf("ON DUPLICATE KEY UPDATE %s", onDuplicateStr)
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

// DoUpdate does "UPDATE ... " statement for the table.
// This function is usually used for custom interface definition, you do not need call it manually.
func (c *Core) DoUpdate(ctx context.Context, link Link, table string, data interface{}, condition string, args ...interface{}) (result sql.Result, err error) {
	table = c.QuotePrefixTableName(table)
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
			fields         []string
			dataMap        = ConvertDataForTableRecord(data)
			counterHandler = func(column string, counter Counter) {
				if counter.Value != 0 {
					var (
						column    = c.QuoteWord(column)
						columnRef = c.QuoteWord(counter.Field)
						columnVal = counter.Value
						operator  = "+"
					)
					if columnVal < 0 {
						operator = "-"
						columnVal = -columnVal
					}
					fields = append(fields, fmt.Sprintf("%s=%s%s?", column, columnRef, operator))
					params = append(params, columnVal)
				}
			}
		)

		for k, v := range dataMap {
			switch value := v.(type) {
			case *Counter:
				counterHandler(k, *value)

			case Counter:
				counterHandler(k, value)

			default:
				if s, ok := v.(Raw); ok {
					fields = append(fields, c.QuoteWord(k)+"="+gconv.String(s))
				} else {
					fields = append(fields, c.QuoteWord(k)+"=?")
					params = append(params, v)
				}
			}
		}
		updates = strings.Join(fields, ",")

	default:
		updates = gconv.String(data)
	}
	if len(updates) == 0 {
		return nil, gerror.NewCode(gcode.CodeMissingParameter, "data cannot be empty")
	}
	if len(params) > 0 {
		args = append(params, args...)
	}
	// If no link passed, it then uses the master link.
	if link == nil {
		if link, err = c.MasterLink(); err != nil {
			return nil, err
		}
	}
	return c.db.DoExec(ctx, link, fmt.Sprintf("UPDATE %s SET %s%s", table, updates, condition), args...)
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
func (c *Core) DoDelete(ctx context.Context, link Link, table string, condition string, args ...interface{}) (result sql.Result, err error) {
	if link == nil {
		if link, err = c.MasterLink(); err != nil {
			return nil, err
		}
	}
	table = c.QuotePrefixTableName(table)
	return c.db.DoExec(ctx, link, fmt.Sprintf("DELETE FROM %s%s", table, condition), args...)
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
		result   = make(Result, 0)
		scanArgs = make([]interface{}, len(values))
	)
	for i := range values {
		scanArgs[i] = &values[i]
	}
	for {
		if err := rows.Scan(scanArgs...); err != nil {
			return result, err
		}
		record := Record{}
		for i, value := range values {
			if value == nil {
				record[columnNames[i]] = gvar.New(nil)
			} else {
				record[columnNames[i]] = gvar.New(c.convertFieldValueToLocalValue(value, columnTypes[i]))
			}
		}
		result = append(result, record)
		if !rows.Next() {
			break
		}
	}
	return result, nil
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
func (c *Core) writeSqlToLogger(ctx context.Context, sql *Sql) {
	var transactionIdStr string
	if sql.IsTransaction {
		if v := ctx.Value(transactionIdForLoggerCtx); v != nil {
			transactionIdStr = fmt.Sprintf(`[%d] `, v.(uint64))
		}
	}
	s := fmt.Sprintf("[%3d ms] [%s] %s%s", sql.End-sql.Start, sql.Group, transactionIdStr, sql.Format)
	if sql.Error != nil {
		s += "\nError: " + sql.Error.Error()
		c.logger.Ctx(ctx).Error(s)
	} else {
		c.logger.Ctx(ctx).Debug(s)
	}
}

// HasTable determine whether the table name exists in the database.
func (c *Core) HasTable(name string) (bool, error) {
	tableList, err := c.db.Tables(c.GetCtx())
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

// isSoftCreatedFieldName checks and returns whether given filed name is an automatic-filled created time.
func (c *Core) isSoftCreatedFieldName(fieldName string) bool {
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
