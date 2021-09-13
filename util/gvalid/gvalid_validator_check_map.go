// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gvalid

import (
	"github.com/gogf/gf/errors/gcode"
	"github.com/gogf/gf/util/gconv"
	"strings"
)

// CheckMap validates map and returns the error result. It returns nil if with successful validation.
// The parameter `params` should be type of map.
func (v *Validator) CheckMap(params interface{}) Error {
	return v.doCheckMap(params)
}

func (v *Validator) doCheckMap(params interface{}) Error {
	// If there's no validation rules, it does nothing and returns quickly.
	if params == nil || v.rules == nil {
		return nil
	}
	var (
		checkRules    = make([]fieldRule, 0)
		customMessage = make(CustomMsg) // map[RuleKey]ErrorMsg.
		errorMaps     = make(map[string]map[string]string)
	)
	switch assertValue := v.rules.(type) {
	// Sequence tag: []sequence tag
	// Sequence has order for error results.
	case []string:
		for _, tag := range assertValue {
			name, rule, msg := parseSequenceTag(tag)
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
	// If there's no validation rules, it does nothing and returns quickly.
	if len(checkRules) == 0 {
		return nil
	}
	data := gconv.Map(params)
	if data == nil {
		return newErrorStr(
			internalParamsErrRuleName,
			"invalid params type: convert to map failed",
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
		value interface{}
	)
	for _, checkRuleItem := range checkRules {
		if len(checkRuleItem.Rule) == 0 {
			continue
		}
		value = nil
		if valueItem, ok := data[checkRuleItem.Name]; ok {
			value = valueItem
		}
		// It checks each rule and its value in loop.
		if validatedError := v.doCheckValue(doCheckValueInput{
			Name:     checkRuleItem.Name,
			Value:    value,
			Rule:     checkRuleItem.Rule,
			Messages: customMessage[checkRuleItem.Name],
			DataRaw:  params,
			DataMap:  data,
		}); validatedError != nil {
			_, errorItem := validatedError.FirstItem()
			// ===========================================================
			// Only in map and struct validations, if value is nil or empty
			// string and has no required* rules, it clears the error message.
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
					// Custom rules are also required in default.
					if f := v.getRuleFunc(ruleKey); f != nil {
						required = true
						break
					}
				}
				if !required {
					continue
				}
			}
			if _, ok := errorMaps[checkRuleItem.Name]; !ok {
				errorMaps[checkRuleItem.Name] = make(map[string]string)
			}
			for ruleKey, errorItemMsgMap := range errorItem {
				errorMaps[checkRuleItem.Name][ruleKey] = errorItemMsgMap
			}
			if v.bail {
				break
			}
		}
	}
	if len(errorMaps) > 0 {
		return newError(gcode.CodeValidationFailed, checkRules, errorMaps)
	}
	return nil
}
