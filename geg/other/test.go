package main

import (
    "fmt"
)

func test(data []byte) {
    data = append(data, byte('A'))
    fmt.Println(data)
}

func main() {
    a := []byte("1")
    test(a)
    fmt.Println(a)
}
