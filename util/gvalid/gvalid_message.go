// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gvalid

// defaultMessages is the default error messages.
var defaultMessages = map[string]string{
	"required":             "The :attribute field is required",
	"required-if":          "The :attribute field is required when :other is :value",
	"required-unless":      "The :attribute field is required unless :other is in :values",
	"required-with":        "The :attribute field is required when :values is present",
	"required-with-all":    "The :attribute field is required when :values is present",
	"required-without":     "The :attribute field is required when :values is not present",
	"required-without-all": "The :attribute field is required when none of :values are present",
	"date":                 "The :attribute is not a valid date",
	"date-format":          "The :attribute does not match the format :format",
	"email":                "The :attribute must be a valid email address",
	"phone":                "The :attribute must be a valid phone number",
	"telephone":            "The :attribute must be a valid telephone number",
	"passport":             "Invalid passport format",
	"password":             "Invalid passport format",
	"password2":            "Invalid passport format",
	"password3":            "Invalid passport format",
	"postcode":             "Invalid postcode format",
	"id-number":            "Invalid id",
	"luhn":                 "The :attribute must be a valid bank card number",
	"qq":                   "The :attribute must be a valid QQ number",
	"ip":                   "The :attribute must be a valid IP address",
	"ipv4":                 "The :attribute must be a valid IPv4 address",
	"ipv6":                 "The :attribute must be a valid IPv6 address",
	"mac":                  "MAC地址格式不正确",
	"url":                  "URL地址格式不正确",
	"domain":               "域名格式不正确",
	"length":               "字段长度为:min到:max个字符",
	"min-length":           "字段最小长度为:min",
	"max-length":           "字段最大长度为:max",
	"between":              "字段大小为:min到:max",
	"min":                  "字段最小值为:min",
	"max":                  "字段最大值为:max",
	"json":                 "字段应当为JSON格式",
	"xml":                  "字段应当为XML格式",
	"array":                "字段应当为数组",
	"integer":              "字段应当为整数",
	"float":                "字段应当为浮点数",
	"boolean":              "字段应当为布尔值",
	"same":                 "字段值不合法",
	"different":            "字段值不合法",
	"in":                   "字段值不合法",
	"not-in":               "字段值不合法",
	"regex":                "字段值不合法",
}

func init() {
	errorMsgMap.Sets(defaultMessages)
}

// SetDefaultErrorMsgs sets the default error messages for package.
func SetDefaultErrorMsgs(msgs map[string]string) {
	errorMsgMap.Sets(msgs)
}
