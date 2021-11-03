// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gvalid

import (
	"context"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/internal/structs"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gutil"
	"strings"
)

// CheckStruct validates struct and returns the error result.
// The parameter `object` should be type of struct/*struct.
func (v *Validator) CheckStruct(ctx context.Context, object interface{}) Error {
	return v.doCheckStruct(ctx, object)
}

func (v *Validator) doCheckStruct(ctx context.Context, object interface{}) Error {
	var (
		errorMaps           = make(map[string]map[string]string) // Returning error.
		fieldToAliasNameMap = make(map[string]string)            // Field names to alias name map.
	)
	fieldMap, err := structs.FieldMap(structs.FieldMapInput{
		Pointer:          object,
		PriorityTagArray: aliasNameTagPriority,
		RecursiveOption:  structs.RecursiveOptionEmbedded,
	})
	if err != nil {
		return newErrorStr(internalObjectErrRuleName, err.Error())
	}
	// It checks the struct recursively if its attribute is an embedded struct.
	for _, field := range fieldMap {
		if field.IsEmbedded() {
			// No validation interface implements check.
			if _, ok := field.Value.Interface().(iNoValidation); ok {
				continue
			}
			if _, ok := field.TagLookup(noValidationTagName); ok {
				continue
			}
			if err = v.doCheckStruct(ctx, field.Value); err != nil {
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
		inputParamMap  map[string]interface{}
		checkRules     = make([]fieldRule, 0)
		nameToRuleMap  = make(map[string]string) // just for internally searching index purpose.
		customMessage  = make(CustomMsg)         // Custom rule error message map.
		checkValueData = v.data                  // Ready to be validated data, which can be type of .
	)
	if checkValueData == nil {
		checkValueData = object
	}
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
			nameToRuleMap[name] = rule
			checkRules = append(checkRules, fieldRule{
				Name: name,
				Rule: rule,
			})
		}

	// Map type rules does not support sequence.
	// Format: map[key]rule
	case map[string]string:
		nameToRuleMap = assertValue
		for name, rule := range assertValue {
			checkRules = append(checkRules, fieldRule{
				Name: name,
				Rule: rule,
			})
		}
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

		if _, ok := nameToRuleMap[name]; !ok {
			if _, ok := nameToRuleMap[fieldName]; ok {
				// If there's alias name,
				// use alias name as its key and remove the field name key.
				nameToRuleMap[name] = nameToRuleMap[fieldName]
				delete(nameToRuleMap, fieldName)
				for index, checkRuleItem := range checkRules {
					if fieldName == checkRuleItem.Name {
						checkRuleItem.Name = name
						checkRules[index] = checkRuleItem
						break
					}
				}
			} else {
				nameToRuleMap[name] = rule
				checkRules = append(checkRules, fieldRule{
					Name: name,
					Rule: rule,
				})
			}
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
		value interface{}
	)
	for _, checkRuleItem := range checkRules {
		_, value = gutil.MapPossibleItemByKey(inputParamMap, checkRuleItem.Name)
		// It checks each rule and its value in loop.
		if validatedError := v.doCheckValue(ctx, doCheckValueInput{
			Name:     checkRuleItem.Name,
			Value:    value,
			Rule:     checkRuleItem.Rule,
			Messages: customMessage[checkRuleItem.Name],
			DataRaw:  checkValueData,
			DataMap:  inputParamMap,
		}); validatedError != nil {
			_, errorItem := validatedError.FirstItem()
			// ===================================================================
			// Only in map and struct validations, if value is nil or empty string
			// and has no required* rules, it clears the error message.
			// ===================================================================
			if value == nil || gconv.String(value) == "" {
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
