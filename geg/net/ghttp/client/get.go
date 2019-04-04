package main

import (
	"fmt"
	"github.com/gogf/gf/g/net/ghttp"
)

func main() {
	c := ghttp.NewClient()
	r, _ := c.Get("http://baidu.com")
	fmt.Println(r.StatusCode)
}
