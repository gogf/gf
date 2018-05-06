package main

import (
    "fmt"
    "net/http"

    "github.com/tabalt/gracehttp"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "hello world")
    })

    err := gracehttp.ListenAndServe(":8888", nil)
    if err != nil {
        fmt.Println(err)
    }
}