package main

import (
	"fmt"
	"github.com/gogf/gf/g/os/grpool"
	"github.com/gogf/gf/g/os/gtimer"
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
	gtimer.SetInterval(time.Second, func() {
		fmt.Println("worker:", pool.Size())
		fmt.Println("  jobs:", pool.Jobs())
		fmt.Println()
		gtimer.Exit()
	})

	select {}
}
