package main

import (
    "fmt"
)

func main() {
    fmt.Println(string([]byte{'\t', '\n', '\v', '\f', '\r', ' ', 0x85, 0xA0}))
}