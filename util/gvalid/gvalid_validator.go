// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gvalid

import (
	"github.com/gogf/gf/i18n/gi18n"
)

// Validator is the validation manager for chaining operations.
type Validator struct {
	i18nManager                      *gi18n.Manager      // I18n manager for error message translation.
	key                              string              // Single validation key.
	value                            interface{}         // Single validation value.
	data                             interface{}         // Validation data, which is usually a map.
	rules                            interface{}         // Custom validation data.
	messages                         interface{}         // Custom validation error messages, which can be string or type of CustomMsg.
	ruleFuncMap                      map[string]RuleFunc // ruleFuncMap stores custom rule functions for current Validator.
	useDataInsteadOfObjectAttributes bool                // Using `data` as its validation source instead of attribute values from `Object`.
	bail                             bool                // Stop validation after the first validation error.
}

// New creates and returns a new Validator.
func New() *Validator {
	return &Validator{
		i18nManager: gi18n.Instance(),          // Use default i18n manager.
		ruleFuncMap: make(map[string]RuleFunc), // Custom rule function storing map.
	}
}

// Clone creates and returns a new Validator which is a shallow copy of current one.
func (v *Validator) Clone() *Validator {
	newValidator := New()
	*newValidator = *v
	return newValidator
}

// I18n sets the i18n manager for the validator.
func (v *Validator) I18n(i18nManager *gi18n.Manager) *Validator {
	newValidator := v.Clone()
	newValidator.i18nManager = i18nManager
	return newValidator
}

// Bail sets the mark for stopping validation after the first validation error.
func (v *Validator) Bail() *Validator {
	newValidator := v.Clone()
	newValidator.bail = true
	return newValidator
}

// Data is a chaining operation function, which sets validation data for current operation.
// The parameter `data` is usually type of map, which specifies the parameter map used in validation.
// Calling this function also sets `useDataInsteadOfObjectAttributes` true no mather the `data` is nil or not.
func (v *Validator) Data(data interface{}) *Validator {
	newValidator := v.Clone()
	newValidator.data = data
	newValidator.useDataInsteadOfObjectAttributes = true
	return newValidator
}

// Rules is a chaining operation function, which sets custom validation rules for current operation.
func (v *Validator) Rules(rules interface{}) *Validator {
	newValidator := v.Clone()
	newValidator.rules = rules
	return newValidator
}

// Messages is a chaining operation function, which sets custom error messages for current operation.
// The parameter `messages` can be type of string/[]string/map[string]string. It supports sequence in error result
// if `rules` is type of []string.
func (v *Validator) Messages(messages interface{}) *Validator {
	newValidator := v.Clone()
	newValidator.messages = messages
	return newValidator
}

// RuleFunc registers one custom rule function to current Validator.
func (v *Validator) RuleFunc(rule string, f RuleFunc) *Validator {
	newValidator := v.Clone()
	newValidator.ruleFuncMap[rule] = f
	return newValidator
}

// RuleFuncMap registers multiple custom rule functions to current Validator.
func (v *Validator) RuleFuncMap(m map[string]RuleFunc) *Validator {
	newValidator := v.Clone()
	for k, v := range m {
		newValidator.ruleFuncMap[k] = v
	}
	return newValidator
}

// getRuleFunc retrieves and returns the custom rule function for specified rule.
func (v *Validator) getRuleFunc(rule string) RuleFunc {
	ruleFunc := v.ruleFuncMap[rule]
	if ruleFunc == nil {
		ruleFunc = customRuleFuncMap[rule]
	}
	return ruleFunc
}
