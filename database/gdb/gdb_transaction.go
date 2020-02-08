// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/gogf/gf/text/gregex"
)

// 数据库事务对象
type TX struct {
	db     DB
	tx     *sql.Tx
	master *sql.DB
}

// 事务操作，提交
func (tx *TX) Commit() error {
	return tx.tx.Commit()
}

// 事务操作，回滚
func (tx *TX) Rollback() error {
	return tx.tx.Rollback()
}

// (事务)数据库sql查询操作，主要执行查询
func (tx *TX) Query(query string, args ...interface{}) (rows *sql.Rows, err error) {
	return tx.db.doQuery(tx.tx, query, args...)
}

// (事务)执行一条sql，并返回执行情况，主要用于非查询操作
func (tx *TX) Exec(query string, args ...interface{}) (sql.Result, error) {
	return tx.db.doExec(tx.tx, query, args...)
}

// sql预处理，执行完成后调用返回值sql.Stmt.Exec完成sql操作
func (tx *TX) Prepare(query string) (*sql.Stmt, error) {
	return tx.db.doPrepare(tx.tx, query)
}

// 数据库查询，获取查询结果集，以列表结构返回
func (tx *TX) GetAll(query string, args ...interface{}) (Result, error) {
	rows, err := tx.Query(query, args...)
	if err != nil || rows == nil {
		return nil, err
	}
	defer rows.Close()
	return tx.db.rowsToResult(rows)
}

// 数据库查询，获取查询结果记录，以关联数组结构返回
func (tx *TX) GetOne(query string, args ...interface{}) (Record, error) {
	list, err := tx.GetAll(query, args...)
	if err != nil {
		return nil, err
	}
	if len(list) > 0 {
		return list[0], nil
	}
	return nil, nil
}

// 数据库查询，获取查询结果记录，自动映射数据到给定的struct对象中
func (tx *TX) GetStruct(obj interface{}, query string, args ...interface{}) error {
	one, err := tx.GetOne(query, args...)
	if err != nil {
		return err
	}
	return one.Struct(obj)
}

// 数据库查询，查询多条记录，并自动转换为指定的slice对象, 如: []struct/[]*struct。
func (tx *TX) GetStructs(objPointerSlice interface{}, query string, args ...interface{}) error {
	all, err := tx.GetAll(query, args...)
	if err != nil {
		return err
	}
	return all.Structs(objPointerSlice)
}

// 将结果转换为指定的struct/*struct/[]struct/[]*struct,
// 参数应该为指针类型，否则返回失败。
// 该方法自动识别参数类型，调用Struct/Structs方法。
func (tx *TX) GetScan(objPointer interface{}, query string, args ...interface{}) error {
	t := reflect.TypeOf(objPointer)
	k := t.Kind()
	if k != reflect.Ptr {
		return fmt.Errorf("params should be type of pointer, but got: %v", k)
	}
	k = t.Elem().Kind()
	switch k {
	case reflect.Array:
	case reflect.Slice:
		return tx.db.GetStructs(objPointer, query, args...)
	case reflect.Struct:
		return tx.db.GetStruct(objPointer, query, args...)
	default:
		return fmt.Errorf("element type should be type of struct/slice, unsupported: %v", k)
	}
	return nil
}

// 数据库查询，获取查询字段值
func (tx *TX) GetValue(query string, args ...interface{}) (Value, error) {
	one, err := tx.GetOne(query, args...)
	if err != nil {
		return nil, err
	}
	for _, v := range one {
		return v, nil
	}
	return nil, nil
}

// 数据库查询，获取查询数量
func (tx *TX) GetCount(query string, args ...interface{}) (int, error) {
	if !gregex.IsMatchString(`(?i)SELECT\s+COUNT\(.+\)\s+FROM`, query) {
		query, _ = gregex.ReplaceString(`(?i)(SELECT)\s+(.+)\s+(FROM)`, `$1 COUNT($2) $3`, query)
	}
	value, err := tx.GetValue(query, args...)
	if err != nil {
		return 0, err
	}
	return value.Int(), nil
}

// CURD操作:单条数据写入, 仅仅执行写入操作，如果存在冲突的主键或者唯一索引，那么报错返回
func (tx *TX) Insert(table string, data interface{}, batch ...int) (sql.Result, error) {
	return tx.db.doInsert(tx.tx, table, data, gINSERT_OPTION_DEFAULT, batch...)
}

func (tx *TX) InsertIgnore(table string, data interface{}, batch ...int) (sql.Result, error) {
	return tx.db.doInsert(tx.tx, table, data, gINSERT_OPTION_IGNORE, batch...)
}

// CURD操作:单条数据写入, 如果数据存在(主键或者唯一索引)，那么删除后重新写入一条
func (tx *TX) Replace(table string, data interface{}, batch ...int) (sql.Result, error) {
	return tx.db.doInsert(tx.tx, table, data, gINSERT_OPTION_REPLACE, batch...)
}

// CURD操作:单条数据写入, 如果数据存在(主键或者唯一索引)，那么更新，否则写入一条新数据
func (tx *TX) Save(table string, data interface{}, batch ...int) (sql.Result, error) {
	return tx.db.doInsert(tx.tx, table, data, gINSERT_OPTION_SAVE, batch...)
}

// CURD操作:批量数据指定批次量写入
func (tx *TX) BatchInsert(table string, list interface{}, batch ...int) (sql.Result, error) {
	return tx.db.doBatchInsert(tx.tx, table, list, gINSERT_OPTION_DEFAULT, batch...)
}

// CURD操作:批量数据指定批次量写入, 如果数据存在(主键或者唯一索引)，那么删除后重新写入一条
func (tx *TX) BatchReplace(table string, list interface{}, batch ...int) (sql.Result, error) {
	return tx.db.doBatchInsert(tx.tx, table, list, gINSERT_OPTION_REPLACE, batch...)
}

// CURD操作:批量数据指定批次量写入, 如果数据存在(主键或者唯一索引)，那么更新，否则写入一条新数据
func (tx *TX) BatchSave(table string, list interface{}, batch ...int) (sql.Result, error) {
	return tx.db.doBatchInsert(tx.tx, table, list, gINSERT_OPTION_SAVE, batch...)
}

// CURD操作:数据更新，统一采用sql预处理,
// data参数支持字符串或者关联数组类型，内部会自行做判断处理.
func (tx *TX) Update(table string, data interface{}, condition interface{}, args ...interface{}) (sql.Result, error) {
	newWhere, newArgs := formatWhere(tx.db, condition, args, false)
	if newWhere != "" {
		newWhere = " WHERE " + newWhere
	}
	return tx.doUpdate(table, data, newWhere, newArgs...)
}

// 与Update方法的区别是不处理条件参数
func (tx *TX) doUpdate(table string, data interface{}, condition string, args ...interface{}) (sql.Result, error) {
	return tx.db.doUpdate(tx.tx, table, data, condition, args...)
}

// CURD操作:删除数据
func (tx *TX) Delete(table string, condition interface{}, args ...interface{}) (sql.Result, error) {
	newWhere, newArgs := formatWhere(tx.db, condition, args, false)
	if newWhere != "" {
		newWhere = " WHERE " + newWhere
	}
	return tx.doDelete(table, newWhere, newArgs...)
}

// 与Delete方法的区别是不处理条件参数
func (tx *TX) doDelete(table string, condition string, args ...interface{}) (sql.Result, error) {
	return tx.db.doDelete(tx.tx, table, condition, args...)
}
