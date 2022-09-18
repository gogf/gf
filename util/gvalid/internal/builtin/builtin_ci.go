// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package builtin

// RuleCi implements `ci` rule: Case-Insensitive configuration for those rules that need value comparison like:
// same, different, in, not-in, etc.
// Format:  ci
type RuleCi struct{}

func init() {
	Register(&RuleCi{})
}

func (r *RuleCi) Name() string {
	return "ci"
}

func (r *RuleCi) Message() string {
	return ""
}

func (r *RuleCi) Run(in RunInput) error {
	return nil
}
