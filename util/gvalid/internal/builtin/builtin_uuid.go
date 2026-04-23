// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package builtin

import (
	"errors"

	"github.com/gogf/gf/v2/text/gregex"
)

// RuleUUID implements `uuid` rule:
// UUID.
//
// Format: uuid
type RuleUUID struct{}

func init() {
	Register(RuleUUID{})
}

func (r RuleUUID) Name() string {
	return "uuid"
}

func (r RuleUUID) Message() string {
	return "The {field} value `{value}` is not a valid UUID"
}

func (r RuleUUID) Run(in RunInput) error {
	ok := gregex.IsMatchString(
		`^[0-9A-Fa-f]{8}-[0-9A-Fa-f]{4}-[1-5][0-9A-Fa-f]{3}-[89AaBb][0-9A-Fa-f]{3}-[0-9A-Fa-f]{12}$`,
		in.Value.String(),
	)
	if ok {
		return nil
	}
	return errors.New(in.Message)
}
