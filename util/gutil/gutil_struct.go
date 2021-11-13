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

// StructToSlice converts struct to slice of which all keys and values are its items.
// Eg: {"K1": "v1", "K2": "v2"} => ["K1", "v1", "K2", "v2"]
func StructToSlice(data interface{}) []interface{} {
	var (
		reflectValue = reflect.ValueOf(data)
		reflectKind  = reflectValue.Kind()
	)
	for reflectKind == reflect.Ptr {
		reflectValue = reflectValue.Elem()
		reflectKind = reflectValue.Kind()
	}
	switch reflectKind {
	case reflect.Struct:
		array := make([]interface{}, 0)
		// Note that, it uses the gconv tag name instead of the attribute name if
		// the gconv tag is fined in the struct attributes.
		for k, v := range gconv.Map(reflectValue) {
			array = append(array, k)
			array = append(array, v)
		}
		return array
	}
	return nil
}
