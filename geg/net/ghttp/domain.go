package main

import "gitee.com/johng/gf/g/net/ghttp"

func Hello1(s *ghttp.Server, r *ghttp.ClientRequest, w *ghttp.ServerResponse) {
    w.WriteString("Hello World1!")
}

func Hello2(s *ghttp.Server, r *ghttp.ClientRequest, w *ghttp.ServerResponse) {
    w.WriteString("Hello World2!")
}

func main() {
    s := ghttp.GetServer()
    s.Domain("127.0.0.1").BindHandler("/", Hello1)
    s.Domain("localhost").BindHandler("/", Hello2)
    s.Run()
}
