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
	"strings"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/internal/reflection"
	"github.com/gogf/gf/v2/util/gconv"
)

func (v *Validator) doCheckMap(ctx context.Context, params interface{}) Error {
	if params == nil {
		return nil
	}
	var (
		checkRules    = make([]fieldRule, 0)
		customMessage = make(CustomMsg) // map[RuleKey]ErrorMsg.
		errorMaps     = make(map[string]map[string]error)
	)
	switch assertValue := v.rules.(type) {
	// Sequence tag: []sequence tag
	// Sequence has order for error results.
	case []string:
		for _, tag := range assertValue {
			name, rule, msg := ParseTagValue(tag)
			if len(name) == 0 {
				continue
			}
			if len(msg) > 0 {
				var (
					msgArray  = strings.Split(msg, "|")
					ruleArray = strings.Split(rule, "|")
				)
				for k, ruleItem := range ruleArray {
					// If length of custom messages is lesser than length of rules,
					// the rest rules use the default error messages.
					if len(msgArray) <= k {
						continue
					}
					if len(msgArray[k]) == 0 {
						continue
					}
					array := strings.Split(ruleItem, ":")
					if _, ok := customMessage[name]; !ok {
						customMessage[name] = make(map[string]string)
					}
					customMessage[name].(map[string]string)[strings.TrimSpace(array[0])] = strings.TrimSpace(msgArray[k])
				}
			}
			checkRules = append(checkRules, fieldRule{
				Name: name,
				Rule: rule,
			})
		}

	// No sequence rules: map[field]rule
	case map[string]string:
		for name, rule := range assertValue {
			checkRules = append(checkRules, fieldRule{
				Name: name,
				Rule: rule,
			})
		}
	}
	inputParamMap := gconv.Map(params)
	if inputParamMap == nil {
		return newValidationErrorByStr(
			internalParamsErrRuleName,
			errors.New("invalid params type: convert to map failed"),
		)
	}
	if msg, ok := v.messages.(CustomMsg); ok && len(msg) > 0 {
		if len(customMessage) > 0 {
			for k, v := range msg {
				customMessage[k] = v
			}
		} else {
			customMessage = msg
		}
	}
	var (
		value     interface{}
		validator = v.Clone()
	)

	// It checks the struct recursively if its attribute is an embedded struct.
	// Ignore inputParamMap, assoc, rules and messages from parent.
	validator.assoc = nil
	validator.rules = nil
	validator.messages = nil
	for _, item := range inputParamMap {
		originTypeAndKind := reflection.OriginTypeAndKind(item)
		switch originTypeAndKind.OriginKind {
		case reflect.Map, reflect.Struct, reflect.Slice, reflect.Array:
			v.doCheckValueRecursively(ctx, doCheckValueRecursivelyInput{
				Value:     item,
				Type:      originTypeAndKind.InputType,
				Kind:      originTypeAndKind.OriginKind,
				ErrorMaps: errorMaps,
			})
		}
		// Bail feature.
		if v.bail && len(errorMaps) > 0 {
			break
		}
	}
	if v.bail && len(errorMaps) > 0 {
		return newValidationError(gcode.CodeValidationFailed, nil, errorMaps)
	}

	// The following logic is the same as some of CheckStruct but without sequence support.
	for _, checkRuleItem := range checkRules {
		if len(checkRuleItem.Rule) == 0 {
			continue
		}
		value = nil
		if valueItem, ok := inputParamMap[checkRuleItem.Name]; ok {
			value = valueItem
		}
		// It checks each rule and its value in loop.
		if validatedError := v.doCheckValue(ctx, doCheckValueInput{
			Name:     checkRuleItem.Name,
			Value:    value,
			Rule:     checkRuleItem.Rule,
			Messages: customMessage[checkRuleItem.Name],
			DataRaw:  params,
			DataMap:  inputParamMap,
		}); validatedError != nil {
			_, errorItem := validatedError.FirstItem()
			// ===========================================================
			// Only in map and struct validations:
			// If value is nil or empty string and has no required* rules,
			// it clears the error message.
			// ===========================================================
			if gconv.String(value) == "" {
				required := false
				// rule => error
				for ruleKey := range errorItem {
					// Default required rules.
					if _, ok := mustCheckRulesEvenValueEmpty[ruleKey]; ok {
						required = true
						break
					}
				}
				if !required {
					continue
				}
			}
			if _, ok := errorMaps[checkRuleItem.Name]; !ok {
				errorMaps[checkRuleItem.Name] = make(map[string]error)
			}
			for ruleKey, ruleError := range errorItem {
				errorMaps[checkRuleItem.Name][ruleKey] = ruleError
			}
			if v.bail {
				break
			}
		}
	}
	if len(errorMaps) > 0 {
		return newValidationError(gcode.CodeValidationFailed, checkRules, errorMaps)
	}
	return nil
}
