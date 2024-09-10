// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv

import (
	"math"
	"reflect"
	"strconv"

	"github.com/gogf/gf/v2/encoding/gbinary"
	"github.com/gogf/gf/v2/util/gconv/internal/localinterface"
)

// Uint converts `any` to uint.
func Uint(any interface{}) uint {
	if any == nil {
		return 0
	}
	if v, ok := any.(uint); ok {
		return v
	}
	return uint(Uint64(any))
}

// Uint8 converts `any` to uint8.
func Uint8(any interface{}) uint8 {
	if any == nil {
		return 0
	}
	if v, ok := any.(uint8); ok {
		return v
	}
	return uint8(Uint64(any))
}

// Uint16 converts `any` to uint16.
func Uint16(any interface{}) uint16 {
	if any == nil {
		return 0
	}
	if v, ok := any.(uint16); ok {
		return v
	}
	return uint16(Uint64(any))
}

// Uint32 converts `any` to uint32.
func Uint32(any interface{}) uint32 {
	if any == nil {
		return 0
	}
	if v, ok := any.(uint32); ok {
		return v
	}
	return uint32(Uint64(any))
}

// Uint64 converts `any` to uint64.
func Uint64(any interface{}) uint64 {
	if any == nil {
		return 0
	}
	rv := reflect.ValueOf(any)
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return uint64(rv.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return uint64(rv.Uint())
	case reflect.Uintptr:
		return uint64(rv.Uint())
	case reflect.Float32, reflect.Float64:
		return uint64(rv.Float())
	case reflect.Bool:
		if rv.Bool() {
			return 1
		}
		return 0
	case reflect.Ptr:
		if rv.IsNil() {
			return 0
		}
		if f, ok := any.(localinterface.IUint64); ok {
			return f.Uint64()
		}
		return Uint64(rv.Elem().Interface())
	case reflect.Slice:
		// TODO：(Map，Array，Slice，Struct) These types should be panic
		if rv.Type().Elem().Kind() == reflect.Uint8 {
			return gbinary.DecodeToUint64(rv.Bytes())
		}
	case reflect.String:
		var (
			s = rv.String()
		)
		// Hexadecimal
		if len(s) > 2 && s[0] == '0' && (s[1] == 'x' || s[1] == 'X') {
			if v, e := strconv.ParseUint(s[2:], 16, 64); e == nil {
				return v
			}
		}
		// Decimal
		if v, e := strconv.ParseUint(s, 10, 64); e == nil {
			return v
		}
		// Float64
		if valueFloat64 := Float64(any); math.IsNaN(valueFloat64) {
			return 0
		} else {
			return uint64(valueFloat64)
		}
	default:
		if f, ok := any.(localinterface.IUint64); ok {
			return f.Uint64()
		}
	}
	return 0
}
