package main

import (
    "fmt"
    "gitee.com/johng/gf/g/util/gconv"
    "gitee.com/johng/gf/g/database/gredis"
)

func main() {
    redis := gredis.New("127.0.0.1:6379", 1)
    redis.Do("SET", "k1", "v1")
    redis.Do("SET", "k2", "v2")
    v1, _ := redis.Do("GET", "k1")
    v2, _ := redis.Do("GET", "k1")
    fmt.Println(gconv.String(v1))
    fmt.Println(gconv.String(v2))
}

