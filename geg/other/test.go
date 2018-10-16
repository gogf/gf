package main

import (
    "fmt"
    "gitee.com/johng/gf/g/encoding/gurl"
    "gitee.com/johng/gf/g/encoding/ghtml"
)

func main() {
    fmt.Println(gurl.Decode("<"))
    fmt.Println(ghtml.SpecialChars("<"))
}
