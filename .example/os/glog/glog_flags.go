package main

import (
	"github.com/jin502437344/gf/os/glog"
)

func main() {
	l := glog.New()
	l.SetFlags(glog.F_TIME_TIME | glog.F_FILE_SHORT)
	l.Println("time and short line number")
	l.SetFlags(glog.F_TIME_MILLI | glog.F_FILE_LONG)
	l.Println("time with millisecond and long line number")
	l.SetFlags(glog.F_TIME_STD | glog.F_FILE_LONG)
	l.Println("standard time format and long line number")
}
