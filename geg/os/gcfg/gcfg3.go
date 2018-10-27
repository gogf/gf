package main

import (
    "fmt"
    "gitee.com/johng/gf/g"
)

func main() {
    fmt.Println(g.Config().GetVar("memcache.0").String())
}

