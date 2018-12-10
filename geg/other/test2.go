package main

import (
    "fmt"
)


func Map() map[int]int {
    return nil
}


func main() {
    v := (interface{})(nil)
    fmt.Println(v == nil)
    if v = Map(); v != nil {
        fmt.Println(v == nil)
    }
    fmt.Println(Map() == nil)
}