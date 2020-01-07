// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gvalid

import "strings"

// 校验错误对象
type Error struct {
	rules     []string          // 校验结果顺序(可能为nil)，可保证返回校验错误的顺序性
	errors    ErrorMap          // 完整的数据校验结果存储(map无序)
	firstKey  string            // 第一条错误项键名(常用操作冗余数据)，默认为空
	firstItem map[string]string // 第一条错误项(常用操作冗余数据)，默认为nil
}

// 校验错误信息: map[键名]map[规则名]错误信息
type ErrorMap map[string]map[string]string

// 创建一个校验错误对象指针(校验错误)
func newError(rules []string, errors map[string]map[string]string) *Error {
	return &Error{
		rules:  rules,
		errors: errors,
	}
}

// 创建一个校验错误对象指针(内部错误)
func newErrorStr(key, err string) *Error {
	return &Error{
		rules: nil,
		errors: map[string]map[string]string{
			"__gvalid__": {
				key: err,
			},
		},
	}
}

// 获得规则与错误信息的map; 当校验结果为多条数据校验时，返回第一条错误map(此时类似FirstItem)
func (e *Error) Map() map[string]string {
	_, m := e.FirstItem()
	return m
}

// 获得原始校验结果ErrorMap
func (e *Error) Maps() ErrorMap {
	return e.errors
}

// 只获取第一个键名的校验错误项
func (e *Error) FirstItem() (key string, msgs map[string]string) {
	if e.firstItem != nil {
		return e.firstKey, e.firstItem
	}
	// 有序
	if len(e.rules) > 0 {
		for _, v := range e.rules {
			name, _, _ := parseSequenceTag(v)
			if m, ok := e.errors[name]; ok {
				e.firstKey = name
				e.firstItem = m
				return name, m
			}
		}
	}
	// 无序
	for k, m := range e.errors {
		e.firstKey = k
		e.firstItem = m
		return k, m
	}
	return "", nil
}

// 只获取第一个校验错误项的规则及错误信息
func (e *Error) FirstRule() (rule string, err string) {
	// 有序
	if len(e.rules) > 0 {
		for _, v := range e.rules {
			name, rule, _ := parseSequenceTag(v)
			if m, ok := e.errors[name]; ok {
				for _, rule := range strings.Split(rule, "|") {
					array := strings.Split(rule, ":")
					rule = strings.TrimSpace(array[0])
					if err, ok := m[rule]; ok {
						return rule, err
					}
				}
			}
		}
	}
	// 无序
	for _, m := range e.errors {
		for k, v := range m {
			return k, v
		}
	}
	return "", ""
}

// 只获取第一个校验错误项的错误信息
func (e *Error) FirstString() (err string) {
	_, err = e.FirstRule()
	return
}

// 将所有错误信息构建称字符串，多个错误信息字符串使用"; "符号分隔
func (e *Error) String() string {
	return strings.Join(e.Strings(), "; ")
}

// Error implements interface of error.Error.
func (e *Error) Error() string {
	return e.String()
}

// 只返回错误信息，构造成字符串数组返回
func (e *Error) Strings() (errs []string) {
	errs = make([]string, 0)
	// 有序
	if len(e.rules) > 0 {
		for _, v := range e.rules {
			name, rule, _ := parseSequenceTag(v)
			if m, ok := e.errors[name]; ok {
				for _, rule := range strings.Split(rule, "|") {
					array := strings.Split(rule, ":")
					rule = strings.TrimSpace(array[0])
					if err, ok := m[rule]; ok {
						errs = append(errs, err)
					}
				}
			}
		}
		return errs
	}
	// 无序
	for _, m := range e.errors {
		for _, err := range m {
			errs = append(errs, err)
		}
	}
	return
}
