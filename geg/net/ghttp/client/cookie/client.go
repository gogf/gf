package main

import (
    "fmt"
    "gitee.com/johng/gf/g/os/glog"
    "gitee.com/johng/gf/g/net/ghttp"
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
