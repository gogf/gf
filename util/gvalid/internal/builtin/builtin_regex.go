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

// RuleRegex implements `regex` rule:
// Value should match custom regular expression pattern.
//
// Format: regex:pattern
type RuleRegex struct{}

func init() {
	Register(RuleRegex{})
}

func (r RuleRegex) Name() string {
	return "regex"
}

func (r RuleRegex) Message() string {
	return "The {field} value `{value}` must be in regex of: {pattern}"
}

func (r RuleRegex) Run(in RunInput) error {
	if !gregex.IsMatchString(in.RulePattern, in.Value.String()) {
		return errors.New(in.Message)
	}
	return nil
}
