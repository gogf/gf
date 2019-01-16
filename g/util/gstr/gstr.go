// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// Package gstr provides useful API for string handling.
// 
// 字符串操作.
package gstr

import (
    "bytes"
    "math"
    "strings"
)

// 字符串替换(大小写敏感)
func Replace(origin, search, replace string, count...int) string {
    n := -1
    if len(count) > 0 {
        n = count[0]
    }
    return strings.Replace(origin, search, replace, n)
}

// 使用map进行字符串替换(大小写敏感)
func ReplaceByMap(origin string, replaces map[string]string) string {
    result := origin
    for k, v := range replaces {
        result = strings.Replace(result, k, v, -1)
    }
    return result
}

// 字符串转换为小写
func ToLower(s string) string {
    return strings.ToLower(s)
}

// 字符串转换为大写
func ToUpper(s string) string {
    return strings.ToUpper(s)
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
    length := len(s)
    if length == 0 {
        return false
    }
    for i := 0; i < len(s); i++ {
        if s[i] < byte('0') || s[i] > byte('9') {
            return false
        }
    }
    return true
}

// 字符串截取，支持中文
func SubStr(str string, start int, length...int) (substr string) {
    // 将字符串的转换成[]rune
    rs  := []rune(str)
    lth := len(rs)
    // 简单的越界判断
    if start < 0 {
        start = 0
    }
    if start >= lth {
        start = lth
    }
    end := lth
    if len(length) > 0 {
        end = start + length[0]
        if end < start {
            end = lth
        }
    }
    if end > lth {
        end = lth
    }
    // 返回子串
    return string(rs[start : end])
}

// 字符串长度截取限制，超过长度限制被截取并在字符串末尾追加指定的内容，支持中文
func StrLimit(str string, length int, suffix...string) (string) {
    rs := []rune(str)
    if len(str) < length {
        return str
    }
    addstr := "..."
    if len(suffix) > 0 {
        addstr = suffix[0]
    }
    return string(rs[0 : length]) + addstr
}

// 按照百分比从字符串中间向两边隐藏字符(主要用于姓名、手机号、邮箱地址、身份证号等的隐藏)，支持utf-8中文，支持email格式。
func HideStr(str string, percent int, hide string) string {
    array := strings.Split(str, "@")
    if len(array) > 1 {
        str = array[0]
    }
    rs       := []rune(str)
    length   := len(rs)
    mid      := math.Floor(float64(length/2))
    hideLen  := int(math.Floor(float64(length) * (float64(percent)/100)))
    start    := int(mid - math.Floor(float64(hideLen) / 2))
    hideStr  := []rune("")
    hideRune := []rune(hide)
    for i := 0; i < int(hideLen); i++ {
        hideStr = append(hideStr, hideRune...)
    }
    buffer := bytes.NewBuffer(nil)
    buffer.WriteString(string(rs[0 : start]))
    buffer.WriteString(string(hideStr))
    buffer.WriteString(string(rs[start + hideLen : ]))
    if len(array) > 1 {
        buffer.WriteString(array[1])
    }
    return buffer.String()
}

// 将\n\r替换为html中的<br>标签。
func Nl2Br(str string) string {
    str = Replace(str, "\r\n", "\n")
    str = Replace(str, "\n\r", "\n")
    str = Replace(str, "\n", "<br />")
    return str
}
