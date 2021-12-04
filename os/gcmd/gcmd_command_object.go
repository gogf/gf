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
	"github.com/gogf/gf/v2/internal/utils"
	"github.com/gogf/gf/v2/os/gstructs"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gmeta"
	"github.com/gogf/gf/v2/util/gutil"
	"github.com/gogf/gf/v2/util/gvalid"
)

const (
	tagNameDc   = `dc`
	tagNameAd   = `ad`
	tagNameEg   = `eg`
	tagNameArgs = `args`
	tagNameRoot = `root`
)

var (
	// defaultValueTags is the struct tag names for default value storing.
	defaultValueTags = []string{"d", "default"}
)

// NewFromObject creates and returns a root command object using given object.
func NewFromObject(object interface{}) (rootCmd Command, err error) {
	originValueAndKind := utils.OriginValueAndKind(object)
	if originValueAndKind.OriginKind != reflect.Struct {
		err = gerror.Newf(
			`input object should be type of struct, but got "%s"`,
			originValueAndKind.InputValue.Type().String(),
		)
		return
	}
	// Root command creating.
	rootCmd, err = newCommandFromObjectMeta(object)
	if err != nil {
		return
	}
	// Sub command creating.
	var (
		nameSet         = gset.NewStrSet()
		rootCommandName = gmeta.Get(object, tagNameRoot).String()
		subCommands     []Command
	)
	for i := 0; i < originValueAndKind.InputValue.NumMethod(); i++ {
		var (
			method        = originValueAndKind.InputValue.Method(i)
			methodCommand Command
		)
		methodCommand, err = newCommandFromMethod(object, method)
		if err != nil {
			return
		}
		if nameSet.Contains(methodCommand.Name) {
			err = gerror.Newf(
				`command name should be unique, found duplicated command name in method "%s"`,
				method.Type().String(),
			)
			return
		}
		if rootCommandName == methodCommand.Name {
			if rootCmd.Func == nil {
				rootCmd.Func = methodCommand.Func
			}
			if rootCmd.FuncWithValue == nil {
				rootCmd.FuncWithValue = methodCommand.FuncWithValue
			}
			if len(rootCmd.Options) == 0 {
				rootCmd.Options = methodCommand.Options
			}
		} else {
			subCommands = append(subCommands, methodCommand)
		}
	}
	if len(subCommands) > 0 {
		err = rootCmd.AddCommand(subCommands...)
	}
	return
}

func newCommandFromObjectMeta(object interface{}) (command Command, err error) {
	var (
		metaData = gmeta.Data(object)
	)
	if len(metaData) == 0 {
		err = gerror.Newf(
			`no meta data found in struct "%s"`,
			reflect.TypeOf(object).String(),
		)
		return
	}
	if err = gconv.Scan(metaData, &command); err != nil {
		return
	}
	// Name filed is necessary.
	if command.Name == "" {
		err = gerror.Newf(
			`command name cannot be empty, "name" tag not found in meta of struct "%s"`,
			reflect.TypeOf(object).String(),
		)
		return
	}
	if !command.NeedArgs {
		command.NeedArgs = gconv.Bool(metaData[tagNameArgs])
	}
	if command.Description == "" {
		command.Description = metaData[tagNameDc]
	}
	if command.Examples == "" {
		command.Examples = metaData[tagNameEg]
	}
	if command.Additional == "" {
		command.Additional = metaData[tagNameAd]
	}
	return
}

func newCommandFromMethod(object interface{}, method reflect.Value) (command Command, err error) {
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
	if command, err = newCommandFromObjectMeta(inputObject.Interface()); err != nil {
		return
	}

	// Options creating.
	if command.Options, err = newOptionsFromInput(inputObject.Interface()); err != nil {
		return
	}

	// Create function that has value return.
	command.FuncWithValue = func(ctx context.Context, parser *Parser) (out interface{}, err error) {
		ctx = context.WithValue(ctx, CtxKeyParser, parser)

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
		// Handle orphan options.
		for _, option := range command.Options {
			if option.Orphan && parser.ContainsOpt(option.Name) {
				data[option.Name] = "true"
			}
		}
		// Default values from struct tag.
		if err = mergeDefaultStructValue(data, inputObject.Interface()); err != nil {
			return nil, err
		}
		// Construct input parameters.
		if len(data) > 0 {
			if inputObject.Kind() == reflect.Ptr {
				err = gconv.Scan(data, inputObject.Interface())
			} else {
				err = gconv.Struct(data, inputObject.Addr().Interface())
			}
			if err != nil {
				return
			}
		}

		// Parameters validation.
		if err = gvalid.New().Bail().Data(inputObject.Interface()).Assoc(data).Run(ctx); err != nil {
			err = gerror.Wrap(gerror.Current(err), `validation failed for command options`)
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

func newOptionsFromInput(object interface{}) (options []Option, err error) {
	var (
		fields []gstructs.Field
	)
	fields, err = gstructs.Fields(gstructs.FieldsInput{
		Pointer:         object,
		RecursiveOption: gstructs.RecursiveOptionEmbeddedNoTag,
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
		if option.Name == helpOptionName {
			return nil, gerror.Newf(
				`option name "%s" is already token by built-in options`,
				option.Name,
			)
		}
		if option.Short == helpOptionNameShort {
			return nil, gerror.Newf(
				`short option name "%s" is already token by built-in options`,
				option.Short,
			)
		}
		options = append(options, option)
	}
	return
}

// mergeDefaultStructValue merges the request parameters with default values from struct tag definition.
func mergeDefaultStructValue(data map[string]interface{}, pointer interface{}) error {
	tagFields, err := gstructs.TagFields(pointer, defaultValueTags)
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
