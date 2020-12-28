// Copyright GoFrame Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gvar

import (
	"github.com/gogf/gf/internal/empty"
	"reflect"
)

// IsNil checks whether <v> is nil.
func (v *Var) IsNil() bool {
	return v.Val() == nil
}

// IsEmpty checks whether <v> is empty.
func (v *Var) IsEmpty() bool {
	return empty.IsEmpty(v.Val())
}

// IsInt checks whether <v> is type of int.
func (v *Var) IsInt() bool {
	switch v.Val().(type) {
	case int, *int, int8, *int8, int16, *int16, int32, *int32, int64, *int64:
		return true
	}
	return false
}

// IsUint checks whether <v> is type of uint.
func (v *Var) IsUint() bool {
	switch v.Val().(type) {
	case uint, *uint, uint8, *uint8, uint16, *uint16, uint32, *uint32, uint64, *uint64:
		return true
	}
	return false
}

// IsFloat checks whether <v> is type of float.
func (v *Var) IsFloat() bool {
	switch v.Val().(type) {
	case float32, *float32, float64, *float64:
		return true
	}
	return false
}

// IsSlice checks whether <v> is type of slice.
func (v *Var) IsSlice() bool {
	var (
		reflectValue = reflect.ValueOf(v.Val())
		reflectKind  = reflectValue.Kind()
	)
	for reflectKind == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}
	switch reflectKind {
	case reflect.Slice, reflect.Array:
		return true
	}
	return false
}

// IsMap checks whether <v> is type of map.
func (v *Var) IsMap() bool {
	var (
		reflectValue = reflect.ValueOf(v.Val())
		reflectKind  = reflectValue.Kind()
	)
	for reflectKind == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}
	switch reflectKind {
	case reflect.Map:
		return true
	}
	return false
}

// IsStruct checks whether <v> is type of struct.
func (v *Var) IsStruct() bool {
	var (
		reflectValue = reflect.ValueOf(v.Val())
		reflectKind  = reflectValue.Kind()
	)
	for reflectKind == reflect.Ptr {
		reflectValue = reflectValue.Elem()
		reflectKind = reflectValue.Kind()
	}
	switch reflectKind {
	case reflect.Struct:
		return true
	}
	return false
}
