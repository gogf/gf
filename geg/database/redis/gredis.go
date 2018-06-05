package main

import (
    "fmt"
    "gitee.com/johng/gf/g/util/gconv"
    "gitee.com/johng/gf/g/database/gredis"
)

func main() {
    redis := gredis.New(gredis.Config{
        Host : "127.0.0.1",
        Port : 6379,
    })
    defer redis.Close()
    redis.Do("SET", "k", "v")
    v, _ := redis.Do("GET", "k")
    fmt.Println(gconv.String(v))
}

