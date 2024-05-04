// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gstructs

import "reflect"

// Signature returns a unique string as this type.
func (t Type) Signature() string {
	return t.PkgPath() + "/" + t.String()
}

// FieldKeys returns the keys of current struct.
func (t Type) FieldKeys() []string {
	if t.Kind() != reflect.Struct {
		return []string{}
	}
	return getStructFields(t.Type)
}

func getStructFields(structType reflect.Type) []string {
	keys := make([]string, 0, structType.NumField())
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		if field.Anonymous {
			if field.Type.Kind() == reflect.Ptr {
				field.Type = field.Type.Elem()
			}
			if field.Type.Kind() == reflect.Struct {
				keys = append(keys, getStructFields(field.Type)...)
				continue
			}
		}
		keys = append(keys, field.Name)
	}
	return keys
}
