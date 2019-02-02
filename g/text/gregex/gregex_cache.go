// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gregex

import (
    "regexp"
    "sync"
)

// 缓存对象，主要用于缓存底层regx对象
var (
    regexMu  = sync.RWMutex{}
    regexMap = make(map[string]*regexp.Regexp)
)

// 根据pattern生成对应的regexp正则对象
func getRegexp(pattern string) (*regexp.Regexp, error) {
    if r := getCache(pattern); r != nil {
        return r, nil
    }
    if r, err := regexp.Compile(pattern); err == nil {
        setCache(pattern, r)
        return r, nil
    } else {
        return nil, err
    }
}

// 获得正则缓存对象
func getCache(pattern string) (regex *regexp.Regexp) {
    regexMu.RLock()
    regex = regexMap[pattern]
    regexMu.RUnlock()
    return
}

// 设置正则缓存对象
func setCache(pattern string, regex *regexp.Regexp) {
    regexMu.Lock()
    regexMap[pattern] = regex
    regexMu.Unlock()
}
