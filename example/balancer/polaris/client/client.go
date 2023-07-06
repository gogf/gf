// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/polarismesh/polaris-go/api"
	"github.com/polarismesh/polaris-go/pkg/config"

	"github.com/gogf/gf/contrib/registry/polaris/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gsel"
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/os/gctx"
)

func main() {
	conf := config.NewDefaultConfiguration([]string{"127.0.0.1:8091"})
	conf.Consumer.LocalCache.SetPersistDir(os.TempDir() + "/polaris/backup")
	if err := api.SetLoggersDir(os.TempDir() + "/polaris/log"); err != nil {
		g.Log().Fatal(context.Background(), err)
	}

	gsvc.SetRegistry(polaris.NewWithConfig(conf, polaris.WithTTL(10)))
	gsel.SetBuilder(gsel.NewBuilderRoundRobin())

	for i := 0; i < 100; i++ {
		res, err := g.Client().Get(gctx.New(), `http://hello-world.svc/`)
		if err != nil {
			panic(err)
		}
		fmt.Println(res.ReadAllString(), " id: ", i, " time: ", time.Now().Format("2006-01-02 15:04:05"), " code: ", res.StatusCode)
		res.Close()
		time.Sleep(time.Second)
	}
}
