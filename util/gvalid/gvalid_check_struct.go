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

// 校验struct对象属性，object参数也可以是一个指向对象的指针，返回值同CheckMap方法。
// struct的数据校验结果信息是顺序的。
func CheckStruct(object interface{}, rules interface{}, msgs ...interface{}) *Error {
	params := make(map[string]interface{})
	checkRules := make(map[string]string)
	customMsgs := make(CustomMsg)
	// 字段别名记录，用于msgs覆盖struct tag的
	fieldAliases := make(map[string]string)
	// 返回的顺序规则
	errorRules := make([]string, 0)
	// 返回的校验错误
	errorMaps := make(ErrorMap)
	// FieldRangeMsg 用于field-in/field-not-in的效验
	fieldRangeMsg := make(FieldRangeMsg)
	// 解析rules参数
	switch v := rules.(type) {
	// 支持校验错误顺序: []sequence tag
	case []string:
		for _, tag := range v {
			name, rule, msg := parseSequenceTag(tag)
			if len(name) == 0 {
				continue
			}
			// 错误提示
			if len(msg) > 0 {
				ruleArray := strings.Split(rule, "|")
				msgArray := strings.Split(msg, "|")
				for k, v := range ruleArray {
					// 如果msg条数比rule少，那么多余的rule使用默认的错误信息
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

	// Map type rules does not support sequence.
	// Format: map[key]rule
	case map[string]string:
		checkRules = v
	}
	// Checks and extends the parameters map with struct alias tag.
	for nameOrTag, field := range structs.MapField(object, aliasNameTagPriority, true) {
		params[nameOrTag] = field.Value()
		params[field.Name()] = field.Value()
	}
	// 首先, 按照属性循环一遍将struct的属性、数值、tag解析
	// It here must use structs.TagFields not structs.MapField to ensure error sequence.
	for _, field := range structs.TagFields(object, structTagPriority, true) {
		fieldName := field.Name()
		// sequence tag == struct tag, 这里的name为别名
		name, rule, msg := parseSequenceTag(field.Tag)
		if len(name) == 0 {
			name = fieldName
		} else {
			fieldAliases[fieldName] = name
		}
		// params参数使用别名**扩容**(而不仅仅使用别名)，仅用于验证使用
		if _, ok := params[name]; !ok {
			params[name] = field.Value()
		}
		// 校验规则
		if _, ok := checkRules[name]; !ok {
			if _, ok := checkRules[fieldName]; ok {
				// tag中存在别名，且rules传入的参数中使用了属性命名时，进行规则替换，并删除该属性的规则
				checkRules[name] = checkRules[fieldName]
				delete(checkRules, fieldName)
			} else {
				checkRules[name] = rule
			}
			errorRules = append(errorRules, name+"@"+rule)
		} else {
			// 传递的rules规则会覆盖struct tag的规则
			continue
		}
		// 错误提示
		if len(msg) > 0 {
			ruleArray := strings.Split(rule, "|")
			msgArray := strings.Split(msg, "|")
			for k, v := range ruleArray {
				// 如果msg条数比rule少，那么多余的rule使用默认的错误信息
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
	}

	// 自定义错误消息，非必须参数，优先级比rules参数中以及struct tag中定义的错误消息更高
	if len(msgs) > 0 {
		for _, msg := range msgs {
			switch msg.(type) {
			case CustomMsg:
				m := msg.(CustomMsg)
				if len(m) > 0 {
					for k, v := range m {
						if a, ok := fieldAliases[k]; ok {
							// 属性的别名存在时，覆盖别名的错误信息
							customMsgs[a] = v
						} else {
							customMsgs[k] = v
						}
					}
				}
			case FieldRangeMsg:
				m := msg.(FieldRangeMsg)
				if len(m) > 0 {
					fieldRangeMsg = m
				}
			}
		}
	}

	/* 以下逻辑和CheckMap相同 */

	// 开始执行校验: 以校验规则作为基础进行遍历校验
	var value interface{}
	// 这里的rule变量为多条校验规则，不包含名字或者错误信息定义
	for key, rule := range checkRules {
		value = nil
		if v, ok := params[key]; ok {
			value = v
		}
		if e := Check(value, rule, customMsgs[key], params, fieldRangeMsg[key]); e != nil {
			_, item := e.FirstItem()
			// 如果值为nil|""，并且不需要require*验证时，其他验证失效
			if value == nil || gconv.String(value) == "" {
				required := false
				// rule => error
				for k := range item {
					if _, ok := mustCheckRulesEvenValueEmpty[k]; ok {
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
