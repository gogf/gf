package main

import (
    "fmt"
    "gitee.com/johng/gf/g/net/ghttp"
)


func main() {
    r, _ := ghttp.Get("http://johng.cn")
    fmt.Println(string(r.ReadAll()))
}