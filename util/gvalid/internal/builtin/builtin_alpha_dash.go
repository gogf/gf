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

// RuleAlphaDash implements `alpha-dash` rule:
// Alpha-numeric characters, hyphens, and underscores (a-z, A-Z, 0-9, -, _).
//
// Format: alpha-dash
type RuleAlphaDash struct{}

func init() {
	Register(RuleAlphaDash{})
}

func (r RuleAlphaDash) Name() string {
	return "alpha-dash"
}

func (r RuleAlphaDash) Message() string {
	return "The {field} value `{value}` must contain only alpha-numeric characters, hyphens, and underscores"
}

func (r RuleAlphaDash) Run(in RunInput) error {
	ok := gregex.IsMatchString(`^[a-zA-Z0-9_\-]+$`, in.Value.String())
	if ok {
		return nil
	}
	return errors.New(in.Message)
}
