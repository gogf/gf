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
	Register(&RuleIn{})
}

func (r *RuleIn) Name() string {
	return "in"
}

func (r *RuleIn) Message() string {
	return "The {attribute} value `{value}` is not in acceptable range: {pattern}"
}

func (r *RuleIn) Run(in RunInput) error {
	var ok bool
	for _, rulePattern := range gstr.SplitAndTrim(in.RulePattern, ",") {
		if in.CaseInsensitive {
			ok = strings.EqualFold(in.Value.String(), strings.TrimSpace(rulePattern))
		} else {
			ok = strings.Compare(in.Value.String(), strings.TrimSpace(rulePattern)) == 0
		}
		if ok {
			return nil
		}
	}
	return errors.New(in.Message)
}
