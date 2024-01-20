// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gutil

import (
	"reflect"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gstructs"
	"github.com/gogf/gf/v2/util/gconv"
)

// StructToSlice converts struct to slice of which all keys and values are its items.
// Eg: {"K1": "v1", "K2": "v2"} => ["K1", "v1", "K2", "v2"]
func StructToSlice(data interface{}) []interface{} {
	var (
		reflectValue = reflect.ValueOf(data)
		reflectKind  = reflectValue.Kind()
	)
	for reflectKind == reflect.Ptr {
		reflectValue = reflectValue.Elem()
		reflectKind = reflectValue.Kind()
	}
	switch reflectKind {
	case reflect.Struct:
		array := make([]interface{}, 0)
		// Note that, it uses the gconv tag name instead of the attribute name if
		// the gconv tag is fined in the struct attributes.
		for k, v := range gconv.Map(reflectValue) {
			array = append(array, k)
			array = append(array, v)
		}
		return array
	}
	return nil
}

// FillStructWithDefault fills  attributes of pointed struct with tag value from `default/d` tag .
// The parameter `structPtr` should be either type of *struct/[]*struct.
func FillStructWithDefault(structPtr interface{}) error {
	var (
		reflectValue reflect.Value
	)
	if rv, ok := structPtr.(reflect.Value); ok {
		reflectValue = rv
	} else {
		reflectValue = reflect.ValueOf(structPtr)
	}
	switch reflectValue.Kind() {
	case reflect.Ptr:
		// Nothing to do.
	case reflect.Array, reflect.Slice:
		if reflectValue.Elem().Kind() != reflect.Ptr {
			return gerror.NewCodef(
				gcode.CodeInvalidParameter,
				`invalid parameter "%s", the element of slice should be type of pointer of struct, but given "%s"`,
				reflectValue.Type().String(), reflectValue.Elem().Type().String(),
			)
		}
	default:
		return gerror.NewCodef(
			gcode.CodeInvalidParameter,
			`invalid parameter "%s", should be type of pointer of struct`,
			reflectValue.Type().String(),
		)
	}
	if reflectValue.IsNil() {
		return gerror.NewCode(
			gcode.CodeInvalidParameter,
			`the pointed struct object should not be nil`,
		)
	}
	if !reflectValue.Elem().IsValid() {
		return gerror.NewCode(
			gcode.CodeInvalidParameter,
			`the pointed struct object should be valid`,
		)
	}
	fields, err := gstructs.Fields(gstructs.FieldsInput{
		Pointer:         reflectValue,
		RecursiveOption: gstructs.RecursiveOptionEmbedded,
	})
	if err != nil {
		return err
	}
	for _, field := range fields {
		if field.OriginalKind() == reflect.Struct {
			err := FillStructWithDefault(field.OriginalValue().Addr())
			if err != nil {
				return err
			}
			continue
		}

		if defaultValue := field.TagDefault(); defaultValue != "" {
			if field.IsEmpty() {
				field.Value.Set(reflect.ValueOf(
					gconv.ConvertWithRefer(defaultValue, field.Value),
				))
			}
		}
	}

	return nil
}
