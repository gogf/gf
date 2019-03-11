package main

import (
    "fmt"
    "github.com/gogf/gf/g/util/gconv"
)

func main() {
    t := gconv.GTime("2010-10-10 00:00:01")
    fmt.Println(t.String())
}