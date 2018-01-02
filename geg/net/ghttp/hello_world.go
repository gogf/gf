package main

import "gitee.com/johng/gf/g/net/ghttp"

func Hello(r *ghttp.Request) {
    r.Response.WriteString("Hello World!")
}
func main() {
    s := ghttp.GetServer()
    s.BindHandler("/", Hello)
    s.Run()
}
