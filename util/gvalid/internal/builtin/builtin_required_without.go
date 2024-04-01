// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package builtin

import (
	"errors"
	"strings"

	"github.com/gogf/gf/v2/internal/empty"
	"github.com/gogf/gf/v2/util/gutil"
)

// RuleRequiredWithout implements `required-without` rule:
// Required if any of given fields are empty.
//
// Format:  required-without:field1,field2,...
// Example: required-without:id,name
type RuleRequiredWithout struct{}

func init() {
	Register(RuleRequiredWithout{})
}

func (r RuleRequiredWithout) Name() string {
	return "required-without"
}

func (r RuleRequiredWithout) Message() string {
	return "The {field} field is required"
}

func (r RuleRequiredWithout) Run(in RunInput) error {
	var (
		required   = false
		array      = strings.Split(in.RulePattern, ",")
		foundValue interface{}
		dataMap    = in.Data.Map()
	)

	for i := 0; i < len(array); i++ {
		_, foundValue = gutil.MapPossibleItemByKey(dataMap, array[i])
		if empty.IsEmpty(foundValue) {
			required = true
			break
		}
	}

	if required && isRequiredEmpty(in.Value.Val()) {
		return errors.New(in.Message)
	}
	return nil
}
