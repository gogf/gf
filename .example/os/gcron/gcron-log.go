package main

import (
	"time"

	"github.com/gogf/gf/v2/os/gcron"
	"github.com/gogf/gf/v2/os/glog"
)

func main() {
	gcron.SetLogLevel(glog.LEVEL_ALL)
	gcron.Add("* * * * * ?", func() {
		glog.Print("test")
	})
	time.Sleep(3 * time.Second)
}
