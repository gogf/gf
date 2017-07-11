package main

import (
    //"time"
    "net/http"
    "io"
    "time"
    "g/ghttp"
    "g"
)

func HelloServer(w http.ResponseWriter, req *http.Request) {
    io.WriteString(w, "hello, world!\n")
}
func HelloServer2(w http.ResponseWriter, req *http.Request) {
    io.WriteString(w, "hello123\n")
}
func main() {
    //s := http.Server{
    //    Addr          : ":8889",
    //    ReadTimeout   : 10 * time.Second,
    //    WriteTimeout  : 10 * time.Second,
    //}
    //g.Http.Server.NewByConfig(s)
    //http.HandleFunc("/hello2", HelloServer2)
    //g.Http.Server.BindHandle("/hello2", HelloServer2)
	//
    //g.Http.Server.BindHandleByMap(map[string]http.HandlerFunc {
    //    "/h":  HelloServer,
    //    "/h1": HelloServer,
    //    "/h2": HelloServer,
    //    "/h3": HelloServer,
    //})
    //dir := "/home/john/Workspace/Go/gf/src/gfexample/static"
    //http.Handle("/static", http.StripPrefix("/static/plugin/agile-lite", http.FileServer(http.Dir(dir))))
    //g.Http.Server.Start(":8199")
    //s := g.Http.Server.NewByAddr(":8199")
    //s.BindHandle("/hello", HelloServer)
    ghttp.Server.SetSetting(g.GHttpServerSetting {
        Addr           : ":8199",
        ReadTimeout    : 10 * time.Second,
        WriteTimeout   : 10 * time.Second,
        IdleTimeout    : 10 * time.Second,
        MaxHeaderBytes : 1024,
        ServerAgent    : "gf",
        ServerRoot     : "/home/john/Workspace/Go/gf/src/gfexample/static/",
    })

    g.Http.Server.BindHandle("/hello", HelloServer2)

    g.Http.Server.Run()
}
