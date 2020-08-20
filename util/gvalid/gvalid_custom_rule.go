// Copyright 2020 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gvalid

import (
	"fmt"
)

// RuleFunc is the custom function for data validation.
// The parameter <value> specifies the value for this rule to validate.
// The parameter <message> specifies the custom error message or configured i18n message for this rule.
// The parameter <params> specifies all the parameters that needs .
type RuleFunc func(value interface{}, message string, params map[string]interface{}) error

var (
	// customRuleFuncMap stores the custom rule functions.
	customRuleFuncMap = make(map[string]RuleFunc)
)

// RegisterRule registers custom validation rule and function for package.
// It returns error if there's already the same rule registered previously.
func RegisterRule(rule string, f RuleFunc) error {
	if _, ok := allSupportedRules[rule]; ok {
		return fmt.Errorf(`validation rule "%s" is already registered`, rule)
	}
	allSupportedRules[rule] = struct{}{}
	customRuleFuncMap[rule] = f
	return nil
}
