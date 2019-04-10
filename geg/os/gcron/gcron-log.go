package main

import (
	"github.com/gogf/gf/g/os/gcron"
	"github.com/gogf/gf/g/os/glog"
	"time"
)


func main() {
	gcron.SetLogLevel(glog.LEVEL_ALL)
	gcron.Add("* * * * * ?", func() {
		glog.Println("test")
	})
	time.Sleep(3 * time.Second)
}
