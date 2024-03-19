// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/wangyougui/gf.

package main

import (
	"github.com/wangyougui/gf/contrib/registry/etcd/v2"
	"github.com/wangyougui/gf/v2/frame/g"
	"github.com/wangyougui/gf/v2/net/ghttp"
	"github.com/wangyougui/gf/v2/net/gsvc"
)

func main() {
	gsvc.SetRegistry(etcd.New(`127.0.0.1:2379`))

	s := g.Server(`hello.svc`)
	s.BindHandler("/", func(r *ghttp.Request) {
		g.Log().Info(r.Context(), `request received`)
		r.Response.Write(`Hello world`)
	})
	s.Run()
}
