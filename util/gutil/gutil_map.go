// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gutil

import (
	"reflect"

	"github.com/gogf/gf/v2/internal/utils"
)

// MapCopy does a shallow copy from map `data` to `copy` for most commonly used map type
// map[string]interface{}.
func MapCopy(data map[string]interface{}) (copy map[string]interface{}) {
	copy = make(map[string]interface{}, len(data))
	for k, v := range data {
		copy[k] = v
	}
	return
}

// MapContains checks whether map `data` contains `key`.
func MapContains(data map[string]interface{}, key string) (ok bool) {
	if len(data) == 0 {
		return
	}
	_, ok = data[key]
	return
}

// MapDelete deletes all `keys` from map `data`.
func MapDelete(data map[string]interface{}, keys ...string) {
	if len(data) == 0 {
		return
	}
	for _, key := range keys {
		delete(data, key)
	}
}

// MapMerge merges all map from `src` to map `dst`.
func MapMerge(dstMap map[string]interface{}, srcMaps ...map[string]interface{}) {
	if dstMap == nil {
		return
	}
	for _, m := range srcMaps {
		for k, v := range m {
			dstMap[k] = v
		}
	}
}

// MapMergeCopy creates and returns a new map which merges all map from `src`.
func MapMergeCopy(maps ...map[string]interface{}) (copy map[string]interface{}) {
	copy = make(map[string]interface{})
	for _, m := range maps {
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
	return utils.MapPossibleItemByKey(data, key)
}

// MapContainsPossibleKey checks if the given `key` is contained in given map `data`.
// It checks the key ignoring cases and symbols.
//
// Note that this function might be of low performance.
func MapContainsPossibleKey(data map[string]interface{}, key string) bool {
	return utils.MapContainsPossibleKey(data, key)
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

// MapToSlice converts map to slice of which all keys and values are its items.
// Eg: {"K1": "v1", "K2": "v2"} => ["K1", "v1", "K2", "v2"]
func MapToSlice(data interface{}) []interface{} {
	var (
		reflectValue = reflect.ValueOf(data)
		reflectKind  = reflectValue.Kind()
	)
	for reflectKind == reflect.Ptr {
		reflectValue = reflectValue.Elem()
		reflectKind = reflectValue.Kind()
	}
	switch reflectKind {
	case reflect.Map:
		array := make([]interface{}, 0)
		for _, key := range reflectValue.MapKeys() {
			array = append(array, key.Interface())
			array = append(array, reflectValue.MapIndex(key).Interface())
		}
		return array
	}
	return nil
}
