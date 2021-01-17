// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gutil

import (
	"github.com/gogf/gf/util/gconv"
	"reflect"
)

// SliceCopy does a shallow copy of slice <data> for most commonly used slice type
// []interface{}.
func SliceCopy(data []interface{}) []interface{} {
	newData := make([]interface{}, len(data))
	copy(newData, data)
	return newData
}

// SliceDelete deletes an element at <index> and returns the new slice.
// It does nothing if the given <index> is invalid.
func SliceDelete(data []interface{}, index int) (newSlice []interface{}) {
	if index < 0 || index >= len(data) {
		return data
	}
	// Determine array boundaries when deleting to improve deletion efficiency.
	if index == 0 {
		return data[1:]
	} else if index == len(data)-1 {
		return data[:index]
	}
	// If it is a non-boundary delete,
	// it will involve the creation of an array,
	// then the deletion is less efficient.
	return append(data[:index], data[index+1:]...)
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
