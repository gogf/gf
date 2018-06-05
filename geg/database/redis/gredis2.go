package main

import (
    "fmt"
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/util/gconv"
)

func main() {
    redis := g.Redis()
    defer redis.Close()
    redis.Do("SET", "k1", "v1")
    redis.Do("SET", "k2", "v2")
    v1, _ := redis.Do("GET", "k1")
    v2, _ := redis.Do("GET", "k1")
    fmt.Println(gconv.String(v1))
    fmt.Println(gconv.String(v2))
}

