// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package main

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/os/gctx"

	"github.com/gogf/gf/contrib/registry/nacos/v2"
)

func main() {
	gsvc.SetRegistry(nacos.New(`127.0.0.1:8848`))
	ctx := gctx.New()
	res := g.Client().GetContent(ctx, `http://hello.svc/`)
	g.Log().Info(ctx, res)
}
