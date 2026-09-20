// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gvalid implements powerful and useful data/form validation functionality.
package gvalid

import (
	"context"
	"reflect"
	"strings"
	"sync"

	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gtag"
)

// CustomMsg is the custom error message type,
// like: map[field] => string|map[rule]string
type CustomMsg = map[string]any

// fieldRule defined the alias name and rule string for specified field.
type fieldRule struct {
	Name      string       // Alias name for the field.
	Rule      string       // Rule string like: "max:6"
	IsMeta    bool         // Is this rule is from gmeta.Meta, which marks it as whole struct rule.
	FieldKind reflect.Kind // Original kind of struct field, which is used for parameter type checks.
	FieldType reflect.Type // Type of struct field, which is used for parameter type checks.
}

// iNoValidation is an interface that marks current struct not validated by package `gvalid`.
type iNoValidation interface {
	NoValidation()
}

const (
	internalRulesErrRuleName  = "InvalidRules"      // rule name for internal invalid rules validation error.
	internalParamsErrRuleName = "InvalidParams"     // rule name for internal invalid params validation error.
	internalObjectErrRuleName = "InvalidObject"     // rule name for internal invalid object validation error.
	internalErrorMapKey       = "__InternalError__" // error map key for internal errors.
	internalDefaultRuleName   = "__default__"       // default rule name for i18n error message format if no i18n message found for specified error rule.
	ruleMessagePrefixForI18n  = "gf.gvalid.rule."   // prefix string for each rule configuration in i18n content.
	noValidationTagName       = gtag.NoValidation   // no validation tag name for struct attribute.
	ruleNameRegex             = "regex"             // the name for rule "regex"
	ruleNameNotRegex          = "not-regex"         // the name for rule "not-regex"
	ruleNameForeach           = "foreach"           // the name for rule "foreach"
	ruleNameBail              = "bail"              // the name for rule "bail"
	ruleNameCi                = "ci"                // the name for rule "ci"
	emptyJsonArrayStr         = "[]"                // Empty json string for array type.
	emptyJsonObjectStr        = "{}"                // Empty json string for object type.
	requiredRulesPrefix       = "required"          // requiredRulesPrefix specifies the rule prefix that must be validated even the value is empty (nil or empty).
)

var (
	// defaultErrorMessages is the default error messages.
	// Note that these messages are synchronized from ./i18n/en/validation.toml .
	defaultErrorMessages = map[string]string{
		internalDefaultRuleName: "The {field} value `{value}` is invalid",
	}

	// structTagPriority specifies the validation tag priority array.
	structTagPriority = []string{gtag.Valid, gtag.ValidShort}

	// aliasNameTagPriority specifies the alias tag priority array.
	aliasNameTagPriority = []string{gtag.Param, gtag.ParamShort}

	// all internal error keys.
	internalErrKeyMap = map[string]string{
		internalRulesErrRuleName:  internalRulesErrRuleName,
		internalParamsErrRuleName: internalParamsErrRuleName,
		internalObjectErrRuleName: internalObjectErrRuleName,
	}

	// decorativeRuleMap defines all rules that are just marked rules which have neither functional meaning
	// nor error messages.
	decorativeRuleMap = map[string]bool{
		ruleNameForeach: true,
		ruleNameBail:    true,
		ruleNameCi:      true,
	}
)

// parsedTagValue is the parsed result of one validation tag value.
type parsedTagValue struct {
	field string
	rule  string
	msg   string
}

// parsedTagValueCache caches the parsed results of the validation tag values.
// The tag values are usually static for the same struct definitions, so their
// parse results are always the same, which makes it safe to cache. It is
// managed by the Validator cache switch, see function New.
var parsedTagValueCache sync.Map // map[string]parsedTagValue

// ParseTagValue parses one sequence tag to field, rule and error message.
// The sequence tag is like: [alias@]rule[...#msg...]
func ParseTagValue(tag string) (field, rule, msg string) {
	field, rule, msg, _ = parseTagValue(tag)
	return
}

// parseTagValueCached performs like ParseTagValue, but it caches the parsed
// results of the successful parses.
func parseTagValueCached(tag string) (field, rule, msg string) {
	if v, ok := parsedTagValueCache.Load(tag); ok {
		var parsed = v.(parsedTagValue)
		return parsed.field, parsed.rule, parsed.msg
	}
	var matched bool
	field, rule, msg, matched = parseTagValue(tag)
	if matched {
		parsedTagValueCache.Store(tag, parsedTagValue{
			field: field,
			rule:  rule,
			msg:   msg,
		})
	}
	return
}

// parseTagValue parses one sequence tag value, and it also returns whether the
// tag value matches the sequence tag pattern.
func parseTagValue(tag string) (field, rule, msg string, matched bool) {
	// Complete sequence tag.
	// Example: name@required|length:2,20|password3|same:password1#||Password strength is insufficient | Passwords are not match
	match, _ := gregex.MatchString(`\s*((\w+)\s*@){0,1}\s*([^#]+)\s*(#\s*(.*)){0,1}\s*`, tag)
	if len(match) > 5 {
		msg = strings.TrimSpace(match[5])
		rule = strings.TrimSpace(match[3])
		field = strings.TrimSpace(match[2])
		return field, rule, msg, true
	}
	intlog.Errorf(context.TODO(), `invalid validation tag value: %s`, tag)
	return "", "", "", false
}

// parseTagValue parses one validation tag value, using the parsed value cache
// unless the cache is disabled for current Validator.
func (v *Validator) parseTagValue(tag string) (field, rule, msg string) {
	if v.cache {
		return parseTagValueCached(tag)
	}
	return ParseTagValue(tag)
}

// parseRuleItem parses one rule item into rule key and rule pattern.
// The rule item is in format of: ruleKey[:rulePattern], for example: "max:6".
// Note that only the first char ':' is treated as the separator, and the rest
// part, which might contain chars like ':' and '|', belongs to the rule pattern.
// Examples:
//
//	"required"           -> ruleKey "required", rulePattern ""
//	"between:1,100"      -> ruleKey "between", rulePattern "1,100"
//	"regex:^\d+:\w+$"    -> ruleKey "regex", rulePattern "^\d+:\w+$"
func parseRuleItem(ruleItem string) (ruleKey, rulePattern string) {
	if index := strings.IndexByte(ruleItem, ':'); index != -1 {
		return gstr.Trim(ruleItem[:index]), gstr.Trim(ruleItem[index+1:])
	}
	return gstr.Trim(ruleItem), ""
}

// GetTags returns the validation tags.
func GetTags() []string {
	return structTagPriority
}
