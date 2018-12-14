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
    "gitee.com/johng/gf/g/os/gcache"
    "gitee.com/johng/gf/g/os/gtime"
    "gitee.com/johng/gf/g/util/gconv"
    "gitee.com/johng/gf/g/util/gstr"
    "reflect"
    "strings"
)

const (
    gDEFAULT_DEBUG_SQL_LENGTH = 1000 // 默认调试模式下记录的SQL条数
)

// 获取已经执行的SQL列表(仅在debug=true时有效)
func (db *dbBase) GetQueriedSqls() []*Sql {
    if db.sqls == nil {
        return nil
    }
    sqls := make([]*Sql, 0)
    db.sqls.Prev()
    db.sqls.RLockIteratorPrev(func(value interface{}) bool {
        if value == nil {
            return false
        }
        sqls = append(sqls, value.(*Sql))
        return true
    })
    return sqls
}

// 打印已经执行的SQL列表(仅在debug=true时有效)
func (db *dbBase) PrintQueriedSqls() {
    sqls := db.GetQueriedSqls()
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
func (db *dbBase) Query(query string, args ...interface{}) (rows *sql.Rows, err error) {
    link, err := db.Slave()
    if err != nil {
        return nil,err
    }
    return doQuery(db.db, link, query, args...)
}

// 执行一条sql，并返回执行情况，主要用于非查询操作
func (db *dbBase) Exec(query string, args ...interface{}) (result sql.Result, err error) {
    link, err := db.Master()
    if err != nil {
        return nil,err
    }
    return doExec(db.db, link, query, args...)
}

// SQL预处理，执行完成后调用返回值sql.Stmt.Exec完成sql操作; 默认执行在Slave上, 通过第二个参数指定执行在Master上
func (db *dbBase) Prepare(query string, execOnMaster...bool) (*sql.Stmt, error) {
    err   := (error)(nil)
    sqldb := (*sql.DB)(nil)
    if len(execOnMaster) > 0 && execOnMaster[0] {
        if sqldb, err = db.Master(); err != nil {
            return nil, err
        }
    } else {
        if sqldb, err = db.Slave(); err != nil {
            return nil, err
        }
    }
    return sqldb.Prepare(query)
}

// 数据库查询，获取查询结果集，以列表结构返回
func (db *dbBase) GetAll(query string, args ...interface{}) (Result, error) {
    rows, err := db.Query(query, args ...)
    if err != nil || rows == nil {
        return nil, err
    }
    return rowsToResult(rows)
}

// 数据库查询，获取查询结果记录，以关联数组结构返回
func (db *dbBase) GetOne(query string, args ...interface{}) (Record, error) {
    list, err := db.GetAll(query, args ...)
    if err != nil {
        return nil, err
    }
    if len(list) > 0 {
        return list[0], nil
    }
    return nil, nil
}

// 数据库查询，获取查询结果记录，自动映射数据到给定的struct对象中
func (db *dbBase) GetStruct(obj interface{}, query string, args ...interface{}) error {
    one, err := db.GetOne(query, args...)
    if err != nil {
        return err
    }
    return one.ToStruct(obj)
}


// 数据库查询，获取查询字段值
func (db *dbBase) GetValue(query string, args ...interface{}) (Value, error) {
    one, err := db.GetOne(query, args ...)
    if err != nil {
        return nil, err
    }
    for _, v := range one {
        return v, nil
    }
    return nil, nil
}

// 数据库查询，获取查询数量
func (db *dbBase) GetCount(query string, args ...interface{}) (int, error) {
    val, err := db.GetValue(query, args ...)
    if err != nil {
        return 0, err
    }
    return gconv.Int(val), nil
}

// 数据表查询，其中tables可以是多个联表查询语句，这种查询方式较复杂，建议使用链式操作
func (db *dbBase) Select(tables, fields string, condition interface{}, groupBy, orderBy string, first, limit int, args ... interface{}) (Result, error) {
    s := fmt.Sprintf("SELECT %s FROM %s ", fields, tables)
    if condition != nil {
        s += fmt.Sprintf("WHERE %s ", db.formatCondition(condition))
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
    return db.GetAll(s, args ... )
}

// ping一下，判断或保持数据库链接(master)
func (db *dbBase) PingMaster() error {
    if master, err := db.Master(); err != nil {
        return err
    } else {
        return master.Ping()
    }
}

// ping一下，判断或保持数据库链接(slave)
func (db *dbBase) PingSlave() error {
    if slave, err := db.Slave(); err != nil {
        return err
    } else {
        return slave.Ping()
    }
}

// 事务操作，开启，会返回一个底层的事务操作对象链接如需要嵌套事务，那么可以使用该对象，否则请忽略
// 只有在tx.Commit/tx.Rollback时，链接会自动Close
func (db *dbBase) Begin() (*TX, error) {
    if master, err := db.Master(); err != nil {
        return nil, err
    } else {
        if tx, err := master.Begin(); err == nil {
            return &TX {
                db     : db.db,
                tx     : tx,
                master : master,
            }, nil
        } else {
            return nil, err
        }
    }
}

// insert、replace, save， ignore操作
// 0: insert:  仅仅执行写入操作，如果存在冲突的主键或者唯一索引，那么报错返回
// 1: replace: 如果数据存在(主键或者唯一索引)，那么删除后重新写入一条
// 2: save:    如果数据存在(主键或者唯一索引)，那么更新，否则写入一条新数据
// 3: ignore:  如果数据存在(主键或者唯一索引)，那么什么也不做
func (db *dbBase) insert(table string, data Map, option uint8) (sql.Result, error) {
    var fields []string
    var values []string
    var params []interface{}
    charl, charr := db.db.getChars()
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
    return db.Exec(
        fmt.Sprintf("%s INTO %s(%s) VALUES(%s) %s",
            operation, table, strings.Join(fields, ","),
            strings.Join(values, ","),
            updatestr),
            params...
    )
}

// CURD操作:单条数据写入, 仅仅执行写入操作，如果存在冲突的主键或者唯一索引，那么报错返回
func (db *dbBase) Insert(table string, data Map) (sql.Result, error) {
    return db.insert(table, data, OPTION_INSERT)
}

// CURD操作:单条数据写入, 如果数据存在(主键或者唯一索引)，那么删除后重新写入一条
func (db *dbBase) Replace(table string, data Map) (sql.Result, error) {
    return db.insert(table, data, OPTION_REPLACE)
}

// CURD操作:单条数据写入, 如果数据存在(主键或者唯一索引)，那么更新，否则写入一条新数据
func (db *dbBase) Save(table string, data Map) (sql.Result, error) {
    return db.insert(table, data, OPTION_SAVE)
}

// 批量写入数据
func (db *dbBase) batchInsert(table string, list List, batch int, option uint8) (sql.Result, error) {
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
    charl, charr   := db.db.getChars()
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
    for i := 0; i < size; i++ {
        for _, k := range keys {
            params = append(params, list[i][k])
        }
        bvalues = append(bvalues, valueHolderStr)
        if len(bvalues) == batch {
            r, err := db.Exec(fmt.Sprintf("%s INTO %s(%s) VALUES%s %s",
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
        r, err := db.Exec(fmt.Sprintf("%s INTO %s(%s) VALUES%s %s",
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
func (db *dbBase) BatchInsert(table string, list List, batch int) (sql.Result, error) {
    return db.batchInsert(table, list, batch, OPTION_INSERT)
}

// CURD操作:批量数据指定批次量写入, 如果数据存在(主键或者唯一索引)，那么删除后重新写入一条
func (db *dbBase) BatchReplace(table string, list List, batch int) (sql.Result, error) {
    return db.batchInsert(table, list, batch, OPTION_REPLACE)
}

// CURD操作:批量数据指定批次量写入, 如果数据存在(主键或者唯一索引)，那么更新，否则写入一条新数据
func (db *dbBase) BatchSave(table string, list List, batch int) (sql.Result, error) {
    return db.batchInsert(table, list, batch, OPTION_SAVE)
}

// CURD操作:数据更新，统一采用sql预处理
// data参数支持字符串或者关联数组类型，内部会自行做判断处理
func (db *dbBase) Update(table string, data interface{}, condition interface{}, args ...interface{}) (sql.Result, error) {
    var params  []interface{}
    var updates string
    charl, charr := db.db.getChars()
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
    return db.Exec(fmt.Sprintf("UPDATE %s SET %s WHERE %s", table, updates, db.formatCondition(condition)), params...)
}

// CURD操作:删除数据
func (db *dbBase) Delete(table string, condition interface{}, args ...interface{}) (sql.Result, error) {
    return db.Exec(fmt.Sprintf("DELETE FROM %s WHERE %s", table, db.formatCondition(condition)), args...)
}

// 格式化SQL查询条件
func (db *dbBase) formatCondition(condition interface{}) (where string) {
    if reflect.ValueOf(condition).Kind() == reflect.Map {
        ks := reflect.ValueOf(condition).MapKeys()
        vs := reflect.ValueOf(condition)
        for _, k := range ks {
            key   := gconv.String(k.Interface())
            value := gconv.String(vs.MapIndex(k).Interface())
            isNum := gstr.IsNumeric(value)
            if len(where) > 0 {
                where += " AND "
            }
            if isNum || value == "?" {
                where += key + "=" + value
            } else {
                where += key + "='" + value + "'"
            }
        }
    } else {
        where += gconv.String(condition)
    }
    return
}

// 获得缓存对象
func (db *dbBase) getCache() *gcache.Cache {
    return db.cache
}

// 记录执行的SQL
func (db *dbBase) putSql(s *Sql) {
    db.sqls.Put(s)
}