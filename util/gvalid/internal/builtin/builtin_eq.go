// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package builtin

// RuleEq implements `eq` rule:
// Value should be the same as value of field.
//
// This rule performs the same as rule `same`.
//
// Format: eq:field
type RuleEq struct{}

func init() {
	Register(RuleEq{})
}

func (r RuleEq) Name() string {
	return "eq"
}

func (r RuleEq) Message() string {
	return "The {field} value `{value}` must be equal to field {pattern}"
}

func (r RuleEq) Run(in RunInput) error {
	return RuleSame{}.Run(in)
}
