// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package main

import (
	"fmt"
	"time"

	_ "github.com/gogf/gf/contrib/nosql/redis/v2"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

func main() {
	var ctx = gctx.New()
	err := g.Redis().SetEX(ctx, "key", "value", 1)
	if err != nil {
		g.Log().Fatal(ctx, err)
	}
	value, err := g.Redis().Get(ctx, "key")
	if err != nil {
		g.Log().Fatal(ctx, err)
	}
	fmt.Println(value.IsNil())
	fmt.Println(value.String())

	time.Sleep(time.Second)

	value, err = g.Redis().Get(ctx, "key")
	if err != nil {
		g.Log().Fatal(ctx, err)
	}
	fmt.Println(value.IsNil())
	fmt.Println(value.Val())
}
