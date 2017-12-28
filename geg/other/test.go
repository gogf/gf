package main

import (
    "gitee.com/johng/gf/g/util/gregx"
    "fmt"
)


func main() {
    s, _ := gregx.Replace(`\w`, []byte("/user/list/page/2"), []byte("-"))
    fmt.Println(string(s))
}