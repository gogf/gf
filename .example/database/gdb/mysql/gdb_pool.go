package main

import (
	"time"

	"github.com/gogf/gf/frame/g"
)

func main() {
	db := g.DB()

	// 开启调试模式，以便于记录所有执行的SQL
	db.SetDebug(true)

	for {
		for i := 0; i < 10; i++ {
			go db.Table("user").All()
		}
		time.Sleep(time.Millisecond * 100)
	}

}
