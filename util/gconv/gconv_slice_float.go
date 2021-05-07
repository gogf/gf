// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv

import "reflect"

// SliceFloat is alias of Floats.
func SliceFloat(any interface{}) []float64 {
	return Floats(any)
}

// SliceFloat32 is alias of Float32s.
func SliceFloat32(any interface{}) []float32 {
	return Float32s(any)
}

// SliceFloat64 is alias of Float64s.
func SliceFloat64(any interface{}) []float64 {
	return Floats(any)
}

// Floats converts `i` to []float64.
func Floats(any interface{}) []float64 {
	return Float64s(any)
}

// Float32s converts `i` to []float32.
func Float32s(any interface{}) []float32 {
	if any == nil {
		return nil
	}
	var array []float32
	switch value := any.(type) {
	case string:
		if value == "" {
			return []float32{}
		}
		return []float32{Float32(value)}
	case []string:
		array = make([]float32, len(value))
		for k, v := range value {
			array[k] = Float32(v)
		}
	case []int:
		array = make([]float32, len(value))
		for k, v := range value {
			array[k] = Float32(v)
		}
	case []int8:
		array = make([]float32, len(value))
		for k, v := range value {
			array[k] = Float32(v)
		}
	case []int16:
		array = make([]float32, len(value))
		for k, v := range value {
			array[k] = Float32(v)
		}
	case []int32:
		array = make([]float32, len(value))
		for k, v := range value {
			array[k] = Float32(v)
		}
	case []int64:
		array = make([]float32, len(value))
		for k, v := range value {
			array[k] = Float32(v)
		}
	case []uint:
		for _, v := range value {
			array = append(array, Float32(v))
		}
	case []uint8:
		array = make([]float32, len(value))
		for k, v := range value {
			array[k] = Float32(v)
		}
	case []uint16:
		array = make([]float32, len(value))
		for k, v := range value {
			array[k] = Float32(v)
		}
	case []uint32:
		array = make([]float32, len(value))
		for k, v := range value {
			array[k] = Float32(v)
		}
	case []uint64:
		array = make([]float32, len(value))
		for k, v := range value {
			array[k] = Float32(v)
		}
	case []bool:
		array = make([]float32, len(value))
		for k, v := range value {
			array[k] = Float32(v)
		}
	case []float32:
		array = value
	case []float64:
		array = make([]float32, len(value))
		for k, v := range value {
			array[k] = Float32(v)
		}
	case []interface{}:
		array = make([]float32, len(value))
		for k, v := range value {
			array[k] = Float32(v)
		}
	default:
		if v, ok := any.(apiFloats); ok {
			return Float32s(v.Floats())
		}
		if v, ok := any.(apiInterfaces); ok {
			return Float32s(v.Interfaces())
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
				slice  = make([]float32, length)
			)
			for i := 0; i < length; i++ {
				slice[i] = Float32(reflectValue.Index(i).Interface())
			}
			return slice

		default:
			return []float32{Float32(any)}
		}
	}
	return array
}

// Float64s converts `i` to []float64.
func Float64s(any interface{}) []float64 {
	if any == nil {
		return nil
	}
	var array []float64
	switch value := any.(type) {
	case string:
		if value == "" {
			return []float64{}
		}
		return []float64{Float64(value)}
	case []string:
		array = make([]float64, len(value))
		for k, v := range value {
			array[k] = Float64(v)
		}
	case []int:
		array = make([]float64, len(value))
		for k, v := range value {
			array[k] = Float64(v)
		}
	case []int8:
		array = make([]float64, len(value))
		for k, v := range value {
			array[k] = Float64(v)
		}
	case []int16:
		array = make([]float64, len(value))
		for k, v := range value {
			array[k] = Float64(v)
		}
	case []int32:
		array = make([]float64, len(value))
		for k, v := range value {
			array[k] = Float64(v)
		}
	case []int64:
		array = make([]float64, len(value))
		for k, v := range value {
			array[k] = Float64(v)
		}
	case []uint:
		for _, v := range value {
			array = append(array, Float64(v))
		}
	case []uint8:
		array = make([]float64, len(value))
		for k, v := range value {
			array[k] = Float64(v)
		}
	case []uint16:
		array = make([]float64, len(value))
		for k, v := range value {
			array[k] = Float64(v)
		}
	case []uint32:
		array = make([]float64, len(value))
		for k, v := range value {
			array[k] = Float64(v)
		}
	case []uint64:
		array = make([]float64, len(value))
		for k, v := range value {
			array[k] = Float64(v)
		}
	case []bool:
		array = make([]float64, len(value))
		for k, v := range value {
			array[k] = Float64(v)
		}
	case []float32:
		array = make([]float64, len(value))
		for k, v := range value {
			array[k] = Float64(v)
		}
	case []float64:
		array = value
	case []interface{}:
		array = make([]float64, len(value))
		for k, v := range value {
			array[k] = Float64(v)
		}
	default:
		if v, ok := any.(apiFloats); ok {
			return v.Floats()
		}
		if v, ok := any.(apiInterfaces); ok {
			return Floats(v.Interfaces())
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
				slice  = make([]float64, length)
			)
			for i := 0; i < length; i++ {
				slice[i] = Float64(reflectValue.Index(i).Interface())
			}
			return slice

		default:
			return []float64{Float64(any)}
		}
	}
	return array

}
