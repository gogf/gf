package main

import (
    "fmt"
    "gitee.com/johng/gf/g"
)

// 使用第二个参数指定读取的配置文件
func main() {
    c              := g.Config()
    redisConfig    := c.GetArray("redis-cache", "redis.toml")
    memConfig      := c.GetArray("",    "memcache.yml")
    fmt.Println(redisConfig)
    fmt.Println(memConfig)
}

