// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv

import (
	"reflect"
	"strconv"

	"github.com/gogf/gf/v2/encoding/gbinary"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/util/gconv/internal/localinterface"
)

// Uint converts `any` to uint.
func Uint(any any) uint {
	v, _ := doUint(any)
	return v
}

func doUint(any any) (uint, error) {
	if v, ok := any.(uint); ok {
		return v, nil
	}
	v, err := doUint64(any)
	if err != nil {
		return 0, err
	}
	return uint(v), nil
}

// Uint8 converts `any` to uint8.
func Uint8(any any) uint8 {
	v, _ := doUint8(any)
	return v
}

func doUint8(any any) (uint8, error) {
	if v, ok := any.(uint8); ok {
		return v, nil
	}
	v, err := doUint64(any)
	if err != nil {
		return 0, err
	}
	return uint8(v), nil
}

// Uint16 converts `any` to uint16.
func Uint16(any any) uint16 {
	v, _ := doUint16(any)
	return v
}

func doUint16(any any) (uint16, error) {
	if v, ok := any.(uint16); ok {
		return v, nil
	}
	v, err := doUint64(any)
	if err != nil {
		return 0, err
	}
	return uint16(v), nil
}

// Uint32 converts `any` to uint32.
func Uint32(any any) uint32 {
	v, _ := doUint32(any)
	return v
}

func doUint32(any any) (uint32, error) {
	if v, ok := any.(uint32); ok {
		return v, nil
	}
	v, err := doUint64(any)
	if err != nil {
		return 0, err
	}
	return uint32(v), nil
}

// Uint64 converts `any` to uint64.
func Uint64(any any) uint64 {
	v, _ := doUint64(any)
	return v
}

func doUint64(any any) (uint64, error) {
	if any == nil {
		return 0, nil
	}
	if v, ok := any.(uint64); ok {
		return v, nil
	}
	rv := reflect.ValueOf(any)
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val := rv.Int()
		if val < 0 {
			return 0, gerror.NewCodef(
				gcode.CodeInvalidParameter,
				`cannot convert negative value "%d" to uint64`,
				val,
			)
		}
		return uint64(val), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return rv.Uint(), nil
	case reflect.Uintptr:
		return rv.Uint(), nil
	case reflect.Float32, reflect.Float64:
		val := rv.Float()
		if val < 0 {
			return 0, gerror.NewCodef(
				gcode.CodeInvalidParameter,
				`cannot convert negative value "%f" to uint64`,
				val,
			)
		}
		return uint64(val), nil
	case reflect.Bool:
		if rv.Bool() {
			return 1, nil
		}
		return 0, nil
	case reflect.Ptr:
		if rv.IsNil() {
			return 0, nil
		}
		if f, ok := any.(localinterface.IUint64); ok {
			return f.Uint64(), nil
		}
		return doUint64(rv.Elem().Interface())
	case reflect.Slice:
		if rv.Type().Elem().Kind() == reflect.Uint8 {
			return gbinary.DecodeToUint64(rv.Bytes()), nil
		}
		return 0, gerror.NewCodef(
			gcode.CodeInvalidParameter,
			`unsupport slice type "%s" for converting to uint64`,
			rv.Type().String(),
		)
	case reflect.String:
		var s = rv.String()
		// Hexadecimal
		if len(s) > 2 && s[0] == '0' && (s[1] == 'x' || s[1] == 'X') {
			v, err := strconv.ParseUint(s[2:], 16, 64)
			if err == nil {
				return v, nil
			}
			return 0, gerror.WrapCodef(
				gcode.CodeInvalidParameter,
				err,
				`cannot convert hexadecimal string "%s" to uint64`,
				s,
			)
		}
		// Decimal
		if v, err := strconv.ParseUint(s, 10, 64); err == nil {
			return v, nil
		}
		// Float64
		if v, err := doFloat64(any); err == nil {
			return uint64(v), nil
		}
	default:
		if f, ok := any.(localinterface.IUint64); ok {
			return f.Uint64(), nil
		}
	}
	return 0, gerror.NewCodef(
		gcode.CodeInvalidParameter,
		`unsupport value type "%s" for converting to uint64`,
		reflect.TypeOf(any).String(),
	)
}
