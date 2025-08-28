// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package utils

import (
	"reflect"

	"github.com/gogf/gf/v2/internal/empty"
)

// IsNil checks whether `value` is nil, especially for any type value.
func IsNil(value any) bool {
	return empty.IsNil(value)
}

// IsEmpty checks whether `value` is empty.
func IsEmpty(value any) bool {
	return empty.IsEmpty(value)
}

// IsInt checks whether `value` is type of int.
func IsInt(value any) bool {
	switch value.(type) {
	case int, *int, int8, *int8, int16, *int16, int32, *int32, int64, *int64:
		return true
	}
	return false
}

// IsUint checks whether `value` is type of uint.
func IsUint(value any) bool {
	switch value.(type) {
	case uint, *uint, uint8, *uint8, uint16, *uint16, uint32, *uint32, uint64, *uint64:
		return true
	}
	return false
}

// IsFloat checks whether `value` is type of float.
func IsFloat(value any) bool {
	switch value.(type) {
	case float32, *float32, float64, *float64:
		return true
	}
	return false
}

// IsSlice checks whether `value` is type of slice.
func IsSlice(value any) bool {
	var (
		reflectValue = reflect.ValueOf(value)
		reflectKind  = reflectValue.Kind()
	)
	for reflectKind == reflect.Pointer {
		reflectValue = reflectValue.Elem()
		reflectKind = reflectValue.Kind()
	}
	switch reflectKind {
	case reflect.Slice, reflect.Array:
		return true
	}
	return false
}

// IsMap checks whether `value` is type of map.
func IsMap(value any) bool {
	var (
		reflectValue = reflect.ValueOf(value)
		reflectKind  = reflectValue.Kind()
	)
	for reflectKind == reflect.Pointer {
		reflectValue = reflectValue.Elem()
		reflectKind = reflectValue.Kind()
	}
	switch reflectKind {
	case reflect.Map:
		return true
	}
	return false
}

// IsStruct checks whether `value` is type of struct.
func IsStruct(value any) bool {
	reflectType := reflect.TypeOf(value)
	if reflectType == nil {
		return false
	}
	reflectKind := reflectType.Kind()
	for reflectKind == reflect.Pointer {
		reflectType = reflectType.Elem()
		reflectKind = reflectType.Kind()
	}
	switch reflectKind {
	case reflect.Struct:
		return true
	}
	return false
}
