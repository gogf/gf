// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gutil

import (
	"github.com/gogf/gf/internal/utils"
	"regexp"
)

var (
	// replaceCharReg is the regular expression object for replacing chars in map keys.
	replaceCharReg, _ = regexp.Compile(`[\-\.\_\s]+`)
)

// MapCopy does a shallow copy from map <data> to <copy> for most commonly used map type
// map[string]interface{}.
func MapCopy(data map[string]interface{}) (copy map[string]interface{}) {
	copy = make(map[string]interface{}, len(data))
	for k, v := range data {
		copy[k] = v
	}
	return
}

// MapContains checks whether map <data> contains <key>.
func MapContains(data map[string]interface{}, key string) (ok bool) {
	_, ok = data[key]
	return
}

// MapPossibleItemByKey tries to find the possible key-value pair for given key with or without
// cases or chars '-'/'_'/'.'/' '.
//
// Note that this function might be of low performance.
func MapPossibleItemByKey(data map[string]interface{}, key string) (foundKey string, foundValue interface{}) {
	if v, ok := data[key]; ok {
		return key, v
	}
	// Loop checking.
	for k, v := range data {
		if utils.EqualFoldWithoutChars(k, key) {
			return k, v
		}
	}
	return "", nil
}

// MapContainsPossibleKey checks if the given <key> is contained in given map <data>.
// It checks the key with or without cases or chars '-'/'_'/'.'/' '.
func MapContainsPossibleKey(data map[string]interface{}, key string) bool {
	if k, _ := MapPossibleItemByKey(data, key); k != "" {
		return true
	}
	return false
}
