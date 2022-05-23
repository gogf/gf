package main

import (
	"github.com/polarismesh/polaris-go/pkg/config"

	"github.com/gogf/gf/contrib/registry/polaris/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/net/gsvc"
)

func main() {
	conf := config.NewDefaultConfiguration([]string{"192.168.100.222:8091"})

	// TTL egt 2*time.Second
	gsvc.SetRegistry(polaris.NewWithConfig(conf, polaris.WithTTL(10)))

	s := g.Server(`hello.svc`)
	s.BindHandler("/", func(r *ghttp.Request) {
		g.Log().Info(r.Context(), `request received`)
		r.Response.Write(`Hello world`)
	})
	s.Run()
}
