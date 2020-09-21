// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gvalid

import (
	"fmt"
	"github.com/gogf/gf/i18n/gi18n"
)

// defaultMessages is the default error messages.
// Note that these messages are synchronized from ./i18n/en/validation.toml .
var defaultMessages = map[string]string{
	"required":             "The :attribute field is required",
	"required-if":          "The :attribute field is required",
	"required-unless":      "The :attribute field is required",
	"required-with":        "The :attribute field is required",
	"required-with-all":    "The :attribute field is required",
	"required-without":     "The :attribute field is required",
	"required-without-all": "The :attribute field is required",
	"date":                 "The :attribute value is not a valid date",
	"date-format":          "The :attribute value does not match the format :format",
	"email":                "The :attribute value must be a valid email address",
	"phone":                "The :attribute value must be a valid phone number",
	"telephone":            "The :attribute value must be a valid telephone number",
	"passport":             "The :attribute value is not a valid passport format",
	"password":             "The :attribute value is not a valid passport format",
	"password2":            "The :attribute value is not a valid passport format",
	"password3":            "The :attribute value is not a valid passport format",
	"postcode":             "The :attribute value is not a valid passport format",
	"resident-id":          "The :attribute value is not a valid resident id number",
	"bank-card":            "The :attribute value must be a valid bank card number",
	"qq":                   "The :attribute value must be a valid QQ number",
	"ip":                   "The :attribute value must be a valid IP address",
	"ipv4":                 "The :attribute value must be a valid IPv4 address",
	"ipv6":                 "The :attribute value must be a valid IPv6 address",
	"mac":                  "The :attribute value must be a valid MAC address",
	"url":                  "The :attribute value must be a valid URL address",
	"domain":               "The :attribute value must be a valid domain format",
	"length":               "The :attribute value length must be between :min and :max",
	"min-length":           "The :attribute value length must be equal or greater than :min",
	"max-length":           "The :attribute value length must be equal or lesser than :max",
	"between":              "The :attribute value must be between :min and :max",
	"min":                  "The :attribute value must be equal or greater than :min",
	"max":                  "The :attribute value must be equal or lesser than :max",
	"json":                 "The :attribute value must be a valid JSON string",
	"xml":                  "The :attribute value must be a valid XML string",
	"array":                "The :attribute value must be an array",
	"integer":              "The :attribute value must be an integer",
	"float":                "The :attribute value must be a float",
	"boolean":              "The :attribute value field must be true or false",
	"same":                 "The :attribute value must be the same as field :field",
	"different":            "The :attribute value must be different from field :field",
	"in":                   "The :attribute value is not in acceptable range",
	"not-in":               "The :attribute value is not in acceptable range",
	"regex":                "The :attribute value is invalid",
	"__default__":          "The :attribute value is invalid",
}

// getErrorMessageByRule retrieves and returns the error message for specified rule.
// It firstly retrieves the message from custom message map, and then checks i18n manager,
// it returns the default error message if it's not found in custom message map or i18n manager.
func getErrorMessageByRule(ruleKey string, customMsgMap map[string]string) string {
	content := customMsgMap[ruleKey]
	if content != "" {
		return content
	}
	content = gi18n.GetContent(fmt.Sprintf(`gf.gvalid.rule.%s`, ruleKey))
	if content == "" {
		content = defaultMessages[ruleKey]
	}
	// If there's no configured rule message, it uses default one.
	if content == "" {
		content = gi18n.GetContent(`gf.gvalid.rule.__default__`)
		if content == "" {
			content = defaultMessages["__default__"]
		}
	}
	return content
}
