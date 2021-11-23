// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//

package gcmd

import (
	"context"
	"reflect"

	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/structs"
	"github.com/gogf/gf/v2/internal/utils"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gmeta"
	"github.com/gogf/gf/v2/util/gutil"
	"github.com/gogf/gf/v2/util/gvalid"
)

const (
	tagNameDc   = `dc`
	tagNameAd   = `ad`
	tagNameRoot = `root`
)

var (
	// defaultValueTags is the struct tag names for default value storing.
	defaultValueTags = []string{"d", "default"}
)

// NewFromObject creates and returns a root command object using given object.
func NewFromObject(object interface{}) (rootCmd *Command, err error) {
	originValueAndKind := utils.OriginValueAndKind(object)
	if originValueAndKind.OriginKind != reflect.Struct {
		return nil, gerror.Newf(
			`input object should be type of struct, but got "%s"`,
			originValueAndKind.InputValue.Type().String(),
		)
	}
	var (
		nameSet     = gset.NewStrSet()
		subCommands []Command
	)
	for i := 0; i < originValueAndKind.InputValue.NumMethod(); i++ {
		var (
			root          bool
			method        = originValueAndKind.InputValue.Method(i)
			methodCommand Command
		)
		methodCommand, root, err = newCommandFromMethod(object, method)
		if err != nil {
			return nil, err
		}
		if nameSet.Contains(methodCommand.Name) {
			return nil, gerror.Newf(
				`command name should be unique, found duplicated command name in method "%s"`,
				method.Type().String(),
			)
		}
		if root {
			if rootCmd != nil {
				return nil, gerror.Newf(
					`there should be only one root command in object, found duplicated in method "%s"`,
					method.Type().String(),
				)
			}
			rootCmd = &methodCommand
		} else {
			subCommands = append(subCommands, methodCommand)
		}
	}
	if rootCmd == nil {
		return nil, gerror.Newf(
			`there should be one root command in object when creating command from object, but found none in object "%s"`,
			originValueAndKind.InputValue.Type().String(),
		)
	}
	if len(subCommands) > 0 {
		err = rootCmd.AddCommand(subCommands...)
	}
	return
}

func newCommandFromMethod(object interface{}, method reflect.Value) (command Command, root bool, err error) {
	var (
		reflectType = method.Type()
	)
	// Necessary validation for input/output parameters and naming.
	if reflectType.NumIn() != 2 || reflectType.NumOut() != 2 {
		if reflectType.PkgPath() != "" {
			err = gerror.NewCodef(
				gcode.CodeInvalidParameter,
				`invalid command: %s.%s.%s defined as "%s", but "func(context.Context, Input)(Output, error)" is required`,
				reflectType.PkgPath(), reflect.TypeOf(object).Name(), reflectType.Name(), reflectType.String(),
			)
		} else {
			err = gerror.NewCodef(
				gcode.CodeInvalidParameter,
				`invalid command: defined as "%s", but "func(context.Context, Input)(Output, error)" is required`,
				reflectType.String(),
			)
		}
		return
	}
	if reflectType.In(0).String() != "context.Context" {
		err = gerror.NewCodef(
			gcode.CodeInvalidParameter,
			`invalid command: defined as "%s", but the first input parameter should be type of "context.Context"`,
			reflectType.String(),
		)
		return
	}
	if reflectType.Out(1).String() != "error" {
		err = gerror.NewCodef(
			gcode.CodeInvalidParameter,
			`invalid command: defined as "%s", but the last output parameter should be type of "error"`,
			reflectType.String(),
		)
		return
	}
	// The input struct should be named as `xxxInput`.
	if !gstr.HasSuffix(reflectType.In(1).String(), `Input`) {
		err = gerror.NewCodef(
			gcode.CodeInvalidParameter,
			`invalid struct naming for input: defined as "%s", but it should be named with "Input" suffix like "xxxInput"`,
			reflectType.In(1).String(),
		)
		return
	}
	// The output struct should be named as `xxxOutput`.
	if !gstr.HasSuffix(reflectType.Out(0).String(), `Output`) {
		err = gerror.NewCodef(
			gcode.CodeInvalidParameter,
			`invalid struct naming for output: defined as "%s", but it should be named with "Output" suffix like "xxxOutput"`,
			reflectType.Out(0).String(),
		)
		return
	}

	var (
		inputObject reflect.Value
	)
	if method.Type().In(1).Kind() == reflect.Ptr {
		inputObject = reflect.New(method.Type().In(1).Elem()).Elem()
	} else {
		inputObject = reflect.New(method.Type().In(1)).Elem()
	}

	// Command creating.
	var (
		metaData = gmeta.Data(inputObject.Interface())
	)
	if err = gconv.Scan(metaData, &command); err != nil {
		return
	}
	root = gconv.Bool(metaData[tagNameRoot])
	// Name filed is necessary.
	if command.Name == "" {
		err = gerror.Newf(
			`command name cannot be empty, "name" tag not found in struct "%s"`,
			inputObject.Type().String(),
		)
		return
	}
	if command.Description == "" {
		command.Description = metaData[tagNameDc]
	}
	if command.Additional == "" {
		command.Additional = metaData[tagNameAd]
	}

	if command.Options, err = newOptionsFromInput(inputObject.Interface()); err != nil {
		return
	}

	// Create function that has value return.
	command.FuncWithValue = func(ctx context.Context, parser *Parser) (out interface{}, err error) {
		defer func() {
			if exception := recover(); exception != nil {
				if v, ok := exception.(error); ok && gerror.HasStack(v) {
					err = v
				} else {
					err = gerror.New(`exception recovered:` + gconv.String(exception))
				}
			}
		}()

		var (
			data        = gconv.Map(parser.GetOptAll())
			inputValues = []reflect.Value{reflect.ValueOf(ctx)}
		)
		if data == nil {
			data = map[string]interface{}{}
		}
		err = mergeDefaultStructValue(data, inputObject.Interface())
		if err != nil {
			return nil, err
		}
		// Construct input parameters.
		if len(data) > 0 {

		}
		if inputObject.Kind() == reflect.Ptr {
			err = gconv.Scan(data, inputObject.Interface())
		} else {
			err = gconv.Struct(data, inputObject.Addr().Interface())
		}
		if err != nil {
			return
		}

		// Parameters validation.
		if err = gvalid.New().Bail().Data(inputObject.Interface()).Assoc(data).Run(ctx); err != nil {
			err = gerror.Current(err)
			return
		}
		inputValues = append(inputValues, inputObject)

		// Call handler with dynamic created parameter values.
		results := method.Call(inputValues)
		out = results[0].Interface()
		if !results[1].IsNil() {
			if v, ok := results[1].Interface().(error); ok {
				err = v
			}
		}
		return
	}
	return
}

// mergeDefaultStructValue merges the request parameters with default values from struct tag definition.
func mergeDefaultStructValue(data map[string]interface{}, pointer interface{}) error {
	tagFields, err := structs.TagFields(pointer, defaultValueTags)
	if err != nil {
		return err
	}
	if len(tagFields) > 0 {
		var (
			foundKey   string
			foundValue interface{}
		)
		for _, field := range tagFields {
			foundKey, foundValue = gutil.MapPossibleItemByKey(data, field.Name())
			if foundKey == "" {
				data[field.Name()] = field.TagValue
			} else {
				if utils.IsEmpty(foundValue) {
					data[foundKey] = field.TagValue
				}
			}
		}
	}
	return nil
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
			metaData = field.TagMap()
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
