// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv

import (
	"reflect"
	"strings"

	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/internal/reflection"
	"github.com/gogf/gf/v2/internal/utils"
	"github.com/gogf/gf/v2/util/gconv/internal/localinterface"
)

// SliceUint is alias of Uints.
func SliceUint(any interface{}) []uint {
	return Uints(any)
}

// SliceUint32 is alias of Uint32s.
func SliceUint32(any interface{}) []uint32 {
	return Uint32s(any)
}

// SliceUint64 is alias of Uint64s.
func SliceUint64(any interface{}) []uint64 {
	return Uint64s(any)
}

// Uints converts `any` to []uint.
func Uints(any interface{}) []uint {
	if any == nil {
		return nil
	}

	var (
		array []uint = nil
	)
	switch value := any.(type) {
	case string:
		value = strings.TrimSpace(value)
		if value == "" {
			return []uint{}
		}
		if utils.IsNumeric(value) {
			return []uint{Uint(value)}
		}

	case []string:
		array = make([]uint, len(value))
		for k, v := range value {
			array[k] = Uint(v)
		}
	case []int8:
		array = make([]uint, len(value))
		for k, v := range value {
			array[k] = uint(v)
		}
	case []int16:
		array = make([]uint, len(value))
		for k, v := range value {
			array[k] = uint(v)
		}
	case []int32:
		array = make([]uint, len(value))
		for k, v := range value {
			array[k] = uint(v)
		}
	case []int64:
		array = make([]uint, len(value))
		for k, v := range value {
			array[k] = uint(v)
		}
	case []uint:
		array = value
	case []uint8:
		if json.Valid(value) {
			_ = json.UnmarshalUseNumber(value, &array)
		} else {
			array = make([]uint, len(value))
			for k, v := range value {
				array[k] = uint(v)
			}
		}
	case []uint16:
		array = make([]uint, len(value))
		for k, v := range value {
			array[k] = uint(v)
		}
	case []uint32:
		array = make([]uint, len(value))
		for k, v := range value {
			array[k] = uint(v)
		}
	case []uint64:
		array = make([]uint, len(value))
		for k, v := range value {
			array[k] = uint(v)
		}
	case []bool:
		array = make([]uint, len(value))
		for k, v := range value {
			if v {
				array[k] = 1
			} else {
				array[k] = 0
			}
		}
	case []float32:
		array = make([]uint, len(value))
		for k, v := range value {
			array[k] = Uint(v)
		}
	case []float64:
		array = make([]uint, len(value))
		for k, v := range value {
			array[k] = Uint(v)
		}
	case []interface{}:
		array = make([]uint, len(value))
		for k, v := range value {
			array[k] = Uint(v)
		}
	case [][]byte:
		array = make([]uint, len(value))
		for k, v := range value {
			array[k] = Uint(v)
		}
	}

	if array != nil {
		return array
	}

	// Default handler.
	if v, ok := any.(localinterface.IUints); ok {
		return v.Uints()
	}
	if v, ok := any.(localinterface.IInterfaces); ok {
		return Uints(v.Interfaces())
	}
	// JSON format string value converting.
	if checkJsonAndUnmarshalUseNumber(any, &array) {
		return array
	}
	// Not a common type, it then uses reflection for conversion.
	originValueAndKind := reflection.OriginValueAndKind(any)
	switch originValueAndKind.OriginKind {
	case reflect.Slice, reflect.Array:
		var (
			length = originValueAndKind.OriginValue.Len()
			slice  = make([]uint, length)
		)
		for i := 0; i < length; i++ {
			slice[i] = Uint(originValueAndKind.OriginValue.Index(i).Interface())
		}
		return slice

	default:
		if originValueAndKind.OriginValue.IsZero() {
			return []uint{}
		}
		return []uint{Uint(any)}
	}
}

// Uint32s converts `any` to []uint32.
func Uint32s(any interface{}) []uint32 {
	if any == nil {
		return nil
	}
	var (
		array []uint32 = nil
	)
	switch value := any.(type) {
	case string:
		value = strings.TrimSpace(value)
		if value == "" {
			return []uint32{}
		}
		if utils.IsNumeric(value) {
			return []uint32{Uint32(value)}
		}
	case []string:
		array = make([]uint32, len(value))
		for k, v := range value {
			array[k] = Uint32(v)
		}
	case []int8:
		array = make([]uint32, len(value))
		for k, v := range value {
			array[k] = uint32(v)
		}
	case []int16:
		array = make([]uint32, len(value))
		for k, v := range value {
			array[k] = uint32(v)
		}
	case []int32:
		array = make([]uint32, len(value))
		for k, v := range value {
			array[k] = uint32(v)
		}
	case []int64:
		array = make([]uint32, len(value))
		for k, v := range value {
			array[k] = uint32(v)
		}
	case []uint:
		array = make([]uint32, len(value))
		for k, v := range value {
			array[k] = uint32(v)
		}
	case []uint8:
		if json.Valid(value) {
			_ = json.UnmarshalUseNumber(value, &array)
		} else {
			array = make([]uint32, len(value))
			for k, v := range value {
				array[k] = uint32(v)
			}
		}
	case []uint16:
		array = make([]uint32, len(value))
		for k, v := range value {
			array[k] = uint32(v)
		}
	case []uint32:
		array = value
	case []uint64:
		array = make([]uint32, len(value))
		for k, v := range value {
			array[k] = uint32(v)
		}
	case []bool:
		array = make([]uint32, len(value))
		for k, v := range value {
			if v {
				array[k] = 1
			} else {
				array[k] = 0
			}
		}
	case []float32:
		array = make([]uint32, len(value))
		for k, v := range value {
			array[k] = Uint32(v)
		}
	case []float64:
		array = make([]uint32, len(value))
		for k, v := range value {
			array[k] = Uint32(v)
		}
	case []interface{}:
		array = make([]uint32, len(value))
		for k, v := range value {
			array[k] = Uint32(v)
		}
	case [][]byte:
		array = make([]uint32, len(value))
		for k, v := range value {
			array[k] = Uint32(v)
		}
	}
	if array != nil {
		return array
	}

	// Default handler.
	if v, ok := any.(localinterface.IUints); ok {
		return Uint32s(v.Uints())
	}
	if v, ok := any.(localinterface.IInterfaces); ok {
		return Uint32s(v.Interfaces())
	}
	// JSON format string value converting.
	if checkJsonAndUnmarshalUseNumber(any, &array) {
		return array
	}
	// Not a common type, it then uses reflection for conversion.
	originValueAndKind := reflection.OriginValueAndKind(any)
	switch originValueAndKind.OriginKind {
	case reflect.Slice, reflect.Array:
		var (
			length = originValueAndKind.OriginValue.Len()
			slice  = make([]uint32, length)
		)
		for i := 0; i < length; i++ {
			slice[i] = Uint32(originValueAndKind.OriginValue.Index(i).Interface())
		}
		return slice

	default:
		if originValueAndKind.OriginValue.IsZero() {
			return []uint32{}
		}
		return []uint32{Uint32(any)}
	}
}

// Uint64s converts `any` to []uint64.
func Uint64s(any interface{}) []uint64 {
	if any == nil {
		return nil
	}
	var (
		array []uint64 = nil
	)
	switch value := any.(type) {
	case string:
		value = strings.TrimSpace(value)
		if value == "" {
			return []uint64{}
		}
		if utils.IsNumeric(value) {
			return []uint64{Uint64(value)}
		}

	case []string:
		array = make([]uint64, len(value))
		for k, v := range value {
			array[k] = Uint64(v)
		}
	case []int8:
		array = make([]uint64, len(value))
		for k, v := range value {
			array[k] = uint64(v)
		}
	case []int16:
		array = make([]uint64, len(value))
		for k, v := range value {
			array[k] = uint64(v)
		}
	case []int32:
		array = make([]uint64, len(value))
		for k, v := range value {
			array[k] = uint64(v)
		}
	case []int64:
		array = make([]uint64, len(value))
		for k, v := range value {
			array[k] = uint64(v)
		}
	case []uint:
		array = make([]uint64, len(value))
		for k, v := range value {
			array[k] = uint64(v)
		}
	case []uint8:
		if json.Valid(value) {
			_ = json.UnmarshalUseNumber(value, &array)
		} else {
			array = make([]uint64, len(value))
			for k, v := range value {
				array[k] = uint64(v)
			}
		}
	case []uint16:
		array = make([]uint64, len(value))
		for k, v := range value {
			array[k] = uint64(v)
		}
	case []uint32:
		array = make([]uint64, len(value))
		for k, v := range value {
			array[k] = uint64(v)
		}
	case []uint64:
		array = value
	case []bool:
		array = make([]uint64, len(value))
		for k, v := range value {
			if v {
				array[k] = 1
			} else {
				array[k] = 0
			}
		}
	case []float32:
		array = make([]uint64, len(value))
		for k, v := range value {
			array[k] = Uint64(v)
		}
	case []float64:
		array = make([]uint64, len(value))
		for k, v := range value {
			array[k] = Uint64(v)
		}
	case []interface{}:
		array = make([]uint64, len(value))
		for k, v := range value {
			array[k] = Uint64(v)
		}
	case [][]byte:
		array = make([]uint64, len(value))
		for k, v := range value {
			array[k] = Uint64(v)
		}
	}
	if array != nil {
		return array
	}
	// Default handler.
	if v, ok := any.(localinterface.IUints); ok {
		return Uint64s(v.Uints())
	}
	if v, ok := any.(localinterface.IInterfaces); ok {
		return Uint64s(v.Interfaces())
	}
	// JSON format string value converting.
	if checkJsonAndUnmarshalUseNumber(any, &array) {
		return array
	}
	// Not a common type, it then uses reflection for conversion.
	originValueAndKind := reflection.OriginValueAndKind(any)
	switch originValueAndKind.OriginKind {
	case reflect.Slice, reflect.Array:
		var (
			length = originValueAndKind.OriginValue.Len()
			slice  = make([]uint64, length)
		)
		for i := 0; i < length; i++ {
			slice[i] = Uint64(originValueAndKind.OriginValue.Index(i).Interface())
		}
		return slice

	default:
		if originValueAndKind.OriginValue.IsZero() {
			return []uint64{}
		}
		return []uint64{Uint64(any)}
	}
}
