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

// RulePassword implements `password` rule:
// Universal password format rule1:
// Containing any visible chars, length between 6 and 18.
//
// Format: password
type RulePassword struct{}

func init() {
	Register(RulePassword{})
}

func (r RulePassword) Name() string {
	return "password"
}

func (r RulePassword) Message() string {
	return "The {field} value `{value}` is not a valid passport format"
}

func (r RulePassword) Run(in RunInput) error {
	if !gregex.IsMatchString(`^[\w\S]{6,18}$`, in.Value.String()) {
		return errors.New(in.Message)
	}
	return nil
}
