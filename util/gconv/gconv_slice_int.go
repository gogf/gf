// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv

import "reflect"

// SliceInt is alias of Ints.
func SliceInt(any interface{}) []int {
	return Ints(any)
}

// SliceInt32 is alias of Int32s.
func SliceInt32(any interface{}) []int32 {
	return Int32s(any)
}

// SliceInt is alias of Int64s.
func SliceInt64(any interface{}) []int64 {
	return Int64s(any)
}

// Ints converts `i` to []int.
func Ints(any interface{}) []int {
	if any == nil {
		return nil
	}
	var array []int
	switch value := any.(type) {
	case string:
		if value == "" {
			return []int{}
		}
		return []int{Int(value)}
	case []string:
		array = make([]int, len(value))
		for k, v := range value {
			array[k] = Int(v)
		}
	case []int:
		array = value
	case []int8:
		array = make([]int, len(value))
		for k, v := range value {
			array[k] = int(v)
		}
	case []int16:
		array = make([]int, len(value))
		for k, v := range value {
			array[k] = int(v)
		}
	case []int32:
		array = make([]int, len(value))
		for k, v := range value {
			array[k] = int(v)
		}
	case []int64:
		array = make([]int, len(value))
		for k, v := range value {
			array[k] = int(v)
		}
	case []uint:
		array = make([]int, len(value))
		for k, v := range value {
			array[k] = int(v)
		}
	case []uint8:
		array = make([]int, len(value))
		for k, v := range value {
			array[k] = int(v)
		}
	case []uint16:
		array = make([]int, len(value))
		for k, v := range value {
			array[k] = int(v)
		}
	case []uint32:
		array = make([]int, len(value))
		for k, v := range value {
			array[k] = int(v)
		}
	case []uint64:
		array = make([]int, len(value))
		for k, v := range value {
			array[k] = int(v)
		}
	case []bool:
		array = make([]int, len(value))
		for k, v := range value {
			if v {
				array[k] = 1
			} else {
				array[k] = 0
			}
		}
	case []float32:
		array = make([]int, len(value))
		for k, v := range value {
			array[k] = Int(v)
		}
	case []float64:
		array = make([]int, len(value))
		for k, v := range value {
			array[k] = Int(v)
		}
	case []interface{}:
		array = make([]int, len(value))
		for k, v := range value {
			array[k] = Int(v)
		}
	case [][]byte:
		array = make([]int, len(value))
		for k, v := range value {
			array[k] = Int(v)
		}
	default:
		if v, ok := any.(apiInts); ok {
			return v.Ints()
		}
		if v, ok := any.(apiInterfaces); ok {
			return Ints(v.Interfaces())
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
				slice  = make([]int, length)
			)
			for i := 0; i < length; i++ {
				slice[i] = Int(reflectValue.Index(i).Interface())
			}
			return slice

		default:
			return []int{Int(any)}
		}
	}
	return array
}

// Int32s converts `i` to []int32.
func Int32s(any interface{}) []int32 {
	if any == nil {
		return nil
	}
	var array []int32
	switch value := any.(type) {
	case string:
		if value == "" {
			return []int32{}
		}
		return []int32{Int32(value)}
	case []string:
		array = make([]int32, len(value))
		for k, v := range value {
			array[k] = Int32(v)
		}
	case []int:
		array = make([]int32, len(value))
		for k, v := range value {
			array[k] = int32(v)
		}
	case []int8:
		array = make([]int32, len(value))
		for k, v := range value {
			array[k] = int32(v)
		}
	case []int16:
		array = make([]int32, len(value))
		for k, v := range value {
			array[k] = int32(v)
		}
	case []int32:
		array = value
	case []int64:
		array = make([]int32, len(value))
		for k, v := range value {
			array[k] = int32(v)
		}
	case []uint:
		array = make([]int32, len(value))
		for k, v := range value {
			array[k] = int32(v)
		}
	case []uint8:
		array = make([]int32, len(value))
		for k, v := range value {
			array[k] = int32(v)
		}
	case []uint16:
		array = make([]int32, len(value))
		for k, v := range value {
			array[k] = int32(v)
		}
	case []uint32:
		array = make([]int32, len(value))
		for k, v := range value {
			array[k] = int32(v)
		}
	case []uint64:
		array = make([]int32, len(value))
		for k, v := range value {
			array[k] = int32(v)
		}
	case []bool:
		array = make([]int32, len(value))
		for k, v := range value {
			if v {
				array[k] = 1
			} else {
				array[k] = 0
			}
		}
	case []float32:
		array = make([]int32, len(value))
		for k, v := range value {
			array[k] = Int32(v)
		}
	case []float64:
		array = make([]int32, len(value))
		for k, v := range value {
			array[k] = Int32(v)
		}
	case []interface{}:
		array = make([]int32, len(value))
		for k, v := range value {
			array[k] = Int32(v)
		}
	case [][]byte:
		array = make([]int32, len(value))
		for k, v := range value {
			array[k] = Int32(v)
		}
	default:
		if v, ok := any.(apiInts); ok {
			return Int32s(v.Ints())
		}
		if v, ok := any.(apiInterfaces); ok {
			return Int32s(v.Interfaces())
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
				slice  = make([]int32, length)
			)
			for i := 0; i < length; i++ {
				slice[i] = Int32(reflectValue.Index(i).Interface())
			}
			return slice

		default:
			return []int32{Int32(any)}
		}
	}
	return array
}

// Int64s converts `i` to []int64.
func Int64s(any interface{}) []int64 {
	if any == nil {
		return nil
	}
	var array []int64
	switch value := any.(type) {
	case string:
		if value == "" {
			return []int64{}
		}
		return []int64{Int64(value)}
	case []string:
		array = make([]int64, len(value))
		for k, v := range value {
			array[k] = Int64(v)
		}
	case []int:
		array = make([]int64, len(value))
		for k, v := range value {
			array[k] = int64(v)
		}
	case []int8:
		array = make([]int64, len(value))
		for k, v := range value {
			array[k] = int64(v)
		}
	case []int16:
		array = make([]int64, len(value))
		for k, v := range value {
			array[k] = int64(v)
		}
	case []int32:
		array = make([]int64, len(value))
		for k, v := range value {
			array[k] = int64(v)
		}
	case []int64:
		array = value
	case []uint:
		array = make([]int64, len(value))
		for k, v := range value {
			array[k] = int64(v)
		}
	case []uint8:
		array = make([]int64, len(value))
		for k, v := range value {
			array[k] = int64(v)
		}
	case []uint16:
		array = make([]int64, len(value))
		for k, v := range value {
			array[k] = int64(v)
		}
	case []uint32:
		array = make([]int64, len(value))
		for k, v := range value {
			array[k] = int64(v)
		}
	case []uint64:
		array = make([]int64, len(value))
		for k, v := range value {
			array[k] = int64(v)
		}
	case []bool:
		array = make([]int64, len(value))
		for k, v := range value {
			if v {
				array[k] = 1
			} else {
				array[k] = 0
			}
		}
	case []float32:
		array = make([]int64, len(value))
		for k, v := range value {
			array[k] = Int64(v)
		}
	case []float64:
		array = make([]int64, len(value))
		for k, v := range value {
			array[k] = Int64(v)
		}
	case []interface{}:
		array = make([]int64, len(value))
		for k, v := range value {
			array[k] = Int64(v)
		}
	case [][]byte:
		array = make([]int64, len(value))
		for k, v := range value {
			array[k] = Int64(v)
		}
	default:
		if v, ok := any.(apiInts); ok {
			return Int64s(v.Ints())
		}
		if v, ok := any.(apiInterfaces); ok {
			return Int64s(v.Interfaces())
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
				slice  = make([]int64, length)
			)
			for i := 0; i < length; i++ {
				slice[i] = Int64(reflectValue.Index(i).Interface())
			}
			return slice

		default:
			return []int64{Int64(any)}
		}
	}
	return array
}
