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

// Float32 converts `any` to float32.
func Float32(any any) float32 {
	v, _ := doFloat32(any)
	return v
}

func doFloat32(any any) (float32, error) {
	if any == nil {
		return 0, nil
	}
	switch value := any.(type) {
	case float32:
		return value, nil
	case float64:
		return float32(value), nil
	case []byte:
		// TODO: It might panic here for these types.
		return gbinary.DecodeToFloat32(value), nil
	default:
		rv := reflect.ValueOf(any)
		switch rv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return float32(rv.Int()), nil
		case reflect.Uintptr, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return float32(rv.Uint()), nil
		case reflect.Float32, reflect.Float64:
			return float32(rv.Float()), nil
		case reflect.Bool:
			if rv.Bool() {
				return 1, nil
			}
			return 0, nil
		case reflect.String:
			f, err := strconv.ParseFloat(rv.String(), 32)
			return float32(f), gerror.WrapCodef(
				gcode.CodeInvalidParameter, err, "converting string to float32 failed for: %v", any,
			)
		case reflect.Ptr:
			if rv.IsNil() {
				return 0, nil
			}
			if f, ok := value.(localinterface.IFloat32); ok {
				return f.Float32(), nil
			}
			return doFloat32(rv.Elem().Interface())
		default:
			if f, ok := value.(localinterface.IFloat32); ok {
				return f.Float32(), nil
			}
			v, err := strconv.ParseFloat(String(any), 32)
			return float32(v), gerror.WrapCodef(
				gcode.CodeInvalidParameter, err, "converting string to float32 failed for: %v", any,
			)
		}
	}
}

// Float64 converts `any` to float64.
func Float64(any any) float64 {
	v, _ := doFloat64(any)
	return v
}

func doFloat64(any any) (float64, error) {
	if any == nil {
		return 0, nil
	}
	switch value := any.(type) {
	case float32:
		return float64(value), nil
	case float64:
		return value, nil
	case []byte:
		// TODO: It might panic here for these types.
		return gbinary.DecodeToFloat64(value), nil
	default:
		rv := reflect.ValueOf(any)
		switch rv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return float64(rv.Int()), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return float64(rv.Uint()), nil
		case reflect.Uintptr:
			return float64(rv.Uint()), nil
		case reflect.Float32, reflect.Float64:
			// Please Note:
			// When the type is float32 or a custom type defined based on float32,
			// switching to float64 may result in a few extra decimal places.
			return rv.Float(), nil
		case reflect.Bool:
			if rv.Bool() {
				return 1, nil
			}
			return 0, nil
		case reflect.String:
			f, err := strconv.ParseFloat(rv.String(), 64)
			return f, gerror.WrapCodef(
				gcode.CodeInvalidParameter, err, "converting string to float64 failed for: %v", any,
			)
		case reflect.Ptr:
			if rv.IsNil() {
				return 0, nil
			}
			if f, ok := value.(localinterface.IFloat64); ok {
				return f.Float64(), nil
			}
			return doFloat64(rv.Elem().Interface())
		default:
			if f, ok := value.(localinterface.IFloat64); ok {
				return f.Float64(), nil
			}
			v, err := strconv.ParseFloat(String(any), 64)
			return v, gerror.WrapCodef(
				gcode.CodeInvalidParameter, err, "converting string to float64 failed for: %v", any,
			)
		}
	}
}
