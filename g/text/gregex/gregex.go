// Copyright 2017-2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gregex provides high performance API for regular expression functionality.
// 
// 正则表达式.
package gregex

import (
    "regexp"
)

// Quote quotes <s> by replacing special chars in <s>
// to match the rules of regular expression pattern.
// And returns the copy.
//
// 转移正则规则字符串，例如：Quote(`[foo]`) 返回 `\[foo\]`
func Quote(s string) string {
    return regexp.QuoteMeta(s)
}

// Validate checks whether given regular expression pattern <pattern> valid.
//
// 校验所给定的正则表达式是否符合规范
func Validate(pattern string) error {
    _, err := getRegexp(pattern)
    return err
}

// IsMatch checks whether given bytes <src> matches <pattern>.
//
// 正则表达式是否匹配
func IsMatch(pattern string, src []byte) bool {
    if r, err := getRegexp(pattern); err == nil {
        return r.Match(src)
    }
    return false
}

// IsMatchString checks whether given string <src> matches <pattern>.
//
// 判断给定的字符串<src>是否满足正则表达式<pattern>.
func IsMatchString(pattern string, src string) bool {
    return IsMatch(pattern, []byte(src))
}

// MatchString return bytes slice that matched <pattern>.
//
// 正则匹配，并返回匹配的列表(参数[]byte)
func Match(pattern string, src []byte) ([][]byte, error) {
    if r, err := getRegexp(pattern); err == nil {
        return r.FindSubmatch(src), nil
    } else {
        return nil, err
    }
}

// MatchString return strings that matched <pattern>.
//
// 正则匹配，并返回匹配的列表(参数[]string)
func MatchString(pattern string, src string) ([]string, error) {
    if r, err := getRegexp(pattern); err == nil {
        return r.FindStringSubmatch(src), nil
    } else {
        return nil, err
    }
}

// MatchAll return all bytes slices that matched <pattern>.
//
// 正则匹配，并返回所有匹配的列表(参数[]string)
func MatchAll(pattern string, src []byte) ([][][]byte, error) {
    if r, err := getRegexp(pattern); err == nil {
        return r.FindAllSubmatch(src, -1), nil
    } else {
        return nil, err
    }
}

// MatchAllString return all strings that matched <pattern>.
//
// 正则匹配，并返回所有匹配的列表(参数[][]string).
func MatchAllString(pattern string, src string) ([][]string, error) {
    if r, err := getRegexp(pattern); err == nil {
        return r.FindAllStringSubmatch(src, -1), nil
    } else {
        return nil, err
    }
}

// ReplaceString replace all matched <pattern> in bytes <src> with bytes <replace>.
//
// 正则替换(全部替换).
func Replace(pattern string, replace, src []byte) ([]byte, error) {
    if r, err := getRegexp(pattern); err == nil {
        return r.ReplaceAll(src, replace), nil
    } else {
        return nil, err
    }
}

// ReplaceString replace all matched <pattern> in string <src> with string <replace>.
//
// 正则替换(全部替换)，字符串
func ReplaceString(pattern, replace, src string) (string, error) {
    r, e := Replace(pattern, []byte(replace), []byte(src))
    return string(r), e
}

// ReplaceFunc replace all matched <pattern> in bytes <src>
// with custom replacement function <replaceFunc>.
//
// 正则替换(全部替换)，给定自定义替换方法
func ReplaceFunc(pattern string, src []byte, replaceFunc func(b []byte) []byte) ([]byte, error) {
    if r, err := getRegexp(pattern); err == nil {
        return r.ReplaceAllFunc(src, replaceFunc), nil
    } else {
        return nil, err
    }
}

// ReplaceStringFunc replace all matched <pattern> in string <src>
// with custom replacement function <replaceFunc>.
//
// 正则替换(全部替换)，给定自定义替换方法
func ReplaceStringFunc(pattern string, src string, replaceFunc func(s string) string) (string, error) {
    bytes, err := ReplaceFunc(pattern, []byte(src), func(bytes []byte) []byte {
        return []byte(replaceFunc(string(bytes)))
    })
    return string(bytes), err
}

// Split slices s into substrings separated by the expression and returns a slice of
// the substrings between those expression matches.
//
// 通过一个正则表达式分隔字符串.
func Split(pattern string, src string) []string {
    if r, err := getRegexp(pattern); err == nil {
        return r.Split(src, -1)
    }
    return nil
}
