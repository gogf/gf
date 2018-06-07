package main

import (
    "fmt"
    "gitee.com/johng/gf/g/net/ghttp"
)

func main() {
    c    := ghttp.NewClient()
    r, _ := c.Get("http://baidu.com")
    fmt.Println(r.StatusCode)
}
