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

// RuleTelephone implements `telephone` rule:
// "XXXX-XXXXXXX"
// "XXXX-XXXXXXXX"
// "XXX-XXXXXXX"
// "XXX-XXXXXXXX"
// "XXXXXXX"
// "XXXXXXXX"
//
// Format:  telephone
type RuleTelephone struct{}

func init() {
	Register(&RuleTelephone{})
}

func (r *RuleTelephone) Name() string {
	return "telephone"
}

func (r *RuleTelephone) Message() string {
	return "The {attribute} value `{value}` is not a valid telephone number"
}

func (r *RuleTelephone) Run(in RunInput) error {
	ok := gregex.IsMatchString(
		`^((\d{3,4})|\d{3,4}-)?\d{7,8}$`,
		in.Value.String(),
	)
	if ok {
		return nil
	}
	return errors.New(in.Message)
}
