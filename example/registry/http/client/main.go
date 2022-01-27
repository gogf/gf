package main

import (
	"fmt"

	"github.com/gogf/gf/contrib/registry/etcd/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/os/gctx"
)

func main() {
	gsvc.SetRegistry(etcd.New(`127.0.0.1:2379`))

	res, err := g.Client().Get(gctx.New(), `http://hello.svc:12345/`)
	if err != nil {
		panic(err)
	}
	defer res.Close()
	fmt.Println(res.ReadAllString())
}
