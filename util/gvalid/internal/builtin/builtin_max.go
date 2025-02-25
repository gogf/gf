// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package builtin

import (
	"errors"
	"strconv"

	"github.com/gogf/gf/v2/text/gstr"
)

// RuleMax implements `max` rule:
// Equal or lesser than :max. It supports both integer and float.
//
// Format: max:max
type RuleMax struct{}

func init() {
	Register(RuleMax{})
}

func (r RuleMax) Name() string {
	return "max"
}

func (r RuleMax) Message() string {
	return "The {field} value `{value}` must be equal or lesser than {max}"
}

func (r RuleMax) Run(in RunInput) error {
	var (
		max, err1    = strconv.ParseFloat(in.RulePattern, 64)
		valueN, err2 = strconv.ParseFloat(in.Value.String(), 64)
	)
	if valueN > max || err1 != nil || err2 != nil {
		return errors.New(gstr.Replace(in.Message, "{max}", strconv.FormatFloat(max, 'f', -1, 64)))
	}
	return nil
}
