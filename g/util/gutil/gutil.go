// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gutil provides some uncategorized util functions.
// 
// 工具包.
package gutil

import (
    "bytes"
    "encoding/json"
    "fmt"
    "github.com/gogf/gf/g/internal/empty"
    "github.com/gogf/gf/g/util/gconv"
    "os"
    "reflect"
    "runtime"
)

// 格式化打印变量
func Dump(i...interface{}) {
    s := Export(i...)
    if s != "" {
        fmt.Println(s)
    }
}

// 格式化导出变量
func Export(i...interface{}) string {
    buffer := bytes.NewBuffer(nil)
    for _, v := range i {
        if b, ok := v.([]byte); ok {
            buffer.Write(b)
        } else {
            // 主要针对 map[interface{}]* 进行处理，json无法进行encode，
            // 这里强制对所有map进行反射处理转换
            refValue := reflect.ValueOf(v)
            if refValue.Kind() == reflect.Map {
                m := make(map[string]interface{})
                keys := refValue.MapKeys()
                for _, k := range keys {
                    m[gconv.String(k.Interface())] = refValue.MapIndex(k).Interface()
                }
                v = m
            }
            // JSON格式化
            encoder := json.NewEncoder(buffer)
            encoder.SetEscapeHTML(false)
            encoder.SetIndent("", "\t")
            if err := encoder.Encode(v); err != nil {
                fmt.Fprintln(os.Stderr, err.Error())
            }
        }
    }
    return buffer.String()
}

// 打印完整的调用回溯信息
func PrintBacktrace() {
    index  := 1
    buffer := bytes.NewBuffer(nil)
    for i := 0; i < 10000; i++ {
        if _, path, line, ok := runtime.Caller(i); ok {
            buffer.WriteString(fmt.Sprintf(`%d. %s:%d%s`, index, path, line, "\n"))
            index++
        } else {
            break
        }
    }
    fmt.Print(buffer.String())
}

// 抛出一个异常
func Throw(exception interface{}) {
    panic(exception)
}

// try...catch...
func TryCatch(try func(), catch ... func(exception interface{})) {
    if len(catch) > 0 {
        defer func() {
            if e := recover(); e != nil {
                catch[0](e)
            }
        }()
    }
    try()
}

// IsEmpty checks given value empty or not.
// false: integer(0), bool(false), slice/map(len=0), nil;
// true : other.
//
// 判断给定的变量是否为空。
// 整型为0, 布尔为false, slice/map长度为0, 其他为nil的情况，都为空。
// 为空时返回true，否则返回false。
func IsEmpty(value interface{}) bool {
    return empty.IsEmpty(value)
}


