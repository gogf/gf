// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//

package gdb

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/gogf/gf/container/gvar"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/text/gregex"
	"github.com/gogf/gf/util/gconv"
)

const (
	gPATH_FILTER_KEY = "/database/gdb/gdb"
)

var (
	// lastOperatorReg is the regular expression object for a string
	// which has operator at its tail.
	lastOperatorReg = regexp.MustCompile(`[<>=]+\s*$`)
)

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
	link, err := c.DB.Slave()
	if err != nil {
		return nil, err
	}
	return c.DB.DoQuery(link, sql, args...)
}

// DoQuery commits the sql string and its arguments to underlying driver
// through given link object and returns the execution result.
func (c *Core) DoQuery(link Link, sql string, args ...interface{}) (rows *sql.Rows, err error) {
	sql, args = formatSql(sql, args)
	sql, args = c.DB.HandleSqlBeforeCommit(link, sql, args)
	if c.DB.GetDebug() {
		mTime1 := gtime.TimestampMilli()
		rows, err = link.Query(sql, args...)
		mTime2 := gtime.TimestampMilli()
		s := &Sql{
			Sql:    sql,
			Args:   args,
			Format: FormatSqlWithArgs(sql, args),
			Error:  err,
			Start:  mTime1,
			End:    mTime2,
		}
		c.writeSqlToLogger(s)
	} else {
		rows, err = link.Query(sql, args...)
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
	link, err := c.DB.Master()
	if err != nil {
		return nil, err
	}
	return c.DB.DoExec(link, sql, args...)
}

// DoExec commits the sql string and its arguments to underlying driver
// through given link object and returns the execution result.
func (c *Core) DoExec(link Link, sql string, args ...interface{}) (result sql.Result, err error) {
	sql, args = formatSql(sql, args)
	sql, args = c.DB.HandleSqlBeforeCommit(link, sql, args)
	if c.DB.GetDebug() {
		mTime1 := gtime.TimestampMilli()
		if !c.DB.GetDryRun() {
			result, err = link.Exec(sql, args...)
		} else {
			result = new(SqlResult)
		}
		mTime2 := gtime.TimestampMilli()
		s := &Sql{
			Sql:    sql,
			Args:   args,
			Format: FormatSqlWithArgs(sql, args),
			Error:  err,
			Start:  mTime1,
			End:    mTime2,
		}
		c.writeSqlToLogger(s)
	} else {
		if !c.DB.GetDryRun() {
			result, err = link.Exec(sql, args...)
		} else {
			result = new(SqlResult)
		}
	}
	return result, formatError(err, sql, args...)
}

// Prepare creates a prepared statement for later queries or executions.
// Multiple queries or executions may be run concurrently from the
// returned statement.
// The caller must call the statement's Close method
// when the statement is no longer needed.
//
// The parameter <execOnMaster> specifies whether executing the sql on master node,
// or else it executes the sql on slave node if master-slave configured.
func (c *Core) Prepare(sql string, execOnMaster ...bool) (*sql.Stmt, error) {
	err := (error)(nil)
	link := (Link)(nil)
	if len(execOnMaster) > 0 && execOnMaster[0] {
		if link, err = c.DB.Master(); err != nil {
			return nil, err
		}
	} else {
		if link, err = c.DB.Slave(); err != nil {
			return nil, err
		}
	}
	return c.DB.DoPrepare(link, sql)
}

// doPrepare calls prepare function on given link object and returns the statement object.
func (c *Core) DoPrepare(link Link, sql string) (*sql.Stmt, error) {
	return link.Prepare(sql)
}

// GetAll queries and returns data records from database.
func (c *Core) GetAll(sql string, args ...interface{}) (Result, error) {
	return c.DB.DoGetAll(nil, sql, args...)
}

// DoGetAll queries and returns data records from database.
func (c *Core) DoGetAll(link Link, sql string, args ...interface{}) (result Result, err error) {
	if link == nil {
		link, err = c.DB.Slave()
		if err != nil {
			return nil, err
		}
	}
	rows, err := c.DB.DoQuery(link, sql, args...)
	if err != nil || rows == nil {
		return nil, err
	}
	defer rows.Close()
	return c.DB.rowsToResult(rows)
}

// GetOne queries and returns one record from database.
func (c *Core) GetOne(sql string, args ...interface{}) (Record, error) {
	list, err := c.DB.GetAll(sql, args...)
	if err != nil {
		return nil, err
	}
	if len(list) > 0 {
		return list[0], nil
	}
	return nil, nil
}

// GetArray queries and returns data values as slice from database.
// Note that if there're multiple columns in the result, it returns just one column values randomly.
func (c *Core) GetArray(sql string, args ...interface{}) ([]Value, error) {
	all, err := c.DB.DoGetAll(nil, sql, args...)
	if err != nil {
		return nil, err
	}
	return all.Array(), nil
}

// GetStruct queries one record from database and converts it to given struct.
// The parameter <pointer> should be a pointer to struct.
func (c *Core) GetStruct(pointer interface{}, sql string, args ...interface{}) error {
	one, err := c.DB.GetOne(sql, args...)
	if err != nil {
		return err
	}
	if len(one) == 0 {
		return ErrNoRows
	}
	return one.Struct(pointer)
}

// GetStructs queries records from database and converts them to given struct.
// The parameter <pointer> should be type of struct slice: []struct/[]*struct.
func (c *Core) GetStructs(pointer interface{}, sql string, args ...interface{}) error {
	all, err := c.DB.GetAll(sql, args...)
	if err != nil {
		return err
	}
	if len(all) == 0 {
		return ErrNoRows
	}
	return all.Structs(pointer)
}

// GetScan queries one or more records from database and converts them to given struct or
// struct array.
//
// If parameter <pointer> is type of struct pointer, it calls GetStruct internally for
// the conversion. If parameter <pointer> is type of slice, it calls GetStructs internally
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
		return c.DB.GetStructs(pointer, sql, args...)
	case reflect.Struct:
		return c.DB.GetStruct(pointer, sql, args...)
	}
	return fmt.Errorf("element type should be type of struct/slice, unsupported: %v", k)
}

// GetValue queries and returns the field value from database.
// The sql should queries only one field from database, or else it returns only one
// field of the result.
func (c *Core) GetValue(sql string, args ...interface{}) (Value, error) {
	one, err := c.DB.GetOne(sql, args...)
	if err != nil {
		return nil, err
	}
	for _, v := range one {
		return v, nil
	}
	return nil, nil
}

// GetCount queries and returns the count from database.
func (c *Core) GetCount(sql string, args ...interface{}) (int, error) {
	// If the query fields do not contains function "COUNT",
	// it replaces the sql string and adds the "COUNT" function to the fields.
	if !gregex.IsMatchString(`(?i)SELECT\s+COUNT\(.+\)\s+FROM`, sql) {
		sql, _ = gregex.ReplaceString(`(?i)(SELECT)\s+(.+)\s+(FROM)`, `$1 COUNT($2) $3`, sql)
	}
	value, err := c.DB.GetValue(sql, args...)
	if err != nil {
		return 0, err
	}
	return value.Int(), nil
}

// PingMaster pings the master node to check authentication or keeps the connection alive.
func (c *Core) PingMaster() error {
	if master, err := c.DB.Master(); err != nil {
		return err
	} else {
		return master.Ping()
	}
}

// PingSlave pings the slave node to check authentication or keeps the connection alive.
func (c *Core) PingSlave() error {
	if slave, err := c.DB.Slave(); err != nil {
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
	if master, err := c.DB.Master(); err != nil {
		return nil, err
	} else {
		if tx, err := master.Begin(); err == nil {
			return &TX{
				db:     c.DB,
				tx:     tx,
				master: master,
			}, nil
		} else {
			return nil, err
		}
	}
}

// Insert does "INSERT INTO ..." statement for the table.
// If there's already one unique record of the data in the table, it returns error.
//
// The parameter <data> can be type of map/gmap/struct/*struct/[]map/[]struct, etc.
// Eg:
// Data(g.Map{"uid": 10000, "name":"john"})
// Data(g.Slice{g.Map{"uid": 10000, "name":"john"}, g.Map{"uid": 20000, "name":"smith"})
//
// The parameter <batch> specifies the batch operation count when given data is slice.
func (c *Core) Insert(table string, data interface{}, batch ...int) (sql.Result, error) {
	return c.DB.DoInsert(nil, table, data, gINSERT_OPTION_DEFAULT, batch...)
}

// InsertIgnore does "INSERT IGNORE INTO ..." statement for the table.
// If there's already one unique record of the data in the table, it ignores the inserting.
//
// The parameter <data> can be type of map/gmap/struct/*struct/[]map/[]struct, etc.
// Eg:
// Data(g.Map{"uid": 10000, "name":"john"})
// Data(g.Slice{g.Map{"uid": 10000, "name":"john"}, g.Map{"uid": 20000, "name":"smith"})
//
// The parameter <batch> specifies the batch operation count when given data is slice.
func (c *Core) InsertIgnore(table string, data interface{}, batch ...int) (sql.Result, error) {
	return c.DB.DoInsert(nil, table, data, gINSERT_OPTION_IGNORE, batch...)
}

// Replace does "REPLACE INTO ..." statement for the table.
// If there's already one unique record of the data in the table, it deletes the record
// and inserts a new one.
//
// The parameter <data> can be type of map/gmap/struct/*struct/[]map/[]struct, etc.
// Eg:
// Data(g.Map{"uid": 10000, "name":"john"})
// Data(g.Slice{g.Map{"uid": 10000, "name":"john"}, g.Map{"uid": 20000, "name":"smith"})
//
// The parameter <data> can be type of map/gmap/struct/*struct/[]map/[]struct, etc.
// If given data is type of slice, it then does batch replacing, and the optional parameter
// <batch> specifies the batch operation count.
func (c *Core) Replace(table string, data interface{}, batch ...int) (sql.Result, error) {
	return c.DB.DoInsert(nil, table, data, gINSERT_OPTION_REPLACE, batch...)
}

// Save does "INSERT INTO ... ON DUPLICATE KEY UPDATE..." statement for the table.
// It updates the record if there's primary or unique index in the saving data,
// or else it inserts a new record into the table.
//
// The parameter <data> can be type of map/gmap/struct/*struct/[]map/[]struct, etc.
// Eg:
// Data(g.Map{"uid": 10000, "name":"john"})
// Data(g.Slice{g.Map{"uid": 10000, "name":"john"}, g.Map{"uid": 20000, "name":"smith"})
//
// If given data is type of slice, it then does batch saving, and the optional parameter
// <batch> specifies the batch operation count.
func (c *Core) Save(table string, data interface{}, batch ...int) (sql.Result, error) {
	return c.DB.DoInsert(nil, table, data, gINSERT_OPTION_SAVE, batch...)
}

// doInsert inserts or updates data for given table.
//
// The parameter <data> can be type of map/gmap/struct/*struct/[]map/[]struct, etc.
// Eg:
// Data(g.Map{"uid": 10000, "name":"john"})
// Data(g.Slice{g.Map{"uid": 10000, "name":"john"}, g.Map{"uid": 20000, "name":"smith"})
//
// The parameter <option> values are as follows:
// 0: insert:  just insert, if there's unique/primary key in the data, it returns error;
// 1: replace: if there's unique/primary key in the data, it deletes it from table and inserts a new one;
// 2: save:    if there's unique/primary key in the data, it updates it or else inserts a new one;
// 3: ignore:  if there's unique/primary key in the data, it ignores the inserting;
func (c *Core) DoInsert(link Link, table string, data interface{}, option int, batch ...int) (result sql.Result, err error) {
	table = c.DB.QuotePrefixTableName(table)
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
		return c.DB.DoBatchInsert(link, table, data, option, batch...)
	case reflect.Map, reflect.Struct:
		dataMap = DataToMapDeep(data)
	default:
		return result, errors.New(fmt.Sprint("unsupported data type:", reflectKind))
	}
	if len(dataMap) == 0 {
		return nil, errors.New("data cannot be empty")
	}
	var (
		charL, charR = c.DB.GetChars()
		operation    = GetInsertOperationByOption(option)
		updateStr    = ""
	)
	for k, v := range dataMap {
		fields = append(fields, charL+k+charR)
		values = append(values, "?")
		params = append(params, v)
	}
	if option == gINSERT_OPTION_SAVE {
		for k, _ := range dataMap {
			// If it's SAVE operation,
			// do not automatically update the creating time.
			if k == gSOFT_FIELD_NAME_CREATE {
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
		if link, err = c.DB.Master(); err != nil {
			return nil, err
		}
	}
	return c.DB.DoExec(
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
// The parameter <list> must be type of slice of map or struct.
func (c *Core) BatchInsert(table string, list interface{}, batch ...int) (sql.Result, error) {
	return c.DB.DoBatchInsert(nil, table, list, gINSERT_OPTION_DEFAULT, batch...)
}

// BatchInsert batch inserts data with ignore option.
// The parameter <list> must be type of slice of map or struct.
func (c *Core) BatchInsertIgnore(table string, list interface{}, batch ...int) (sql.Result, error) {
	return c.DB.DoBatchInsert(nil, table, list, gINSERT_OPTION_IGNORE, batch...)
}

// BatchReplace batch replaces data.
// The parameter <list> must be type of slice of map or struct.
func (c *Core) BatchReplace(table string, list interface{}, batch ...int) (sql.Result, error) {
	return c.DB.DoBatchInsert(nil, table, list, gINSERT_OPTION_REPLACE, batch...)
}

// BatchSave batch replaces data.
// The parameter <list> must be type of slice of map or struct.
func (c *Core) BatchSave(table string, list interface{}, batch ...int) (sql.Result, error) {
	return c.DB.DoBatchInsert(nil, table, list, gINSERT_OPTION_SAVE, batch...)
}

// DoBatchInsert batch inserts/replaces/saves data.
func (c *Core) DoBatchInsert(link Link, table string, list interface{}, option int, batch ...int) (result sql.Result, err error) {
	table = c.DB.QuotePrefixTableName(table)
	var (
		keys    []string
		values  []string
		params  []interface{}
		listMap List
	)
	switch v := list.(type) {
	case Result:
		listMap = v.List()
	case Record:
		listMap = List{v.Map()}
	case List:
		listMap = v
	case Map:
		listMap = List{v}
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
				listMap[i] = DataToMapDeep(rv.Index(i).Interface())
			}
		case reflect.Map, reflect.Struct:
			listMap = List{DataToMapDeep(list)}
		default:
			return result, errors.New(fmt.Sprint("unsupported list type:", kind))
		}
	}
	if len(listMap) < 1 {
		return result, errors.New("data list cannot be empty")
	}
	if link == nil {
		if link, err = c.DB.Master(); err != nil {
			return
		}
	}
	// Handle the field names and place holders.
	holders := []string(nil)
	for k, _ := range listMap[0] {
		keys = append(keys, k)
		holders = append(holders, "?")
	}
	// Prepare the batch result pointer.
	var (
		charL, charR   = c.DB.GetChars()
		batchResult    = new(SqlResult)
		keysStr        = charL + strings.Join(keys, charR+","+charL) + charR
		valueHolderStr = "(" + strings.Join(holders, ",") + ")"
		operation      = GetInsertOperationByOption(option)
		updateStr      = ""
	)
	if option == gINSERT_OPTION_SAVE {
		for _, k := range keys {
			// If it's SAVE operation,
			// do not automatically update the creating time.
			if k == gSOFT_FIELD_NAME_CREATE {
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
	batchNum := gDEFAULT_BATCH_NUM
	if len(batch) > 0 && batch[0] > 0 {
		batchNum = batch[0]
	}
	listMapLen := len(listMap)
	for i := 0; i < listMapLen; i++ {
		// Note that the map type is unordered,
		// so it should use slice+key to retrieve the value.
		for _, k := range keys {
			params = append(params, listMap[i][k])
		}
		values = append(values, valueHolderStr)
		if len(values) == batchNum || (i == listMapLen-1 && len(values) > 0) {
			r, err := c.DB.DoExec(
				link,
				fmt.Sprintf(
					"%s INTO %s(%s) VALUES%s %s",
					operation,
					table,
					keysStr,
					strings.Join(values, ","),
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
			values = values[:0]
		}
	}
	return batchResult, nil
}

// Update does "UPDATE ... " statement for the table.
//
// The parameter <data> can be type of string/map/gmap/struct/*struct, etc.
// Eg: "uid=10000", "uid", 10000, g.Map{"uid": 10000, "name":"john"}
//
// The parameter <condition> can be type of string/map/gmap/slice/struct/*struct, etc.
// It is commonly used with parameter <args>.
// Eg:
// "uid=10000",
// "uid", 10000
// "money>? AND name like ?", 99999, "vip_%"
// "status IN (?)", g.Slice{1,2,3}
// "age IN(?,?)", 18, 50
// User{ Id : 1, UserName : "john"}
func (c *Core) Update(table string, data interface{}, condition interface{}, args ...interface{}) (sql.Result, error) {
	newWhere, newArgs := formatWhere(c.DB, condition, args, false)
	if newWhere != "" {
		newWhere = " WHERE " + newWhere
	}
	return c.DB.DoUpdate(nil, table, data, newWhere, newArgs...)
}

// doUpdate does "UPDATE ... " statement for the table.
// Also see Update.
func (c *Core) DoUpdate(link Link, table string, data interface{}, condition string, args ...interface{}) (result sql.Result, err error) {
	table = c.DB.QuotePrefixTableName(table)
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
			dataMap = DataToMapDeep(data)
		)
		for k, v := range dataMap {
			fields = append(fields, c.DB.QuoteWord(k)+"=?")
			params = append(params, v)
		}
		updates = strings.Join(fields, ",")
	default:
		updates = gconv.String(data)
	}
	if len(updates) == 0 {
		return nil, errors.New("data cannot be empty")
	}
	if len(params) > 0 {
		args = append(params, args...)
	}
	// If no link passed, it then uses the master link.
	if link == nil {
		if link, err = c.DB.Master(); err != nil {
			return nil, err
		}
	}
	return c.DB.DoExec(
		link,
		fmt.Sprintf("UPDATE %s SET %s%s", table, updates, condition),
		args...,
	)
}

// Delete does "DELETE FROM ... " statement for the table.
//
// The parameter <condition> can be type of string/map/gmap/slice/struct/*struct, etc.
// It is commonly used with parameter <args>.
// Eg:
// "uid=10000",
// "uid", 10000
// "money>? AND name like ?", 99999, "vip_%"
// "status IN (?)", g.Slice{1,2,3}
// "age IN(?,?)", 18, 50
// User{ Id : 1, UserName : "john"}
func (c *Core) Delete(table string, condition interface{}, args ...interface{}) (result sql.Result, err error) {
	newWhere, newArgs := formatWhere(c.DB, condition, args, false)
	if newWhere != "" {
		newWhere = " WHERE " + newWhere
	}
	return c.DB.DoDelete(nil, table, newWhere, newArgs...)
}

// DoDelete does "DELETE FROM ... " statement for the table.
// Also see Delete.
func (c *Core) DoDelete(link Link, table string, condition string, args ...interface{}) (result sql.Result, err error) {
	if link == nil {
		if link, err = c.DB.Master(); err != nil {
			return nil, err
		}
	}
	table = c.DB.QuotePrefixTableName(table)
	return c.DB.DoExec(link, fmt.Sprintf("DELETE FROM %s%s", table, condition), args...)
}

// rowsToResult converts underlying data record type sql.Rows to Result type.
func (c *Core) rowsToResult(rows *sql.Rows) (Result, error) {
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
	values := make([]sql.RawBytes, len(columnNames))
	records := make(Result, 0)
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	for {
		if err := rows.Scan(scanArgs...); err != nil {
			return records, err
		}
		// Creates a new row object.
		row := make(Record)
		// Note that the internal looping variable <value> is type of []byte,
		// which points to the same memory address. So it should do a copy.
		for i, value := range values {
			if value == nil {
				row[columnNames[i]] = gvar.New(nil)
			} else {
				// As sql.RawBytes is type of slice,
				// it should do a copy of it.
				v := make([]byte, len(value))
				copy(v, value)
				row[columnNames[i]] = gvar.New(c.DB.convertValue(v, columnTypes[i]))
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
// It is enabled when configuration "debug" is true.
func (c *Core) writeSqlToLogger(v *Sql) {
	s := fmt.Sprintf("[%3d ms] %s", v.End-v.Start, v.Format)
	if v.Error != nil {
		s += "\nError: " + v.Error.Error()
		c.logger.StackWithFilter(gPATH_FILTER_KEY).Error(s)
	} else {
		c.logger.StackWithFilter(gPATH_FILTER_KEY).Debug(s)
	}
}
