package main

import (
    "fmt"
)

type T struct{
    name string
}
func main() {
    var i interface{} = T{"john"}
    switch v := i.(type) {
    case T:
    default:
        fmt.Println(v)
    }
}