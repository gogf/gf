package main

import (
	"fmt"

	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/util/gconv"
)

// 使用原生gredis.New操作redis，但是注意需要自己调用Close方法关闭redis链接池
func main() {
	config := &gredis.Config{
		Host: "127.0.0.1",
		Port: 6379,
	}
	redis := gredis.New(config)
	defer redis.Close()
	redis.Do("SET", "k", "v")
	v, _ := redis.Do("GET", "k")
	fmt.Println(gconv.String(v))
}
