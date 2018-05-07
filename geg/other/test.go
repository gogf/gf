package main

import (
    "fmt"
    "net/http"

    "github.com/tabalt/gracehttp"
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
    test()
}