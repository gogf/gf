// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package builtin implements built-in validation rules.
package builtin

import (
	"github.com/gogf/gf/v2/container/gvar"
)

type Rule interface {
	// Name returns the builtin name of the rule.
	Name() string

	// Message returns the default error message of the rule.
	Message() string

	// Run starts running the rule, it returns nil if successful, or else an error.
	Run(in RunInput) error
}

type RunInput struct {
	RuleKey         string    // RuleKey is like the "max" in rule "max: 6"
	RulePattern     string    // RulePattern is like "6" in rule:"max:6"
	Message         string    // Message specifies the custom error message or configured i18n message for this rule.
	Value           *gvar.Var // Value specifies the value for this rule to validate.
	Data            *gvar.Var // Data specifies the `data` which is passed to the Validator.
	CaseInsensitive bool      // CaseInsensitive indicates that it does Case-Insensitive comparison in string.
}

var (
	ruleMap = map[string]Rule{}
)

func Register(rule Rule) {
	ruleMap[rule.Name()] = rule
}

func GetRule(name string) Rule {
	return ruleMap[name]
}
