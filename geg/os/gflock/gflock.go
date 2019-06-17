package main

import (
	"github.com/gogf/gf/g/os/gflock"
	"github.com/gogf/gf/g/os/glog"
	"github.com/gogf/gf/g/os/gproc"
	"time"
)

func main() {
	l := gflock.New("demo.lock")
	l.Lock()
	glog.Printf("locked by pid: %d", gproc.Pid())
	time.Sleep(3 * time.Second)
	l.UnLock()
	glog.Printf("unlocked by pid: %d", gproc.Pid())
}
