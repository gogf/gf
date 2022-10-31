// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gredis

import (
	"reflect"

	"github.com/gogf/gf/v2/os/gstructs"
)

func mustMergeOptionToArgs(args []interface{}, option interface{}) []interface{} {
	if option == nil {
		return args
	}
	var (
		err        error
		optionArgs []interface{}
	)
	optionArgs, err = convertOptionToArgs(option)
	if err != nil {
		panic(err)
	}
	return append(args, optionArgs...)
}

func convertOptionToArgs(option interface{}) ([]interface{}, error) {
	if option == nil {
		return nil, nil
	}
	var (
		err       error
		args      = make([]interface{}, 0)
		fields    []gstructs.Field
		subFields []gstructs.Field
	)
	fields, err = gstructs.Fields(gstructs.FieldsInput{
		Pointer:         option,
		RecursiveOption: gstructs.RecursiveOptionEmbeddedNoTag,
	})
	if err != nil {
		return nil, err
	}
	for _, field := range fields {
		switch field.OriginalKind() {
		// See SetOption
		case reflect.Bool:
			if field.Value.Bool() {
				args = append(args, field.Name())
			}

		// See ZRangeOption
		case reflect.Struct:
			if field.Value.IsNil() {
				continue
			}
			if !field.IsEmbedded() {
				args = append(args, field.Name())
			}
			subFields, err = gstructs.Fields(gstructs.FieldsInput{
				Pointer:         option,
				RecursiveOption: gstructs.RecursiveOptionEmbeddedNoTag,
			})
			if err != nil {
				return nil, err
			}
			for _, subField := range subFields {
				args = append(args, subField.Value.Interface())
			}

		// See TTLOption
		default:
			fieldValue := field.Value.Interface()
			if field.Value.Kind() == reflect.Ptr {
				if field.Value.IsNil() {
					continue
				}
				fieldValue = field.Value.Elem().Interface()
			}
			args = append(args, field.Name(), fieldValue)
		}
	}
	return args, nil
}
