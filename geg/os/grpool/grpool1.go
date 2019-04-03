package main

import (
	"fmt"
	"github.com/gogf/gf/g/os/grpool"
	"github.com/gogf/gf/g/os/gtime"
	"time"
)

func job() {
	time.Sleep(1 * time.Second)
}

func main() {
	pool := grpool.New(100)
	for i := 0; i < 1000; i++ {
		pool.Add(job)
	}
	fmt.Println("worker:", pool.Size())
	fmt.Println("  jobs:", pool.Jobs())
	gtime.SetInterval(time.Second, func() bool {
		fmt.Println("worker:", pool.Size())
		fmt.Println("  jobs:", pool.Jobs())
		fmt.Println()
		return true
	})

	select {}
}
