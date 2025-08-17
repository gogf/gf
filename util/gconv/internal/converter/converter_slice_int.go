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

// SliceInt converts `any` to []int.
func (c *Converter) SliceInt(any any, option ...SliceOption) ([]int, error) {
	if empty.IsNil(any) {
		return nil, nil
	}
	var (
		err         error
		ii          int
		array       []int = nil
		sliceOption       = c.getSliceOption(option...)
	)
	switch value := any.(type) {
	case []string:
		array = make([]int, len(value))
		for k, v := range value {
			ii, err = c.Int(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = ii
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
		if json.Valid(value) {
			if err = json.UnmarshalUseNumber(value, &array); array != nil {
				return array, err
			}
		}
		array = make([]int, len(value))
		for k, v := range value {
			array[k] = int(v)
		}
	case string:
		byteValue := []byte(value)
		if json.Valid(byteValue) {
			if err = json.UnmarshalUseNumber(byteValue, &array); array != nil {
				return array, err
			}
		}
		if value == "" {
			return []int{}, err
		}
		if utils.IsNumeric(value) {
			ii, err = c.Int(value)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			return []int{ii}, err
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
			ii, err = c.Int(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = ii
		}
	case []float64:
		array = make([]int, len(value))
		for k, v := range value {
			ii, err = c.Int(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = ii
		}
	case []interface{}:
		array = make([]int, len(value))
		for k, v := range value {
			ii, err = c.Int(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = ii
		}
	case [][]byte:
		array = make([]int, len(value))
		for k, v := range value {
			ii, err = c.Int(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = ii
		}
	}
	if array != nil {
		return array, err
	}
	if v, ok := any.(localinterface.IInts); ok {
		return v.Ints(), err
	}
	if v, ok := any.(localinterface.IInterfaces); ok {
		return c.SliceInt(v.Interfaces(), option...)
	}
	// Not a common type, it then uses reflection for conversion.
	originValueAndKind := reflection.OriginValueAndKind(any)
	switch originValueAndKind.OriginKind {
	case reflect.Slice, reflect.Array:
		var (
			length = originValueAndKind.OriginValue.Len()
			slice  = make([]int, length)
		)
		for i := 0; i < length; i++ {
			ii, err = c.Int(originValueAndKind.OriginValue.Index(i).Interface())
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			slice[i] = ii
		}
		return slice, err

	default:
		if originValueAndKind.OriginValue.IsZero() {
			return []int{}, err
		}
		ii, err = c.Int(any)
		if err != nil && !sliceOption.ContinueOnError {
			return nil, err
		}
		return []int{ii}, err
	}
}

// SliceInt32 converts `any` to []int32.
func (c *Converter) SliceInt32(any any, option ...SliceOption) ([]int32, error) {
	if empty.IsNil(any) {
		return nil, nil
	}
	var (
		err         error
		ii          int32
		array       []int32 = nil
		sliceOption         = c.getSliceOption(option...)
	)
	switch value := any.(type) {
	case []string:
		array = make([]int32, len(value))
		for k, v := range value {
			ii, err = c.Int32(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = ii
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
		if json.Valid(value) {
			if err = json.UnmarshalUseNumber(value, &array); array != nil {
				return array, err
			}
		}
		array = make([]int32, len(value))
		for k, v := range value {
			array[k] = int32(v)
		}
	case string:
		byteValue := []byte(value)
		if json.Valid(byteValue) {
			if err = json.UnmarshalUseNumber(byteValue, &array); array != nil {
				return array, err
			}
		}
		if value == "" {
			return []int32{}, err
		}
		if utils.IsNumeric(value) {
			ii, err = c.Int32(value)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			return []int32{ii}, err
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
			ii, err = c.Int32(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = ii
		}
	case []float64:
		array = make([]int32, len(value))
		for k, v := range value {
			ii, err = c.Int32(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = ii
		}
	case []interface{}:
		array = make([]int32, len(value))
		for k, v := range value {
			ii, err = c.Int32(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = ii
		}
	case [][]byte:
		array = make([]int32, len(value))
		for k, v := range value {
			ii, err = c.Int32(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = ii
		}
	}
	if array != nil {
		return array, err
	}
	if v, ok := any.(localinterface.IInts); ok {
		return c.SliceInt32(v.Ints(), option...)
	}
	if v, ok := any.(localinterface.IInterfaces); ok {
		return c.SliceInt32(v.Interfaces(), option...)
	}
	// Not a common type, it then uses reflection for conversion.
	originValueAndKind := reflection.OriginValueAndKind(any)
	switch originValueAndKind.OriginKind {
	case reflect.Slice, reflect.Array:
		var (
			length = originValueAndKind.OriginValue.Len()
			slice  = make([]int32, length)
		)
		for i := 0; i < length; i++ {
			ii, err = c.Int32(originValueAndKind.OriginValue.Index(i).Interface())
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			slice[i] = ii
		}
		return slice, err

	default:
		if originValueAndKind.OriginValue.IsZero() {
			return []int32{}, err
		}
		ii, err = c.Int32(any)
		if err != nil && !sliceOption.ContinueOnError {
			return nil, err
		}
		return []int32{ii}, err
	}
}

// SliceInt64 converts `any` to []int64.
func (c *Converter) SliceInt64(any any, option ...SliceOption) ([]int64, error) {
	if empty.IsNil(any) {
		return nil, nil
	}
	var (
		err         error
		ii          int64
		array       []int64 = nil
		sliceOption         = c.getSliceOption(option...)
	)
	switch value := any.(type) {
	case []string:
		array = make([]int64, len(value))
		for k, v := range value {
			ii, err = c.Int64(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = ii
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
		if json.Valid(value) {
			if err = json.UnmarshalUseNumber(value, &array); array != nil {
				return array, err
			}
		}
		array = make([]int64, len(value))
		for k, v := range value {
			array[k] = int64(v)
		}
	case string:
		byteValue := []byte(value)
		if json.Valid(byteValue) {
			if err = json.UnmarshalUseNumber(byteValue, &array); array != nil {
				return array, err
			}
		}
		if value == "" {
			return []int64{}, err
		}
		if utils.IsNumeric(value) {
			ii, err = c.Int64(value)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			return []int64{ii}, err
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
			ii, err = c.Int64(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = ii
		}
	case []float64:
		array = make([]int64, len(value))
		for k, v := range value {
			ii, err = c.Int64(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = ii
		}
	case []interface{}:
		array = make([]int64, len(value))
		for k, v := range value {
			ii, err = c.Int64(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = ii
		}
	case [][]byte:
		array = make([]int64, len(value))
		for k, v := range value {
			ii, err = c.Int64(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = ii
		}
	}
	if array != nil {
		return array, err
	}
	if v, ok := any.(localinterface.IInts); ok {
		return c.SliceInt64(v.Ints(), option...)
	}
	if v, ok := any.(localinterface.IInterfaces); ok {
		return c.SliceInt64(v.Interfaces(), option...)
	}
	// Not a common type, it then uses reflection for conversion.
	originValueAndKind := reflection.OriginValueAndKind(any)
	switch originValueAndKind.OriginKind {
	case reflect.Slice, reflect.Array:
		var (
			length = originValueAndKind.OriginValue.Len()
			slice  = make([]int64, length)
		)
		for i := 0; i < length; i++ {
			ii, err = c.Int64(originValueAndKind.OriginValue.Index(i).Interface())
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			slice[i] = ii
		}
		return slice, err

	default:
		if originValueAndKind.OriginValue.IsZero() {
			return []int64{}, err
		}
		ii, err = c.Int64(any)
		if err != nil && !sliceOption.ContinueOnError {
			return nil, err
		}
		return []int64{ii}, err
	}
}
