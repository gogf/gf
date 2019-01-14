// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
//

package gdb

import (
    "database/sql"
    "errors"
    "fmt"
    "gitee.com/johng/gf/g/container/gvar"
    "gitee.com/johng/gf/g/os/gcache"
    "gitee.com/johng/gf/g/os/gtime"
    "gitee.com/johng/gf/g/util/gconv"
    "gitee.com/johng/gf/g/util/gregex"
    "reflect"
    "strings"
)

const (
    gDEFAULT_DEBUG_SQL_LENGTH = 1000 // 默认调试模式下记录的SQL条数
)

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
        fmt.Println(len(sqls) - k, ":")
        fmt.Println("    Sql  :", v.Sql)
        fmt.Println("    Args :", v.Args)
        fmt.Println("    Error:", v.Error)
        fmt.Println("    Start:", gtime.NewFromTimeStamp(v.Start).Format("Y-m-d H:i:s.u"))
        fmt.Println("    End  :", gtime.NewFromTimeStamp(v.End).Format("Y-m-d H:i:s.u"))
        fmt.Println("    Cost :", v.End - v.Start, "ms")
    }
}

// 数据库sql查询操作，主要执行查询
func (bs *dbBase) Query(query string, args ...interface{}) (rows *sql.Rows, err error) {
    link, err := bs.db.Slave()
    if err != nil {
        return nil,err
    }
    return bs.db.doQuery(link, query, args...)
}

// 数据库sql查询操作，主要执行查询
func (bs *dbBase) doQuery(link dbLink, query string, args ...interface{}) (rows *sql.Rows, err error) {
    query = bs.db.handleSqlBeforeExec(query)
    if bs.db.getDebug() {
        mTime1    := gtime.Millisecond()
        rows, err  = link.Query(query, args...)
        mTime2    := gtime.Millisecond()
        s         := &Sql {
            Sql   : query,
            Args  : args,
            Error : err,
            Start : mTime1,
            End   : mTime2,
        }
        bs.sqls.Put(s)
        printSql(s)
    } else {
        rows, err = link.Query(query, args ...)
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
        return nil,err
    }
    return bs.db.doExec(link, query, args...)
}

// 执行一条sql，并返回执行情况，主要用于非查询操作
func (bs *dbBase) doExec(link dbLink, query string, args ...interface{}) (result sql.Result, err error) {
    query = bs.db.handleSqlBeforeExec(query)
    if bs.db.getDebug() {
        mTime1     := gtime.Millisecond()
        result, err = link.Exec(query, args ...)
        mTime2     := gtime.Millisecond()
        s := &Sql{
            Sql   : query,
            Args  : args,
            Error : err,
            Start : mTime1,
            End   : mTime2,
        }
        bs.sqls.Put(s)
        printSql(s)
    } else {
        result, err = link.Exec(query, args ...)
    }
    return result, formatError(err, query, args...)
}

// SQL预处理，执行完成后调用返回值sql.Stmt.Exec完成sql操作; 默认执行在Slave上, 通过第二个参数指定执行在Master上
func (bs *dbBase) Prepare(query string, execOnMaster...bool) (*sql.Stmt, error) {
    err   := (error)(nil)
    link  := (dbLink)(nil)
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
    rows, err := bs.Query(query, args ...)
    if err != nil || rows == nil {
        return nil, err
    }
    defer rows.Close()
    return bs.db.rowsToResult(rows)
}

// 数据库查询，获取查询结果记录，以关联数组结构返回
func (bs *dbBase) GetOne(query string, args ...interface{}) (Record, error) {
    list, err := bs.GetAll(query, args ...)
    if err != nil {
        return nil, err
    }
    if len(list) > 0 {
        return list[0], nil
    }
    return nil, nil
}

// 数据库查询，获取查询结果记录，自动映射数据到给定的struct对象中
func (bs *dbBase) GetStruct(obj interface{}, query string, args ...interface{}) error {
    one, err := bs.GetOne(query, args...)
    if err != nil {
        return err
    }
    return one.ToStruct(obj)
}

// 数据库查询，获取查询字段值
func (bs *dbBase) GetValue(query string, args ...interface{}) (Value, error) {
    one, err := bs.GetOne(query, args ...)
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
    value, err := bs.GetValue(query, args ...)
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
            return &TX {
                db     : bs.db,
                tx     : tx,
                master : master,
            }, nil
        } else {
            return nil, err
        }
    }
}

// CURD操作:单条数据写入, 仅仅执行写入操作，如果存在冲突的主键或者唯一索引，那么报错返回
func (bs *dbBase) Insert(table string, data Map) (sql.Result, error) {
    return bs.db.doInsert(nil, table, data, OPTION_INSERT)
}

// CURD操作:单条数据写入, 如果数据存在(主键或者唯一索引)，那么删除后重新写入一条
func (bs *dbBase) Replace(table string, data Map) (sql.Result, error) {
    return bs.db.doInsert(nil, table, data, OPTION_REPLACE)
}

// CURD操作:单条数据写入, 如果数据存在(主键或者唯一索引)，那么更新，否则写入一条新数据
func (bs *dbBase) Save(table string, data Map) (sql.Result, error) {
    return bs.db.doInsert(nil, table, data, OPTION_SAVE)
}

// insert、replace, save， ignore操作
// 0: insert:  仅仅执行写入操作，如果存在冲突的主键或者唯一索引，那么报错返回
// 1: replace: 如果数据存在(主键或者唯一索引)，那么删除后重新写入一条
// 2: save:    如果数据存在(主键或者唯一索引)，那么更新，否则写入一条新数据
// 3: ignore:  如果数据存在(主键或者唯一索引)，那么什么也不做
func (bs *dbBase) doInsert(link dbLink, table string, data Map, option int) (result sql.Result, err error) {
    var fields []string
    var values []string
    var params []interface{}
    charl, charr := bs.db.getChars()
    for k, v := range data {
        fields = append(fields,   charl + k + charr)
        values = append(values, "?")
        params = append(params, v)
    }
    operation := getInsertOperationByOption(option)
    updatestr := ""
    if option == OPTION_SAVE {
        var updates []string
        for k, _ := range data {
            updates = append(updates,
                fmt.Sprintf("%s%s%s=VALUES(%s%s%s)",
                    charl, k, charr,
                    charl, k, charr,
                ),
            )
        }
        updatestr = fmt.Sprintf("ON DUPLICATE KEY UPDATE %s", strings.Join(updates, ","))
    }
    if link == nil {
        if link, err = bs.db.Master(); err != nil {
            return nil, err
        }
    }
    return bs.db.doExec(link, fmt.Sprintf("%s INTO %s(%s) VALUES(%s) %s",
        operation, table, strings.Join(fields, ","),
        strings.Join(values, ","), updatestr),
        params...)
}

// CURD操作:批量数据指定批次量写入
func (bs *dbBase) BatchInsert(table string, list List, batch int) (sql.Result, error) {
    return bs.db.doBatchInsert(nil, table, list, batch, OPTION_INSERT)
}

// CURD操作:批量数据指定批次量写入, 如果数据存在(主键或者唯一索引)，那么删除后重新写入一条
func (bs *dbBase) BatchReplace(table string, list List, batch int) (sql.Result, error) {
    return bs.db.doBatchInsert(nil, table, list, batch, OPTION_REPLACE)
}

// CURD操作:批量数据指定批次量写入, 如果数据存在(主键或者唯一索引)，那么更新，否则写入一条新数据
func (bs *dbBase) BatchSave(table string, list List, batch int) (sql.Result, error) {
    return bs.db.doBatchInsert(nil, table, list, batch, OPTION_SAVE)
}

// 批量写入数据
func (bs *dbBase) doBatchInsert(link dbLink, table string, list List, batch int, option int) (result sql.Result, err error) {
    var keys    []string
    var values  []string
    var bvalues []string
    var params  []interface{}
    // 判断长度
    if len(list) < 1 {
        return result, errors.New("empty data list")
    }
    if link == nil {
        if link, err = bs.db.Master(); err != nil {
            return
        }
    }
    // 首先获取字段名称及记录长度
    for k, _ := range list[0] {
        keys   = append(keys,   k)
        values = append(values, "?")
    }
    charl, charr   := bs.db.getChars()
    keyStr         := charl + strings.Join(keys, charl + "," + charr) + charr
    valueHolderStr := "(" + strings.Join(values, ",") + ")"
    // 操作判断
    operation := getInsertOperationByOption(option)
    updatestr := ""
    if option == OPTION_SAVE {
        var updates []string
        for _, k := range keys {
            updates = append(updates,
                fmt.Sprintf("%s%s%s=VALUES(%s%s%s)",
                    charl, k, charr,
                    charl, k, charr,
                ),
            )
        }
        updatestr = fmt.Sprintf(" ON DUPLICATE KEY UPDATE %s", strings.Join(updates, ","))
    }
    // 构造批量写入数据格式(注意map的遍历是无序的)
    for i := 0; i < len(list); i++ {
        for _, k := range keys {
            params = append(params, list[i][k])
        }
        bvalues = append(bvalues, valueHolderStr)
        if len(bvalues) == batch {
            r, err := bs.db.doExec(link, fmt.Sprintf("%s INTO %s(%s) VALUES%s %s",
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
        r, err := bs.db.doExec(link, fmt.Sprintf("%s INTO %s(%s) VALUES%s %s",
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

// CURD操作:数据更新，统一采用sql预处理
// data参数支持字符串或者关联数组类型，内部会自行做判断处理
func (bs *dbBase) Update(table string, data interface{}, condition interface{}, args ...interface{}) (sql.Result, error) {
    link, err := bs.db.Master()
    if err != nil {
        return nil, err
    }
    return bs.db.doUpdate(link, table, data, condition, args ...)
}

// CURD操作:数据更新，统一采用sql预处理
// data参数支持字符串或者关联数组类型，内部会自行做判断处理
func (bs *dbBase) doUpdate(link dbLink, table string, data interface{}, condition interface{}, args ...interface{}) (result sql.Result, err error) {
    params       := ([]interface{})(nil)
    updates      := ""
    charl, charr := bs.db.getChars()
    refValue     := reflect.ValueOf(data)
    if refValue.Kind() == reflect.Map {
        var fields []string
        keys := refValue.MapKeys()
        for _, k := range keys {
            fields = append(fields, fmt.Sprintf("%s%s%s=?", charl, k, charr))
            params = append(params, gconv.String(refValue.MapIndex(k).Interface()))
        }
        updates = strings.Join(fields, ",")
    } else {
        updates = gconv.String(data)
    }
    for _, v := range args {
        params = append(params, gconv.String(v))
    }
    if link == nil {
        if link, err = bs.db.Master(); err != nil {
            return nil, err
        }
    }
    newWhere, newArgs := formatCondition(condition, params)
    return bs.db.doExec(link, fmt.Sprintf("UPDATE %s SET %s WHERE %s", table, updates, newWhere), newArgs...)
}

// CURD操作:删除数据
func (bs *dbBase) Delete(table string, condition interface{}, args ...interface{}) (result sql.Result, err error) {
    link, err := bs.db.Master()
    if err != nil {
        return nil, err
    }
    return bs.db.doDelete(link, table, condition, args ...)
}

// CURD操作:删除数据
func (bs *dbBase) doDelete(link dbLink, table string, condition interface{}, args ...interface{}) (result sql.Result, err error) {
    newWhere, newArgs := formatCondition(condition, args)
    return bs.db.doExec(link, fmt.Sprintf("DELETE FROM %s WHERE %s", table, newWhere), newArgs...)
}

// 获得缓存对象
func (bs *dbBase) getCache() *gcache.Cache {
    return bs.cache
}

// 将数据查询的列表数据*sql.Rows转换为Result类型
func (bs *dbBase) rowsToResult(rows *sql.Rows) (Result, error) {
    // 列信息列表, 名称与类型
    types          := make([]string, 0)
    columns        := make([]string, 0)
    columnTypes, _ := rows.ColumnTypes()
    for _, t := range columnTypes {
        types   = append(types, t.DatabaseTypeName())
        columns = append(columns, t.Name())
    }
    // 返回结构组装
    values   := make([]sql.RawBytes, len(columns))
    scanArgs := make([]interface{}, len(values))
    records  := make(Result, 0)
    for i := range values {
        scanArgs[i] = &values[i]
    }
    for rows.Next() {
        if err := rows.Scan(scanArgs...); err != nil {
            return records, err
        }
        row := make(Record)
        // 注意col字段是一个[]byte类型(slice类型本身是一个指针)，多个记录循环时该变量指向的是同一个内存地址
        for i, col := range values {
            if col == nil {
                row[columns[i]] = gvar.New(nil, true)
            } else {
                // 由于 sql.RawBytes 是slice类型, 这里必须使用值复制
                v := make([]byte, len(col))
                copy(v, col)
                row[columns[i]] = gvar.New(bs.db.convertValue(v, types[i]), true)
            }
        }
        records = append(records, row)
    }
    return records, nil
}
