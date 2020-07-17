package main

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"time"
)

func main() {
	db := g.DB()
	db.SetDebug(true)
	for {
		r, err := db.Table("user").All()
		fmt.Println(err)
		fmt.Println(r)
		time.Sleep(time.Second * 10)
	}
}
