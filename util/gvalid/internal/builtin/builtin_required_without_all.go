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

// RuleRequiredWithoutAll implements `required-without-all` rule: Required if all given fields are empty.
// Format:  required-without-all:field1,field2,...
// Example: required-without-all:id,name
type RuleRequiredWithoutAll struct{}

func init() {
	Register(&RuleRequiredWithoutAll{})
}

func (r *RuleRequiredWithoutAll) Name() string {
	return "required-without-all"
}

func (r *RuleRequiredWithoutAll) Message() string {
	return "The {attribute} field is required"
}

func (r *RuleRequiredWithoutAll) Run(in RunInput) error {
	var (
		required   = true
		array      = strings.Split(in.RulePattern, ",")
		foundValue interface{}
	)
	for i := 0; i < len(array); i++ {
		_, foundValue = gutil.MapPossibleItemByKey(in.Data.Map(), array[i])
		if !empty.IsEmpty(foundValue) {
			required = false
			break
		}
	}

	if required && isRequiredEmpty(in.Value.Val()) {
		return errors.New(in.Message)
	}
	return nil
}
