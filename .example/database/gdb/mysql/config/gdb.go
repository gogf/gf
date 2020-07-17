package main

import (
	"github.com/gogf/gf/database/gdb"
	"sync"
	"time"
)

var db gdb.DB

func init() {
	gdb.AddDefaultConfigNode(gdb.ConfigNode{
		Host:             "127.0.0.1",
		Port:             "3306",
		User:             "root",
		Pass:             "12345678",
		Name:             "test",
		Type:             "mysql",
		Role:             "master",
		Charset:          "utf8",
		MaxOpenConnCount: 100,
	})
	db, _ = gdb.New()
}

func main() {
	wg := sync.WaitGroup{}
	for i := 0; i < 100000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(10 * time.Second)
			db.Table("user").Where("id=1").All()
		}()
	}
	wg.Wait()
}
