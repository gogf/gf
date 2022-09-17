// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package builtin

import (
	"errors"

	"github.com/gogf/gf/v2/net/gipv4"
	"github.com/gogf/gf/v2/net/gipv6"
)

type RuleIp struct{}

func init() {
	Register(&RuleIp{})
}

func (r *RuleIp) Name() string {
	return "ip"
}

func (r *RuleIp) Message() string {
	return "The {attribute} value `{value}` is not a valid IP address"
}

func (r *RuleIp) Run(in RunInput) error {
	var (
		ok    bool
		value = in.Value.String()
	)
	if ok = gipv4.Validate(value) || gipv6.Validate(value); !ok {
		return errors.New(in.Message)
	}
	return nil
}
