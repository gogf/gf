// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package builtin

import (
	"reflect"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/util/gconv"
)

func getValueLength(value *gvar.Var) int {
	if value == nil || value.Val() == nil {
		return 0
	}
	reflectValue := reflect.ValueOf(value.Val())
	for reflectValue.IsValid() && (reflectValue.Kind() == reflect.Pointer || reflectValue.Kind() == reflect.Interface) {
		if reflectValue.IsNil() {
			return 0
		}
		reflectValue = reflectValue.Elem()
	}
	if reflectValue.IsValid() {
		switch reflectValue.Kind() {
		case reflect.String:
			return len(gconv.Runes(reflectValue.String()))

		case reflect.Array, reflect.Slice:
			return reflectValue.Len()
		}
	}
	return len(gconv.Runes(value.String()))
}
