package main

import (
	"fmt"
	"time"

	"github.com/gogf/gf/v2/os/grpool"
)

func main() {
	p := grpool.New(1)
	for i := 0; i < 10; i++ {
		v := i
		p.Add(func() {
			fmt.Println(v)
			time.Sleep(3 * time.Second)
		})
	}
	time.Sleep(time.Minute)
}
