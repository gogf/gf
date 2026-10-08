// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package builtin implements built-in validation rules.
//
// Referred to Laravel validation:
// https://laravel.com/docs/master/validation#available-validation-rules
package builtin

import (
	"reflect"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/util/gconv"
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
	RuleKey     string       // RuleKey is like the "max" in rule "max: 6"
	RulePattern string       // RulePattern is like "6" in rule:"max:6"
	Field       string       // The field name of Value.
	ValueType   reflect.Type // ValueType specifies the type of the value, which might be nil.
	Value       *gvar.Var    // Value specifies the value for this rule to validate.
	Data        *gvar.Var    // Data specifies the `data` which is passed to the Validator.
	Message     string       // Message specifies the custom error message or configured i18n message for this rule.
	Option      RunOption    // Option provides extra configuration for validation rule.
}

type RunOption struct {
	CaseInsensitive bool // CaseInsensitive indicates that it does Case-Insensitive comparison in string.
}

var (
	// ruleMap stores all builtin validation rules.
	ruleMap = map[string]Rule{}
)

// Register registers builtin rule into manager.
func Register(rule Rule) {
	ruleMap[rule.Name()] = rule
}

// GetRule retrieves and returns rule by `name`.
func GetRule(name string) Rule {
	return ruleMap[name]
}

// valueLength returns the length of `value` for the length related rules, that is
// `length`, `min-length`, `max-length` and `size`.
//
// For a slice, an array or a map it returns the number of elements. That is the same
// notion of length that `required` already applies to the very same value, see
// isRequiredEmpty, so that both rules in a tag like `required|min-length:2` agree on
// what the length of a container is.
//
// For anything else it returns the number of unicode runes of its string form, which
// keeps the historical behavior: one chinese character or letter both has the length
// of 1.
//
// A byte slice is deliberately measured as a string. gconv converts it to a string
// directly instead of serializing it, so `[]byte("GoFrame")` has always had the length
// of its text and keeps it here.
func valueLength(value any) int {
	reflectValue := reflect.ValueOf(value)
	for reflectValue.Kind() == reflect.Pointer {
		reflectValue = reflectValue.Elem()
	}
	switch reflectValue.Kind() {
	case reflect.Slice, reflect.Array:
		if reflectValue.Type().Elem().Kind() == reflect.Uint8 {
			break
		}
		return reflectValue.Len()

	case reflect.Map:
		return reflectValue.Len()

	default:
	}
	return len(gconv.Runes(gconv.String(value)))
}
