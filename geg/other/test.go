package main

import (
    "fmt"
    "runtime"
)

func main() {
    fmt.Println(runtime.GOMAXPROCS(0))
}