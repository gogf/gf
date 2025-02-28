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
	"github.com/gogf/gf/v2/internal/utils"
	"github.com/gogf/gf/v2/util/gconv/internal/localinterface"
)

// SliceInt is alias of Ints.
func SliceInt(any interface{}) []int {
	return Ints(any)
}

// SliceInt32 is alias of Int32s.
func SliceInt32(any interface{}) []int32 {
	return Int32s(any)
}

// SliceInt64 is alias of Int64s.
func SliceInt64(any interface{}) []int64 {
	return Int64s(any)
}

// Ints converts `any` to []int.
func Ints(any interface{}) []int {
	if empty.IsNil(any) {
		return nil
	}
	var (
		array []int = nil
	)
	switch value := any.(type) {
	case []string:
		array = make([]int, len(value))
		for k, v := range value {
			array[k] = Int(v)
		}
	case []int:
		array = value
	case []int8:
		array = make([]int, len(value))
		for k, v := range value {
			array[k] = int(v)
		}
	case []int16:
		array = make([]int, len(value))
		for k, v := range value {
			array[k] = int(v)
		}
	case []int32:
		array = make([]int, len(value))
		for k, v := range value {
			array[k] = int(v)
		}
	case []int64:
		array = make([]int, len(value))
		for k, v := range value {
			array[k] = int(v)
		}
	case []uint:
		array = make([]int, len(value))
		for k, v := range value {
			array[k] = int(v)
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
		array = make([]int, len(value))
		for k, v := range value {
			array[k] = int(v)
		}
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
			return []int{}
		}
		if utils.IsNumeric(value) {
			return []int{Int(value)}
		}
	case []uint16:
		array = make([]int, len(value))
		for k, v := range value {
			array[k] = int(v)
		}
	case []uint32:
		array = make([]int, len(value))
		for k, v := range value {
			array[k] = int(v)
		}
	case []uint64:
		array = make([]int, len(value))
		for k, v := range value {
			array[k] = int(v)
		}
	case []bool:
		array = make([]int, len(value))
		for k, v := range value {
			if v {
				array[k] = 1
			} else {
				array[k] = 0
			}
		}
	case []float32:
		array = make([]int, len(value))
		for k, v := range value {
			array[k] = Int(v)
		}
	case []float64:
		array = make([]int, len(value))
		for k, v := range value {
			array[k] = Int(v)
		}
	case []interface{}:
		array = make([]int, len(value))
		for k, v := range value {
			array[k] = Int(v)
		}
	case [][]byte:
		array = make([]int, len(value))
		for k, v := range value {
			array[k] = Int(v)
		}
	}
	if array != nil {
		return array
	}
	if v, ok := any.(localinterface.IInts); ok {
		return v.Ints()
	}
	if v, ok := any.(localinterface.IInterfaces); ok {
		return Ints(v.Interfaces())
	}
	// Not a common type, it then uses reflection for conversion.
	originValueAndKind := reflection.OriginValueAndKind(any)
	switch originValueAndKind.OriginKind {
	case reflect.Slice, reflect.Array:
		var (
			length = originValueAndKind.OriginValue.Len()
			slice  = make([]int, length)
		)
		for i := 0; i < length; i++ {
			slice[i] = Int(originValueAndKind.OriginValue.Index(i).Interface())
		}
		return slice

	default:
		if originValueAndKind.OriginValue.IsZero() {
			return []int{}
		}
		return []int{Int(any)}
	}
}

// Int32s converts `any` to []int32.
func Int32s(any interface{}) []int32 {
	if empty.IsNil(any) {
		return nil
	}
	var (
		array []int32 = nil
	)
	switch value := any.(type) {
	case []string:
		array = make([]int32, len(value))
		for k, v := range value {
			array[k] = Int32(v)
		}
	case []int:
		array = make([]int32, len(value))
		for k, v := range value {
			array[k] = int32(v)
		}
	case []int8:
		array = make([]int32, len(value))
		for k, v := range value {
			array[k] = int32(v)
		}
	case []int16:
		array = make([]int32, len(value))
		for k, v := range value {
			array[k] = int32(v)
		}
	case []int32:
		array = value
	case []int64:
		array = make([]int32, len(value))
		for k, v := range value {
			array[k] = int32(v)
		}
	case []uint:
		array = make([]int32, len(value))
		for k, v := range value {
			array[k] = int32(v)
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
		array = make([]int32, len(value))
		for k, v := range value {
			array[k] = int32(v)
		}
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
			return []int32{}
		}
		if utils.IsNumeric(value) {
			return []int32{Int32(value)}
		}
	case []uint16:
		array = make([]int32, len(value))
		for k, v := range value {
			array[k] = int32(v)
		}
	case []uint32:
		array = make([]int32, len(value))
		for k, v := range value {
			array[k] = int32(v)
		}
	case []uint64:
		array = make([]int32, len(value))
		for k, v := range value {
			array[k] = int32(v)
		}
	case []bool:
		array = make([]int32, len(value))
		for k, v := range value {
			if v {
				array[k] = 1
			} else {
				array[k] = 0
			}
		}
	case []float32:
		array = make([]int32, len(value))
		for k, v := range value {
			array[k] = Int32(v)
		}
	case []float64:
		array = make([]int32, len(value))
		for k, v := range value {
			array[k] = Int32(v)
		}
	case []interface{}:
		array = make([]int32, len(value))
		for k, v := range value {
			array[k] = Int32(v)
		}
	case [][]byte:
		array = make([]int32, len(value))
		for k, v := range value {
			array[k] = Int32(v)
		}
	}
	if array != nil {
		return array
	}
	if v, ok := any.(localinterface.IInts); ok {
		return Int32s(v.Ints())
	}
	if v, ok := any.(localinterface.IInterfaces); ok {
		return Int32s(v.Interfaces())
	}
	// Not a common type, it then uses reflection for conversion.
	originValueAndKind := reflection.OriginValueAndKind(any)
	switch originValueAndKind.OriginKind {
	case reflect.Slice, reflect.Array:
		var (
			length = originValueAndKind.OriginValue.Len()
			slice  = make([]int32, length)
		)
		for i := 0; i < length; i++ {
			slice[i] = Int32(originValueAndKind.OriginValue.Index(i).Interface())
		}
		return slice

	default:
		if originValueAndKind.OriginValue.IsZero() {
			return []int32{}
		}
		return []int32{Int32(any)}
	}
}

// Int64s converts `any` to []int64.
func Int64s(any interface{}) []int64 {
	if empty.IsNil(any) {
		return nil
	}
	var (
		array []int64 = nil
	)
	switch value := any.(type) {
	case []string:
		array = make([]int64, len(value))
		for k, v := range value {
			array[k] = Int64(v)
		}
	case []int:
		array = make([]int64, len(value))
		for k, v := range value {
			array[k] = int64(v)
		}
	case []int8:
		array = make([]int64, len(value))
		for k, v := range value {
			array[k] = int64(v)
		}
	case []int16:
		array = make([]int64, len(value))
		for k, v := range value {
			array[k] = int64(v)
		}
	case []int32:
		array = make([]int64, len(value))
		for k, v := range value {
			array[k] = int64(v)
		}
	case []int64:
		array = value
	case []uint:
		array = make([]int64, len(value))
		for k, v := range value {
			array[k] = int64(v)
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
		array = make([]int64, len(value))
		for k, v := range value {
			array[k] = int64(v)
		}
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
			return []int64{}
		}
		if utils.IsNumeric(value) {
			return []int64{Int64(value)}
		}
	case []uint16:
		array = make([]int64, len(value))
		for k, v := range value {
			array[k] = int64(v)
		}
	case []uint32:
		array = make([]int64, len(value))
		for k, v := range value {
			array[k] = int64(v)
		}
	case []uint64:
		array = make([]int64, len(value))
		for k, v := range value {
			array[k] = int64(v)
		}
	case []bool:
		array = make([]int64, len(value))
		for k, v := range value {
			if v {
				array[k] = 1
			} else {
				array[k] = 0
			}
		}
	case []float32:
		array = make([]int64, len(value))
		for k, v := range value {
			array[k] = Int64(v)
		}
	case []float64:
		array = make([]int64, len(value))
		for k, v := range value {
			array[k] = Int64(v)
		}
	case []interface{}:
		array = make([]int64, len(value))
		for k, v := range value {
			array[k] = Int64(v)
		}
	case [][]byte:
		array = make([]int64, len(value))
		for k, v := range value {
			array[k] = Int64(v)
		}
	}
	if array != nil {
		return array
	}
	if v, ok := any.(localinterface.IInts); ok {
		return Int64s(v.Ints())
	}
	if v, ok := any.(localinterface.IInterfaces); ok {
		return Int64s(v.Interfaces())
	}
	// Not a common type, it then uses reflection for conversion.
	originValueAndKind := reflection.OriginValueAndKind(any)
	switch originValueAndKind.OriginKind {
	case reflect.Slice, reflect.Array:
		var (
			length = originValueAndKind.OriginValue.Len()
			slice  = make([]int64, length)
		)
		for i := 0; i < length; i++ {
			slice[i] = Int64(originValueAndKind.OriginValue.Index(i).Interface())
		}
		return slice

	default:
		if originValueAndKind.OriginValue.IsZero() {
			return []int64{}
		}
		return []int64{Int64(any)}
	}
}
