// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package converter

import (
	"reflect"
	"strings"

	"github.com/gogf/gf/v2/internal/empty"
	"github.com/gogf/gf/v2/util/gconv/internal/localinterface"
)

// Bool converts `any` to bool.
func (c *Converter) Bool(any any) (bool, error) {
	if empty.IsNil(any) {
		return false, nil
	}
	switch value := any.(type) {
	case bool:
		return value, nil
	case []byte:
		if _, ok := emptyStringMap[strings.ToLower(string(value))]; ok {
			return false, nil
		}
		return true, nil
	case string:
		if _, ok := emptyStringMap[strings.ToLower(value)]; ok {
			return false, nil
		}
		return true, nil
	default:
		if f, ok := value.(localinterface.IBool); ok {
			return f.Bool(), nil
		}
		rv := reflect.ValueOf(any)
		switch rv.Kind() {
		case reflect.Ptr:
			if rv.IsNil() {
				return false, nil
			}
			if rv.Type().Elem().Kind() == reflect.Bool {
				return rv.Elem().Bool(), nil
			}
			return c.Bool(rv.Elem().Interface())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return rv.Int() != 0, nil
		case reflect.Uintptr, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return rv.Uint() != 0, nil
		case reflect.Float32, reflect.Float64:
			return rv.Float() != 0, nil
		case reflect.Bool:
			return rv.Bool(), nil
		// TODO：(Map，Array，Slice，Struct) It might panic here for these types.
		case reflect.Map, reflect.Array:
			fallthrough
		case reflect.Slice:
			return rv.Len() != 0, nil
		case reflect.Struct:
			return true, nil
		default:
			s, err := c.String(any)
			if err != nil {
				return false, err
			}
			if _, ok := emptyStringMap[strings.ToLower(s)]; ok {
				return false, nil
			}
			return true, nil
		}
	}
}
