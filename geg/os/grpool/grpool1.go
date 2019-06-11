package main

import (
<<<<<<< HEAD
    "time"
    "fmt"
    "gitee.com/johng/gf/g/os/gtime"
    "gitee.com/johng/gf/g/os/grpool"
)

func job() {
    time.Sleep(1*time.Second)
}

func main() {
    grpool.SetSize(10)
    for i := 0; i < 1000; i++ {
        grpool.Add(job)
    }
    gtime.SetInterval(2*time.Second, func() bool {
        fmt.Println("size:", grpool.Size())
        fmt.Println("jobs:", grpool.Jobs())
        return true
    })
    select {}
=======
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
>>>>>>> upstream/master
}
