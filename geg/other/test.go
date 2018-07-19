package main

import (
    "fmt"
    "time"
    "errors"
)

func test() (err error) {
    defer func() {
        fmt.Println(err)
    }()
    time.Sleep(time.Second)
    return errors.New("111")
}
func main() {
    test()
}