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

// SliceStr converts `any` to []string.
func (c *Converter) SliceStr(any interface{}, option ...SliceOption) ([]string, error) {
	if empty.IsNil(any) {
		return nil, nil
	}
	var (
		err         error
		s           string
		array       []string = nil
		sliceOption          = c.getSliceOption(option...)
	)
	switch value := any.(type) {
	case []int:
		array = make([]string, len(value))
		for k, v := range value {
			s, err = c.String(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = s
		}
	case []int8:
		array = make([]string, len(value))
		for k, v := range value {
			s, err = c.String(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = s
		}
	case []int16:
		array = make([]string, len(value))
		for k, v := range value {
			s, err = c.String(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = s
		}
	case []int32:
		array = make([]string, len(value))
		for k, v := range value {
			s, err = c.String(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = s
		}
	case []int64:
		array = make([]string, len(value))
		for k, v := range value {
			s, err = c.String(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = s
		}
	case []uint:
		array = make([]string, len(value))
		for k, v := range value {
			s, err = c.String(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = s
		}
	case []uint8:
		if json.Valid(value) {
			if err = json.UnmarshalUseNumber(value, &array); array != nil {
				return array, err
			}
		}
		array = make([]string, len(value))
		for k, v := range value {
			s, err = c.String(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = s
		}
		return array, err
	case string:
		byteValue := []byte(value)
		if json.Valid(byteValue) {
			if err = json.UnmarshalUseNumber(byteValue, &array); array != nil {
				return array, err
			}
		}
		if value == "" {
			return []string{}, err
		}
		return []string{value}, err
	case []uint16:
		array = make([]string, len(value))
		for k, v := range value {
			s, err = c.String(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = s
		}
	case []uint32:
		array = make([]string, len(value))
		for k, v := range value {
			s, err = c.String(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = s
		}
	case []uint64:
		array = make([]string, len(value))
		for k, v := range value {
			s, err = c.String(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = s
		}
	case []bool:
		array = make([]string, len(value))
		for k, v := range value {
			s, err = c.String(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = s
		}
	case []float32:
		array = make([]string, len(value))
		for k, v := range value {
			s, err = c.String(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = s
		}
	case []float64:
		array = make([]string, len(value))
		for k, v := range value {
			s, err = c.String(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = s
		}
	case []interface{}:
		array = make([]string, len(value))
		for k, v := range value {
			s, err = c.String(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = s
		}
	case []string:
		array = value
	case [][]byte:
		array = make([]string, len(value))
		for k, v := range value {
			s, err = c.String(v)
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			array[k] = s
		}
	}
	if array != nil {
		return array, err
	}
	if v, ok := any.(localinterface.IStrings); ok {
		return v.Strings(), err
	}
	if v, ok := any.(localinterface.IInterfaces); ok {
		return c.SliceStr(v.Interfaces(), option...)
	}
	// Not a common type, it then uses reflection for conversion.
	originValueAndKind := reflection.OriginValueAndKind(any)
	switch originValueAndKind.OriginKind {
	case reflect.Slice, reflect.Array:
		var (
			length = originValueAndKind.OriginValue.Len()
			slice  = make([]string, length)
		)
		for i := 0; i < length; i++ {
			s, err = c.String(originValueAndKind.OriginValue.Index(i).Interface())
			if err != nil && !sliceOption.ContinueOnError {
				return nil, err
			}
			slice[i] = s
		}
		return slice, err

	default:
		if originValueAndKind.OriginValue.IsZero() {
			return []string{}, err
		}
		s, err = c.String(any)
		if err != nil && !sliceOption.ContinueOnError {
			return nil, err
		}
		return []string{s}, err
	}
}
