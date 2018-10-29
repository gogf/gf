// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gdb

import (
    "fmt"
    "errors"
    "strings"
    "reflect"
    "database/sql"
    "gitee.com/johng/gf/g/os/gtime"
    "gitee.com/johng/gf/g/util/gconv"
    _ "gitee.com/johng/gf/third/github.com/go-sql-driver/mysql"
    "gitee.com/johng/gf/g/container/gvar"
)

// 数据库事务对象
type Tx struct {
    db     *Db
    tx     *sql.Tx
    master *sql.DB
}

// 事务操作，提交
func (tx *Tx) Commit() error {
    err := tx.tx.Commit()
    tx.master.Close()
    return err
}

// 事务操作，回滚
func (tx *Tx) Rollback() error {
    err := tx.tx.Rollback()
    tx.master.Close()
    return err
}

// (事务)数据库sql查询操作，主要执行查询
func (tx *Tx) Query(query string, args ...interface{}) (*sql.Rows, error) {
    var err  error
    var rows *sql.Rows
    p := tx.db.link.handleSqlBeforeExec(&query)
    if tx.db.debug.Val() {
        militime1 := gtime.Millisecond()
        rows, err  = tx.tx.Query(*p, args ...)
        militime2 := gtime.Millisecond()
        s := &Sql{
            Sql   : *p,
            Args  : args,
            Error : err,
            Start : militime1,
            End   : militime2,
            Func  : "TX:Query",
        }
        tx.db.sqls.Put(s)
        tx.db.printSql(s)
    } else {
        rows, err  = tx.tx.Query(*p, args ...)
    }
    if err == nil {
        return rows, nil
    } else {
        err = tx.db.formatError(err, p, args...)
    }
    return nil, err
}

// (事务)执行一条sql，并返回执行情况，主要用于非查询操作
func (tx *Tx) Exec(query string, args ...interface{}) (sql.Result, error) {
    var err    error
    var result sql.Result
    p := tx.db.link.handleSqlBeforeExec(&query)
    if tx.db.debug.Val() {
        militime1  := gtime.Millisecond()
        result, err = tx.tx.Exec(*p, args ...)
        militime2  := gtime.Millisecond()
        s := &Sql{
            Sql   : *p,
            Args  : args,
            Error : err,
            Start : militime1,
            End   : militime2,
            Func  : "TX:Exec",
        }
        tx.db.sqls.Put(s)
        tx.db.printSql(s)
    } else {
        result, err = tx.tx.Exec(*p, args ...)
    }
    return result, tx.db.formatError(err, p, args...)
}

// 数据库查询，获取查询结果集，以列表结构返回
func (tx *Tx) GetAll(query string, args ...interface{}) (Result, error) {
    // 执行sql
    rows, err := tx.Query(query, args ...)
    if err != nil || rows == nil {
        return nil, err
    }
    // 列名称列表
    columns, err := rows.Columns()
    if err != nil {
        return nil, err
    }
    // 返回结构组装
    values   := make([]sql.RawBytes, len(columns))
    scanArgs := make([]interface{}, len(values))
    records  := make(Result, 0)
    for i := range values {
        scanArgs[i] = &values[i]
    }
    for rows.Next() {
        err = rows.Scan(scanArgs...)
        if err != nil {
            return records, err
        }
        row := make(Record)
        // 注意col字段是一个[]byte类型(slice类型本身是一个指针)，多个记录循环时该变量指向的是同一个内存地址
        for i, col := range values {
            v := make([]byte, len(col))
            copy(v, col)
            row[columns[i]] = gvar.New(v)
        }
        //fmt.Printf("%p\n", row["typeid"])
        records = append(records, row)
    }
    return records, nil
}

// 数据库查询，获取查询结果记录，以关联数组结构返回
func (tx *Tx) GetOne(query string, args ...interface{}) (Record, error) {
    list, err := tx.GetAll(query, args ...)
    if err != nil {
        return nil, err
    }
    if len(list) > 0 {
        return list[0], nil
    }
    return nil, nil
}

// 数据库查询，获取查询结果记录，自动映射数据到给定的struct对象中
func (tx *Tx) GetStruct(obj interface{}, query string, args ...interface{}) error {
    one, err := tx.GetOne(query, args...)
    if err != nil {
        return err
    }
    return one.ToStruct(obj)
}


// 数据库查询，获取查询字段值
func (tx *Tx) GetValue(query string, args ...interface{}) (Value, error) {
    one, err := tx.GetOne(query, args ...)
    if err != nil {
        return nil, err
    }
    for _, v := range one {
        return v, nil
    }
    return nil, nil
}

// 数据库查询，获取查询数量
func (tx *Tx) GetCount(query string, args ...interface{}) (int, error) {
    val, err := tx.GetValue(query, args ...)
    if err != nil {
        return 0, err
    }
    return gconv.Int(val), nil
}

// 数据表查询，其中tables可以是多个联表查询语句，这种查询方式较复杂，建议使用链式操作
func (tx *Tx) Select(tables, fields string, condition interface{}, groupBy, orderBy string, first, limit int, args ... interface{}) (Result, error) {
    s := fmt.Sprintf("SELECT %s FROM %s ", fields, tables)
    if condition != nil {
        s += fmt.Sprintf("WHERE %s ", tx.db.formatCondition(condition))
    }
    if len(groupBy) > 0 {
        s += fmt.Sprintf("GROUP BY %s ", groupBy)
    }
    if len(orderBy) > 0 {
        s += fmt.Sprintf("ORDER BY %s ", orderBy)
    }
    if limit > 0 {
        s += fmt.Sprintf("LIMIT %d,%d ", first, limit)
    }
    return tx.GetAll(s, args ... )
}

// sql预处理，执行完成后调用返回值sql.Stmt.Exec完成sql操作
// 记得调用sql.Stmt.Close关闭操作对象
func (tx *Tx) Prepare(query string) (*sql.Stmt, error) {
    return tx.tx.Prepare(query)
}

// insert、replace, save， ignore操作
// 0: insert:  仅仅执行写入操作，如果存在冲突的主键或者唯一索引，那么报错返回
// 1: replace: 如果数据存在(主键或者唯一索引)，那么删除后重新写入一条
// 2: save:    如果数据存在(主键或者唯一索引)，那么更新，否则写入一条新数据
// 3: ignore:  如果数据存在(主键或者唯一索引)，那么什么也不做
func (tx *Tx) insert(table string, data Map, option uint8) (sql.Result, error) {
    var keys   []string
    var values []string
    var params []interface{}
    for k, v := range data {
        keys   = append(keys,   tx.db.charl + k + tx.db.charr)
        values = append(values, "?")
        params = append(params, v)
    }
    operation := tx.db.getInsertOperationByOption(option)
    updatestr := ""
    if option == OPTION_SAVE {
        var updates []string
        for k, _ := range data {
            updates = append(updates, fmt.Sprintf("%s%s%s=VALUES(%s)", tx.db.charl, k, tx.db.charr, k))
        }
        updatestr = fmt.Sprintf(" ON DUPLICATE KEY UPDATE %s", strings.Join(updates, ","))
    }
    return tx.Exec(
        fmt.Sprintf("%s INTO %s(%s) VALUES(%s) %s",
            operation, table, strings.Join(keys, ","),
            strings.Join(values, ","),
            updatestr),
            params...
    )
}

// CURD操作:单条数据写入, 仅仅执行写入操作，如果存在冲突的主键或者唯一索引，那么报错返回
func (tx *Tx) Insert(table string, data Map) (sql.Result, error) {
    return tx.insert(table, data, OPTION_INSERT)
}

// CURD操作:单条数据写入, 如果数据存在(主键或者唯一索引)，那么删除后重新写入一条
func (tx *Tx) Replace(table string, data Map) (sql.Result, error) {
    return tx.insert(table, data, OPTION_REPLACE)
}

// CURD操作:单条数据写入, 如果数据存在(主键或者唯一索引)，那么更新，否则写入一条新数据
func (tx *Tx) Save(table string, data Map) (sql.Result, error) {
    return tx.insert(table, data, OPTION_SAVE)
}

// 批量写入数据
func (tx *Tx) batchInsert(table string, list List, batch int, option uint8) (sql.Result, error) {
    var keys    []string
    var values  []string
    var bvalues []string
    var params  []interface{}
    var result  sql.Result
    var size = len(list)
    // 判断长度
    if size < 1 {
        return result, errors.New("empty data list")
    }
    // 首先获取字段名称及记录长度
    for k, _ := range list[0] {
        keys   = append(keys,   k)
        values = append(values, "?")
    }
    keyStr         := tx.db.charl + strings.Join(keys, tx.db.charl + "," + tx.db.charr) + tx.db.charr
    valueHolderStr := "(" + strings.Join(values, ",") + ")"
    // 操作判断
    operation := tx.db.getInsertOperationByOption(option)
    updatestr := ""
    if option == OPTION_SAVE {
        var updates []string
        for _, k := range keys {
            updates = append(updates, fmt.Sprintf("%s%s%s=VALUES(%s)", tx.db.charl, k, tx.db.charr, k))
        }
        updatestr = fmt.Sprintf(" ON DUPLICATE KEY UPDATE %s", strings.Join(updates, ","))
    }
    // 构造批量写入数据格式(注意map的遍历是无序的)
    for i := 0; i < size; i++ {
        for _, k := range keys {
            params = append(params, list[i][k])
        }
        bvalues = append(bvalues, valueHolderStr)
        if len(bvalues) == batch {
            r, err := tx.Exec(fmt.Sprintf("%s INTO %s(%s) VALUES%s %s",
                operation, table, keyStr, strings.Join(bvalues, ","),
                updatestr),
                params...)
            if err != nil {
                return result, err
            }
            result  = r
            params  = params[:0]
            bvalues = bvalues[:0]
        }
    }
    // 处理最后不构成指定批量的数据
    if len(bvalues) > 0 {
        r, err := tx.Exec(fmt.Sprintf("%s INTO %s(%s) VALUES%s %s",
            operation, table, keyStr, strings.Join(bvalues, ","),
            updatestr),
            params...)
        if err != nil {
            return result, err
        }
        result = r
    }
    return result, nil
}

// CURD操作:批量数据指定批次量写入
func (tx *Tx) BatchInsert(table string, list List, batch int) (sql.Result, error) {
    return tx.batchInsert(table, list, batch, OPTION_INSERT)
}

// CURD操作:批量数据指定批次量写入, 如果数据存在(主键或者唯一索引)，那么删除后重新写入一条
func (tx *Tx) BatchReplace(table string, list List, batch int) (sql.Result, error) {
    return tx.batchInsert(table, list, batch, OPTION_REPLACE)
}

// CURD操作:批量数据指定批次量写入, 如果数据存在(主键或者唯一索引)，那么更新，否则写入一条新数据
func (tx *Tx) BatchSave(table string, list List, batch int) (sql.Result, error) {
    return tx.batchInsert(table, list, batch, OPTION_SAVE)
}

// CURD操作:数据更新，统一采用sql预处理
// data参数支持字符串或者关联数组类型，内部会自行做判断处理
func (tx *Tx) Update(table string, data interface{}, condition interface{}, args ...interface{}) (sql.Result, error) {
    var params  []interface{}
    var updates string
    refValue := reflect.ValueOf(data)
    if refValue.Kind() == reflect.Map {
        var fields []string
        keys := refValue.MapKeys()
        for _, k := range keys {
            fields = append(fields, fmt.Sprintf("%s%s%s=?", tx.db.charl, k, tx.db.charr))
            params = append(params, gconv.String(refValue.MapIndex(k).Interface()))
            updates = strings.Join(fields,   ",")
        }
    } else {
        updates = gconv.String(data)
    }
    for _, v := range args {
        params = append(params, gconv.String(v))
    }
    return tx.Exec(fmt.Sprintf("UPDATE %s SET %s WHERE %s", table, updates, tx.db.formatCondition(condition)), params...)
}

// CURD操作:删除数据
func (tx *Tx) Delete(table string, condition interface{}, args ...interface{}) (sql.Result, error) {
    return tx.Exec(fmt.Sprintf("DELETE FROM %s WHERE %s", table, tx.db.formatCondition(condition)), args...)
}

