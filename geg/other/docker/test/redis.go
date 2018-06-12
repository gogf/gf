package main

import (
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/net/ghttp"
    "gitee.com/johng/gf/g/database/gredis"
)

func main() {
    s := g.Server()
    s.BindHandler("/", func(r *ghttp.Request) {
        redis := gredis.New(gredis.Config{
            Host: "redis-service",
            Port: 9999,
        })
        defer redis.Close()

        v, err := redis.Do("GET", "k")
        r.Response.Writeln(v)
        r.Response.Writeln(err)
        v, err  = redis.Do("SET", "k", "v")
        r.Response.Writeln(v)
        r.Response.Writeln(err)

    })
    s.Run()
}