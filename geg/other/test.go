package main

import (
    "fmt"
)

func test() {
    defer fmt.Println(1)
    fmt.Println(2)
}

func main() {
    test()
}
