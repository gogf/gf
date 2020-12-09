// Copyright 2017-2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gvalid

import (
	"strings"

	"github.com/gogf/gf/util/gconv"
)

// CheckMap validates map and returns the error result. It returns nil if with successful validation.
//
// The parameter <rules> can be type of []string/map[string]string. It supports sequence in error result
// if <rules> is type of []string.
// The optional parameter <messages> specifies the custom error messages for specified keys and rules.
func CheckMap(params interface{}, rules interface{}, messages ...CustomMsg) *Error {
	// If there's no validation rules, it does nothing and returns quickly.
	if params == nil || rules == nil {
		return nil
	}
	var (
		checkRules = make(map[string]string)
		customMsgs = make(CustomMsg)
		errorRules = make([]string, 0)
		errorMaps  = make(ErrorMap)
	)
	switch v := rules.(type) {
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
			"invalid_params",
			"invalid params type: convert to map failed",
		)
	}
	if len(messages) > 0 && len(messages[0]) > 0 {
		if len(customMsgs) > 0 {
			for k, v := range messages[0] {
				customMsgs[k] = v
			}
		} else {
			customMsgs = messages[0]
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
		if e := doCheck(key, value, rule, customMsgs[key], data); e != nil {
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
					if _, ok := customRuleFuncMap[k]; ok {
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
