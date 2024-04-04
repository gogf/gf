// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package builtin

import (
	"errors"
	"strings"

	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gutil"
)

// RuleRequiredIfAnyEq implements `required-if-any-eq` rule:
// Required if any given field and its value are equal.
//
// Format:  required-if-any-eq:field,value,...
// Example: required-if-any-eq: id,1,age,18
type RuleRequiredIfAnyEq struct{}

func init() {
	Register(RuleRequiredIfAnyEq{})
}

func (r RuleRequiredIfAnyEq) Name() string {
	return "required-if-any-eq"
}

func (r RuleRequiredIfAnyEq) Message() string {
	return "The {field} field is required"
}

func (r RuleRequiredIfAnyEq) Run(in RunInput) error {
	var (
		required   = false
		array      = strings.Split(in.RulePattern, ",")
		foundValue interface{}
		dataMap    = in.Data.Map()
	)
	// It supports multiple field and value pairs.
	if len(array)%2 == 0 {
		for i := 0; i < len(array); {
			tk := array[i]
			tv := array[i+1]
			var eq bool
			_, foundValue = gutil.MapPossibleItemByKey(dataMap, tk)
			if in.Option.CaseInsensitive {
				eq = strings.EqualFold(tv, gconv.String(foundValue))
			} else {
				eq = strings.Compare(tv, gconv.String(foundValue)) == 0
			}
			if eq {
				required = true
				break
			}
			i += 2
		}
	}

	if required && isRequiredEmpty(in.Value.Val()) {
		return errors.New(in.Message)
	}
	return nil
}
