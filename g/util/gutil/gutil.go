<<<<<<< HEAD
// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 其他工具包
package gutil

import (
    "reflect"
    "gitee.com/johng/gf/g/util/gstr"
    "gitee.com/johng/gf/g/util/gconv"
)

// 将map键值对映射到对应的struct对象属性上，需要注意：
// 1、第二个参数为struct对象指针；
// 2、struct对象的公开属性才能被映射赋值；
// 3、map中的键名可以为小写，映射转换时会自动将键名首字母转为大写做匹配映射，如果无法匹配则忽略；
func MapToStruct(m map[string]interface{}, o interface{}) error {
    for k, v := range m {
        _MapToStructSetField(o, k, v)
    }
    return nil
}
func _MapToStructSetField(obj interface{}, name string, value interface{}) {
    structValue      := reflect.ValueOf(obj).Elem()
    structFieldValue := structValue.FieldByName(gstr.UcFirst(name))
    // 键名与对象属性匹配检测
    if !structFieldValue.IsValid() {
        //return fmt.Errorf("No such field: %s in obj", name)
        return
    }
    // CanSet的属性必须为公开属性(首字母大写)
    if !structFieldValue.CanSet() {
        //return fmt.Errorf("Cannot set %s field value", name)
        return
    }
    // 必须将value转换为struct属性的数据类型，这里必须用到gconv包
    structFieldValue.Set(reflect.ValueOf(gconv.Convert(value, structFieldValue.Type().String())))
}
=======
// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gutil provides utility functions.
package gutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/g/internal/empty"
	"github.com/gogf/gf/g/util/gconv"
	"os"
	"runtime"
)

// Dump prints variables <i...> to stdout with more manually readable.
func Dump(i...interface{}) {
    s := Export(i...)
    if s != "" {
        fmt.Println(s)
    }
}

// Export returns variables <i...> as a string with more manually readable.
func Export(i...interface{}) string {
    buffer := bytes.NewBuffer(nil)
    for _, v := range i {
        if b, ok := v.([]byte); ok {
            buffer.Write(b)
        } else {
            if m := gconv.Map(v); m != nil {
            	v = m
            }
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

// PrintBacktrace prints the caller backtrace to stdout.
func PrintBacktrace() {
    index  := 1
    buffer := bytes.NewBuffer(nil)
    for i := 1; i < 10000; i++ {
        if _, path, line, ok := runtime.Caller(i); ok {
            buffer.WriteString(fmt.Sprintf(`%d. %s:%d%s`, index, path, line, "\n"))
            index++
        } else {
            break
        }
    }
    fmt.Print(buffer.String())
}

// Throw throws out an exception, which can be caught be TryCatch or recover.
func Throw(exception interface{}) {
    panic(exception)
}

// TryCatch implements try...catch... logistics.
func TryCatch(try func(), catch ... func(exception interface{})) {
    if len(catch) > 0 {
    	// If <catch> is given, it's used to handle the exception.
        defer func() {
            if e := recover(); e != nil {
                catch[0](e)
            }
        }()
    } else {
    	// If no <catch> function passed, it filters the exception.
	    defer func() {
		    recover()
	    }()
    }
    try()
}

// IsEmpty checks given <value> empty or not.
// It returns false if <value> is: integer(0), bool(false), slice/map(len=0), nil;
// or else returns true.
func IsEmpty(value interface{}) bool {
    return empty.IsEmpty(value)
}


>>>>>>> upstream/master
