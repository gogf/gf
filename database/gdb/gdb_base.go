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
	"github.com/gogf/gf/os/gcache"
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

// Query commits one query SQL to underlying driver and returns the execution result.
// It is most commonly used for data querying.
func (bs *dbBase) Query(query string, args ...interface{}) (rows *sql.Rows, err error) {
	link, err := bs.db.Slave()
	if err != nil {
		return nil, err
	}
	return bs.db.doQuery(link, query, args...)
}

// doQuery commits the query string and its arguments to underlying driver
// through given link object and returns the execution result.
func (bs *dbBase) doQuery(link dbLink, query string, args ...interface{}) (rows *sql.Rows, err error) {
	query, args = formatQuery(query, args)
	query = bs.db.handleSqlBeforeExec(query)
	if bs.db.getDebug() {
		mTime1 := gtime.TimestampMilli()
		rows, err = link.Query(query, args...)
		mTime2 := gtime.TimestampMilli()
		s := &Sql{
			Sql:    query,
			Args:   args,
			Format: bindArgsToQuery(query, args),
			Error:  err,
			Start:  mTime1,
			End:    mTime2,
		}
		bs.printSql(s)
	} else {
		rows, err = link.Query(query, args...)
	}
	if err == nil {
		return rows, nil
	} else {
		err = formatError(err, query, args...)
	}
	return nil, err
}

// Exec commits one query SQL to underlying driver and returns the execution result.
// It is most commonly used for data inserting and updating.
func (bs *dbBase) Exec(query string, args ...interface{}) (result sql.Result, err error) {
	link, err := bs.db.Master()
	if err != nil {
		return nil, err
	}
	return bs.db.doExec(link, query, args...)
}

// doExec commits the query string and its arguments to underlying driver
// through given link object and returns the execution result.
func (bs *dbBase) doExec(link dbLink, query string, args ...interface{}) (result sql.Result, err error) {
	query, args = formatQuery(query, args)
	query = bs.db.handleSqlBeforeExec(query)
	if bs.db.getDebug() {
		mTime1 := gtime.TimestampMilli()
		result, err = link.Exec(query, args...)
		mTime2 := gtime.TimestampMilli()
		s := &Sql{
			Sql:    query,
			Args:   args,
			Format: bindArgsToQuery(query, args),
			Error:  err,
			Start:  mTime1,
			End:    mTime2,
		}
		bs.printSql(s)
	} else {
		result, err = link.Exec(query, args...)
	}
	return result, formatError(err, query, args...)
}

// Prepare creates a prepared statement for later queries or executions.
// Multiple queries or executions may be run concurrently from the
// returned statement.
// The caller must call the statement's Close method
// when the statement is no longer needed.
//
// The parameter <execOnMaster> specifies whether executing the sql on master node,
// or else it executes the sql on slave node if master-slave configured.
func (bs *dbBase) Prepare(query string, execOnMaster ...bool) (*sql.Stmt, error) {
	err := (error)(nil)
	link := (dbLink)(nil)
	if len(execOnMaster) > 0 && execOnMaster[0] {
		if link, err = bs.db.Master(); err != nil {
			return nil, err
		}
	} else {
		if link, err = bs.db.Slave(); err != nil {
			return nil, err
		}
	}
	return bs.db.doPrepare(link, query)
}

// doPrepare calls prepare function on given link object and returns the statement object.
func (bs *dbBase) doPrepare(link dbLink, query string) (*sql.Stmt, error) {
	return link.Prepare(query)
}

// GetAll queries and returns data records from database.
func (bs *dbBase) GetAll(query string, args ...interface{}) (Result, error) {
	return bs.db.doGetAll(nil, query, args...)
}

// doGetAll queries and returns data records from database.
func (bs *dbBase) doGetAll(link dbLink, query string, args ...interface{}) (result Result, err error) {
	if link == nil {
		link, err = bs.db.Slave()
		if err != nil {
			return nil, err
		}
	}
	rows, err := bs.doQuery(link, query, args...)
	if err != nil || rows == nil {
		return nil, err
	}
	defer rows.Close()
	return bs.db.rowsToResult(rows)
}

// GetOne queries and returns one record from database.
func (bs *dbBase) GetOne(query string, args ...interface{}) (Record, error) {
	list, err := bs.GetAll(query, args...)
	if err != nil {
		return nil, err
	}
	if len(list) > 0 {
		return list[0], nil
	}
	return nil, nil
}

// GetStruct queries one record from database and converts it to given struct.
// The parameter <pointer> should be a pointer to struct.
func (bs *dbBase) GetStruct(pointer interface{}, query string, args ...interface{}) error {
	one, err := bs.GetOne(query, args...)
	if err != nil {
		return err
	}
	if len(one) == 0 {
		return sql.ErrNoRows
	}
	return one.Struct(pointer)
}

// GetStructs queries records from database and converts them to given struct.
// The parameter <pointer> should be type of struct slice: []struct/[]*struct.
func (bs *dbBase) GetStructs(pointer interface{}, query string, args ...interface{}) error {
	all, err := bs.GetAll(query, args...)
	if err != nil {
		return err
	}
	if len(all) == 0 {
		return sql.ErrNoRows
	}
	return all.Structs(pointer)
}

// GetScan queries one or more records from database and converts them to given struct.
// If parameter <pointer> is type of struct pointer, it calls GetStruct internally for
// the conversion. If parameter <pointer> is type of slice, it calls GetStructs internally
// for conversion.
func (bs *dbBase) GetScan(pointer interface{}, query string, args ...interface{}) error {
	t := reflect.TypeOf(pointer)
	k := t.Kind()
	if k != reflect.Ptr {
		return fmt.Errorf("params should be type of pointer, but got: %v", k)
	}
	k = t.Elem().Kind()
	switch k {
	case reflect.Array, reflect.Slice:
		return bs.db.GetStructs(pointer, query, args...)
	case reflect.Struct:
		return bs.db.GetStruct(pointer, query, args...)
	}
	return fmt.Errorf("element type should be type of struct/slice, unsupported: %v", k)
}

// GetValue queries and returns the field value from database.
// The sql should queries only one field from database, or else it returns only one
// field of the result.
func (bs *dbBase) GetValue(query string, args ...interface{}) (Value, error) {
	one, err := bs.GetOne(query, args...)
	if err != nil {
		return nil, err
	}
	for _, v := range one {
		return v, nil
	}
	return nil, nil
}

// GetCount queries and returns the count from database.
func (bs *dbBase) GetCount(query string, args ...interface{}) (int, error) {
	// If the query fields do not contains function "COUNT",
	// it replaces the query string and adds the "COUNT" function to the fields.
	if !gregex.IsMatchString(`(?i)SELECT\s+COUNT\(.+\)\s+FROM`, query) {
		query, _ = gregex.ReplaceString(`(?i)(SELECT)\s+(.+)\s+(FROM)`, `$1 COUNT($2) $3`, query)
	}
	value, err := bs.GetValue(query, args...)
	if err != nil {
		return 0, err
	}
	return value.Int(), nil
}

// ping一下，判断或保持数据库链接(master)
func (bs *dbBase) PingMaster() error {
	if master, err := bs.db.Master(); err != nil {
		return err
	} else {
		return master.Ping()
	}
}

// ping一下，判断或保持数据库链接(slave)
func (bs *dbBase) PingSlave() error {
	if slave, err := bs.db.Slave(); err != nil {
		return err
	} else {
		return slave.Ping()
	}
}

// 事务操作，开启，会返回一个底层的事务操作对象链接如需要嵌套事务，那么可以使用该对象，否则请忽略
// 只有在tx.Commit/tx.Rollback时，链接会自动Close
func (bs *dbBase) Begin() (*TX, error) {
	if master, err := bs.db.Master(); err != nil {
		return nil, err
	} else {
		if tx, err := master.Begin(); err == nil {
			return &TX{
				db:     bs.db,
				tx:     tx,
				master: master,
			}, nil
		} else {
			return nil, err
		}
	}
}

// CURD操作:单条数据写入, 仅仅执行写入操作，如果存在冲突的主键或者唯一索引，那么报错返回。
// 参数data支持map/struct/*struct/slice类型，
// 当为slice(例如[]map/[]struct/[]*struct)类型时，batch参数生效，并自动切换为批量操作。
func (bs *dbBase) Insert(table string, data interface{}, batch ...int) (sql.Result, error) {
	return bs.db.doInsert(nil, table, data, gINSERT_OPTION_DEFAULT, batch...)
}

func (bs *dbBase) InsertIgnore(table string, data interface{}, batch ...int) (sql.Result, error) {
	return bs.db.doInsert(nil, table, data, gINSERT_OPTION_IGNORE, batch...)
}

// CURD操作:单条数据写入, 如果数据存在(主键或者唯一索引)，那么删除后重新写入一条。
// 参数data支持map/struct/*struct/slice类型，
// 当为slice(例如[]map/[]struct/[]*struct)类型时，batch参数生效，并自动切换为批量操作。
func (bs *dbBase) Replace(table string, data interface{}, batch ...int) (sql.Result, error) {
	return bs.db.doInsert(nil, table, data, gINSERT_OPTION_REPLACE, batch...)
}

// CURD操作:单条数据写入, 如果数据存在(主键或者唯一索引)，那么更新，否则写入一条新数据。
// 参数data支持map/struct/*struct/slice类型，
// 当为slice(例如[]map/[]struct/[]*struct)类型时，batch参数生效，并自动切换为批量操作。
func (bs *dbBase) Save(table string, data interface{}, batch ...int) (sql.Result, error) {
	return bs.db.doInsert(nil, table, data, gINSERT_OPTION_SAVE, batch...)
}

// 支持insert、replace, save， ignore操作。
// 0: insert:  仅仅执行写入操作，如果存在冲突的主键或者唯一索引，那么报错返回;
// 1: replace: 如果数据存在(主键或者唯一索引)，那么删除后重新写入一条;
// 2: save:    如果数据存在(主键或者唯一索引)，那么更新，否则写入一条新数据;
// 3: ignore:  如果数据存在(主键或者唯一索引)，那么什么也不做;
//
// 参数data支持map/struct/*struct/slice类型，
// 当为slice(例如[]map/[]struct/[]*struct)类型时，batch参数生效，并自动切换为批量操作。
func (bs *dbBase) doInsert(link dbLink, table string, data interface{}, option int, batch ...int) (result sql.Result, err error) {
	var fields []string
	var values []string
	var params []interface{}
	var dataMap Map
	table = bs.db.handleTableName(table)
	// 使用反射判断data数据类型，如果为slice类型，那么自动转为批量操作
	rv := reflect.ValueOf(data)
	kind := rv.Kind()
	if kind == reflect.Ptr {
		rv = rv.Elem()
		kind = rv.Kind()
	}
	switch kind {
	case reflect.Slice, reflect.Array:
		return bs.db.doBatchInsert(link, table, data, option, batch...)
	case reflect.Map, reflect.Struct:
		dataMap = varToMapDeep(data)
	default:
		return result, errors.New(fmt.Sprint("unsupported data type:", kind))
	}
	if len(dataMap) == 0 {
		return nil, errors.New("data cannot be empty")
	}
	charL, charR := bs.db.getChars()
	for k, v := range dataMap {
		fields = append(fields, charL+k+charR)
		values = append(values, "?")
		params = append(params, v)
	}
	operation := getInsertOperationByOption(option)
	updateStr := ""
	if option == gINSERT_OPTION_SAVE {
		for k, _ := range dataMap {
			if len(updateStr) > 0 {
				updateStr += ","
			}
			updateStr += fmt.Sprintf("%s%s%s=VALUES(%s%s%s)",
				charL, k, charR,
				charL, k, charR,
			)
		}
		updateStr = fmt.Sprintf("ON DUPLICATE KEY UPDATE %s", updateStr)
	}
	if link == nil {
		if link, err = bs.db.Master(); err != nil {
			return nil, err
		}
	}
	return bs.db.doExec(link, fmt.Sprintf("%s INTO %s(%s) VALUES(%s) %s",
		operation, table, strings.Join(fields, ","),
		strings.Join(values, ","), updateStr),
		params...)
}

// CURD操作:批量数据指定批次量写入
func (bs *dbBase) BatchInsert(table string, list interface{}, batch ...int) (sql.Result, error) {
	return bs.db.doBatchInsert(nil, table, list, gINSERT_OPTION_DEFAULT, batch...)
}

// CURD操作:批量数据指定批次量写入, 如果数据存在(主键或者唯一索引)，那么删除后重新写入一条
func (bs *dbBase) BatchReplace(table string, list interface{}, batch ...int) (sql.Result, error) {
	return bs.db.doBatchInsert(nil, table, list, gINSERT_OPTION_REPLACE, batch...)
}

// CURD操作:批量数据指定批次量写入, 如果数据存在(主键或者唯一索引)，那么更新，否则写入一条新数据
func (bs *dbBase) BatchSave(table string, list interface{}, batch ...int) (sql.Result, error) {
	return bs.db.doBatchInsert(nil, table, list, gINSERT_OPTION_SAVE, batch...)
}

// 批量写入数据, 参数list支持slice类型，例如: []map/[]struct/[]*struct。
func (bs *dbBase) doBatchInsert(link dbLink, table string, list interface{}, option int, batch ...int) (result sql.Result, err error) {
	var keys, values []string
	var params []interface{}
	table = bs.db.handleTableName(table)
	listMap := (List)(nil)
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
		rv := reflect.ValueOf(list)
		kind := rv.Kind()
		if kind == reflect.Ptr {
			rv = rv.Elem()
			kind = rv.Kind()
		}
		switch kind {
		// If it's slice type, it then converts it to List type.
		case reflect.Slice, reflect.Array:
			listMap = make(List, rv.Len())
			for i := 0; i < rv.Len(); i++ {
				listMap[i] = varToMapDeep(rv.Index(i).Interface())
			}
		case reflect.Map, reflect.Struct:
			listMap = List{varToMapDeep(list)}
		default:
			return result, errors.New(fmt.Sprint("unsupported list type:", kind))
		}
	}
	if len(listMap) < 1 {
		return result, errors.New("data list cannot be empty")
	}
	if link == nil {
		if link, err = bs.db.Master(); err != nil {
			return
		}
	}
	// Handle the field names and place holders.
	holders := []string(nil)
	for k, _ := range listMap[0] {
		keys = append(keys, k)
		holders = append(holders, "?")
	}
	// Prepare the result pointer.
	batchResult := new(batchSqlResult)
	charL, charR := bs.db.getChars()
	keysStr := charL + strings.Join(keys, charR+","+charL) + charR
	valueHolderStr := "(" + strings.Join(holders, ",") + ")"

	operation := getInsertOperationByOption(option)
	updateStr := ""
	if option == gINSERT_OPTION_SAVE {
		for _, k := range keys {
			if len(updateStr) > 0 {
				updateStr += ","
			}
			updateStr += fmt.Sprintf("%s%s%s=VALUES(%s%s%s)",
				charL, k, charR,
				charL, k, charR,
			)
		}
		updateStr = fmt.Sprintf("ON DUPLICATE KEY UPDATE %s", updateStr)
	}
	// 构造批量写入数据格式(注意map的遍历是无序的)
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
			r, err := bs.db.doExec(
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
				batchResult.lastResult = r
				batchResult.rowsAffected += n
			}
			params = params[:0]
			values = values[:0]
		}
	}
	return batchResult, nil
}

// CURD操作:数据更新，统一采用sql预处理。
// data参数支持string/map/struct/*struct类型。
func (bs *dbBase) Update(table string, data interface{}, condition interface{}, args ...interface{}) (sql.Result, error) {
	newWhere, newArgs := formatWhere(bs.db, condition, args, false)
	if newWhere != "" {
		newWhere = " WHERE " + newWhere
	}
	return bs.db.doUpdate(nil, table, data, newWhere, newArgs...)
}

// CURD操作:数据更新，统一采用sql预处理。
// data参数支持string/map/struct/*struct类型类型。
func (bs *dbBase) doUpdate(link dbLink, table string, data interface{}, condition string, args ...interface{}) (result sql.Result, err error) {
	table = bs.db.handleTableName(table)
	updates := ""
	// 使用反射进行类型判断
	rv := reflect.ValueOf(data)
	kind := rv.Kind()
	if kind == reflect.Ptr {
		rv = rv.Elem()
		kind = rv.Kind()
	}
	params := []interface{}(nil)
	switch kind {
	case reflect.Map, reflect.Struct:
		var fields []string
		for k, v := range varToMapDeep(data) {
			fields = append(fields, bs.db.quoteWord(k)+"=?")
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
		if link, err = bs.db.Master(); err != nil {
			return nil, err
		}
	}
	return bs.db.doExec(link, fmt.Sprintf("UPDATE %s SET %s%s", table, updates, condition), args...)
}

// CURD操作:删除数据
func (bs *dbBase) Delete(table string, condition interface{}, args ...interface{}) (result sql.Result, err error) {
	newWhere, newArgs := formatWhere(bs.db, condition, args, false)
	if newWhere != "" {
		newWhere = " WHERE " + newWhere
	}
	return bs.db.doDelete(nil, table, newWhere, newArgs...)
}

// CURD操作:删除数据
func (bs *dbBase) doDelete(link dbLink, table string, condition string, args ...interface{}) (result sql.Result, err error) {
	if link == nil {
		if link, err = bs.db.Master(); err != nil {
			return nil, err
		}
	}
	table = bs.db.handleTableName(table)
	return bs.db.doExec(link, fmt.Sprintf("DELETE FROM %s%s", table, condition), args...)
}

// 获得缓存对象
func (bs *dbBase) getCache() *gcache.Cache {
	return bs.cache
}

// 获得表名前缀
func (bs *dbBase) getPrefix() string {
	return bs.prefix
}

// 将数据查询的列表数据*sql.Rows转换为Result类型
func (bs *dbBase) rowsToResult(rows *sql.Rows) (Result, error) {
	if !rows.Next() {
		return nil, nil
	}
	// 列信息列表, 名称与类型
	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}
	types := make([]string, len(columnTypes))
	columns := make([]string, len(columnTypes))
	for k, v := range columnTypes {
		types[k] = v.DatabaseTypeName()
		columns[k] = v.Name()
	}
	// 返回结构组装
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	records := make(Result, 0)
	for i := range values {
		scanArgs[i] = &values[i]
	}
	for {
		if err := rows.Scan(scanArgs...); err != nil {
			return records, err
		}
		row := make(Record)
		// 注意col字段是一个[]byte类型(slice类型本身是一个引用类型)，
		// 多个记录循环时该变量指向的是同一个内存地址
		for i, column := range values {
			if column == nil {
				row[columns[i]] = gvar.New(nil)
			} else {
				// 由于 sql.RawBytes 是slice类型, 这里必须使用值复制
				v := make([]byte, len(column))
				copy(v, column)
				row[columns[i]] = gvar.New(bs.db.convertValue(v, types[i]))
			}
		}
		records = append(records, row)
		if !rows.Next() {
			break
		}
	}
	return records, nil
}

// handleTableName adds prefix string and quote chars for the table. It handles table string like:
// "user", "user u", "user,user_detail", "user u, user_detail ut", "user as u, user_detail as ut".
//
// Note that, this will automatically checks the table prefix whether already added, if true it does
// nothing to the table name, or else adds the prefix to the table name.
func (bs *dbBase) handleTableName(table string) string {
	charLeft, charRight := bs.db.getChars()
	prefix := bs.db.getPrefix()
	return doHandleTableName(table, prefix, charLeft, charRight)
}

// quoteWord checks given string <s> a word, if true quotes it with security chars of the database
// and returns the quoted string; or else return <s> without any change.
func (bs *dbBase) quoteWord(s string) string {
	charLeft, charRight := bs.db.getChars()
	return doQuoteWord(s, charLeft, charRight)
}

// quoteString quotes string with quote chars. Strings like:
// "user", "user u", "user,user_detail", "user u, user_detail ut", "u.id asc".
func (bs *dbBase) quoteString(s string) string {
	charLeft, charRight := bs.db.getChars()
	return doQuoteString(s, charLeft, charRight)
}

// 打印SQL对象(仅在debug=true时有效)
func (bs *dbBase) printSql(v *Sql) {
	s := fmt.Sprintf("[%d ms] %s", v.End-v.Start, v.Format)
	if v.Error != nil {
		s += "\nError: " + v.Error.Error()
		bs.logger.StackWithFilter(gPATH_FILTER_KEY).Error(s)
	} else {
		bs.logger.StackWithFilter(gPATH_FILTER_KEY).Debug(s)
	}
}
