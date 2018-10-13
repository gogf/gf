package main

import (
    "fmt"
    "gitee.com/johng/gf/g/util/grand"
    "os"
)

func main() {
    fmt.Println(uint(-1))
    os.Exit(0)
    for i := 0; i < 10; i++ {
        fmt.Println(grand.Rand(100, 200))
    }
}
