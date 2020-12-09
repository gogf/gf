// Copyright 2017-2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gvalid

import (
	"errors"
	"github.com/gogf/gf/internal/json"
	"github.com/gogf/gf/net/gipv4"
	"github.com/gogf/gf/net/gipv6"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/text/gregex"
	"github.com/gogf/gf/util/gconv"
	"regexp"
	"strconv"
	"strings"
)

const (
	// regular expression pattern for single validation rule.
	gSINGLE_RULE_PATTERN = `^([\w-]+):{0,1}(.*)`
)

var (
	// regular expression object for single rule
	// which is compiled just once and of repeatable usage.
	ruleRegex, _ = regexp.Compile(gSINGLE_RULE_PATTERN)

	// mustCheckRulesEvenValueEmpty specifies some rules that must be validated
	// even the value is empty (nil or empty).
	mustCheckRulesEvenValueEmpty = map[string]struct{}{
		"required":             {},
		"required-if":          {},
		"required-unless":      {},
		"required-with":        {},
		"required-with-all":    {},
		"required-without":     {},
		"required-without-all": {},
		//"same":                 {},
		//"different":            {},
		//"in":                   {},
		//"not-in":               {},
		//"regex":                {},
	}
	// allSupportedRules defines all supported rules that is used for quick checks.
	allSupportedRules = map[string]struct{}{
		"required":             {},
		"required-if":          {},
		"required-unless":      {},
		"required-with":        {},
		"required-with-all":    {},
		"required-without":     {},
		"required-without-all": {},
		"date":                 {},
		"date-format":          {},
		"email":                {},
		"phone":                {},
		"phone-loose":          {},
		"telephone":            {},
		"passport":             {},
		"password":             {},
		"password2":            {},
		"password3":            {},
		"postcode":             {},
		"resident-id":          {},
		"bank-card":            {},
		"qq":                   {},
		"ip":                   {},
		"ipv4":                 {},
		"ipv6":                 {},
		"mac":                  {},
		"url":                  {},
		"domain":               {},
		"length":               {},
		"min-length":           {},
		"max-length":           {},
		"between":              {},
		"min":                  {},
		"max":                  {},
		"json":                 {},
		"integer":              {},
		"float":                {},
		"boolean":              {},
		"same":                 {},
		"different":            {},
		"in":                   {},
		"not-in":               {},
		"regex":                {},
	}
	// boolMap defines the boolean values.
	boolMap = map[string]struct{}{
		"1":     {},
		"true":  {},
		"on":    {},
		"yes":   {},
		"":      {},
		"0":     {},
		"false": {},
		"off":   {},
		"no":    {},
	}
)

// Check checks single value with specified rules.
// It returns nil if successful validation.
//
// The parameter <value> can be any type of variable, which will be converted to string
// for validation.
// The parameter <rules> can be one or more rules, multiple rules joined using char '|'.
// The parameter <messages> specifies the custom error messages, which can be type of:
// string/map/struct/*struct.
// The optional parameter <params> specifies the extra validation parameters for some rules
// like: required-*、same、different, etc.
func Check(value interface{}, rules string, messages interface{}, params ...interface{}) *Error {
	return doCheck("", value, rules, messages, params...)
}

// doCheck does the really rules validation for single key-value.
func doCheck(key string, value interface{}, rules string, messages interface{}, params ...interface{}) *Error {
	// If there's no validation rules, it does nothing and returns quickly.
	if rules == "" {
		return nil
	}
	// It converts value to string and then does the validation.
	var (
		// Do not trim it as the space is also part of the value.
		data          = make(map[string]string)
		errorMsgArray = make(map[string]string)
	)
	if len(params) > 0 {
		for k, v := range gconv.Map(params[0]) {
			data[k] = gconv.String(v)
		}
	}
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
		if !ok && customRuleFuncMap[array[0]] == nil {
			if i > 0 && ruleItems[i-1][:5] == "regex" {
				ruleItems[i-1] += "|" + ruleItems[i]
				ruleItems = append(ruleItems[:i], ruleItems[i+1:]...)
			} else {
				return newErrorStr(
					"invalid_rules",
					"invalid rules: "+rules,
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
			err         error
			match       = false
			results     = ruleRegex.FindStringSubmatch(ruleItems[index])
			ruleKey     = strings.TrimSpace(results[1])
			rulePattern = strings.TrimSpace(results[2])
		)
		if len(msgArray) > index {
			customMsgMap[ruleKey] = strings.TrimSpace(msgArray[index])
		}

		if f, ok := customRuleFuncMap[ruleKey]; ok {
			// It checks custom validation rules with most priority.
			var (
				dataMap map[string]interface{}
				message = getErrorMessageByRule(ruleKey, customMsgMap)
			)
			if len(params) > 0 {
				dataMap = gconv.Map(params[0])
			}
			if err := f(ruleItems[index], value, message, dataMap); err != nil {
				match = false
				errorMsgArray[ruleKey] = err.Error()
			} else {
				match = true
			}
		} else {
			// It checks build-in validation rules if there's no custom rule.
			match, err = doCheckBuildInRules(index, value, ruleKey, rulePattern, ruleItems, data, customMsgMap)
			if !match && err != nil {
				errorMsgArray[ruleKey] = err.Error()
			}
		}

		// Error message handling.
		if !match {
			// It does nothing if the error message for this rule
			// is already set in previous validation.
			if _, ok := errorMsgArray[ruleKey]; !ok {
				errorMsgArray[ruleKey] = getErrorMessageByRule(ruleKey, customMsgMap)
			}
		}
		index++
	}
	if len(errorMsgArray) > 0 {
		return newError([]string{rules}, ErrorMap{
			key: errorMsgArray,
		})
	}
	return nil
}

func doCheckBuildInRules(
	index int,
	value interface{},
	ruleKey string,
	rulePattern string,
	ruleItems []string,
	dataMap map[string]string,
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
		match = checkRequired(valueStr, ruleKey, rulePattern, dataMap)

	// Length rules.
	// It also supports length of unicode string.
	case
		"length",
		"min-length",
		"max-length":
		if msg := checkLength(valueStr, ruleKey, rulePattern, customMsgMap); msg != "" {
			return match, errors.New(msg)
		} else {
			match = true
		}

	// Range rules.
	case
		"min",
		"max",
		"between":
		if msg := checkRange(valueStr, ruleKey, rulePattern, customMsgMap); msg != "" {
			return match, errors.New(msg)
		} else {
			match = true
		}

	// Custom regular expression.
	case "regex":
		// It here should check the rule as there might be special char '|' in it.
		for i := index + 1; i < len(ruleItems); i++ {
			if !gregex.IsMatchString(gSINGLE_RULE_PATTERN, ruleItems[i]) {
				rulePattern += "|" + ruleItems[i]
				index++
			}
		}
		match = gregex.IsMatchString(rulePattern, valueStr)

	// Date rules.
	case "date":
		// Standard date string, which must contain char '-' or '.'.
		if _, err := gtime.StrToTime(valueStr); err == nil {
			match = true
			break
		}
		// Date that not contains char '-' or '.'.
		if _, err := gtime.StrToTime(valueStr, "Ymd"); err == nil {
			match = true
			break
		}

	// Date rule with specified format.
	case "date-format":
		if _, err := gtime.StrToTimeFormat(valueStr, rulePattern); err == nil {
			match = true
		} else {
			var msg string
			msg = getErrorMessageByRule(ruleKey, customMsgMap)
			msg = strings.Replace(msg, ":format", rulePattern, -1)
			return match, errors.New(msg)
		}

	// Values of two fields should be equal as string.
	case "same":
		if v, ok := dataMap[rulePattern]; ok {
			if strings.Compare(valueStr, v) == 0 {
				match = true
			}
		}
		if !match {
			var msg string
			msg = getErrorMessageByRule(ruleKey, customMsgMap)
			msg = strings.Replace(msg, ":field", rulePattern, -1)
			return match, errors.New(msg)
		}

	// Values of two fields should not be equal as string.
	case "different":
		match = true
		if v, ok := dataMap[rulePattern]; ok {
			if strings.Compare(valueStr, v) == 0 {
				match = false
			}
		}
		if !match {
			var msg string
			msg = getErrorMessageByRule(ruleKey, customMsgMap)
			msg = strings.Replace(msg, ":field", rulePattern, -1)
			return match, errors.New(msg)
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
		match = checkResidentId(valueStr)

	// Bank card number using LUHN algorithm.
	case "bank-card":
		match = checkLuHn(valueStr)

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
		return match, errors.New("Invalid rule name: " + ruleKey)
	}
	return match, nil
}
