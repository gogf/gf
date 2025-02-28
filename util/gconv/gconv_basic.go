// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/encoding/gbinary"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/internal/reflection"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv/internal/localinterface"
)

// Byte converts `any` to byte.
func Byte(any any) byte {
	v, _ := doByte(any)
	return v
}

func doByte(any any) (byte, error) {
	if v, ok := any.(byte); ok {
		return v, nil
	}
	return doUint8(any)
}

// Bytes converts `any` to []byte.
func Bytes(any any) []byte {
	v, _ := doBytes(any)
	return v
}

func doBytes(any any) ([]byte, error) {
	if any == nil {
		return nil, nil
	}
	switch value := any.(type) {
	case string:
		return []byte(value), nil

	case []byte:
		return value, nil

	default:
		if f, ok := value.(localinterface.IBytes); ok {
			return f.Bytes(), nil
		}
		originValueAndKind := reflection.OriginValueAndKind(any)
		switch originValueAndKind.OriginKind {
		case reflect.Map:
			bytes, err := json.Marshal(any)
			if err != nil {
				return nil, err
			}
			return bytes, nil

		case reflect.Array, reflect.Slice:
			var (
				ok    = true
				bytes = make([]byte, originValueAndKind.OriginValue.Len())
			)
			for i := range bytes {
				int32Value, err := doInt32(originValueAndKind.OriginValue.Index(i).Interface())
				if err != nil {
					return nil, err
				}
				if int32Value < 0 || int32Value > math.MaxUint8 {
					ok = false
					break
				}
				bytes[i] = byte(int32Value)
			}
			if ok {
				return bytes, nil
			}
		default:
		}
		return gbinary.Encode(any), nil
	}
}

// Rune converts `any` to rune.
func Rune(any any) rune {
	v, _ := doRune(any)
	return v
}

func doRune(any any) (rune, error) {
	if v, ok := any.(rune); ok {
		return v, nil
	}
	v, err := doInt32(any)
	if err != nil {
		return 0, err
	}
	return v, nil
}

// Runes converts `any` to []rune.
func Runes(any any) []rune {
	v, _ := doRunes(any)
	return v
}

func doRunes(any any) ([]rune, error) {
	if v, ok := any.([]rune); ok {
		return v, nil
	}
	s, err := doString(any)
	if err != nil {
		return nil, err
	}
	return []rune(s), nil
}

// String converts `any` to string.
// It's most commonly used converting function.
func String(any any) string {
	v, _ := doString(any)
	return v
}

func doString(any any) (string, error) {
	if any == nil {
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
			return doString(rv.Elem().Interface())
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

// Bool converts `any` to bool.
// It returns false if `any` is: false, "", 0, "false", "off", "no", empty slice/map.
func Bool(any any) bool {
	v, _ := doBool(any)
	return v
}

func doBool(any any) (bool, error) {
	if any == nil {
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
			return doBool(rv.Elem().Interface())
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
			s, err := doString(any)
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
