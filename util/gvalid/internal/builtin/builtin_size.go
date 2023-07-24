// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package builtin

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/util/gconv"
)

// RuleSize implements `size` rule:
// Length must be :size.
// The length is calculated using unicode string, which means one chinese character or letter both has the length of 1.
//
// Format: size:size
type RuleSize struct{}

func init() {
	Register(RuleSize{})
}

func (r RuleSize) Name() string {
	return "size"
}

func (r RuleSize) Message() string {
	return "The {field} value `{value}` length must be {size}"
}

func (r RuleSize) Run(in RunInput) error {
	var (
		valueRunes = gconv.Runes(in.Value.String())
		valueLen   = len(valueRunes)
	)
	size, err := strconv.Atoi(in.RulePattern)
	if valueLen != size || err != nil {
		return errors.New(strings.Replace(in.Message, "{size}", strconv.Itoa(size), -1))
	}
	return nil
}
