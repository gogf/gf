// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gconv implements powerful and convenient converting functionality for any types of variables.
//
// This package should keep much fewer dependencies with other packages.
package gconv

import (
	"context"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/encoding/gbinary"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/internal/reflection"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv/internal/localinterface"
	"github.com/gogf/gf/v2/util/gconv/internal/structcache"
)

var (
	// Empty strings.
	emptyStringMap = map[string]struct{}{
		"":      {},
		"0":     {},
		"no":    {},
		"off":   {},
		"false": {},
	}
)

// IUnmarshalValue is the interface for custom defined types customizing value assignment.
// Note that only pointer can implement interface IUnmarshalValue.
type IUnmarshalValue = localinterface.IUnmarshalValue

func init() {
	// register common converters for internal usage.
	structcache.RegisterCommonConverter(structcache.CommonConverter{
		Int64:   Int64,
		Uint64:  Uint64,
		String:  String,
		Float32: Float32,
		Float64: Float64,
		Time:    Time,
		GTime:   GTime,
		Bytes:   Bytes,
		Bool:    Bool,
	})
}

// Byte converts `any` to byte.
func Byte(any any) byte {
	if v, ok := any.(byte); ok {
		return v
	}
	return Uint8(any)
}

// Bytes converts `any` to []byte.
func Bytes(any any) []byte {
	if any == nil {
		return nil
	}
	switch value := any.(type) {
	case string:
		return []byte(value)

	case []byte:
		return value

	default:
		if f, ok := value.(localinterface.IBytes); ok {
			return f.Bytes()
		}
		originValueAndKind := reflection.OriginValueAndKind(any)
		switch originValueAndKind.OriginKind {
		case reflect.Map:
			bytes, err := json.Marshal(any)
			if err != nil {
				intlog.Errorf(context.TODO(), `%+v`, err)
			}
			return bytes

		case reflect.Array, reflect.Slice:
			var (
				ok    = true
				bytes = make([]byte, originValueAndKind.OriginValue.Len())
			)
			for i := range bytes {
				int32Value := Int32(originValueAndKind.OriginValue.Index(i).Interface())
				if int32Value < 0 || int32Value > math.MaxUint8 {
					ok = false
					break
				}
				bytes[i] = byte(int32Value)
			}
			if ok {
				return bytes
			}
		}
		return gbinary.Encode(any)
	}
}

// Rune converts `any` to rune.
func Rune(any any) rune {
	if v, ok := any.(rune); ok {
		return v
	}
	return Int32(any)
}

// Runes converts `any` to []rune.
func Runes(any any) []rune {
	if v, ok := any.([]rune); ok {
		return v
	}
	return []rune(String(any))
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
	if any == nil {
		return false
	}
	switch value := any.(type) {
	case bool:
		return value
	case []byte:
		if _, ok := emptyStringMap[strings.ToLower(string(value))]; ok {
			return false
		}
		return true
	case string:
		if _, ok := emptyStringMap[strings.ToLower(value)]; ok {
			return false
		}
		return true
	default:
		if f, ok := value.(localinterface.IBool); ok {
			return f.Bool()
		}
		rv := reflect.ValueOf(any)
		switch rv.Kind() {
		case reflect.Ptr:
			if rv.IsNil() {
				return false
			}
			if rv.Type().Elem().Kind() == reflect.Bool {
				return rv.Elem().Bool()
			}
			return Bool(rv.Elem().Interface())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return rv.Int() != 0
		case reflect.Uintptr, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return rv.Uint() != 0
		case reflect.Float32, reflect.Float64:
			return rv.Float() != 0
		case reflect.Bool:
			return rv.Bool()
		// TODO：(Map，Array，Slice，Struct) It might panic here for these types.
		case reflect.Map, reflect.Array:
			fallthrough
		case reflect.Slice:
			return rv.Len() != 0
		case reflect.Struct:
			return true
		default:
			s := strings.ToLower(String(any))
			if _, ok := emptyStringMap[s]; ok {
				return false
			}
			return true
		}
	}
}

// checkJsonAndUnmarshalUseNumber checks if given `any` is JSON formatted string value and does converting using `json.UnmarshalUseNumber`.
func checkJsonAndUnmarshalUseNumber(any any, target any) bool {
	switch r := any.(type) {
	case []byte:
		if json.Valid(r) {
			if err := json.UnmarshalUseNumber(r, &target); err != nil {
				return false
			}
			return true
		}

	case string:
		anyAsBytes := []byte(r)
		if json.Valid(anyAsBytes) {
			if err := json.UnmarshalUseNumber(anyAsBytes, &target); err != nil {
				return false
			}
			return true
		}
	}
	return false
}
