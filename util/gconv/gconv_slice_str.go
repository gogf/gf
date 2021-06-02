// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv

import "reflect"

// SliceStr is alias of Strings.
func SliceStr(any interface{}) []string {
	return Strings(any)
}

// Strings converts `any` to []string.
func Strings(any interface{}) []string {
	if any == nil {
		return nil
	}
	var array []string
	switch value := any.(type) {
	case []int:
		array = make([]string, len(value))
		for k, v := range value {
			array[k] = String(v)
		}
	case []int8:
		array = make([]string, len(value))
		for k, v := range value {
			array[k] = String(v)
		}
	case []int16:
		array = make([]string, len(value))
		for k, v := range value {
			array[k] = String(v)
		}
	case []int32:
		array = make([]string, len(value))
		for k, v := range value {
			array[k] = String(v)
		}
	case []int64:
		array = make([]string, len(value))
		for k, v := range value {
			array[k] = String(v)
		}
	case []uint:
		array = make([]string, len(value))
		for k, v := range value {
			array[k] = String(v)
		}
	case []uint8:
		array = make([]string, len(value))
		for k, v := range value {
			array[k] = String(v)
		}
	case []uint16:
		array = make([]string, len(value))
		for k, v := range value {
			array[k] = String(v)
		}
	case []uint32:
		array = make([]string, len(value))
		for k, v := range value {
			array[k] = String(v)
		}
	case []uint64:
		array = make([]string, len(value))
		for k, v := range value {
			array[k] = String(v)
		}
	case []bool:
		array = make([]string, len(value))
		for k, v := range value {
			array[k] = String(v)
		}
	case []float32:
		array = make([]string, len(value))
		for k, v := range value {
			array[k] = String(v)
		}
	case []float64:
		array = make([]string, len(value))
		for k, v := range value {
			array[k] = String(v)
		}
	case []interface{}:
		array = make([]string, len(value))
		for k, v := range value {
			array[k] = String(v)
		}
	case []string:
		array = value
	case [][]byte:
		array = make([]string, len(value))
		for k, v := range value {
			array[k] = String(v)
		}
	default:
		if v, ok := any.(apiStrings); ok {
			return v.Strings()
		}
		if v, ok := any.(apiInterfaces); ok {
			return Strings(v.Interfaces())
		}
		// Not a common type, it then uses reflection for conversion.
		var reflectValue reflect.Value
		if v, ok := value.(reflect.Value); ok {
			reflectValue = v
		} else {
			reflectValue = reflect.ValueOf(value)
		}
		reflectKind := reflectValue.Kind()
		for reflectKind == reflect.Ptr {
			reflectValue = reflectValue.Elem()
			reflectKind = reflectValue.Kind()
		}
		switch reflectKind {
		case reflect.Slice, reflect.Array:
			var (
				length = reflectValue.Len()
				slice  = make([]string, length)
			)
			for i := 0; i < length; i++ {
				slice[i] = String(reflectValue.Index(i).Interface())
			}
			return slice

		default:
			return []string{String(any)}
		}
	}
	return array
}
