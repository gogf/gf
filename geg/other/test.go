package main

import (
    "fmt"
    "net/http"

    "github.com/tabalt/gracehttp"
    "os"
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
    fmt.Println(os.NewFile(11111, ""))
    fmt.Println(os.NewFile(111111111, ""))
    fmt.Println(os.NewFile(33333333333333, ""))
}
