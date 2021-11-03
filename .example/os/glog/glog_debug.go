package main

import (
	"github.com/gogf/gf/v2/frame/g"
	"time"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/os/gtimer"
)

func main() {
	gtimer.SetTimeout(3*time.Second, func() {
		g.Log().SetDebug(false)
	})
	for {
		g.Log().Debug(gtime.Datetime())
		time.Sleep(time.Second)
	}
}
