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

// RuleMaxLength implements `max-length` rule:
// Length is equal or lesser than :max.
// The length is calculated using unicode string, which means one chinese character or letter both has the length of 1.
//
// Format: max-length:max
type RuleMaxLength struct{}

func init() {
	Register(&RuleMaxLength{})
}

func (r *RuleMaxLength) Name() string {
	return "max-length"
}

func (r *RuleMaxLength) Message() string {
	return "The {attribute} value `{value}` length must be equal or lesser than {max}"
}

func (r *RuleMaxLength) Run(in RunInput) error {
	var (
		valueRunes = gconv.Runes(in.Value.String())
		valueLen   = len(valueRunes)
	)
	max, err := strconv.Atoi(in.RulePattern)
	if valueLen > max || err != nil {
		return errors.New(gstr.Replace(in.Message, "{max}", strconv.Itoa(max)))
	}
	return nil
}
