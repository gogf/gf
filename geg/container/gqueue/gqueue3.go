package main

import (
	"fmt"
	"github.com/gogf/gf/g/container/gqueue"
	"github.com/gogf/gf/g/os/gtime"
	"github.com/gogf/gf/g/os/gtimer"
	"time"
)

func main() {
	queue := gqueue.New()
	// 数据生产者，每隔1秒往队列写数据
	gtimer.SetInterval(time.Second, func() {
		queue.Push(gtime.Now().String())
	})

	// 消费者，不停读取队列数据并输出到终端
	for {
		select {
		case v := <-queue.C:
			if v != nil {
				fmt.Println(v)
			} else {
				return
			}
		}
	}
}
