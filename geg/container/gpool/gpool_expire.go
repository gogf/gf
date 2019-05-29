package main

import (
	"fmt"
	"github.com/gogf/gf/g/container/gpool"
	"time"
)

func main() {
	p := gpool.New(60000, nil, func(i interface{}) {
		fmt.Println("expired")
	})
	p.Put(1)
	time.Sleep(10000*time.Second)
	fmt.Println(p.Get())
}
