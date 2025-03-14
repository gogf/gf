// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package converter

import (
	"reflect"

	"github.com/gogf/gf/v2/internal/empty"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/internal/reflection"
	"github.com/gogf/gf/v2/internal/utils"
	"github.com/gogf/gf/v2/util/gconv/internal/localinterface"
)

// SliceFloat32 converts `any` to []float32.
func (c *Converter) SliceFloat32(any interface{}, option ...SliceOption) ([]float32, error) {
	if empty.IsNil(any) {
		return nil, nil
	}
	var (
		err         error
		f           float32
		array       []float32 = nil
		sliceOption           = c.getSliceOption(option...)
	)
	switch value := any.(type) {
	case []string:
		array = make([]float32, len(value))
		for k, v := range value {
			f, err = c.Float32(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = f
		}
	case []int:
		array = make([]float32, len(value))
		for k, v := range value {
			f, err = c.Float32(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = f
		}
	case []int8:
		array = make([]float32, len(value))
		for k, v := range value {
			f, err = c.Float32(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = f
		}
	case []int16:
		array = make([]float32, len(value))
		for k, v := range value {
			f, err = c.Float32(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = f
		}
	case []int32:
		array = make([]float32, len(value))
		for k, v := range value {
			f, err = c.Float32(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = f
		}
	case []int64:
		array = make([]float32, len(value))
		for k, v := range value {
			f, err = c.Float32(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = f
		}
	case []uint:
		array = make([]float32, len(value))
		for k, v := range value {
			f, err = c.Float32(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = f
		}
	case []uint8:
		if json.Valid(value) {
			if err = json.UnmarshalUseNumber(value, &array); array != nil {
				return array, err
			}
		}
		array = make([]float32, len(value))
		for k, v := range value {
			f, err = c.Float32(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = f
		}
	case string:
		byteValue := []byte(value)
		if json.Valid(byteValue) {
			if err = json.UnmarshalUseNumber(byteValue, &array); array != nil {
				return array, err
			}
		}
		if value == "" {
			return []float32{}, err
		}
		if utils.IsNumeric(value) {
			f, err = c.Float32(value)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			return []float32{f}, err
		}
	case []uint16:
		array = make([]float32, len(value))
		for k, v := range value {
			f, err = c.Float32(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = f
		}
	case []uint32:
		array = make([]float32, len(value))
		for k, v := range value {
			f, err = c.Float32(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = f
		}
	case []uint64:
		array = make([]float32, len(value))
		for k, v := range value {
			f, err = c.Float32(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = f
		}
	case []bool:
		array = make([]float32, len(value))
		for k, v := range value {
			f, err = c.Float32(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = f
		}
	case []float32:
		array = value
	case []float64:
		array = make([]float32, len(value))
		for k, v := range value {
			f, err = c.Float32(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = f
		}
	case []interface{}:
		array = make([]float32, len(value))
		for k, v := range value {
			f, err = c.Float32(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = f
		}
	}
	if array != nil {
		return array, err
	}
	if v, ok := any.(localinterface.IFloats); ok {
		return c.SliceFloat32(v.Floats(), option...)
	}
	if v, ok := any.(localinterface.IInterfaces); ok {
		return c.SliceFloat32(v.Interfaces(), option...)
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
			f, err = c.Float32(originValueAndKind.OriginValue.Index(i).Interface())
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			slice[i] = f
		}
		return slice, err

	default:
		if originValueAndKind.OriginValue.IsZero() {
			return []float32{}, err
		}
		f, err = c.Float32(any)
		if err != nil && !sliceOption.ContinueOnError {
			return nil, err
		}
		return []float32{f}, err
	}
}

// SliceFloat64 converts `any` to []float64.
func (c *Converter) SliceFloat64(any interface{}, option ...SliceOption) ([]float64, error) {
	if empty.IsNil(any) {
		return nil, nil
	}
	var (
		err         error
		f           float64
		array       []float64 = nil
		sliceOption           = c.getSliceOption(option...)
	)
	switch value := any.(type) {
	case []string:
		array = make([]float64, len(value))
		for k, v := range value {
			f, err = c.Float64(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = f
		}
	case []int:
		array = make([]float64, len(value))
		for k, v := range value {
			f, err = c.Float64(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = f
		}
	case []int8:
		array = make([]float64, len(value))
		for k, v := range value {
			f, err = c.Float64(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = f
		}
	case []int16:
		array = make([]float64, len(value))
		for k, v := range value {
			f, err = c.Float64(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = f
		}
	case []int32:
		array = make([]float64, len(value))
		for k, v := range value {
			f, err = c.Float64(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = f
		}
	case []int64:
		array = make([]float64, len(value))
		for k, v := range value {
			f, err = c.Float64(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = f
		}
	case []uint:
		array = make([]float64, len(value))
		for k, v := range value {
			f, err = c.Float64(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = f
		}
	case []uint8:
		if json.Valid(value) {
			if err = json.UnmarshalUseNumber(value, &array); array != nil {
				return array, err
			}
		}
		array = make([]float64, len(value))
		for k, v := range value {
			f, err = c.Float64(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = f
		}
	case string:
		byteValue := []byte(value)
		if json.Valid(byteValue) {
			if err = json.UnmarshalUseNumber(byteValue, &array); array != nil {
				return array, err
			}
		}
		if value == "" {
			return []float64{}, err
		}
		if utils.IsNumeric(value) {
			f, err = c.Float64(value)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			return []float64{f}, err
		}
	case []uint16:
		array = make([]float64, len(value))
		for k, v := range value {
			f, err = c.Float64(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = f
		}
	case []uint32:
		array = make([]float64, len(value))
		for k, v := range value {
			f, err = c.Float64(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = f
		}
	case []uint64:
		array = make([]float64, len(value))
		for k, v := range value {
			f, err = c.Float64(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = f
		}
	case []bool:
		array = make([]float64, len(value))
		for k, v := range value {
			f, err = c.Float64(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = f
		}
	case []float32:
		array = make([]float64, len(value))
		for k, v := range value {
			f, err = c.Float64(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = f
		}
	case []float64:
		array = value
	case []interface{}:
		array = make([]float64, len(value))
		for k, v := range value {
			f, err = c.Float64(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = f
		}
	}
	if array != nil {
		return array, err
	}
	if v, ok := any.(localinterface.IFloats); ok {
		return v.Floats(), err
	}
	if v, ok := any.(localinterface.IInterfaces); ok {
		return c.SliceFloat64(v.Interfaces(), option...)
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
			f, err = c.Float64(originValueAndKind.OriginValue.Index(i).Interface())
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			slice[i] = f
		}
		return slice, err

	default:
		if originValueAndKind.OriginValue.IsZero() {
			return []float64{}, err
		}
		f, err = c.Float64(any)
		if err != nil && !sliceOption.ContinueOnError {
			return nil, err
		}
		return []float64{f}, err
	}
}
