package main

import (
    "fmt"
    "gitee.com/johng/gf/g/net/reuseport"
    "gitee.com/johng/gf/g/os/gproc"
    "net/http"
)

func main() {
    listener, err := reuseport.Listen("tcp", ":8881")
    if err != nil {
        panic(err)
    }
    defer listener.Close()

    server := &http.Server{}
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "%d\n", gproc.Pid())
    })

    panic(server.Serve(listener))
}