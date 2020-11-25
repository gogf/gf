// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gvalid

import (
	"errors"
	"github.com/gogf/gf/text/gregex"
	"strings"
)

// Error is the validation error for validation result.
type Error struct {
	rules     []string          // Rules by sequence, which is used for keeping error sequence.
	errors    ErrorMap          // Error map.
	firstKey  string            // The first error rule key(nil in default).
	firstItem map[string]string // The first error rule value(nil in default).
}

// ErrorMap is the validation error map:
// map[field]map[rule]message
type ErrorMap map[string]map[string]string

// newError creates and returns a validation error.
func newError(rules []string, errors map[string]map[string]string) *Error {
	for field, m := range errors {
		for k, v := range m {
			v = strings.Replace(v, ":attribute", field, -1)
			m[k], _ = gregex.ReplaceString(`\s{2,}`, ` `, v)
		}
		errors[field] = m
	}
	return &Error{
		rules:  rules,
		errors: errors,
	}
}

// newErrorStr creates and returns a validation error by string.
func newErrorStr(key, err string) *Error {
	return newError(nil, map[string]map[string]string{
		"__gvalid__": {
			key: err,
		},
	})
}

// Map returns the first error message as map.
func (e *Error) Map() map[string]string {
	if e == nil {
		return map[string]string{}
	}
	_, m := e.FirstItem()
	return m
}

// Maps returns all error messages as map.
func (e *Error) Maps() ErrorMap {
	if e == nil {
		return nil
	}
	return e.errors
}

// FirstItem returns the field name and error messages for the first validation rule error.
func (e *Error) FirstItem() (key string, messages map[string]string) {
	if e == nil {
		return "", map[string]string{}
	}
	if e.firstItem != nil {
		return e.firstKey, e.firstItem
	}
	// By sequence.
	if len(e.rules) > 0 {
		for _, v := range e.rules {
			name, _, _ := parseSequenceTag(v)
			if m, ok := e.errors[name]; ok {
				e.firstKey = name
				e.firstItem = m
				return name, m
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
func (e *Error) FirstRule() (rule string, err string) {
	if e == nil {
		return "", ""
	}
	// By sequence.
	if len(e.rules) > 0 {
		for _, v := range e.rules {
			name, rule, _ := parseSequenceTag(v)
			if m, ok := e.errors[name]; ok {
				for _, rule := range strings.Split(rule, "|") {
					array := strings.Split(rule, ":")
					rule = strings.TrimSpace(array[0])
					if err, ok := m[rule]; ok {
						return rule, err
					}
				}
			}
		}
	}
	// No sequence.
	for _, m := range e.errors {
		for k, v := range m {
			return k, v
		}
	}
	return "", ""
}

// FirstString returns the first error message as string.
// Note that the returned message might be different if it has no sequence.
func (e *Error) FirstString() (err string) {
	if e == nil {
		return ""
	}
	_, err = e.FirstRule()
	return
}

// Current is alis of FirstString, which implements interface gerror.ApiCurrent.
func (e *Error) Current() error {
	if e == nil {
		return nil
	}
	_, err := e.FirstRule()
	return errors.New(err)
}

// String returns all error messages as string, multiple error messages joined using char ';'.
func (e *Error) String() string {
	if e == nil {
		return ""
	}
	return strings.Join(e.Strings(), "; ")
}

// Error implements interface of error.Error.
func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	return e.String()
}

// Strings returns all error messages as string array.
func (e *Error) Strings() (errs []string) {
	if e == nil {
		return []string{}
	}
	errs = make([]string, 0)
	// By sequence.
	if len(e.rules) > 0 {
		for _, v := range e.rules {
			name, rule, _ := parseSequenceTag(v)
			if m, ok := e.errors[name]; ok {
				for _, rule := range strings.Split(rule, "|") {
					array := strings.Split(rule, ":")
					rule = strings.TrimSpace(array[0])
					if err, ok := m[rule]; ok {
						errs = append(errs, err)
					}
				}
			}
		}
		return errs
	}
	// No sequence.
	for _, m := range e.errors {
		for _, err := range m {
			errs = append(errs, err)
		}
	}
	return
}
