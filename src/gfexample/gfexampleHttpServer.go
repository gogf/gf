package main

import (
    "gf"
    //"time"
    "net/http"
    "io"
    "time"
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
    //gf.Http.Server.NewByConfig(s)
    //http.HandleFunc("/hello2", HelloServer2)
    //gf.Http.Server.BindHandle("/hello2", HelloServer2)
	//
    //gf.Http.Server.BindHandleByMap(map[string]http.HandlerFunc {
    //    "/h":  HelloServer,
    //    "/h1": HelloServer,
    //    "/h2": HelloServer,
    //    "/h3": HelloServer,
    //})
    //dir := "/home/john/Workspace/Go/gf/src/gfexample/static"
    //http.Handle("/static", http.StripPrefix("/static/plugin/agile-lite", http.FileServer(http.Dir(dir))))
    //gf.Http.Server.Start(":8199")
    //s := gf.Http.Server.NewByAddr(":8199")
    //s.BindHandle("/hello", HelloServer)
    gf.Http.Server.SetSetting(gf.GstHttpServerSetting {
        Addr           : ":8199",
        ReadTimeout    : 10 * time.Second,
        WriteTimeout   : 10 * time.Second,
        IdleTimeout    : 10 * time.Second,
        MaxHeaderBytes : 1024,
        ServerAgent    : "gf",
        ServerRoot     : "/home/john/Workspace/Go/gf/src/gfexample/static/",
    })

    gf.Http.Server.BindHandle("/hello", HelloServer2)

    gf.Http.Server.Run()
}
