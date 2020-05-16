// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gvalid

import (
	"github.com/gogf/gf/util/gconv"
	"strconv"
	"strings"
)

// checkLength checks <value> using length rules.
// The length is calculated using unicode string, which means one chinese character or letter
// both has the length of 1.
func checkLength(value, ruleKey, ruleVal string, customMsgMap map[string]string) string {
	var (
		msg       = ""
		runeArray = gconv.Runes(value)
		valueLen  = len(runeArray)
	)
	switch ruleKey {
	case "length":
		var (
			min   = 0
			max   = 0
			array = strings.Split(ruleVal, ",")
		)
		if len(array) > 0 {
			if v, err := strconv.Atoi(strings.TrimSpace(array[0])); err == nil {
				min = v
			}
		}
		if len(array) > 1 {
			if v, err := strconv.Atoi(strings.TrimSpace(array[1])); err == nil {
				max = v
			}
		}
		if valueLen < min || valueLen > max {
			msg = getErrorMessageByRule(ruleKey, customMsgMap)
			msg = strings.Replace(msg, ":min", strconv.Itoa(min), -1)
			msg = strings.Replace(msg, ":max", strconv.Itoa(max), -1)
			return msg
		}

	case "min-length":
		min, err := strconv.Atoi(ruleVal)
		if valueLen < min || err != nil {
			msg = getErrorMessageByRule(ruleKey, customMsgMap)
			msg = strings.Replace(msg, ":min", strconv.Itoa(min), -1)
		}

	case "max-length":
		max, err := strconv.Atoi(ruleVal)
		if valueLen > max || err != nil {
			msg = getErrorMessageByRule(ruleKey, customMsgMap)
			msg = strings.Replace(msg, ":max", strconv.Itoa(max), -1)
		}
	}
	return msg
}
