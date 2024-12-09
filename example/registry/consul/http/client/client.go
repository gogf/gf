// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package main

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/os/gctx"

	"github.com/gogf/gf/contrib/registry/consul/v2"
)

func main() {
	registry, err := consul.New(consul.WithAddress("127.0.0.1:8500"))
	if err != nil {
		g.Log().Fatal(context.Background(), err)
	}
	gsvc.SetRegistry(registry)

	ctx := gctx.New()
	res := g.Client().GetContent(ctx, `http://hello.svc/`)
	g.Log().Info(ctx, res)
}
