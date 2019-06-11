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

var (
    regexMu  = sync.RWMutex{}
    regexMap = make(map[string]*regexp.Regexp)
)

// getRegexp returns *regexp.Regexp object with given <pattern>.
// It uses cache to enhance the performance for compiling regular expression pattern,
// which means, it will return the same *regexp.Regexp object with the same regular
// expression pattern.
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

// getCache returns *regexp.Regexp object from cache by given <pattern>, for internal usage.
func getCache(pattern string) (regex *regexp.Regexp) {
    regexMu.RLock()
    regex = regexMap[pattern]
    regexMu.RUnlock()
    return
}

// setCache stores *regexp.Regexp object into cache, for internal usage.
func setCache(pattern string, regex *regexp.Regexp) {
    regexMu.Lock()
    regexMap[pattern] = regex
    regexMu.Unlock()
}
