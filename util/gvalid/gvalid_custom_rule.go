// Copyright 2020 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gvalid

// RuleFunc is the custom function for data validation.
// The parameter <rule> specifies the validation rule string, like "required", "between:1,100", etc.
// The parameter <value> specifies the value for this rule to validate.
// The parameter <message> specifies the custom error message or configured i18n message for this rule.
// The parameter <params> specifies all the parameters that needs. You can ignore parameter <params> if
// you do not really need it in your custom validation rule.
type RuleFunc func(rule string, value interface{}, message string, params map[string]interface{}) error

var (
	// customRuleFuncMap stores the custom rule functions.
	// map[Rule]RuleFunc
	customRuleFuncMap = make(map[string]RuleFunc)
)

// RegisterRule registers custom validation rule and function for package.
// It returns error if there's already the same rule registered previously.
func RegisterRule(rule string, f RuleFunc) error {
	customRuleFuncMap[rule] = f
	return nil
}

// DeleteRule deletes custom defined validation rule and its function from global package.
func DeleteRule(rule string) {
	delete(customRuleFuncMap, rule)
}
