package main

import (
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/os/gctx"
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
	var (
		wg  = sync.WaitGroup{}
		ctx = gctx.New()
	)
	for i := 0; i < 100000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(10 * time.Second)
			db.Ctx(ctx).Model("user").Where("id=1").All()
		}()
	}
	wg.Wait()
}
