// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gconv implements powerful and convenient converting functionality for any types of variables.
//
// This package should keep much less dependencies with other packages.
package gconv

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/encoding/gbinary"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/internal/utils"
	"github.com/gogf/gf/v2/os/gtime"
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

	// StructTagPriority defines the default priority tags for Map*/Struct* functions.
	// Note, the "gconv", "param", "params" tags are used by old version of package.
	// It is strongly recommended using short tag "c" or "p" instead in the future.
	StructTagPriority = []string{"gconv", "param", "params", "c", "p", "json"}
)

// Convert converts the variable `fromValue` to the type `toTypeName`, the type `toTypeName` is specified by string.
// The optional parameter `extraParams` is used for additional necessary parameter for this conversion.
// It supports common types conversion as its conversion based on type name string.
func Convert(fromValue interface{}, toTypeName string, extraParams ...interface{}) interface{} {
	return doConvert(doConvertInput{
		FromValue:  fromValue,
		ToTypeName: toTypeName,
		ReferValue: nil,
		Extra:      extraParams,
	})
}

type doConvertInput struct {
	FromValue  interface{}   // Value that is converted from.
	ToTypeName string        // Target value type name in string.
	ReferValue interface{}   // Referred value, a value in type `ToTypeName`.
	Extra      []interface{} // Extra values for implementing the converting.
}

// doConvert does commonly use types converting.
func doConvert(in doConvertInput) interface{} {
	switch in.ToTypeName {
	case "int":
		return Int(in.FromValue)
	case "*int":
		if _, ok := in.FromValue.(*int); ok {
			return in.FromValue
		}
		v := Int(in.FromValue)
		return &v

	case "int8":
		return Int8(in.FromValue)
	case "*int8":
		if _, ok := in.FromValue.(*int8); ok {
			return in.FromValue
		}
		v := Int8(in.FromValue)
		return &v

	case "int16":
		return Int16(in.FromValue)
	case "*int16":
		if _, ok := in.FromValue.(*int16); ok {
			return in.FromValue
		}
		v := Int16(in.FromValue)
		return &v

	case "int32":
		return Int32(in.FromValue)
	case "*int32":
		if _, ok := in.FromValue.(*int32); ok {
			return in.FromValue
		}
		v := Int32(in.FromValue)
		return &v

	case "int64":
		return Int64(in.FromValue)
	case "*int64":
		if _, ok := in.FromValue.(*int64); ok {
			return in.FromValue
		}
		v := Int64(in.FromValue)
		return &v

	case "uint":
		return Uint(in.FromValue)
	case "*uint":
		if _, ok := in.FromValue.(*uint); ok {
			return in.FromValue
		}
		v := Uint(in.FromValue)
		return &v

	case "uint8":
		return Uint8(in.FromValue)
	case "*uint8":
		if _, ok := in.FromValue.(*uint8); ok {
			return in.FromValue
		}
		v := Uint8(in.FromValue)
		return &v

	case "uint16":
		return Uint16(in.FromValue)
	case "*uint16":
		if _, ok := in.FromValue.(*uint16); ok {
			return in.FromValue
		}
		v := Uint16(in.FromValue)
		return &v

	case "uint32":
		return Uint32(in.FromValue)
	case "*uint32":
		if _, ok := in.FromValue.(*uint32); ok {
			return in.FromValue
		}
		v := Uint32(in.FromValue)
		return &v

	case "uint64":
		return Uint64(in.FromValue)
	case "*uint64":
		if _, ok := in.FromValue.(*uint64); ok {
			return in.FromValue
		}
		v := Uint64(in.FromValue)
		return &v

	case "float32":
		return Float32(in.FromValue)
	case "*float32":
		if _, ok := in.FromValue.(*float32); ok {
			return in.FromValue
		}
		v := Float32(in.FromValue)
		return &v

	case "float64":
		return Float64(in.FromValue)
	case "*float64":
		if _, ok := in.FromValue.(*float64); ok {
			return in.FromValue
		}
		v := Float64(in.FromValue)
		return &v

	case "bool":
		return Bool(in.FromValue)
	case "*bool":
		if _, ok := in.FromValue.(*bool); ok {
			return in.FromValue
		}
		v := Bool(in.FromValue)
		return &v

	case "string":
		return String(in.FromValue)
	case "*string":
		if _, ok := in.FromValue.(*string); ok {
			return in.FromValue
		}
		v := String(in.FromValue)
		return &v

	case "[]byte":
		return Bytes(in.FromValue)
	case "[]int":
		return Ints(in.FromValue)
	case "[]int32":
		return Int32s(in.FromValue)
	case "[]int64":
		return Int64s(in.FromValue)
	case "[]uint":
		return Uints(in.FromValue)
	case "[]uint8":
		return Bytes(in.FromValue)
	case "[]uint32":
		return Uint32s(in.FromValue)
	case "[]uint64":
		return Uint64s(in.FromValue)
	case "[]float32":
		return Float32s(in.FromValue)
	case "[]float64":
		return Float64s(in.FromValue)
	case "[]string":
		return Strings(in.FromValue)

	case "Time", "time.Time":
		if len(in.Extra) > 0 {
			return Time(in.FromValue, String(in.Extra[0]))
		}
		return Time(in.FromValue)
	case "*time.Time":
		var v interface{}
		if len(in.Extra) > 0 {
			v = Time(in.FromValue, String(in.Extra[0]))
		} else {
			if _, ok := in.FromValue.(*time.Time); ok {
				return in.FromValue
			}
			v = Time(in.FromValue)
		}
		return &v

	case "GTime", "gtime.Time":
		if len(in.Extra) > 0 {
			if v := GTime(in.FromValue, String(in.Extra[0])); v != nil {
				return *v
			} else {
				return *gtime.New()
			}
		}
		if v := GTime(in.FromValue); v != nil {
			return *v
		} else {
			return *gtime.New()
		}
	case "*gtime.Time":
		if len(in.Extra) > 0 {
			if v := GTime(in.FromValue, String(in.Extra[0])); v != nil {
				return v
			} else {
				return gtime.New()
			}
		}
		if v := GTime(in.FromValue); v != nil {
			return v
		} else {
			return gtime.New()
		}

	case "Duration", "time.Duration":
		return Duration(in.FromValue)
	case "*time.Duration":
		if _, ok := in.FromValue.(*time.Duration); ok {
			return in.FromValue
		}
		v := Duration(in.FromValue)
		return &v

	case "map[string]string":
		return MapStrStr(in.FromValue)

	case "map[string]interface{}":
		return Map(in.FromValue)

	case "[]map[string]interface{}":
		return Maps(in.FromValue)

	default:
		if in.ReferValue != nil {
			var referReflectValue reflect.Value
			if v, ok := in.ReferValue.(reflect.Value); ok {
				referReflectValue = v
			} else {
				referReflectValue = reflect.ValueOf(in.ReferValue)
			}
			in.ToTypeName = referReflectValue.Kind().String()
			in.ReferValue = nil
			return reflect.ValueOf(doConvert(in)).Convert(referReflectValue.Type()).Interface()
		}
		return in.FromValue
	}
}

// Byte converts `any` to byte.
func Byte(any interface{}) byte {
	if v, ok := any.(byte); ok {
		return v
	}
	return Uint8(any)
}

// Bytes converts `any` to []byte.
func Bytes(any interface{}) []byte {
	if any == nil {
		return nil
	}
	switch value := any.(type) {
	case string:
		return []byte(value)

	case []byte:
		return value

	default:
		if f, ok := value.(iBytes); ok {
			return f.Bytes()
		}
		originValueAndKind := utils.OriginValueAndKind(any)
		switch originValueAndKind.OriginKind {
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
func Rune(any interface{}) rune {
	if v, ok := any.(rune); ok {
		return v
	}
	return Int32(any)
}

// Runes converts `any` to []rune.
func Runes(any interface{}) []rune {
	if v, ok := any.([]rune); ok {
		return v
	}
	return []rune(String(any))
}

// String converts `any` to string.
// It's most commonly used converting function.
func String(any interface{}) string {
	if any == nil {
		return ""
	}
	switch value := any.(type) {
	case int:
		return strconv.Itoa(value)
	case int8:
		return strconv.Itoa(int(value))
	case int16:
		return strconv.Itoa(int(value))
	case int32:
		return strconv.Itoa(int(value))
	case int64:
		return strconv.FormatInt(value, 10)
	case uint:
		return strconv.FormatUint(uint64(value), 10)
	case uint8:
		return strconv.FormatUint(uint64(value), 10)
	case uint16:
		return strconv.FormatUint(uint64(value), 10)
	case uint32:
		return strconv.FormatUint(uint64(value), 10)
	case uint64:
		return strconv.FormatUint(value, 10)
	case float32:
		return strconv.FormatFloat(float64(value), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(value, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(value)
	case string:
		return value
	case []byte:
		return string(value)
	case time.Time:
		if value.IsZero() {
			return ""
		}
		return value.String()
	case *time.Time:
		if value == nil {
			return ""
		}
		return value.String()
	case gtime.Time:
		if value.IsZero() {
			return ""
		}
		return value.String()
	case *gtime.Time:
		if value == nil {
			return ""
		}
		return value.String()
	default:
		// Empty checks.
		if value == nil {
			return ""
		}
		if f, ok := value.(iString); ok {
			// If the variable implements the String() interface,
			// then use that interface to perform the conversion
			return f.String()
		}
		if f, ok := value.(iError); ok {
			// If the variable implements the Error() interface,
			// then use that interface to perform the conversion
			return f.Error()
		}
		// Reflect checks.
		var (
			rv   = reflect.ValueOf(value)
			kind = rv.Kind()
		)
		switch kind {
		case reflect.Chan,
			reflect.Map,
			reflect.Slice,
			reflect.Func,
			reflect.Ptr,
			reflect.Interface,
			reflect.UnsafePointer:
			if rv.IsNil() {
				return ""
			}
		case reflect.String:
			return rv.String()
		}
		if kind == reflect.Ptr {
			return String(rv.Elem().Interface())
		}
		// Finally, we use json.Marshal to convert.
		if jsonContent, err := json.Marshal(value); err != nil {
			return fmt.Sprint(value)
		} else {
			return string(jsonContent)
		}
	}
}

// Bool converts `any` to bool.
// It returns false if `any` is: false, "", 0, "false", "off", "no", empty slice/map.
func Bool(any interface{}) bool {
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
		if f, ok := value.(iBool); ok {
			return f.Bool()
		}
		rv := reflect.ValueOf(any)
		switch rv.Kind() {
		case reflect.Ptr:
			return !rv.IsNil()
		case reflect.Map:
			fallthrough
		case reflect.Array:
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
func checkJsonAndUnmarshalUseNumber(any interface{}, target interface{}) bool {
	switch r := any.(type) {
	case []byte:
		if json.Valid(r) {
			_ = json.UnmarshalUseNumber(r, &target)
			return true
		}

	case string:
		anyAsBytes := []byte(r)
		if json.Valid(anyAsBytes) {
			_ = json.UnmarshalUseNumber(anyAsBytes, &target)
			return true
		}
	}
	return false
}
