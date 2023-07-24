// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package builtin

import (
	"errors"

	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gutil"
)

// RuleBefore implements `before` rule:
// The datetime value should be after the value of field `field`.
//
// Format: before:field
type RuleBefore struct{}

func init() {
	Register(RuleBefore{})
}

func (r RuleBefore) Name() string {
	return "before"
}

func (r RuleBefore) Message() string {
	return "The {field} value `{value}` must be before field {field1} value `{value1}`"
}

func (r RuleBefore) Run(in RunInput) error {
	var (
		fieldName, fieldValue = gutil.MapPossibleItemByKey(in.Data.Map(), in.RulePattern)
		valueDatetime         = in.Value.Time()
		fieldDatetime         = gconv.Time(fieldValue)
	)
	if valueDatetime.Before(fieldDatetime) {
		return nil
	}
	return errors.New(gstr.ReplaceByMap(in.Message, map[string]string{
		"{field1}": fieldName,
		"{value1}": gconv.String(fieldValue),
	}))
}
