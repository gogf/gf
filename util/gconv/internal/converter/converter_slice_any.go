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

// SliceOption is the option for Slice type converting.
type SliceOption struct {
	// ContinueOnError specifies whether to continue converting the next element
	// if one element converting fails.
	ContinueOnError bool
}

func (c *Converter) getSliceOption(option ...SliceOption) SliceOption {
	if len(option) > 0 {
		return option[0]
	}
	return SliceOption{}
}

// SliceAny converts `any` to []any.
func (c *Converter) SliceAny(any interface{}, _ ...SliceOption) ([]any, error) {
	if empty.IsNil(any) {
		return nil, nil
	}
	var (
		err   error
		array []interface{}
	)
	switch value := any.(type) {
	case []interface{}:
		array = value
	case []string:
		array = make([]interface{}, len(value))
		for k, v := range value {
			array[k] = v
		}
	case []int:
		array = make([]interface{}, len(value))
		for k, v := range value {
			array[k] = v
		}
	case []int8:
		array = make([]interface{}, len(value))
		for k, v := range value {
			array[k] = v
		}
	case []int16:
		array = make([]interface{}, len(value))
		for k, v := range value {
			array[k] = v
		}
	case []int32:
		array = make([]interface{}, len(value))
		for k, v := range value {
			array[k] = v
		}
	case []int64:
		array = make([]interface{}, len(value))
		for k, v := range value {
			array[k] = v
		}
	case []uint:
		array = make([]interface{}, len(value))
		for k, v := range value {
			array[k] = v
		}
	case []uint8:
		if json.Valid(value) {
			if err = json.UnmarshalUseNumber(value, &array); array != nil {
				return array, err
			}
		}
		array = make([]interface{}, len(value))
		for k, v := range value {
			array[k] = v
		}
	case string:
		byteValue := []byte(value)
		if json.Valid(byteValue) {
			if err = json.UnmarshalUseNumber(byteValue, &array); array != nil {
				return array, err
			}
		}

	case []uint16:
		array = make([]interface{}, len(value))
		for k, v := range value {
			array[k] = v
		}
	case []uint32:
		for _, v := range value {
			array = append(array, v)
		}
	case []uint64:
		array = make([]interface{}, len(value))
		for k, v := range value {
			array[k] = v
		}
	case []bool:
		array = make([]interface{}, len(value))
		for k, v := range value {
			array[k] = v
		}
	case []float32:
		array = make([]interface{}, len(value))
		for k, v := range value {
			array[k] = v
		}
	case []float64:
		array = make([]interface{}, len(value))
		for k, v := range value {
			array[k] = v
		}
	}
	if array != nil {
		return array, err
	}
	if v, ok := any.(localinterface.IInterfaces); ok {
		return v.Interfaces(), err
	}

	// Not a common type, it then uses reflection for conversion.
	originValueAndKind := reflection.OriginValueAndKind(any)
	switch originValueAndKind.OriginKind {
	case reflect.Slice, reflect.Array:
		var (
			length = originValueAndKind.OriginValue.Len()
			slice  = make([]interface{}, length)
		)
		for i := 0; i < length; i++ {
			slice[i] = originValueAndKind.OriginValue.Index(i).Interface()
		}
		return slice, err

	default:
		return []interface{}{any}, err
	}
}
