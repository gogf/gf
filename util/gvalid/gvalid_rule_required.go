// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gvalid

import (
	"strings"
)

// 判断必须字段
func checkRequired(value, ruleKey, ruleVal string, params map[string]string) bool {
	required := false
	switch ruleKey {
	// 必须字段
	case "required":
		required = true

	// 必须字段(当任意所给定字段值与所给值相等时)
	case "required-if":
		required = false
		array := strings.Split(ruleVal, ",")
		// 必须为偶数，才能是键值对匹配
		if len(array)%2 == 0 {
			for i := 0; i < len(array); {
				tk := array[i]
				tv := array[i+1]
				if v, ok := params[tk]; ok {
					if strings.Compare(tv, v) == 0 {
						required = true
						break
					}
				}
				i += 2
			}
		}

	// 必须字段(当所给定字段值与所给值都不相等时)
	case "required-unless":
		required = true
		array := strings.Split(ruleVal, ",")
		// 必须为偶数，才能是键值对匹配
		if len(array)%2 == 0 {
			for i := 0; i < len(array); {
				tk := array[i]
				tv := array[i+1]
				if v, ok := params[tk]; ok {
					if strings.Compare(tv, v) == 0 {
						required = false
						break
					}
				}
				i += 2
			}
		}

	// 必须字段(当所给定任意字段值不为空时)
	case "required-with":
		required = false
		array := strings.Split(ruleVal, ",")
		for i := 0; i < len(array); i++ {
			if params[array[i]] != "" {
				required = true
				break
			}
		}

	// 必须字段(当所给定所有字段值都不为空时)
	case "required-with-all":
		required = true
		array := strings.Split(ruleVal, ",")
		for i := 0; i < len(array); i++ {
			if params[array[i]] == "" {
				required = false
				break
			}
		}

	// 必须字段(当所给定任意字段值为空时)
	case "required-without":
		required = false
		array := strings.Split(ruleVal, ",")
		for i := 0; i < len(array); i++ {
			if params[array[i]] == "" {
				required = true
				break
			}
		}

	// 必须字段(当所给定所有字段值都为空时)
	case "required-without-all":
		required = true
		array := strings.Split(ruleVal, ",")
		for i := 0; i < len(array); i++ {
			if params[array[i]] != "" {
				required = false
				break
			}
		}
	}
	if required {
		return !(value == "")
	} else {
		return true
	}
}
