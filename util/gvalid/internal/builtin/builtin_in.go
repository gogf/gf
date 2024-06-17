// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package builtin

import (
	"errors"
	"strings"

	"github.com/gogf/gf/v2/text/gstr"
)

// RuleIn implements `in` rule:
// Value should be in: value1,value2,...
//
// Format: in:value1,value2,...
type RuleIn struct{}

func init() {
	Register(RuleIn{})
}

func (r RuleIn) Name() string {
	return "in"
}

func (r RuleIn) Message() string {
	return "The {field} value `{value}` is not in acceptable range: {pattern}"
}

func (r RuleIn) Run(in RunInput) error {
	var (
		ok               bool
		inputValueString = in.Value.String()
	)

	for _, rulePattern := range gstr.SplitAndTrim(in.RulePattern, ",") {
		if in.Option.CaseInsensitive {
			ok = strings.EqualFold(inputValueString, strings.TrimSpace(rulePattern))
		} else {
			ok = strings.Compare(inputValueString, strings.TrimSpace(rulePattern)) == 0
		}
		if ok {
			return nil
		}
	}
	return errors.New(in.Message)
}
