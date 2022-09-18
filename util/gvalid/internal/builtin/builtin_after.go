// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package builtin

import (
	"errors"

	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gutil"
)

// RuleAfter implements `after` rule:
// The datetime value should be after the value of field `field`.
//
// Format: after:field
type RuleAfter struct{}

func init() {
	Register(RuleAfter{})
}

func (r RuleAfter) Name() string {
	return "after"
}

func (r RuleAfter) Message() string {
	return "The {field} value `{value}` must be after field {pattern}"
}

func (r RuleAfter) Run(in RunInput) error {
	var (
		_, fieldValue = gutil.MapPossibleItemByKey(in.Data.Map(), in.RulePattern)
		valueDatetime = in.Value.Time()
		fieldDatetime = gconv.Time(fieldValue)
	)
	if valueDatetime.IsZero() || fieldDatetime.IsZero() {
		return errors.New(in.Message)
	}
	if valueDatetime.After(fieldDatetime) {
		return nil
	}
	return errors.New(in.Message)
}
