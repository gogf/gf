package main

import (
	"fmt"
	"time"

	"github.com/gogf/gf/v2/os/gtimer"
)

func main() {
	now := time.Now()
	interval := 510 * time.Millisecond
	gtimer.Add(interval, func() {
		fmt.Println(time.Now(), time.Duration(time.Now().UnixNano()-now.UnixNano()))
		now = time.Now()
	})
	time.Sleep(time.Hour)
	select {}
}
