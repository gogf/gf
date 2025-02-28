// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv

import (
	"bytes"
	"reflect"
	"strings"

	"github.com/gogf/gf/v2/internal/empty"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/internal/reflection"
	"github.com/gogf/gf/v2/internal/utils"
	"github.com/gogf/gf/v2/util/gconv/internal/localinterface"
)

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

// Floats converts `any` to []float64.
func Floats(any interface{}) []float64 {
	return Float64s(any)
}

// Float32s converts `any` to []float32.
func Float32s(any interface{}) []float32 {
	if empty.IsNil(any) {
		return nil
	}
	var (
		array []float32 = nil
	)
	switch value := any.(type) {
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
		if json.Valid(value) {
			if _ = json.UnmarshalUseNumber(value, &array); array != nil {
				return array
			}
			if bytes.EqualFold([]byte("null"), value) {
				return nil
			}
		}
		array = make([]float32, len(value))
		for k, v := range value {
			array[k] = Float32(v)
		}
	case string:
		byteValue := []byte(value)
		if json.Valid(byteValue) {
			if _ = json.UnmarshalUseNumber(byteValue, &array); array != nil {
				return array
			}
			if strings.EqualFold(value, "null") {
				return nil
			}
		}
		if value == "" {
			return []float32{}
		}
		if utils.IsNumeric(value) {
			return []float32{Float32(value)}
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
	}
	if array != nil {
		return array
	}
	if v, ok := any.(localinterface.IFloats); ok {
		return Float32s(v.Floats())
	}
	if v, ok := any.(localinterface.IInterfaces); ok {
		return Float32s(v.Interfaces())
	}
	// Not a common type, it then uses reflection for conversion.
	originValueAndKind := reflection.OriginValueAndKind(any)
	switch originValueAndKind.OriginKind {
	case reflect.Slice, reflect.Array:
		var (
			length = originValueAndKind.OriginValue.Len()
			slice  = make([]float32, length)
		)
		for i := 0; i < length; i++ {
			slice[i] = Float32(originValueAndKind.OriginValue.Index(i).Interface())
		}
		return slice

	default:
		if originValueAndKind.OriginValue.IsZero() {
			return []float32{}
		}
		return []float32{Float32(any)}
	}
}

// Float64s converts `any` to []float64.
func Float64s(any interface{}) []float64 {
	if empty.IsNil(any) {
		return nil
	}
	var (
		array []float64 = nil
	)
	switch value := any.(type) {
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
		if json.Valid(value) {
			if _ = json.UnmarshalUseNumber(value, &array); array != nil {
				return array
			}
			if bytes.EqualFold([]byte("null"), value) {
				return nil
			}
		}
		array = make([]float64, len(value))
		for k, v := range value {
			array[k] = Float64(v)
		}
	case string:
		byteValue := []byte(value)
		if json.Valid(byteValue) {
			if _ = json.UnmarshalUseNumber(byteValue, &array); array != nil {
				return array
			}
			if strings.EqualFold(value, "null") {
				return nil
			}
		}
		if value == "" {
			return []float64{}
		}
		if utils.IsNumeric(value) {
			return []float64{Float64(value)}
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
	}
	if array != nil {
		return array
	}
	if v, ok := any.(localinterface.IFloats); ok {
		return v.Floats()
	}
	if v, ok := any.(localinterface.IInterfaces); ok {
		return Floats(v.Interfaces())
	}
	// Not a common type, it then uses reflection for conversion.
	originValueAndKind := reflection.OriginValueAndKind(any)
	switch originValueAndKind.OriginKind {
	case reflect.Slice, reflect.Array:
		var (
			length = originValueAndKind.OriginValue.Len()
			slice  = make([]float64, length)
		)
		for i := 0; i < length; i++ {
			slice[i] = Float64(originValueAndKind.OriginValue.Index(i).Interface())
		}
		return slice

	default:
		if originValueAndKind.OriginValue.IsZero() {
			return []float64{}
		}
		return []float64{Float64(any)}
	}
}
