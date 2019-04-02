package main

import (
	"fmt"
	"github.com/gogf/gf/g/net/ghttp"
	"github.com/gogf/gf/g/os/glog"
)

func main() {
	c := ghttp.NewClient()
	c.SetHeader("Cookie", "name=john; score=100")
	if r, e := c.Get("http://127.0.0.1:8199/"); e != nil {
		glog.Error(e)
	} else {
		fmt.Println(string(r.ReadAll()))
	}
}
