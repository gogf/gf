package main

import (
    "fmt"
    "gitee.com/johng/gf/g/encoding/gparser"
)

func main() {
    f := gparser.New()
    f.Set("name", "john")
    f.Set("name", "john2")
    c, e := f.ToJson()
    fmt.Println(e)
    fmt.Println(string(c))
}