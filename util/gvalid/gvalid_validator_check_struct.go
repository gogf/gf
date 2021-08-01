// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gvalid

import (
	"github.com/gogf/gf/errors/gerror"
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
		errorMaps           = make(map[string]map[string]string) // Returning error.
		fieldToAliasNameMap = make(map[string]string)            // Field name to alias name map.
	)
	fieldMap, err := structs.FieldMap(object, aliasNameTagPriority, true)
	if err != nil {
		return newErrorStr(internalObjectErrRuleName, err.Error())
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
		} else {
			if field.TagValue != "" {
				fieldToAliasNameMap[field.Name()] = field.TagValue
			}
		}
	}
	// It here must use structs.TagFields not structs.FieldMap to ensure error sequence.
	tagField, err := structs.TagFields(object, structTagPriority)
	if err != nil {
		return newErrorStr(internalObjectErrRuleName, err.Error())
	}
	// If there's no struct tag and validation rules, it does nothing and returns quickly.
	if len(tagField) == 0 && v.messages == nil {
		return nil
	}

	var (
		inputParamMap   map[string]interface{}
		checkRuleStrMap = make(map[string]string) // Complete rules map of struct: map[name]rule, the rule is complete pattern like: Name@RuleStr#Message
		customMessage   = make(CustomMsg)         // Custom rule error message map.
		checkValueData  = v.data                  // Ready to be validated data, which can be type of .
		errorRules      = make([]string, 0)       // Sequence rules.
	)
	if checkValueData == nil {
		checkValueData = object
	}
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
			checkRuleStrMap[name] = rule
			errorRules = append(errorRules, name+"@"+rule)
		}

	// Map type rules does not support sequence.
	// Format: map[key]rule
	case map[string]string:
		checkRuleStrMap = v
	}
	// If there's no struct tag and validation rules, it does nothing and returns quickly.
	if len(tagField) == 0 && len(checkRuleStrMap) == 0 {
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
		var (
			fieldName       = field.Name()                     // Attribute name.
			name, rule, msg = parseSequenceTag(field.TagValue) // The `name` is different from `attribute alias`, which is used for validation only.
		)
		if len(name) == 0 {
			if v, ok := fieldToAliasNameMap[fieldName]; ok {
				// It uses alias name of the attribute if its alias name tag exists.
				name = v
			} else {
				// It or else uses the attribute name directly.
				name = fieldName
			}
		} else {
			// It uses the alias name from validation rule.
			fieldToAliasNameMap[fieldName] = name
		}
		// It here extends the params map using alias names.
		// Note that the variable `name` might be alias name or attribute name.
		if _, ok := inputParamMap[name]; !ok {
			if !v.useDataInsteadOfObjectAttributes {
				inputParamMap[name] = field.Value.Interface()
			} else {
				if name != fieldName {
					if foundKey, foundValue := gutil.MapPossibleItemByKey(inputParamMap, fieldName); foundKey != "" {
						inputParamMap[name] = foundValue
					}
				}
			}
		}
		if _, ok := checkRuleStrMap[name]; !ok {
			if _, ok := checkRuleStrMap[fieldName]; ok {
				// If there's alias name,
				// use alias name as its key and remove the field name key.
				checkRuleStrMap[name] = checkRuleStrMap[fieldName]
				delete(checkRuleStrMap, fieldName)
			} else {
				checkRuleStrMap[name] = rule
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
			if a, ok := fieldToAliasNameMap[k]; ok {
				// Overwrite the key of field name.
				customMessage[a] = v
			} else {
				customMessage[k] = v
			}
		}
	}

	// The following logic is the same as some of CheckMap but with sequence support.
	var (
		value       interface{}
		hasBailRule bool
	)
	for key, rule := range checkRuleStrMap {
		_, value = gutil.MapPossibleItemByKey(inputParamMap, key)
		// It checks each rule and its value in loop.
		if validatedError := v.doCheckValue(key, value, rule, customMessage[key], checkValueData, inputParamMap); validatedError != nil {
			_, item := validatedError.FirstItem()
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
			if hasBailRule {
				break
			}
		}
	}
	if len(errorMaps) > 0 {
		return newError(gerror.CodeValidationFailed, errorRules, errorMaps)
	}
	return nil
}
