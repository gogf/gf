package main

import (
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/util/gutil"
)

func main() {
	gutil.Dump(gfile.ScanDir("/Users/john/Documents", "*.*"))
	gutil.Dump(gfile.ScanDir("/home/john/temp/newproject", "*", true))
}
