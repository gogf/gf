// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package utils

import "reflect"

// IsArray checks whether given value is array/slice.
// Note that it uses reflect internally implementing this feature.
func IsArray(value interface{}) bool {
	rv := reflect.ValueOf(value)
	kind := rv.Kind()
	if kind == reflect.Ptr {
		rv = rv.Elem()
		kind = rv.Kind()
	}
	switch kind {
	case reflect.Array, reflect.Slice:
		return true
	default:
		return false
	}
}
