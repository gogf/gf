<<<<<<< HEAD
// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
=======
// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
>>>>>>> upstream/master
//

package gdb

import (
<<<<<<< HEAD
    "fmt"
    "errors"
    "strings"
    "reflect"
    "database/sql"
    "gitee.com/johng/gf/g/util/gstr"
    "gitee.com/johng/gf/g/util/gconv"
)

// 关闭链接
func (db *Db) Close() error {
    if db.master != nil {
        if err := db.master.Close(); err == nil {
            db.master = nil
        } else {
            return err
        }
    }
    if db.slave != nil {
        if err := db.slave.Close(); err == nil {
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
=======
    "database/sql"
    "errors"
    "fmt"
    "github.com/gogf/gf/g/container/gvar"
    "github.com/gogf/gf/g/os/gcache"
    "github.com/gogf/gf/g/os/gtime"
    "github.com/gogf/gf/g/text/gregex"
    "github.com/gogf/gf/g/util/gconv"
    "reflect"
    "strings"
)

const (
    gDEFAULT_DEBUG_SQL_LENGTH = 1000 // 默认调试模式下记录的SQL条数
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
>>>>>>> upstream/master
    }
    return nil, err
}

// 执行一条sql，并返回执行情况，主要用于非查询操作
<<<<<<< HEAD
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
func (db *Db) GetAll(query string, args ...interface{}) (Result, error) {
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
            k := columns[i]
            v := make([]byte, len(col))
            copy(v, col)
            row[k] = v
        }
        //fmt.Printf("%p\n", row["typeid"])
        records = append(records, row)
    }
    return records, nil
}

// 数据库查询，获取查询结果记录，以关联数组结构返回
func (db *Db) GetOne(query string, args ...interface{}) (Record, error) {
    list, err := db.GetAll(query, args ...)
=======
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
>>>>>>> upstream/master
    if err != nil {
        return nil, err
    }
    if len(list) > 0 {
        return list[0], nil
    }
    return nil, nil
}

<<<<<<< HEAD
// 数据库查询，获取查询结果记录，自动映射数据到给定的struct对象中
func (db *Db) GetStruct(obj interface{}, query string, args ...interface{}) error {
    one, err := db.GetOne(query, args...)
    if err != nil {
        return err
    }
    return one.ToStruct(obj)
}


// 数据库查询，获取查询字段值
func (db *Db) GetValue(query string, args ...interface{}) (Value, error) {
    one, err := db.GetOne(query, args ...)
=======
// 数据库查询，查询单条记录，自动映射数据到给定的struct对象中
func (bs *dbBase) GetStruct(objPointer interface{}, query string, args ...interface{}) error {
    one, err := bs.GetOne(query, args...)
    if err != nil {
        return err
    }
    return one.ToStruct(objPointer)
}

// 数据库查询，查询多条记录，并自动转换为指定的slice对象, 如: []struct/[]*struct。
func (bs *dbBase) GetStructs(objPointerSlice interface{}, query string, args ...interface{}) error {
    all, err := bs.GetAll(query, args...)
    if err != nil {
        return err
    }
    return all.ToStructs(objPointerSlice)
}

// 将结果转换为指定的struct/*struct/[]struct/[]*struct,
// 参数应该为指针类型，否则返回失败。
// 该方法自动识别参数类型，调用Struct/Structs方法。
func (bs *dbBase) GetScan(objPointer interface{}, query string, args ...interface{}) error {
    t := reflect.TypeOf(objPointer)
    k := t.Kind()
    if k != reflect.Ptr {
        return fmt.Errorf("params should be type of pointer, but got: %v", k)
    }
    k = t.Elem().Kind()
    switch k {
        case reflect.Array:
        case reflect.Slice:
            return bs.db.GetStructs(objPointer, query, args ...)
        case reflect.Struct:
            return bs.db.GetStruct(objPointer, query, args ...)
        default:
            return fmt.Errorf("element type should be type of struct/slice, unsupported: %v", k)
    }
    return nil
}

// 数据库查询，获取查询字段值
func (bs *dbBase) GetValue(query string, args ...interface{}) (Value, error) {
    one, err := bs.GetOne(query, args ...)
>>>>>>> upstream/master
    if err != nil {
        return nil, err
    }
    for _, v := range one {
        return v, nil
    }
    return nil, nil
}

// 数据库查询，获取查询数量
<<<<<<< HEAD
func (db *Db) GetCount(query string, args ...interface{}) (int, error) {
    val, err := db.GetValue(query, args ...)
    if err != nil {
        return 0, err
    }
    return gconv.Int(val), nil
}

// 数据表查询，其中tables可以是多个联表查询语句，这种查询方式较复杂，建议使用链式操作
func (db *Db) Select(tables, fields string, condition interface{}, groupBy, orderBy string, first, limit int, args ... interface{}) (Result, error) {
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
    refValue := reflect.ValueOf(data)
    if refValue.Kind() == reflect.Map {
        var fields []string
        keys := refValue.MapKeys()
        for _, k := range keys {
            fields  = append(fields, fmt.Sprintf("%s%s%s=?", db.charl, k, db.charr))
            params  = append(params, gconv.String(refValue.MapIndex(k).Interface()))
            updates = strings.Join(fields,   ",")
        }
    } else {
        updates = gconv.String(data)
    }
    for _, v := range args {
        params = append(params, gconv.String(v))
    }
    return db.Exec(fmt.Sprintf("UPDATE %s%s%s SET %s WHERE %s", db.charl, table, db.charr, updates, db.formatCondition(condition)), params...)
}

// CURD操作:删除数据
func (db *Db) Delete(table string, condition interface{}, args ...interface{}) (sql.Result, error) {
    return db.Exec(fmt.Sprintf("DELETE FROM %s%s%s WHERE %s", db.charl, table, db.charr, db.formatCondition(condition)), args...)
}

// 格式化SQL查询条件
func (db *Db) formatCondition(condition interface{}) (where string) {
    if reflect.ValueOf(condition).Kind() == reflect.Map {
        ks := reflect.ValueOf(condition).MapKeys()
        vs := reflect.ValueOf(condition)
        for _, k := range ks {
            key   := gconv.String(k.Interface())
            value := gconv.String(vs.MapIndex(k).Interface())
            isNum := gstr.IsNumeric(value)
            if len(where) > 0 {
                where += " and "
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

=======
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

// CURD操作:单条数据写入, 仅仅执行写入操作，如果存在冲突的主键或者唯一索引，那么报错返回。
// 参数data支持map/struct/*struct/slice类型，
// 当为slice(例如[]map/[]struct/[]*struct)类型时，batch参数生效，并自动切换为批量操作。
func (bs *dbBase) Insert(table string, data interface{}, batch...int) (sql.Result, error) {
    return bs.db.doInsert(nil, table, data, OPTION_INSERT, batch...)
}

// CURD操作:单条数据写入, 如果数据存在(主键或者唯一索引)，那么删除后重新写入一条。
// 参数data支持map/struct/*struct/slice类型，
// 当为slice(例如[]map/[]struct/[]*struct)类型时，batch参数生效，并自动切换为批量操作。
func (bs *dbBase) Replace(table string, data interface{}, batch...int) (sql.Result, error) {
    return bs.db.doInsert(nil, table, data, OPTION_REPLACE, batch...)
}

// CURD操作:单条数据写入, 如果数据存在(主键或者唯一索引)，那么更新，否则写入一条新数据。
// 参数data支持map/struct/*struct/slice类型，
// 当为slice(例如[]map/[]struct/[]*struct)类型时，batch参数生效，并自动切换为批量操作。
func (bs *dbBase) Save(table string, data interface{}, batch...int) (sql.Result, error) {
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
func (bs *dbBase) doInsert(link dbLink, table string, data interface{}, option int, batch...int) (result sql.Result, err error) {
    var fields  []string
    var values  []string
    var params  []interface{}
    var dataMap Map
    // 使用反射判断data数据类型，如果为slice类型，那么自动转为批量操作
    rv   := reflect.ValueOf(data)
    kind := rv.Kind()
    if kind == reflect.Ptr {
        rv   = rv.Elem()
        kind = rv.Kind()
    }
    switch kind {
        case reflect.Slice: fallthrough
        case reflect.Array:
            return bs.db.doBatchInsert(link, table, data, option, batch...)
        case reflect.Map:   fallthrough
        case reflect.Struct:
            dataMap = structToMap(data)
        default:
            return result, errors.New(fmt.Sprint("unsupported data type:", kind))
    }
    charL, charR := bs.db.getChars()
    for k, v := range dataMap {
        fields = append(fields, charL + k + charR)
        values = append(values, "?")
        params = append(params, convertParam(v))
    }
    operation := getInsertOperationByOption(option)
    updateStr := ""
    if option == OPTION_SAVE {
        var updates []string
        for k, _ := range dataMap {
            updates = append(updates,
                fmt.Sprintf("%s%s%s=VALUES(%s%s%s)",
                    charL, k, charR,
                    charL, k, charR,
                ),
            )
        }
        updateStr = fmt.Sprintf("ON DUPLICATE KEY UPDATE %s", strings.Join(updates, ","))
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
func (bs *dbBase) BatchInsert(table string, list interface{}, batch...int) (sql.Result, error) {
    return bs.db.doBatchInsert(nil, table, list, OPTION_INSERT, batch...)
}

// CURD操作:批量数据指定批次量写入, 如果数据存在(主键或者唯一索引)，那么删除后重新写入一条
func (bs *dbBase) BatchReplace(table string, list interface{}, batch...int) (sql.Result, error) {
    return bs.db.doBatchInsert(nil, table, list, OPTION_REPLACE, batch...)
}

// CURD操作:批量数据指定批次量写入, 如果数据存在(主键或者唯一索引)，那么更新，否则写入一条新数据
func (bs *dbBase) BatchSave(table string, list interface{}, batch...int) (sql.Result, error) {
    return bs.db.doBatchInsert(nil, table, list, OPTION_SAVE, batch...)
}

// 批量写入数据, 参数list支持slice类型，例如: []map/[]struct/[]*struct。
func (bs *dbBase) doBatchInsert(link dbLink, table string, list interface{}, option int, batch...int) (result sql.Result, err error) {
    var keys   []string
    var values []string
    var params []interface{}
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
            rv   := reflect.ValueOf(list)
            kind := rv.Kind()
            if kind == reflect.Ptr {
                rv   = rv.Elem()
                kind = rv.Kind()
            }
            switch kind {
                // 如果是slice，那么转换为List类型
                case reflect.Slice: fallthrough
                case reflect.Array:
                    listMap = make(List, rv.Len())
                    for i := 0; i < rv.Len(); i++ {
                        listMap[i] = structToMap(rv.Index(i).Interface())
                    }
                case reflect.Map:   fallthrough
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
        keys    = append(keys,   k)
        holders = append(holders, "?")
    }
    batchResult    := new(batchSqlResult)
    charL, charR   := bs.db.getChars()
    keyStr         := charL + strings.Join(keys, charL + "," + charR) + charR
    valueHolderStr := "(" + strings.Join(holders, ",") + ")"
    // 操作判断
    operation := getInsertOperationByOption(option)
    updateStr := ""
    if option == OPTION_SAVE {
        var updates []string
        for _, k := range keys {
            updates = append(updates,
                fmt.Sprintf("%s%s%s=VALUES(%s%s%s)",
                    charL, k, charR,
                    charL, k, charR,
                ),
            )
        }
        updateStr = fmt.Sprintf(" ON DUPLICATE KEY UPDATE %s", strings.Join(updates, ","))
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
            if n, err := r.RowsAffected(); err != nil  {
                return r, err
            } else {
                batchResult.lastResult    = r
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
        if n, err := r.RowsAffected(); err != nil  {
            return r, err
        } else {
            batchResult.lastResult    = r
            batchResult.rowsAffected += n
        }
    }
    return batchResult, nil
}

// CURD操作:数据更新，统一采用sql预处理。
// data参数支持string/map/struct/*struct类型。
func (bs *dbBase) Update(table string, data interface{}, condition interface{}, args ...interface{}) (sql.Result, error) {
    newWhere, newArgs := formatCondition(condition, args)
    return bs.db.doUpdate(nil, table, data, newWhere, newArgs ...)
}

// CURD操作:数据更新，统一采用sql预处理。
// data参数支持string/map/struct/*struct类型类型。
func (bs *dbBase) doUpdate(link dbLink, table string, data interface{}, condition string, args ...interface{}) (result sql.Result, err error) {
    updates      := ""
    charL, charR := bs.db.getChars()
    // 使用反射进行类型判断
    rv   := reflect.ValueOf(data)
    kind := rv.Kind()
    if kind == reflect.Ptr {
        rv   = rv.Elem()
        kind = rv.Kind()
    }
    params := []interface{}(nil)
    switch kind {
        case reflect.Map:   fallthrough
        case reflect.Struct:
            var fields []string
            for k, v := range structToMap(data) {
                fields = append(fields, fmt.Sprintf("%s%s%s=?", charL, k, charR))
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
    if len(condition) == 0 {
        return bs.db.doExec(link, fmt.Sprintf("UPDATE %s SET %s", table, updates), args...)
    }
    return bs.db.doExec(link, fmt.Sprintf("UPDATE %s SET %s WHERE %s", table, updates, condition), args...)
}

// CURD操作:删除数据
func (bs *dbBase) Delete(table string, condition interface{}, args ...interface{}) (result sql.Result, err error) {
    newWhere, newArgs := formatCondition(condition, args)
    return bs.db.doDelete(nil, table, newWhere, newArgs ...)
}

// CURD操作:删除数据
func (bs *dbBase) doDelete(link dbLink, table string, condition string, args ...interface{}) (result sql.Result, err error) {
    if link == nil {
        if link, err = bs.db.Master(); err != nil {
            return nil, err
        }
    }
    if len(condition) == 0 {
        return bs.db.doExec(link, fmt.Sprintf("DELETE FROM %s", table), args...)
    }
    return bs.db.doExec(link, fmt.Sprintf("DELETE FROM %s WHERE %s", table, condition), args...)
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
>>>>>>> upstream/master
