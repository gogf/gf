package main

import (
	"time"

	"github.com/jin502437344/gf/os/gcron"
	"github.com/jin502437344/gf/os/glog"
)

func main() {
	gcron.SetLogLevel(glog.LEVEL_ALL)
	gcron.Add("* * * * * ?", func() {
		glog.Println("test")
	})
	time.Sleep(3 * time.Second)
}
