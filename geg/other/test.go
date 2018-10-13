package main

import (
    "fmt"
    "gitee.com/johng/gf/g/util/grand"
)

func main() {
    for i := 0; i < 10; i++ {
        //fmt.Println(grand.RandStr(3))
        fmt.Println(grand.Rand(100, 200))
    }
}
