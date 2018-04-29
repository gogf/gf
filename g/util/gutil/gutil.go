// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 其他工具包
package gutil

import (
    "reflect"
    "gitee.com/johng/gf/g/util/gconv"
)

// 字符串首字母转换为大写
func UcFirst(s string) string {
    if len(s) == 0 {
        return s
    }
    if IsLetterLower(s[0]) {
        return string(s[0] - 32) + s[1 :]
    }
    return s
}

// 字符串首字母转换为小写
func LcFirst(s string) string {
    if len(s) == 0 {
        return s
    }
    if IsLetterUpper(s[0]) {
        return string(s[0] + 32) + s[1 :]
    }
    return s
}

// 便利数组查找字符串索引位置，如果不存在则返回-1，使用完整遍历查找
func StringSearch (a []string, s string) int {
    for i, v := range a {
        if s == v {
            return i
        }
    }
    return -1
}

// 判断字符串是否在数组中
func StringInArray (a []string, s string) bool {
    return StringSearch(a, s) != -1
}

// 判断给定字符是否小写
func IsLetterLower(b byte) bool {
    if b >= byte('a') && b <= byte('z') {
        return true
    }
    return false
}

// 判断给定字符是否大写
func IsLetterUpper(b byte) bool {
    if b >= byte('A') && b <= byte('Z') {
        return true
    }
    return false
}

// 判断锁给字符串是否为数字
func IsNumeric(s string) bool {
    for i := 0; i < len(s); i++ {
        if s[i] < byte('0') && s[i] > byte('9') {
            return false
        }
    }
    return true
}

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
    structFieldValue := structValue.FieldByName(UcFirst(name))
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