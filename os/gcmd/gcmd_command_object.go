// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//

package gcmd

import (
	"context"
	"fmt"
	"reflect"

	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/internal/reflection"
	"github.com/gogf/gf/v2/internal/utils"
	"github.com/gogf/gf/v2/os/gstructs"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gmeta"
	"github.com/gogf/gf/v2/util/gtag"
	"github.com/gogf/gf/v2/util/gutil"
	"github.com/gogf/gf/v2/util/gvalid"
)

var (
	// defaultValueTags is the struct tag names for default value storing.
	defaultValueTags = []string{"d", "default"}
)

// NewFromObject creates and returns a root command object using given object.
func NewFromObject(object interface{}) (rootCmd *Command, err error) {
	switch c := object.(type) {
	case Command:
		return &c, nil
	case *Command:
		return c, nil
	}

	originValueAndKind := reflection.OriginValueAndKind(object)
	if originValueAndKind.OriginKind != reflect.Struct {
		err = gerror.Newf(
			`input object should be type of struct, but got "%s"`,
			originValueAndKind.InputValue.Type().String(),
		)
		return
	}
	var reflectValue = originValueAndKind.InputValue
	// If given `object` is not pointer, it then creates a temporary one,
	// of which the value is `reflectValue`.
	// It then can retrieve all the methods both of struct/*struct.
	if reflectValue.Kind() == reflect.Struct {
		newValue := reflect.New(reflectValue.Type())
		newValue.Elem().Set(reflectValue)
		reflectValue = newValue
	}

	// Root command creating.
	rootCmd, err = newCommandFromObjectMeta(object, "")
	if err != nil {
		return
	}
	// Sub command creating.
	var (
		nameSet         = gset.NewStrSet()
		rootCommandName = gmeta.Get(object, gtag.Root).String()
		subCommands     []*Command
	)
	if rootCommandName == "" {
		rootCommandName = rootCmd.Name
	}
	for i := 0; i < reflectValue.NumMethod(); i++ {
		var (
			method      = reflectValue.Type().Method(i)
			methodValue = reflectValue.Method(i)
			methodType  = methodValue.Type()
			methodCmd   *Command
		)
		methodCmd, err = newCommandFromMethod(object, method, methodValue, methodType)
		if err != nil {
			return
		}
		if nameSet.Contains(methodCmd.Name) {
			err = gerror.Newf(
				`command name should be unique, found duplicated command name in method "%s"`,
				methodType.String(),
			)
			return
		}
		if rootCommandName == methodCmd.Name {
			methodToRootCmdWhenNameEqual(rootCmd, methodCmd)
		} else {
			subCommands = append(subCommands, methodCmd)
		}
	}
	if len(subCommands) > 0 {
		err = rootCmd.AddCommand(subCommands...)
	}
	return
}

func methodToRootCmdWhenNameEqual(rootCmd *Command, methodCmd *Command) {
	if rootCmd.Usage == "" {
		rootCmd.Usage = methodCmd.Usage
	}
	if rootCmd.Brief == "" {
		rootCmd.Brief = methodCmd.Brief
	}
	if rootCmd.Description == "" {
		rootCmd.Description = methodCmd.Description
	}
	if rootCmd.Examples == "" {
		rootCmd.Examples = methodCmd.Examples
	}
	if rootCmd.Func == nil {
		rootCmd.Func = methodCmd.Func
	}
	if rootCmd.FuncWithValue == nil {
		rootCmd.FuncWithValue = methodCmd.FuncWithValue
	}
	if rootCmd.HelpFunc == nil {
		rootCmd.HelpFunc = methodCmd.HelpFunc
	}
	if len(rootCmd.Arguments) == 0 {
		rootCmd.Arguments = methodCmd.Arguments
	}
	if !rootCmd.Strict {
		rootCmd.Strict = methodCmd.Strict
	}
	if rootCmd.Config == "" {
		rootCmd.Config = methodCmd.Config
	}
}

// The `object` is the Meta attribute from business object, and the `name` is the command name,
// commonly from method name, which is used when no name tag is defined in Meta.
func newCommandFromObjectMeta(object interface{}, name string) (command *Command, err error) {
	var metaData = gmeta.Data(object)
	if err = gconv.Scan(metaData, &command); err != nil {
		return
	}
	// Name field is necessary.
	if command.Name == "" {
		if name == "" {
			err = gerror.Newf(
				`command name cannot be empty, "name" tag not found in meta of struct "%s"`,
				reflect.TypeOf(object).String(),
			)
			return
		}
		command.Name = name
	}
	if command.Brief == "" {
		for _, tag := range []string{gtag.Summary, gtag.SummaryShort, gtag.SummaryShort2} {
			command.Brief = metaData[tag]
			if command.Brief != "" {
				break
			}
		}
	}
	if command.Description == "" {
		command.Description = metaData[gtag.DescriptionShort]
	}
	if command.Brief == "" && command.Description != "" {
		command.Brief = command.Description
		command.Description = ""
	}
	if command.Examples == "" {
		command.Examples = metaData[gtag.ExampleShort]
	}
	if command.Additional == "" {
		command.Additional = metaData[gtag.AdditionalShort]
	}
	return
}

func newCommandFromMethod(
	object interface{}, method reflect.Method, methodValue reflect.Value, methodType reflect.Type,
) (command *Command, err error) {
	// Necessary validation for input/output parameters and naming.
	if methodType.NumIn() != 2 || methodType.NumOut() != 2 {
		if methodType.PkgPath() != "" {
			err = gerror.NewCodef(
				gcode.CodeInvalidParameter,
				`invalid command: %s.%s.%s defined as "%s", but "func(context.Context, Input)(Output, error)" is required`,
				methodType.PkgPath(), reflect.TypeOf(object).Name(), method.Name, methodType.String(),
			)
		} else {
			err = gerror.NewCodef(
				gcode.CodeInvalidParameter,
				`invalid command: %s.%s defined as "%s", but "func(context.Context, Input)(Output, error)" is required`,
				reflect.TypeOf(object).Name(), method.Name, methodType.String(),
			)
		}
		return
	}
	if !methodType.In(0).Implements(reflect.TypeOf((*context.Context)(nil)).Elem()) {
		err = gerror.NewCodef(
			gcode.CodeInvalidParameter,
			`invalid command: %s.%s defined as "%s", but the first input parameter should be type of "context.Context"`,
			reflect.TypeOf(object).Name(), method.Name, methodType.String(),
		)
		return
	}
	if !methodType.Out(1).Implements(reflect.TypeOf((*error)(nil)).Elem()) {
		err = gerror.NewCodef(
			gcode.CodeInvalidParameter,
			`invalid command: %s.%s defined as "%s", but the last output parameter should be type of "error"`,
			reflect.TypeOf(object).Name(), method.Name, methodType.String(),
		)
		return
	}
	// The input struct should be named as `xxxInput`.
	if !gstr.HasSuffix(methodType.In(1).String(), `Input`) {
		err = gerror.NewCodef(
			gcode.CodeInvalidParameter,
			`invalid struct naming for input: defined as "%s", but it should be named with "Input" suffix like "xxxInput"`,
			methodType.In(1).String(),
		)
		return
	}
	// The output struct should be named as `xxxOutput`.
	if !gstr.HasSuffix(methodType.Out(0).String(), `Output`) {
		err = gerror.NewCodef(
			gcode.CodeInvalidParameter,
			`invalid struct naming for output: defined as "%s", but it should be named with "Output" suffix like "xxxOutput"`,
			methodType.Out(0).String(),
		)
		return
	}

	var inputObject reflect.Value
	if methodType.In(1).Kind() == reflect.Ptr {
		inputObject = reflect.New(methodType.In(1).Elem()).Elem()
	} else {
		inputObject = reflect.New(methodType.In(1)).Elem()
	}

	// Command creating.
	if command, err = newCommandFromObjectMeta(inputObject.Interface(), method.Name); err != nil {
		return
	}

	// Options creating.
	if command.Arguments, err = newArgumentsFromInput(inputObject.Interface()); err != nil {
		return
	}

	// For input struct converting using priority tag.
	var priorityTag = gstr.Join([]string{tagNameName, tagNameShort}, ",")

	// =============================================================================================
	// Create function that has value return.
	// =============================================================================================
	command.FuncWithValue = func(ctx context.Context, parser *Parser) (out interface{}, err error) {
		ctx = context.WithValue(ctx, CtxKeyParser, parser)
		var (
			data        = gconv.Map(parser.GetOptAll())
			argIndex    = 0
			arguments   = parser.GetArgAll()
			inputValues = []reflect.Value{reflect.ValueOf(ctx)}
		)
		if value := ctx.Value(CtxKeyArgumentsIndex); value != nil {
			argIndex = value.(int)
			// Use the left args to assign to input struct object.
			if argIndex < len(arguments) {
				arguments = arguments[argIndex:]
			}
		}
		if data == nil {
			data = map[string]interface{}{}
		}
		// Handle orphan options.
		for _, arg := range command.Arguments {
			if arg.IsArg {
				// Read argument from command line index.
				if argIndex < len(arguments) {
					data[arg.Name] = arguments[argIndex]
					argIndex++
				}
			} else {
				// Read argument from command line option name.
				if arg.Orphan {
					if orphanValue := parser.GetOpt(arg.Name); orphanValue != nil {
						if orphanValue.String() == "" {
							// Example: gf -f
							data[arg.Name] = "true"
							if arg.Short != "" {
								data[arg.Short] = "true"
							}
						} else {
							// Adapter with common user habits.
							// Eg:
							// `gf -f=0`: which parameter `f` is parsed as false
							// `gf -f=1`: which parameter `f` is parsed as true
							data[arg.Name] = orphanValue.Bool()
						}
					}
				}
			}
		}
		// Default values from struct tag.
		if err = mergeDefaultStructValue(data, inputObject.Interface()); err != nil {
			return nil, err
		}
		// Construct input parameters.
		if len(data) > 0 {
			intlog.PrintFunc(ctx, func() string {
				return fmt.Sprintf(`input command data map: %s`, gjson.MustEncode(data))
			})
			if inputObject.Kind() == reflect.Ptr {
				err = gconv.StructTag(data, inputObject.Interface(), priorityTag)
			} else {
				err = gconv.StructTag(data, inputObject.Addr().Interface(), priorityTag)
			}
			intlog.PrintFunc(ctx, func() string {
				return fmt.Sprintf(`input object assigned data: %s`, gjson.MustEncode(inputObject.Interface()))
			})
			if err != nil {
				return
			}
		}

		// Parameters validation.
		if err = gvalid.New().Bail().Data(inputObject.Interface()).Assoc(data).Run(ctx); err != nil {
			err = gerror.Wrapf(gerror.Current(err), `arguments validation failed for command "%s"`, command.Name)
			return
		}
		inputValues = append(inputValues, inputObject)

		// Call handler with dynamic created parameter values.
		results := methodValue.Call(inputValues)
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

func newArgumentsFromInput(object interface{}) (args []Argument, err error) {
	var (
		fields   []gstructs.Field
		nameSet  = gset.NewStrSet()
		shortSet = gset.NewStrSet()
	)
	fields, err = gstructs.Fields(gstructs.FieldsInput{
		Pointer:         object,
		RecursiveOption: gstructs.RecursiveOptionEmbeddedNoTag,
	})
	for _, field := range fields {
		var (
			arg      = Argument{}
			metaData = field.TagMap()
		)
		if err = gconv.Scan(metaData, &arg); err != nil {
			return nil, err
		}
		if arg.Name == "" {
			arg.Name = field.Name()
		}
		if arg.Name == helpOptionName {
			return nil, gerror.Newf(
				`argument name "%s" defined in "%s.%s" is already token by built-in arguments`,
				arg.Name, reflect.TypeOf(object).String(), field.Name(),
			)
		}
		if arg.Short == helpOptionNameShort {
			return nil, gerror.Newf(
				`short argument name "%s" defined in "%s.%s" is already token by built-in arguments`,
				arg.Short, reflect.TypeOf(object).String(), field.Name(),
			)
		}
		if arg.Brief == "" {
			arg.Brief = field.TagDescription()
		}
		if v, ok := metaData[gtag.Arg]; ok {
			arg.IsArg = gconv.Bool(v)
		}
		if nameSet.Contains(arg.Name) {
			return nil, gerror.Newf(
				`argument name "%s" defined in "%s.%s" is already token by other argument`,
				arg.Name, reflect.TypeOf(object).String(), field.Name(),
			)
		}
		nameSet.Add(arg.Name)

		if arg.Short != "" {
			if shortSet.Contains(arg.Short) {
				return nil, gerror.Newf(
					`short argument name "%s" defined in "%s.%s" is already token by other argument`,
					arg.Short, reflect.TypeOf(object).String(), field.Name(),
				)
			}
			shortSet.Add(arg.Short)
		}

		args = append(args, arg)
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
			var (
				nameValue  = field.Tag(tagNameName)
				shortValue = field.Tag(tagNameShort)
			)
			// If it already has value, it then ignores the default value.
			if value, ok := data[nameValue]; ok {
				data[field.Name()] = value
				continue
			}
			if value, ok := data[shortValue]; ok {
				data[field.Name()] = value
				continue
			}
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
