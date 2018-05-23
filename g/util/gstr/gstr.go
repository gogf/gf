// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 字符串操作.
package gstr

import "strings"

// 字符串替换
func Replace(origin, search, replace string) string {
    return strings.Replace(origin, search, replace, -1)
}

// 使用map进行字符串替换
func ReplaceByMap(origin string, replaces map[string]string) string {
    result := origin
    for k, v := range replaces {
        result = strings.Replace(result, k, v, -1)
    }
    return result
}

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
func SearchArray (a []string, s string) int {
    for i, v := range a {
        if s == v {
            return i
        }
    }
    return -1
}

// 判断字符串是否在数组中
func InArray (a []string, s string) bool {
    return SearchArray(a, s) != -1
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
        if s[i] < byte('0') || s[i] > byte('9') {
            return false
        }
    }
    return true
}