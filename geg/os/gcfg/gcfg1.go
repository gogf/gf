package main

import (
    "fmt"
    "gitee.com/johng/gf/g/os/gcfg"
)

func main() {
    c              := gcfg.New("/home/john/Workspace/Go/GOPATH/src/gitee.com/johng/gf/geg/os/gcfg")
    redisConfig    := c.GetArray("redis-cache", "redis.toml")
    memConfig      := c.GetArray("",    "memcache.yml")
    fmt.Println(redisConfig)
    fmt.Println(memConfig)
}

