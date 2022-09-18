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

// RulePassword2 implements `password2` rule:
// Universal password format rule2:
// Must meet password rule1, must contain lower and upper letters and numbers.
//
// Format: password2
type RulePassword2 struct{}

func init() {
	Register(&RulePassword2{})
}

func (r *RulePassword2) Name() string {
	return "password2"
}

func (r *RulePassword2) Message() string {
	return "The {attribute} value `{value}` is not a valid passport format"
}

func (r *RulePassword2) Run(in RunInput) error {
	var value = in.Value.String()
	if gregex.IsMatchString(`^[\w\S]{6,18}$`, value) &&
		gregex.IsMatchString(`[a-z]+`, value) &&
		gregex.IsMatchString(`[A-Z]+`, value) &&
		gregex.IsMatchString(`\d+`, value) {
		return nil
	}
	return errors.New(in.Message)
}
