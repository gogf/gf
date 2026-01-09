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

// RuleNumeric implements `numeric` rule:
// Numeric string (0-9).
//
// Format: numeric
type RuleNumeric struct{}

func init() {
	Register(RuleNumeric{})
}

func (r RuleNumeric) Name() string {
	return "numeric"
}

func (r RuleNumeric) Message() string {
	return "The {field} value `{value}` must be numeric"
}

func (r RuleNumeric) Run(in RunInput) error {
	ok := gregex.IsMatchString(`^[0-9]+$`, in.Value.String())
	if ok {
		return nil
	}
	return errors.New(in.Message)
}
