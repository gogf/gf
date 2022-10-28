// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gredis

import (
	"reflect"

	"github.com/gogf/gf/v2/os/gstructs"
	"github.com/gogf/gf/v2/util/gutil"
)

func mustMergeOptionToArgs(args []interface{}, option interface{}) []interface{} {
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
		args        = make([]interface{}, 0)
		fields, err = gstructs.Fields(gstructs.FieldsInput{
			Pointer:         option,
			RecursiveOption: gstructs.RecursiveOptionEmbeddedNoTag,
		})
	)
	if err != nil {
		return nil, err
	}
	for _, field := range fields {
		switch field.OriginalKind() {
		// See SetOption
		case reflect.Bool:
			args = append(args, field.Name())

		// See ZRangeOption
		case reflect.Struct:
			args = append(args, field.Name())
			args = append(args, gutil.Values(field.Value.Interface())...)

		// See TTLOption
		default:
			args = append(args, field.Name(), field.Value.Interface())
		}
	}
	return args, nil
}
