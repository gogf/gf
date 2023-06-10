// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package builtin

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/internal/reflection"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gtag"
)

// RuleEnums implements `enums` rule:
// Value should be in enums of its constant type.
//
// Format: enums
type RuleEnums struct{}

func init() {
	Register(RuleEnums{})
}

func (r RuleEnums) Name() string {
	return "enums"
}

func (r RuleEnums) Message() string {
	return "The {field} value `{value}` should be in enums of: {enums}"
}

func (r RuleEnums) Run(in RunInput) error {
	originTypeAndKind := reflection.OriginTypeAndKind(in.Data.Val())
	switch originTypeAndKind.OriginKind {
	case reflect.Struct:
		for i := 0; i < originTypeAndKind.OriginType.NumField(); i++ {
			field := originTypeAndKind.OriginType.Field(i)
			if in.Field == field.Name {
				var (
					typeId   = fmt.Sprintf(`%s.%s`, field.Type.PkgPath(), field.Type.Name())
					tagEnums = gtag.GetEnumsByType(typeId)
				)
				if tagEnums == "" {
					return gerror.NewCodef(
						gcode.CodeInvalidOperation,
						`no enums found for type "%s"`,
						typeId,
					)
				}
				var enumsValues = make([]interface{}, 0)
				if err := json.Unmarshal([]byte(tagEnums), &enumsValues); err != nil {
					return err
				}
				if !gstr.InArray(gconv.Strings(enumsValues), in.Value.String()) {
					return errors.New(gstr.Replace(
						in.Message, `{enums}`, tagEnums,
					))
				}
			}
		}

	default:
		return gerror.NewCode(
			gcode.CodeInvalidOperation,
			`"enums" validation rule can only be used in struct validation currently`,
		)
	}
	return nil
}
