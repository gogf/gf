package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/os/gtimer"
	"time"
)

func GetList() {
START:
	for {
		res, err := g.Redis().DoVar("RPOP", "mill")
		if err != nil {
			glog.Debug("Rpop:", err)
			break
		}
		glog.Debug(res)
		if res.IsEmpty() {
			glog.Debug("nil")
			continue START
		}
		interval := 50 * time.Second
		gtimer.AddOnce(interval, func() {
			glog.Debug("end------:", res, gtime.Now().Format("Y-m-d H:i:s"))
		})
	}
}

func main() {
	g.Redis().SetMaxActive(2)
	//g.Redis().SetMaxIdle(100)
	GetList()
}
