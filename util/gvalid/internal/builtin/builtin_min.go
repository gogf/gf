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
)

type RuleMin struct{}

func init() {
	Register(&RuleMin{})
}

func (r *RuleMin) Name() string {
	return "min"
}

func (r *RuleMin) Message() string {
	return "The {attribute} value `{value}` must be equal or greater than {min}"
}

func (r *RuleMin) Run(in RunInput) error {
	var (
		min, err1    = strconv.ParseFloat(in.RulePattern, 10)
		valueN, err2 = strconv.ParseFloat(in.Value.String(), 10)
	)
	if valueN < min || err1 != nil || err2 != nil {
		return errors.New(gstr.Replace(in.Message, "{min}", strconv.FormatFloat(min, 'f', -1, 64)))
	}
	return nil
}
