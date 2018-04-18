// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
//

package gdb

import (
    "fmt"
    "errors"
    "strings"
    "database/sql"
    "gitee.com/johng/gf/g/util/gconv"
)

// 关闭链接
func (db *Db) Close() error {
    if db.master != nil {
        err := db.master.Close()
        if err == nil {
            db.master = nil
        } else {
            return err
        }
    }
    if db.slave != nil {
        err := db.slave.Close()
        if err == nil {
            db.slave = nil
        } else {
            return err
        }
    }
    return nil
}

// 数据库sql查询操作，主要执行查询
func (db *Db) Query(query string, args ...interface{}) (*sql.Rows, error) {
    p         := db.link.handleSqlBeforeExec(&query)
    rows, err := db.slave.Query(*p, args ...)
    err        = db.formatError(err, p, args...)
    if err == nil {
        return rows, nil
    }
    return nil, err
}

// 执行一条sql，并返回执行情况，主要用于非查询操作
func (db *Db) Exec(query string, args ...interface{}) (sql.Result, error) {
    p      := db.link.handleSqlBeforeExec(&query)
    r, err := db.master.Exec(*p, args ...)
    err     = db.formatError(err, p, args...)
    return r, err
}

// 格式化错误信息
func (db *Db) formatError(err error, query *string, args ...interface{}) error {
    if err != nil {
        errstr := fmt.Sprintf("DB ERROR: %s\n", err.Error())
        errstr += fmt.Sprintf("DB QUERY: %s\n", *query)
        if len(args) > 0 {
            errstr += fmt.Sprintf("DB PARAM: %v\n", args)
        }
        err = errors.New(errstr)
    }
    return err
}


// 数据库查询，获取查询结果集，以列表结构返回
func (db *Db) GetAll(query string, args ...interface{}) (List, error) {
    // 执行sql
    rows, err := db.Query(query, args ...)
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
    var list List
    for i := range values {
        scanArgs[i] = &values[i]
    }
    for rows.Next() {
        err = rows.Scan(scanArgs...)
        if err != nil {
            return list, err
        }
        row := make(Map)
        for i, col := range values {
            row[columns[i]] = string(col)
        }
        list = append(list, row)
    }
    return list, nil
}

// 数据库查询，获取查询结果集，以关联数组结构返回
func (db *Db) GetOne(query string, args ...interface{}) (Map, error) {
    list, err := db.GetAll(query, args ...)
    if err != nil {
        return nil, err
    }
    if len(list) > 0 {
        return list[0], nil
    }
    return nil, nil
}

// 数据库查询，获取查询字段值
func (db *Db) GetValue(query string, args ...interface{}) (interface{}, error) {
    one, err := db.GetOne(query, args ...)
    if err != nil {
        return nil, err
    }
    for _, v := range one {
        return v, nil
    }
    return nil, nil
}

// sql预处理，执行完成后调用返回值sql.Stmt.Exec完成sql操作
// 记得调用sql.Stmt.Close关闭操作对象
func (db *Db) Prepare(query string) (*sql.Stmt, error) {
    return db.master.Prepare(query)
}

// ping一下，判断或保持数据库链接(master)
func (db *Db) PingMaster() error {
    err := db.master.Ping();
    return err
}

// ping一下，判断或保持数据库链接(slave)
func (db *Db) PingSlave() error {
    err := db.slave.Ping();
    return err
}

// 设置数据库连接池中空闲链接的大小
func (db *Db) SetMaxIdleConns(n int) {
    db.master.SetMaxIdleConns(n);
}

// 设置数据库连接池最大打开的链接数量
func (db *Db) SetMaxOpenConns(n int) {
    db.master.SetMaxOpenConns(n);
}

// 事务操作，开启，会返回一个底层的事务操作对象链接如需要嵌套事务，那么可以使用该对象，否则请忽略
func (db *Db) Begin() (*Tx, error) {
    if tx, err := db.master.Begin(); err == nil {
        return &Tx {
            db : db,
            tx : tx,
        }, nil
    } else {
        return nil, err
    }
}

// 根据insert选项获得操作名称
func (db *Db) getInsertOperationByOption(option uint8) string {
    oper := "INSERT"
    switch option {
        case OPTION_INSERT:
        case OPTION_REPLACE:
            oper = "REPLACE"
        case OPTION_SAVE:
        case OPTION_IGNORE:
            oper = "INSERT IGNORE"
    }
    return oper
}

// insert、replace, save， ignore操作
// 0: insert:  仅仅执行写入操作，如果存在冲突的主键或者唯一索引，那么报错返回
// 1: replace: 如果数据存在(主键或者唯一索引)，那么删除后重新写入一条
// 2: save:    如果数据存在(主键或者唯一索引)，那么更新，否则写入一条新数据
// 3: ignore:  如果数据存在(主键或者唯一索引)，那么什么也不做
func (db *Db) insert(table string, data Map, option uint8) (sql.Result, error) {
    var keys   []string
    var values []string
    var params []interface{}
    for k, v := range data {
        keys   = append(keys,   db.charl + k + db.charr)
        values = append(values, "?")
        params = append(params, v)
    }
    operation := db.getInsertOperationByOption(option)
    updatestr := ""
    if option == OPTION_SAVE {
        var updates []string
        for k, _ := range data {
            updates = append(updates, fmt.Sprintf("%s%s%s=VALUES(%s)", db.charl, k, db.charr, k))
        }
        updatestr = fmt.Sprintf(" ON DUPLICATE KEY UPDATE %s", strings.Join(updates, ","))
    }
    return db.Exec(
        fmt.Sprintf("%s INTO %s%s%s(%s) VALUES(%s) %s",
            operation, db.charl, table, db.charr, strings.Join(keys, ","), strings.Join(values, ","), updatestr), params...
    )
}

// CURD操作:单条数据写入, 仅仅执行写入操作，如果存在冲突的主键或者唯一索引，那么报错返回
func (db *Db) Insert(table string, data Map) (sql.Result, error) {
    return db.insert(table, data, OPTION_INSERT)
}

// CURD操作:单条数据写入, 如果数据存在(主键或者唯一索引)，那么删除后重新写入一条
func (db *Db) Replace(table string, data Map) (sql.Result, error) {
    return db.insert(table, data, OPTION_REPLACE)
}

// CURD操作:单条数据写入, 如果数据存在(主键或者唯一索引)，那么更新，否则写入一条新数据
func (db *Db) Save(table string, data Map) (sql.Result, error) {
    return db.insert(table, data, OPTION_SAVE)
}

// 批量写入数据
func (db *Db) batchInsert(table string, list List, batch int, option uint8) (sql.Result, error) {
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
    var kstr = db.charl + strings.Join(keys, db.charl + "," + db.charr) + db.charr
    // 操作判断
    operation := db.getInsertOperationByOption(option)
    updatestr := ""
    if option == OPTION_SAVE {
        var updates []string
        for _, k := range keys {
            updates = append(updates, fmt.Sprintf("%s=VALUES(%s)", db.charl, k, db.charr, k))
        }
        updatestr = fmt.Sprintf(" ON DUPLICATE KEY UPDATE %s", strings.Join(updates, ","))
    }
    // 构造批量写入数据格式(注意map的遍历是无序的)
    for i := 0; i < size; i++ {
        for _, k := range keys {
            params = append(params, list[i][k])
        }
        bvalues = append(bvalues, "(" + strings.Join(values, ",") + ")")
        if len(bvalues) == batch {
            r, err := db.Exec(fmt.Sprintf("%s INTO %s%s%s(%s) VALUES%s %s", operation, db.charl, table, db.charr, kstr, strings.Join(bvalues, ","), updatestr), params...)
            if err != nil {
                return result, err
            }
            result  = r
            bvalues = bvalues[:0]
        }
    }
    // 处理最后不构成指定批量的数据
    if len(bvalues) > 0 {
        r, err := db.Exec(fmt.Sprintf("%s INTO %s%s%s(%s) VALUES%s %s", operation, db.charl, table, db.charr, kstr, strings.Join(bvalues, ","), updatestr), params...)
        if err != nil {
            return result, err
        }
        result = r
    }
    return result, nil
}

// CURD操作:批量数据指定批次量写入
func (db *Db) BatchInsert(table string, list List, batch int) (sql.Result, error) {
    return db.batchInsert(table, list, batch, OPTION_INSERT)
}

// CURD操作:批量数据指定批次量写入, 如果数据存在(主键或者唯一索引)，那么删除后重新写入一条
func (db *Db) BatchReplace(table string, list List, batch int) (sql.Result, error) {
    return db.batchInsert(table, list, batch, OPTION_REPLACE)
}

// CURD操作:批量数据指定批次量写入, 如果数据存在(主键或者唯一索引)，那么更新，否则写入一条新数据
func (db *Db) BatchSave(table string, list List, batch int) (sql.Result, error) {
    return db.batchInsert(table, list, batch, OPTION_SAVE)
}

// CURD操作:数据更新，统一采用sql预处理
// data参数支持字符串或者关联数组类型，内部会自行做判断处理
func (db *Db) Update(table string, data interface{}, condition interface{}, args ...interface{}) (sql.Result, error) {
    var params  []interface{}
    var updates string
    switch value := data.(type) {
        case string:
            updates = value
        case Map:
            var keys []string
            for k, v := range value {
                keys   = append(keys,   fmt.Sprintf("%s%s%s=?", db.charl, k, db.charr))
                params = append(params, v)
            }
            updates = strings.Join(keys,   ",")

        default:
            return nil, errors.New("invalid data type for 'data' field, string or Map expected")
    }
    for _, v := range args {
        params = append(params, gconv.String(v))
    }
    return db.Exec(fmt.Sprintf("UPDATE %s%s%s SET %s WHERE %s", db.charl, table, db.charr, updates, condition), params...)
}

// CURD操作:删除数据
func (db *Db) Delete(table string, condition interface{}, args ...interface{}) (sql.Result, error) {
    return db.Exec(fmt.Sprintf("DELETE FROM %s WHERE %s", db.charl, table, db.charr, condition), args...)
}

