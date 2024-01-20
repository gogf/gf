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

// RulePassword3 implements `password3` rule:
// Universal password format rule3:
// Must meet password rule1, must contain lower and upper letters, numbers and special chars.
//
// Format: password3
type RulePassword3 struct{}

func init() {
	Register(RulePassword3{})
}

func (r RulePassword3) Name() string {
	return "password3"
}

func (r RulePassword3) Message() string {
	return "The {field} value `{value}` is not a valid password3 format"
}

func (r RulePassword3) Run(in RunInput) error {
	var value = in.Value.String()
	if gregex.IsMatchString(`^[\w\S]{6,18}$`, value) &&
		gregex.IsMatchString(`[a-z]+`, value) &&
		gregex.IsMatchString(`[A-Z]+`, value) &&
		gregex.IsMatchString(`\d+`, value) &&
		gregex.IsMatchString(`[^a-zA-Z0-9]+`, value) {
		return nil
	}
	return errors.New(in.Message)
}
