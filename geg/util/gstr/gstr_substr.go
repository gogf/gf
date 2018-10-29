package main

import (
    "fmt"
    "gitee.com/johng/gf/g/util/gstr"
)

func main() {
    fmt.Println(gstr.SubStr("我是中国人", 2))
    fmt.Println(gstr.SubStr("我是中国人", 2, 2))
}
