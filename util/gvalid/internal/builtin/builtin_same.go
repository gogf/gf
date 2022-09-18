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

// RuleSame implements `same` rule:
// Value should be the same as value of field.
//
// Format: same:field
type RuleSame struct{}

func init() {
	Register(&RuleSame{})
}

func (r *RuleSame) Name() string {
	return "same"
}

func (r *RuleSame) Message() string {
	return "The {attribute} value `{value}` must be the same as field {pattern}"
}

func (r *RuleSame) Run(in RunInput) error {
	var (
		ok    bool
		value = in.Value.String()
	)
	_, foundValue := gutil.MapPossibleItemByKey(in.Data.Map(), in.RulePattern)
	if foundValue != nil {
		if in.CaseInsensitive {
			ok = strings.EqualFold(value, gconv.String(foundValue))
		} else {
			ok = strings.Compare(value, gconv.String(foundValue)) == 0
		}
	}
	if !ok {
		return errors.New(in.Message)
	}
	return nil
}
