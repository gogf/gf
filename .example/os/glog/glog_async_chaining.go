package main

import (
	"github.com/gogf/gf/frame/g"
	"time"
)

func main() {
	for i := 0; i < 10; i++ {
		g.Log().Async().Print("async log", i)
	}
	time.Sleep(time.Second)
}
