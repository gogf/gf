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
	"github.com/gogf/gf/v2/util/gconv"
)

// RuleMinLength implements `min-length` rule:
// Length is equal or greater than :min.
// The length is calculated using unicode string, which means one chinese character or letter both has the length of 1.
//
// Format: min-length:min
type RuleMinLength struct{}

func init() {
	Register(&RuleMinLength{})
}

func (r *RuleMinLength) Name() string {
	return "min-length"
}

func (r *RuleMinLength) Message() string {
	return "The {attribute} value `{value}` length must be equal or greater than {min}"
}

func (r *RuleMinLength) Run(in RunInput) error {
	var (
		valueRunes = gconv.Runes(in.Value.String())
		valueLen   = len(valueRunes)
	)
	min, err := strconv.Atoi(in.RulePattern)
	if valueLen < min || err != nil {
		return errors.New(gstr.Replace(in.Message, "{min}", strconv.Itoa(min)))
	}
	return nil
}
