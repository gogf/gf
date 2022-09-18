// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package builtin

import (
	"errors"
	"time"

	"github.com/gogf/gf/v2/text/gregex"
)

type RuleDate struct{}

func init() {
	Register(&RuleDate{})
}

func (r *RuleDate) Name() string {
	return "date"
}

func (r *RuleDate) Message() string {
	return "The {attribute} value `{value}` is not a valid date"
}

func (r *RuleDate) Run(in RunInput) error {
	type iTime interface {
		Date() (year int, month time.Month, day int)
		IsZero() bool
	}
	// support for time value, eg: gtime.Time/*gtime.Time, time.Time/*time.Time.
	if obj, ok := in.Value.Val().(iTime); ok {
		if obj.IsZero() {
			return errors.New(in.Message)
		}
	}
	if !gregex.IsMatchString(
		`\d{4}[\.\-\_/]{0,1}\d{2}[\.\-\_/]{0,1}\d{2}`,
		in.Value.String(),
	) {
		return errors.New(in.Message)
	}
	return nil
}
