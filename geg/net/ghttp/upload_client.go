package main

import (
    "fmt"
    "gitee.com/johng/gf/g/net/ghttp"
)

func main() {
    _, e := ghttp.Post("http://127.0.0.1:8199/upload", "name=john&age=18")
    fmt.Println(e)
}