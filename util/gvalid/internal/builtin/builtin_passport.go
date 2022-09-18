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

// RulePassport implements `passport` rule:
// Universal passport format rule:
// Starting with letter, containing only numbers or underscores, length between 6 and 18
//
// Format:  passport
type RulePassport struct{}

func init() {
	Register(&RulePassport{})
}

func (r *RulePassport) Name() string {
	return "passport"
}

func (r *RulePassport) Message() string {
	return "The {attribute} value `{value}` is not a valid passport format"
}

func (r *RulePassport) Run(in RunInput) error {
	ok := gregex.IsMatchString(
		`^[a-zA-Z]{1}\w{5,17}$`,
		in.Value.String(),
	)
	if ok {
		return nil
	}
	return errors.New(in.Message)
}
