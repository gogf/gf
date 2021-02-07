// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package structs

import (
	"github.com/gogf/gf/errors/gerror"
	"reflect"
)

// StructType retrieves and returns the Type of specified struct/*struct.
func StructType(structOrPointer interface{}) (*Type, error) {
	var (
		reflectValue reflect.Value
		reflectKind  reflect.Kind
		reflectType  reflect.Type
	)
	if rv, ok := structOrPointer.(reflect.Value); ok {
		reflectValue = rv
	} else {
		reflectValue = reflect.ValueOf(structOrPointer)
	}
	reflectKind = reflectValue.Kind()
	for reflectKind == reflect.Ptr {
		if !reflectValue.IsValid() || reflectValue.IsNil() {
			// If pointer is type of *struct and nil, then automatically create a temporary struct.
			reflectValue = reflect.New(reflectValue.Type().Elem()).Elem()
			reflectKind = reflectValue.Kind()
		} else {
			reflectValue = reflectValue.Elem()
			reflectKind = reflectValue.Kind()
		}
	}
	for reflectKind == reflect.Ptr {
		reflectValue = reflectValue.Elem()
		reflectKind = reflectValue.Kind()
	}
	if reflectKind != reflect.Struct {
		return nil, gerror.Newf(
			`invalid parameter kind "%s", kind of "struct" is required`,
			reflectKind,
		)
	}
	reflectType = reflectValue.Type()
	return &Type{
		Type: reflectType,
	}, nil
}

func (t *Type) Signature() string {
	return t.PkgPath() + "/" + t.String()
}
