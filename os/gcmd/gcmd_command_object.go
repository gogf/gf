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
	"github.com/gogf/gf/v2/internal/utils"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gmeta"
)

const (
	tagNameName        = `name`
	tagNameUsage       = `usage`
	tagNameBrief       = `brief`
	tagNameShort       = `short`
	tagNameOrphan      = `orphan`
	tagNameDescription = `description`
	tagNameDc          = `dc`
	tagNameAddition    = `additional`
	tagNameAd          = `ad`
)

func (c *Command) AddObject(objects ...interface{}) (err error) {
	for _, object := range objects {
		if err = c.doAddObject(object); err != nil {
			return err
		}
	}
	return nil
}

func (c *Command) doAddObject(object interface{}) error {
	originValueAndKind := utils.OriginValueAndKind(object)
	if originValueAndKind.OriginKind != reflect.Struct {
		return gerror.Newf(
			`input object should be type of struct, but got "%s"`,
			originValueAndKind.InputValue.Type().String(),
		)
	}
	for i := 0; i < originValueAndKind.InputValue.NumMethod(); i++ {
		method := originValueAndKind.InputValue.Method(i)
	}
	for _, field := range fields {

	}
	return nil
}

func newCommandFromMethod(object, method reflect.Value) (cmd *Command, err error) {
	var (
		reflectType = method.Type()
	)
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
		return
	}

	if reflectType.In(0).String() != "context.Context" {
		err = gerror.NewCodef(
			gcode.CodeInvalidParameter,
			`invalid handler: defined as "%s", but the first input parameter should be type of "context.Context"`,
			reflectType.String(),
		)
		return
	}

	if reflectType.Out(1).String() != "error" {
		err = gerror.NewCodef(
			gcode.CodeInvalidParameter,
			`invalid handler: defined as "%s", but the last output parameter should be type of "error"`,
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
		metaMap = gmeta.Data()
	)
}
