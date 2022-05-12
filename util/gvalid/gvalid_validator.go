// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gvalid

import (
	"context"
	"errors"
	"reflect"

	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/internal/reflection"
	"github.com/gogf/gf/v2/internal/utils"
	"github.com/gogf/gf/v2/util/gconv"
)

// Validator is the validation manager for chaining operations.
type Validator struct {
	i18nManager                       *gi18n.Manager      // I18n manager for error message translation.
	data                              interface{}         // Validation data, which can be a map, struct or a certain value to be validated.
	assoc                             interface{}         // Associated data, which is usually a map, for union validation.
	rules                             interface{}         // Custom validation data.
	messages                          interface{}         // Custom validation error messages, which can be string or type of CustomMsg.
	ruleFuncMap                       map[string]RuleFunc // ruleFuncMap stores custom rule functions for current Validator.
	useAssocInsteadOfObjectAttributes bool                // Using `assoc` as its validation source instead of attribute values from `Object`.
	bail                              bool                // Stop validation after the first validation error.
	caseInsensitive                   bool                // Case-Insensitive configuration for those rules that need value comparison.
}

// New creates and returns a new Validator.
func New() *Validator {
	return &Validator{
		i18nManager: gi18n.Instance(),          // Use default i18n manager.
		ruleFuncMap: make(map[string]RuleFunc), // Custom rule function storing map.
	}
}

// Run starts validating the given data with rules and messages.
func (v *Validator) Run(ctx context.Context) Error {
	if v.data == nil {
		return newValidationErrorByStr(
			internalParamsErrRuleName,
			errors.New(`no data passed for validation`),
		)
	}

	originValueAndKind := reflection.OriginValueAndKind(v.data)
	switch originValueAndKind.OriginKind {
	case reflect.Map:
		isMapValidation := false
		if v.rules == nil {
			isMapValidation = true
		} else if utils.IsMap(v.rules) || utils.IsSlice(v.rules) {
			isMapValidation = true
		}
		if isMapValidation {
			return v.doCheckMap(ctx, v.data)
		}

	case reflect.Struct:
		isStructValidation := false
		if v.rules == nil {
			isStructValidation = true
		} else if utils.IsMap(v.rules) || utils.IsSlice(v.rules) {
			isStructValidation = true
		}
		if isStructValidation {
			return v.doCheckStruct(ctx, v.data)
		}
	}

	return v.doCheckValue(ctx, doCheckValueInput{
		Name:     "",
		Value:    v.data,
		Rule:     gconv.String(v.rules),
		Messages: v.messages,
		DataRaw:  v.assoc,
		DataMap:  gconv.Map(v.assoc),
	})
}

// Clone creates and returns a new Validator which is a shallow copy of current one.
func (v *Validator) Clone() *Validator {
	newValidator := New()
	*newValidator = *v
	return newValidator
}

// I18n sets the i18n manager for the validator.
func (v *Validator) I18n(i18nManager *gi18n.Manager) *Validator {
	if i18nManager == nil {
		return v
	}
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

// Ci sets the mark for Case-Insensitive for those rules that need value comparison.
func (v *Validator) Ci() *Validator {
	newValidator := v.Clone()
	newValidator.caseInsensitive = true
	return newValidator
}

// Data is a chaining operation function, which sets validation data for current operation.
func (v *Validator) Data(data interface{}) *Validator {
	if data == nil {
		return v
	}
	newValidator := v.Clone()
	newValidator.data = data
	return newValidator
}

// Assoc is a chaining operation function, which sets associated validation data for current operation.
// The optional parameter `assoc` is usually type of map, which specifies the parameter map used in union validation.
// Calling this function with `assoc` also sets `useAssocInsteadOfObjectAttributes` true.
func (v *Validator) Assoc(assoc interface{}) *Validator {
	if assoc == nil {
		return v
	}
	newValidator := v.Clone()
	newValidator.assoc = assoc
	newValidator.useAssocInsteadOfObjectAttributes = true
	return newValidator
}

// Rules is a chaining operation function, which sets custom validation rules for current operation.
func (v *Validator) Rules(rules interface{}) *Validator {
	if rules == nil {
		return v
	}
	newValidator := v.Clone()
	newValidator.rules = rules
	return newValidator
}

// Messages is a chaining operation function, which sets custom error messages for current operation.
// The parameter `messages` can be type of string/[]string/map[string]string. It supports sequence in error result
// if `rules` is type of []string.
func (v *Validator) Messages(messages interface{}) *Validator {
	if messages == nil {
		return v
	}
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
	if m == nil {
		return v
	}
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
