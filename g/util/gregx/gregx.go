// Copyright 2017-2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 正则表达式.
package gregx

import (
    "regexp"
    "gitee.com/johng/gf/g/os/gcache"
)

// 缓存对象，主要用于缓存底层regx对象
var regxCache = gcache.New()

// 根据pattern生成对应的regexp正则对象
func getRegexp(pattern string) (*regexp.Regexp, error) {
    if v := regxCache.Get(pattern); v != nil {
        return v.(*regexp.Regexp), nil
    }
    if r, err := regexp.Compile(pattern); err == nil {
        regxCache.Set(pattern, r, 0)
        return r, nil
    } else {
        return nil, err
    }
}

// 校验所给定的正则表达式是否符合规范
func Validate(pattern string) error {
    _, err := getRegexp(pattern)
    return err
}

// 正则表达式是否匹配
func IsMatch(pattern string, src []byte) bool {
    if r, err := getRegexp(pattern); err == nil {
        return r.Match(src)
    }
    return false
}

func IsMatchString(pattern string, src string) bool {
    return IsMatch(pattern, []byte(src))
}

// 正则匹配，并返回匹配的列表
func MatchString(pattern string, src string) ([]string, error) {
    if r, err := getRegexp(pattern); err == nil {
        return r.FindStringSubmatch(src), nil
    } else {
        return nil, err
    }
}

func MatchAllString(pattern string, src string) ([][]string, error) {
    if r, err := getRegexp(pattern); err == nil {
        return r.FindAllStringSubmatch(src, -1), nil
    } else {
        return nil, err
    }
}

// 正则替换(全部替换)
func Replace(pattern string, replace, src []byte) ([]byte, error) {
    if r, err := getRegexp(pattern); err == nil {
        return r.ReplaceAll(src, replace), nil
    } else {
        return nil, err
    }
}

// 正则替换(全部替换)，字符串
func ReplaceString(pattern, replace, src string) (string, error) {
    r, e := Replace(pattern, []byte(replace), []byte(src))
    return string(r), e
}