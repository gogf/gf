// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gvalid

import (
	"strconv"
	"strings"
)

// 对字段值大小进行检测
func checkRange(value, ruleKey, ruleVal string, customMsgMap map[string]string) string {
	msg := ""
	switch ruleKey {
	// 大小范围
	case "between":
		array := strings.Split(ruleVal, ",")
		min := float64(0)
		max := float64(0)
		if len(array) > 0 {
			if v, err := strconv.ParseFloat(strings.TrimSpace(array[0]), 10); err == nil {
				min = v
			}
		}
		if len(array) > 1 {
			if v, err := strconv.ParseFloat(strings.TrimSpace(array[1]), 10); err == nil {
				max = v
			}
		}
		if v, err := strconv.ParseFloat(value, 10); err == nil {
			if v < min || v > max {
				if v, ok := customMsgMap[ruleKey]; !ok {
					msg = getDefaultErrorMessageByRule(ruleKey)
				} else {
					msg = v
				}
				msg = strings.Replace(msg, ":min", strconv.FormatFloat(min, 'f', -1, 64), -1)
				msg = strings.Replace(msg, ":max", strconv.FormatFloat(max, 'f', -1, 64), -1)
			}
		} else {
			msg = "输入参数[" + value + "]应当为数字类型"
		}

	// 最小值
	case "min":
		if min, err := strconv.ParseFloat(ruleVal, 10); err == nil {
			if v, err := strconv.ParseFloat(value, 10); err == nil {
				if v < min {
					if v, ok := customMsgMap[ruleKey]; !ok {
						msg = getDefaultErrorMessageByRule(ruleKey)
					} else {
						msg = v
					}
					msg = strings.Replace(msg, ":min", strconv.FormatFloat(min, 'f', -1, 64), -1)
				}
			} else {
				msg = "输入参数[" + value + "]应当为数字类型"
			}
		} else {
			msg = "校验参数[" + ruleVal + "]应当为数字类型"
		}

	// 最大值
	case "max":
		if max, err := strconv.ParseFloat(ruleVal, 10); err == nil {
			if v, err := strconv.ParseFloat(value, 10); err == nil {
				if v > max {
					if v, ok := customMsgMap[ruleKey]; !ok {
						msg = getDefaultErrorMessageByRule(ruleKey)
					} else {
						msg = v
					}
					msg = strings.Replace(msg, ":max", strconv.FormatFloat(max, 'f', -1, 64), -1)
				}
			} else {
				msg = "输入参数[" + value + "]应当为数字类型"
			}
		} else {
			msg = "校验参数[" + ruleVal + "]应当为数字类型"
		}
	}
	return msg
}
