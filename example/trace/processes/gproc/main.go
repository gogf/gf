// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package main

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gproc"
)

func main() {
	ctx := gctx.GetInitCtx()
	g.Log().Debug(ctx, `this is main process`)
	if err := gproc.ShellRun(ctx, `go run sub/sub.go`); err != nil {
		panic(err)
	}
}
