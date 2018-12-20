package main

import (
    "fmt"
    "gitee.com/johng/gf/g/util/grand"
)

func main() {
    for i := 0; i < 100; i++ {
        fmt.Println(grand.Rand(0, 99999))
    }
}
