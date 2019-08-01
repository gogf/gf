package main

import (
	"time"

	"github.com/gogf/gf/os/gflock"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/os/gproc"
)

func main() {
	l := gflock.New("demo.lock")
	l.Lock()
	glog.Printf("locked by pid: %d", gproc.Pid())
	time.Sleep(10 * time.Second)
	l.UnLock()
	glog.Printf("unlocked by pid: %d", gproc.Pid())
}
