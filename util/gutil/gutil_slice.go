// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gutil

import (
	"reflect"

	"github.com/gogf/gf/v2/util/gconv"
)

// SliceCopy does a shallow copy of slice `data` for most commonly used slice type
// []interface{}.
func SliceCopy(slice []interface{}) []interface{} {
	newSlice := make([]interface{}, len(slice))
	copy(newSlice, slice)
	return newSlice
}

// SliceInsertBefore inserts the `values` to the front of `index` and returns a new slice.
func SliceInsertBefore(slice []interface{}, index int, values ...interface{}) (newSlice []interface{}) {
	if index < 0 || index >= len(slice) {
		return slice
	}
	newSlice = make([]interface{}, len(slice)+len(values))
	copy(newSlice, slice[0:index])
	copy(newSlice[index:], values)
	copy(newSlice[index+len(values):], slice[index:])
	return
}

// SliceInsertAfter inserts the `values` to the back of `index` and returns a new slice.
func SliceInsertAfter(slice []interface{}, index int, values ...interface{}) (newSlice []interface{}) {
	if index < 0 || index >= len(slice) {
		return slice
	}
	newSlice = make([]interface{}, len(slice)+len(values))
	copy(newSlice, slice[0:index+1])
	copy(newSlice[index+1:], values)
	copy(newSlice[index+1+len(values):], slice[index+1:])
	return
}

// SliceDelete deletes an element at `index` and returns the new slice.
// It does nothing if the given `index` is invalid.
func SliceDelete(slice []interface{}, index int) (newSlice []interface{}) {
	if index < 0 || index >= len(slice) {
		return slice
	}
	// Determine array boundaries when deleting to improve deletion efficiency.
	if index == 0 {
		return slice[1:]
	} else if index == len(slice)-1 {
		return slice[:index]
	}
	// If it is a non-boundary delete,
	// it will involve the creation of an array,
	// then the deletion is less efficient.
	return append(slice[:index], slice[index+1:]...)
}

// SliceToMap converts slice type variable `slice` to `map[string]interface{}`.
// Note that if the length of `slice` is not an even number, it returns nil.
// Eg:
// ["K1", "v1", "K2", "v2"] => {"K1": "v1", "K2": "v2"}
// ["K1", "v1", "K2"]       => nil
func SliceToMap(slice interface{}) map[string]interface{} {
	var (
		reflectValue = reflect.ValueOf(slice)
		reflectKind  = reflectValue.Kind()
	)
	for reflectKind == reflect.Ptr {
		reflectValue = reflectValue.Elem()
		reflectKind = reflectValue.Kind()
	}
	switch reflectKind {
	case reflect.Slice, reflect.Array:
		length := reflectValue.Len()
		if length%2 != 0 {
			return nil
		}
		data := make(map[string]interface{})
		for i := 0; i < reflectValue.Len(); i += 2 {
			data[gconv.String(reflectValue.Index(i).Interface())] = reflectValue.Index(i + 1).Interface()
		}
		return data
	}
	return nil
}

// SliceToMapWithColumnAsKey converts slice type variable `slice` to `map[interface{}]interface{}`
// The value of specified column use as the key for returned map.
// Eg:
// SliceToMapWithColumnAsKey([{"K1": "v1", "K2": 1}, {"K1": "v2", "K2": 2}], "K1") => {"v1": {"K1": "v1", "K2": 1}, "v2": {"K1": "v2", "K2": 2}}
// SliceToMapWithColumnAsKey([{"K1": "v1", "K2": 1}, {"K1": "v2", "K2": 2}], "K2") => {1: {"K1": "v1", "K2": 1}, 2: {"K1": "v2", "K2": 2}}
func SliceToMapWithColumnAsKey(slice interface{}, key interface{}) map[interface{}]interface{} {
	var (
		reflectValue = reflect.ValueOf(slice)
		reflectKind  = reflectValue.Kind()
	)
	for reflectKind == reflect.Ptr {
		reflectValue = reflectValue.Elem()
		reflectKind = reflectValue.Kind()
	}
	data := make(map[interface{}]interface{})
	switch reflectKind {
	case reflect.Slice, reflect.Array:
		for i := 0; i < reflectValue.Len(); i++ {
			if k, ok := ItemValue(reflectValue.Index(i), key); ok {
				data[k] = reflectValue.Index(i).Interface()
			}
		}
	}
	return data
}
