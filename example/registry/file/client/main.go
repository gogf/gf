package main

import (
	"fmt"
	"time"

	"github.com/gogf/gf/contrib/registry/file/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
)

func main() {
	gsvc.SetRegistry(file.New(gfile.Temp("gsvc")))

	client := g.Client()
	for i := 0; i < 100; i++ {
		res, err := client.Get(gctx.New(), `http://hello.svc/`)
		if err != nil {
			panic(err)
		}
		fmt.Println(res.ReadAllString())
		res.Close()
		time.Sleep(time.Second)
	}
}
