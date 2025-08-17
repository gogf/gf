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

// RuleUrl implements `url` rule:
// URL.
//
// Format: url
type RuleUrl struct{}

func init() {
	Register(RuleUrl{})
}

func (r RuleUrl) Name() string {
	return "url"
}

func (r RuleUrl) Message() string {
	return "The {field} value `{value}` is not a valid URL address"
}

func (r RuleUrl) Run(in RunInput) error {
	ok := gregex.IsMatchString(
		`(https?|ftp|file)://[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|]`,
		in.Value.String(),
	)
	if ok {
		return nil
	}
	return errors.New(in.Message)
}
