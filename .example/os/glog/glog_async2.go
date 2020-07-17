package main

import (
	"time"

	"github.com/jin502437344/gf/os/glog"
)

func main() {
	glog.SetAsync(true)
	for i := 0; i < 10; i++ {
		glog.Async().Print("async log", i)
	}
	time.Sleep(time.Second)
}
