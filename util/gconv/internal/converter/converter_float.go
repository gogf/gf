// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package converter

import (
	"reflect"
	"strconv"

	"github.com/gogf/gf/v2/encoding/gbinary"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/empty"
	"github.com/gogf/gf/v2/util/gconv/internal/localinterface"
)

// Float32 converts `any` to float32.
func (c *Converter) Float32(any any) (float32, error) {
	if empty.IsNil(any) {
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
			if err != nil {
				return 0, gerror.WrapCodef(
					gcode.CodeInvalidParameter, err, "converting string to float32 failed for: %v", any,
				)
			}
			return float32(f), nil
		case reflect.Ptr:
			if rv.IsNil() {
				return 0, nil
			}
			if f, ok := value.(localinterface.IFloat32); ok {
				return f.Float32(), nil
			}
			return c.Float32(rv.Elem().Interface())
		default:
			if f, ok := value.(localinterface.IFloat32); ok {
				return f.Float32(), nil
			}
			s, err := c.String(any)
			if err != nil {
				return 0, err
			}
			v, err := strconv.ParseFloat(s, 32)
			if err != nil {
				return 0, gerror.WrapCodef(
					gcode.CodeInvalidParameter, err, "converting string to float32 failed for: %v", any,
				)
			}
			return float32(v), nil
		}
	}
}

// Float64 converts `any` to float64.
func (c *Converter) Float64(any any) (float64, error) {
	if empty.IsNil(any) {
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
			if err != nil {
				return 0, gerror.WrapCodef(
					gcode.CodeInvalidParameter, err, "converting string to float64 failed for: %v", any,
				)
			}
			return f, nil
		case reflect.Ptr:
			if rv.IsNil() {
				return 0, nil
			}
			if f, ok := value.(localinterface.IFloat64); ok {
				return f.Float64(), nil
			}
			return c.Float64(rv.Elem().Interface())
		default:
			if f, ok := value.(localinterface.IFloat64); ok {
				return f.Float64(), nil
			}
			s, err := c.String(any)
			if err != nil {
				return 0, err
			}
			v, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return 0, gerror.WrapCodef(
					gcode.CodeInvalidParameter, err, "converting string to float64 failed for: %v", any,
				)
			}
			return v, nil
		}
	}
}
