// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package builtin

import (
	"errors"
)

type RuleRequired struct{}

func init() {
	Register(&RuleRequired{})
}

func (r *RuleRequired) Name() string {
	return "required"
}

func (r *RuleRequired) Message() string {
	return "The {attribute} field is required"
}

func (r *RuleRequired) Run(in RunInput) error {
	if in.Value.IsEmpty() {
		return errors.New(in.Message)
	}
	return nil
}
