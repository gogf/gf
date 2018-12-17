// Copyright 2017-2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gdb

import (
    "bytes"
    "database/sql"
    "errors"
    "fmt"
    "gitee.com/johng/gf/g/container/gvar"
    "gitee.com/johng/gf/g/os/glog"
    "gitee.com/johng/gf/g/os/gtime"
    "gitee.com/johng/gf/g/util/gconv"
    "gitee.com/johng/gf/g/util/gregex"
    "gitee.com/johng/gf/g/util/gstr"
    _ "gitee.com/johng/gf/third/github.com/go-sql-driver/mysql"
    "reflect"
    "strings"
)

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
            if col == nil {
                row[columns[i]] = gvar.New(nil, false)
            } else {
                v := make([]byte, len(col))
                copy(v, col)
                row[columns[i]] = gvar.New(v, false)
            }
        }
        records = append(records, row)
    }
    return records, nil
}

// 格式化SQL查询条件
func formatCondition(where interface{}, args []interface{}) (string, []interface{}) {
    // 条件字符串处理
    buffer := bytes.NewBuffer(nil)
    if reflect.ValueOf(where).Kind() == reflect.Map {
        ks := reflect.ValueOf(where).MapKeys()
        vs := reflect.ValueOf(where)
        for _, k := range ks {
            key   := gconv.String(k.Interface())
            value := gconv.String(vs.MapIndex(k).Interface())
            if buffer.Len() > 0 {
                buffer.WriteString(" AND ")
            }
            if gstr.IsNumeric(value) || value == "?" {
                buffer.WriteString(key + "=" + value)
            } else {
                buffer.WriteString(key + "='" + value + "'")
            }
        }
    } else {
        buffer.Write(gconv.Bytes(where))
    }
    if buffer.Len() == 0 {
        buffer.WriteString("1")
    }
    // 查询条件处理
    newWhere := buffer.String()
    newArgs  := make([]interface{}, 0)
    if len(args) > 0 {
        for index, arg := range args {
            rv   := reflect.ValueOf(arg)
            kind := rv.Kind()
            if kind == reflect.Ptr {
                rv   = rv.Elem()
                kind = rv.Kind()
            }
            switch kind {
                case reflect.Slice: fallthrough
                case reflect.Array:
                    for i := 0; i < rv.Len(); i++ {
                        newArgs = append(newArgs, rv.Index(i).Interface())
                    }
                    counter    := 0
                    newWhere, _ = gregex.ReplaceStringFunc(`\?`, newWhere, func(s string) string {
                        counter++
                        if counter == index + 1 {
                            return "?" + strings.Repeat(",?", rv.Len() - 1)
                        }
                        return s
                    })
                default:
                    newArgs = append(newArgs, arg)
            }
        }
    }
    return newWhere, newArgs
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
func getInsertOperationByOption(option int) string {
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
