package main

import (
	"fmt"
	"github.com/gogf/gf/g"
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
