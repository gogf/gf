package main

import (
    "fmt"
    "gitee.com/johng/gf/g"
)

func main() {
    c              := g.Config()
    redisConfig    := c.GetArray("redis-cache", "redis.toml")
    memConfig      := c.GetArray("",    "memcache.yml")
    fmt.Println(redisConfig)
    fmt.Println(memConfig)
}

