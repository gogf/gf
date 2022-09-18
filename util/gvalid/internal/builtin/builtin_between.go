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

	"github.com/gogf/gf/v2/text/gstr"
)

// RuleBetween implements `between` rule:
// Range between :min and :max. It supports both integer and float.
//
// Format: between:min,max
type RuleBetween struct{}

func init() {
	Register(&RuleBetween{})
}

func (r *RuleBetween) Name() string {
	return "between"
}

func (r *RuleBetween) Message() string {
	return "The {attribute} value `{value}` must be between {min} and {max}"
}

func (r *RuleBetween) Run(in RunInput) error {
	var (
		array = strings.Split(in.RulePattern, ",")
		min   = float64(0)
		max   = float64(0)
	)
	if len(array) > 0 {
		if v, err := strconv.ParseFloat(strings.TrimSpace(array[0]), 10); err == nil {
			min = v
		}
	}
	if len(array) > 1 {
		if v, err := strconv.ParseFloat(strings.TrimSpace(array[1]), 10); err == nil {
			max = v
		}
	}
	valueF, err := strconv.ParseFloat(in.Value.String(), 10)
	if valueF < min || valueF > max || err != nil {
		return errors.New(gstr.ReplaceByMap(in.Message, map[string]string{
			"{min}": strconv.FormatFloat(min, 'f', -1, 64),
			"{max}": strconv.FormatFloat(max, 'f', -1, 64),
		}))
	}
	return nil
}
