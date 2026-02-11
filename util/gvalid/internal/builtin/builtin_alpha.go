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

// RuleAlpha implements `alpha` rule:
// Alpha characters (a-z, A-Z).
//
// Format: alpha
type RuleAlpha struct{}

func init() {
	Register(RuleAlpha{})
}

func (r RuleAlpha) Name() string {
	return "alpha"
}

func (r RuleAlpha) Message() string {
	return "The {field} value `{value}` must contain only alphabetic characters"
}

func (r RuleAlpha) Run(in RunInput) error {
	ok := gregex.IsMatchString(`^[a-zA-Z]+$`, in.Value.String())
	if ok {
		return nil
	}
	return errors.New(in.Message)
}
