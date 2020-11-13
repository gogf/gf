// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gutil

import (
	"github.com/gogf/gf/internal/utils"
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
	if len(data) == 0 {
		return
	}
	_, ok = data[key]
	return
}

// MapDelete deletes all <keys> from map <data>.
func MapDelete(data map[string]interface{}, keys ...string) {
	if len(data) == 0 {
		return
	}
	for _, key := range keys {
		delete(data, key)
	}
}

// MapMerge merges all map from <src> to map <dst>.
func MapMerge(dst map[string]interface{}, src ...map[string]interface{}) {
	if dst == nil {
		return
	}
	for _, m := range src {
		for k, v := range m {
			dst[k] = v
		}
	}
}

// MapMergeCopy creates and returns a new map which merges all map from <src>.
func MapMergeCopy(src ...map[string]interface{}) (copy map[string]interface{}) {
	copy = make(map[string]interface{})
	for _, m := range src {
		for k, v := range m {
			copy[k] = v
		}
	}
	return
}

// MapPossibleItemByKey tries to find the possible key-value pair for given key ignoring cases and symbols.
//
// Note that this function might be of low performance.
func MapPossibleItemByKey(data map[string]interface{}, key string) (foundKey string, foundValue interface{}) {
	if len(data) == 0 {
		return
	}
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
// It checks the key ignoring cases and symbols.
//
// Note that this function might be of low performance.
func MapContainsPossibleKey(data map[string]interface{}, key string) bool {
	if k, _ := MapPossibleItemByKey(data, key); k != "" {
		return true
	}
	return false
}

// MapOmitEmpty deletes all empty values from given map.
func MapOmitEmpty(data map[string]interface{}) {
	if len(data) == 0 {
		return
	}
	for k, v := range data {
		if IsEmpty(v) {
			delete(data, k)
		}
	}
}
