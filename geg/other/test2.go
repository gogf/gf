package main

import (
    "fmt"
)

func main() {
    //var v interface{}
    m := map[string]int {
        "age" : 18,
    }
    //v  = m
    p := &m
    (*p)["age"] = 16
    //fmt.Println(v)
    fmt.Println(m)
}