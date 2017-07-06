package main

import (
    "fmt"
)

func test(data interface{}) {
    switch data.(type) {
    case string:
        fmt.Println("string")
    case map[string]string:
        fmt.Println("map[string]string")
    case []interface{}:
        fmt.Println("[]interface{}")
    default:
        fmt.Println("default")
    }
}
func main() {
    test(map[string]string {"k": "v"})
}