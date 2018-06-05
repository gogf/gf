package main

import (
    "fmt"
    "gitee.com/johng/gf/g/util/gregx"
)

func main() {
    a , e := gregx.MatchString(`(.+):(\d+),{0,1}(\d*),{0,1}(.*)`, "127.0.0.1:12333")
    fmt.Println(e)
    for k, v := range a {
        fmt.Printf("%d:%v\n", k, v)
    }
}