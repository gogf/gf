// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package main

import (
	"context"

	"github.com/polarismesh/polaris-go/api"
	"github.com/polarismesh/polaris-go/pkg/config"

	"github.com/gogf/gf/contrib/registry/polaris/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/net/gsvc"
)

func main() {
	conf := config.NewDefaultConfiguration([]string{"183.47.111.80:8091"})
	conf.Consumer.LocalCache.SetPersistDir("/tmp/polaris/backup")
	if err := api.SetLoggersDir("/tmp/polaris/log"); err != nil {
		g.Log().Fatal(context.Background(), err)
	}

	// TTL egt 2*time.Second
	gsvc.SetRegistry(polaris.NewWithConfig(conf, polaris.WithTTL(10)))

	s := g.Server(`hello-world.svc`)
	s.BindHandler("/", func(r *ghttp.Request) {
		g.Log().Info(r.Context(), `request received`)
		r.Response.Write(`Hello world`)
	})
	s.Run()
}
