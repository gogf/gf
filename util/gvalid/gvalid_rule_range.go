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

// checkRange checks <value> using range rules.
func checkRange(value, ruleKey, ruleVal string, customMsgMap map[string]string) string {
	msg := ""
	switch ruleKey {
	// Value range.
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
		v, err := strconv.ParseFloat(value, 10)
		if v < min || v > max || err != nil {
			msg = getErrorMessageByRule(ruleKey, customMsgMap)
			msg = strings.Replace(msg, ":min", strconv.FormatFloat(min, 'f', -1, 64), -1)
			msg = strings.Replace(msg, ":max", strconv.FormatFloat(max, 'f', -1, 64), -1)
		}

	// Min value.
	case "min":
		var (
			min, err1    = strconv.ParseFloat(ruleVal, 10)
			valueN, err2 = strconv.ParseFloat(value, 10)
		)
		if valueN < min || err1 != nil || err2 != nil {
			msg = getErrorMessageByRule(ruleKey, customMsgMap)
			msg = strings.Replace(msg, ":min", strconv.FormatFloat(min, 'f', -1, 64), -1)
		}

	// Max value.
	case "max":
		var (
			max, err1    = strconv.ParseFloat(ruleVal, 10)
			valueN, err2 = strconv.ParseFloat(value, 10)
		)
		if valueN > max || err1 != nil || err2 != nil {
			msg = getErrorMessageByRule(ruleKey, customMsgMap)
			msg = strings.Replace(msg, ":max", strconv.FormatFloat(max, 'f', -1, 64), -1)
		}

	}
	return msg
}
