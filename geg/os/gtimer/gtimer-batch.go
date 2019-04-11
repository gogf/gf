package main

import (
	"fmt"
	"github.com/gogf/gf/g/os/gtimer"
	"time"
)

func main() {
	for i := 0; i < 100000; i++ {
		gtimer.Add(time.Second, func() {

		})
	}
	fmt.Println("start")
	time.Sleep(48*time.Hour)
}
