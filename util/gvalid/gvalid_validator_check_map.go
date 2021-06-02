// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gvalid

import (
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
		checkRules = make(map[string]string)
		customMsgs = make(CustomMsg)
		errorRules = make([]string, 0)
		errorMaps  = make(map[string]map[string]string)
	)
	switch v := v.rules.(type) {
	// Sequence tag: []sequence tag
	// Sequence has order for error results.
	case []string:
		for _, tag := range v {
			name, rule, msg := parseSequenceTag(tag)
			if len(name) == 0 {
				continue
			}
			if len(msg) > 0 {
				var (
					msgArray  = strings.Split(msg, "|")
					ruleArray = strings.Split(rule, "|")
				)
				for k, v := range ruleArray {
					// If length of custom messages is lesser than length of rules,
					// the rest rules use the default error messages.
					if len(msgArray) <= k {
						continue
					}
					if len(msgArray[k]) == 0 {
						continue
					}
					array := strings.Split(v, ":")
					if _, ok := customMsgs[name]; !ok {
						customMsgs[name] = make(map[string]string)
					}
					customMsgs[name].(map[string]string)[strings.TrimSpace(array[0])] = strings.TrimSpace(msgArray[k])
				}
			}
			checkRules[name] = rule
			errorRules = append(errorRules, name+"@"+rule)
		}

	// No sequence rules: map[field]rule
	case map[string]string:
		checkRules = v
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
		if len(customMsgs) > 0 {
			for k, v := range msg {
				customMsgs[k] = v
			}
		} else {
			customMsgs = msg
		}
	}
	var value interface{}
	for key, rule := range checkRules {
		if len(rule) == 0 {
			continue
		}
		value = nil
		if v, ok := data[key]; ok {
			value = v
		}
		// It checks each rule and its value in loop.
		if e := v.doCheckValue(key, value, rule, customMsgs[key], params, data); e != nil {
			_, item := e.FirstItem()
			// ===========================================================
			// Only in map and struct validations, if value is nil or empty
			// string and has no required* rules, it clears the error message.
			// ===========================================================
			if gconv.String(value) == "" {
				required := false
				// rule => error
				for k := range item {
					// Default required rules.
					if _, ok := mustCheckRulesEvenValueEmpty[k]; ok {
						required = true
						break
					}
					// Custom rules are also required in default.
					if f := v.getRuleFunc(k); f != nil {
						required = true
						break
					}
				}
				if !required {
					continue
				}
			}
			if _, ok := errorMaps[key]; !ok {
				errorMaps[key] = make(map[string]string)
			}
			for k, v := range item {
				errorMaps[key][k] = v
			}
		}
	}
	if len(errorMaps) > 0 {
		return newError(errorRules, errorMaps)
	}
	return nil
}
