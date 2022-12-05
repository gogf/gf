/**
// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
*/

package gconv

import (
	"github.com/shopspring/decimal"
	"reflect"

	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/internal/reflection"
)

// SliceDecimal is alias of Decimals.
func SliceDecimal(any interface{}) []decimal.Decimal {
	return Decimals(any)
}

// Decimals converts `any` to []decimal.Decimal.
func Decimals(any interface{}) []decimal.Decimal {
	if any == nil {
		return nil
	}
	var (
		array []decimal.Decimal = nil
	)
	switch value := any.(type) {
	case string:
		if value == "" {
			return []decimal.Decimal{}
		}
		return []decimal.Decimal{Decimal(value)}
	case []string:
		array = make([]decimal.Decimal, len(value))
		for k, v := range value {
			array[k] = Decimal(v)
		}
	case []int:
		array = make([]decimal.Decimal, len(value))
		for k, v := range value {
			array[k] = Decimal(v)
		}
	case []int8:
		array = make([]decimal.Decimal, len(value))
		for k, v := range value {
			array[k] = Decimal(v)
		}
	case []int16:
		array = make([]decimal.Decimal, len(value))
		for k, v := range value {
			array[k] = Decimal(v)
		}
	case []int32:
		array = make([]decimal.Decimal, len(value))
		for k, v := range value {
			array[k] = Decimal(v)
		}
	case []int64:
		array = make([]decimal.Decimal, len(value))
		for k, v := range value {
			array[k] = Decimal(v)
		}
	case []uint:
		for _, v := range value {
			array = append(array, Decimal(v))
		}
	case []uint8:
		if json.Valid(value) {
			var floatArray []float64 = nil
			_ = json.UnmarshalUseNumber(value, &floatArray)
			array = make([]decimal.Decimal, len(floatArray))
			for k, v := range floatArray {
				array[k] = Decimal(v)
			}
		} else {
			array = make([]decimal.Decimal, len(value))
			for k, v := range value {
				array[k] = Decimal(v)
			}
		}
	case []uint16:
		array = make([]decimal.Decimal, len(value))
		for k, v := range value {
			array[k] = Decimal(v)
		}
	case []uint32:
		array = make([]decimal.Decimal, len(value))
		for k, v := range value {
			array[k] = Decimal(v)
		}
	case []uint64:
		array = make([]decimal.Decimal, len(value))
		for k, v := range value {
			array[k] = Decimal(v)
		}
	case []bool:
		array = make([]decimal.Decimal, len(value))
		for k, v := range value {
			array[k] = Decimal(v)
		}
	case []float32:
		array = make([]decimal.Decimal, len(value))
		for k, v := range value {
			array[k] = Decimal(v)
		}
	case []float64:
		array = make([]decimal.Decimal, len(value))
		for k, v := range value {
			array[k] = Decimal(v)
		}
	case []interface{}:
		array = make([]decimal.Decimal, len(value))
		for k, v := range value {
			array[k] = Decimal(v)
		}
	}
	if array != nil {
		return array
	}
	if v, ok := any.(iDecimals); ok {
		return v.Decimals()
	}
	if v, ok := any.(iInterfaces); ok {
		return Decimals(v.Interfaces())
	}
	// JSON format string value converting.
	if checkJsonAndUnmarshalUseNumber(any, &array) {
		return array
	}
	// Not a common type, it then uses reflection for conversion.
	originValueAndKind := reflection.OriginValueAndKind(any)
	switch originValueAndKind.OriginKind {
	case reflect.Slice, reflect.Array:
		var (
			length = originValueAndKind.OriginValue.Len()
			slice  = make([]decimal.Decimal, length)
		)
		for i := 0; i < length; i++ {
			slice[i] = Decimal(originValueAndKind.OriginValue.Index(i).Interface())
		}
		return slice

	default:
		if originValueAndKind.OriginValue.IsZero() {
			return []decimal.Decimal{}
		}
		return []decimal.Decimal{Decimal(any)}
	}
}
