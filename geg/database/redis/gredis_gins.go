package main

import (
    "fmt"
    "gitee.com/johng/gf/g/frame/gins"
    "gitee.com/johng/gf/g/util/gconv"
)

func main() {
    gins.Config().SetPath("/home/john/Workspace/Go/GOPATH/src/gitee.com/johng/gf/geg/frame")
    
    redis := gins.Redis("cache")
    redis.Do("SET", "k", "v")
    v, _ := redis.Do("GET", "k")
    fmt.Println(gconv.String(v))
}

