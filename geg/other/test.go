package main

import (
    "fmt"
    "gitee.com/johng/gf/g/text/gstr"
)

func main() {
    fmt.Println(gstr.TrimLeftStr("gogo我爱gogo", "go"))
    fmt.Println(gstr.TrimRightStr("gogo我爱gogo", "go"))
}
