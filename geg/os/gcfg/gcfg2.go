package main

import (
    "fmt"
    "gitee.com/johng/gf/g"
)

// 使用默认的config.toml配置文件读取配置
func main() {
    c := g.Config()
    fmt.Println(c.GetArray("memcache"))
}

