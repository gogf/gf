package main

import (
    "gitee.com/johng/gf/g/net/ghttp"
)

func main() {
    s := ghttp.GetServer()
    s.BindHandler("/template2", func(r *ghttp.Request){
        //panic("123")
    })
    s.SetAccessLogEnabled(true)
    s.SetPort(8199)
    s.Run()
}