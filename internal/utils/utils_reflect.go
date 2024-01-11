// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package utils

import (
	"reflect"
)

// CanCallIsNil Can reflect.Value call reflect.Value.IsNil.
// It can avoid reflect.Value.IsNil panics.
func CanCallIsNil(v interface{}) bool {
	rv, ok := v.(reflect.Value)
	if !ok {
		return false
	}
	switch rv.Kind() {
	case reflect.Interface, reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.Slice, reflect.UnsafePointer:
		return true
	default:
		return false
	}
}
