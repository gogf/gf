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

// Loose mobile phone number verification(宽松的手机号验证)
// As long as the 11 digits numbers beginning with
// 13, 14, 15, 16, 17, 18, 19 can pass the verification
// (只要满足 13、14、15、16、17、18、19开头的11位数字都可以通过验证).

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
		`^13[\d]{9}$|^14[5,7]{1}\d{8}$|^15[^4]{1}\d{8}$|^16[\d]{9}$|^17[0,2,3,5,6,7,8]{1}\d{8}$|^18[\d]{9}$|^19[\d]{9}$`,
		in.Value.String(),
	)
	if ok {
		return nil
	}
	return errors.New(in.Message)
}
