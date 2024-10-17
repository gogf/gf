// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package redis

import (
	"github.com/gogf/gf/v2/frame/g"
)

type redisOption interface {
	OptionToArgs() []interface{}
}

func mustMergeOptionToArgs(args []interface{}, opt redisOption) []interface{} {
	if g.IsNil(opt) {
		return args
	}
	return append(args, opt.OptionToArgs()...)
}
