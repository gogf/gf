// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//

package gdb

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/gogf/gf/g/text/gstr"

	"github.com/gogf/gf/g/container/gvar"
	"github.com/gogf/gf/g/os/gcache"
	"github.com/gogf/gf/g/os/gtime"
	"github.com/gogf/gf/g/text/gregex"
	"github.com/gogf/gf/g/util/gconv"
)

const (
	// 默认调试模式下记录的SQL条数
	gDEFAULT_DEBUG_SQL_LENGTH = 1000
)

var (
	wordReg         = regexp.MustCompile(`^[a-zA-Z0-9\-_]+$`)
	lastOperatorReg = regexp.MustCompile(`[<>=]+\s*$`)
)

// 获取最近一条执行的sql
func (bs *dbBase) GetLastSql() *Sql {
	if bs.sqls == nil {
		return nil
	}
	if v := bs.sqls.Val(); v != nil {
		return v.(*Sql)
	}
	return nil
}

// 获取已经执行的SQL列表(仅在debug=true时有效)
func (bs *dbBase) GetQueriedSqls() []*Sql {
	if bs.sqls == nil {
		return nil
	}
	sqls := make([]*Sql, 0)
	bs.sqls.Prev()
	bs.sqls.RLockIteratorPrev(func(value interface{}) bool {
		if value == nil {
			return false
		}
		sqls = append(sqls, value.(*Sql))
		return true
	})
	return sqls
}

// 打印已经执行的SQL列表(仅在debug=true时有效)
func (bs *dbBase) PrintQueriedSqls() {
	sqls := bs.GetQueriedSqls()
	for k, v := range sqls {
		fmt.Println(len(sqls)-k, ":")
		fmt.Println("    Sql  :", v.Sql)
		fmt.Println("    Args :", v.Args)
		fmt.Println("    Error:", v.Error)
		fmt.Println("    Start:", gtime.NewFromTimeStamp(v.Start).Format("Y-m-d H:i:s.u"))
		fmt.Println("    End  :", gtime.NewFromTimeStamp(v.End).Format("Y-m-d H:i:s.u"))
		fmt.Println("    Cost :", v.End-v.Start, "ms")
	}
}

// 数据库sql查询操作，主要执行查询
func (bs *dbBase) Query(query string, args ...interface{}) (rows *sql.Rows, err error) {
	link, err := bs.db.Slave()
	if err != nil {
		return nil, err
	}
	return bs.db.doQuery(link, query, args...)
}

// 数据库sql查询操作，主要执行查询
func (bs *dbBase) doQuery(link dbLink, query string, args ...interface{}) (rows *sql.Rows, err error) {
	query, args = formatQuery(query, args)
	query = bs.db.handleSqlBeforeExec(query)
	if bs.db.getDebug() {
		mTime1 := gtime.Millisecond()
		rows, err = link.Query(query, args...)
		mTime2 := gtime.Millisecond()
		s := &Sql{
			Sql:   query,
			Args:  args,
			Error: err,
			Start: mTime1,
			End:   mTime2,
		}
		bs.sqls.Put(s)
		printSql(s)
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

// 执行一条sql，并返回执行情况，主要用于非查询操作
func (bs *dbBase) Exec(query string, args ...interface{}) (result sql.Result, err error) {
	link, err := bs.db.Master()
	if err != nil {
		return nil, err
	}
	return bs.db.doExec(link, query, args...)
}

// 执行一条sql，并返回执行情况，主要用于非查询操作
func (bs *dbBase) doExec(link dbLink, query string, args ...interface{}) (result sql.Result, err error) {
	query, args = formatQuery(query, args)
	query = bs.db.handleSqlBeforeExec(query)
	if bs.db.getDebug() {
		mTime1 := gtime.Millisecond()
		result, err = link.Exec(query, args...)
		mTime2 := gtime.Millisecond()
		s := &Sql{
			Sql:   query,
			Args:  args,
			Error: err,
			Start: mTime1,
			End:   mTime2,
		}
		bs.sqls.Put(s)
		printSql(s)
	} else {
		result, err = link.Exec(query, args...)
	}
	return result, formatError(err, query, args...)
}

// SQL预处理，执行完成后调用返回值sql.Stmt.Exec完成sql操作; 默认执行在Slave上, 通过第二个参数指定执行在Master上
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

// SQL预处理，执行完成后调用返回值sql.Stmt.Exec完成sql操作
func (bs *dbBase) doPrepare(link dbLink, query string) (*sql.Stmt, error) {
	return link.Prepare(query)
}

// 数据库查询，获取查询结果集，以列表结构返回
func (bs *dbBase) GetAll(query string, args ...interface{}) (Result, error) {
	rows, err := bs.Query(query, args...)
	if err != nil || rows == nil {
		return nil, err
	}
	defer rows.Close()
	return bs.db.rowsToResult(rows)
}

// 数据库查询，获取查询结果记录，以关联数组结构返回
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

// 数据库查询，查询单条记录，自动映射数据到给定的struct对象中
func (bs *dbBase) GetStruct(pointer interface{}, query string, args ...interface{}) error {
	one, err := bs.GetOne(query, args...)
	if err != nil {
		return err
	}
	return one.ToStruct(pointer)
}

// 数据库查询，查询多条记录，并自动转换为指定的slice对象, 如: []struct/[]*struct。
func (bs *dbBase) GetStructs(pointer interface{}, query string, args ...interface{}) error {
	all, err := bs.GetAll(query, args...)
	if err != nil {
		return err
	}
	return all.ToStructs(pointer)
}

// 将结果转换为指定的struct/*struct/[]struct/[]*struct,
// 参数应该为指针类型，否则返回失败。
// 该方法自动识别参数类型，调用Struct/Structs方法。
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

// 数据库查询，获取查询字段值
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

// 数据库查询，获取查询数量
func (bs *dbBase) GetCount(query string, args ...interface{}) (int, error) {
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
	return bs.db.doInsert(nil, table, data, OPTION_INSERT, batch...)
}

// CURD操作:单条数据写入, 如果数据存在(主键或者唯一索引)，那么删除后重新写入一条。
// 参数data支持map/struct/*struct/slice类型，
// 当为slice(例如[]map/[]struct/[]*struct)类型时，batch参数生效，并自动切换为批量操作。
func (bs *dbBase) Replace(table string, data interface{}, batch ...int) (sql.Result, error) {
	return bs.db.doInsert(nil, table, data, OPTION_REPLACE, batch...)
}

// CURD操作:单条数据写入, 如果数据存在(主键或者唯一索引)，那么更新，否则写入一条新数据。
// 参数data支持map/struct/*struct/slice类型，
// 当为slice(例如[]map/[]struct/[]*struct)类型时，batch参数生效，并自动切换为批量操作。
func (bs *dbBase) Save(table string, data interface{}, batch ...int) (sql.Result, error) {
	return bs.db.doInsert(nil, table, data, OPTION_SAVE, batch...)
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
	table = bs.db.quoteWord(table)
	// 使用反射判断data数据类型，如果为slice类型，那么自动转为批量操作
	rv := reflect.ValueOf(data)
	kind := rv.Kind()
	if kind == reflect.Ptr {
		rv = rv.Elem()
		kind = rv.Kind()
	}
	switch kind {
	case reflect.Slice:
		fallthrough
	case reflect.Array:
		return bs.db.doBatchInsert(link, table, data, option, batch...)
	case reflect.Map:
		fallthrough
	case reflect.Struct:
		dataMap = structToMap(data)
	default:
		return result, errors.New(fmt.Sprint("unsupported data type:", kind))
	}
	charL, charR := bs.db.getChars()
	for k, v := range dataMap {
		fields = append(fields, charL+k+charR)
		values = append(values, "?")
		params = append(params, convertParam(v))
	}
	operation := getInsertOperationByOption(option)
	updateStr := ""
	if option == OPTION_SAVE {
		for k, _ := range dataMap {
			if len(updateStr) > 0 {
				updateStr += ","
			}
			updateStr += fmt.Sprintf("%s%s%s=VALUES(%s%s%s)",
				charL, k, charR,
				charL, k, charR,
			)
		}
		updateStr = fmt.Sprintf(" ON DUPLICATE KEY UPDATE %s", updateStr)
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
	return bs.db.doBatchInsert(nil, table, list, OPTION_INSERT, batch...)
}

// CURD操作:批量数据指定批次量写入, 如果数据存在(主键或者唯一索引)，那么删除后重新写入一条
func (bs *dbBase) BatchReplace(table string, list interface{}, batch ...int) (sql.Result, error) {
	return bs.db.doBatchInsert(nil, table, list, OPTION_REPLACE, batch...)
}

// CURD操作:批量数据指定批次量写入, 如果数据存在(主键或者唯一索引)，那么更新，否则写入一条新数据
func (bs *dbBase) BatchSave(table string, list interface{}, batch ...int) (sql.Result, error) {
	return bs.db.doBatchInsert(nil, table, list, OPTION_SAVE, batch...)
}

// 批量写入数据, 参数list支持slice类型，例如: []map/[]struct/[]*struct。
func (bs *dbBase) doBatchInsert(link dbLink, table string, list interface{}, option int, batch ...int) (result sql.Result, err error) {
	var keys []string
	var values []string
	var params []interface{}
	table = bs.db.quoteWord(table)
	listMap := (List)(nil)
	switch v := list.(type) {
	case Result:
		listMap = v.ToList()
	case Record:
		listMap = List{v.ToMap()}
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
		// 如果是slice，那么转换为List类型
		case reflect.Slice:
			fallthrough
		case reflect.Array:
			listMap = make(List, rv.Len())
			for i := 0; i < rv.Len(); i++ {
				listMap[i] = structToMap(rv.Index(i).Interface())
			}
		case reflect.Map:
			fallthrough
		case reflect.Struct:
			listMap = List{Map(structToMap(list))}
		default:
			return result, errors.New(fmt.Sprint("unsupported list type:", kind))
		}
	}
	// 判断长度
	if len(listMap) < 1 {
		return result, errors.New("empty data list")
	}
	if link == nil {
		if link, err = bs.db.Master(); err != nil {
			return
		}
	}
	// 首先获取字段名称及记录长度
	holders := []string(nil)
	for k, _ := range listMap[0] {
		keys = append(keys, k)
		holders = append(holders, "?")
	}
	batchResult := new(batchSqlResult)
	charL, charR := bs.db.getChars()
	keyStr := charL + strings.Join(keys, charR+","+charL) + charR
	valueHolderStr := "(" + strings.Join(holders, ",") + ")"
	// 操作判断
	operation := getInsertOperationByOption(option)
	updateStr := ""
	if option == OPTION_SAVE {
		for _, k := range keys {
			if len(updateStr) > 0 {
				updateStr += ","
			}
			updateStr += fmt.Sprintf("%s%s%s=VALUES(%s%s%s)",
				charL, k, charR,
				charL, k, charR,
			)
		}
		updateStr = fmt.Sprintf(" ON DUPLICATE KEY UPDATE %s", updateStr)
	}
	// 构造批量写入数据格式(注意map的遍历是无序的)
	batchNum := gDEFAULT_BATCH_NUM
	if len(batch) > 0 {
		batchNum = batch[0]
	}
	for i := 0; i < len(listMap); i++ {
		for _, k := range keys {
			params = append(params, convertParam(listMap[i][k]))
		}
		values = append(values, valueHolderStr)
		if len(values) == batchNum {
			r, err := bs.db.doExec(link, fmt.Sprintf("%s INTO %s(%s) VALUES%s %s",
				operation, table, keyStr, strings.Join(values, ","),
				updateStr),
				params...)
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
	// 处理最后不构成指定批量的数据
	if len(values) > 0 {
		r, err := bs.db.doExec(link, fmt.Sprintf("%s INTO %s(%s) VALUES%s %s",
			operation, table, keyStr, strings.Join(values, ","),
			updateStr),
			params...)
		if err != nil {
			return r, err
		}
		if n, err := r.RowsAffected(); err != nil {
			return r, err
		} else {
			batchResult.lastResult = r
			batchResult.rowsAffected += n
		}
	}
	return batchResult, nil
}

// CURD操作:数据更新，统一采用sql预处理。
// data参数支持string/map/struct/*struct类型。
func (bs *dbBase) Update(table string, data interface{}, condition interface{}, args ...interface{}) (sql.Result, error) {
	newWhere, newArgs := bs.db.formatWhere(condition, args)
	if newWhere != "" {
		newWhere = " WHERE " + newWhere
	}
	return bs.db.doUpdate(nil, table, data, newWhere, newArgs...)
}

// CURD操作:数据更新，统一采用sql预处理。
// data参数支持string/map/struct/*struct类型类型。
func (bs *dbBase) doUpdate(link dbLink, table string, data interface{}, condition string, args ...interface{}) (result sql.Result, err error) {
	table = bs.db.quoteWord(table)
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
	case reflect.Map:
		fallthrough
	case reflect.Struct:
		var fields []string
		for k, v := range structToMap(data) {
			fields = append(fields, bs.db.quoteWord(k)+"=?")
			params = append(params, convertParam(v))
		}
		updates = strings.Join(fields, ",")
	default:
		updates = gconv.String(data)
	}
	if len(params) > 0 {
		args = append(params, args...)
	}
	// 如果没有传递link，那么使用默认的写库对象
	if link == nil {
		if link, err = bs.db.Master(); err != nil {
			return nil, err
		}
	}
	return bs.db.doExec(link, fmt.Sprintf("UPDATE %s SET %s%s", table, updates, condition), args...)
}

// CURD操作:删除数据
func (bs *dbBase) Delete(table string, condition interface{}, args ...interface{}) (result sql.Result, err error) {
	newWhere, newArgs := bs.db.formatWhere(condition, args)
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
	table = bs.db.quoteWord(table)
	return bs.db.doExec(link, fmt.Sprintf("DELETE FROM %s%s", table, condition), args...)
}

// 获得缓存对象
func (bs *dbBase) getCache() *gcache.Cache {
	return bs.cache
}

// 将数据查询的列表数据*sql.Rows转换为Result类型
func (bs *dbBase) rowsToResult(rows *sql.Rows) (Result, error) {
	if !rows.Next() {
		return nil, sql.ErrNoRows
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
				row[columns[i]] = gvar.New(nil, true)
			} else {
				// 由于 sql.RawBytes 是slice类型, 这里必须使用值复制
				v := make([]byte, len(column))
				copy(v, column)
				//fmt.Println(columns[i], types[i], string(v), v, bs.db.convertValue(v, types[i]))
				row[columns[i]] = gvar.New(bs.db.convertValue(v, types[i]), true)
			}
		}
		records = append(records, row)
		if !rows.Next() {
			break
		}
	}
	return records, nil
}

// 格式化Where查询条件。
func (bs *dbBase) formatWhere(where interface{}, args []interface{}) (newWhere string, newArgs []interface{}) {
	// 条件字符串处理
	buffer := bytes.NewBuffer(nil)
	// 使用反射进行类型判断
	rv := reflect.ValueOf(where)
	kind := rv.Kind()
	if kind == reflect.Ptr {
		rv = rv.Elem()
		kind = rv.Kind()
	}
	switch kind {
	// map/struct类型
	case reflect.Map:
		fallthrough
	case reflect.Struct:
		for key, value := range structToMap(where) {
			// 字段安全符号判断
			key = bs.db.quoteWord(key)
			if buffer.Len() > 0 {
				buffer.WriteString(" AND ")
			}
			// 支持slice键值/属性，如果只有一个?占位符号，那么作为IN查询，否则打散作为多个查询参数
			rv := reflect.ValueOf(value)
			switch rv.Kind() {
			case reflect.Slice:
				fallthrough
			case reflect.Array:
				count := gstr.Count(key, "?")
				if count == 0 {
					buffer.WriteString(key + " IN(?)")
					newArgs = append(newArgs, value)
				} else if count != rv.Len() {
					buffer.WriteString(key)
					newArgs = append(newArgs, value)
				} else {
					buffer.WriteString(key)
					// 如果键名/属性名称中带有多个?占位符号，那么将参数打散
					newArgs = append(newArgs, gconv.Interfaces(value)...)
				}
			default:
				if value == nil {
					buffer.WriteString(key)
				} else {
					// 支持key带操作符号
					if gstr.Pos(key, "?") == -1 {
						if gstr.Pos(key, "<") == -1 && gstr.Pos(key, ">") == -1 && gstr.Pos(key, "=") == -1 {
							buffer.WriteString(key + "=?")
						} else {
							buffer.WriteString(key + "?")
						}
					} else {
						buffer.WriteString(key)
					}
					newArgs = append(newArgs, value)
				}
			}
		}

	default:
		buffer.WriteString(gconv.String(where))
	}
	// 没有任何条件查询参数，直接返回
	if buffer.Len() == 0 {
		return "", args
	}
	newArgs = append(newArgs, args...)
	newWhere = buffer.String()
	// 查询条件参数处理，主要处理slice参数类型
	if len(newArgs) > 0 {
		// 支持例如 Where/And/Or("uid", 1) , Where/And/Or("uid>=", 1) 这种格式
		if gstr.Pos(newWhere, "?") == -1 {
			if lastOperatorReg.MatchString(newWhere) {
				newWhere += "?"
			} else if wordReg.MatchString(newWhere) {
				newWhere += "=?"
			}
		}
	}
	return
}

// 使用关键字操作符转义给定字符串。
// 如果给定的字符串不为单词，那么不转义，直接返回该字符串。
func (bs *dbBase) quoteWord(s string) string {
	charLeft, charRight := bs.db.getChars()
	if wordReg.MatchString(s) && !gstr.ContainsAny(s, charLeft+charRight) {
		return charLeft + s + charRight
	}
	return s
}

// 动态切换数据库
func (bs *dbBase) setSchema(sqlDb *sql.DB, schema string) error {
	_, err := sqlDb.Exec("USE " + schema)
	return err
}
