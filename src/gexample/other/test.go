package main

import (
    "fmt"
)

func main() {
    if (false) {
        goto T1
    } else {
        goto T2
    }
    T1:
    fmt.Println(1)
    T2:
    fmt.Println(2)



}