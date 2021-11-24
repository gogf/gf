// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gtag providing tag content storing for struct.
//
// Note that calling functions of this package is concurrently safe.
package gtag

import (
	"regexp"
	"sync"
)

var (
	mu    sync.RWMutex
	data  = make(map[string]string)
	regex = regexp.MustCompile(`\{(.+?)\}`)
)

// Set sets tag content for specified name.
func Set(name, value string) {
	mu.Lock()
	defer mu.Unlock()
	data[name] = value
}

// Sets sets multiple tag content by map.
func Sets(m map[string]string) {
	mu.Lock()
	defer mu.Unlock()
	for k, v := range m {
		data[k] = v
	}
}

// Get retrieves and returns the stored tag content for specified name.
func Get(name string) string {
	mu.RLock()
	defer mu.RUnlock()
	return data[name]
}

// Parse parses and returns the content by replacing all tag name variable to
// its content for given `content`.
// Eg:
// If "Demo:content" in tag mapping,
// Parse(`This is {Demo}`) -> `This is content`.
func Parse(content string) string {
	mu.RLock()
	defer mu.RUnlock()
	return regex.ReplaceAllStringFunc(content, func(s string) string {
		if v, ok := data[s[1:len(s)-1]]; ok {
			return v
		}
		return s
	})
}
