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

// StructType retrieves and returns the struct Type of specified struct/*struct.
// The parameter `object` should be either type of struct/*struct/[]struct/[]*struct.
func StructType(object interface{}) (*Type, error) {
	var (
		reflectValue reflect.Value
		reflectKind  reflect.Kind
		reflectType  reflect.Type
	)
	if rv, ok := object.(reflect.Value); ok {
		reflectValue = rv
	} else {
		reflectValue = reflect.ValueOf(object)
	}
	reflectKind = reflectValue.Kind()
	for {
		switch reflectKind {
		case reflect.Ptr:
			if !reflectValue.IsValid() || reflectValue.IsNil() {
				// If pointer is type of *struct and nil, then automatically create a temporary struct.
				reflectValue = reflect.New(reflectValue.Type().Elem()).Elem()
				reflectKind = reflectValue.Kind()
			} else {
				reflectValue = reflectValue.Elem()
				reflectKind = reflectValue.Kind()
			}
		case reflect.Array, reflect.Slice:
			reflectValue = reflect.New(reflectValue.Type().Elem()).Elem()
			reflectKind = reflectValue.Kind()
		default:
			goto exitLoop
		}
	}
exitLoop:
	for reflectKind == reflect.Ptr {
		reflectValue = reflectValue.Elem()
		reflectKind = reflectValue.Kind()
	}
	if reflectKind != reflect.Struct {
		return nil, gerror.Newf(
			`invalid object kind "%s", kind of "struct" is required`,
			reflectKind,
		)
	}
	reflectType = reflectValue.Type()
	return &Type{
		Type: reflectType,
	}, nil
}

// Signature returns an unique string as this type.
func (t Type) Signature() string {
	return t.PkgPath() + "/" + t.String()
}

// FieldKeys returns the keys of current struct/map.
func (t Type) FieldKeys() []string {
	keys := make([]string, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		keys[i] = t.Field(i).Name
	}
	return keys
}
