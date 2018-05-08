package main

import (
    "fmt"
    "net/http"

    "github.com/tabalt/gracehttp"
    "os"
    "gitee.com/johng/gf/g/os/gpm"
    "time"
)


func test() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "hello world")
    })
    s := gracehttp.NewServer(":8888", nil, gracehttp.DEFAULT_READ_TIMEOUT, gracehttp.DEFAULT_WRITE_TIMEOUT)
    err := s.ListenAndServe()
    if err != nil {
        fmt.Println(err)
    }
}

func main() {
    m    := gpm.New()
    args := os.Args
    args  = append(args, "--child=1")
    p    := m.NewProcess(args[0], args, nil)
    p.Run()
    time.Sleep(100*time.Second)
}
