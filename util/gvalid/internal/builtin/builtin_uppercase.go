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

// RuleUppercase implements `uppercase` rule:
// Uppercase alphabetic characters (A-Z).
//
// Format: uppercase
type RuleUppercase struct{}

func init() {
	Register(RuleUppercase{})
}

func (r RuleUppercase) Name() string {
	return "uppercase"
}

func (r RuleUppercase) Message() string {
	return "The {field} value `{value}` must be uppercase"
}

func (r RuleUppercase) Run(in RunInput) error {
	ok := gregex.IsMatchString(`^[A-Z]+$`, in.Value.String())
	if ok {
		return nil
	}
	return errors.New(in.Message)
}
