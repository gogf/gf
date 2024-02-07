// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package main

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/os/gcron"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
)

func main() {
	fmt.Println("start:", gtime.Now())
	var (
		err      error
		pattern1 = "# * * * * *"
		pattern2 = "# */2 * * * *"
	)
	_, err = gcron.Add(gctx.New(), pattern1, func(ctx context.Context) {
		fmt.Println(pattern1, gtime.Now())
	})
	if err != nil {
		panic(err)
	}
	_, err = gcron.Add(gctx.New(), pattern2, func(ctx context.Context) {
		fmt.Println(pattern2, gtime.Now())
	})
	if err != nil {
		panic(err)
	}

	select {}
}
