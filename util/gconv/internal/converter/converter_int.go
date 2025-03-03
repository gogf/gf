// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package converter

import (
	"math"
	"reflect"
	"strconv"

	"github.com/gogf/gf/v2/encoding/gbinary"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/empty"
	"github.com/gogf/gf/v2/util/gconv/internal/localinterface"
)

// Int converts `any` to int.
func (c *Converter) Int(any any) (int, error) {
	if v, ok := any.(int); ok {
		return v, nil
	}
	v, err := c.Int64(any)
	if err != nil {
		return 0, err
	}
	return int(v), nil
}

// Int8 converts `any` to int8.
func (c *Converter) Int8(any any) (int8, error) {
	if v, ok := any.(int8); ok {
		return v, nil
	}
	v, err := c.Int64(any)
	if err != nil {
		return 0, err
	}
	return int8(v), nil
}

// Int16 converts `any` to int16.
func (c *Converter) Int16(any any) (int16, error) {
	if v, ok := any.(int16); ok {
		return v, nil
	}
	v, err := c.Int64(any)
	if err != nil {
		return 0, err
	}
	return int16(v), nil
}

// Int32 converts `any` to int32.
func (c *Converter) Int32(any any) (int32, error) {
	if v, ok := any.(int32); ok {
		return v, nil
	}
	v, err := c.Int64(any)
	if err != nil {
		return 0, err
	}
	return int32(v), nil
}

// Int64 converts `any` to int64.
func (c *Converter) Int64(any any) (int64, error) {
	if empty.IsNil(any) {
		return 0, nil
	}
	if v, ok := any.(int64); ok {
		return v, nil
	}
	rv := reflect.ValueOf(any)
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return rv.Int(), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int64(rv.Uint()), nil
	case reflect.Uintptr:
		return int64(rv.Uint()), nil
	case reflect.Float32, reflect.Float64:
		return int64(rv.Float()), nil
	case reflect.Bool:
		if rv.Bool() {
			return 1, nil
		}
		return 0, nil
	case reflect.Ptr:
		if rv.IsNil() {
			return 0, nil
		}
		if f, ok := any.(localinterface.IInt64); ok {
			return f.Int64(), nil
		}
		return c.Int64(rv.Elem().Interface())
	case reflect.Slice:
		// TODO: It might panic here for these types.
		if rv.Type().Elem().Kind() == reflect.Uint8 {
			return gbinary.DecodeToInt64(rv.Bytes()), nil
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
					return -v, nil
				}
				return v, nil
			}
		}
		// Decimal.
		if v, e := strconv.ParseInt(s, 10, 64); e == nil {
			if isMinus {
				return -v, nil
			}
			return v, nil
		}
		// Float64.
		valueInt64, err := c.Float64(s)
		if err != nil {
			return 0, err
		}
		if math.IsNaN(valueInt64) {
			return 0, nil
		} else {
			if isMinus {
				return -int64(valueInt64), nil
			}
			return int64(valueInt64), nil
		}
	default:
		if f, ok := any.(localinterface.IInt64); ok {
			return f.Int64(), nil
		}
	}
	return 0, gerror.NewCodef(
		gcode.CodeInvalidParameter,
		`unsupport value type for converting to int64: %v`,
		reflect.TypeOf(any),
	)
}
