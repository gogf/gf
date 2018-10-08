package main

import (
    "fmt"
    "gitee.com/johng/gf/g"
)

func main() {
    fmt.Println(g.Config().Get("redis"))

    type RedisConfig struct {
        Disk  string
        Cache string
    }

    redisCfg := new(RedisConfig)
    fmt.Println(g.Config().GetToStruct("redis", redisCfg))
    fmt.Println(redisCfg)
}



