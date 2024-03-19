// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/wangyougui/gf.

package main

import (
	"context"

	"github.com/wangyougui/gf/v2/frame/g"
	"github.com/wangyougui/gf/v2/os/gcmd"
	"github.com/wangyougui/gf/v2/os/gctx"
)

var (
	Sub = &gcmd.Command{
		Name:  "sub",
		Brief: "sub process",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			g.Log().Debug(ctx, `this is sub process`)
			return nil
		},
	}
)

func main() {
	Sub.Run(gctx.GetInitCtx())
}
