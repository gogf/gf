package main

import (
    "gitee.com/johng/gf/g/os/gtime"
    "fmt"
)

func main() {
    t := gtime.Now()
    err := t.ToZone("Asia/Aden")
    fmt.Println(err)
    fmt.Println(t.String())
    //fmt.Println(string([]byte{112,108,97,121,101,114,105,100}))
}
