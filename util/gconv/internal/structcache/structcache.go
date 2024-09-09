// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package structcache provides struct and field info cache feature to enhance performance for struct converting.
package structcache

import (
	"reflect"

	"github.com/gogf/gf/v2/util/gconv/internal/localinterface"
)

var (
	// customConvertTypeMap is used to store whether field types are registered to custom conversions
	// For example:
	// func (src *TypeA) (dst *TypeB,err error)
	// This map will store `TypeB` for quick judgment during assignment.
	customConvertTypeMap = map[reflect.Type]struct{}{}
)

// RegisterCustomConvertType registers custom
func RegisterCustomConvertType(fieldType reflect.Type) {
	if fieldType.Kind() == reflect.Ptr {
		fieldType = fieldType.Elem()
	}
	customConvertTypeMap[fieldType] = struct{}{}
}

var (
	implUnmarshalText  = reflect.TypeOf((*localinterface.IUnmarshalText)(nil)).Elem()
	implUnmarshalJson  = reflect.TypeOf((*localinterface.IUnmarshalJSON)(nil)).Elem()
	implUnmarshalValue = reflect.TypeOf((*localinterface.IUnmarshalValue)(nil)).Elem()
)

func checkTypeIsCommonInterface(field reflect.StructField) bool {
	isCommonInterface := false
	switch field.Type.String() {
	case "time.Time", "*time.Time":
		// default convert.

	case "gtime.Time", "*gtime.Time":
		// default convert.

	default:
		// Implemented three types of interfaces that must be pointer types, otherwise it is meaningless
		if field.Type.Kind() != reflect.Ptr {
			field.Type = reflect.PointerTo(field.Type)
		}
		switch {
		case field.Type.Implements(implUnmarshalText):
			isCommonInterface = true

		case field.Type.Implements(implUnmarshalJson):
			isCommonInterface = true

		case field.Type.Implements(implUnmarshalValue):
			isCommonInterface = true
		}
	}
	return isCommonInterface
}
