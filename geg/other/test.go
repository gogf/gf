package main

import "fmt"

func test() (err interface{}) {
    defer func() {
        err = recover()
    }()
    panic(1)
    return
}

func main() {

    switch err := test(); err {
    default:
        fmt.Println(err)
    }
}
