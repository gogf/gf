package main

import (
    "fmt"
    "gitee.com/johng/gf/g/string/gstr"
)

func main() {
    fmt.Println(gstr.PosI("abcdEfg", "eF", 0))
}