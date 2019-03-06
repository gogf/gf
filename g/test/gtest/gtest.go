// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gtest provides simple and useful test utils.
// 
// 测试模块.
package gtest

import (
    "fmt"
    "github.com/gogf/gf/g/util/gconv"
    "os"
    "reflect"
    "regexp"
    "runtime"
    "testing"
)

// 封装一个测试用例
func Case(t *testing.T, f func()) {
    defer func() {
        if err := recover(); err != nil {
            fmt.Fprintf(os.Stderr, "%v\n%s", err, getBacktrace())
            t.Fail()
        }
    }()
    f()
}

// 断言判断, 相等
func Assert(value, expect interface{}) {
    rvValue  := reflect.ValueOf(value)
    rvExpect := reflect.ValueOf(expect)
    if rvValue.Kind() == reflect.Ptr {
        if rvValue.IsNil() {
            value = nil
        }
    }
    if rvExpect.Kind() == reflect.Map {
        if err := compareMap(value, expect); err != nil {
            panic(err)
        }
        return
    }
    if fmt.Sprintf("%v", value) != fmt.Sprintf("%v", expect) {
        panic(fmt.Sprintf(`[ASSERT] EXPECT %v == %v`, value, expect))
    }
}

// 断言判断, 相等, 包括数据类型
func AssertEQ(value, expect interface{}) {
    // 类型判断
    t1 := reflect.TypeOf(value)
    t2 := reflect.TypeOf(expect)
    if t1 != t2 {
        panic(fmt.Sprintf(`[ASSERT] EXPECT TYPE %v == %v`, t1, t2))
    }
    rvValue  := reflect.ValueOf(value)
    rvExpect := reflect.ValueOf(expect)
    if rvValue.Kind() == reflect.Ptr {
        if rvValue.IsNil() {
            value = nil
        }
    }
    if rvExpect.Kind() == reflect.Map {
        if err := compareMap(value, expect); err != nil {
            panic(err)
        }
        return
    }
    if fmt.Sprintf("%v", value) != fmt.Sprintf("%v", expect) {
        panic(fmt.Sprintf(`[ASSERT] EXPECT %v == %v`, value, expect))
    }
}

// 断言判断, 不相等
func AssertNE(value, expect interface{}) {
    rvValue  := reflect.ValueOf(value)
    rvExpect := reflect.ValueOf(expect)
    if rvValue.Kind() == reflect.Ptr {
        if rvValue.IsNil() {
            value = nil
        }
    }
    if rvExpect.Kind() == reflect.Map {
        if err := compareMap(value, expect); err == nil {
            panic(fmt.Sprintf(`[ASSERT] EXPECT %v != %v`, value, expect))
        }
        return
    }
    if fmt.Sprintf("%v", value) == fmt.Sprintf("%v", expect) {
        panic(fmt.Sprintf(`[ASSERT] EXPECT %v != %v`, value, expect))
    }
}

// 断言判断, value > expect; 注意: 仅有字符串、整形、浮点型才可以比较
func AssertGT(value, expect interface{}) {
    passed := false
    switch reflect.ValueOf(expect).Kind() {
        case reflect.String:
            passed = gconv.String(value) > gconv.String(expect)

        case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
            passed = gconv.Int(value) > gconv.Int(expect)

        case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
            passed = gconv.Uint(value) > gconv.Uint(expect)

        case reflect.Float32, reflect.Float64:
            passed = gconv.Float64(value) > gconv.Float64(expect)
    }
    if !passed {
        panic(fmt.Sprintf(`[ASSERT] EXPECT %v > %v`, value, expect))
    }
}

// 断言判断, value >= expect; 注意: 仅有字符串、整形、浮点型才可以比较
func AssertGTE(value, expect interface{}) {
    passed := false
    switch reflect.ValueOf(expect).Kind() {
        case reflect.String:
            passed = gconv.String(value) >= gconv.String(expect)

        case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
            passed = gconv.Int(value) >= gconv.Int(expect)

        case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
            passed = gconv.Uint(value) >= gconv.Uint(expect)

        case reflect.Float32, reflect.Float64:
            passed = gconv.Float64(value) >= gconv.Float64(expect)
    }
    if !passed {
        panic(fmt.Sprintf(`[ASSERT] EXPECT %v >= %v`, value, expect))
    }
}

// 断言判断, value < expect; 注意: 仅有字符串、整形、浮点型才可以比较
func AssertLT(value, expect interface{}) {
    passed := false
    switch reflect.ValueOf(expect).Kind() {
        case reflect.String:
            passed = gconv.String(value) < gconv.String(expect)

        case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
            passed = gconv.Int(value) < gconv.Int(expect)

        case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
            passed = gconv.Uint(value) < gconv.Uint(expect)

        case reflect.Float32, reflect.Float64:
            passed = gconv.Float64(value) < gconv.Float64(expect)
    }
    if !passed {
        panic(fmt.Sprintf(`[ASSERT] EXPECT %v < %v`, value, expect))
    }
}

// 断言判断, value <= expect; 注意: 仅有字符串、整形、浮点型才可以比较
func AssertLTE(value, expect interface{}) {
    passed := false
    switch reflect.ValueOf(expect).Kind() {
        case reflect.String:
            passed = gconv.String(value) <= gconv.String(expect)

        case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
            passed = gconv.Int(value) <= gconv.Int(expect)

        case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
            passed = gconv.Uint(value) <= gconv.Uint(expect)

        case reflect.Float32, reflect.Float64:
            passed = gconv.Float64(value) <= gconv.Float64(expect)
    }
    if !passed {
        panic(fmt.Sprintf(`[ASSERT] EXPECT %v <= %v`, value, expect))
    }
}


// 断言判断, value IN expect; 注意: expect必须为slice类型。
// 注意：value参数可以为普通变量，也可以为slice类型。
func AssertIN(value, expect interface{}) {
    passed := true
    switch reflect.ValueOf(expect).Kind() {
        case reflect.Slice, reflect.Array:
            for _, v1 := range gconv.Interfaces(value) {
                result := false
                for _, v2 := range gconv.Interfaces(expect) {
                    if v1 == v2 {
                        result = true
                        break
                    }
                }
                if !result {
                    passed = false
                    break
                }
            }
    }
    if !passed {
        panic(fmt.Sprintf(`[ASSERT] EXPECT %v IN %v`, value, expect))
    }
}

// 断言判断, value NOT IN expect; 注意: expect必须为slice类型
func AssertNI(value, expect interface{}) {
    passed := true
    switch reflect.ValueOf(expect).Kind() {
        case reflect.Slice, reflect.Array:
            for _, v1 := range gconv.Interfaces(value) {
                result := true
                for _, v2 := range gconv.Interfaces(expect) {
                    if v1 == v2 {
                        result = false
                        break
                    }
                }
                if !result {
                    passed = false
                    break
                }
            }
    }
    if !passed {
        panic(fmt.Sprintf(`[ASSERT] EXPECT %v NOT IN %v`, value, expect))
    }
}

// 提示错误不退出进程执行
func Error(message...interface{}) {
    fmt.Fprintf(os.Stderr, "[ERROR] %s\n%s", fmt.Sprint(message...), getBacktrace())
}

// 提示错误并退出进程执行
func Fatal(message...interface{}) {
    fmt.Fprintf(os.Stderr, "[FATAL] %s\n%s", fmt.Sprint(message...), getBacktrace())
    os.Exit(1)
}

// Map比较，如果相等返回nil，否则返回错误信息.
func compareMap(value, expect interface{}) error {
    rvValue  := reflect.ValueOf(value)
    rvExpect := reflect.ValueOf(expect)
    if rvValue.Kind() == reflect.Ptr {
        if rvValue.IsNil() {
            value = nil
        }
    }
    if rvExpect.Kind() == reflect.Map {
        if rvValue.Kind() == reflect.Map {
            if rvExpect.Len() == rvValue.Len() {
                // 将两个map类型转换为同一个map类型, 才能执行比较,
                // 直接使用 rvValue.MapIndex(key).Interface() 当key类型不一致时会报错。
                mValue   := make(map[string]string)
                mExpect  := make(map[string]string)
                ksValue  := rvValue.MapKeys()
                ksExpect := rvExpect.MapKeys()
                for _, key := range ksValue {
                    mValue[gconv.String(key.Interface())] = gconv.String(rvValue.MapIndex(key).Interface())
                }
                for _, key := range ksExpect {
                    mExpect[gconv.String(key.Interface())] = gconv.String(rvExpect.MapIndex(key).Interface())
                }
                for k, v := range mExpect {
                    if v != mValue[k] {
                        return fmt.Errorf(`[ASSERT] EXPECT VALUE map["%v"]:%v == %v`, k, mValue[k], v)
                    }
                }
            } else {
                return fmt.Errorf(`[ASSERT] EXPECT MAP LENGTH %d == %d`, rvExpect.Len(), rvValue.Len())
            }
        } else {
            return fmt.Errorf(`[ASSERT] EXPECT VALUE TO BE A MAP`)
        }
    }
    return nil
}

// 获取文件调用回溯字符串，参数skip表示调用端往上多少级开始回溯
func getBacktrace(skip...int) string {
    customSkip := 0
    if len(skip) > 0 {
        customSkip = skip[0]
    }
    backtrace := ""
    index     := 1
    from      := 0
    // 首先定位业务文件开始位置
    for i := 0; i < 10; i++ {
        if _, file, _, ok := runtime.Caller(i); ok {
            if reg, _  := regexp.Compile(`gtest\.go$`); !reg.MatchString(file) {
                from = i
                break
            }
        }
    }
    // 从业务文件开始位置根据自定义的skip开始backtrace
    goRoot := runtime.GOROOT()
    for i := from + customSkip; i < 10000; i++ {
        if _, file, cline, ok := runtime.Caller(i); ok && file != "" {
            if reg, _  := regexp.Compile(`<autogenerated>`); reg.MatchString(file) {
                continue
            }
            if reg, _  := regexp.Compile(`gtest\.go$`); reg.MatchString(file) {
                continue
            }
            if goRoot != "" {
                if reg, _  := regexp.Compile("^" + goRoot); reg.MatchString(file) {
                    continue
                }
            }
            backtrace += fmt.Sprintf(`%d. %s:%d%s`, index, file, cline, "\n")
            index++
        } else {
            break
        }
    }
    return backtrace
}