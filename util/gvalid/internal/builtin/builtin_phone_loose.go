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

// RulePhoneLoose implements `phone-loose` rule:
// Loose mobile phone number verification(宽松的手机号验证)
// As long as the 11 digits numbers beginning with
// 13, 14, 15, 16, 17, 18, 19 can pass the verification
// (只要满足 13、14、15、16、17、18、19开头的11位数字都可以通过验证).
//
// Format: phone-loose
type RulePhoneLoose struct{}

func init() {
	Register(&RulePhoneLoose{})
}

func (r *RulePhoneLoose) Name() string {
	return "phone-loose"
}

func (r *RulePhoneLoose) Message() string {
	return "The {attribute} value `{value}` is not a valid phone number"
}

func (r *RulePhoneLoose) Run(in RunInput) error {
	ok := gregex.IsMatchString(
		`^1(3|4|5|6|7|8|9)\d{9}$`,
		in.Value.String(),
	)
	if ok {
		return nil
	}
	return errors.New(in.Message)
}
