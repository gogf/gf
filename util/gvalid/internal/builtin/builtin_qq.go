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

type RuleQQ struct{}

func init() {
	Register(&RuleQQ{})
}

func (r *RuleQQ) Name() string {
	return "qq"
}

func (r *RuleQQ) Message() string {
	return "The {attribute} value `{value}` is not a valid QQ number"
}

func (r *RuleQQ) Run(in RunInput) error {
	ok := gregex.IsMatchString(
		`^[1-9][0-9]{4,}$`,
		in.Value.String(),
	)
	if ok {
		return nil
	}
	return errors.New(in.Message)
}
