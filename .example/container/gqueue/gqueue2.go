package main

import (
	"fmt"
	"time"

	"github.com/jin502437344/gf/container/gqueue"
	"github.com/jin502437344/gf/os/gtimer"
)

func main() {
	q := gqueue.New()
	// 数据生产者，每隔1秒往队列写数据
	gtimer.SetInterval(time.Second, func() {
		for i := 0; i < 10; i++ {
			q.Push(i)
		}
	})

	// 消费者，不停读取队列数据并输出到终端
	for {
		if v := q.Pop(); v != nil {
			fmt.Println(" Pop:", v)
		} else {
			break
		}
	}
}
