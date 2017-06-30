package main

import (
    "gf"
    //"time"
    "net/http"
    "io"
)

func init() {
    // gf.Http.Server.BindHandle("/hello2", HelloServer2)
    http.HandleFunc("/hello2", HelloServer2)
}

func HelloServer(w http.ResponseWriter, req *http.Request) {
    io.WriteString(w, "hello, world!\n")
}
func HelloServer2(w http.ResponseWriter, req *http.Request) {
    io.WriteString(w, "hello, world2!\n")
}
func main() {
    //s := http.Server{
    //    Addr          : ":8889",
    //    ReadTimeout   : 10 * time.Second,
    //    WriteTimeout  : 10 * time.Second,
    //}
    //gf.Http.Server.NewByConfig(s)

    gf.Http.Server.New(":8199", nil)
}
