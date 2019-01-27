package main

import (
    "fmt"
    "gitee.com/johng/gf/g/net/greuseport"
    "net/http"
    "os"
)

// 创建**两个**进程，并通过HTTP访问，查看返回结果。
func main() {
    listener, err := greuseport.Listen("tcp", ":8881")
    if err != nil {
        panic(err)
    }
    defer listener.Close()

    server := &http.Server{}
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "gid: %d, pid: %d\n", os.Getgid(), os.Getpid())
    })

    panic(server.Serve(listener))
}
