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
	"github.com/gogf/gf/v2/util/gutil"
)

// RuleLTE implements `lte` rule:
// Lesser than or equal to `field`.
// It supports both integer and float.
//
// Format: lte:field
type RuleLTE struct{}

func init() {
	Register(RuleLTE{})
}

func (r RuleLTE) Name() string {
	return "lte"
}

func (r RuleLTE) Message() string {
	return "The {field} value `{value}` must be lesser than or equal to field {field1} value `{value1}`"
}

func (r RuleLTE) Run(in RunInput) error {
	var (
		fieldName, fieldValue = gutil.MapPossibleItemByKey(in.Data.Map(), in.RulePattern)
		fieldValueN, err1     = strconv.ParseFloat(gconv.String(fieldValue), 64)
		valueN, err2          = strconv.ParseFloat(in.Value.String(), 64)
	)

	if valueN > fieldValueN || err1 != nil || err2 != nil {
		return errors.New(gstr.ReplaceByMap(in.Message, map[string]string{
			"{field1}": fieldName,
			"{value1}": gconv.String(fieldValue),
		}))
	}
	return nil
}
