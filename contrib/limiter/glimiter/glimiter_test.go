// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package glimiter_test

import (
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"

	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/os/gctx"
)

var (
	ctx    = gctx.GetInitCtx()
	config = &gredis.Config{
		Address: `:6379`,
		Db:      1,
	}
	re *gredis.Redis
)

func init() {
	r, err := gredis.New(config)
	if err != nil {
		panic(err)
	}
	re = r
}
