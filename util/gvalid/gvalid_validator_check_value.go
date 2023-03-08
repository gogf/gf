// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gvalid

import (
	"context"
	"errors"
	"github.com/gogf/gf/v2/text/gregex"
	"reflect"
	"strings"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gvalid/internal/builtin"
)

type doCheckValueInput struct {
	Name     string                 // Name specifies the name of parameter `value`.
	Value    interface{}            // Value specifies the value for the rules to be validated.
	Rule     string                 // Rule specifies the validation rules string, like "required", "required|between:1,100", etc.
	Messages interface{}            // Messages specifies the custom error messages for this rule from parameters input, which is usually type of map/slice.
	DataRaw  interface{}            // DataRaw specifies the `raw data` which is passed to the Validator. It might be type of map/struct or a nil value.
	DataMap  map[string]interface{} // DataMap specifies the map that is converted from `dataRaw`. It is usually used internally
}

// doCheckSingleValue does the really rules validation for single key-value.
func (v *Validator) doCheckValue(ctx context.Context, in doCheckValueInput) Error {
	// If there's no validation rules, it does nothing and returns quickly.
	if in.Rule == "" {
		return nil
	}
	// It converts value to string and then does the validation.
	var (
		// Do not trim it as the space is also part of the value.
		ruleErrorMap = make(map[string]error)
	)
	// Custom error messages handling.
	var (
		msgArray     = make([]string, 0)
		customMsgMap = make(map[string]string)
	)
	switch messages := in.Messages.(type) {
	case string:
		msgArray = strings.Split(messages, "|")

	default:
		for k, message := range gconv.Map(in.Messages) {
			customMsgMap[k] = gconv.String(message)
		}
	}
	// Handle the char '|' in the rule,
	// which makes this rule separated into multiple rules.
	ruleItems := strings.Split(strings.TrimSpace(in.Rule), "|")
	for i := 0; ; {
		array := strings.Split(ruleItems[i], ":")
		if builtin.GetRule(array[0]) == nil && v.getCustomRuleFunc(array[0]) == nil {
			// ============================ SPECIAL ============================
			// Special `regex` and `not-regex` rules.
			// Merge the regex pattern if there are special chars, like ':', '|', in pattern.
			// ============================ SPECIAL ============================
			var (
				ruleNameRegexLengthMatch    bool
				ruleNameNotRegexLengthMatch bool
			)
			if i > 0 {
				ruleItem := ruleItems[i-1]
				if len(ruleItem) >= len(ruleNameRegex) && ruleItem[:len(ruleNameRegex)] == ruleNameRegex {
					ruleNameRegexLengthMatch = true
				}
				if len(ruleItem) >= len(ruleNameNotRegex) && ruleItem[:len(ruleNameNotRegex)] == ruleNameNotRegex {
					ruleNameNotRegexLengthMatch = true
				}
			}
			if i > 0 && (ruleNameRegexLengthMatch || ruleNameNotRegexLengthMatch) {
				ruleItems[i-1] += "|" + ruleItems[i]
				ruleItems = append(ruleItems[:i], ruleItems[i+1:]...)
			} else {
				return newValidationErrorByStr(
					internalRulesErrRuleName,
					errors.New(internalRulesErrRuleName+": "+ruleItems[i]),
				)
			}
		} else {
			i++
		}
		if i == len(ruleItems) {
			break
		}
	}
	var (
		hasBailRule        = v.bail
		hasForeachRule     = v.foreach
		hasCaseInsensitive = v.caseInsensitive
	)
	for index := 0; index < len(ruleItems); {
		var (
			err         error
			results     = ruleRegex.FindStringSubmatch(ruleItems[index]) // split single rule.
			ruleKey     = gstr.Trim(results[1])                          // rule key like "max" in rule "max: 6"
			rulePattern = gstr.Trim(results[2])                          // rule pattern is like "6" in rule:"max:6"
		)

		if !hasBailRule && ruleKey == ruleNameBail {
			hasBailRule = true
		}
		if !hasForeachRule && ruleKey == ruleNameForeach {
			hasForeachRule = true
		}
		if !hasCaseInsensitive && ruleKey == ruleNameCi {
			hasCaseInsensitive = true
		}

		// Ignore logic executing for marked rules.
		if decorativeRuleMap[ruleKey] {
			index++
			continue
		}

		if len(msgArray) > index {
			customMsgMap[ruleKey] = strings.TrimSpace(msgArray[index])
		}

		var (
			message        = v.getErrorMessageByRule(ctx, ruleKey, customMsgMap)
			customRuleFunc = v.getCustomRuleFunc(ruleKey)
			builtinRule    = builtin.GetRule(ruleKey)
			foreachValues  = []interface{}{in.Value}
		)
		if hasForeachRule {
			// As it marks `foreach`, so it converts the value to slice.
			foreachValues = gconv.Interfaces(in.Value)
			// Reset `foreach` rule as it only takes effect just once for next rule.
			hasForeachRule = false
		}

		for _, value := range foreachValues {
			switch {
			// Custom validation rules.
			case customRuleFunc != nil:
				if err = customRuleFunc(ctx, RuleFuncInput{
					Rule:    ruleItems[index],
					Message: message,
					Field:   in.Name,
					Value:   gvar.New(value),
					Data:    gvar.New(in.DataRaw),
				}); err != nil {
					// The error should have stack info to indicate the error position.
					if !gerror.HasStack(err) {
						err = gerror.New(err.Error())
					}
				}

			// Builtin validation rules.
			case customRuleFunc == nil && builtinRule != nil:
				err = builtinRule.Run(builtin.RunInput{
					RuleKey:     ruleKey,
					RulePattern: rulePattern,
					Field:       in.Name,
					Value:       gvar.New(value),
					Data:        gvar.New(in.DataRaw),
					Message:     message,
					Option: builtin.RunOption{
						CaseInsensitive: hasCaseInsensitive,
					},
				})

			default:
				// It never comes across here.
			}

			// Error handling.
			if err != nil {
				// The error should have error code that is `gcode.CodeValidationFailed`.
				if gerror.Code(err) == gcode.CodeNil {
					if e, ok := err.(*gerror.Error); ok {
						e.SetCode(gcode.CodeValidationFailed)
					}
				}

				// Error variable replacement for error message.
				s := gstr.ReplaceByMap(err.Error(), map[string]string{
					"{field}":     in.Name,             // Field name of the `value`.
					"{value}":     gconv.String(value), // Current validating value.
					"{pattern}":   rulePattern,         // The variable part of the rule.
					"{attribute}": in.Name,             // The same as `{field}`. It is deprecated.
				})
				s, _ = gregex.ReplaceString(`\s{2,}`, ` `, s)
				err = gerror.New(s)

				ruleErrorMap[ruleKey] = err

				// If it is with error and there's bail rule,
				// it then does not continue validating for left rules.
				if hasBailRule {
					goto CheckDone
				}
			}
		}
		index++
	}

CheckDone:
	if len(ruleErrorMap) > 0 {
		return newValidationError(
			gcode.CodeValidationFailed,
			[]fieldRule{{Name: in.Name, Rule: in.Rule}},
			map[string]map[string]error{
				in.Name: ruleErrorMap,
			},
		)
	}
	return nil
}

type doCheckValueRecursivelyInput struct {
	Value               interface{}                 // Value to be validated.
	Type                reflect.Type                // Struct/map/slice type which to be recursively validated.
	Kind                reflect.Kind                // Struct/map/slice kind to be asserted in following switch case.
	ErrorMaps           map[string]map[string]error // The validated failed error map.
	ResultSequenceRules *[]fieldRule                // The validated failed rule in sequence.
}

func (v *Validator) doCheckValueRecursively(ctx context.Context, in doCheckValueRecursivelyInput) {
	switch in.Kind {
	case reflect.Ptr:
		v.doCheckValueRecursively(ctx, doCheckValueRecursivelyInput{
			Value:               in.Value,
			Type:                in.Type.Elem(),
			Kind:                in.Type.Elem().Kind(),
			ErrorMaps:           in.ErrorMaps,
			ResultSequenceRules: in.ResultSequenceRules,
		})

	case reflect.Struct:
		// Ignore data, assoc, rules and messages from parent.
		var (
			validator           = v.Clone()
			toBeValidatedObject interface{}
		)
		if in.Type.Kind() == reflect.Ptr {
			toBeValidatedObject = reflect.New(in.Type.Elem()).Interface()
		} else {
			toBeValidatedObject = reflect.New(in.Type).Interface()
		}
		validator.assoc = nil
		validator.rules = nil
		validator.messages = nil
		if err := validator.Data(toBeValidatedObject).Assoc(in.Value).Run(ctx); err != nil {
			// It merges the errors into single error map.
			for k, m := range err.(*validationError).errors {
				in.ErrorMaps[k] = m
			}
			if in.ResultSequenceRules != nil {
				*in.ResultSequenceRules = append(*in.ResultSequenceRules, err.(*validationError).rules...)
			}
		}

	case reflect.Map:
		var (
			dataMap     = gconv.Map(in.Value)
			mapTypeElem = in.Type.Elem()
			mapTypeKind = mapTypeElem.Kind()
		)
		for _, item := range dataMap {
			v.doCheckValueRecursively(ctx, doCheckValueRecursivelyInput{
				Value:               item,
				Type:                mapTypeElem,
				Kind:                mapTypeKind,
				ErrorMaps:           in.ErrorMaps,
				ResultSequenceRules: in.ResultSequenceRules,
			})
			// Bail feature.
			if v.bail && len(in.ErrorMaps) > 0 {
				break
			}
		}

	case reflect.Slice, reflect.Array:
		var array []interface{}
		if gjson.Valid(in.Value) {
			array = gconv.Interfaces(gconv.Bytes(in.Value))
		} else {
			array = gconv.Interfaces(in.Value)
		}
		if len(array) == 0 {
			return
		}
		for _, item := range array {
			v.doCheckValueRecursively(ctx, doCheckValueRecursivelyInput{
				Value:               item,
				Type:                in.Type.Elem(),
				Kind:                in.Type.Elem().Kind(),
				ErrorMaps:           in.ErrorMaps,
				ResultSequenceRules: in.ResultSequenceRules,
			})
			// Bail feature.
			if v.bail && len(in.ErrorMaps) > 0 {
				break
			}
		}
	}
}
