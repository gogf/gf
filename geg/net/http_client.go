package main

import (
    "fmt"
    "gitee.com/johng/gf/g/net/ghttp"
)


func main() {
    c    := ghttp.NewClient()
    r, _ := c.Get("http://192.168.2.124")
    fmt.Println(r.StatusCode)
}
