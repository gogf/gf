package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/gogf/gf/os/gtime"
)

func main() {
	start := gtime.TimestampMilli()
	wg := sync.WaitGroup{}
	for i := 0; i < 100000; i++ {
		wg.Add(1)
		go func() {
			time.Sleep(time.Second)
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println("time spent:", gtime.TimestampMilli()-start)
}
