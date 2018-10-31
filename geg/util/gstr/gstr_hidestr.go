package main

import (
    "fmt"
    "gitee.com/johng/gf/g/util/gstr"
)

func main() {
    fmt.Println(gstr.HideStr("热爱GF热爱生活", 20, "*"))
    fmt.Println(gstr.HideStr("热爱GF热爱生活", 50, "*"))
}
