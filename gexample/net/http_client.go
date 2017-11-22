package main

import (
    "fmt"
    "g/net/ghttp"
)


func main() {
    c := ghttp.NewClient()
    r := c.Get("http://192.168.2.124")

    fmt.Println(r.StatusCode)
}
