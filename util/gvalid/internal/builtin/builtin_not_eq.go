// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package builtin

// RuleNotEq implements `not-eq` rule:
// Value should be different from value of field.
//
// Format: not-eq:field
type RuleNotEq struct{}

func init() {
	Register(RuleNotEq{})
}

func (r RuleNotEq) Name() string {
	return "not-eq"
}

func (r RuleNotEq) Message() string {
	return "The {field} value `{value}` must not be equal to field {pattern}"
}

func (r RuleNotEq) Run(in RunInput) error {
	return RuleDifferent{}.Run(in)
}
