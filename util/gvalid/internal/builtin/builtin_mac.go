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

type RuleMac struct{}

func init() {
	Register(&RuleMac{})
}

func (r *RuleMac) Name() string {
	return "mac"
}

func (r *RuleMac) Message() string {
	return "The {attribute} value `{value}` is not a valid MAC address"
}

func (r *RuleMac) Run(in RunInput) error {
	ok := gregex.IsMatchString(
		`^([0-9A-Fa-f]{2}[\-:]){5}[0-9A-Fa-f]{2}$`,
		in.Value.String(),
	)
	if ok {
		return nil
	}
	return errors.New(in.Message)
}
