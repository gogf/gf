// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package builtin

import (
	"errors"

	"github.com/gogf/gf/v2/internal/json"
)

type RuleJson struct{}

func init() {
	Register(&RuleJson{})
}

func (r *RuleJson) Name() string {
	return "json"
}

func (r *RuleJson) Message() string {
	return "The {attribute} value `{value}` is not a valid JSON string"
}

func (r *RuleJson) Run(in RunInput) error {
	if json.Valid(in.Value.Bytes()) {
		return nil
	}
	return errors.New(in.Message)
}
