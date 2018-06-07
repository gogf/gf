package main

import (
    "fmt"
    "gitee.com/johng/gf/g/util/gregx"
)

func main() {
    a , e := gregx.MatchString(`<\?xml.*encoding\s*=\s*['|"](.*?)['|"].*\?>`, `<?xml version= '1.0' encoding = "utf-8" ?>`)
    fmt.Println(e)
    for k, v := range a {
        fmt.Printf("%d:%v\n", k, v)
    }
}