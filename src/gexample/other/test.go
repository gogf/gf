package main

import (
    "fmt"
)

type T struct {
    I int
    J string
}

func main() {
    m := make(map[string]string)
    fmt.Printf("%T", m)
}