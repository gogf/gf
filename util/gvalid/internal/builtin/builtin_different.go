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

// RuleDifferent implements `different` rule:
// Value should be different from value of field.
//
// Format: different:field
type RuleDifferent struct{}

func init() {
	Register(RuleDifferent{})
}

func (r RuleDifferent) Name() string {
	return "different"
}

func (r RuleDifferent) Message() string {
	return "The {field} value `{value}` must be different from field {pattern}"
}

func (r RuleDifferent) Run(in RunInput) error {
	var (
		ok    = true
		value = in.Value.String()
	)
	_, foundValue := gutil.MapPossibleItemByKey(in.Data.Map(), in.RulePattern)
	if foundValue != nil {
		if in.Option.CaseInsensitive {
			ok = !strings.EqualFold(value, gconv.String(foundValue))
		} else {
			ok = strings.Compare(value, gconv.String(foundValue)) != 0
		}
	}
	if !ok {
		return errors.New(in.Message)
	}
	return nil
}
