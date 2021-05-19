// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gvalid

import (
	"github.com/gogf/gf/internal/structs"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/gutil"
	"strings"
)

// CheckStruct validates struct and returns the error result.
// The parameter `object` should be type of struct/*struct.
func (v *Validator) CheckStruct(object interface{}) Error {
	return v.doCheckStruct(object)
}

func (v *Validator) doCheckStruct(object interface{}) Error {
	var (
		// Returning error.
		errorMaps = make(map[string]map[string]string)
	)
	fieldMap, err := structs.FieldMap(object, aliasNameTagPriority, true)
	if err != nil {
		return newErrorStr("invalid_object", err.Error())
	}
	// It checks the struct recursively the its attribute is an embedded struct.
	for _, field := range fieldMap {
		if field.IsEmbedded() {
			// No validation interface implements check.
			if _, ok := field.Value.Interface().(apiNoValidation); ok {
				continue
			}
			if _, ok := field.TagLookup(noValidationTagName); ok {
				continue
			}
			if err := v.doCheckStruct(field.Value); err != nil {
				// It merges the errors into single error map.
				for k, m := range err.(*validationError).errors {
					errorMaps[k] = m
				}
			}
		}
	}
	// It here must use structs.TagFields not structs.FieldMap to ensure error sequence.
	tagField, err := structs.TagFields(object, structTagPriority)
	if err != nil {
		return newErrorStr("invalid_object", err.Error())
	}
	// If there's no struct tag and validation rules, it does nothing and returns quickly.
	if len(tagField) == 0 && v.messages == nil {
		return nil
	}

	var (
		inputParamMap map[string]interface{}
		checkRules    = make(map[string]string)
		customMessage = make(CustomMsg)
		fieldAliases  = make(map[string]string) // Alias names for `messages` overwriting struct tag names.
		errorRules    = make([]string, 0)       // Sequence rules.
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
	// Input parameter map handling.
	if v.data == nil || !v.useDataInsteadOfObjectAttributes {
		inputParamMap = make(map[string]interface{})
	} else {
		inputParamMap = gconv.Map(v.data)
	}
	// Checks and extends the parameters map with struct alias tag.
	if !v.useDataInsteadOfObjectAttributes {
		for nameOrTag, field := range fieldMap {
			inputParamMap[nameOrTag] = field.Value.Interface()
			if nameOrTag != field.Name() {
				inputParamMap[field.Name()] = field.Value.Interface()
			}
		}
	}
	// Merge the custom validation rules with rules in struct tag.
	// The custom rules has the most high priority that can overwrite the struct tag rules.
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
		if _, ok := inputParamMap[name]; !ok {
			if !v.useDataInsteadOfObjectAttributes {
				inputParamMap[name] = field.Value.Interface()
			}
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
			// The input rules can overwrite the rules in struct tag.
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
	// which have the most priority than `rules` and struct tag.
	if msg, ok := v.messages.(CustomMsg); ok && len(msg) > 0 {
		for k, v := range msg {
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
		_, value = gutil.MapPossibleItemByKey(inputParamMap, key)
		// It checks each rule and its value in loop.
		if e := v.doCheckValue(key, value, rule, customMessage[key], inputParamMap); e != nil {
			_, item := e.FirstItem()
			// ===================================================================
			// Only in map and struct validations, if value is nil or empty string
			// and has no required* rules, it clears the error message.
			// ===================================================================
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
