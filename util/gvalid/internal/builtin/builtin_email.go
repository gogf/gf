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

type RuleEmail struct{}

func init() {
	Register(&RuleEmail{})
}

func (r *RuleEmail) Name() string {
	return "email"
}

func (r *RuleEmail) Message() string {
	return "The {attribute} value `{value}` is not a valid email address"
}

func (r *RuleEmail) Run(in RunInput) error {
	ok := gregex.IsMatchString(
		`^[a-zA-Z0-9_\-\.]+@[a-zA-Z0-9_\-]+(\.[a-zA-Z0-9_\-]+)+$`,
		in.Value.String(),
	)
	if ok {
		return nil
	}
	return errors.New(in.Message)
}
