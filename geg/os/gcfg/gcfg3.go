package main

import (
    "fmt"
    "gitee.com/johng/gf/g"
)

// 使用GetVar获取动态变量
func main() {
    fmt.Println(g.Config().GetVar("memcache.0").String())
}

