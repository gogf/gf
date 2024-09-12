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

// Int converts `any` to int.
func Int(any interface{}) int {
	if any == nil {
		return 0
	}
	if v, ok := any.(int); ok {
		return v
	}
	return int(Int64(any))
}

// Int8 converts `any` to int8.
func Int8(any interface{}) int8 {
	if any == nil {
		return 0
	}
	if v, ok := any.(int8); ok {
		return v
	}
	return int8(Int64(any))
}

// Int16 converts `any` to int16.
func Int16(any interface{}) int16 {
	if any == nil {
		return 0
	}
	if v, ok := any.(int16); ok {
		return v
	}
	return int16(Int64(any))
}

// Int32 converts `any` to int32.
func Int32(any interface{}) int32 {
	if any == nil {
		return 0
	}
	if v, ok := any.(int32); ok {
		return v
	}
	return int32(Int64(any))
}

// Int64 converts `any` to int64.
func Int64(any interface{}) int64 {
	if any == nil {
		return 0
	}
	rv := reflect.ValueOf(any)
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return int64(rv.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int64(rv.Uint())
	case reflect.Uintptr:
		return int64(rv.Uint())
	case reflect.Float32, reflect.Float64:
		return int64(rv.Float())
	case reflect.Bool:
		if rv.Bool() {
			return 1
		}
		return 0
	case reflect.Ptr:
		if rv.IsNil() {
			return 0
		}
		if f, ok := any.(localinterface.IInt64); ok {
			return f.Int64()
		}
		return Int64(rv.Elem().Interface())
	case reflect.Slice:
		// TODO: It might panic here for these types.
		if rv.Type().Elem().Kind() == reflect.Uint8 {
			return gbinary.DecodeToInt64(rv.Bytes())
		}
	case reflect.String:
		var (
			s       = rv.String()
			isMinus = false
		)
		if len(s) > 0 {
			if s[0] == '-' {
				isMinus = true
				s = s[1:]
			} else if s[0] == '+' {
				s = s[1:]
			}
		}
		// Hexadecimal.
		if len(s) > 2 && s[0] == '0' && (s[1] == 'x' || s[1] == 'X') {
			if v, e := strconv.ParseInt(s[2:], 16, 64); e == nil {
				if isMinus {
					return -v
				}
				return v
			}
		}
		// Decimal.
		if v, e := strconv.ParseInt(s, 10, 64); e == nil {
			if isMinus {
				return -v
			}
			return v
		}
		// Float64.
		if valueInt64 := Float64(s); math.IsNaN(valueInt64) {
			return 0
		} else {
			return int64(valueInt64)
		}
	default:
		if f, ok := any.(localinterface.IInt64); ok {
			return f.Int64()
		}
	}
	return 0
}
