package main

import (
    "fmt"
)

func main() {
    fmt.Println(len(make(chan int, 10)))
}