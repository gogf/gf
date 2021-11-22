// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gvalid

import (
	"context"
	"errors"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/internal/utils"
	"github.com/gogf/gf/v2/net/gipv4"
	"github.com/gogf/gf/v2/net/gipv6"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gutil"
)

type iTime interface {
	Date() (year int, month time.Month, day int)
	IsZero() bool
}

// CheckValue checks single value with specified rules.
// It returns nil if successful validation.
func (v *Validator) CheckValue(ctx context.Context, value interface{}) Error {
	return v.doCheckValue(ctx, doCheckValueInput{
		Name:     "",
		Value:    value,
		Rule:     gconv.String(v.rules),
		Messages: v.messages,
		DataRaw:  v.data,
		DataMap:  gconv.Map(v.data),
	})
}

type doCheckValueInput struct {
	Name     string                 // Name specifies the name of parameter `value`.
	Value    interface{}            // Value specifies the value for the rules to be validated.
	Rule     string                 // Rule specifies the validation rules string, like "required", "required|between:1,100", etc.
	Messages interface{}            // Messages specifies the custom error messages for this rule from parameters input, which is usually type of map/slice.
	DataRaw  interface{}            // DataRaw specifies the `raw data` which is passed to the Validator. It might be type of map/struct or a nil value.
	DataMap  map[string]interface{} // DataMap specifies the map that is converted from `dataRaw`. It is usually used internally
}

// doCheckSingleValue does the really rules validation for single key-value.
func (v *Validator) doCheckValue(ctx context.Context, in doCheckValueInput) Error {
	// If there's no validation rules, it does nothing and returns quickly.
	if in.Rule == "" {
		return nil
	}
	// It converts value to string and then does the validation.
	var (
		// Do not trim it as the space is also part of the value.
		ruleErrorMap = make(map[string]error)
	)
	// Custom error messages handling.
	var (
		msgArray     = make([]string, 0)
		customMsgMap = make(map[string]string)
	)
	switch messages := in.Messages.(type) {
	case string:
		msgArray = strings.Split(messages, "|")

	default:
		for k, message := range gconv.Map(in.Messages) {
			customMsgMap[k] = gconv.String(message)
		}
	}
	// Handle the char '|' in the rule,
	// which makes this rule separated into multiple rules.
	ruleItems := strings.Split(strings.TrimSpace(in.Rule), "|")
	for i := 0; ; {
		array := strings.Split(ruleItems[i], ":")
		_, ok := allSupportedRules[array[0]]
		if !ok && v.getRuleFunc(array[0]) == nil {
			if i > 0 && ruleItems[i-1][:5] == "regex" {
				ruleItems[i-1] += "|" + ruleItems[i]
				ruleItems = append(ruleItems[:i], ruleItems[i+1:]...)
			} else {
				return newValidationErrorByStr(
					internalRulesErrRuleName,
					errors.New(internalRulesErrRuleName+": "+in.Rule),
				)
			}
		} else {
			i++
		}
		if i == len(ruleItems) {
			break
		}
	}
	var (
		hasBailRule        = v.bail
		hasCaseInsensitive = v.caseInsensitive
	)
	for index := 0; index < len(ruleItems); {
		var (
			err            error
			match          = false                                          // whether this rule is matched(has no error)
			results        = ruleRegex.FindStringSubmatch(ruleItems[index]) // split single rule.
			ruleKey        = gstr.Trim(results[1])                          // rule key like "max" in rule "max: 6"
			rulePattern    = gstr.Trim(results[2])                          // rule pattern is like "6" in rule:"max:6"
			customRuleFunc RuleFunc
		)

		if !hasBailRule && ruleKey == ruleNameBail {
			hasBailRule = true
		}

		if !hasCaseInsensitive && ruleKey == ruleNameCi {
			hasCaseInsensitive = true
		}

		// Ignore logic executing for marked rules.
		if markedRuleMap[ruleKey] {
			index++
			continue
		}

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
			message := v.getErrorMessageByRule(ctx, ruleKey, customMsgMap)
			if err = customRuleFunc(ctx, RuleFuncInput{
				Rule:    ruleItems[index],
				Message: message,
				Value:   gvar.New(in.Value),
				Data:    gvar.New(in.DataRaw),
			}); err != nil {
				match = false
				// The error should have stack info to indicate the error position.
				if !gerror.HasStack(err) {
					err = gerror.NewCodeSkip(gcode.CodeValidationFailed, 1, err.Error())
				}
				// The error should have error code that is `gcode.CodeValidationFailed`.
				if gerror.Code(err) == gcode.CodeNil {
					if e, ok := err.(*gerror.Error); ok {
						e.SetCode(gcode.CodeValidationFailed)
					}
				}
				ruleErrorMap[ruleKey] = err
			} else {
				match = true
			}
		} else {
			// It checks build-in validation rules if there's no custom rule.
			match, err = v.doCheckSingleBuildInRules(
				ctx,
				doCheckBuildInRulesInput{
					Index:           index,
					Value:           in.Value,
					RuleKey:         ruleKey,
					RulePattern:     rulePattern,
					RuleItems:       ruleItems,
					DataMap:         in.DataMap,
					CustomMsgMap:    customMsgMap,
					CaseInsensitive: hasCaseInsensitive,
				},
			)
			if !match && err != nil {
				ruleErrorMap[ruleKey] = err
			}
		}

		// Error message handling.
		if !match {
			// It does nothing if the error message for this rule
			// is already set in previous validation.
			if _, ok := ruleErrorMap[ruleKey]; !ok {
				ruleErrorMap[ruleKey] = errors.New(v.getErrorMessageByRule(ctx, ruleKey, customMsgMap))
			}

			// Error variable replacement for error message.
			if err = ruleErrorMap[ruleKey]; !gerror.HasStack(err) {
				var s string
				s = gstr.ReplaceByMap(err.Error(), map[string]string{
					"{value}":     gconv.String(in.Value),
					"{pattern}":   rulePattern,
					"{attribute}": in.Name,
				})
				s, _ = gregex.ReplaceString(`\s{2,}`, ` `, s)
				ruleErrorMap[ruleKey] = errors.New(s)
			}

			// If it is with error and there's bail rule,
			// it then does not continue validating for left rules.
			if hasBailRule {
				break
			}
		}
		index++
	}
	if len(ruleErrorMap) > 0 {
		return newValidationError(
			gcode.CodeValidationFailed,
			[]fieldRule{{Name: in.Name, Rule: in.Rule}},
			map[string]map[string]error{
				in.Name: ruleErrorMap,
			},
		)
	}
	return nil
}

type doCheckBuildInRulesInput struct {
	Index           int                    // Index of RuleKey in RuleItems.
	Value           interface{}            // Value to be validated.
	RuleKey         string                 // RuleKey is like the "max" in rule "max: 6"
	RulePattern     string                 // RulePattern is like "6" in rule:"max:6"
	RuleItems       []string               // RuleItems are all the rules that should be validated on single field, like: []string{"required", "min:1"}
	DataMap         map[string]interface{} // Parameter map.
	CustomMsgMap    map[string]string      // Custom error message map.
	CaseInsensitive bool                   // Case-Insensitive comparison.
}

func (v *Validator) doCheckSingleBuildInRules(ctx context.Context, in doCheckBuildInRulesInput) (match bool, err error) {
	valueStr := gconv.String(in.Value)
	switch in.RuleKey {
	// Required rules.
	case
		"required",
		"required-if",
		"required-unless",
		"required-with",
		"required-with-all",
		"required-without",
		"required-without-all":
		match = v.checkRequired(checkRequiredInput{
			Value:           in.Value,
			RuleKey:         in.RuleKey,
			RulePattern:     in.RulePattern,
			DataMap:         in.DataMap,
			CaseInsensitive: in.CaseInsensitive,
		})

	// Length rules.
	// It also supports length of unicode string.
	case
		"length",
		"min-length",
		"max-length",
		"size":
		if msg := v.checkLength(ctx, valueStr, in.RuleKey, in.RulePattern, in.CustomMsgMap); msg != "" {
			return match, errors.New(msg)
		} else {
			match = true
		}

	// Range rules.
	case
		"min",
		"max",
		"between":
		if msg := v.checkRange(ctx, valueStr, in.RuleKey, in.RulePattern, in.CustomMsgMap); msg != "" {
			return match, errors.New(msg)
		} else {
			match = true
		}

	// Custom regular expression.
	case "regex":
		// It here should check the rule as there might be special char '|' in it.
		for i := in.Index + 1; i < len(in.RuleItems); i++ {
			if !gregex.IsMatchString(singleRulePattern, in.RuleItems[i]) {
				in.RulePattern += "|" + in.RuleItems[i]
				in.Index++
			}
		}
		match = gregex.IsMatchString(in.RulePattern, valueStr)

	// Date rules.
	case "date":
		// support for time value, eg: gtime.Time/*gtime.Time, time.Time/*time.Time.
		if value, ok := in.Value.(iTime); ok {
			return !value.IsZero(), nil
		}
		match = gregex.IsMatchString(`\d{4}[\.\-\_/]{0,1}\d{2}[\.\-\_/]{0,1}\d{2}`, valueStr)

	// Datetime rule.
	case "datetime":
		// support for time value, eg: gtime.Time/*gtime.Time, time.Time/*time.Time.
		if value, ok := in.Value.(iTime); ok {
			return !value.IsZero(), nil
		}
		if _, err = gtime.StrToTimeFormat(valueStr, `Y-m-d H:i:s`); err == nil {
			match = true
		}

	// Date rule with specified format.
	case "date-format":
		// support for time value, eg: gtime.Time/*gtime.Time, time.Time/*time.Time.
		if value, ok := in.Value.(iTime); ok {
			return !value.IsZero(), nil
		}
		if _, err = gtime.StrToTimeFormat(valueStr, in.RulePattern); err == nil {
			match = true
		} else {
			var (
				msg string
			)
			msg = v.getErrorMessageByRule(ctx, in.RuleKey, in.CustomMsgMap)
			return match, errors.New(msg)
		}

	// Values of two fields should be equal as string.
	case "same":
		_, foundValue := gutil.MapPossibleItemByKey(in.DataMap, in.RulePattern)
		if foundValue != nil {
			if in.CaseInsensitive {
				match = strings.EqualFold(valueStr, gconv.String(foundValue))
			} else {
				match = strings.Compare(valueStr, gconv.String(foundValue)) == 0
			}
		}
		if !match {
			var msg string
			msg = v.getErrorMessageByRule(ctx, in.RuleKey, in.CustomMsgMap)
			return match, errors.New(msg)
		}

	// Values of two fields should not be equal as string.
	case "different":
		match = true
		_, foundValue := gutil.MapPossibleItemByKey(in.DataMap, in.RulePattern)
		if foundValue != nil {
			if in.CaseInsensitive {
				match = !strings.EqualFold(valueStr, gconv.String(foundValue))
			} else {
				match = strings.Compare(valueStr, gconv.String(foundValue)) != 0
			}
		}
		if !match {
			var msg string
			msg = v.getErrorMessageByRule(ctx, in.RuleKey, in.CustomMsgMap)
			return match, errors.New(msg)
		}

	// Field value should be in range of.
	case "in":
		for _, value := range gstr.SplitAndTrim(in.RulePattern, ",") {
			if in.CaseInsensitive {
				match = strings.EqualFold(valueStr, strings.TrimSpace(value))
			} else {
				match = strings.Compare(valueStr, strings.TrimSpace(value)) == 0
			}
			if match {
				break
			}
		}

	// Field value should not be in range of.
	case "not-in":
		match = true
		for _, value := range gstr.SplitAndTrim(in.RulePattern, ",") {
			if in.CaseInsensitive {
				match = !strings.EqualFold(valueStr, strings.TrimSpace(value))
			} else {
				match = strings.Compare(valueStr, strings.TrimSpace(value)) != 0
			}
			if !match {
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
		if _, err = strconv.Atoi(valueStr); err == nil {
			match = true
		}

	// Float.
	case "float":
		if _, err = strconv.ParseFloat(valueStr, 10); err == nil {
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
		return match, errors.New("Invalid rule name: " + in.RuleKey)
	}
	return match, nil
}

type doCheckValueRecursivelyInput struct {
	Value               interface{}
	Type                reflect.Type
	OriginKind          reflect.Kind
	ErrorMaps           map[string]map[string]error
	ResultSequenceRules *[]fieldRule
}

func (v *Validator) doCheckValueRecursively(ctx context.Context, in doCheckValueRecursivelyInput) {
	switch in.OriginKind {
	case reflect.Struct:
		// Ignore data, rules and messages from parent.
		validator := v.Clone()
		validator.rules = nil
		validator.messages = nil
		if err := validator.Data(in.Value).doCheckStruct(ctx, reflect.New(in.Type).Interface()); err != nil {
			// It merges the errors into single error map.
			for k, m := range err.(*validationError).errors {
				in.ErrorMaps[k] = m
			}
			if in.ResultSequenceRules != nil {
				*in.ResultSequenceRules = append(*in.ResultSequenceRules, err.(*validationError).rules...)
			}
		}

	case reflect.Map:
		var (
			dataMap   = gconv.Map(in.Value)
			validator = v.Clone()
		)
		// Ignore data, rules and messages from parent.
		validator.rules = nil
		validator.messages = nil
		for _, item := range dataMap {
			originTypeAndKind := utils.OriginTypeAndKind(item)
			v.doCheckValueRecursively(ctx, doCheckValueRecursivelyInput{
				Value:               item,
				Type:                originTypeAndKind.InputType,
				OriginKind:          originTypeAndKind.OriginKind,
				ErrorMaps:           in.ErrorMaps,
				ResultSequenceRules: in.ResultSequenceRules,
			})
			// Bail feature.
			if v.bail && len(in.ErrorMaps) > 0 {
				break
			}
		}

	case reflect.Slice, reflect.Array:
		array := gconv.Interfaces(in.Value)
		if len(array) == 0 {
			return
		}
		for _, item := range array {
			originTypeAndKind := utils.OriginTypeAndKind(item)
			v.doCheckValueRecursively(ctx, doCheckValueRecursivelyInput{
				Value:               item,
				Type:                originTypeAndKind.InputType,
				OriginKind:          originTypeAndKind.OriginKind,
				ErrorMaps:           in.ErrorMaps,
				ResultSequenceRules: in.ResultSequenceRules,
			})
			// Bail feature.
			if v.bail && len(in.ErrorMaps) > 0 {
				break
			}
		}
	}
}
