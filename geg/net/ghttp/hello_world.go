package main

import "gitee.com/johng/gf/g/net/ghttp"

func Hello(s *ghttp.Server, r *ghttp.ClientRequest, w *ghttp.ServerResponse) {
    w.WriteString("Hello World!")
}
func main() {
    s := ghttp.GetServer()
    s.BindHandler("/", Hello)
    s.Run()
}
