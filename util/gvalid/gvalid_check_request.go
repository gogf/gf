// Copyright 2017-2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gvalid

import (
	"github.com/gogf/gf/net/ghttp"
	"reflect"
	"strings"

	"github.com/gogf/gf/util/gconv"
)

// 检测键值对参数Map，
// rules参数支持 []string / map[string]string 类型，前面一种类型支持返回校验结果顺序(具体格式参考struct tag)，后一种不支持；
// rules参数中得 map[string]string 是一个2维的关联数组，第一维键名为参数键名，第二维为带有错误的校验规则名称，值为错误信息。
func CheckRequest(cls interface{}, r *ghttp.Request, rules interface{}, msgs ...CustomMsg) *Error {
	//先获取request的入参
	params := r.GetRequestMap()
	// 将参数转换为 map[string]interface{}类型
	data := gconv.Map(params)
	if data == nil {
		return newErrorStr("invalid_params", "invalid params type: convert to map[string]interface{} failed")
	}
	// 真实校验规则数据结构
	checkRules := make(map[string]string)
	// 真实自定义错误信息数据结构
	customMsgs := make(CustomMsg)
	// 返回的顺序规则
	errorRules := make([]string, 0)
	// 返回的校验错误
	errorMaps := make(ErrorMap)
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

	// 不支持校验错误顺序: map[键名]校验规则
	case map[string]string:
		checkRules = v
	}
	// 自定义错误消息，非必须参数，优先级比rules参数中定义的错误消息更高
	if len(msgs) > 0 && len(msgs[0]) > 0 {
		if len(customMsgs) > 0 {
			for k, v := range msgs[0] {
				customMsgs[k] = v
			}
		} else {
			customMsgs = msgs[0]
		}
	}
	// 开始执行校验: 以校验规则作为基础进行遍历校验
	var value interface{}
	// 这里的rule变量为多条校验规则，不包含名字或者错误信息定义
	for key, rule := range checkRules {
		// 如果规则为空，那么不执行校验
		if len(rule) == 0 {
			continue
		}
		value = nil
		if v, ok := data[key]; ok {
			value = v
		}
		//自定义验证开始
		// 查看rule中是否有func[关键字 rule:func[customFunc]
		if strings.Contains(rule, "func[") {
			//正则获取funcName会不会更优雅些？
			funcRule := strings.Replace(rule, "func[", "",1)
			funcName := strings.Replace(funcRule,"]","",1)
			if len(funcName) > 0{
				res, errMsg := CallCustomValidator(cls, funcName, value)
				if !res {
					errorMaps[key] = make(map[string]string)
					errorMaps[key][rule] = errMsg
				}
			}
		}else{
			if e := Check(value, rule, customMsgs[key], data); e != nil {
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

	}
	if len(errorMaps) > 0 {
		return newError(errorRules, errorMaps)
	}
	return nil
}
// 用反射的方法调用对象的方法
// 要求自定义验证器必须返回两个参数第一个返回bool类型，第二个参数为验证失败时的错误信息(string)类型
func CallCustomValidator(any interface{}, name string, args ...interface{}) (bool, string) {
	inputs := make([]reflect.Value, len(args))
	for i := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}

	if v := reflect.ValueOf(any).MethodByName(name); v.String() != "<invalid Value>" {
		res := v.Call(inputs)
		return res[0].Bool(), res[1].String()
	}
	return false ,""
}