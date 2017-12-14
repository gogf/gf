package gdb

import (
    "fmt"
    "errors"
    "strings"
    "database/sql"
    "gitee.com/johng/gf/g/os/glog"
)

// 关闭链接
func (l *dbLink) Close() error {
    if l.master != nil {
        err := l.master.Close()
        if (err == nil) {
            l.master = nil
        } else {
            glog.Fatal(err)
            return err
        }
    }
    if l.slave != nil {
        err := l.slave.Close()
        if (err == nil) {
            l.slave = nil
        } else {
            glog.Fatal(err)
            return err
        }
    }
    return nil
}

// 数据库sql查询操作，主要执行查询
func (l *dbLink) Query(q string, args ...interface{}) (*sql.Rows, error) {
    p         := l.link.handleSqlBeforeExec(&q)
    rows, err := l.slave.Query(*p, args ...)
    err        = l.formatError(err, p, args...)
    if (err == nil) {
        return rows, nil
    }
    return nil, err
}

// 执行一条sql，并返回执行情况，主要用于非查询操作
func (l *dbLink) Exec(q string, args ...interface{}) (sql.Result, error) {
    //fmt.Println(q)
    //fmt.Println(args)
    p      := l.link.handleSqlBeforeExec(&q)
    r, err := l.master.Exec(*p, args ...)
    err     = l.formatError(err, p, args...)
    return r, err
}

// 格式化错误信息
func (l *dbLink) formatError(err error, q *string, args ...interface{}) error {
    if err != nil {
        errstr := fmt.Sprintf("DB ERROR: %s\n", err.Error())
        errstr += fmt.Sprintf("DB QUERY: %s\n", *q)
        if len(args) > 0 {
            errstr += fmt.Sprintf("DB PARAM: %v\n", args)
        }
        err     = errors.New(errstr)
    }
    return err
}


// 数据库查询，获取查询结果集，以列表结构返回
func (l *dbLink) GetAll(q string, args ...interface{}) (*List, error) {
    // 执行sql
    rows, err := l.Query(q, args ...)
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
            return &list, err
        }
        row := make(Map)
        for i, col := range values {
            row[columns[i]] = string(col)
        }
        list = append(list, row)
    }
    return &list, nil
}

// 数据库查询，获取查询结果集，以关联数组结构返回
func (l *dbLink) GetOne(q string, args ...interface{}) (*Map, error) {
    list, err := l.GetAll(q, args ...)
    if err != nil {
        return nil, err
    }
    return &(*list)[0], nil
}

// 数据库查询，获取查询字段值
func (l *dbLink) GetValue(q string, args ...interface{}) (interface{}, error) {
    one, err := l.GetOne(q, args ...)
    if err != nil {
        return "", err
    }
    for _, v := range *one {
        return v, nil
    }
    return "", nil
}

// sql预处理，执行完成后调用返回值sql.Stmt.Exec完成sql操作
// 记得调用sql.Stmt.Close关闭操作对象
func (l *dbLink) Prepare(q string) (*sql.Stmt, error) {
    return l.master.Prepare(q)
}

// ping一下，判断或保持数据库链接(master)
func (l *dbLink) PingMaster() error {
    err := l.master.Ping();
    return err
}

// ping一下，判断或保持数据库链接(slave)
func (l *dbLink) PingSlave() error {
    err := l.slave.Ping();
    return err
}

// 设置数据库连接池中空闲链接的大小
func (l *dbLink) SetMaxIdleConns(n int) {
    l.master.SetMaxIdleConns(n);
}

// 设置数据库连接池最大打开的链接数量
func (l *dbLink) SetMaxOpenConns(n int) {
    l.master.SetMaxOpenConns(n);
}

// 事务操作，开启，会返回一个底层的事务操作对象链接如需要嵌套事务，那么可以使用该对象，否则请忽略
func (l *dbLink) Begin() (*sql.Tx, error) {
    tx, err := l.master.Begin()
    if err == nil {
        l.transaction = tx
    }
    return tx, err
}

// 事务操作，提交
func (l *dbLink) Commit() error {
    if l.transaction == nil {
        return errors.New("transaction not start")
    }
    err := l.transaction.Commit()
    return err
}

// 事务操作，回滚
func (l *dbLink) Rollback() error {
    if l.transaction == nil {
        return errors.New("transaction not start")
    }
    err := l.transaction.Rollback()
    return err
}

// 根据insert选项获得操作名称
func (l *dbLink) getInsertOperationByOption(option uint8) string {
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
func (l *dbLink) insert(table string, data *Map, option uint8) (sql.Result, error) {
    var keys   []string
    var values []string
    var params []interface{}
    for k, v := range *data {
        keys   = append(keys,   l.charl + k + l.charr)
        values = append(values, "?")
        params = append(params, v)
    }
    operation := l.getInsertOperationByOption(option)
    updatestr := ""
    if option == OPTION_SAVE {
        var updates []string
        for k, _ := range *data {
            updates = append(updates, fmt.Sprintf("%s%s%s=VALUES(%s)", l.charl, k, l.charr, k))
        }
        updatestr = fmt.Sprintf(" ON DUPLICATE KEY UPDATE %s", strings.Join(updates, ","))
    }
    return l.Exec(
        fmt.Sprintf("%s INTO %s%s%s(%s) VALUES(%s) %s",
            operation, l.charl, table, l.charr, strings.Join(keys, ","), strings.Join(values, ","), updatestr), params...
    )
}

// CURD操作:单条数据写入, 仅仅执行写入操作，如果存在冲突的主键或者唯一索引，那么报错返回
func (l *dbLink) Insert(table string, data *Map) (sql.Result, error) {
    return l.link.insert(table, data, OPTION_INSERT)
}

// CURD操作:单条数据写入, 如果数据存在(主键或者唯一索引)，那么删除后重新写入一条
func (l *dbLink) Replace(table string, data *Map) (sql.Result, error) {
    return l.link.insert(table, data, OPTION_REPLACE)
}

// CURD操作:单条数据写入, 如果数据存在(主键或者唯一索引)，那么更新，否则写入一条新数据
func (l *dbLink) Save(table string, data *Map) (sql.Result, error) {
    return l.link.insert(table, data, OPTION_SAVE)
}

// 批量写入数据
func (l *dbLink) batchInsert(table string, list *List, batch int, option uint8) error {
    var keys    []string
    var values  []string
    var bvalues []string
    var params  []interface{}
    var size int = len(*list)
    // 判断长度
    if size < 1 {
        return errors.New("empty data list")
    }
    // 首先获取字段名称及记录长度
    for k, _ := range (*list)[0] {
        keys   = append(keys,   k)
        values = append(values, "?")
    }
    var kstr = l.charl + strings.Join(keys, l.charl + "," + l.charr) + l.charr
    // 操作判断
    operation := l.getInsertOperationByOption(option)
    updatestr := ""
    if option == OPTION_SAVE {
        var updates []string
        for _, k := range keys {
            updates = append(updates, fmt.Sprintf("%s=VALUES(%s)", l.charl, k, l.charr, k))
        }
        updatestr = fmt.Sprintf(" ON DUPLICATE KEY UPDATE %s", strings.Join(updates, ","))
    }
    // 构造批量写入数据格式(注意map的遍历是无序的)
    for i := 0; i < size; i++ {
        for _, k := range keys {
            params = append(params, (*list)[i][k])
        }
        bvalues = append(bvalues, "(" + strings.Join(values, ",") + ")")
        if len(bvalues) == batch {
            _, err := l.Exec(fmt.Sprintf("%s INTO %s%s%s(%s) VALUES%s %s", operation, l.charl, table, l.charr, kstr, strings.Join(bvalues, ","), updatestr), params...)
            if err != nil {
                return err
            }
            bvalues = bvalues[:0]
        }
    }
    // 处理最后不构成指定批量的数据
    if (len(bvalues) > 0) {
        _, err := l.Exec(fmt.Sprintf("%s INTO %s%s%s(%s) VALUES%s %s", operation, l.charl, table, l.charr, kstr, strings.Join(bvalues, ","), updatestr), params...)
        if err != nil {
            return err
        }
    }
    return nil
}

// CURD操作:批量数据指定批次量写入
func (l *dbLink) BatchInsert(table string, list *List, batch int) error {
    return l.link.batchInsert(table, list, batch, OPTION_INSERT)
}

// CURD操作:批量数据指定批次量写入, 如果数据存在(主键或者唯一索引)，那么删除后重新写入一条
func (l *dbLink) BatchReplace(table string, list *List, batch int) error {
    return l.link.batchInsert(table, list, batch, OPTION_REPLACE)
}

// CURD操作:批量数据指定批次量写入, 如果数据存在(主键或者唯一索引)，那么更新，否则写入一条新数据
func (l *dbLink) BatchSave(table string, list *List, batch int) error {
    return l.link.batchInsert(table, list, batch, OPTION_SAVE)
}

// CURD操作:数据更新，统一采用sql预处理
// data参数支持字符串或者关联数组类型，内部会自行做判断处理
func (l *dbLink) Update(table string, data interface{}, condition interface{}, args ...interface{}) (sql.Result, error) {
    var params  []interface{}
    var updates string
    switch data.(type) {
        case string:
            updates = data.(string)
        case *Map:
            var keys []string
            for k, v := range *data.(*Map) {
                keys   = append(keys,   fmt.Sprintf("%s%s%s=?", l.charl, k, l.charr))
                params = append(params, v)
            }
            updates = strings.Join(keys,   ",")

        default:
            return nil, errors.New("invalid data type for 'data' field, string or *Map expected")
    }
    for _, v := range args {
        if r, ok := v.(string); ok {
            params = append(params, r)
        } else if r, ok := v.(int); ok {
            params = append(params, string(r))
        } else {

        }
    }
    return l.Exec(fmt.Sprintf("UPDATE %s%s%s SET %s WHERE %s", l.charl, table, l.charr, updates, condition), params...)
}

// CURD操作:删除数据
func (l *dbLink) Delete(table string, condition interface{}, args ...interface{}) (sql.Result, error) {
    return l.Exec(fmt.Sprintf("DELETE FROM %s WHERE %s", l.charl, table, l.charr, condition), args...)
}

