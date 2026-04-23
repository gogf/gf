// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package builtin

import (
	"errors"
	"strings"
	"unicode/utf8"
)

// RuleContainsAny implements `contains-any` rule:
// Value should contain any of the specified characters.
//
// Format: contains-any:chars
type RuleContainsAny struct{}

func init() {
	Register(RuleContainsAny{})
	Register(RuleContainsRune{})
	Register(RuleExcludesAll{})
	Register(RuleStartsWith{})
	Register(RuleStartsNotWith{})
}

func (r RuleContainsAny) Name() string {
	return "contains-any"
}

func (r RuleContainsAny) Message() string {
	return "The {field} value `{value}` must contain any of {pattern}"
}

func (r RuleContainsAny) Run(in RunInput) error {
	if stringContainsAny(in.Value.String(), in.RulePattern, in.Option.CaseInsensitive) {
		return nil
	}
	return errors.New(in.Message)
}

// RuleContainsRune implements `contains-rune` rule:
// Value should contain the specified rune.
//
// Format: contains-rune:rune
type RuleContainsRune struct{}

func (r RuleContainsRune) Name() string {
	return "contains-rune"
}

func (r RuleContainsRune) Message() string {
	return "The {field} value `{value}` must contain rune {pattern}"
}

func (r RuleContainsRune) Run(in RunInput) error {
	if stringContainsRune(in.Value.String(), in.RulePattern, in.Option.CaseInsensitive) {
		return nil
	}
	return errors.New(in.Message)
}

// RuleExcludesAll implements `excludes-all` rule:
// Value should not contain any of the specified characters.
//
// Format: excludes-all:chars
type RuleExcludesAll struct{}

func (r RuleExcludesAll) Name() string {
	return "excludes-all"
}

func (r RuleExcludesAll) Message() string {
	return "The {field} value `{value}` must not contain any of {pattern}"
}

func (r RuleExcludesAll) Run(in RunInput) error {
	if !stringContainsAny(in.Value.String(), in.RulePattern, in.Option.CaseInsensitive) {
		return nil
	}
	return errors.New(in.Message)
}

// RuleStartsWith implements `starts-with` rule:
// Value should start with the specified pattern.
//
// Format: starts-with:pattern
type RuleStartsWith struct{}

func (r RuleStartsWith) Name() string {
	return "starts-with"
}

func (r RuleStartsWith) Message() string {
	return "The {field} value `{value}` must start with {pattern}"
}

func (r RuleStartsWith) Run(in RunInput) error {
	if stringHasPrefix(in.Value.String(), in.RulePattern, in.Option.CaseInsensitive) {
		return nil
	}
	return errors.New(in.Message)
}

// RuleStartsNotWith implements `starts-not-with` rule:
// Value should not start with the specified pattern.
//
// Format: starts-not-with:pattern
type RuleStartsNotWith struct{}

func (r RuleStartsNotWith) Name() string {
	return "starts-not-with"
}

func (r RuleStartsNotWith) Message() string {
	return "The {field} value `{value}` must not start with {pattern}"
}

func (r RuleStartsNotWith) Run(in RunInput) error {
	if !stringHasPrefix(in.Value.String(), in.RulePattern, in.Option.CaseInsensitive) {
		return nil
	}
	return errors.New(in.Message)
}

func normalizeStringMatch(value string, caseInsensitive bool) string {
	if caseInsensitive {
		return strings.ToLower(value)
	}
	return value
}

func stringContainsAny(value, pattern string, caseInsensitive bool) bool {
	value = normalizeStringMatch(value, caseInsensitive)
	pattern = normalizeStringMatch(pattern, caseInsensitive)
	return strings.ContainsAny(value, pattern)
}

func stringContainsRune(value, pattern string, caseInsensitive bool) bool {
	value = normalizeStringMatch(value, caseInsensitive)
	pattern = normalizeStringMatch(pattern, caseInsensitive)
	r, _ := utf8.DecodeRuneInString(pattern)
	if r == utf8.RuneError && pattern == "" {
		return false
	}
	return strings.ContainsRune(value, r)
}

func stringHasPrefix(value, pattern string, caseInsensitive bool) bool {
	value = normalizeStringMatch(value, caseInsensitive)
	pattern = normalizeStringMatch(pattern, caseInsensitive)
	return strings.HasPrefix(value, pattern)
}
