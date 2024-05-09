// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package builtin

import (
	"errors"
	"strings"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gutil"
)

// RuleRequiredIf implements `required-if` rule:
// Required if any of given field and its value are equal.
//
// Format:  required-if:field,value,...
// Example: required-if:id,1,age,18
type RuleRequiredIf struct{}

func init() {
	Register(RuleRequiredIf{})
}

func (r RuleRequiredIf) Name() string {
	return "required-if"
}

func (r RuleRequiredIf) Message() string {
	return "The {field} field is required"
}

func (r RuleRequiredIf) Run(in RunInput) error {
	var (
		required   = false
		array      = strings.Split(in.RulePattern, ",")
		foundValue interface{}
		dataMap    = in.Data.Map()
	)
	if len(array)%2 != 0 {
		return gerror.NewCodef(
			gcode.CodeInvalidParameter,
			`invalid "%s" rule pattern: %s`,
			r.Name(),
			in.RulePattern,
		)
	}
	// It supports multiple field and value pairs.
	for i := 0; i < len(array); {
		var (
			tk = array[i]
			tv = array[i+1]
		)
		_, foundValue = gutil.MapPossibleItemByKey(dataMap, tk)
		if in.Option.CaseInsensitive {
			required = strings.EqualFold(tv, gconv.String(foundValue))
		} else {
			required = strings.Compare(tv, gconv.String(foundValue)) == 0
		}
		if required {
			break
		}
		i += 2
	}
	if required && isRequiredEmpty(in.Value.Val()) {
		return errors.New(in.Message)
	}
	return nil
}
