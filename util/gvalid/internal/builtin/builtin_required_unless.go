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

// RuleRequiredUnless implements `required-unless` rule: Required unless all given field and its value are not equal.
// Format:  required-unless:field,value,...
// Example: required-unless: id,1,age,18
type RuleRequiredUnless struct{}

func init() {
	Register(&RuleRequiredUnless{})
}

func (r *RuleRequiredUnless) Name() string {
	return "required-unless"
}

func (r *RuleRequiredUnless) Message() string {
	return "The {attribute} field is required"
}

func (r *RuleRequiredUnless) Run(in RunInput) error {
	var (
		required   = true
		array      = strings.Split(in.RulePattern, ",")
		foundValue interface{}
	)
	// It supports multiple field and value pairs.
	if len(array)%2 == 0 {
		for i := 0; i < len(array); {
			tk := array[i]
			tv := array[i+1]
			_, foundValue = gutil.MapPossibleItemByKey(in.Data.Map(), tk)
			if in.CaseInsensitive {
				required = !strings.EqualFold(tv, gconv.String(foundValue))
			} else {
				required = strings.Compare(tv, gconv.String(foundValue)) != 0
			}
			if !required {
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
