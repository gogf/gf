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
	"github.com/gogf/gf/v2/util/gconv/internal/localinterface"
)

// SliceBool converts `any` to []bool.
func (c *Converter) SliceBool(anyInput any, option ...SliceOption) ([]bool, error) {
	if empty.IsNil(anyInput) {
		return nil, nil
	}
	var (
		err         error
		bb          bool
		array       []bool
		sliceOption = c.getSliceOption(option...)
	)
	switch value := anyInput.(type) {
	case []string:
		array = make([]bool, len(value))
		for k, v := range value {
			bb, err = c.Bool(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = bb
		}
	case []int:
		array = make([]bool, len(value))
		for k, v := range value {
			array[k] = v != 0
		}
	case []int8:
		array = make([]bool, len(value))
		for k, v := range value {
			array[k] = v != 0
		}
	case []int16:
		array = make([]bool, len(value))
		for k, v := range value {
			array[k] = v != 0
		}
	case []int32:
		array = make([]bool, len(value))
		for k, v := range value {
			array[k] = v != 0
		}
	case []int64:
		array = make([]bool, len(value))
		for k, v := range value {
			array[k] = v != 0
		}
	case []uint:
		array = make([]bool, len(value))
		for k, v := range value {
			array[k] = v != 0
		}
	case []uint8:
		if json.Valid(value) {
			if err = json.UnmarshalUseNumber(value, &array); array != nil {
				return array, err
			}
		}
		array = make([]bool, len(value))
		for k, v := range value {
			array[k] = v != 0
		}
	case []uint16:
		array = make([]bool, len(value))
		for k, v := range value {
			array[k] = v != 0
		}
	case []uint32:
		array = make([]bool, len(value))
		for k, v := range value {
			array[k] = v != 0
		}
	case []uint64:
		array = make([]bool, len(value))
		for k, v := range value {
			array[k] = v != 0
		}
	case []bool:
		array = value
	case []float32:
		array = make([]bool, len(value))
		for k, v := range value {
			array[k] = v != 0
		}
	case []float64:
		array = make([]bool, len(value))
		for k, v := range value {
			array[k] = v != 0
		}
	case []any:
		array = make([]bool, len(value))
		for k, v := range value {
			bb, err = c.Bool(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = bb
		}
	case [][]byte:
		array = make([]bool, len(value))
		for k, v := range value {
			bb, err = c.Bool(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = bb
		}
	case string:
		byteValue := []byte(value)
		if json.Valid(byteValue) {
			if err = json.UnmarshalUseNumber(byteValue, &array); array != nil {
				return array, err
			}
		}
		if value == "" {
			return []bool{}, err
		}
		bb, err = c.Bool(value)
		if err != nil && !sliceOption.ContinueOnError {
			return nil, err
		}
		return []bool{bb}, err
	}
	if array != nil {
		return array, err
	}
	if v, ok := anyInput.(localinterface.IInterfaces); ok {
		return c.SliceBool(v.Interfaces(), option...)
	}
	// Not a common type, it then uses reflection for conversion.
	originValueAndKind := reflection.OriginValueAndKind(anyInput)
	switch originValueAndKind.OriginKind {
	case reflect.Slice, reflect.Array:
		var (
			length = originValueAndKind.OriginValue.Len()
			slice  = make([]bool, length)
		)
		for i := 0; i < length; i++ {
			bb, err = c.Bool(originValueAndKind.OriginValue.Index(i).Interface())
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			slice[i] = bb
		}
		return slice, err

	default:
		if originValueAndKind.OriginValue.IsZero() {
			return []bool{}, err
		}
		bb, err = c.Bool(anyInput)
		if err != nil && !sliceOption.ContinueOnError {
			return nil, err
		}
		return []bool{bb}, err
	}
}
