// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gutil

import (
	"regexp"
	"strings"
)

var (
	// replaceCharReg is the regular expression object for replacing chars in map keys.
	replaceCharReg, _ = regexp.Compile(`[\-\.\_\s]+`)
)

// CopyMap does a shallow copy from map <data> to <copy> for most commonly used map type
// map[string]interface{}.
func CopyMap(data map[string]interface{}) (copy map[string]interface{}) {
	copy = make(map[string]interface{}, len(data))
	for k, v := range data {
		copy[k] = v
	}
	return
}

// MapPossibleItemByKey tries to find the possible key-value pair for given key with or without
// cases or chars '-'/'_'/'.'/' '.
//
// Note that this function might be of low performance.
func MapPossibleItemByKey(data map[string]interface{}, key string) (string, interface{}) {
	if v, ok := data[key]; ok {
		return key, v
	}
	replacedKey := replaceCharReg.ReplaceAllString(key, "")
	if v, ok := data[replacedKey]; ok {
		return replacedKey, v
	}
	// Loop for check.
	for k, v := range data {
		if strings.EqualFold(replaceCharReg.ReplaceAllString(k, ""), replacedKey) {
			return k, v
		}
	}
	return "", nil
}
