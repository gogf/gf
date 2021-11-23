// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//

package gcmd

import (
	"reflect"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/structs"
	"github.com/gogf/gf/v2/internal/utils"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gmeta"
	"github.com/gogf/gf/v2/util/gutil"
)

const (
	tagNameDc = `dc`
	tagNameAd = `ad`
)

func CommandsFromObject(object interface{}) (commands []Command, err error) {
	originValueAndKind := utils.OriginValueAndKind(object)
	if originValueAndKind.OriginKind != reflect.Struct {
		return nil, gerror.Newf(
			`input object should be type of struct, but got "%s"`,
			originValueAndKind.InputValue.Type().String(),
		)
	}
	//for i := 0; i < originValueAndKind.InputValue.NumMethod(); i++ {
	//	method := originValueAndKind.InputValue.Method(i)
	//}
	//for _, field := range fields {
	//
	//}
	return
}

func newCommandFromMethod(object, method reflect.Value) (*Command, error) {
	var (
		err         error
		reflectType = method.Type()
	)
	// Necessary validation for input/output parameters and naming.
	if reflectType.NumIn() != 2 || reflectType.NumOut() != 2 {
		if reflectType.PkgPath() != "" {
			err = gerror.NewCodef(
				gcode.CodeInvalidParameter,
				`invalid handler: %s.%s.%s defined as "%s", but "func(context.Context, Input)(Output, error)" is required`,
				reflectType.PkgPath(), object.Type().Name(), reflectType.Name(), reflectType.String(),
			)
		} else {
			err = gerror.NewCodef(
				gcode.CodeInvalidParameter,
				`invalid handler: defined as "%s", but "func(context.Context, Input)(Output, error)" is required`,
				reflectType.String(),
			)
		}
		return nil, err
	}
	if reflectType.In(0).String() != "context.Context" {
		err = gerror.NewCodef(
			gcode.CodeInvalidParameter,
			`invalid handler: defined as "%s", but the first input parameter should be type of "context.Context"`,
			reflectType.String(),
		)
		return nil, err
	}
	if reflectType.Out(1).String() != "error" {
		err = gerror.NewCodef(
			gcode.CodeInvalidParameter,
			`invalid handler: defined as "%s", but the last output parameter should be type of "error"`,
			reflectType.String(),
		)
		return nil, err
	}
	// The input struct should be named as `xxxInput`.
	if !gstr.HasSuffix(reflectType.In(1).String(), `Input`) {
		err = gerror.NewCodef(
			gcode.CodeInvalidParameter,
			`invalid struct naming for input: defined as "%s", but it should be named with "Input" suffix like "xxxInput"`,
			reflectType.In(1).String(),
		)
		return nil, err
	}
	// The output struct should be named as `xxxOutput`.
	if !gstr.HasSuffix(reflectType.Out(0).String(), `Output`) {
		err = gerror.NewCodef(
			gcode.CodeInvalidParameter,
			`invalid struct naming for output: defined as "%s", but it should be named with "Output" suffix like "xxxOutput"`,
			reflectType.Out(0).String(),
		)
		return nil, err
	}

	var (
		inputObject  reflect.Value
		outputObject reflect.Value
	)
	if method.Type().In(1).Kind() == reflect.Ptr {
		inputObject = reflect.New(method.Type().In(1).Elem()).Elem()
	} else {
		inputObject = reflect.New(method.Type().In(1)).Elem()
	}

	if method.Type().Out(1).Kind() == reflect.Ptr {
		outputObject = reflect.New(method.Type().Out(0).Elem()).Elem()
	} else {
		outputObject = reflect.New(method.Type().Out(0)).Elem()
	}
	// Command creating.
	var (
		cmd      = Command{}
		metaData = gmeta.Data(inputObject.Interface())
	)
	if err = gconv.Scan(metaData, &cmd); err != nil {
		return nil, err
	}
	// Name filed is necessary.
	if cmd.Name == "" {
		return nil, gerror.Newf(
			`command name cannot be empty, "name" tag not found in struct "%s"`,
			inputObject.Type().String(),
		)
	}
	if cmd.Description == "" {
		cmd.Description = metaData[tagNameDc]
	}
	if cmd.Additional == "" {
		cmd.Additional = metaData[tagNameAd]
	}

	if cmd.Options, err = newOptionsFromInput(inputObject.Interface()); err != nil {
		return nil, err
	}
	return &cmd, nil
}

func newOptionsFromInput(object interface{}) (options []Option, err error) {
	var (
		fields []structs.Field
	)
	fields, err = structs.Fields(structs.FieldsInput{
		Pointer:         object,
		RecursiveOption: structs.RecursiveOptionEmbeddedNoTag,
	})
	for _, field := range fields {
		var (
			option   = Option{}
			metaData = gmeta.Data(field.Value.Interface())
		)
		if err = gconv.Scan(metaData, &option); err != nil {
			return nil, err
		}
		if option.Name == "" {
			option.Name = field.Name()
		}
		options = append(options, option)
	}
	return
}
