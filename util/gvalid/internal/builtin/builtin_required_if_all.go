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

// RuleRequiredIfAll implements `required-if-all` rule:
// Required if all given field and its value are equal.
//
// Format:  required-if-all:field,value,...
// Example: required-if-all:id,1,age,18
type RuleRequiredIfAll struct{}

func init() {
	Register(RuleRequiredIfAll{})
}

func (r RuleRequiredIfAll) Name() string {
	return "required-if-all"
}

func (r RuleRequiredIfAll) Message() string {
	return "The {field} field is required"
}

func (r RuleRequiredIfAll) Run(in RunInput) error {
	var (
		required   = true
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
	for i := 0; i < len(array); {
		var (
			tk = array[i]
			tv = array[i+1]
			eq bool
		)
		_, foundValue = gutil.MapPossibleItemByKey(dataMap, tk)
		if in.Option.CaseInsensitive {
			eq = strings.EqualFold(tv, gconv.String(foundValue))
		} else {
			eq = strings.Compare(tv, gconv.String(foundValue)) == 0
		}
		if !eq {
			required = false
			break
		}
		i += 2
	}
	if required && isRequiredEmpty(in.Value.Val()) {
		return errors.New(in.Message)
	}
	return nil
}
