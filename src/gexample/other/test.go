package main

import (
    "g/net/ghttp"
    "fmt"
)


func main() {
    c := ghttp.NewClient()
    r := c.Get("http://baidu.com")
    fmt.Println(r.Close)
    r.ReadAll()
    fmt.Println(r.Close)


}