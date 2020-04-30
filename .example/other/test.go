package main

import (
	"fmt"
	"github.com/gogf/gf/os/gtimer"
	"time"
)

func main() {
	tr := gtimer.New(100, 1*time.Second, 10)
	tr.Add(1*time.Second, func() { fmt.Println("hello") })
	select {}
}
