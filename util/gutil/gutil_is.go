// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gutil

import (
	"reflect"

	"github.com/gogf/gf/v2/internal/empty"
)

// IsEmpty checks given `value` empty or not.
// It returns false if `value` is: integer(0), bool(false), slice/map(len=0), nil;
// or else returns true.
func IsEmpty(value interface{}) bool {
	return empty.IsEmpty(value)
}

// IsTypeOf checks and returns whether the type of `value` and `valueInExpectType` equal.
func IsTypeOf(value, valueInExpectType interface{}) bool {
	return reflect.TypeOf(value) == reflect.TypeOf(valueInExpectType)
}
