// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gvalid

import (
	"github.com/gogf/gf/errors/gerror"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/internal/json"
	"github.com/gogf/gf/net/gipv4"
	"github.com/gogf/gf/net/gipv6"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/text/gregex"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/gutil"
)

type apiTime interface {
	Date() (year int, month time.Month, day int)
	IsZero() bool
}

// CheckValue checks single value with specified rules.
// It returns nil if successful validation.
func (v *Validator) CheckValue(value interface{}) Error {
	return v.doCheckValue("", value, gconv.String(v.rules), v.messages, v.data, gconv.Map(v.data))
}

// doCheckSingleValue does the really rules validation for single key-value.
//
// The parameter `rules` specifies the validation rules string, like "required", "required|between:1,100", etc.
// The parameter `value` specifies the value for this rules to be validated.
// The parameter `messages` specifies the custom error messages for this rule, which is usually type of map/slice.
// The parameter `dataRaw` specifies the `raw data` which is passed to the Validator. It might be type of map/struct or a nil value.
// The parameter `dataMap` specifies the map that is converted from `dataRaw`. It is usually used internally
func (v *Validator) doCheckValue(
	key string,
	value interface{},
	rules string,
	messages interface{},
	dataRaw interface{},
	dataMap map[string]interface{},
) Error {
	// If there's no validation rules, it does nothing and returns quickly.
	if rules == "" {
		return nil
	}
	// It converts value to string and then does the validation.
	var (
		// Do not trim it as the space is also part of the value.
		errorMsgArray = make(map[string]string)
	)
	// Custom error messages handling.
	var (
		msgArray     = make([]string, 0)
		customMsgMap = make(map[string]string)
	)
	switch v := messages.(type) {
	case string:
		msgArray = strings.Split(v, "|")
	default:
		for k, v := range gconv.Map(messages) {
			customMsgMap[k] = gconv.String(v)
		}
	}
	// Handle the char '|' in the rule,
	// which makes this rule separated into multiple rules.
	ruleItems := strings.Split(strings.TrimSpace(rules), "|")
	for i := 0; ; {
		array := strings.Split(ruleItems[i], ":")
		_, ok := allSupportedRules[array[0]]
		if !ok && v.getRuleFunc(array[0]) == nil {
			if i > 0 && ruleItems[i-1][:5] == "regex" {
				ruleItems[i-1] += "|" + ruleItems[i]
				ruleItems = append(ruleItems[:i], ruleItems[i+1:]...)
			} else {
				return newErrorStr(
					internalRulesErrRuleName,
					internalRulesErrRuleName+": "+rules,
				)
			}
		} else {
			i++
		}
		if i == len(ruleItems) {
			break
		}
	}
	for index := 0; index < len(ruleItems); {
		var (
			err            error
			match          = false
			results        = ruleRegex.FindStringSubmatch(ruleItems[index])
			ruleKey        = strings.TrimSpace(results[1])
			rulePattern    = strings.TrimSpace(results[2])
			customRuleFunc RuleFunc
		)
		if len(msgArray) > index {
			customMsgMap[ruleKey] = strings.TrimSpace(msgArray[index])
		}

		// Custom rule handling.
		// 1. It firstly checks and uses the custom registered rules functions in the current Validator.
		// 2. It secondly checks and uses the globally registered rules functions.
		// 3. It finally checks and uses the build-in rules functions.
		customRuleFunc = v.getRuleFunc(ruleKey)
		if customRuleFunc != nil {
			// It checks custom validation rules with most priority.
			message := v.getErrorMessageByRule(ruleKey, customMsgMap)
			if err := customRuleFunc(v.ctx, ruleItems[index], value, message, dataRaw); err != nil {
				match = false
				errorMsgArray[ruleKey] = err.Error()
			} else {
				match = true
			}
		} else {
			// It checks build-in validation rules if there's no custom rule.
			match, err = v.doCheckBuildInRules(index, value, ruleKey, rulePattern, ruleItems, dataMap, customMsgMap)
			if !match && err != nil {
				errorMsgArray[ruleKey] = err.Error()
			}
		}

		// Error message handling.
		if !match {
			// It does nothing if the error message for this rule
			// is already set in previous validation.
			if _, ok := errorMsgArray[ruleKey]; !ok {
				errorMsgArray[ruleKey] = v.getErrorMessageByRule(ruleKey, customMsgMap)
			}
		}
		index++
	}
	if len(errorMsgArray) > 0 {
		return newError([]string{rules}, map[string]map[string]string{
			key: errorMsgArray,
		})
	}
	return nil
}

func (v *Validator) doCheckBuildInRules(
	index int,
	value interface{},
	ruleKey string,
	rulePattern string,
	ruleItems []string,
	dataMap map[string]interface{},
	customMsgMap map[string]string,
) (match bool, err error) {
	valueStr := gconv.String(value)
	switch ruleKey {
	// Required rules.
	case
		"required",
		"required-if",
		"required-unless",
		"required-with",
		"required-with-all",
		"required-without",
		"required-without-all":
		match = v.checkRequired(value, ruleKey, rulePattern, dataMap)

	// Length rules.
	// It also supports length of unicode string.
	case
		"length",
		"min-length",
		"max-length",
		"size":
		if msg := v.checkLength(valueStr, ruleKey, rulePattern, customMsgMap); msg != "" {
			return match, gerror.New(msg)
		} else {
			match = true
		}

	// Range rules.
	case
		"min",
		"max",
		"between":
		if msg := v.checkRange(valueStr, ruleKey, rulePattern, customMsgMap); msg != "" {
			return match, gerror.New(msg)
		} else {
			match = true
		}

	// Custom regular expression.
	case "regex":
		// It here should check the rule as there might be special char '|' in it.
		for i := index + 1; i < len(ruleItems); i++ {
			if !gregex.IsMatchString(singleRulePattern, ruleItems[i]) {
				rulePattern += "|" + ruleItems[i]
				index++
			}
		}
		match = gregex.IsMatchString(rulePattern, valueStr)

	// Date rules.
	case "date":
		// support for time value, eg: gtime.Time/*gtime.Time, time.Time/*time.Time.
		if v, ok := value.(apiTime); ok {
			return !v.IsZero(), nil
		}
		match = gregex.IsMatchString(`\d{4}[\.\-\_/]{0,1}\d{2}[\.\-\_/]{0,1}\d{2}`, valueStr)

	// Date rule with specified format.
	case "date-format":
		// support for time value, eg: gtime.Time/*gtime.Time, time.Time/*time.Time.
		if v, ok := value.(apiTime); ok {
			return !v.IsZero(), nil
		}
		if _, err := gtime.StrToTimeFormat(valueStr, rulePattern); err == nil {
			match = true
		} else {
			var msg string
			msg = v.getErrorMessageByRule(ruleKey, customMsgMap)
			msg = strings.Replace(msg, ":format", rulePattern, -1)
			return match, gerror.New(msg)
		}

	// Values of two fields should be equal as string.
	case "same":
		_, foundValue := gutil.MapPossibleItemByKey(dataMap, rulePattern)
		if foundValue != nil {
			if strings.Compare(valueStr, gconv.String(foundValue)) == 0 {
				match = true
			}
		}
		if !match {
			var msg string
			msg = v.getErrorMessageByRule(ruleKey, customMsgMap)
			msg = strings.Replace(msg, ":field", rulePattern, -1)
			return match, gerror.New(msg)
		}

	// Values of two fields should not be equal as string.
	case "different":
		match = true
		_, foundValue := gutil.MapPossibleItemByKey(dataMap, rulePattern)
		if foundValue != nil {
			if strings.Compare(valueStr, gconv.String(foundValue)) == 0 {
				match = false
			}
		}
		if !match {
			var msg string
			msg = v.getErrorMessageByRule(ruleKey, customMsgMap)
			msg = strings.Replace(msg, ":field", rulePattern, -1)
			return match, gerror.New(msg)
		}

	// Field value should be in range of.
	case "in":
		array := strings.Split(rulePattern, ",")
		for _, v := range array {
			if strings.Compare(valueStr, strings.TrimSpace(v)) == 0 {
				match = true
				break
			}
		}

	// Field value should not be in range of.
	case "not-in":
		match = true
		array := strings.Split(rulePattern, ",")
		for _, v := range array {
			if strings.Compare(valueStr, strings.TrimSpace(v)) == 0 {
				match = false
				break
			}
		}

	// Phone format validation.
	// 1. China Mobile:
	//    134, 135, 136, 137, 138, 139, 150, 151, 152, 157, 158, 159, 182, 183, 184, 187, 188,
	//    178(4G), 147(Net)；
	//    172
	//
	// 2. China Unicom:
	//    130, 131, 132, 155, 156, 185, 186 ,176(4G), 145(Net), 175
	//
	// 3. China Telecom:
	//    133, 153, 180, 181, 189, 177(4G)
	//
	// 4. Satelite:
	//    1349
	//
	// 5. Virtual:
	//    170, 173
	//
	// 6. 2018:
	//    16x, 19x
	case "phone":
		match = gregex.IsMatchString(`^13[\d]{9}$|^14[5,7]{1}\d{8}$|^15[^4]{1}\d{8}$|^16[\d]{9}$|^17[0,2,3,5,6,7,8]{1}\d{8}$|^18[\d]{9}$|^19[\d]{9}$`, valueStr)

	// Loose mobile phone number verification(宽松的手机号验证)
	// As long as the 11 digit numbers beginning with
	// 13, 14, 15, 16, 17, 18, 19 can pass the verification (只要满足 13、14、15、16、17、18、19开头的11位数字都可以通过验证)
	case "phone-loose":
		match = gregex.IsMatchString(`^1(3|4|5|6|7|8|9)\d{9}$`, valueStr)

	// Telephone number:
	// "XXXX-XXXXXXX"
	// "XXXX-XXXXXXXX"
	// "XXX-XXXXXXX"
	// "XXX-XXXXXXXX"
	// "XXXXXXX"
	// "XXXXXXXX"
	case "telephone":
		match = gregex.IsMatchString(`^((\d{3,4})|\d{3,4}-)?\d{7,8}$`, valueStr)

	// QQ number: from 10000.
	case "qq":
		match = gregex.IsMatchString(`^[1-9][0-9]{4,}$`, valueStr)

	// Postcode number.
	case "postcode":
		match = gregex.IsMatchString(`^\d{6}$`, valueStr)

	// China resident id number.
	//
	// xxxxxx yyyy MM dd 375 0  十八位
	// xxxxxx   yy MM dd  75 0  十五位
	//
	// 地区：     [1-9]\d{5}
	// 年的前两位：(18|19|([23]\d))  1800-2399
	// 年的后两位：\d{2}
	// 月份：     ((0[1-9])|(10|11|12))
	// 天数：     (([0-2][1-9])|10|20|30|31) 闰年不能禁止29+
	//
	// 三位顺序码：\d{3}
	// 两位顺序码：\d{2}
	// 校验码：   [0-9Xx]
	//
	// 十八位：^[1-9]\d{5}(18|19|([23]\d))\d{2}((0[1-9])|(10|11|12))(([0-2][1-9])|10|20|30|31)\d{3}[0-9Xx]$
	// 十五位：^[1-9]\d{5}\d{2}((0[1-9])|(10|11|12))(([0-2][1-9])|10|20|30|31)\d{3}$
	//
	// 总：
	// (^[1-9]\d{5}(18|19|([23]\d))\d{2}((0[1-9])|(10|11|12))(([0-2][1-9])|10|20|30|31)\d{3}[0-9Xx]$)|(^[1-9]\d{5}\d{2}((0[1-9])|(10|11|12))(([0-2][1-9])|10|20|30|31)\d{3}$)
	case "resident-id":
		match = v.checkResidentId(valueStr)

	// Bank card number using LUHN algorithm.
	case "bank-card":
		match = v.checkLuHn(valueStr)

	// Universal passport format rule:
	// Starting with letter, containing only numbers or underscores, length between 6 and 18.
	case "passport":
		match = gregex.IsMatchString(`^[a-zA-Z]{1}\w{5,17}$`, valueStr)

	// Universal password format rule1:
	// Containing any visible chars, length between 6 and 18.
	case "password":
		match = gregex.IsMatchString(`^[\w\S]{6,18}$`, valueStr)

	// Universal password format rule2:
	// Must meet password rule1, must contain lower and upper letters and numbers.
	case "password2":
		if gregex.IsMatchString(`^[\w\S]{6,18}$`, valueStr) &&
			gregex.IsMatchString(`[a-z]+`, valueStr) &&
			gregex.IsMatchString(`[A-Z]+`, valueStr) &&
			gregex.IsMatchString(`\d+`, valueStr) {
			match = true
		}

	// Universal password format rule3:
	// Must meet password rule1, must contain lower and upper letters, numbers and special chars.
	case "password3":
		if gregex.IsMatchString(`^[\w\S]{6,18}$`, valueStr) &&
			gregex.IsMatchString(`[a-z]+`, valueStr) &&
			gregex.IsMatchString(`[A-Z]+`, valueStr) &&
			gregex.IsMatchString(`\d+`, valueStr) &&
			gregex.IsMatchString(`[^a-zA-Z0-9]+`, valueStr) {
			match = true
		}

	// Json.
	case "json":
		if json.Valid([]byte(valueStr)) {
			match = true
		}

	// Integer.
	case "integer":
		if _, err := strconv.Atoi(valueStr); err == nil {
			match = true
		}

	// Float.
	case "float":
		if _, err := strconv.ParseFloat(valueStr, 10); err == nil {
			match = true
		}

	// Boolean(1,true,on,yes:true | 0,false,off,no,"":false).
	case "boolean":
		match = false
		if _, ok := boolMap[strings.ToLower(valueStr)]; ok {
			match = true
		}

	// Email.
	case "email":
		match = gregex.IsMatchString(`^[a-zA-Z0-9_\-\.]+@[a-zA-Z0-9_\-]+(\.[a-zA-Z0-9_\-]+)+$`, valueStr)

	// URL
	case "url":
		match = gregex.IsMatchString(`(https?|ftp|file)://[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|]`, valueStr)

	// Domain
	case "domain":
		match = gregex.IsMatchString(`^([0-9a-zA-Z][0-9a-zA-Z\-]{0,62}\.)+([a-zA-Z]{0,62})$`, valueStr)

	// IP(IPv4/IPv6).
	case "ip":
		match = gipv4.Validate(valueStr) || gipv6.Validate(valueStr)

	// IPv4.
	case "ipv4":
		match = gipv4.Validate(valueStr)

	// IPv6.
	case "ipv6":
		match = gipv6.Validate(valueStr)

	// MAC.
	case "mac":
		match = gregex.IsMatchString(`^([0-9A-Fa-f]{2}[\-:]){5}[0-9A-Fa-f]{2}$`, valueStr)

	default:
		return match, gerror.New("Invalid rule name: " + ruleKey)
	}
	return match, nil
}
