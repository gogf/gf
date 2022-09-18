// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package builtin

import (
	"errors"
	"strconv"

	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gutil"
)

// RuleGTE implements `gte` rule:
// Greater than or equal to `field`.
// It supports both integer and float.
//
// Format: gte:field
type RuleGTE struct{}

func init() {
	Register(RuleGTE{})
}

func (r RuleGTE) Name() string {
	return "gte"
}

func (r RuleGTE) Message() string {
	return "The {field} value `{value}` must be greater than or equal to field {pattern}"
}

func (r RuleGTE) Run(in RunInput) error {
	var (
		_, fieldValue     = gutil.MapPossibleItemByKey(in.Data.Map(), in.RulePattern)
		fieldValueN, err1 = strconv.ParseFloat(gconv.String(fieldValue), 10)
		valueN, err2      = strconv.ParseFloat(in.Value.String(), 10)
	)

	if valueN < fieldValueN || err1 != nil || err2 != nil {
		return errors.New(in.Message)
	}
	return nil
}
