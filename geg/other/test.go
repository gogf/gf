package main

import (
    "fmt"
)
type T struct {

}

func main() {
    var i interface{}
    i = "s"
    i = make([]string, 100)

    if i == "s" {
        fmt.Println(1)
    }
}