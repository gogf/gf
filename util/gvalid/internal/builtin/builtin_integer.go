// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package builtin

import (
	"errors"
	"strconv"
)

// RuleInteger implements `integer` rule:
// Integer.
//
// Format: integer
type RuleInteger struct{}

func init() {
	Register(&RuleInteger{})
}

func (r *RuleInteger) Name() string {
	return "integer"
}

func (r *RuleInteger) Message() string {
	return "The {attribute} value `{value}` is not an integer"
}

func (r *RuleInteger) Run(in RunInput) error {
	if _, err := strconv.Atoi(in.Value.String()); err == nil {
		return nil
	}
	return errors.New(in.Message)
}
