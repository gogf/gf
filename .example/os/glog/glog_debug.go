package main

import (
	"github.com/gogf/gf/frame/g"
	"time"

	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/os/gtimer"
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
