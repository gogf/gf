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

// RuleLowercase implements `lowercase` rule:
// Lowercase alphabetic characters (a-z).
//
// Format: lowercase
type RuleLowercase struct{}

func init() {
	Register(RuleLowercase{})
}

func (r RuleLowercase) Name() string {
	return "lowercase"
}

func (r RuleLowercase) Message() string {
	return "The {field} value `{value}` must be lowercase"
}

func (r RuleLowercase) Run(in RunInput) error {
	ok := gregex.IsMatchString(`^[a-z]+$`, in.Value.String())
	if ok {
		return nil
	}
	return errors.New(in.Message)
}
