// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package builtin

// RuleBail implements `bail` rule:
// Stop validating when this field's validation failed.
//
// Format: bail
type RuleBail struct{}

func init() {
	Register(RuleBail{})
}

func (r RuleBail) Name() string {
	return "bail"
}

func (r RuleBail) Message() string {
	return ""
}

func (r RuleBail) Run(in RunInput) error {
	return nil
}
