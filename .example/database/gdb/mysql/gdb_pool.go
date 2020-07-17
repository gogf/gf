package main

import (
	"time"

	"github.com/jin502437344/gf/frame/g"
)

func main() {
	db := g.DB()
	db.SetMaxIdleConnCount(10)
	db.SetMaxOpenConnCount(10)
	db.SetMaxConnLifetime(time.Minute)

	// 开启调试模式，以便于记录所有执行的SQL
	db.SetDebug(true)

	for {
		for i := 0; i < 10; i++ {
			go db.Table("user").All()
		}
		time.Sleep(time.Second)
	}

}
