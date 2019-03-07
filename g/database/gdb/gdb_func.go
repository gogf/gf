// Copyright 2017-2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
    "bytes"
    "errors"
    "fmt"
    "github.com/gogf/gf/g/os/glog"
    "github.com/gogf/gf/g/os/gtime"
    "github.com/gogf/gf/g/text/gregex"
    "github.com/gogf/gf/g/util/gconv"
    _ "github.com/gogf/gf/third/github.com/go-sql-driver/mysql"
    "reflect"
    "strings"
)

// 格式化SQL查询条件
func formatCondition(where interface{}, args []interface{}) (newWhere string, newArgs []interface{}) {
    // 条件字符串处理
    buffer := bytes.NewBuffer(nil)
    // 使用反射进行类型判断
    rv   := reflect.ValueOf(where)
    kind := rv.Kind()
    if kind == reflect.Ptr {
        rv   = rv.Elem()
        kind = rv.Kind()
    }
    switch kind {
        // 注意当where为map/struct类型时，args参数必须为空。
        case reflect.Map:   fallthrough
        case reflect.Struct:
            for k, v := range gconv.Map(where) {
                if buffer.Len() > 0 {
                    buffer.WriteString(" AND ")
                }
                buffer.WriteString(k + "=?")
                newArgs = append(newArgs, v)
            }
            newWhere = buffer.String()
        default:
            buffer.WriteString(gconv.String(where))
    }
    if buffer.Len() == 0 {
        buffer.WriteString("1=1")
    }
    // 查询条件参数处理，主要处理slice参数类型
    newWhere = buffer.String()
    if len(args) > 0 {
        for index, arg := range args {
            rv   := reflect.ValueOf(arg)
            kind := rv.Kind()
            if kind == reflect.Ptr {
                rv   = rv.Elem()
                kind = rv.Kind()
            }
            switch kind {
                // '?'占位符支持slice类型,
                // 这里会将slice参数拆散，并更新原有占位符'?'为多个'?'，使用','符号连接。
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
    return
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
    operator := "INSERT"
    switch option {
        case OPTION_REPLACE:
            operator = "REPLACE"
        case OPTION_SAVE:
        case OPTION_IGNORE:
            operator = "INSERT IGNORE"
    }
    return operator
}
