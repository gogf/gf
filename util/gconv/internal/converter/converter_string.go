// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package converter

import (
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/empty"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv/internal/localinterface"
)

func (c *Converter) String(any any) (string, error) {
	if empty.IsNil(any) {
		return "", nil
	}
	switch value := any.(type) {
	case int:
		return strconv.Itoa(value), nil
	case int8:
		return strconv.Itoa(int(value)), nil
	case int16:
		return strconv.Itoa(int(value)), nil
	case int32:
		return strconv.Itoa(int(value)), nil
	case int64:
		return strconv.FormatInt(value, 10), nil
	case uint:
		return strconv.FormatUint(uint64(value), 10), nil
	case uint8:
		return strconv.FormatUint(uint64(value), 10), nil
	case uint16:
		return strconv.FormatUint(uint64(value), 10), nil
	case uint32:
		return strconv.FormatUint(uint64(value), 10), nil
	case uint64:
		return strconv.FormatUint(value, 10), nil
	case float32:
		return strconv.FormatFloat(float64(value), 'f', -1, 32), nil
	case float64:
		return strconv.FormatFloat(value, 'f', -1, 64), nil
	case bool:
		return strconv.FormatBool(value), nil
	case string:
		return value, nil
	case []byte:
		return string(value), nil
	case complex64, complex128:
		return fmt.Sprintf("%v", value), nil
	case time.Time:
		if value.IsZero() {
			return "", nil
		}
		return value.String(), nil
	case *time.Time:
		if value == nil {
			return "", nil
		}
		return value.String(), nil
	case gtime.Time:
		if value.IsZero() {
			return "", nil
		}
		return value.String(), nil
	case *gtime.Time:
		if value == nil {
			return "", nil
		}
		return value.String(), nil
	default:
		if f, ok := value.(localinterface.IString); ok {
			// If the variable implements the String() interface,
			// then use that interface to perform the conversion
			return f.String(), nil
		}
		if f, ok := value.(localinterface.IError); ok {
			// If the variable implements the Error() interface,
			// then use that interface to perform the conversion
			return f.Error(), nil
		}
		// Reflect checks.
		var (
			rv   = reflect.ValueOf(value)
			kind = rv.Kind()
		)
		switch kind {
		case
			reflect.Chan,
			reflect.Map,
			reflect.Slice,
			reflect.Func,
			reflect.Interface,
			reflect.UnsafePointer:
			if rv.IsNil() {
				return "", nil
			}
		case reflect.String:
			return rv.String(), nil
		case reflect.Ptr:
			if rv.IsNil() {
				return "", nil
			}
			return c.String(rv.Elem().Interface())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return strconv.FormatInt(rv.Int(), 10), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return strconv.FormatUint(rv.Uint(), 10), nil
		case reflect.Uintptr:
			return strconv.FormatUint(rv.Uint(), 10), nil
		case reflect.Float32, reflect.Float64:
			return strconv.FormatFloat(rv.Float(), 'f', -1, 64), nil
		case reflect.Bool:
			return strconv.FormatBool(rv.Bool()), nil
		default:
		}
		// Finally, we use json.Marshal to convert.
		jsonContent, err := json.Marshal(value)
		if err != nil {
			return fmt.Sprint(value), gerror.WrapCodef(
				gcode.CodeInvalidParameter, err, "error marshaling value to JSON for: %v", value,
			)
		}
		return string(jsonContent), nil
	}
}
