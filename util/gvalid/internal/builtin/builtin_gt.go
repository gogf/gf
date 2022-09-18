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

// RuleGT implements `gt` rule:
// Greater than `field`.
// It supports both integer and float.
//
// Format: gt:field
type RuleGT struct{}

func init() {
	Register(RuleGT{})
}

func (r RuleGT) Name() string {
	return "gt"
}

func (r RuleGT) Message() string {
	return "The {field} value `{value}` must be greater than field {pattern}"
}

func (r RuleGT) Run(in RunInput) error {
	var (
		_, fieldValue     = gutil.MapPossibleItemByKey(in.Data.Map(), in.RulePattern)
		fieldValueN, err1 = strconv.ParseFloat(gconv.String(fieldValue), 10)
		valueN, err2      = strconv.ParseFloat(in.Value.String(), 10)
	)

	if valueN <= fieldValueN || err1 != nil || err2 != nil {
		return errors.New(in.Message)
	}
	return nil
}
