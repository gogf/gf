package main

import (
    "fmt"
    "gitee.com/johng/gf/g/encoding/gurl"
)


type T struct {
    name string
}

func (t *T)Test() {
    fmt.Println(t.name)
}

func main() {
    fmt.Println(gurl.Encode("@123"))
}