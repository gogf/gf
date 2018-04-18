package main

import (
    "fmt"
    "gitee.com/johng/gf/g/database/gredis"
)

func main() {
    redis := gredis.New("127.0.0.1:6379", 1)
    fmt.Println(redis.Do("SET", "k1", "v1"))
    fmt.Println(redis.Do("SET", "k2", "v3"))
    fmt.Println(redis.Do("GET", "k2"))
    fmt.Println(redis.Stats())
}

