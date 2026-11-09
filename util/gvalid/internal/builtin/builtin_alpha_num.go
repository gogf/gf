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

// RuleAlphaNum implements `alpha-num` rule:
// Alpha-numeric characters (a-z, A-Z, 0-9).
//
// Format: alpha-num
type RuleAlphaNum struct{}

func init() {
	Register(RuleAlphaNum{})
}

func (r RuleAlphaNum) Name() string {
	return "alpha-num"
}

func (r RuleAlphaNum) Message() string {
	return "The {field} value `{value}` must contain only alpha-numeric characters"
}

func (r RuleAlphaNum) Run(in RunInput) error {
	ok := gregex.IsMatchString(`^[a-zA-Z0-9]+$`, in.Value.String())
	if ok {
		return nil
	}
	return errors.New(in.Message)
}
