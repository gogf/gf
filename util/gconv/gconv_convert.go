// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv

import (
	"reflect"
	"time"

	"github.com/gogf/gf/v2/os/gtime"
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
	// Marks that the value is already converted and set to `ReferValue`. Caller can ignore the returned result.
	// It is an attribute for internal usage purpose.
	alreadySetToReferValue bool
}

// doConvert does commonly use types converting.
func doConvert(in doConvertInput) (convertedValue interface{}) {
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

	case "json.RawMessage":
		return Bytes(in.FromValue)

	default:
		if in.ReferValue != nil {
			var referReflectValue reflect.Value
			if v, ok := in.ReferValue.(reflect.Value); ok {
				referReflectValue = v
			} else {
				referReflectValue = reflect.ValueOf(in.ReferValue)
			}
			defer func() {
				if recover() != nil {
					if err := bindVarToReflectValue(referReflectValue, in.FromValue, nil); err == nil {
						in.alreadySetToReferValue = true
						convertedValue = referReflectValue.Interface()
					}
				}
			}()
			in.ToTypeName = referReflectValue.Kind().String()
			in.ReferValue = nil
			return reflect.ValueOf(doConvert(in)).Convert(referReflectValue.Type()).Interface()
		}
		return in.FromValue
	}
}

func doConvertWithReflectValueSet(reflectValue reflect.Value, in doConvertInput) {
	convertedValue := doConvert(in)
	if !in.alreadySetToReferValue {
		reflectValue.Set(reflect.ValueOf(convertedValue))
	}
}
