// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv

import (
	"reflect"
)

// Convert converts the variable `fromValue` to the type `toTypeName`, the type `toTypeName` is specified by string.
//
// The optional parameter `extraParams` is used for additional necessary parameter for this conversion.
// It supports common basic types conversion as its conversion based on type name string.
func Convert(fromValue any, toTypeName string, extraParams ...any) any {
	result, _ := defaultConverter.doConvert(
		doConvertInput{
			FromValue:  fromValue,
			ToTypeName: toTypeName,
			ReferValue: nil,
			Extra:      extraParams,
		},
	)
	return result
}

// ConvertWithRefer converts the variable `fromValue` to the type referred by value `referValue`.
//
// The optional parameter `extraParams` is used for additional necessary parameter for this conversion.
// It supports common basic types conversion as its conversion based on type name string.
func ConvertWithRefer(fromValue any, referValue any, extraParams ...any) any {
	var referValueRf reflect.Value
	if v, ok := referValue.(reflect.Value); ok {
		referValueRf = v
	} else {
		referValueRf = reflect.ValueOf(referValue)
	}
	result, _ := defaultConverter.doConvert(doConvertInput{
		FromValue:  fromValue,
		ToTypeName: referValueRf.Type().String(),
		ReferValue: referValue,
		Extra:      extraParams,
	})
	return result
}
