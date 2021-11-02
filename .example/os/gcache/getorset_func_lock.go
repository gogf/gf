package main

import (
	"fmt"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gctx"
	"time"
)

func main() {
	var (
		ch    = make(chan struct{}, 0)
		ctx   = gctx.New()
		key   = `key`
		value = `value`
	)
	for i := 0; i < 10; i++ {
		go func(index int) {
			<-ch
			_, _ = gcache.Ctx(ctx).GetOrSetFuncLock(key, func() (interface{}, error) {
				fmt.Println(index, "entered")
				return value, nil
			}, 0)
		}(i)
	}
	close(ch)
	time.Sleep(time.Second)
}
