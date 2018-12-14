// Copyright 2017-2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gdb

import (
    "database/sql"
    "errors"
    "fmt"
    "gitee.com/johng/gf/g/container/gvar"
    "gitee.com/johng/gf/g/os/glog"
    "gitee.com/johng/gf/g/os/gtime"
    _ "gitee.com/johng/gf/third/github.com/go-sql-driver/mysql"
    "strings"
)

type dbLink interface {
    Query(query string, args ...interface{}) (*sql.Rows, error)
    Exec(sql string, args ...interface{}) (sql.Result, error)
    Prepare(sql string) (*sql.Stmt, error)
}

// 数据库sql查询操作，主要执行查询
func doQuery(db DB, link dbLink, query string, args ...interface{}) (rows *sql.Rows, err error) {
    query = db.handleSqlBeforeExec(query)
    if db.getDebug() {
        mTime1    := gtime.Millisecond()
        rows, err  = link.Query(query, args...)
        mTime2    := gtime.Millisecond()
        s         := &Sql{
            Sql   : query,
            Args  : args,
            Error : err,
            Start : mTime1,
            End   : mTime2,
        }
        db.putSql(s)
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
func doExec(db DB, link dbLink, query string, args ...interface{}) (result sql.Result, err error) {
    query = db.handleSqlBeforeExec(query)
    if db.getDebug() {
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
        db.putSql(s)
        printSql(s)
    } else {
        result, err = link.Exec(query, args ...)
    }
    return result, formatError(err, query, args...)
}

// 将数据查询的列表数据*sql.Rows转换为Result类型
func rowsToResult(rows *sql.Rows) (Result, error) {
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
            row[columns[i]] = gvar.New(v, false)
        }
        records = append(records, row)
    }
    return records, nil
}

func formatInsertQuery(db DB, table string, data Map, option uint8) (string, []interface{}) {
    var fields []string
    var values []string
    var params []interface{}
    charl, charr := db.getChars()
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
    return fmt.Sprintf("%s INTO %s(%s) VALUES(%s) %s",
            operation, table, strings.Join(fields, ","),
            strings.Join(values, ","), updatestr),
            params
}

// 打印SQL对象(仅在debug=true时有效)
func printSql(v *Sql) {
    s := fmt.Sprintf("%s, %v, %s, %s, %d ms, %s", v.Sql, v.Args,
        gtime.NewFromTimeStamp(v.Start).Format("Y-m-d H:i:s.u"),
        gtime.NewFromTimeStamp(v.End).Format("Y-m-d H:i:s.u"),
        v.End - v.Start,
        v.Func,
    )
    if v.Error != nil {
        s += "\nError: " + v.Error.Error()
        glog.Backtrace(true, 2).Error(s)
    } else {
        glog.Debug(s)
    }
}

// 格式化错误信息
func formatError(err error, query string, args ...interface{}) error {
    if err != nil {
        errstr := fmt.Sprintf("DB ERROR: %s\n", err.Error())
        errstr += fmt.Sprintf("DB QUERY: %s\n", query)
        if len(args) > 0 {
            errstr += fmt.Sprintf("DB PARAM: %v\n", args)
        }
        err = errors.New(errstr)
    }
    return err
}

// 根据insert选项获得操作名称
func getInsertOperationByOption(option uint8) string {
    oper := "INSERT"
    switch option {
        case OPTION_REPLACE:
            oper = "REPLACE"
        case OPTION_SAVE:
        case OPTION_IGNORE:
            oper = "INSERT IGNORE"
    }
    return oper
}
