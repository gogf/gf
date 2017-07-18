package main

import (
    "net/http"
    "io"
    "g/net/ghttp"
)

func HelloServer1(w http.ResponseWriter, r *http.Request) {
    io.WriteString(w, "hello1!\n")
}
func HelloServer2(w http.ResponseWriter, r *http.Request) {
    io.WriteString(w, "hello2\n")
}
func main() {
    s := ghttp.New()
    s.SetAddr(":8199")
    s.SetIndexFolder(true)
    s.SetServerRoot("/home/john/Workspace/")
    s.BindHandleByMap(ghttp.HandlerMap {
        "/h":  HelloServer1,
        "/h1": HelloServer1,
        "/h2": HelloServer1,
        "/h3": HelloServer1,
    })
    s.BindHandle("/hello1", HelloServer1)
    s.BindHandle("/hello2", HelloServer2)
    s.Run()
}
