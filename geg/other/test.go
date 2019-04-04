package main

import (
	"fmt"
	"github.com/gogf/gf/g/database/gredis"
	"github.com/gogf/gf/g/util/gconv"
)

var (
	config = gredis.Config{
		Host : "127.0.0.1",
		Port : 6379,
		Db   : 1,
	}
)

func main() {
	group := "test"
	gredis.SetConfig(config, group)

	redis := gredis.Instance(group)
	defer redis.Close()

	_, err := redis.Do("SET", "k", "v")
	if err != nil {
		panic(err)
	}

	r, err := redis.Do("GET", "k")
	if err != nil {
		panic(err)
	}
	fmt.Println(gconv.String(r))

}