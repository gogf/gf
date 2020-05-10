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
	"passport":             "The :attribute value is not a valid passport format",
	"password":             "The :attribute value is not a valid passport format",
	"password2":            "The :attribute value is not a valid passport format",
	"password3":            "The :attribute value is not a valid passport format",
	"postcode":             "The :attribute value is not a valid passport format",
	"resident-id":          "The :attribute value is not a valid resident id number",
	"bank-card":            "The :attribute must be a valid bank card number",
	"qq":                   "The :attribute must be a valid QQ number",
	"ip":                   "The :attribute must be a valid IP address",
	"ipv4":                 "The :attribute must be a valid IPv4 address",
	"ipv6":                 "The :attribute must be a valid IPv6 address",
	"mac":                  "The :attribute must be a valid MAC address",
	"url":                  "The :attribute must be a valid URL address",
	"domain":               "The :attribute must be a valid domain format",
	"length":               "The :attribute length must be between :min and :max",
	"min-length":           "The :attribute length must be equal or greater than :min",
	"max-length":           "The :attribute length must be equal or lesser than :max",
	"between":              "The :attribute value must be between :min and :max",
	"min":                  "The :attribute value must be equal or greater than :min",
	"max":                  "The :attribute value must be equal or lesser than :max",
	"json":                 "The :attribute must be a valid JSON string",
	"xml":                  "The :attribute must be a valid XML string",
	"array":                "The :attribute must be an array",
	"integer":              "The :attribute must be an integer",
	"float":                "The :attribute must be a float",
	"boolean":              "The :attribute field must be true or false",
	"same":                 "The :attribute value must be the same as field :other",
	"different":            "The :attribute value must be different from field :other",
	"in":                   "The :attribute value is not in acceptable range",
	"not-in":               "The :attribute value is not in acceptable range",
	"regex":                "The :attribute value is invalid",
}

// getDefaultErrorMessageByRule retrieves and returns the default error message
// for specified rule. It firstly retrieves the message from i18n manager, it returns
// from default error messages if it's not found in i18n manager.
func getDefaultErrorMessageByRule(rule string) string {
	i18nKey := fmt.Sprintf(`gf.gvalid.%s`, rule)
	content := gi18n.GetContent(i18nKey)
	if content == "" {
		content = defaultMessages[rule]
	}
	return content
}
