package main

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

func main() {
	ctx := context.Background()
	mylog := g.Log()
	for {
		mylog.Debug(ctx, "debug")
		time.Sleep(time.Second)
		mylog.Info(ctx, "info")
		time.Sleep(time.Second)
		mylog.Warning(ctx, "warning")
		time.Sleep(time.Second)
		mylog.Error(ctx, "error")
		time.Sleep(time.Second)
	}
}
