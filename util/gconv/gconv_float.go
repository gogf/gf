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
	"github.com/gogf/gf/v2/util/gconv/internal/localinterface"
)

// Float32 converts `any` to float32.
func Float32(any interface{}) float32 {
	if any == nil {
		return 0
	}
	switch value := any.(type) {
	case float32:
		return value
	case float64:
		return float32(value)
	case []byte:
		return gbinary.DecodeToFloat32(value)
	default:
		if f, ok := value.(localinterface.IFloat32); ok {
			return f.Float32()
		}
		rv := reflect.ValueOf(any)
		switch rv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return float32(rv.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return float32(rv.Uint())
		case reflect.Float32, reflect.Float64:
			return float32(rv.Float())
		case reflect.Bool:
			if rv.Bool() {
				return 1
			}
			return 0
		case reflect.String:
			f, _ := strconv.ParseFloat(rv.String(), 32)
			return float32(f)
		case reflect.Ptr:
			return Float32(rv.Elem().Interface())
		default:
			v, _ := strconv.ParseFloat(String(any), 64)
			return float32(v)
		}

	}
}

// Float64 converts `any` to float64.
func Float64(any interface{}) float64 {
	if any == nil {
		return 0
	}
	switch value := any.(type) {
	case float32:
		return float64(value)
	case float64:
		return value
	case []byte:
		return gbinary.DecodeToFloat64(value)
	default:
		if f, ok := value.(localinterface.IFloat64); ok {
			return f.Float64()
		}
		rv := reflect.ValueOf(any)
		switch rv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return float64(rv.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return float64(rv.Uint())
		case reflect.Float32, reflect.Float64:
			return rv.Float()
		case reflect.Bool:
			if rv.Bool() {
				return 1
			}
			return 0
		case reflect.String:
			f, _ := strconv.ParseFloat(rv.String(), 64)
			return f
		case reflect.Ptr:
			return Float64(rv.Elem().Interface())
		default:
			v, _ := strconv.ParseFloat(String(any), 64)
			return v
		}
	}
}
