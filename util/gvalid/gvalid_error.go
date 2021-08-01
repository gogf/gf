// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gvalid

import (
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/text/gregex"
	"github.com/gogf/gf/text/gstr"
	"strings"
)

// Error is the validation error for validation result.
type Error interface {
	Code() int
	Current() error
	Error() string
	FirstItem() (key string, messages map[string]string)
	FirstRule() (rule string, err string)
	FirstString() (err string)
	Items() (items []map[string]map[string]string)
	Map() map[string]string
	Maps() map[string]map[string]string
	String() string
	Strings() (errs []string)
}

// validationError is the validation error for validation result.
type validationError struct {
	code      int                          // Error code.
	rules     []fieldRule                  // Rules by sequence, which is used for keeping error sequence.
	errors    map[string]map[string]string // Error map:map[field]map[rule]message
	firstKey  string                       // The first error rule key(empty in default).
	firstItem map[string]string            // The first error rule value(nil in default).
}

// newError creates and returns a validation error.
func newError(code int, rules []fieldRule, errors map[string]map[string]string) *validationError {
	for field, m := range errors {
		for k, v := range m {
			v = strings.Replace(v, ":attribute", field, -1)
			v, _ = gregex.ReplaceString(`\s{2,}`, ` `, v)
			v = gstr.Trim(v)
			m[k] = v
		}
		errors[field] = m
	}
	return &validationError{
		code:   code,
		rules:  rules,
		errors: errors,
	}
}

// newErrorStr creates and returns a validation error by string.
func newErrorStr(key, err string) *validationError {
	return newError(gerror.CodeInternalError, nil, map[string]map[string]string{
		internalErrorMapKey: {
			key: err,
		},
	})
}

// Code returns the error code of current validation error.
func (e *validationError) Code() int {
	if e == nil {
		return gerror.CodeNil
	}
	return e.code
}

// Map returns the first error message as map.
func (e *validationError) Map() map[string]string {
	if e == nil {
		return map[string]string{}
	}
	_, m := e.FirstItem()
	return m
}

// Maps returns all error messages as map.
func (e *validationError) Maps() map[string]map[string]string {
	if e == nil {
		return nil
	}
	return e.errors
}

// Items retrieves and returns error items array in sequence if possible,
// or else it returns error items with no sequence .
func (e *validationError) Items() (items []map[string]map[string]string) {
	if e == nil {
		return []map[string]map[string]string{}
	}
	items = make([]map[string]map[string]string, 0)
	// By sequence.
	if len(e.rules) > 0 {
		for _, v := range e.rules {
			if errorItemMap, ok := e.errors[v.Name]; ok {
				items = append(items, map[string]map[string]string{
					v.Name: errorItemMap,
				})
			}
		}
		return items
	}
	// No sequence.
	for name, errorRuleMap := range e.errors {
		items = append(items, map[string]map[string]string{
			name: errorRuleMap,
		})
	}
	return
}

// FirstItem returns the field name and error messages for the first validation rule error.
func (e *validationError) FirstItem() (key string, messages map[string]string) {
	if e == nil {
		return "", map[string]string{}
	}
	if e.firstItem != nil {
		return e.firstKey, e.firstItem
	}
	// By sequence.
	if len(e.rules) > 0 {
		for _, v := range e.rules {
			if errorItemMap, ok := e.errors[v.Name]; ok {
				e.firstKey = v.Name
				e.firstItem = errorItemMap
				return v.Name, errorItemMap
			}
		}
	}
	// No sequence.
	for k, m := range e.errors {
		e.firstKey = k
		e.firstItem = m
		return k, m
	}
	return "", nil
}

// FirstRule returns the first error rule and message string.
func (e *validationError) FirstRule() (rule string, err string) {
	if e == nil {
		return "", ""
	}
	// By sequence.
	if len(e.rules) > 0 {
		for _, v := range e.rules {
			if errorItemMap, ok := e.errors[v.Name]; ok {
				for _, ruleItem := range strings.Split(v.Rule, "|") {
					array := strings.Split(ruleItem, ":")
					ruleItem = strings.TrimSpace(array[0])
					if err, ok = errorItemMap[ruleItem]; ok {
						return ruleItem, err
					}
				}
			}
		}
	}
	// No sequence.
	for _, errorItemMap := range e.errors {
		for k, v := range errorItemMap {
			return k, v
		}
	}
	return "", ""
}

// FirstString returns the first error message as string.
// Note that the returned message might be different if it has no sequence.
func (e *validationError) FirstString() (err string) {
	if e == nil {
		return ""
	}
	_, err = e.FirstRule()
	return
}

// Current is alis of FirstString, which implements interface gerror.ApiCurrent.
func (e *validationError) Current() error {
	if e == nil {
		return nil
	}
	_, err := e.FirstRule()
	return gerror.NewCode(e.code, err)
}

// String returns all error messages as string, multiple error messages joined using char ';'.
func (e *validationError) String() string {
	if e == nil {
		return ""
	}
	return strings.Join(e.Strings(), "; ")
}

// Error implements interface of error.Error.
func (e *validationError) Error() string {
	if e == nil {
		return ""
	}
	return e.String()
}

// Strings returns all error messages as string array.
func (e *validationError) Strings() (errs []string) {
	if e == nil {
		return []string{}
	}
	errs = make([]string, 0)
	// By sequence.
	if len(e.rules) > 0 {
		for _, v := range e.rules {
			if errorItemMap, ok := e.errors[v.Name]; ok {
				// validation error checks.
				for _, ruleItem := range strings.Split(v.Rule, "|") {
					ruleItem = strings.TrimSpace(strings.Split(ruleItem, ":")[0])
					if err, ok := errorItemMap[ruleItem]; ok {
						errs = append(errs, err)
					}
				}
				// internal error checks.
				for k, _ := range internalErrKeyMap {
					if err, ok := errorItemMap[k]; ok {
						errs = append(errs, err)
					}
				}
			}
		}
		return errs
	}
	// No sequence.
	for _, errorItemMap := range e.errors {
		for _, err := range errorItemMap {
			errs = append(errs, err)
		}
	}
	return
}
