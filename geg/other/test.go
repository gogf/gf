package main

import (
    "fmt"
)

func main() {
    methodMap := make(map[string]bool)
    methodMap["t"] = true
    fmt.Println(methodMap["t"])
}