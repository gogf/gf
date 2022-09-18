// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package builtin

import (
	"errors"

	"github.com/gogf/gf/v2/text/gregex"
)

// RuleNotRegex implements `not-regex` rule:
// Value should not match custom regular expression pattern.
//
// Format: not-regex:pattern
type RuleNotRegex struct{}

func init() {
	Register(&RuleNotRegex{})
}

func (r *RuleNotRegex) Name() string {
	return "not-regex"
}

func (r *RuleNotRegex) Message() string {
	return "The {attribute} value `{value}` should not be in regex of: {pattern}"
}

func (r *RuleNotRegex) Run(in RunInput) error {
	if gregex.IsMatchString(in.RulePattern, in.Value.String()) {
		return errors.New(in.Message)
	}
	return nil
}
