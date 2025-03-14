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

// SliceUint converts `any` to []uint.
func (c *Converter) SliceUint(any interface{}, option ...SliceOption) ([]uint, error) {
	if empty.IsNil(any) {
		return nil, nil
	}
	var (
		err         error
		ui          uint
		array       []uint = nil
		sliceOption        = c.getSliceOption(option...)
	)
	switch value := any.(type) {
	case []string:
		array = make([]uint, len(value))
		for k, v := range value {
			ui, err = c.Uint(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = ui
		}
	case []int8:
		array = make([]uint, len(value))
		for k, v := range value {
			array[k] = uint(v)
		}
	case []int16:
		array = make([]uint, len(value))
		for k, v := range value {
			array[k] = uint(v)
		}
	case []int32:
		array = make([]uint, len(value))
		for k, v := range value {
			array[k] = uint(v)
		}
	case []int64:
		array = make([]uint, len(value))
		for k, v := range value {
			array[k] = uint(v)
		}
	case []uint:
		array = value
	case []uint8:
		if json.Valid(value) {
			if err = json.UnmarshalUseNumber(value, &array); array != nil {
				return array, err
			}
		}
		array = make([]uint, len(value))
		for k, v := range value {
			array[k] = uint(v)
		}
	case string:
		byteValue := []byte(value)
		if json.Valid(byteValue) {
			if err = json.UnmarshalUseNumber(byteValue, &array); array != nil {
				return array, err
			}
		}
		if value == "" {
			return []uint{}, err
		}
		if utils.IsNumeric(value) {
			ui, err = c.Uint(value)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			return []uint{ui}, err
		}
	case []uint16:
		array = make([]uint, len(value))
		for k, v := range value {
			array[k] = uint(v)
		}
	case []uint32:
		array = make([]uint, len(value))
		for k, v := range value {
			array[k] = uint(v)
		}
	case []uint64:
		array = make([]uint, len(value))
		for k, v := range value {
			array[k] = uint(v)
		}
	case []bool:
		array = make([]uint, len(value))
		for k, v := range value {
			if v {
				array[k] = 1
			} else {
				array[k] = 0
			}
		}
	case []float32:
		array = make([]uint, len(value))
		for k, v := range value {
			ui, err = c.Uint(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = ui
		}
	case []float64:
		array = make([]uint, len(value))
		for k, v := range value {
			ui, err = c.Uint(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = ui
		}
	case []interface{}:
		array = make([]uint, len(value))
		for k, v := range value {
			ui, err = c.Uint(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = ui
		}
	case [][]byte:
		array = make([]uint, len(value))
		for k, v := range value {
			ui, err = c.Uint(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = ui
		}
	}

	if array != nil {
		return array, err
	}

	// Default handler.
	if v, ok := any.(localinterface.IUints); ok {
		return v.Uints(), err
	}
	if v, ok := any.(localinterface.IInterfaces); ok {
		return c.SliceUint(v.Interfaces(), option...)
	}
	// Not a common type, it then uses reflection for conversion.
	originValueAndKind := reflection.OriginValueAndKind(any)
	switch originValueAndKind.OriginKind {
	case reflect.Slice, reflect.Array:
		var (
			length = originValueAndKind.OriginValue.Len()
			slice  = make([]uint, length)
		)
		for i := 0; i < length; i++ {
			ui, err = c.Uint(originValueAndKind.OriginValue.Index(i).Interface())
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			slice[i] = ui
		}
		return slice, err

	default:
		if originValueAndKind.OriginValue.IsZero() {
			return []uint{}, err
		}
		ui, err = c.Uint(any)
		if err != nil && !sliceOption.ContinueOnError {
			return nil, err
		}
		return []uint{ui}, err
	}
}

// SliceUint32 converts `any` to []uint32.
func (c *Converter) SliceUint32(any interface{}, option ...SliceOption) ([]uint32, error) {
	if empty.IsNil(any) {
		return nil, nil
	}
	var (
		err         error
		ui          uint32
		array       []uint32 = nil
		sliceOption          = c.getSliceOption(option...)
	)
	switch value := any.(type) {
	case []string:
		array = make([]uint32, len(value))
		for k, v := range value {
			ui, err = c.Uint32(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = ui
		}
	case []int8:
		array = make([]uint32, len(value))
		for k, v := range value {
			array[k] = uint32(v)
		}
	case []int16:
		array = make([]uint32, len(value))
		for k, v := range value {
			array[k] = uint32(v)
		}
	case []int32:
		array = make([]uint32, len(value))
		for k, v := range value {
			array[k] = uint32(v)
		}
	case []int64:
		array = make([]uint32, len(value))
		for k, v := range value {
			array[k] = uint32(v)
		}
	case []uint:
		array = make([]uint32, len(value))
		for k, v := range value {
			array[k] = uint32(v)
		}
	case []uint8:
		if json.Valid(value) {
			if err = json.UnmarshalUseNumber(value, &array); array != nil {
				return array, err
			}
		}
		array = make([]uint32, len(value))
		for k, v := range value {
			array[k] = uint32(v)
		}
	case string:
		byteValue := []byte(value)
		if json.Valid(byteValue) {
			if err = json.UnmarshalUseNumber(byteValue, &array); array != nil {
				return array, err
			}
		}
		if value == "" {
			return []uint32{}, err
		}
		if utils.IsNumeric(value) {
			ui, err = c.Uint32(value)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			return []uint32{ui}, err
		}
	case []uint16:
		array = make([]uint32, len(value))
		for k, v := range value {
			array[k] = uint32(v)
		}
	case []uint32:
		array = value
	case []uint64:
		array = make([]uint32, len(value))
		for k, v := range value {
			array[k] = uint32(v)
		}
	case []bool:
		array = make([]uint32, len(value))
		for k, v := range value {
			if v {
				array[k] = 1
			} else {
				array[k] = 0
			}
		}
	case []float32:
		array = make([]uint32, len(value))
		for k, v := range value {
			ui, err = c.Uint32(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = ui
		}
	case []float64:
		array = make([]uint32, len(value))
		for k, v := range value {
			ui, err = c.Uint32(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = ui
		}
	case []interface{}:
		array = make([]uint32, len(value))
		for k, v := range value {
			ui, err = c.Uint32(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = ui
		}
	case [][]byte:
		array = make([]uint32, len(value))
		for k, v := range value {
			ui, err = c.Uint32(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = ui
		}
	}
	if array != nil {
		return array, err
	}

	// Default handler.
	if v, ok := any.(localinterface.IUints); ok {
		return c.SliceUint32(v.Uints(), option...)
	}
	if v, ok := any.(localinterface.IInterfaces); ok {
		return c.SliceUint32(v.Interfaces(), option...)
	}
	// Not a common type, it then uses reflection for conversion.
	originValueAndKind := reflection.OriginValueAndKind(any)
	switch originValueAndKind.OriginKind {
	case reflect.Slice, reflect.Array:
		var (
			length = originValueAndKind.OriginValue.Len()
			slice  = make([]uint32, length)
		)
		for i := 0; i < length; i++ {
			ui, err = c.Uint32(originValueAndKind.OriginValue.Index(i).Interface())
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			slice[i] = ui
		}
		return slice, err

	default:
		if originValueAndKind.OriginValue.IsZero() {
			return []uint32{}, err
		}
		ui, err = c.Uint32(any)
		if err != nil && !sliceOption.ContinueOnError {
			return nil, err
		}
		return []uint32{ui}, err
	}
}

// SliceUint64 converts `any` to []uint64.
func (c *Converter) SliceUint64(any interface{}, option ...SliceOption) ([]uint64, error) {
	if empty.IsNil(any) {
		return nil, nil
	}
	var (
		err         error
		ui          uint64
		array       []uint64 = nil
		sliceOption          = c.getSliceOption(option...)
	)
	switch value := any.(type) {
	case []string:
		array = make([]uint64, len(value))
		for k, v := range value {
			ui, err = c.Uint64(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = ui
		}
	case []int8:
		array = make([]uint64, len(value))
		for k, v := range value {
			array[k] = uint64(v)
		}
	case []int16:
		array = make([]uint64, len(value))
		for k, v := range value {
			array[k] = uint64(v)
		}
	case []int32:
		array = make([]uint64, len(value))
		for k, v := range value {
			array[k] = uint64(v)
		}
	case []int64:
		array = make([]uint64, len(value))
		for k, v := range value {
			array[k] = uint64(v)
		}
	case []uint:
		array = make([]uint64, len(value))
		for k, v := range value {
			array[k] = uint64(v)
		}
	case []uint8:
		if json.Valid(value) {
			if err = json.UnmarshalUseNumber(value, &array); array != nil {
				return array, err
			}
		}
		array = make([]uint64, len(value))
		for k, v := range value {
			array[k] = uint64(v)
		}
	case string:
		byteValue := []byte(value)
		if json.Valid(byteValue) {
			if err = json.UnmarshalUseNumber(byteValue, &array); array != nil {
				return array, err
			}
		}
		if value == "" {
			return []uint64{}, err
		}
		if utils.IsNumeric(value) {
			ui, err = c.Uint64(value)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			return []uint64{ui}, err
		}
	case []uint16:
		array = make([]uint64, len(value))
		for k, v := range value {
			array[k] = uint64(v)
		}
	case []uint32:
		array = make([]uint64, len(value))
		for k, v := range value {
			array[k] = uint64(v)
		}
	case []uint64:
		array = value
	case []bool:
		array = make([]uint64, len(value))
		for k, v := range value {
			if v {
				array[k] = 1
			} else {
				array[k] = 0
			}
		}
	case []float32:
		array = make([]uint64, len(value))
		for k, v := range value {
			ui, err = c.Uint64(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = ui
		}
	case []float64:
		array = make([]uint64, len(value))
		for k, v := range value {
			ui, err = c.Uint64(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = ui
		}
	case []interface{}:
		array = make([]uint64, len(value))
		for k, v := range value {
			ui, err = c.Uint64(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = ui
		}
	case [][]byte:
		array = make([]uint64, len(value))
		for k, v := range value {
			ui, err = c.Uint64(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = ui
		}
	}
	if array != nil {
		return array, err
	}
	// Default handler.
	if v, ok := any.(localinterface.IUints); ok {
		return c.SliceUint64(v.Uints(), option...)
	}
	if v, ok := any.(localinterface.IInterfaces); ok {
		return c.SliceUint64(v.Interfaces(), option...)
	}
	// Not a common type, it then uses reflection for conversion.
	originValueAndKind := reflection.OriginValueAndKind(any)
	switch originValueAndKind.OriginKind {
	case reflect.Slice, reflect.Array:
		var (
			length = originValueAndKind.OriginValue.Len()
			slice  = make([]uint64, length)
		)
		for i := 0; i < length; i++ {
			ui, err = c.Uint64(originValueAndKind.OriginValue.Index(i).Interface())
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			slice[i] = ui
		}
		return slice, err

	default:
		if originValueAndKind.OriginValue.IsZero() {
			return []uint64{}, err
		}
		ui, err = c.Uint64(any)
		if err != nil && !sliceOption.ContinueOnError {
			return nil, err
		}
		return []uint64{ui}, err
	}
}
