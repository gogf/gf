// Copyright 2017-2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gvalid

import (
	"strings"

	"github.com/gogf/gf/internal/structs"
	"github.com/gogf/gf/util/gconv"
)

var (
	structTagPriority    = []string{"gvalid", "valid", "v"} // structTagPriority specifies the validation tag priority array.
	aliasNameTagPriority = []string{"param", "params", "p"} // aliasNameTagPriority specifies the alias tag priority array.
)

// CheckStruct validates strcut and returns the error result.
//
// The parameter <object> should be type of struct/*struct.
// The parameter <rules> can be type of []string/map[string]string. It supports sequence in error result
// if <rules> is type of []string.
// The optional parameter <messages> specifies the custom error messages for specified keys and rules.
func CheckStruct(object interface{}, rules interface{}, messages ...CustomMsg) *Error {
	// It here must use structs.TagFields not structs.MapField to ensure error sequence.
	tagField, err := structs.TagFields(object, structTagPriority)
	if err != nil {
		return newErrorStr("invalid_object", err.Error())
	}
	// If there's no struct tag and validation rules, it does nothing and returns quickly.
	if len(tagField) == 0 && rules == nil {
		return nil
	}
	var (
		params        = make(map[string]interface{})
		checkRules    = make(map[string]string)
		customMessage = make(CustomMsg)
		fieldAliases  = make(map[string]string) // Alias names for <messages> overwriting struct tag names.
		errorRules    = make([]string, 0)       // Sequence rules.
		errorMaps     = make(ErrorMap)          // Returned error
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
					if _, ok := customMessage[name]; !ok {
						customMessage[name] = make(map[string]string)
					}
					customMessage[name].(map[string]string)[strings.TrimSpace(array[0])] = strings.TrimSpace(msgArray[k])
				}
			}
			checkRules[name] = rule
			errorRules = append(errorRules, name+"@"+rule)
		}

	// Map type rules does not support sequence.
	// Format: map[key]rule
	case map[string]string:
		checkRules = v
	}
	// If there's no struct tag and validation rules, it does nothing and returns quickly.
	if len(tagField) == 0 && len(checkRules) == 0 {
		return nil
	}
	// Checks and extends the parameters map with struct alias tag.
	mapField, err := structs.MapField(object, aliasNameTagPriority)
	if err != nil {
		return newErrorStr("invalid_object", err.Error())
	}
	for nameOrTag, field := range mapField {
		params[nameOrTag] = field.Value()
		params[field.Name()] = field.Value()
	}
	for _, field := range tagField {
		fieldName := field.Name()
		// sequence tag == struct tag
		// The name here is alias of field name.
		name, rule, msg := parseSequenceTag(field.TagValue)
		if len(name) == 0 {
			name = fieldName
		} else {
			fieldAliases[fieldName] = name
		}
		// It here extends the params map using alias names.
		if _, ok := params[name]; !ok {
			params[name] = field.Value()
		}
		if _, ok := checkRules[name]; !ok {
			if _, ok := checkRules[fieldName]; ok {
				// If there's alias name,
				// use alias name as its key and remove the field name key.
				checkRules[name] = checkRules[fieldName]
				delete(checkRules, fieldName)
			} else {
				checkRules[name] = rule
			}
			errorRules = append(errorRules, name+"@"+rule)
		} else {
			// The passed rules can overwrite the rules in struct tag.
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
				if _, ok := customMessage[name]; !ok {
					customMessage[name] = make(map[string]string)
				}
				customMessage[name].(map[string]string)[strings.TrimSpace(array[0])] = strings.TrimSpace(msgArray[k])
			}
		}
	}

	// Custom error messages,
	// which have the most priority than <rules> and struct tag.
	if len(messages) > 0 && len(messages[0]) > 0 {
		for k, v := range messages[0] {
			if a, ok := fieldAliases[k]; ok {
				// Overwrite the key of field name.
				customMessage[a] = v
			} else {
				customMessage[k] = v
			}
		}
	}

	// The following logic is the same as some of CheckMap.
	var value interface{}
	for key, rule := range checkRules {
		value = nil
		if v, ok := params[key]; ok {
			value = v
		}
		// It checks each rule and its value in loop.
		if e := doCheck(key, value, rule, customMessage[key], params); e != nil {
			_, item := e.FirstItem()
			// ===========================================================
			// Only in map and struct validations, if value is nil or empty
			// string and has no required* rules, it clears the error message.
			// ===========================================================
			if value == nil || gconv.String(value) == "" {
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
