package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gogf/gf/contrib/registry/polaris/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/polarismesh/polaris-go/api"
	"github.com/polarismesh/polaris-go/pkg/config"
)

func main() {
	conf := config.NewDefaultConfiguration([]string{"192.168.100.222:8091"})
	conf.Consumer.LocalCache.SetPersistDir("/tmp/polaris/backup")
	if err := api.SetLoggersDir("/tmp/polaris/log"); err != nil {
		g.Log().Fatal(context.Background(), err)
	}

	gsvc.SetRegistry(polaris.NewWithConfig(conf, polaris.WithTTL(10)))

	for i := 0; i < 100; i++ {
		res, err := g.Client().Get(gctx.New(), `http://hello.svc/`)
		if err != nil {
			panic(err)
		}
		fmt.Println(res.ReadAllString())
		res.Close()
		time.Sleep(time.Second)
	}
}
