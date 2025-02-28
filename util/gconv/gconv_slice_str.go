// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv

import (
	"bytes"
	"reflect"
	"strings"

	"github.com/gogf/gf/v2/internal/empty"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/internal/reflection"
	"github.com/gogf/gf/v2/util/gconv/internal/localinterface"
)

// SliceStr is alias of Strings.
func SliceStr(any interface{}) []string {
	return Strings(any)
}

// Strings converts `any` to []string.
func Strings(any interface{}) []string {
	if empty.IsNil(any) {
		return nil
	}
	var (
		array []string = nil
	)
	switch value := any.(type) {
	case []int:
		array = make([]string, len(value))
		for k, v := range value {
			array[k] = String(v)
		}
	case []int8:
		array = make([]string, len(value))
		for k, v := range value {
			array[k] = String(v)
		}
	case []int16:
		array = make([]string, len(value))
		for k, v := range value {
			array[k] = String(v)
		}
	case []int32:
		array = make([]string, len(value))
		for k, v := range value {
			array[k] = String(v)
		}
	case []int64:
		array = make([]string, len(value))
		for k, v := range value {
			array[k] = String(v)
		}
	case []uint:
		array = make([]string, len(value))
		for k, v := range value {
			array[k] = String(v)
		}
	case []uint8:
		if json.Valid(value) {
			if _ = json.UnmarshalUseNumber(value, &array); array != nil {
				return array
			}
			if bytes.EqualFold([]byte("null"), value) {
				return nil
			}
		}
		array = make([]string, len(value))
		for k, v := range value {
			array[k] = String(v)
		}
		return array
	case string:
		byteValue := []byte(value)
		if json.Valid(byteValue) {
			if _ = json.UnmarshalUseNumber(byteValue, &array); array != nil {
				return array
			}
			if strings.EqualFold(value, "null") {
				return nil
			}
		}
		if value == "" {
			return []string{}
		}
		return []string{value}
	case []uint16:
		array = make([]string, len(value))
		for k, v := range value {
			array[k] = String(v)
		}
	case []uint32:
		array = make([]string, len(value))
		for k, v := range value {
			array[k] = String(v)
		}
	case []uint64:
		array = make([]string, len(value))
		for k, v := range value {
			array[k] = String(v)
		}
	case []bool:
		array = make([]string, len(value))
		for k, v := range value {
			array[k] = String(v)
		}
	case []float32:
		array = make([]string, len(value))
		for k, v := range value {
			array[k] = String(v)
		}
	case []float64:
		array = make([]string, len(value))
		for k, v := range value {
			array[k] = String(v)
		}
	case []interface{}:
		array = make([]string, len(value))
		for k, v := range value {
			array[k] = String(v)
		}
	case []string:
		array = value
	case [][]byte:
		array = make([]string, len(value))
		for k, v := range value {
			array[k] = String(v)
		}
	}
	if array != nil {
		return array
	}
	if v, ok := any.(localinterface.IStrings); ok {
		return v.Strings()
	}
	if v, ok := any.(localinterface.IInterfaces); ok {
		return Strings(v.Interfaces())
	}
	// Not a common type, it then uses reflection for conversion.
	originValueAndKind := reflection.OriginValueAndKind(any)
	switch originValueAndKind.OriginKind {
	case reflect.Slice, reflect.Array:
		var (
			length = originValueAndKind.OriginValue.Len()
			slice  = make([]string, length)
		)
		for i := 0; i < length; i++ {
			slice[i] = String(originValueAndKind.OriginValue.Index(i).Interface())
		}
		return slice

	default:
		if originValueAndKind.OriginValue.IsZero() {
			return []string{}
		}
		return []string{String(any)}
	}
}
